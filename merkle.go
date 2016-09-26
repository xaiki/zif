// a merkle tree
// look into securing against second preimage attacks

package zif

import (
	"errors"

	"golang.org/x/crypto/sha3"
)

// what is being stored!
type MerkleEntry struct {
	Hash []byte
	Id   uint
}

type MerkleNode struct {
	MerkleEntry
	Left   *MerkleNode
	Right  *MerkleNode
	Parent *MerkleNode

	dirty bool
}

func NewMerkleNode(hash []byte) *MerkleNode {
	mn := MerkleNode{}

	if hash == nil {
		hash = make([]byte, 32)
	}
	mn.Hash = hash
	mn.dirty = true

	return &mn
}

func (mn *MerkleNode) Update() {
	// If this is a leaf that has had it's hash changed (more likely just been
	// added :))

	data := make([]byte, 0, 64)

	if mn.Left != nil {
		if mn.Left.dirty {
			mn.Left.Update()
		}
		mn.Left.dirty = false
		data = append(data, mn.Left.Hash...)
	}

	if mn.Right != nil {
		if mn.Right.dirty {
			mn.Right.Update()
		}
		mn.Right.dirty = false
		data = append(data, mn.Right.Hash...)
	}

	// then this is a leaf!
	if mn.Left == nil && mn.Right == nil && mn.Hash != nil {
		return
	}

	hash := sha3.Sum256(data)

	if len(mn.Hash) != len(hash) {
		mn.Hash = make([]byte, 32)
	}

	copy(mn.Hash, hash[:])
}

func (mn *MerkleNode) Insert(entry MerkleEntry) error {
	if entry.Hash == nil {
		return errors.New("Invalid entry")
	}

	//finish this bit and write a test!

	return nil
}

func (mn *MerkleNode) InsertLeft(left *MerkleNode) {
	mn.Left = left
	mn.Left.Parent = mn

	mn.dirty = true
}

func (mn *MerkleNode) InsertRight(right *MerkleNode) {
	mn.Right = right
	mn.Right.Parent = mn

	mn.dirty = true
}
