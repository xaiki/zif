package zif

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/sha3"
)

func GetHash(data []byte) []byte {
	hash := sha3.Sum256((data))
	return hash[:]
}

func HashMismatchCheck(hash1, hash2 []byte, t *testing.T) {
	if !bytes.Equal(hash1, hash2) {
		t.Error("Hash mismatch")
	}
}

func TestMerkleUpdate(t *testing.T) {
	root := NewMerkleNode(nil)

	root.InsertLeft(NewMerkleNode([]byte("abc")))
	root.InsertRight(NewMerkleNode([]byte("cba")))

	root.Update()

	HashMismatchCheck(root.Hash, GetHash([]byte("abccba")), t)

	root.Left.InsertLeft(NewMerkleNode([]byte("abc")))
	root.Update()

	left_hash := GetHash([]byte("abc"))
	root_hash := GetHash(append(left_hash, []byte("cba")...))

	HashMismatchCheck(left_hash, root.Left.Hash, t)
	HashMismatchCheck(root_hash, root.Hash, t)

	root.Left.InsertRight(NewMerkleNode([]byte("cab")))

	left_hash = GetHash([]byte("abccab"))
	root.Update()

	HashMismatchCheck(left_hash, root.Left.Hash, t)

	root.Right.InsertRight(NewMerkleNode([]byte("cba")))
	root.Update()

	HashMismatchCheck(root.Right.Hash, GetHash([]byte("cba")), t)

	root.Right.InsertLeft(NewMerkleNode([]byte("bca")))

	right_hash := GetHash([]byte("bcacba"))
	root_hash = GetHash(append(left_hash, right_hash...))
	root.Update()

	HashMismatchCheck(right_hash, root.Right.Hash, t)
	HashMismatchCheck(root_hash, root.Hash, t)
}
