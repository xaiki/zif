// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.
// Just a bit of a wrapper for the client really, that contains most of the
// networking code, this mostly has the data and a few other things.

package main

import "golang.org/x/crypto/ed25519"
import log "github.com/sirupsen/logrus"
import "github.com/hashicorp/yamux"

import "strconv"

type Peer struct {
	ZifAddress    Address
	PublicAddress string
	client        Client
	publicKey     ed25519.PublicKey
	localPeer     *LocalPeer
}

func (p *Peer) OpenStream() (Client, error) {
	log.Debug("Opening new stream for ", p.ZifAddress.Encode())
	var ret Client
	client, err := p.localPeer.CreateServer(p.ZifAddress.Encode())

	if err != nil {
		log.Error(err.Error())
		return ret, err
	}

	stream, err := client.OpenStream()

	if err != nil {
		log.Error(err.Error())
		return ret, err
	}

	ret.conn = stream

	return ret, nil
}

func (p *Peer) Close() {
	p.client.Close()
}

func (p *Peer) Ping() {
	log.Info("Pinging ", p.ZifAddress.Encode())
	p.client.Ping()
}

func (p *Peer) Pong() {
	log.Debug("Ping from", p.ZifAddress.Encode())

	return_stream, err := p.OpenStream()

	if err != nil {
		log.Error(err.Error())
		return
	}

	return_stream.Pong()
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
		peer, err := p.localPeer.ConnectPeerDirect(i.PublicAddress + ":" + strconv.Itoa(i.Port))

		if err != nil ||
			i.ZifAddress.Equals(&entry.ZifAddress) {

			continue
		}

		peer.client.conn.Write(proto_dht_announce)
		peer.client.SendEntry(&entry, sig)
	}
}

// Very much the same as the counterpart in Server, just a little different as
// this peer is the one that *started* the TCP connection.
func (p *Peer) ListenStream(header ProtocolHeader, client *yamux.Session) {
	msg := make([]byte, 2)

	for {
		stream, err := client.Accept()

		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Debug("Client accepted new stream from ", header.zifAddress.Encode())

		net_recvall(msg, stream)

		RouteMessage(msg, p.localPeer.CreatePeer(stream, header), p.localPeer)
	}
}
