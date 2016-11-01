package zif

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/ed25519"

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
func (c *Client) Search(search string) ([]*Post, error) {
	log.Info("Querying for ", search)

	msg := &Message{
		Header:  ProtoSearch,
		Content: []byte(search),
	}

	c.WriteMessage(msg)

	var posts []*Post

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

func (c *Client) Recent(page int) ([]*Post, error) {
	log.Info("Fetching recent posts from peer")

	page_s := strconv.Itoa(page)

	msg := &Message{
		Header:  ProtoRecent,
		Content: []byte(page_s),
	}

	err := c.WriteMessage(msg)

	if err != nil {
		return nil, err
	}

	posts_msg, err := c.ReadMessage()

	if err != nil {
		return nil, err
	}

	var posts []*Post
	posts_msg.Decode(&posts)

	log.Info("Recieved ", len(posts), " recent posts")

	return posts, nil
}

func (c *Client) Popular(page int) ([]*Post, error) {
	log.Info("Fetching popular posts from peer")

	page_s := strconv.Itoa(page)

	msg := &Message{
		Header:  ProtoPopular,
		Content: []byte(page_s),
	}

	err := c.WriteMessage(msg)

	if err != nil {
		return nil, err
	}

	posts_msg, err := c.ReadMessage()

	if err != nil {
		return nil, err
	}

	var posts []*Post
	posts_msg.Decode(&posts)

	log.Info("Recieved ", len(posts), " popular posts")

	return posts, nil
}

// Download a hash list for a peer. Expects said hash list to be valid and
// signed.
func (c *Client) HashList(address Address, pk ed25519.PublicKey) ([]byte, error) {
	log.WithField("for", address.Encode()).Info("Sending request for a collection")

	msg := &Message{
		Header:  ProtoRequestHashList,
		Content: address.Bytes,
	}

	c.WriteMessage(msg)

	hl, err := c.ReadMessage()

	if err != nil {
		return nil, err
	}

	mhl, err := MessageHashListDecode(hl.Content)

	if err != nil {
		return nil, err
	}

	err = mhl.Verify(pk)

	if err != nil {
		return nil, err
	}

	return mhl.HashList, nil
}
