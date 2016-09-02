// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.
// Just a bit of a wrapper for the client really, that contains most of the
// networking code, this mostly has the data and a few other things.

package main

import "golang.org/x/crypto/ed25519"

type Peer struct {
	ZifAddress    Address
	PublicAddress string
	client        Client
	publicKey     ed25519.PublicKey
}

func (p *Peer) Connect(addr string) {
	p.client.Connect(addr)
}

func (p *Peer) Ping() {
	p.client.Ping()
}

func (p *Peer) Pong() {
	p.client.Pong()
}

func (p *Peer) Announce(entry *Entry) string {
	return ""
}
