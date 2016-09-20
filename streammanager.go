// Keeps track of open TCP connections, as well as yamux sessions

package main

import (
	"errors"
	"net"

	"github.com/hashicorp/yamux"
	log "github.com/sirupsen/logrus"
)

type StreamManager struct {
	connection ConnHeader

	// Open yamux servers
	server *yamux.Session

	// Open yamux clients
	client *yamux.Session

	// Open yamux streams
	clients []Client
}

func (sm *StreamManager) Setup(lp *LocalPeer) {
	sm.server = nil
	sm.client = nil
	sm.clients = make([]Client, 0, 10)
}

func (sm *StreamManager) OpenTCP(addr string, lp *LocalPeer) (ConnHeader, error) {
	if sm.connection.conn != nil {
		return sm.connection, nil
	}

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return ConnHeader{conn, ProtocolHeader{}}, err
	}

	header, err := sm.Handshake(conn, lp)

	if err != nil {
		return ConnHeader{conn, ProtocolHeader{}}, err
	}

	pair := ConnHeader{conn, header}
	sm.connection = pair

	return pair, nil
}

func (sm *StreamManager) Handshake(conn net.Conn, lp *LocalPeer) (ProtocolHeader, error) {
	// I use the term "server" somewhat loosely. It's the "server" part of a peer.
	err := handshake_send(conn, lp)

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

func (sm *StreamManager) ConnectClient() (*yamux.Session, error) {
	// If there is already a client connected, return that.
	if sm.client != nil {
		return sm.client, nil
	}

	if sm.server != nil {
		return nil, errors.New("There is already a server connected to that socket")
	}

	client, err := yamux.Client(sm.connection.conn, nil)

	if err != nil {
		return client, err
	}

	sm.client = client

	return client, nil
}

func (sm *StreamManager) ConnectServer() (*yamux.Session, error) {
	// If there is already a server connected, return that.
	if sm.server != nil {
		return sm.server, nil
	}

	if sm.client != nil {
		return nil, errors.New("There is already a client connected to that socket")
	}

	server, err := yamux.Server(sm.connection.conn, nil)

	if err != nil {
		return server, err
	}

	sm.server = server

	return server, nil
}

func (sm *StreamManager) Close() {
	session := sm.GetSession()

	if session != nil {
		session.Close()
	}

	if sm.connection.conn != nil {
		sm.connection.conn.Close()
	}
}

func (sm *StreamManager) GetSession() *yamux.Session {
	if sm.server != nil {
		return sm.server
	}

	if sm.client != nil {
		return sm.client
	}

	return nil
}

func (sm *StreamManager) OpenStream() (Client, error) {
	var ret Client
	var err error
	session := sm.GetSession()

	if session == nil {
		return ret, errors.New("Cannot open stream, no session")
	}

	ret.conn, err = session.OpenStream()

	if err != nil {
		return ret, err
	}

	log.Debug("Opened stream (", session.NumStreams(), " total)")
	return ret, nil
}

// These streams should be coming from Server.ListenStream, as they will be started
// by the peer.
func (sm *StreamManager) AddStream(conn net.Conn) {
	var ret Client
	ret.conn = conn
	sm.clients = append(sm.clients, ret)
}

func (sm *StreamManager) GetStream(conn net.Conn) *Client {
	id := conn.(*yamux.Stream).StreamID()

	for _, c := range sm.clients {
		if c.conn.(*yamux.Stream).StreamID() == id {
			return &c
		}
	}

	return nil
}

func (sm *StreamManager) RemoveStream(conn net.Conn) {
	id := conn.(*yamux.Stream).StreamID()

	for i, c := range sm.clients {
		if c.conn.(*yamux.Stream).StreamID() == id {
			sm.clients = append(sm.clients[:i], sm.clients[i+1:]...)
		}
	}
}
