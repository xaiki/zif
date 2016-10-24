package zif

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	EntryLengthMax = 1024
	MaxPageSize    = 25
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

func (c *Client) Terminate() {
	//c.conn.Write(proto_terminate)
}

func (c *Client) Close() (err error) {
	if c.conn != nil {
		err = c.conn.Close()
	}
	return
}

func (c *Client) WriteMessage(v interface{}) error {
	err := json.NewEncoder(c.conn).Encode(v)

	return err
}

func (c *Client) ReadMessage() (*Message, error) {
	var msg Message

	if err := json.NewDecoder(c.conn).Decode(&msg); err != nil {
		return nil, err
	}

	msg.Stream = c.conn

	return &msg, nil
}

func (c *Client) Ping(timeout time.Duration) bool {
	/*c.conn.Write(proto_ping)

	tchan := make(chan bool)

	go func() {
		buf := make([]byte, 2)
		net_recvall(buf, c.conn)

		tchan <- true
	}()

	select {
	case <-tchan:
		return true
	case <-time.After(timeout):
		return false
	}*/
	return true
}

func (c *Client) Pong() {
	//c.conn.Write(proto_pong)
}

func (c *Client) SendEntry(e *Entry) error {
	json, err := EntryToJson(e)
	msg := Message{Header: ProtoEntry, Content: json}

	if err != nil {
		c.conn.Close()
		return err
	}

	c.WriteMessage(msg)

	return nil
}

func (c *Client) Announce(e *Entry) error {
	json, err := EntryToJson(e)

	if err != nil {
		c.conn.Close()
		return err
	}

	msg := &Message{
		Header:  ProtoDhtAnnounce,
		Content: json,
	}

	err = c.WriteMessage(msg)

	if err != nil {
		return err
	}

	ok, err := c.ReadMessage()

	if err != nil {
		return err
	}

	if !ok.Ok() {
		return errors.New("Peer did not respond with ok")
	}

	return nil
}

func (c *Client) Query(address string) ([]Entry, error) {
	// TODO: LimitReader

	msg := &Message{
		Header:  ProtoDhtQuery,
		Content: []byte(address),
	}

	// Tell the peer the address we are looking for
	err := c.WriteMessage(msg)

	if err != nil {
		return nil, err
	}

	// Make sure the peer accepts the address
	recv, err := c.ReadMessage()

	if err != nil {
		return nil, err
	}

	if !recv.Ok() {
		return nil, errors.New("Peer refused query address")
	}

	closest, err := c.ReadMessage()

	if err != nil {
		return nil, err
	}

	var entries []Entry
	err = closest.Decode(&entries)

	log.WithField("entries", len(entries)).Info("Query complete")
	return entries, err
}

func (c *Client) Bootstrap(rt *RoutingTable, address Address) error {
	peers, err := c.Query(address.Encode())

	if err != nil {
		return err
	}

	// add them all to our routing table! :D
	for _, e := range peers {
		if len(e.ZifAddress.Bytes) != AddressBinarySize {
			continue
		}
		rt.Update(e)
	}

	if len(peers) > 1 {
		log.Info("Bootstrapped with ", len(peers), " new peers")
	} else if len(peers) == 1 {
		log.Info("Bootstrapped with 1 new peer")
	}

	return nil
}

// TODO: Paginate searches
func (c *Client) Search(search string) ([]Post, error) {
	log.Info("Querying for ", search)

	msg := &Message{
		Header:  ProtoSearch,
		Content: []byte(search),
	}

	c.WriteMessage(msg)

	var posts []Post

	recv, err := c.ReadMessage()

	if err != nil {
		return nil, err
	}

	err = recv.Decode(&posts)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (c *Client) Recent(page uint64) ([]Post, error) {
	/*log.Info("Fetching recent posts from peer")

	c.conn.Write(proto_recent)
	err := net_sendlength(c.conn, page)

	if err != nil {
		return nil, err
	}

	length, err := net_recvlength(c.conn)

	if err != nil {
		return nil, err
	}

	if length > MaxPageSize {
		return nil, errors.New("Peer returned a page that was too large")
	}

	posts := make([]*Post, 0, length)

	for i := uint64(0); i < length; i++ {
		post, err := net_recvpost(c.conn)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	log.Info("Recieved ", len(posts), " recent posts")*/

	return nil, nil
}
