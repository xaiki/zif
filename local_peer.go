// The local peer. This runs on the current node, so we have access to its
// private key, database, etc.

package main

import (
	"crypto"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ed25519"
)

type LocalPeer struct {
	Peer
	Router       Router
	DHTServer    DHTServer
	Entry        Entry
	RoutingTable RoutingTable
	privateKey   ed25519.PrivateKey
}

func (lp *LocalPeer) Setup() {
	lp.ZifAddress.Generate(lp.publicKey)

	// Would rather not have to do this :/
	lp.Router.LocalPeer = lp
	lp.DHTServer.localPeer = lp
}

// address, router (TCP) port, dht (udp) port
func (lp *LocalPeer) Listen(addr string, router int, dht int) {
	lp.RouterAddress = fmt.Sprintf("%s:%v", addr, router)
	lp.DHTAddress = fmt.Sprintf("%s:%v", addr, dht)

	go lp.Router.Listen(addr, router)
	go lp.DHTServer.Listen(addr, dht)
}

func (lp *LocalPeer) GenerateKey() {
	var err error

	lp.publicKey, lp.privateKey, err = ed25519.GenerateKey(nil)

	if err != nil {
		panic(err)
	}
}

// Writes the private key to a file, in this way persisting your identity -
// all the other addresses can be generated from this, no need to save them.
func (lp *LocalPeer) WriteKey() error {
	if len(lp.privateKey) == 0 {
		return errors.
			New("LocalPeer does not have a private key, please generate")
	}

	err := ioutil.WriteFile("identity.dat", lp.privateKey, 0400)

	return err
}

func (lp *LocalPeer) ReadKey() error {
	pk, err := ioutil.ReadFile("identity.dat")

	if err != nil {
		return err
	}

	lp.privateKey = pk
	lp.publicKey = lp.privateKey.Public().(ed25519.PublicKey)

	return nil
}

func (lp *LocalPeer) Handshake(p *Peer) {
	header := make([]byte, 0, 100)

	if len(lp.RouterAddress) > 64 {
		panic(errors.New("Local peer address too long, cannot handshake (>64)"))
	}

	var router_addr [64]byte
	var dht_addr [64]byte

	if lp.RouterAddress == "" {
		local := lp.Router.listener.Addr().String()
		copy(router_addr[:], local)
	} else {
		copy(router_addr[:], lp.RouterAddress)
	}

	if lp.DHTAddress == "" {
		local := lp.DHTServer.addr.String()
		copy(dht_addr[:], local)
	} else {
		copy(dht_addr[:], lp.DHTAddress)
	}

	header = append(header, proto_zif...)
	header = append(header, proto_version...)
	header = append(header, router_addr[:]...)
	header = append(header, dht_addr[:]...)
	header = append(header, lp.publicKey...)

	p.router_conn.Write(header)

	cookie := make([]byte, 20)
	net_recvall(cookie, p.router_conn)

	sig, err := lp.privateKey.Sign(nil, cookie, crypto.Hash(0))

	if err != nil {
		panic(err)
	}

	p.router_conn.Write(sig)
}
