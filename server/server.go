package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"encoding/gob"
	"github.com/260by/sysmonitor/model"
)

const (
	connHost = ""
	connPort = "5000"
	connType = "tcp"
	msgLength = 1024
)

func main()  {
	address := fmt.Sprintf("%s:%s", connHost, connPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		handleRequest(conn)
	}

}

func handleRequest(conn net.Conn)  {
	// var employees = []Employee{}
	var monitors []model.Monitor
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&monitors)
	checkError(err)
	for _, monitor := range monitors {
		rAddr := conn.RemoteAddr()
		monitor.IP = strings.Split(rAddr.String(), ":")[0]
		fmt.Println(monitor)
	}
	
	conn.Close()
}

func checkError(err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}