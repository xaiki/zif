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

		log.Debug("Handshaking new connection")
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

	//s.Handle(s.localPeer.CreatePeer(conn, header))

	s.localPeer.connections[header.zifAddress.Encode()] = ConnHeader{conn, header}

	go s.ListenStream(conn, header)

	return nil
}

func (s *Server) ListenStream(conn net.Conn, header ProtocolHeader) {
	session, err := s.localPeer.CreateServer(header.zifAddress.Encode())
	conn.Write(proto_ok)

	if err != nil {
		log.Error(err.Error())
		conn.Close()
		return
	}

	for {
		log.Debug("Session listening for streams")
		stream, err := session.Accept()

		if err != nil {
			log.Error(err.Error())
			return
		}

		go s.Handle(s.localPeer.CreatePeer(stream, header))
	}
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
