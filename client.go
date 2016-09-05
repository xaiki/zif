package main

import (
	"encoding/binary"
	"net"

	log "github.com/sirupsen/logrus"
)

const EntryLengthMax = 1024

type Client struct {
	conn net.Conn
}

func (c *Client) Terminate() {
	c.conn.Write(proto_terminate)
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Ping() {
	c.conn.Write(proto_ping)
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
