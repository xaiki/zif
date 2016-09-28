package zif

import (
	"encoding/json"
	"errors"
	"net"
)

const (
	TitleMax    = 144
	TagsMax     = 256
	MaxPostSize = TitleMax + TagsMax + 1024
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

func NewPost(ih, title string, seeders, leechers, uploaddate int, source []byte) Post {
	var p Post

	p.InfoHash = ih
	p.Title = title
	p.Seeders = seeders
	p.Leechers = leechers
	p.UploadDate = uploaddate
	copy(p.Source[:], source)

	return p
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
