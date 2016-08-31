// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.

package main

import "golang.org/x/crypto/ed25519"

type Peer struct {
	ZifAddress    Address
	PublicAddress string
	publicKey     ed25519.PublicKey
}

func CreatePeer(pub ed25519.PublicKey, addr string, port int) Peer {
	var ret Peer

	ret.publicKey = pub
	ret.ZifAddress.Generate(ret.publicKey)

	return ret
}

func (p *Peer) Connect() {
}

func (p *Peer) Ping() {
}

func (p *Peer) Announce(entry *Entry) string {
	return ""
}
