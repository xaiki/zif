package zif

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
)

func EntriesToJson(entries []*Entry) ([]byte, error) {
	data, err := json.Marshal(entries)

	return data, err
}

func EntryToJson(e *Entry) ([]byte, error) {
	data, err := json.Marshal(e)

	return data, err
}

func PostsToJson(posts []*Post) ([]byte, error) {
	data, err := json.Marshal(posts)

	return data, err
}

func JsonToEntry(data []byte) (Entry, error) {
	var e Entry
	err := json.Unmarshal(data, &e)

	return e, err
}

// This is signed, *not* the JSON.
func EntryToBytes(e *Entry) []byte {
	var str string

	str += e.Name
	str += e.Desc
	str += string(e.PublicKey)
	str += string(e.Port)
	str += string(e.PublicAddress)
	str += string(e.ZifAddress.Encode())
	str += string(e.PostCount)

	return []byte(str)
}

func ValidateEntry(entry *Entry) error {
	if len(entry.PublicKey) < ed25519.PublicKeySize {
		return errors.New(fmt.Sprintf("Public key too small: %d", len(entry.PublicKey)))
	}

	if len(entry.Signature) < ed25519.SignatureSize {
		return errors.New("Signature too small")
	}

	verified := ed25519.Verify(entry.PublicKey, EntryToBytes(entry), entry.Signature[:])

	if !verified {
		return errors.New("Failed to verify signature")
	}

	if len(entry.PublicAddress) == 0 {
		return errors.New("Public address must be set")
	}

	// 253 is the maximum length of a domain name
	if len(entry.PublicAddress) >= 253 {
		return errors.New("Public address is too large (253 char max)")
	}

	if entry.Port > 65535 {
		return errors.New("Port too large (" + string(entry.Port) + ")")
	}

	return nil
}

// Takes a database, starting id, and piece size. Generates a hash list.
func CreateCollection(db *Database, start, pieceSize int) (*Collection, error) {
	col := Collection{}
	col.Setup()

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

	return &col, nil
}
