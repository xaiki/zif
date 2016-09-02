package main

import "bytes"

// routes and handles tcp messages

func RouteMessage(msg_type []byte, peer Peer) {
	if bytes.Equal(msg_type, proto_ping) {
		peer.Pong()
	}
}
