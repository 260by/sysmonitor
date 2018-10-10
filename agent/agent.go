package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"strings"
	"encoding/gob"
	"github.com/260by/tools/sys/cpu"
	"github.com/260by/tools/sys/mem"
	"github.com/260by/tools/sys/disk"
	"github.com/260by/tools/sys/load"
	network "github.com/260by/tools/sys/net"
	"github.com/260by/sysmonitor/model"
	// "github.com/robfig/cron"
)

func main()  {
	var addr = flag.String("server", "127.0.0.1:5000", "Server address")
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err !=nil {
			fmt.Println("Connect to Server error: ", err)
			time.Sleep(time.Second * 120)
			continue
		}

		var monitors []model.Monitor
		for i := 0; i < 2; i++ {
			var monitor model.Monitor

			localAddr := conn.LocalAddr()
			monitor.IP = strings.Split(localAddr.String(), ":")[0]

			getMonitorData(&monitor)
			monitors = append(monitors, monitor)
			time.Sleep(time.Second * 12)
		}
		
		// fmt.Println(monitors)
		enc := gob.NewEncoder(conn)
		err = enc.Encode(monitors)
		if err != nil {
			fmt.Println("Send monitor data to server error: ", err)
			continue
		}
		conn.Close()
	}
}

func getMonitorData(monitor *model.Monitor) {
	var err error
	monitor.CreateTime = time.Now().Unix()
	monitor.HostName, err = os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	monitor.CPUPercent = cpu.Usage()
	monitor.MemoryPercent = mem.Usage()
	
	var diskPercent string
	for k, v := range disk.Usage() {
		diskPercent += fmt.Sprintf("%s:%f ", k, v)
	}
	monitor.DisksPercent = diskPercent

	loadavg := load.Avg()
	monitor.Load1 = loadavg[0]
	monitor.Load5 = loadavg[1]
	monitor.Load15 = loadavg[2]

	tcpStats := network.TCPState()
	monitor.TCPEstablished = tcpStats["ESTABLISHED"]
	monitor.TCPTimeWait = tcpStats["TIME_WAIT"]
}