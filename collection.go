package zif

import (
	"errors"
	"hash"
	"io/ioutil"
	"math"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

// A collection of pieces, by extension a structure containing all posts this
// peer has. Whether or not the pieces are *actually* there is optional, if not
// this is essentially a hash list.
type Collection struct {
	Pieces   []*Piece
	hashlist []byte
	hash     hash.Hash
}

// Create a new collection, set all it's members to the correct default values.
func NewCollection() *Collection {
	col := &Collection{}

	col.hash = sha3.New256()
	col.Pieces = make([]*Piece, 0, 2)
	col.hashlist = make([]byte, 0)

	return col
}

// Takes a database, starting id, and piece size. From this we create a
// collection, except it does not contain any posts - consider making this optional.
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

// Loads a collection from file.
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

// Save the collection hash list to the given path, with permissions 0777.
func (c *Collection) Save(path string) {
	ioutil.WriteFile(path, c.hashlist, 0777)
}

// Add a piece to the collection, storing it in c.Pieces and appending it's hash
// to the hash list.
func (c *Collection) Add(piece *Piece) {
	c.Pieces = append(c.Pieces, piece)
	c.hashlist = append(c.hashlist, piece.Hash()...)

	c.hash.Write(piece.Hash())
}

// Add a post to the collection. Automatically assigns to the correct piece,
// allocates a new piece if needed! Optional whether or not the actual post is
// stored, if not just it's hash is.
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

// Return the hash of the hash list, which can then go on to be signed by the
// LocalPeer. This allows proper validation of an entire collection, but the
// localpeer only needs to sign a single hash.
func (c *Collection) Hash() []byte {
	var ret []byte

	ret = c.hash.Sum(nil)

	return ret
}

// Regenerates the root hash from the hash list we have.
func (c *Collection) Rehash() {
	c.hash = sha3.New256()

	for i := 0; i < len(c.hashlist)/32; i++ {
		c.hash.Write(c.hashlist[i*32 : i*32+32])
	}
}

// Simply return the hash list.
func (c *Collection) HashList() []byte {
	return c.hashlist
}
