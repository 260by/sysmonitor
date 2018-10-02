package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"encoding/gob"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/cpu"
	// psnet "github.com/shirou/gopsutil/net"
	"github.com/260by/sysmonitor/model"
	// "strconv"
	// "time"
)




func main()  {
	var addr = flag.String("server", "127.0.0.1:5000", "Server address")
	flag.Parse()

	for {
		tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err !=nil {
			fmt.Println("Connect to Server error: ", err)
			time.Sleep(time.Second * 120)
			continue
		}

		var monitors []model.Monitor
		for i := 0; i < 2; i++ {
			var monitor model.Monitor
			monitor.CreateTime = time.Now().Unix()

			getMonitorData(&monitor)
			// getIP(&monitor)
			monitors = append(monitors, monitor)
			time.Sleep(time.Second * 15)
		}
		
		fmt.Println(monitors)
		enc := gob.NewEncoder(conn)
		err = enc.Encode(monitors)
		if err != nil {
			fmt.Println("Send monitor data to server error: ", err)
			continue
		}
		conn.Close()
	}
}

func getMonitorData(monitor *model.Monitor)  {
	getHostName(monitor)
	getCPU(monitor)
	getMem(monitor)
	getDisk(monitor)
}

func getHostName(monitor *model.Monitor)  {
	hostStats, err := host.Info()
	if err != nil {
		fmt.Println("Get host info err: ", err)
		return
	}
	monitor.HostName = hostStats.Hostname
}

func getCPU(monitor *model.Monitor)  {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("Get CPU percent err: ", err)
		return
	}
	monitor.CPUPercent = cpuPercent[0]
}

func getMem(monitor *model.Monitor)  {
	m, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}

	monitor.MemoryPercent = m.UsedPercent
}

func getDisk(monitor *model.Monitor)  {
	partition, _ := disk.Partitions(true)

	for _, p := range partition {
		if p.Fstype == "ext3" || p.Fstype == "ext4" || p.Fstype == "xfs" {
			diskInfo, err := disk.Usage(p.Mountpoint)
			if err != nil {
				panic(err)
			}
			
			monitor.DisksPercent += fmt.Sprintf("%s:%v ", p.Mountpoint, diskInfo.UsedPercent)
		}
	}
}

func getIP(monitor *model.Monitor)  {
	// interStats, err := psnet.Interfaces()
	// if err != nil {
	// 	fmt.Println("Get interfaces err: ", err)
	// 	return
	// }

	// fmt.Println(psnet.Connections(kind))1
	os.Exit(0)
}