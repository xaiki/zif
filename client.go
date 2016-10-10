package zif

import (
	"compress/gzip"
	"encoding/json"
	"net"
	"time"
)

/*import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	EntryLengthMax = 1024
	MaxPageSize    = 25
)*/

type Client struct {
	conn net.Conn
}

func NewClient(stream net.Conn) Client {
	return Client{conn}
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
	w := gzip.NewWriter(c.conn)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return w.Close()
}

func (c *Client) ReadMessage() (*Message, error) {
	r, err := gzip.NewReader(c.conn)
	if err != nil {
		return nil, err
	}
	r.Multistream(false)
	var msg Message
	if err := json.NewDecoder(r).Decode(&msg); err != nil {
		return nil, err
	}
	if err := r.Close(); err != nil {
		return nil, err
	}
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
	/*json, err := EntryToJson(e)

	if err != nil {
		c.conn.Close()
		return err
	}

	length := len(json)
	length_b := make([]byte, 8)
	binary.PutVarint(length_b, int64(length))

	c.conn.Write(length_b)
	c.conn.Write(json)*/

	return nil
}

func (c *Client) Announce(e *Entry) error {
	/*c.conn.Write(proto_dht_announce)
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
	}*/

	return nil
}

func (c *Client) Query(address string) ([]Entry, error) {
	/*c.conn.Write(proto_dht_query)

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
	}*/

	return nil, nil
}

func (c *Client) Bootstrap(rt *RoutingTable, address Address) error {
	/*peers, err := c.Query(address.Encode())

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
	}*/

	return nil
}

// TODO: Paginate searches
func (c *Client) Search(search string) ([]*Post, error) {
	/*log.Info("Querying for ", search)

	c.conn.Write(proto_search)
	err := net_sendlength(c.conn, uint64(len(search)))

	if err != nil {
		return nil, err
	}

	ok := make([]byte, 2)

	if !bytes.Equal(proto_ok, ok) {
		return nil, errors.New("Peer did not accept search")
	}

	net_recvall(ok, c.conn)

	c.conn.Write([]byte(search))

	length, err := net_recvlength(c.conn)

	log.Info("Query returned ", length, " results")

	posts := make([]*Post, 0, length)

	for i := uint64(0); i < length; i++ {
		post, err := net_recvpost(c.conn)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}*/

	return nil, nil
}

func (c *Client) Recent(page uint64) ([]*Post, error) {
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
