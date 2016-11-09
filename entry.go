package zif

// This is an entry into the DHT. It is used to connect to a peer given just
// it's Zif address.
type Entry struct {
	ZifAddress    Address
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
	Port      int

	// Used in the FindClosest function, for sorting.
	distance Address
}

func (e *Entry) SetLocalPeer(lp *LocalPeer) {
	e.ZifAddress = lp.ZifAddress
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
