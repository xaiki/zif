package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"strings"

	log "github.com/sirupsen/logrus"
)

func SetupLocalPeer(addr string, newAddr bool) *LocalPeer {
	var lp LocalPeer

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
	var newAddr = flag.Bool("new", false, "Ignore identity file and create a new address")

	var http = flag.String("http", "127.0.0.1:8080", "HTTP address and port")

	flag.Parse()

	port, _ := strconv.Atoi(strings.Split(*addr, ":")[1])

	lp := SetupLocalPeer(fmt.Sprintf("%s:%v", *addr), *newAddr)
	lp.Entry.Name = "Zif"
	lp.Entry.Desc = "Decentralize all the things! :D"
	lp.Entry.Port = port
	lp.Entry.PublicAddress = ""
	lp.Entry.ZifAddress = lp.ZifAddress
	lp.Entry.PublicKey = lp.publicKey
	lp.Entry.PublicAddress = "127.0.0.1"
	lp.SignEntry()

	lp.Listen(*addr)

	log.Info("My address: ", lp.ZifAddress.Encode())

	var httpServer HTTPServer
	httpServer.localPeer = lp
	go httpServer.ListenHTTP(*http)

	// Listen for SIGINT
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	for _ = range sigchan {
		lp.RoutingTable.Save()

		os.Exit(0)
	}
}
