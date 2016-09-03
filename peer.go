// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.
// Just a bit of a wrapper for the client really, that contains most of the
// networking code, this mostly has the data and a few other things.

package main

import "golang.org/x/crypto/ed25519"
import log "github.com/sirupsen/logrus"
import "strconv"

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
	log.Debug("Pinging", p.ZifAddress.Encode())
	p.client.Ping()
}

func (p *Peer) Pong() {
	log.Debug("Ping from", p.ZifAddress.Encode())
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

func (p *Peer) Announce() {
	log.Debug("Sending announce to ", p.ZifAddress.Encode())

	if p.localPeer.Entry.PublicAddress == "" {
		p.localPeer.Entry.PublicAddress = external_ip()
	}

	p.localPeer.SignEntry()

	p.client.Announce(&p.localPeer.Entry, p.localPeer.entrySig[:])
}

func (p *Peer) RecievedAnnounce() {
	log.Debug("Recieved announce")
	entry, sig, err := recieve_entry(p.client.conn)

	if err != nil {
		p.Close()
		log.Error(err.Error())
		return
	}

	var addr Address
	addr.Generate(entry.PublicKey[:])

	log.Debug("Announce from ", addr.Encode())

	saved := p.localPeer.RoutingTable.Update(entry)

	if saved {
		log.Info("Saved new peer ", addr.Encode())
	}

	// next up, tell other people!
	closest := p.localPeer.RoutingTable.FindClosest(addr, BucketSize)

	// TODO: Parallize this
	for _, i := range closest {
		peer, err := p.localPeer.ConnectPeer(i.PublicAddress + ":" + strconv.Itoa(i.Port))

		if err != nil ||
			i.ZifAddress.Equals(&entry.ZifAddress) {

			continue
		}

		peer.client.conn.Write(proto_dht_announce)
		peer.client.SendEntry(&entry, sig)
	}
}
