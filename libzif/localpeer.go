// The local peer. This runs on the current node, so we have access to its
// private key, database, etc.

package libzif

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/streamrail/concurrent-map"
	data "github.com/wjh/zif/libzif/data"
	"github.com/wjh/zif/libzif/dht"
	"github.com/wjh/zif/libzif/proto"
	"golang.org/x/crypto/ed25519"
)

const ResolveListSize = 1

type LocalPeer struct {
	Peer
	Entry         Entry
	RoutingTable  *dht.RoutingTable
	Server        proto.Server
	Collection    *data.Collection
	Database      *data.Database
	PublicAddress string
	// These are the databases of all of the peers that we have mirrored.
	Databases   cmap.ConcurrentMap
	Collections cmap.ConcurrentMap

	SearchProvider *data.SearchProvider

	// a map of currently connected peers
	// also use to cancel reconnects :)
	Peers cmap.ConcurrentMap
	// A map of public address to Zif address
	PublicToZif cmap.ConcurrentMap

	privateKey ed25519.PrivateKey

	Tor bool
}

func (lp *LocalPeer) Setup() {
	var err error

	lp.Entry.Signature = make([]byte, ed25519.SignatureSize)

	lp.Databases = cmap.New()
	lp.Collections = cmap.New()
	lp.Peers = cmap.New()
	lp.PublicToZif = cmap.New()

	lp.Address().Generate(lp.PublicKey())

	lp.RoutingTable, err = dht.LoadRoutingTable("dht", *lp.Address())

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

	lp.SearchProvider = data.NewSearchProvider()
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

	lp.Peers.Set(peer.Address().String(), peer)
	lp.PublicToZif.Set(addr, peer.Address().String())

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
		return nil, data.AddressResolutionError{addr}
	}

	// now should have an entry for the peer, connect to it!
	log.Debug("Connecting to ", entry.Address.String())

	peer, err = lp.ConnectPeerDirect(entry.PublicAddress + ":" + strconv.Itoa(entry.Port))

	// Caller can go on to choose a seed to connect to, not quite the end of the
	// world :P
	if err != nil {
		log.WithField("peer", addr).Info("Failed to connect")

		return nil, err
	}

	return peer, nil
}

func (lp *LocalPeer) SignEntry() {
	data, _ := lp.Entry.Bytes()
	copy(lp.Entry.Signature, ed25519.Sign(lp.privateKey, data))
}

// Sign any bytes.
func (lp *LocalPeer) Sign(msg []byte) []byte {
	return ed25519.Sign(lp.privateKey, msg)
}

// Pass the address to listen on. This is for the Zif connection.
func (lp *LocalPeer) Listen(addr string) {
	go lp.Server.Listen(addr, lp)
}

// Generate a ed25519 keypair.
func (lp *LocalPeer) GenerateKey() {
	var err error

	lp.publicKey, lp.privateKey, err = ed25519.GenerateKey(nil)

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
	lp.publicKey = lp.privateKey.Public().(ed25519.PublicKey)

	return nil
}

func (lp *LocalPeer) Resolve(addr string) (*Entry, error) {
	log.Debug("Resolving ", addr)

	if addr == lp.Address().String() {
		return &lp.Entry, nil
	}

	address := dht.DecodeAddress(addr)

	closest := lp.RoutingTable.FindClosest(address, dht.MaxBucketSize)
	current := make(map[string]bool)

	for _, i := range closest {
		entry, err := JsonToEntry(i.Value)

		if err != nil {
			continue
		}

		if i.Key.Equals(&address) {
			return entry, nil
		}

		current[i.Key.String()] = true
	}

	if len(closest) < 1 {
		return nil, errors.New("Routing table is empty")
	}

	// Create a worker pool of goroutines working on resolving an address, then
	// proceed to block on a result.

	workers := 3
	addresses := make(chan string, dht.MaxBucketSize)
	results := make(chan dht.Pairs, dht.MaxBucketSize*workers)

	defer close(results)
	defer close(addresses)

	// Setup the workers
	for i := 0; i < workers; i++ {
		go lp.worker(i, addr, addresses, results)
	}

	// Feed in the initial addresses
	for _, i := range closest {
		entry, err := JsonToEntry(i.Value)

		if err != nil {
			continue
		}

		addresses <- fmt.Sprintf("%s:%s", entry.PublicAddress, entry.Port)
	}

	// Listen for results from workers, feeding addresses we have not seen before
	// back into the system to be queried. Terminates when we have found what we
	// are looking for.
	for i := range results {
		for _, j := range i {
			// If this is a new address we have not yet seen
			if _, ok := current[j.Key.String()]; !ok {
				entry, err := JsonToEntry(j.Value)

				if err != nil {
					continue
				}

				if j.Key.Equals(&address) {
					return entry, nil
				}

				addresses <- fmt.Sprintf("%s:%s", entry.PublicAddress, entry.Port)

				closest = append(closest, j)
				current[j.Key.String()] = true
			}
		}
	}

	return nil, errors.New("Failed to resolve entry")
}

type workResult struct {
	id int
}

// Pass this the id of the worker, the address we are looking for, a channel
// that will be sending peers to attempt to query, and a channel to send query
// results on. Note that the addresses being passed in via channel are those
// of public internet addresses and not Zif addresses. They should have
// already been resolved :)
func (lp *LocalPeer) worker(id int, address string, addresses <-chan string, results chan<- dht.Pairs) {

	// If any errors occur, just skip that peer and attempt to work with the
	// next. No point terminating if we meet one dodgy peer.

	found := make(map[string]bool)

	for i := range addresses {
		p := lp.GetPeer(i)

		if p == nil {
			p = &Peer{}

			err := p.Connect(i, lp)

			if err != nil {
				continue
			}

			client, res, err := p.Query(address)

			if err != nil {
				continue
			}

			defer client.Close()

			results <- res
		}
	}
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
