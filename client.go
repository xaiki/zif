package main

import (
	"bytes"
	"encoding/binary"
	"net"

	log "github.com/sirupsen/logrus"
)

const EntryLengthMax = 1024

type Client struct {
	conn net.Conn
}

// Attempt to connect a client to an address, return true on success.
func (c *Client) Connect(addr string) bool {
	var err error

	c.conn, err = net.Dial("tcp", addr)

	if err != nil {
		log.Error(err.Error())
		return false
	}

	return true
}

func (c *Client) Close() {
	c.conn.Write(proto_terminate)
	c.conn.Close()
}

func (c *Client) Handshake(lp *LocalPeer) (ProtocolHeader, error) {
	// I use the term "server" somewhat loosely. It's the "server" part of a peer.
	err := handshake_send(c.conn, lp)

	// server now knows that we are definitely who we say we are.
	// but...
	// is the server who we think it is?
	// better check!
	server_header, err := handshake_recieve(c.conn)

	if err != nil {
		return server_header, err
	}

	server_header.zifAddress.Generate(server_header.PublicKey[:])

	return server_header, nil
}

func (c *Client) Ping() bool {
	c.conn.Write(proto_ping)

	buf := make([]byte, 2)
	net_recvall(buf, c.conn)

	return bytes.Equal(buf, proto_pong)
}

func (c *Client) Pong() {
	c.conn.Write(proto_pong)
}

func (c *Client) Who() (Entry, error) {
	c.conn.Write(proto_who)

	entry, _, err := recieve_entry(c.conn)

	if err != nil {
		c.Close()
	}

	return entry, err
}

func (c *Client) SendEntry(e *Entry, sig []byte) {
	json, err := EntryToJson(e)

	if err != nil {
		log.Error(err.Error())
		c.conn.Close()
		return
	}

	length := len(json)
	length_b := make([]byte, 8)
	binary.PutVarint(length_b, int64(length))

	c.conn.Write(length_b)
	c.conn.Write(json)
	c.conn.Write(sig)
}

func (c *Client) Announce(e *Entry, sig []byte) {
	c.conn.Write(proto_dht_announce)
	c.SendEntry(e, sig)
}
