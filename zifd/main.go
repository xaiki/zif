package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"strings"

	zif "github.com/wjh/zif/libzif"
	data "github.com/wjh/zif/libzif/data"

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

	return &lp
}

func main() {

	log.SetLevel(log.DebugLevel)
	formatter := new(log.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "15:04:05"
	log.SetFormatter(formatter)

	os.Mkdir("./data", 0777)

	var addr = flag.String("address", "0.0.0.0:5050", "Bind address")
	var db_path = flag.String("database", "./data/posts.db", "Posts database path")
	var newAddr = flag.Bool("new", false, "Ignore identity file and create a new address")
	var tor = flag.Bool("tor", false, "Start hidden service and proxy connections through tor")

	var http = flag.String("http", "127.0.0.1:8080", "HTTP address and port")

	flag.Parse()

	port, _ := strconv.Atoi(strings.Split(*addr, ":")[1])

	lp := SetupLocalPeer(fmt.Sprintf("%s:%v", *addr), *newAddr)

	if *tor {
		_, onion, err := zif.SetupZifTorService(5050, 9051, "./tor/cookie")

		if err == nil {
			lp.PublicAddress = onion
			lp.Entry.PublicAddress = onion
			lp.Tor = true
			lp.Peer.Streams().Tor = true
		}
	}

	lp.Entry.Name = "Zif"
	lp.Entry.Desc = "Decentralize all the things! :D"
	lp.Entry.Port = port
	lp.Entry.SetLocalPeer(lp)
	lp.SignEntry()

	post := data.Post{}
	post.InfoHash = "foo"
	post.Title = "Foo"
	log.Info(data.PostToString(&post, "|", ""))

	lp.Database = data.NewDatabase(*db_path)

	err := lp.Database.Connect()

	if err != nil {
		log.Fatal(err.Error())
	}

	lp.Listen(*addr)

	log.Info("My address: ", lp.ZifAddress.Encode())

	var commandServer zif.CommandServer
	commandServer.LocalPeer = lp
	var httpServer zif.HttpServer
	httpServer.CommandServer = &commandServer
	go httpServer.ListenHttp(*http)

	// Listen for SIGINT
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	for _ = range sigchan {
		lp.Close()

		os.Exit(0)
	}
}
