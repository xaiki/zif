package main

// tcp server

import (
	"bytes"
	"net"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	listener  net.Listener
	localPeer *LocalPeer
}

func (s *Server) Listen(addr string) {
	var err error

	s.listener, err = net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	log.Info("Listening on ", addr)

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Error(err.Error())
		}

		go s.Handshake(conn)
	}
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) Handshake(conn net.Conn) error {
	header, err := handshake_recieve(conn)

	if err != nil {
		return err
	}

	err = handshake_send(conn, s.localPeer)

	if err != nil {
		return err
	}

	s.Handle(s.localPeer.CreatePeer(conn, header))
	return nil
}

func (s *Server) Handle(peer Peer) {
	msg := make([]byte, 2)
	for {
		err := net_recvall(msg, peer.client.conn)
		if err != nil {
			log.Error(err.Error())
			return
		}

		if bytes.Equal(msg, proto_terminate) {
			log.Debug(peer.ZifAddress.Encode(), " closed connection")
			return
		}

		RouteMessage(msg, peer, s.localPeer)
	}
}
