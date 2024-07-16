package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

type HardwareStatus struct {
	CPUUsage  []float64
	MemUsage  *mem.VirtualMemoryStat
	DiskUsage []disk.UsageStat
	NetIO     []net.IOCountersStat
	Uptime    uint64
	GPUStatus []GPUStatus
}

type GPUStatus struct {
	Name        string
	Utilization string
	MemoryUsage string
}

var previousStatus HardwareStatus

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	bgRed       = "\033[41m"
	bgGreen     = "\033[42m"
	bgYellow    = "\033[43m"
	bgBlue      = "\033[44m"
	bgPurple    = "\033[45m"
	bgCyan      = "\033[46m"
	bgWhite     = "\033[47m"
	cpuColor    = "\033[44;37m" // Blue background with white text
	memColor    = "\033[42;37m" // Green background with white text
	diskColor   = "\033[46;37m" // Cyan background with white text
	netColor    = "\033[43;37m" // Yellow background with white text
	sysColor    = "\033[45;37m" // Purple background with white text
	gpuColor    = "\033[41;37m" // Red background with white text
)

func main() {
	for {
		currentStatus := getHardwareStatus()
		displayStatus(currentStatus)
		previousStatus = currentStatus
		time.Sleep(500 * time.Millisecond)
	}
}

func getHardwareStatus() HardwareStatus {
	cpuUsage, _ := cpu.Percent(0, true)
	memUsage, _ := mem.VirtualMemory()
	partitions, _ := disk.Partitions(true)
	var diskUsages []disk.UsageStat
	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		diskUsages = append(diskUsages, *usage)
	}
	netIO, _ := net.IOCounters(true)
	uptime, _ := host.Uptime()
	gpuStatus := getGPUStatus()

	return HardwareStatus{
		CPUUsage:  cpuUsage,
		MemUsage:  memUsage,
		DiskUsage: diskUsages,
		NetIO:     netIO,
		Uptime:    uptime,
		GPUStatus: gpuStatus,
	}
}

func getGPUStatus() []GPUStatus {
	cmd := exec.Command("nvidia-smi", "--query-gpu=name,utilization.gpu,memory.used", "--format=csv,noheader,nounits")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to execute nvidia-smi: %v", err)
	}

	scanner := bufio.NewScanner(&out)
	var statuses []GPUStatus
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ", ")
		if len(fields) == 3 {
			statuses = append(statuses, GPUStatus{
				Name:        fields[0],
				Utilization: fields[1] + "%",
				MemoryUsage: fields[2] + " MiB",
			})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to parse nvidia-smi output: %v", err)
	}

	return statuses
}

func displayStatus(status HardwareStatus) {
	clearScreen()
	drawTableTopLine()
	drawTableHeader()

	for i, usage := range status.CPUUsage {
		color := getColor(previousStatus.CPUUsage != nil && previousStatus.CPUUsage[i] != usage, cpuColor)
		drawTableRow("CPU", fmt.Sprintf("Core %d Usage", i), fmt.Sprintf("%s%.2f%%%s", color, usage, colorReset))
	}

	color := getColor(previousStatus.MemUsage != nil && previousStatus.MemUsage.UsedPercent != status.MemUsage.UsedPercent, memColor)
	drawTableRow("Memory", "Usage", fmt.Sprintf("%s%.2f%%%s", color, status.MemUsage.UsedPercent, colorReset))

	for _, usage := range status.DiskUsage {
		var prevUsagePercent float64
		for _, prevUsage := range previousStatus.DiskUsage {
			if prevUsage.Path == usage.Path {
				prevUsagePercent = prevUsage.UsedPercent
			}
		}
		color := getColor(prevUsagePercent != usage.UsedPercent, diskColor)
		drawTableRow("Disk", fmt.Sprintf("Usage (%s)", usage.Path), fmt.Sprintf("%s%.2f%%%s", color, usage.UsedPercent, colorReset))
	}

	for _, io := range status.NetIO {
		var prevIO net.IOCountersStat
		for _, prev := range previousStatus.NetIO {
			if prev.Name == io.Name {
				prevIO = prev
			}
		}
		colorSent := getColor(prevIO.BytesSent != io.BytesSent, netColor)
		colorRecv := getColor(prevIO.BytesRecv != io.BytesRecv, netColor)
		drawTableRow("Network", fmt.Sprintf("Bytes Sent (%s)", io.Name), fmt.Sprintf("%s%d bytes%s", colorSent, io.BytesSent, colorReset))
		drawTableRow("Network", fmt.Sprintf("Bytes Received (%s)", io.Name), fmt.Sprintf("%s%d bytes%s", colorRecv, io.BytesRecv, colorReset))
	}

	color = getColor(previousStatus.Uptime != status.Uptime, sysColor)
	drawTableRow("System", "Uptime", fmt.Sprintf("%s%d seconds%s", color, status.Uptime, colorReset))

	for i, gpu := range status.GPUStatus {
		var prevGPU GPUStatus
		if previousStatus.GPUStatus != nil && len(previousStatus.GPUStatus) > i {
			prevGPU = previousStatus.GPUStatus[i]
		}
		colorUtilization := getColor(prevGPU.Utilization != gpu.Utilization, gpuColor)
		colorMemory := getColor(prevGPU.MemoryUsage != gpu.MemoryUsage, gpuColor)
		drawTableRow("GPU", fmt.Sprintf("Utilization (%s)", gpu.Name), fmt.Sprintf("%s%s%s", colorUtilization, gpu.Utilization, colorReset))
		drawTableRow("GPU", fmt.Sprintf("Memory Usage (%s)", gpu.Name), fmt.Sprintf("%s%s%s", colorMemory, gpu.MemoryUsage, colorReset))
	}

	drawTableBottomLine()
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[0;0H")
}

func drawTableTopLine() {
	fmt.Println("┌──────────────────────┬────────────────────────────────────────────────────┬────────────────────┐")
}

func drawTableHeader() {
	fmt.Println("│ Component            │ Metric                                             │ Value              │")
	drawTableMidLine()
}

func drawTableMidLine() {
	fmt.Println("├──────────────────────┼────────────────────────────────────────────────────┼────────────────────┤")
}

func drawTableBottomLine() {
	fmt.Println("└──────────────────────┴────────────────────────────────────────────────────┴────────────────────┘")
}

func drawTableRow(component, metric, value string) {
	fmt.Printf("│ %-20s │ %-50s │ %-30s │\n", component, metric, value)
}

func getColor(changed bool, baseColor string) string {
	if changed {
		return bgWhite + colorRed // White background with red text for updates
	}
	return baseColor // Base color for each component type
}
