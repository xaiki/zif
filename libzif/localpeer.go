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
	data "github.com/wjh/zif/libzif/data"
	"golang.org/x/crypto/ed25519"
)

const ResolveListSize = 1

type LocalPeer struct {
	Peer
	Entry         Entry
	RoutingTable  *RoutingTable
	Server        Server
	Collection    *data.Collection
	Database      *data.Database
	PublicAddress string
	// These are the databases of all of the peers that we have mirrored.
	Databases cmap.ConcurrentMap

	// a map of currently connected peers
	// also use to cancel reconnects :)
	Peers cmap.ConcurrentMap
	// A map of public address to Zif address
	PublicToZif cmap.ConcurrentMap

	privateKey ed25519.PrivateKey

	MsgChan chan Message

	Tor bool
}

func (lp *LocalPeer) Setup() {
	var err error

	lp.Entry.Signature = make([]byte, ed25519.SignatureSize)

	lp.Databases = cmap.New()
	lp.Peers = cmap.New()
	lp.PublicToZif = cmap.New()

	lp.ZifAddress.Generate(lp.PublicKey)

	lp.Server.localPeer = lp

	lp.MsgChan = make(chan Message)

	lp.RoutingTable, err = LoadRoutingTable("dht", lp.ZifAddress)

	if err != nil {
		panic(err)
	}

	lp.Collection, err = data.LoadCollection("./data/collection.dat")

	if err != nil {
		lp.Collection = data.NewCollection()
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

			db := data.NewDatabase(path)

			err = db.Connect()

			if err != nil {
				return err
			}

			if len(addr) < 2 {
				return nil
			}

			lp.Databases.Set(addr[1], db)
		}
		return nil
	}

	filepath.Walk("./data", handler)

	// TODO: This does not work without internet xD
	/*if lp.Entry.PublicAddress == "" {
		log.Debug("Local peer public address is nil, attempting to fetch")
		ip := external_ip()
		log.Debug("External IP is ", ip)
		lp.Entry.PublicAddress = ip
	}*/
}

// Given a direct address, for instance an IP or domain, connect to the peer there.
// This can be used for something like bootstrapping, or for something like
// connecting to a peer whose Zif address we have just resolved.
func (lp *LocalPeer) ConnectPeerDirect(addr string) (*Peer, error) {
	var peer *Peer
	var err error

	if lp.PublicToZif.Has(addr) {
		return nil, errors.New("Already connected")
	}

	peer = &Peer{}

	if err != nil {
		return nil, err
	}

	if lp.Tor {
		peer.streams.Tor = true
	}

	err = peer.Connect(addr, lp)

	if err != nil {
		return nil, err
	}

	peer.ConnectClient(lp)

	lp.Peers.Set(peer.ZifAddress.Encode(), peer)
	lp.PublicToZif.Set(addr, peer.ZifAddress.Encode())

	return peer, nil
}

func (lp *LocalPeer) GetPeer(addr string) *Peer {
	peer, has := lp.Peers.Get(addr)

	if !has {
		return nil
	}

	return peer.(*Peer)
}

// Resolved a Zif address into an entry, connects to the peer at the
// PublicAddress in the Entry, then return it. The peer is also stored in a map.
func (lp *LocalPeer) ConnectPeer(addr string) (*Peer, error) {
	var peer *Peer

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

	return peer, nil
}

func (lp *LocalPeer) SignEntry() {
	copy(lp.Entry.Signature, ed25519.Sign(lp.privateKey, EntryToBytes(&lp.Entry)))
}

// Sign any bytes.
func (lp *LocalPeer) Sign(msg []byte) []byte {
	return ed25519.Sign(lp.privateKey, msg)
}

// Pass the address to listen on. This is for the Zif connection.
func (lp *LocalPeer) Listen(addr string) {
	go lp.Server.Listen(addr)
}

// Generate a ed25519 keypair.
func (lp *LocalPeer) GenerateKey() {
	var err error

	lp.PublicKey, lp.privateKey, err = ed25519.GenerateKey(nil)

	if err != nil {
		panic(err)
	}
}

// Writes the private key to a file, in this way persisting your identity -
// all the other addresses can be generated from this, no need to save them.
// By default this file is "identity.dat"
func (lp *LocalPeer) WriteKey() error {
	if len(lp.privateKey) == 0 {
		return errors.
			New("LocalPeer does not have a private key, please generate")
	}

	err := ioutil.WriteFile("identity.dat", lp.privateKey, 0400)

	return err
}

// Read the private key from file. This is the "identity.dat" file. The public
// key is also then generated from the private key.
func (lp *LocalPeer) ReadKey() error {
	pk, err := ioutil.ReadFile("identity.dat")

	if err != nil {
		return err
	}

	lp.privateKey = pk
	lp.PublicKey = lp.privateKey.Public().(ed25519.PublicKey)

	return nil
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
// May well change, I'm unsure really. Pretty happy with it at the moment though.
func (lp *LocalPeer) Resolve(addr string) (*Entry, error) {
	log.Debug("Resolving ", addr)

	if addr == lp.ZifAddress.Encode() {
		return &lp.Entry, nil
	}

	address := DecodeAddress(addr)

	// First, find the closest peers in our routing table.
	// Satisfying if we already have the address :D
	var closest *Entry
	closest_returned := lp.RoutingTable.FindClosest(address, 1)

	if len(closest_returned) < 1 {
		return nil, errors.New("Routing table is empty")
	}

	closest = closest_returned[0]

	log.Info(len(closest_returned))

	iterations := 0

	defer func() {
		// If the peer was not in our routing table, add them!
		if iterations > 0 {
			lp.RoutingTable.Update(*closest)
		}
	}()

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

		peer = &Peer{}
		peer.streams.Tor = lp.Tor

		err := peer.Connect(closest.PublicAddress+":"+strconv.Itoa(closest.Port), lp)

		if err != nil {
			return nil, err
		}

		_, err = peer.ConnectClient(lp)

		if err != nil {
			return nil, err
		}

		// Query the peer we just connected to for the address we are hunting
		// for.
		client, results, err := peer.Query(addr)

		closest = &results[0]
		defer client.Close()

		if err != nil {
			return nil, errors.New("Peer query failed: " + err.Error())
		}

		iterations++
	}

	return nil, errors.New("Address could not be resolved")
}

func (lp *LocalPeer) Close() {
	lp.CloseStreams()
	lp.Server.Close()
	lp.Database.Close()
	lp.RoutingTable.Save("dht")
	lp.Collection.Save("./data/collection.dat")
}

func (lp *LocalPeer) AddPost(p data.Post, store bool) {
	log.Info("Adding post with title ", p.Title)

	lp.Collection.AddPost(p, store)
	err := lp.Database.InsertPost(p)

	if err != nil {
		log.Error(err.Error())
	}
}
