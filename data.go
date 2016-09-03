package main

import "encoding/json"

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
