package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Ping() {
	c.conn.Write(proto_ping)
}

func (c *Client) Pong() {
	c.conn.Write(proto_pong)
}

func (c *Client) SendEntry(e *Entry) error {
	json, err := EntryToJson(e)

	if err != nil {
		c.conn.Close()
		return err
	}

	length := len(json)
	length_b := make([]byte, 8)
	binary.PutVarint(length_b, int64(length))

	c.conn.Write(length_b)
	c.conn.Write(json)

	return nil
}

func (c *Client) Announce(e *Entry) error {
	c.conn.Write(proto_dht_announce)
	err := c.SendEntry(e)

	if err != nil {
		return err
	}

	buf := make([]byte, 2)
	err = net_recvall(buf, c.conn)

	if err != nil {
		return err
	}

	if !bytes.Equal(buf, proto_ok) {
		return errors.New("Peer did not ok announce")
	}

	return nil
}

func (c *Client) Query(address string) ([]Entry, error) {
	c.conn.Write(proto_dht_query)

	net_sendlength(c.conn, uint64(len(address)))
	c.conn.Write([]byte(address))

	length, err := net_recvlength(c.conn)

	if length > EntryLengthMax*BucketSize {
		c.Close()
		return nil, errors.New("Peer sent too much data")
	}

	closest_json := make([]byte, length)
	net_recvall(closest_json, c.conn)

	var closest []Entry
	err = json.Unmarshal(closest_json, &closest)

	if err != nil {
		return nil, err
	}

	return closest, nil
}

func (c *Client) Bootstrap(rt *RoutingTable, address Address) error {
	peers, err := c.Query(address.Encode())

	if err != nil {
		return err
	}

	// add them all to our routing table! :D
	for _, e := range peers {
		rt.Update(e)
	}

	if len(peers) > 1 {
		log.Info("Bootstrapped with ", len(peers), " new peers")
	} else if len(peers) == 1 {
		log.Info("Bootstrapped with 1 new peer")
	}

	return nil
}
