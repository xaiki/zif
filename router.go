// Handles peer connections

package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net"

	"golang.org/x/crypto/ed25519"
)

type Router struct {
	listener  net.Listener
	LocalPeer *LocalPeer
}

// Will block. If you don't want that, run it in a goroutine.
func (r *Router) Listen(host string, port int) {
	var err error
	r.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%v", host, port))

	if err != nil {
		fmt.Println("Error listening")
	}

	fmt.Println("Started TCP server listening on", fmt.Sprintf("%s:%v", host, port))

	for {
		// TODO: Log and handle error. Should not terminate program.
		conn, err := r.listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection")
		}

		go r.HandleConnection(conn)
	}
}

func (r *Router) HandleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 164)
	net_recvall(buf, conn)

	// verify the peer actually has the private key, the cookie is very unlikely
	// to be known in advance by an attacker.
	cookie := make([]byte, 20)
	_, err := rand.Read(cookie)
	conn.Write(cookie)

	if err != nil {
		panic(err)
	}

	sig := make([]byte, 64)
	net_recvall(sig, conn)

	verified := ed25519.Verify(buf[132:], cookie, sig)

	if !verified {
		return
	}

	r.HandlePeer(CreatePeerConn(conn, buf[132:], string(buf[4:68]), string(buf[68:132])))
}

// Once HandleConnection knows the details of the peer from the header, it
// creates a peer and calls this.
func (r *Router) HandlePeer(p Peer) {
	fmt.Println("Peer connected: ", fmt.Sprintf("%s:%s", p.RouterAddress, p.DHTAddress))

	// get a message from the peer
	// these are all 2 bytes, defined in protocol_info.go
	buf := make([]byte, 2)
	net_recvall(buf, p.router_conn)

	// Can't use a switch :(
	if bytes.Equal(buf, proto_ping) {
		fmt.Println("PING: ", p.RouterAddress)
		p.router_conn.Write(proto_pong)
	}
}
