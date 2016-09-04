// The local peer. This runs on the current node, so we have access to its
// private key, database, etc.

package main

import (
	"errors"
	"io/ioutil"
	"net"

	"github.com/hashicorp/yamux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
)

type LocalPeer struct {
	Peer
	Entry        Entry
	RoutingTable RoutingTable
	Server       Server

	privateKey ed25519.PrivateKey
	entrySig   [64]byte

	// Currently open TCP connections
	connections map[string]ConnHeader

	// Open yamux servers
	servers map[string]*yamux.Session

	// Open yamux clients
	clients map[string]*yamux.Session
}

func (lp *LocalPeer) Setup() {
	lp.connections = make(map[string]ConnHeader)
	lp.servers = make(map[string]*yamux.Session)
	lp.clients = make(map[string]*yamux.Session)

	lp.ZifAddress.Generate(lp.publicKey)
}

func (lp *LocalPeer) CreatePeer(conn net.Conn, header ProtocolHeader) Peer {
	var ret Peer

	ret.ZifAddress = header.zifAddress
	ret.publicKey = header.PublicKey[:]
	ret.client.conn = conn
	ret.localPeer = lp

	return ret
}

func (lp *LocalPeer) ConnectDirect(addr string) (ConnHeader, error) {
	if c, ok := lp.connections[addr]; ok {
		return c, nil
	}

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return ConnHeader{conn, ProtocolHeader{}}, err
	}

	header, err := lp.Handshake(conn)
	pair := ConnHeader{conn, header}
	lp.connections[header.zifAddress.Encode()] = pair

	return pair, nil
}

func (lp LocalPeer) ConnectPeerDirect(addr string) (Peer, error) {
	var ret Peer
	ret.localPeer = &lp

	pair, err := lp.ConnectDirect(addr)

	if err != nil {
		return ret, err
	}

	check_ok(pair.conn)

	client, err := lp.CreateClient(pair.header.zifAddress.Encode())
	log.Debug("Created local client")

	if err != nil || client == nil {
		return ret, err
	}

	stream, err := client.OpenStream()
	log.Debug("Opened client stream")

	if err != nil {
		return ret, err
	}

	ret.client.conn = stream

	return ret, nil
}

func (lp *LocalPeer) Handshake(conn net.Conn) (ProtocolHeader, error) {
	// I use the term "server" somewhat loosely. It's the "server" part of a peer.
	err := handshake_send(conn, lp)

	// server now knows that we are definitely who we say we are.
	// but...
	// is the server who we think it is?
	// better check!
	server_header, err := handshake_recieve(conn)

	if err != nil {
		return server_header, err
	}

	server_header.zifAddress.Generate(server_header.PublicKey[:])

	return server_header, nil
}

func (lp *LocalPeer) CreateClient(addr string) (*yamux.Session, error) {
	if c, ok := lp.clients[addr]; ok {
		return c, nil
	}

	client, err := yamux.Client(lp.connections[addr].conn, nil)

	if err != nil {
		return client, err
	}

	lp.clients[addr] = client

	return client, nil
}

func (lp *LocalPeer) CreateServer(addr string) (*yamux.Session, error) {
	if s, ok := lp.servers[addr]; ok {
		return s, nil
	}

	server, err := yamux.Server(lp.connections[addr].conn, nil)

	if err != nil {
		return server, err
	}

	lp.servers[addr] = server

	return server, nil
}

func (lp *LocalPeer) SignEntry() {
	copy(lp.entrySig[:], ed25519.Sign(lp.privateKey, EntryToBytes(&lp.Entry)))
}

func (lp *LocalPeer) Sign(msg []byte) []byte {
	return ed25519.Sign(lp.privateKey, msg)
}

func (lp *LocalPeer) ProtocolHeader() ProtocolHeader {
	var ph ProtocolHeader

	copy(ph.Zif[:], proto_zif)
	copy(ph.Version[:], proto_version)
	copy(ph.PublicKey[:], lp.publicKey[:])

	return ph
}

// address, router (TCP) port, dht (udp) port
func (lp *LocalPeer) Listen(addr string) {
	go lp.Server.Listen(addr)

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
