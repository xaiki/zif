package main

import "flag"
import "fmt"

func SetupLocalPeer(addr string, newAddr bool) LocalPeer {
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

	lp.Listen(addr)

	return lp
}

func main() {

	var addr = flag.String("address", "0.0.0.0", "Bind address")
	var port = flag.Int("port", 5050, "TCP port")
	var newAddr = flag.Bool("new", false, "Ignore identity file and create a new address")

	var http = flag.String("http", "127.0.0.1:8080", "HTTP address and port")

	flag.Parse()

	lp := SetupLocalPeer(fmt.Sprintf("%s:%v", *addr, *port), *newAddr)
	lp.Entry.Name = "Zif"
	lp.Entry.Desc = "Decentralize all the things! :D"
	lp.Entry.Port = *port
	lp.Entry.PublicAddress = ""
	lp.Entry.ZifAddress = lp.ZifAddress
	lp.Entry.PublicKey = lp.publicKey

	fmt.Println("My address:", lp.ZifAddress.Encode())

	var httpServer HTTPServer
	httpServer.localPeer = &lp
	httpServer.ListenHTTP(*http)
}
