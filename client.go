package zif

import (
	"compress/gzip"
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

	decoder *json.Decoder
	encoder *json.Encoder
}

// Creates a new client, automatically setting up the json encoder/decoder.
func NewClient(conn net.Conn) *Client {
	return &Client{conn, json.NewDecoder(conn), json.NewEncoder(conn)}
}

func (c *Client) Terminate() {
	//c.conn.Write(proto_terminate)
}

// Close the client connection.
func (c *Client) Close() (err error) {
	if c.conn != nil {
		err = c.conn.Close()
	}
	return
}

// Encodes v as json and writes it to c.conn.
func (c *Client) WriteMessage(v interface{}) error {
	if c.encoder == nil {
		c.encoder = json.NewEncoder(c.conn)
	}

	err := c.encoder.Encode(v)

	return err
}

// Blocks until a message is read from c.conn, decodes it into a *Message and
// returns.
func (c *Client) ReadMessage() (*Message, error) {
	var msg Message

	if c.decoder == nil {
		c.decoder = json.NewDecoder(c.conn)
	}

	if err := c.decoder.Decode(&msg); err != nil {
		return nil, err
	}

	msg.Stream = c.conn

	return &msg, nil
}

// Pings a client with a specified timeout, returns true/false depending on
// if it recieves a reply.
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

// Replies to a Ping request.
func (c *Client) Pong() {
	//c.conn.Write(proto_pong)
}

// Sends a DHT entry to a peer.
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

// Announce the given DHT entry to a peer, passes on this peers details,
// meaning that it can be reached by other peers on the network.
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

// Given a Zif address, attempts to resolve it for a DHT entry. Returns the k
// closest peers to the address. It only returns the closest entries that the
// peer knows about, so more Queries may well be needed.
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

// Adds the initial entries into the given routing table. Essentially queries for
// both it's own and the peers address, storing the result. This means that after
// a bootstrap, it should be possible to connect to *any* peer!
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
func (c *Client) Search(search string, page int) ([]*Post, error) {
	log.Info("Querying for ", search)

	sq := MessageSearchQuery{search, page}
	data, err := sq.Encode()

	if err != nil {
		return nil, err
	}

	msg := &Message{
		Header:  ProtoSearch,
		Content: data,
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
func (c *Client) Collection(address Address, pk ed25519.PublicKey) (*MessageCollection, error) {
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

	mhl := MessageCollection{}
	err = hl.Decode(&mhl)

	if err != nil {
		return nil, err
	}

	err = mhl.Verify(pk)

	if err != nil {
		return nil, err
	}

	log.WithField("pieces", mhl.Size).Info("Recieved valid collection")

	return &mhl, nil
}

// Download a piece from a peer, given the address and id of the piece we want.
func (c *Client) Pieces(address Address, id, length int) chan *Piece {
	log.WithFields(log.Fields{
		"address": address.Encode(),
		"id":      id,
		"length":  length,
	}).Info("Sending request for piece")

	ret := make(chan *Piece, 100)

	mrp := MessageRequestPiece{address.Encode(), id, length}
	data, err := mrp.Encode()

	if err != nil {
		return nil
	}

	msg := &Message{
		Header:  ProtoRequestPiece,
		Content: data,
	}

	c.WriteMessage(msg)

	// Convert a string to an int, prevents endless error checks below.
	convert := func(val string) int {
		var ret int
		ret, err := strconv.Atoi(val)

		if err != nil {
			return 0
		}

		return ret
	}

	go func() {
		defer close(ret)

		gzr, err := gzip.NewReader(c.conn)
		if err != nil {
			return
		}

		errReader := NewErrorReader(gzr)

		for i := 0; i < length; i++ {
			piece := Piece{}
			piece.Setup()

			count := 0
			for {
				if count >= PieceSize {
					break
				}

				id := convert(errReader.ReadString('|'))

				if id == -1 {
					break
				}

				ih := errReader.ReadString('|')
				title := errReader.ReadString('|')
				size := convert(errReader.ReadString('|'))
				filecount := convert(errReader.ReadString('|'))
				seeders := convert(errReader.ReadString('|'))
				leechers := convert(errReader.ReadString('|'))
				uploaddate := convert(errReader.ReadString('|'))
				tags := errReader.ReadString('|')

				if errReader.err != nil {
					log.Error("Failed to read post: ", errReader.err.Error())
					break
				}

				if err != nil {
					log.Error(err.Error())
				}

				post := Post{
					Id:         id,
					InfoHash:   ih,
					Title:      title,
					Size:       size,
					FileCount:  filecount,
					Seeders:    seeders,
					Leechers:   leechers,
					UploadDate: uploaddate,
					Tags:       tags,
				}

				piece.Add(post, true)
				count++
			}
			ret <- &piece
		}
	}()

	return ret
}
