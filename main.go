package main

import (
	"fmt"
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
}

var previousStatus HardwareStatus

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

	return HardwareStatus{
		CPUUsage:  cpuUsage,
		MemUsage:  memUsage,
		DiskUsage: diskUsages,
		NetIO:     netIO,
		Uptime:    uptime,
	}
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
}

func displayCPUUsage(cpuUsage []float64) {
	for i, usage := range cpuUsage {
		fmt.Printf("CPU Core %d Usage: %0.2f%%", i, usage)
		if previousStatus.CPUUsage != nil && previousStatus.CPUUsage[i] != usage {
			fmt.Print(" (updated)")
		}
		fmt.Println()
	}
}

func displayMemUsage(memUsage *mem.VirtualMemoryStat) {
	fmt.Printf("Memory Usage: %0.2f%%", memUsage.UsedPercent)
	if previousStatus.MemUsage != nil && previousStatus.MemUsage.UsedPercent != memUsage.UsedPercent {
		fmt.Print(" (updated)")
	}
	fmt.Println()
}

func displayDiskUsage(diskUsage []disk.UsageStat) {
	for _, usage := range diskUsage {
		fmt.Printf("Disk (%s) Usage: %0.2f%%", usage.Path, usage.UsedPercent)
		if previousStatus.DiskUsage != nil {
			for _, prevUsage := range previousStatus.DiskUsage {
				if prevUsage.Path == usage.Path && prevUsage.UsedPercent != usage.UsedPercent {
					fmt.Print(" (updated)")
				}
			}
		}
		fmt.Println()
	}
}

func displayNetIO(netIO []net.IOCountersStat) {
	for _, io := range netIO {
		fmt.Printf("Net IO (%s) - Bytes Sent: %d, Bytes Received: %d", io.Name, io.BytesSent, io.BytesRecv)
		if previousStatus.NetIO != nil {
			for _, prevIO := range previousStatus.NetIO {
				if prevIO.Name == io.Name {
					if prevIO.BytesSent != io.BytesSent {
						fmt.Print(" (sent updated)")
					}
					if prevIO.BytesRecv != io.BytesRecv {
						fmt.Print(" (received updated)")
					}
				}
			}
		}
		fmt.Println()
	}
}

func displayUptime(uptime uint64) {
	fmt.Printf("System Uptime: %d seconds", uptime)
	if previousStatus.Uptime != uptime {
		fmt.Print(" (updated)")
	}
	fmt.Println()
}
