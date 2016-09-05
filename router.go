package main

import "bytes"
import log "github.com/sirupsen/logrus"

// routes and handles tcp messages

func RouteMessage(msg_type []byte, peer Peer, lp *LocalPeer) {
	if bytes.Equal(msg_type, proto_ping) {
		peer.Pong()
	} else if bytes.Equal(msg_type, proto_pong) {
		log.Debug("Pong from ", peer.ZifAddress.Encode())
	} else if bytes.Equal(msg_type, proto_who) {
		peer.SendWho()
	} else if bytes.Equal(msg_type, proto_dht_announce) {
		peer.RecievedAnnounce()
	}
}
