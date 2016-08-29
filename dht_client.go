package main

import (
	"encoding/json"
	"fmt"
	"net"
)

// The maximum number of addresses a query will return
const DHTQuerySize = 8

type DHTClient struct {
	remote_addr *net.UDPAddr
	local_addr  *net.UDPAddr
	conn        *net.UDPConn
}

func (dc *DHTClient) Connect(addr string) {
	var err error

	// pain in the ass writing this over and over again.
	check := func(err error) {
		if err != nil {
			fmt.Printf("Error connecting: ", err.Error())
			return
		}
	}

	dc.remote_addr, err = net.ResolveUDPAddr("udp", addr)
	check(err)

	dc.local_addr, err = net.ResolveUDPAddr("udp", "0.0.0.0:0")
	check(err)

	dc.conn, err = net.DialUDP("udp", dc.local_addr, dc.remote_addr)
	check(err)
}

func (dc *DHTClient) Ping(from *LocalPeer) {
	cookie := RandBytes(20)
	msg := MakePacket(proto_dht_ping, cookie, from)

	fmt.Println("Sending ping to", dc.remote_addr.String())
	dc.conn.Write(msg.Bytes())
}

func (dc *DHTClient) Query(from *LocalPeer, target Address) {
	cookie := RandBytes(20)
	packet := MakePacket(proto_dht_query, cookie, from)
	packet.SetData(target.Bytes)

	fmt.Println("Querying for", target.Encode())

	dc.conn.Write(packet.Bytes())
}

func (dc *DHTClient) Announce(from *LocalPeer, entry Entry) {
	cookie := RandBytes(20)
	packet := MakePacket(proto_dht_announce, cookie, from)

	e, err := json.Marshal(entry)

	if err != nil {
		fmt.Println("Announce failed:", err.Error())
		return
	}

	packet.SetData(e)

	dc.conn.Write(packet.Bytes())
}
