// The local peer. This runs on the current node, so we have access to its
// private key, database, etc.

package zif

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
	Collection   *Collection
	Database     *Database
	Databases    cmap.ConcurrentMap

	privateKey ed25519.PrivateKey

	// a map of zif addresses to peers
	peers cmap.ConcurrentMap

	// maps public addresses to zif address
	public_to_zif cmap.ConcurrentMap

	MsgChan chan Message
}

func (lp *LocalPeer) Setup() {
	var err error

	lp.Entry.Signature = make([]byte, ed25519.SignatureSize)
	lp.peers = cmap.New()
	lp.public_to_zif = cmap.New()
	lp.Databases = cmap.New()
	lp.ZifAddress.Generate(lp.PublicKey)

	lp.Server.localPeer = lp

	lp.MsgChan = make(chan Message)

	lp.RoutingTable.Setup(lp.ZifAddress)

	lp.Collection, err = LoadCollection("./data/collection.dat")

	if err != nil {
		lp.Collection = NewCollection()
		log.Error(err.Error())
		log.Info("Created new collection")
	}

	// Loop through all the databases of other peers in ./data, load them.
	handler := func(path string, info os.FileInfo, err error) error {
		if path != "data/posts.db" && info.Name() == "posts.db" {
			r, err := regexp.Compile("data/(\\w+)/.+")

			if err != nil {
				return err
			}

			addr := r.FindStringSubmatch(path)

			db := NewDatabase(path)

			err = db.Connect()

			if err != nil {
				return err
			}

			lp.Databases.Set(addr[1], db)
		}
		return nil
	}

	filepath.Walk("./data", handler)
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
		return nil, AddressResolutionError{addr}
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

// address, router (TCP) port, dht (udp) port
func (lp *LocalPeer) Listen(addr string) {
	go lp.Server.Listen(addr)
}

func (lp *LocalPeer) GenerateKey() {
	var err error

	lp.PublicKey, lp.privateKey, err = ed25519.GenerateKey(nil)

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
	lp.PublicKey = lp.privateKey.Public().(ed25519.PublicKey)

	return nil
}

func (lp *LocalPeer) CheckSessions() {
	log.Debug("Checking peer sessions")

	// TODO: Stick this in a wait group
	for p := range lp.peers.Iter() {
		peer := p.Val.(*Peer)

		session := peer.GetSession()

		if session == nil {
			log.Debug("Peer has no session")
			return
		}

		_, err := session.Ping()

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

	if addr == lp.ZifAddress.Encode() {
		return &lp.Entry, nil
	}

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

			_, err = peer.ConnectClient(lp)

			if err != nil {
				return nil, err
			}
		}

		client, results, err := peer.Query(closest.ZifAddress.Encode())

		if closest.ZifAddress.Equals(&results[0].ZifAddress) {
			return nil, nil
		}

		closest = &results[0]
		defer client.Close()

		if err != nil {
			return nil, errors.New("Peer query failed: " + err.Error())
		}
	}

	return nil, errors.New("Address could not be resolved")
}

func (lp *LocalPeer) Close() {
	lp.CloseStreams()
	lp.Server.Close()
	lp.Database.Close()
	lp.RoutingTable.Save()
	lp.Collection.Save("./data/collection.dat")
}

func (lp *LocalPeer) AddPost(p Post, store bool) {
	log.Info("Adding post with title ", p.Title)

	lp.Collection.AddPost(p, store)
	err := lp.Database.InsertPost(p)

	if err != nil {
		log.Error(err.Error())
	}
}
