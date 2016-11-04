package zif

import (
	"errors"
	"hash"
	"io/ioutil"
	"math"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

// a collection of pieces, by extension a structure containing all posts this
// peer has.

type Collection struct {
	Pieces   []*Piece
	hashlist []byte
	hash     hash.Hash
}

func NewCollection() *Collection {
	col := &Collection{}

	col.hash = sha3.New256()
	col.Pieces = make([]*Piece, 0, 2)
	col.hashlist = make([]byte, 0)

	return col
}

// Takes a database, starting id, and piece size. Generates a hash list.
func CreateCollection(db *Database, start, pieceSize int) (*Collection, error) {
	col := NewCollection()

	postCount := db.PostCount()
	pieceCount := int(math.Ceil(float64(postCount) / float64(pieceSize)))

	log.Info("Piece count ", pieceCount)

	for i := 0; i < pieceCount; i++ {
		piece, err := db.QueryPiece(i, false)

		if err != nil {
			return nil, err
		}

		col.Add(piece)
	}

	return col, nil
}

// Loads a collection from file
// This essentially loads the hash list, the data of pieces themselves is just
// left. It's all in the database if it is really needed.
func LoadCollection(path string) (col *Collection, err error) {
	col = NewCollection()

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	if len(data)%32 != 0 {
		err = errors.New("Invalid collection data file")
		return
	}

	col.hashlist = data
	col.Rehash()

	return
}

func (c *Collection) Save(path string) {
	ioutil.WriteFile(path, c.hashlist, 0777)
}

func (c *Collection) Add(piece *Piece) {
	c.Pieces = append(c.Pieces, piece)
	c.hashlist = append(c.hashlist, piece.Hash()...)

	c.hash.Write(piece.Hash())
}

func (c *Collection) AddPost(post Post, store bool) {
	if len(c.Pieces) == 0 || len(c.Pieces[len(c.Pieces)-1].Posts) == PieceSize {
		piece := &Piece{}
		piece.Setup()
		c.Add(piece)

		c.hashlist = append(c.hashlist, piece.Hash()...)
	}

	lastIndex := len(c.Pieces) - 1
	last := c.Pieces[lastIndex]
	last.Add(post, store)

	copy(c.hashlist[lastIndex*32:lastIndex*32+32], last.Hash())
}

func (c *Collection) Hash() []byte {
	var ret []byte

	ret = c.hash.Sum(nil)

	return ret
}

func (c *Collection) Rehash() {
	c.hash = sha3.New256()

	for i := 0; i < len(c.hashlist)/32; i++ {
		c.hash.Write(c.hashlist[i*32 : i*32+i])
	}
}

func (c *Collection) HashList() []byte {
	return c.hashlist
}
