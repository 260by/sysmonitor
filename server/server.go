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
	// "github.com/go-xorm/xorm"
	"github.com/260by/sysmonitor/server/handlefunc"
	"github.com/260by/sysmonitor/config"
)

var configFile = flag.String("config", "config.toml", "Configration file")
var migrate = flag.Bool("migrate", false, "Sync database table structure")
var initUser = flag.Bool("init", false, "Init admin user")

func main()  {
	flag.Parse()

	conf, err := config.Parse(*configFile)
	if err != nil {
		panic(err)
	}

	// 同步数据结构
	if *migrate {
		orm, err := model.Connect(conf.Database.Driver, conf.Database.Dsn, conf.Database.ShowSQL)
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
		orm, err := model.Connect(conf.Database.Driver, conf.Database.Dsn, conf.Database.ShowSQL)
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
	
	httpAddr := fmt.Sprintf("%s:%v", conf.HTTPServer.IP, conf.HTTPServer.Port)
	go startWebServer(httpAddr)
	
	tcpAddr := fmt.Sprintf("%s:%v", conf.Monitor.IP, conf.Monitor.Port)
	listener, err := net.Listen("tcp", tcpAddr)
	fmt.Printf("Starting server...\nListening on %s.\n", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, conf)
	}
}

func handleConnection(conn net.Conn, c *config.Config)  {
	// 连接数据库
	orm, err := model.Connect(c.Database.Driver, c.Database.Dsn, c.Database.ShowSQL)
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
}

func startWebServer(addr string) {
	http.HandleFunc("/", handlefunc.Index)

	http.HandleFunc("/api/assets", func(w http.ResponseWriter, r *http.Request)  {
		fmt.Fprintln(w, "assets list")
	})
	
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Listening on %s.", addr)
}