// The local peer. This runs on the current node, so we have access to its
// private key, database, etc.

package main

import (
	"errors"
	"io/ioutil"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/streamrail/concurrent-map"
	"golang.org/x/crypto/ed25519"
)

const ResolveListSize = 1

type LocalPeer struct {
	Peer
	Entry        Entry
	RoutingTable RoutingTable
	Server       Server

	privateKey ed25519.PrivateKey

	// a map of zif addresses to peers
	peers cmap.ConcurrentMap

	// maps public addresses to zif address
	public_to_zif cmap.ConcurrentMap
}

func (lp *LocalPeer) Setup() {
	lp.Entry.Signature = make([]byte, ed25519.SignatureSize)
	lp.peers = cmap.New()
	lp.public_to_zif = cmap.New()
	lp.ZifAddress.Generate(lp.publicKey)
}

// Creates a peer, connects to a public address
func (lp *LocalPeer) ConnectPeerDirect(addr string) (*Peer, error) {
	lp.CheckSessions()

	zif_address, ok := lp.public_to_zif.Get(addr)

	if ok {
		return lp.GetPeer(zif_address.(string)), nil
	}

	var peer Peer
	err := peer.Connect(addr, lp)

	if err != nil {
		return nil, err
	}

	return lp.AddPeer(&peer), nil
}

// Creates a peer, resolves a zif address then connects to the assosciated
// public address
func (lp *LocalPeer) ConnectPeer(addr string) (*Peer, error) {
	lp.CheckSessions()

	peer := lp.GetPeer(addr)

	if peer != nil {
		return peer, nil
	}

	entry, err := lp.Resolve(addr)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if entry == nil {

	}

	// now should have an entry for the peer, connect to it!
	log.Debug("Connecting to ", entry.ZifAddress.Encode())
	peer, err = lp.ConnectPeerDirect(entry.PublicAddress + ":" + strconv.Itoa(entry.Port))

	if err != nil {
		return nil, err
	}

	return lp.AddPeer(peer), nil
}

func (lp *LocalPeer) AddPeer(peer *Peer) *Peer {
	lp.peers.Set(peer.ZifAddress.Encode(), peer)
	lp.public_to_zif.Set(peer.PublicAddress, peer.ZifAddress.Encode())

	return peer
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

	log.Debug("Peer sessions checked")
}

// At the moment just query for the closest known peer

// This takes a Zif address as a string and attempts to resolve it to an entry.
// This may be fast, may be a little slower. Will recurse its way through as
// many Queries as needed, getting closer to the target until either it cannot
// be found or is found.
// Cannot be found if a Query returns nothing, in this case the address does not
// exist on the DHT. Otherwise we should get to a peer that either has the entry,
// or one that IS the peer we are hunting.

// Takes a string as the API will just be passing a Zif address as a string.
// May well change, I'm unsure really.
func (lp *LocalPeer) Resolve(addr string) (*Entry, error) {
	log.Debug("Resolving ", addr)
	address := DecodeAddress(addr)

	// First, find the closest peers in our routing table.
	// Satisfying if we already have the address :D
	var closest *Entry
	closest_returned := lp.RoutingTable.FindClosest(address, ResolveListSize)

	if len(closest_returned) < 1 {
		return nil, errors.New("Routing table is empty")
	}

	closest = closest_returned[0]

	for {
		// Check the current closest known peers. First iteration this will be
		// the ones from our routing table.
		if closest == nil {
			return nil, errors.New("Address could not be resolved")
			// The first in the slice is the closest, if we have this entry in our table
			// then this will be it.
		} else if closest.ZifAddress.Equals(&address) {
			log.Debug("Found ", closest.ZifAddress.Encode())
			return closest, nil
		}

		var peer *Peer

		// If the peer is not already connected, then connect.
		if peer = lp.GetPeer(closest.ZifAddress.Encode()); peer == nil {

			var peer Peer
			err := peer.Connect(closest.PublicAddress+":"+strconv.Itoa(closest.Port), lp)

			if err != nil {
				return nil, err
			}

			_, err = peer.ConnectClient()

			if err != nil {
				return nil, err
			}
		}

		client, results, err := peer.Query(closest.ZifAddress.Encode())
		closest = &results[0]
		defer client.Close()

		if err != nil {
			return nil, errors.New("Peer query failed: " + err.Error())
		}
	}

	return nil, errors.New("Address could not be resolved")
}
