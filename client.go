package main

import (
	"fmt"
	"net"
)

type Client struct {
	conn      net.Conn
	localPeer *LocalPeer
}

// Returns a client and true on success
func ConnectClient(addr string, lp *LocalPeer) (Client, bool) {
	var client Client
	client.localPeer = lp
	return client, client.Connect(addr)
}

// Attempt to connect a client to an address, return true on success.
func (c *Client) Connect(addr string) bool {
	var err error

	c.conn, err = net.Dial("tcp", addr)

	if err != nil {
		fmt.Println("Error:", err.Error())
		return false
	}

	return true
}

func (c *Client) Handshake() {
	fmt.Println("Handshaking with", c.conn.RemoteAddr().String())
	//ph := c.localPeer.ProtocolHeader()

	header := c.localPeer.ProtocolHeader()
	c.conn.Write(header.Bytes())

	// The server will want us to sign this. Proof of identity and all that.
	cookie := make([]byte, 20)
	net_recvall(cookie, c.conn)

	sig := c.localPeer.Sign(cookie)
	c.conn.Write(sig)
}
