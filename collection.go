package zif

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

// a collection of pieces, by extension a structure containing all posts this
// peer has.

type Collection struct {
	Pieces []*Piece
	hash   hash.Hash
}

func (c *Collection) Setup() {
	c.hash = sha3.New256()
	c.Pieces = make([]*Piece, 0, 2)
}

func (c *Collection) Add(piece *Piece) {
	c.Pieces = append(c.Pieces, piece)

	c.hash.Write(piece.Hash())
}

func (c *Collection) AddPost(post Post, store bool) {
	if len(c.Pieces) == 0 || len(c.Pieces[len(c.Pieces)-1].Posts) == PieceSize {
		piece := &Piece{}
		piece.Setup()
		c.Add(piece)
	}

	c.Pieces[len(c.Pieces)-1].Add(post, store)
}

func (c *Collection) Hash() []byte {
	var ret []byte

	ret = c.hash.Sum(nil)

	return ret
}

func (c *Collection) HashList() []byte {
	hash_list := make([]byte, 0, len(c.Pieces))

	for _, h := range c.Pieces {
		hash_list = append(hash_list, h.Hash()...)
	}

	return hash_list
}
