package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// TODO: While I think about it, move all these TODOs to issues or a separate
// file/issue tracker or something.

// Querying peer sends a Zif address
// This peer will respond with a list of the k closest peers, ordered by distance.
// The top peer may well be the one that is being queried for :)
func (lp *LocalPeer) HandleQuery(stream net.Conn) error {
	address_length, err := net_recvlength(stream)

	if err != nil {
		return err
	}

	address_bin := make([]byte, address_length)
	err = net_recvall(address_bin, stream)

	if err != nil {
		return err
	}

	address := DecodeAddress(string(address_bin))
	log.Info("Recieved query for ", address.Encode())

	closest := lp.RoutingTable.FindClosest(address, BucketSize)

	closest_json, err := json.Marshal(closest)

	if err != nil {
		return errors.New("Failed to encode closest peers to json")
	}

	net_sendlength(stream, uint64(len(closest_json)))
	stream.Write(closest_json)

	return nil
}

// TODO: Rate limit this to prevent announce flooding.
func (lp *LocalPeer) HandleAnnounce(stream net.Conn, from *Peer) {
	log.Debug("Recieved announce")
	lp.CheckSessions()

	defer stream.Close()

	entry, err := recieve_entry(stream)

	if err != nil {
		log.Error(err.Error())
		return
	}

	var addr Address
	addr.Generate(entry.PublicKey[:])

	log.Debug("Announce from ", from.ZifAddress.Encode())

	saved := lp.RoutingTable.Update(entry)

	if saved {
		log.Info("Saved new peer ", addr.Encode())
	}

	// next up, tell other people!
	closest := lp.RoutingTable.FindClosest(addr, BucketSize)

	// TODO: Parallize this
	for _, i := range closest {
		if i.ZifAddress.Equals(&entry.ZifAddress) || i.ZifAddress.Equals(&from.ZifAddress) {
			continue
		}

		peer := lp.GetPeer(i.ZifAddress.Encode())

		if peer == nil {
			log.Debug("Connecting to new peer")
			peer = NewPeer(lp)
			err = peer.Connect(i.PublicAddress + ":" + strconv.Itoa(i.Port))

			if err != nil {
				log.Warn("Failed to connect to peer: ", err.Error())
				continue
			}

			peer.ConnectClient()
		}

		peer_stream, err := peer.OpenStream()
		defer peer_stream.Close()

		if err != nil {
			log.Error(err.Error())
			continue
		}

		peer_stream.conn.Write(proto_dht_announce)
		peer_stream.SendEntry(&entry)
	}

}

func (lp *LocalPeer) HandleStream(peer *Peer, stream net.Conn) {
	log.Debug("Handling stream")
	msg := make([]byte, 2)
	for {
		err := net_recvall(msg, stream)

		if err != nil {
			if err.Error() == "EOF" {
				log.Info("Closed stream from ", peer.ZifAddress.Encode())
			} else {
				log.Error(err.Error())
			}

			peer.RemoveStream(stream)

			return
		}

		if bytes.Equal(msg, proto_terminate) {
			peer.Terminate()
			log.Debug("Terminated connection with ", peer.ZifAddress.Encode())
			return
		}

		lp.RouteMessage(msg, peer, stream)
	}
}
