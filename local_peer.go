// The local peer. This runs on the current node, so we have access to its
// private key, database, etc.

package main

import (
	"errors"
	"io/ioutil"
	"net"

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

	streams StreamManager
}

func (lp *LocalPeer) Setup() {
	lp.ZifAddress.Generate(lp.publicKey)
}

func (lp *LocalPeer) CreatePeer(conn net.Conn, header ProtocolHeader) Peer {
	var ret Peer
	ret.Setup()

	ret.ZifAddress = header.zifAddress
	ret.publicKey = header.PublicKey[:]
	ret.clients[0].conn = conn
	ret.localPeer = lp

	return ret
}

func (lp LocalPeer) ConnectPeerDirect(addr string) (Peer, error) {
	var ret Peer
	ret.Setup()

	ret.localPeer = &lp

	pair, err := lp.streams.OpenTCP(addr)

	if err != nil {
		return ret, err
	}

	check_ok(pair.conn)

	ret.ZifAddress = pair.header.zifAddress

	client, err := lp.streams.ConnectClient(pair)
	log.Debug("Created local client")

	if err != nil || client == nil {
		return ret, err
	}

	stream, err := client.OpenStream()
	log.Debug("Opened client stream #", client.NumStreams())

	if err != nil {
		return ret, err
	}

	ret.clients[0].conn = stream

	// Next up we need to make a server listen as well, this way connections can
	// be initiated by our peer.
	go ret.ListenStream(pair.header, client)

	return ret, nil
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
