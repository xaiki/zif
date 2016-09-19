package main

import (
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

	if len(closest) < 1 {
		return nil, errors.New("Query returned no results")
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
