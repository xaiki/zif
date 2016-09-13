package main

import (
	"bytes"
	"net"

	log "github.com/sirupsen/logrus"
)

// routes and handles tcp messages

func RouteMessage(msg_type []byte, peer *Peer, stream net.Conn) {
	//log.Debug("Routing message ", msg_type)

	if bytes.Equal(msg_type, proto_ping) {
		peer.Pong()
	} else if bytes.Equal(msg_type, proto_pong) {
		log.Debug("Pong from ", peer.ZifAddress.Encode())
	} else if bytes.Equal(msg_type, proto_dht_announce) {
		peer.RecievedAnnounce(stream, peer)
	} else if bytes.Equal(msg_type, proto_dht_query) {
		peer.RecieveQuery(stream)
	} else if bytes.Equal(msg_type, proto_bootstrap) {
		peer.RecieveBootstrap(stream)
	}
}
