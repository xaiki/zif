package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"strings"

	"github.com/wjh/zif"

	log "github.com/sirupsen/logrus"
)

func SetupLocalPeer(addr string, newAddr bool) *zif.LocalPeer {
	var lp zif.LocalPeer

	if !newAddr {
		if lp.ReadKey() != nil {
			lp.GenerateKey()
			lp.WriteKey()
		}
	} else {
		lp.GenerateKey()
	}

	lp.Setup()
	lp.RoutingTable.Setup(lp.ZifAddress)

	return &lp
}

func main() {

	log.SetLevel(log.DebugLevel)
	formatter := new(log.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "15:04:05"
	log.SetFormatter(formatter)

	var addr = flag.String("address", "0.0.0.0:5050", "Bind address")
	var db_path = flag.String("database", "./data/posts.db", "Posts database path")
	var newAddr = flag.Bool("new", false, "Ignore identity file and create a new address")

	var http = flag.String("http", "127.0.0.1:8080", "HTTP address and port")

	flag.Parse()

	port, _ := strconv.Atoi(strings.Split(*addr, ":")[1])

	lp := SetupLocalPeer(fmt.Sprintf("%s:%v", *addr), *newAddr)
	lp.Entry.Name = "Zif"
	lp.Entry.Desc = "Decentralize all the things! :D"
	lp.Entry.Port = port
	lp.Entry.PublicAddress = ""
	lp.Entry.PublicAddress = "127.0.0.1"
	lp.Entry.SetLocalPeer(lp)
	lp.SignEntry()

	lp.Database = zif.NewDatabase(*db_path)

	err := lp.Database.Connect()

	if err != nil {
		log.Fatal(err.Error())
	}

	lp.Listen(*addr)

	log.Info("My address: ", lp.ZifAddress.Encode())

	var httpServer zif.HTTPServer
	httpServer.LocalPeer = lp
	go httpServer.ListenHTTP(*http)

	// Listen for SIGINT
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	for _ = range sigchan {
		lp.RoutingTable.Save()
		lp.Database.Close()

		os.Exit(0)
	}
}
