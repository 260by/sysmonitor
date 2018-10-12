package main

import (
	"fmt"
	"flag"
	"net"
	"net/http"
	"os"
	// "strings"
	"encoding/gob"
	"github.com/260by/sysmonitor/model"
	"github.com/260by/tools/gconfig"
	// "github.com/go-xorm/xorm"
)

const (
	connType = "tcp"
)

type Config struct {
	Monitor struct {
		IP string
		Port int
	}
	Database struct {
		Driver string
		Dsn string
		ShowSQL bool
		Migrate bool
	}
	HTTPServer struct {
		IP string
		Port int
	}
}

func main()  {
	var configFile = flag.String("config", "config.toml", "Configration file")
	var migrate = flag.Bool("migrate", false, "Sync database table structure")
	var initUser = flag.Bool("init", false, "Init admin user")
	flag.Parse()

	var config = Config{}
	err := gconfig.Parse(*configFile, &config)
	if err != nil {
		panic(err)
	}

	// 同步数据结构
	if *migrate {
		orm, err := model.Connect(config.Database.Driver, config.Database.Dsn, config.Database.ShowSQL)
		if err != nil {
			panic(err)
		}
		err = model.Migrate(orm)
		if err != nil {
			panic(err)
		}
		if err == nil {
			fmt.Println("Sync database table structure is success.")
			os.Exit(0)
		}
		defer orm.Close()
	}

	if *initUser {
		orm, err := model.Connect(config.Database.Driver, config.Database.Dsn, config.Database.ShowSQL)
		if err != nil {
			panic(err)
		}
		err = model.InitUser(orm)
		if err != nil {
			panic(err)
		}
		if err == nil {
			fmt.Println("Init admin user is success.")
			os.Exit(0)
		}
		defer orm.Close()
	}

	address := fmt.Sprintf("%s:%v", config.Monitor.IP, config.Monitor.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer listener.Close()

	httpListenAddr := fmt.Sprintf("%s:%v", config.HTTPServer.IP, config.HTTPServer.Port)
	go startWebServer(httpListenAddr, config)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		handleRequest(conn, config)
	}

}

func handleRequest(conn net.Conn, config Config)  {
	// 连接数据库
	orm, err := model.Connect(config.Database.Driver, config.Database.Dsn, config.Database.ShowSQL)
	if err != nil {
		panic(err)
	}

	var monitors []model.Monitor
	dec := gob.NewDecoder(conn)
	err = dec.Decode(&monitors)
	if err != nil {
		fmt.Println("Accept data err: ", err)
	}
	for _, monitor := range monitors {
		// rAddr := conn.RemoteAddr()
		// monitor.IP = strings.Split(rAddr.String(), ":")[0]
		fmt.Println(monitor)
	}
	orm.Insert(&monitors)
	orm.Close()
	
	conn.Close()
}

func checkError(err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func startWebServer(address string, config Config)  {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)  {
		fmt.Fprintf(w, "Hello world\n")
	})

	http.HandleFunc("/api/assets", func(w http.ResponseWriter, r *http.Request)  {
		fmt.Fprintln(w, "assets list")
	})
	
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}