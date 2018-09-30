package main

import (
	"fmt"
	"net"
	"os"
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
	var statsList []model.Stats
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&statsList)
	checkError(err)
	for _, stats := range statsList {
		fmt.Println(stats)
	}
	
	conn.Close()
}

func checkError(err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}