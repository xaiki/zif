package zif

import (
	"errors"
	"hash"

	log "github.com/sirupsen/logrus"
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

func (p *Piece) Add(post Post, store bool) error {
	if len(p.Posts) > PieceSize {
		return errors.New("Piece full")
	}

	if store {
		p.Posts = append(p.Posts, post)
	}

	data := PostToString(&post, "|", "")
	p.hash.Write([]byte(data))

	return nil
}

func (p *Piece) Hash() []byte {
	var ret []byte

	ret = p.hash.Sum(nil)

	return ret
}

func (p *Piece) Rehash() ([]byte, error) {
	p.hash = sha3.New256()

	for _, i := range p.Posts {
		data := PostToString(&i, "|", "")
		p.hash.Write([]byte(data))
	}

	log.Info("Piece rehashed")

	return p.hash.Sum(nil), nil
}
