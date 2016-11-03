package zif

import (
	"encoding/json"
	"errors"

	"golang.org/x/crypto/ed25519"
)

// This contains the more "complex" structures that will be sent in message
// data fields.

type MessageCollection struct {
	Hash      []byte
	HashList  []byte
	Size      int
	Signature []byte
}

type MessageSearchQuery struct {
	Query string
	Page  int
}

type MessageRequestPiece struct {
	Address string
	Id      int
}

func (mhl *MessageCollection) Verify(pk ed25519.PublicKey) error {
	// TODO: Check length of pk/hashlist/etc

	verified := ed25519.Verify(pk, mhl.HashList, mhl.Signature)

	if !verified {
		return errors.New("Invalid signature")
	}

	return nil
}

func (mhl *MessageCollection) Encode() ([]byte, error) {
	data, err := json.Marshal(mhl)
	return data, err
}

func (sq *MessageSearchQuery) Encode() ([]byte, error) {
	data, err := json.Marshal(sq)
	return data, err
}

func (mrp *MessageRequestPiece) Encode() ([]byte, error) {
	data, err := json.Marshal(mrp)
	return data, err
}
