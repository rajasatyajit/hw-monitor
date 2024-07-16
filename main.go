package main

/*
#cgo CFLAGS: -IC:\\PROGRA~1\\NVIDIA~2\\CUDA\\v12.5\\include
#cgo LDFLAGS: -LC:\\PROGRA~1\\NVIDIA~2\\CUDA\\v12.5\\lib\\x64 -lnvidia-ml
*/

import (
	"fmt"
	"log"
	"time"

	"github.com/rajasatyajit/hw-monitor/nvml"
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
	Utilization uint
	MemoryUsage uint64
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
)

func main() {
	err := nvml.Initialize()
	if err != nil {
		log.Fatalf("Could not initialize NVML: %v", err)
	}
	defer nvml.Shutdown()

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
	deviceCount, _ := nvml.DeviceCount()
	var statuses []GPUStatus
	for i := 0; i < int(deviceCount); i++ {
		device, _ := nvml.DeviceHandleByIndex(uint(i))
		name, _ := device.Name()
		utilization, _, _ := device.UtilizationRates()
		memory, _, _ := device.MemoryInfo()
		statuses = append(statuses, GPUStatus{
			Name:        name,
			Utilization: utilization,
			MemoryUsage: memory,
		})
	}
	return statuses
}

func displayStatus(status HardwareStatus) {
	fmt.Print("\033[H\033[2J") // Clear the terminal
	fmt.Print("\033[0;0H")     // Move cursor to the top left corner

	fmt.Println("Hardware Status:")
	displayCPUUsage(status.CPUUsage)
	displayMemUsage(status.MemUsage)
	displayDiskUsage(status.DiskUsage)
	displayNetIO(status.NetIO)
	displayUptime(status.Uptime)
	displayGPUStatus(status.GPUStatus)
}

func displayCPUUsage(cpuUsage []float64) {
	for i, usage := range cpuUsage {
		color := getColor(previousStatus.CPUUsage != nil && previousStatus.CPUUsage[i] != usage)
		fmt.Printf("%sCPU Core %d Usage: %0.2f%%%s\n", color, i, usage, colorReset)
	}
}

func displayMemUsage(memUsage *mem.VirtualMemoryStat) {
	color := getColor(previousStatus.MemUsage != nil && previousStatus.MemUsage.UsedPercent != memUsage.UsedPercent)
	fmt.Printf("%sMemory Usage: %0.2f%%%s\n", color, memUsage.UsedPercent, colorReset)
}

func displayDiskUsage(diskUsage []disk.UsageStat) {
	for _, usage := range diskUsage {
		var prevUsagePercent float64
		for _, prevUsage := range previousStatus.DiskUsage {
			if prevUsage.Path == usage.Path {
				prevUsagePercent = prevUsage.UsedPercent
			}
		}
		color := getColor(prevUsagePercent != usage.UsedPercent)
		fmt.Printf("%sDisk (%s) Usage: %0.2f%%%s\n", color, usage.Path, usage.UsedPercent, colorReset)
	}
}

func displayNetIO(netIO []net.IOCountersStat) {
	for _, io := range netIO {
		var prevIO net.IOCountersStat
		for _, prev := range previousStatus.NetIO {
			if prev.Name == io.Name {
				prevIO = prev
			}
		}
		colorSent := getColor(prevIO.BytesSent != io.BytesSent)
		colorRecv := getColor(prevIO.BytesRecv != io.BytesRecv)
		fmt.Printf("%sNet IO (%s) - Bytes Sent: %d%s, %sBytes Received: %d%s\n", colorSent, io.Name, io.BytesSent, colorReset, colorRecv, io.BytesRecv, colorReset)
	}
}

func displayUptime(uptime uint64) {
	color := getColor(previousStatus.Uptime != uptime)
	fmt.Printf("%sSystem Uptime: %d seconds%s\n", color, uptime, colorReset)
}

func displayGPUStatus(gpuStatus []GPUStatus) {
	for i, status := range gpuStatus {
		var prevStatus GPUStatus
		if previousStatus.GPUStatus != nil && len(previousStatus.GPUStatus) > i {
			prevStatus = previousStatus.GPUStatus[i]
		}
		colorUtilization := getColor(prevStatus.Utilization != status.Utilization)
		colorMemory := getColor(prevStatus.MemoryUsage != status.MemoryUsage)
		fmt.Printf("%sGPU %s - Utilization: %d%%%s, %sMemory Usage: %d bytes%s\n", colorUtilization, status.Name, status.Utilization, colorReset, colorMemory, status.MemoryUsage, colorReset)
	}
}

func getColor(changed bool) string {
	if changed {
		return bgRed + colorWhite // Red background with white text for updates
	}
	return colorReset // Default color
}
