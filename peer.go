// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.
// Just a bit of a wrapper for the client really, that contains most of the networking code, this mostly has the data and a few other things.

package main

import (
	"errors"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
)

import "github.com/hashicorp/yamux"

type Peer struct {
	ZifAddress    Address
	PublicAddress string
	publicKey     ed25519.PublicKey
	streams       StreamManager

	entry *Entry
}

func (p *Peer) Announce(lp *LocalPeer) {
	log.Debug("Sending announce to ", p.ZifAddress.Encode())

	if lp.Entry.PublicAddress == "" {
		log.Debug("Local peer public address is nil, attempting to fetch")
		ip := external_ip()
		log.Debug("External IP is ", ip)
		lp.Entry.PublicAddress = ip
	}

	lp.SignEntry()

	stream, _ := p.OpenStream()
	defer stream.Close()
	stream.Announce(&lp.Entry)
}

func (p *Peer) Connect(addr string, lp *LocalPeer) error {
	pair, err := p.streams.OpenTCP(addr, lp)

	if err != nil {
		return err
	}

	p.publicKey = pair.header.PublicKey[:]
	p.ZifAddress = pair.header.zifAddress

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

	return client, err
}

func (p *Peer) GetSession() *yamux.Session {
	return p.streams.GetSession()
}

func (p *Peer) Terminate() {
	p.streams.Close()
}

func (p *Peer) OpenStream() (Client, error) {
	if p.GetSession() == nil {
		return Client{}, errors.New("Peer session nil")
	}

	if p.GetSession().IsClosed() {
		return Client{}, errors.New("Peer session closed")
	}
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

func (p *Peer) Entry() (*Entry, error) {
	if p.entry != nil {
		return p.entry, nil
	}

	client, entries, err := p.Query(p.ZifAddress.Encode())
	defer client.Close()

	if err != nil {
		return nil, err
	}

	if len(entries) < 1 {
		return nil, errors.New("Query did not return an entry")
	}

	p.entry = &entries[0]

	return &entries[0], nil
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

func (p *Peer) Bootstrap(rt *RoutingTable) (*Client, error) {
	log.Info("Bootstrapping from ", p.streams.connection.conn.RemoteAddr())

	initial, err := p.Entry()

	if err != nil {
		return nil, err
	}
	rt.Update(*initial)

	stream, _ := p.OpenStream()

	return &stream, stream.Bootstrap(rt, rt.LocalAddress)
}

func (p *Peer) Query(address string) (*Client, []Entry, error) {
	log.Info("Querying for ", address)

	stream, _ := p.OpenStream()
	entry, err := stream.Query(address)
	return &stream, entry, err
}
