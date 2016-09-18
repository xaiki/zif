// The local peer. This runs on the current node, so we have access to its
// private key, database, etc.

package main

import (
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/streamrail/concurrent-map"
	"golang.org/x/crypto/ed25519"
)

type LocalPeer struct {
	Peer
	Entry        Entry
	RoutingTable RoutingTable
	Server       Server

	privateKey ed25519.PrivateKey

	peers         cmap.ConcurrentMap
	public_to_zif cmap.ConcurrentMap
}

func (lp *LocalPeer) Setup() {
	lp.Entry.Signature = make([]byte, ed25519.SignatureSize)
	lp.peers = cmap.New()
	lp.public_to_zif = cmap.New()
	lp.ZifAddress.Generate(lp.publicKey)
}

func (lp *LocalPeer) GetPeer(addr string) *Peer {
	if p, ok := lp.peers.Get(addr); ok {
		peer := p.(*Peer)
		return peer
	}

	return nil
}

func (lp *LocalPeer) SignEntry() {
	copy(lp.Entry.Signature, ed25519.Sign(lp.privateKey, EntryToBytes(&lp.Entry)))
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

func (lp *LocalPeer) CheckSessions() {
	log.Debug("Checking peer sessions")

	// TODO: Stick this in a wait group
	for p := range lp.peers.Iter() {
		peer := p.Val.(*Peer)

		_, err := peer.GetSession().Ping()

		if err != nil {
			log.Debug(err.Error())
			log.Debug("Removing ", peer.ZifAddress.Encode(), " from map")
			lp.peers.Remove(peer.ZifAddress.Encode())
			return
		}

		if peer.GetSession().IsClosed() {
			log.Warn("TCP session has closed")
			log.Debug("Removing ", peer.ZifAddress.Encode(), " from map")
			lp.peers.Remove(peer.ZifAddress.Encode())
			return
		}
	}
}

func (lp *LocalPeer) ListenStream(peer *Peer) {
	var err error
	session := peer.GetSession()

	if session == nil {
		log.Info("Peer has no active session, starting server")
		session, err = peer.ConnectServer()

		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	for {
		stream, err := session.Accept()

		if err != nil {
			if err.Error() == "EOF" {
				log.Info("Peer closed connection")
			} else {
				log.Error(err.Error())
			}

			peer.localPeer.CheckSessions()

			return
		}

		log.Debug("Accepted stream (", session.NumStreams(), " total)")

		peer.AddStream(stream)

		go lp.HandleStream(peer, stream)
	}
}
