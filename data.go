package zif

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

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

// This is signed, *not* the JSON. This is needed because otherwise the order of
// the posts encoded is not actually guaranteed, which can lead to invalid
// signatures. Plus we can only sign data that is actually needed.
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

// Convert a post to a string, with an optional separator between fields, and
// an optional terminating value (appended to the end of the post string).
// This is *actually* signed, to allow for different encoders encoding in a
// different order. (relying on a json encoding to always encode the same way
// is not the best idea)
func PostToString(p *Post, sep, term string) string {
	buf := bytes.Buffer{}

	WritePost(p, sep, term, &buf)

	return buf.String()
}

func WritePost(p *Post, sep, term string, w io.Writer) {
	bw := bufio.NewWriter(w)

	bw.WriteString(strconv.Itoa(p.Id) + sep)
	bw.WriteString(p.InfoHash + sep)
	bw.WriteString(p.Title + sep)
	bw.WriteString(strconv.Itoa(p.Size) + sep)
	bw.WriteString(strconv.Itoa(p.FileCount) + sep)
	bw.WriteString(strconv.Itoa(p.Seeders) + sep)
	bw.WriteString(strconv.Itoa(p.Leechers) + sep)
	bw.WriteString(strconv.Itoa(p.UploadDate) + sep)
	bw.WriteString(p.Tags + sep)
	bw.WriteString(term)
}

func ReadPost(r io.Reader, delim byte) {
	br := bufio.NewReader(r)

	br.ReadString(delim)
}

// Ensures that all the members of an entry struct fit the requirements for the
// Zif protocol. If an entry passes this, then we should be able to perform
// most operations on it.
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
