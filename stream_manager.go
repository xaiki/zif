// Keeps track of open TCP connections, as well as yamux sessions

package main

import (
	"errors"
	"net"

	"github.com/hashicorp/yamux"
	log "github.com/sirupsen/logrus"
)

type StreamManager struct {
	// Currently open TCP connections
	connections map[string]ConnHeader

	// Open yamux servers
	servers map[string]*yamux.Session

	// Open yamux clients
	clients map[string]*yamux.Session

	// Open yamux streams
	streams map[string][]net.Conn

	local_peer *LocalPeer
}

func (sm *StreamManager) Setup(lp *LocalPeer) {
	sm.connections = make(map[string]ConnHeader)
	sm.servers = make(map[string]*yamux.Session)
	sm.clients = make(map[string]*yamux.Session)

	sm.local_peer = lp
}

func (sm *StreamManager) OpenTCP(addr string) (ConnHeader, error) {
	if c, ok := sm.connections[addr]; ok {
		return c, nil
	}

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return ConnHeader{conn, ProtocolHeader{}}, err
	}

	header, err := sm.Handshake(conn)
	pair := ConnHeader{conn, header}
	sm.connections[header.zifAddress.Encode()] = pair

	return pair, nil
}

func (sm *StreamManager) Handshake(conn net.Conn) (ProtocolHeader, error) {
	// I use the term "server" somewhat loosely. It's the "server" part of a peer.
	err := handshake_send(conn, sm.local_peer)

	// server now knows that we are definitely who we say we are.
	// but...
	// is the server who we think it is?
	// better check!
	server_header, err := handshake_recieve(conn)

	if err != nil {
		return server_header, err
	}

	server_header.zifAddress.Generate(server_header.PublicKey[:])

	return server_header, nil
}

func (sm *StreamManager) ConnectClient(pair ConnHeader) (*yamux.Session, error) {
	addr := pair.header.zifAddress.Encode()

	// If there is already a client connected, return that.
	if c, ok := sm.clients[addr]; ok {
		return c, nil
	}

	if s, ok := sm.servers[addr]; ok {
		return nil, errors.New("There is already a server connected to that socket")
	}

	client, err := yamux.Client(sm.connections[addr].conn, nil)

	if err != nil {
		return client, err
	}

	sm.clients[addr] = client

	return client, nil
}

func (sm *StreamManager) ConnectServer(pair ConnHeader) (*yamux.Session, error) {
	addr := pair.header.zifAddress.Encode()

	// If there is already a client connected, return that.
	if s, ok := sm.servers[addr]; ok {
		return s, nil
	}

	if c, ok := sm.clients[addr]; ok {
		return nil, errors.New("There is already a client connected to that socket")
	}

	server, err := yamux.Server(sm.connections[addr].conn, nil)

	if err != nil {
		return server, err
	}

	sm.servers[addr] = server

	return server, nil
}

func (sm *StreamManager) GetSession(addr string) *yamux.Session {
	if s, ok := sm.servers[addr]; ok {
		return s
	}

	if c, ok := sm.clients[addr]; ok {
		return c
	}

	return nil
}

func (sm *StreamManager) OpenStream(addr string) (net.Conn, error) {
	log.Debug("Opening new stream for ", addr)

	session := sm.GetSession(addr)

	if session == nil {
		return nil, errors.New("Cannot open stream, no session")
	}

	return session.OpenStream()
}
