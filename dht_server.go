// a UDP server that responds to DHT requests

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"net"
	"strings"
)

type DHTServer struct {
	server    *net.UDPConn
	addr      *net.UDPAddr
	publicKey ed25519.PublicKey
	localPeer *LocalPeer
}

func (ds *DHTServer) Listen(addr string, port int) {
	var err error

	ds.addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", addr, port))

	if err != nil {
		fmt.Println("Error: ", err.Error())
		panic(err)
	}

	ds.server, err = net.ListenUDP("udp", ds.addr)
	fmt.Println("Started UDP server listening on", ds.addr.String())
	defer ds.server.Close()

	buf := make([]byte, 4096)

	for {
		n, addr, err := ds.server.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error: ", err.Error())
		}

		go ds.Handle(n, buf, addr)
	}
}

func (ds *DHTServer) Handle(n int, buf []byte, addr *net.UDPAddr) {
	fmt.Println("Packet recieved")
	if len(buf) < 2 {
		return
	}
	packet := PacketFromBytes(buf)

	sender_ip := strings.Split(addr.String(), ":")[0]

	var client DHTClient
	client.Connect(fmt.Sprintf("%s:%v", sender_ip, packet.UDPPort()))

	var address Address
	address.Generate(packet.PublicKey())

	if bytes.Equal(packet.protoMsg[:], proto_dht_ping) {
		ds.handlePing(&client, &packet, addr)
	} else if bytes.Equal(packet.protoMsg[:], proto_dht_pong) {
		// TODO: Map requests to replies
		fmt.Println("Pong from", addr.String())
	} else if bytes.Equal(packet.protoMsg[:], proto_dht_query) {
		ds.handleQuery(&client, &packet, addr)
	} else if bytes.Equal(packet.protoMsg[:], proto_dht_announce) {
		ds.handleAnnounce(&client, &packet, addr)
	}

}

func (ds *DHTServer) handlePing(client *DHTClient, packet *Packet, from *net.UDPAddr) {
	fmt.Println("Ping from", client.remote_addr.String())
	pong := MakePacket(proto_dht_pong, packet.cookie[:], ds.localPeer)

	client.conn.Write(pong.Bytes())
}

func (ds *DHTServer) handleQuery(client *DHTClient, packet *Packet, from *net.UDPAddr) {
	var target Address
	target.Bytes = packet.data

	fmt.Println("Query for", target.Encode())

	// TODO: Finish this
}

func (ds *DHTServer) handleAnnounce(client *DHTClient, packet *Packet, from *net.UDPAddr) {
	var entry Entry
	err := json.Unmarshal(packet.data, &entry)

	if err != nil {
		fmt.Println("Failed to decode announce:", err.Error())
		return
	}

	var packetAddress Address
	packetAddress.Generate(packet.PublicKey())

	if len(entry.PublicAddress) == 0 {
		entry.PublicAddress = from.IP.String()
	}

	fmt.Println("Announce for", entry.Name)

	// (maybe) insert into our routing table
	ds.localPeer.RoutingTable.Update(entry)

	// announce this to the closest k peers in the table
	closest := ds.localPeer.RoutingTable.FindClosest(entry.ZifAddress, BucketSize)

	for _, i := range closest {
		if i.ZifAddress.Equals(&entry.ZifAddress) || packetAddress.Equals(&i.ZifAddress) {
			continue
		}
		peer := CreateUDPPeer(i.PublicKey, fmt.Sprintf("%s:%v", i.PublicAddress, i.UDPPort))
		peer.Announce(ds.localPeer, entry)
	}
}
