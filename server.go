package main

// tcp server

import (
	"bytes"
	"fmt"
	"net"

	"golang.org/x/crypto/ed25519"
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

	fmt.Println("Listening on", addr)

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			fmt.Println("Error accepting:", err.Error())
		}

		go s.Handshake(conn)
	}
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) Handshake(conn net.Conn) {

	check := func(e error) bool {
		if e != nil {
			fmt.Println("Error:", e.Error())
			conn.Close()
			return true
		}

		return false
	}

	header := make([]byte, ProtocolHeaderSize)
	err := net_recvall(header, conn)
	if check(err) {
		conn.Write(proto_no)
		return
	}

	pHeader, err := ProtocolHeaderFromBytes(header)
	if check(err) {
		conn.Write(proto_no)
		return
	}

	conn.Write(proto_no)

	fmt.Println("Incoming connection from", pHeader.zifAddress.Encode())

	// Send the client a cookie for them to sign, this proves they have the
	// private key, and it is highly unlikely an attacker has a signed cookie
	// cached.
	cookie, err := RandBytes(20)
	if check(err) {
		return
	}

	conn.Write(cookie)

	sig := make([]byte, ed25519.SignatureSize)
	net_recvall(sig, conn)

	verified := ed25519.Verify(pHeader.PublicKey[:], cookie, sig)

	if !verified {
		fmt.Println("Failed to verify peer", pHeader.zifAddress.Encode())
		conn.Write(proto_no)
		conn.Close()
		return
	}

	conn.Write(proto_ok)

	fmt.Println(fmt.Sprintf("%s verified", pHeader.zifAddress.Encode()))

	s.Handle(s.localPeer.CreatePeer(conn, pHeader))
}

func (s *Server) Handle(peer Peer) {
	msg := make([]byte, 2)

	for {
		net_recvall(msg, peer.client.conn)

		if bytes.Equal(msg, proto_terminate) {
			fmt.Println("Peer closed connection")
			return
		}

		RouteMessage(msg, peer)
	}
}
