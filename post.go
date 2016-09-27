package zif

import (
	"encoding/json"
	"errors"
	"net"
)

const (
	TitleMax = 144
	TagsMax  = 256
)

type Post struct {
	Id         int
	InfoHash   string
	Title      string
	Size       int
	FileCount  int
	Seeders    int
	Leechers   int
	UploadDate int
	Source     []byte
	Tags       string
}

func (p Post) Json() ([]byte, error) {
	json, err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	return json, nil
}

func (p *Post) NetSend(conn net.Conn) error {
	return errors.New("not implemented")
}
