package zif

import (
	"encoding/json"
	"net"
)

type Message struct {
	Header  int
	Content []byte

	From   *Peer
	Stream net.Conn
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
