// Stores things like message codes, etc.

package main

var (
	// Protocol header, so we know this is a zif client.
	// Version should follow.
	proto_zif     = []byte{0x7a, 0x66}
	proto_version = []byte{0x00, 0x00} //version 0 atm :D

	// inform a peer on the status of the latest request
	proto_ok        = []byte{0x00, 0x00}
	proto_no        = []byte{0x00, 0x01}
	proto_terminate = []byte{0x00, 0x02}

	proto_msg_latest = []byte{0x01, 0x00}

	proto_ping = []byte{0x02, 0x00}
	proto_pong = []byte{0x02, 0x01}

	proto_dht_query    = []byte{0x03, 0x00}
	proto_dht_announce = []byte{0x03, 0x01}
)
