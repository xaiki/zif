package main

import (
	"encoding/binary"
	"errors"
	"net"

	log "github.com/sirupsen/logrus"
)

const EntryLengthMax = 1024

type Client struct {
	conn net.Conn
}

func NewClient(stream net.Conn) Client {
	var ret Client
	ret.conn = stream
	return ret
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

func (c *Client) SendEntry(e *Entry) {
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
}

func (c *Client) Announce(e *Entry) {
	c.conn.Write(proto_dht_announce)
	c.SendEntry(e)
}

func (c *Client) Query(address string) {
	c.conn.Write(proto_dht_query)
	c.conn.Write([]byte(address))
}

func (c *Client) Bootstrap() error {
	c.conn.Write(proto_bootstrap)

	length_b := make([]byte, 8)
	err := net_recvall(length_b, c.conn)

	if err != nil {
		c.Close()
		return err
	}

	length, _ := binary.Uvarint(length_b)

	if length > EntryLengthMax*BucketSize {
		c.Close()
		return errors.New("Peer sent too much data")
	}

	closest_json := make([]byte, length)
	net_recvall(closest_json, c.conn)
	log.Debug(string(closest_json))

	return nil
}
