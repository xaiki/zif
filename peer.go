// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.
// Just a bit of a wrapper for the client really, that contains most of the
// networking code, this mostly has the data and a few other things.

package main

import (
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
)

import "github.com/hashicorp/yamux"

type Peer struct {
	ZifAddress    Address
	PublicAddress string
	publicKey     ed25519.PublicKey
	localPeer     *LocalPeer
	streams       StreamManager
}

func NewPeer(local_peer *LocalPeer) *Peer {
	var ret Peer

	ret.streams.local_peer = local_peer
	ret.localPeer = local_peer

	return &ret
}

func (p *Peer) Connect(addr string) error {
	pair, err := p.streams.OpenTCP(addr)

	if err != nil {
		return err
	}

	p.publicKey = pair.header.PublicKey[:]
	p.ZifAddress = pair.header.zifAddress

	p.localPeer.peers.Set(p.ZifAddress.Encode(), p)
	p.localPeer.public_to_zif.Set(addr, p.ZifAddress.Encode())

	return nil
}

func (p *Peer) SetTCP(pair ConnHeader) {
	p.streams.connection = pair

	p.publicKey = pair.header.PublicKey[:]
	p.ZifAddress = pair.header.zifAddress
}

func (p *Peer) ConnectServer() (*yamux.Session, error) {
	return p.streams.ConnectServer()
}

func (p *Peer) ConnectClient() (*yamux.Session, error) {
	client, err := p.streams.ConnectClient()

	if err != nil {
		return client, err
	}

	go listen_stream(p)

	return client, err
}

func (p *Peer) GetSession() *yamux.Session {
	return p.streams.GetSession()
}

func (p *Peer) Terminate() {
	p.streams.Close()
}

func (p *Peer) OpenStream() (Client, error) {
	return p.streams.OpenStream()
}

func (p *Peer) AddStream(conn net.Conn) {
	p.streams.AddStream(conn)
}

func (p *Peer) RemoveStream(conn net.Conn) {
	p.streams.RemoveStream(conn)
}

func (p *Peer) GetStream(conn net.Conn) *Client {
	return p.streams.GetStream(conn)
}

func (p *Peer) CloseStreams() {
	p.streams.Close()
}

func (p *Peer) Ping() *Client {
	stream, err := p.OpenStream()

	if err != nil {
		log.Error(err.Error())
	}

	log.Info("Pinging ", p.ZifAddress.Encode())
	stream.Ping()

	return &stream
}

func (p *Peer) Pong() *Client {
	log.Debug("Ping from ", p.ZifAddress.Encode())

	stream, _ := p.OpenStream()
	stream.Pong()

	return &stream
}
func (p *Peer) Who() (*Client, Entry, error) {
	log.Info("Sending who query to ", p.ZifAddress.Encode())

	stream, _ := p.OpenStream()
	entry, err := stream.Who()

	return &stream, entry, err
}

func (p *Peer) Announce() *Client {
	log.Debug("Sending announce to ", p.ZifAddress.Encode())

	if p.localPeer.Entry.PublicAddress == "" {
		log.Debug("Local peer public address is nil, attempting to fetch")
		ip := external_ip()
		log.Debug("External IP is ", ip)
		p.localPeer.Entry.PublicAddress = ip
	}

	p.localPeer.SignEntry()

	stream, _ := p.OpenStream()
	stream.Announce(&p.localPeer.Entry, p.localPeer.entrySig[:])

	return &stream
}

func (p *Peer) RecievedAnnounce(stream net.Conn, from *Peer) {
	log.Debug("Recieved announce")
	defer stream.Close()

	entry, sig, err := recieve_entry(stream)

	if err != nil {
		log.Error(err.Error())
		return
	}

	var addr Address
	addr.Generate(entry.PublicKey[:])

	log.Debug("Announce from ", from.ZifAddress.Encode())

	saved := p.localPeer.RoutingTable.Update(entry)

	if saved {
		log.Info("Saved new peer ", addr.Encode())
	}

	// next up, tell other people!
	closest := p.localPeer.RoutingTable.FindClosest(addr, BucketSize)

	// TODO: Parallize this
	for _, i := range closest {
		if i.ZifAddress.Equals(&entry.ZifAddress) || i.ZifAddress.Equals(&from.ZifAddress) {
			continue
		}

		peer := p.localPeer.GetPeer(i.ZifAddress.Encode())

		if peer == nil {
			log.Debug("Connecting to new peer")
			peer := NewPeer(p.localPeer)
			err = peer.Connect(i.PublicAddress + ":" + strconv.Itoa(i.Port))
		}

		if err != nil {
			log.Warn("Failed to connect to peer: ", err.Error())
			continue
		}

		peer_stream, err := peer.OpenStream()
		defer peer_stream.Close()

		if err != nil {
			log.Error(err.Error())
		}

		peer_stream.conn.Write(proto_dht_announce)
		peer_stream.SendEntry(&entry, sig)
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

		//RouteMessage(msg, p.localPeer.CreatePeer(header), p.localPeer)
	}
}
