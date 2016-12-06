package libzif

import "github.com/wjh/zif/libzif/dht"

// This is an entry into the DHT. It is used to connect to a peer given just
// it's Zif address.
type Entry struct {
	Address       dht.Address
	Name          string
	Desc          string
	PublicAddress string
	PublicKey     []byte
	PostCount     int

	// The owner of this entry should have signed it, we need to store the
	// sigature. It's actually okay as we can verify that a peer owns a public
	// key by generating an address from it - if the address is not the peers,
	// then Mallory is just using someone elses entry for their own address.
	Signature []byte
	// Signature of the root hash of a hash list representing all of the posts
	// a peer has.
	CollectionSig []byte
	Port          int

	// Essentially just a list of other peers who have this entry in their table.
	// They may or may not actually have pieces, so mirror/piece requests may go
	// awry.
	// This is... weird. It is not signed. It is not verified.
	// While the idea of doing the above irks me somewhat, as potentially bad
	// actors can make themselves become seeds - they will fail with piece
	// requests. Hashes of pieces will not match.
	// Removing the requirement for this to be both signed and verified means
	// that any peer can become a seed without that much work at all. Again,
	// seed lists can then be updated *without* the requirement that the origin
	// peer actually be online in the first place.
	// TODO: Switch this to be a struct containing the last time this peer was
	// announced as a peer, then the list can be periodically culled.
	Seeds [][]byte

	// Used in the FindClosest function, for sorting.
	distance dht.Address
}

func (e *Entry) SetLocalPeer(lp *LocalPeer) {
	e.Address = lp.Address
	e.PublicKey = lp.PublicKey
}

type Entries []*Entry

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Entries) Less(i, j int) bool {
	return e[i].distance.Less(&e[j].distance)
}
