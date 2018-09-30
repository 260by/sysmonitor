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
	"github.com/260by/sysmonitor/model"
	// "strconv"
	// "time"
)




func main()  {
	var addr = flag.String("server", "127.0.0.1:5000", "Server address")
	flag.Parse()

	for {
		tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
		checkError(err)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		checkError(err)

		var statsList []model.Stats
		var stats model.Stats
		getDisk(&stats)
		getMem(&stats)
		
		// enc := gob.NewEncoder(conn)
		// err = enc.Encode(statsList)
		// checkError(err)

		time.Sleep(time.Second * 15)
	}

	// tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	// checkError(err)
	// conn, err := net.DialTCP("tcp", nil, tcpAddr)
	// checkError(err)

	// // _, err = conn.Write([]byte("Hello World!")) // send string
	// enc := gob.NewEncoder(conn)
	// err = enc.Encode(employee)
	// checkError(err)

	// conn.Close()
	// os.Exit(0)
}

func checkError(err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
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