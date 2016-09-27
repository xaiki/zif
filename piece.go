package zif

import (
	"errors"
	"hash"

	"golang.org/x/crypto/sha3"
)

const PieceSize = 1000

type Piece struct {
	Posts []Post
	hash  hash.Hash
}

func (p *Piece) Setup() {
	p.hash = sha3.New256()
}

func (p *Piece) Add(post Post) error {
	if len(p.Posts) > PieceSize {
		return errors.New("Piece full")
	}

	p.Posts = append(p.Posts, post)

	json, err := post.Json()

	if err != nil {
		return err
	}

	p.hash.Write(json)

	return nil
}

func (p *Piece) Hash() []byte {
	var ret []byte

	ret = p.hash.Sum(nil)

	return ret
}
