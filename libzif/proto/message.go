package proto

import (
	"encoding/json"
	"net"
)

type Message struct {
	Header  int
	Content []byte

	Stream net.Conn
	Client *Client
}

func (m *Message) Json() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Decode(iface interface{}) error {
	err := json.Unmarshal(m.Content, iface)

	return err
}

// Ok() is just an easier way to check if the peer has sent an "ok" response,
// rather than comparing the header member to a constant repeatedly.
func (m *Message) Ok() bool {
	return m.Header == ProtoOk
}
