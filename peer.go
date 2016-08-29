// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.

package main

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ed25519"
)

type Peer struct {
	RouterAddress string
	DHTAddress    string
	ZifAddress    Address
	publicKey     ed25519.PublicKey

	// TODO: Split out router connections like DHTClient -> RouterClient
	// TODO: Also think of a better name than Router. It sucks.
	router_conn net.Conn
	dht_client  DHTClient

	// Private key will only exist if this is the local peer.
}

func CreatePeerConn(conn net.Conn, pub ed25519.PublicKey, router_addr, dht_addr string) Peer {
	var ret Peer

	ret.router_conn = conn
	ret.publicKey = pub
	ret.RouterAddress = router_addr
	ret.DHTAddress = dht_addr
	ret.ZifAddress.Generate(ret.publicKey)

	return ret
}

func CreatePeer(pub ed25519.PublicKey, router_addr, dht_addr string) Peer {
	var ret Peer

	ret.publicKey = pub
	ret.RouterAddress = router_addr
	ret.DHTAddress = dht_addr
	ret.ZifAddress.Generate(ret.publicKey)

	return ret
}

func CreateUDPPeer(pub ed25519.PublicKey, addr string) Peer {
	var ret Peer

	ret.publicKey = pub
	ret.DHTAddress = addr
	ret.ZifAddress.Generate(ret.publicKey)

	ret.ConnectUDP()

	return ret
}

func (p *Peer) ConnectUDP() {
	p.dht_client.Connect(p.DHTAddress)
}

func (p *Peer) Connect() {
	// TODO: Check if an onion address. If so, DO NOT resolve. Let SOCKS do
	// that.
	addr, err := net.ResolveTCPAddr("tcp", p.RouterAddress)

	if err != nil {
		fmt.Println("Error resolving address")
		return
	}

	p.router_conn, err = net.DialTCP("tcp", nil, addr)

	if err != nil {
		fmt.Println("Error connecting")
		return
	}

	p.ConnectUDP()
}

func (p *Peer) Ping(from *LocalPeer) {
	p.dht_client.Ping(from)
}

func (p *Peer) Announce(from *LocalPeer, entry Entry) {
	p.dht_client.Announce(from, entry)
}

func (p *Peer) Query(from *LocalPeer, target Address) {
	p.dht_client.Query(from, target)
}
