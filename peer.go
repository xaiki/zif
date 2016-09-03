// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.
// Just a bit of a wrapper for the client really, that contains most of the
// networking code, this mostly has the data and a few other things.

package main

import "golang.org/x/crypto/ed25519"
import "fmt"
import log "github.com/sirupsen/logrus"

type Peer struct {
	ZifAddress    Address
	PublicAddress string
	client        Client
	publicKey     ed25519.PublicKey
	localPeer     *LocalPeer
}

func (p *Peer) Connect(addr string) {
	p.client.Connect(addr)
}

func (p *Peer) Close() {
	p.client.Close()
}

func (p *Peer) Ping() {
	fmt.Println("Pinging", p.ZifAddress.Encode())
	p.client.Ping()
}

func (p *Peer) Pong() {
	fmt.Println("Ping from", p.ZifAddress.Encode())
	p.client.Pong()
}

func (p *Peer) Handshake() error {
	header, err := p.client.Handshake(p.localPeer)

	if err != nil {
		return err
	}

	p.ZifAddress.Bytes = make([]byte, AddressBinarySize)
	copy(p.ZifAddress.Bytes[:], header.zifAddress.Bytes)
	copy(p.publicKey, header.PublicKey[:])

	return err
}

func (p *Peer) Who() (Entry, error) {
	return p.client.Who()
}

func (p *Peer) SendWho() {
	p.client.SendEntry(&p.localPeer.Entry, p.localPeer.entrySig[:])
}

func (p *Peer) Announce(entry *Entry) {
	p.client.Announce(entry, p.localPeer.entrySig[:])
}

func (p *Peer) RecievedAnnounce() {
	entry, err := recieve_entry(p.client.conn)

	if err != nil {
		p.Close()
		fmt.Println("Error recieving announce:", err.Error())
		return
	}

	var addr Address
	addr.Generate(entry.PublicKey[:])

	log.Debug("Announce from ", addr.Encode())
}
