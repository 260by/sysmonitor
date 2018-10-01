package main

import (
	"flag"
	"fmt"
	"net"
	// "os"
	"time"
	"encoding/gob"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/cpu"
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

		var statsList []model.Stats
		for i := 0; i < 2; i++ {
			var stats model.Stats
			stats.CreateTime = time.Now().Unix()
			hostStats, err := host.Info()
			if err != nil {
				fmt.Println("Get host info err: ", err)
				continue
			}
			stats.HostName = hostStats.Hostname
			// cpuStats, err := cpu.Info()
			// if err != nil {
			// 	fmt.Println("Get CPU info err: ", err)
			// 	continue
			// }
			cpuPercent, err := cpu.Percent(0, false)
			if err != nil {
				fmt.Println("Get CPU percent err: ", err)
				continue
			}
			stats.CPUPercent = cpuPercent[0]
			getDisk(&stats)
			getMem(&stats)
			statsList = append(statsList, stats)
			time.Sleep(time.Second * 15)
		}
		
		fmt.Println(statsList)
		enc := gob.NewEncoder(conn)
		err = enc.Encode(statsList)
		if err != nil {
			fmt.Println("Send monitor data to server error: ", err)
			continue
		}
		conn.Close()
	}
}

func getMem(stats *model.Stats)  {
	m, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}

	stats.Memory.Total = m.Total
	stats.Memory.Free = m.Free
	stats.Memory.UsePercent = m.UsedPercent
}

func getDisk(stats *model.Stats)  {
	partition, _ := disk.Partitions(true)

	for _, p := range partition {
		var d model.Disk
		if p.Fstype == "ext3" || p.Fstype == "ext4" || p.Fstype == "xfs" {
			diskInfo, err := disk.Usage(p.Mountpoint)
			if err != nil {
				panic(err)
			}
			d.MountPoint = p.Mountpoint
			d.Total = diskInfo.Total
			d.Free = diskInfo.Free
			d.Used = diskInfo.Used
			d.UsePercent = diskInfo.UsedPercent
			stats.Disks = append(stats.Disks, d)
			// fmt.Printf("挂载点:%s\n磁盘总容量:%v\n使用容量:%v\n使用率:%.2f%%\n", p.Mountpoint, diskInfo.Total>>30, diskInfo.Used>>30,diskInfo.UsedPercent)
		}
	}
}