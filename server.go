package main

// tcp server

import (
	"fmt"
	"golang.org/x/crypto/ed25519"
	"net"
)

type Server struct {
	listener net.Listener
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
		return
	}
	fmt.Println("new peer")

	pHeader, err := ProtocolHeaderFromBytes(header)
	if check(err) {
		return
	}

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
		conn.Close()
		return
	}

	fmt.Println(fmt.Sprintf("%s verified", pHeader.zifAddress.Encode()))
}

func (s *Server) Handle(conn net.Conn) {

}
