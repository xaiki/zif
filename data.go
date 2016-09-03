package main

import "encoding/json"
import "golang.org/x/crypto/ed25519"

func EntryToJson(e *Entry) ([]byte, error) {
	data, err := json.Marshal(e)

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

	return []byte(str)
}

func ValidateEntry(entry *Entry, sig []byte) bool {
	verified := ed25519.Verify(entry.PublicKey, EntryToBytes(entry), sig)

	if !verified {
		return false
	}

	// 253 is the maximum length of a domain name
	return len(entry.PublicAddress) > 0 && len(entry.PublicAddress) < 253 &&
		entry.Port < 65535

}
