package main

// tcp server

import (
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

		log.Debug("Handshaking new connection")
		go s.Handshake(conn)
	}
}

func (s *Server) Handshake(conn net.Conn) {
	header, err := handshake(conn, s.localPeer)

	if err != nil {
		log.Error(err.Error())
		return
	}

	peer := NewPeer(s.localPeer)
	peer.SetTCP(ConnHeader{conn, header})

	s.localPeer.peers.Set(peer.ZifAddress.Encode(), peer)

	listen_stream(peer)
}

func (s *Server) Close() {
	s.listener.Close()
}
