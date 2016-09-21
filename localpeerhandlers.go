package zif

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// TODO: Move this into some sort of handler object, can handle general requests.

// TODO: While I think about it, move all these TODOs to issues or a separate
// file/issue tracker or something.

// Querying peer sends a Zif address
// This peer will respond with a list of the k closest peers, ordered by distance.
// The top peer may well be the one that is being queried for :)
func (lp *LocalPeer) HandleQuery(stream net.Conn) error {
	log.Debug(lp.Entry.Desc)

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
	log.WithField("target", address.Encode()).Info("Recieved query")

	var closest []*Entry

	if address.Equals(&lp.ZifAddress) {
		log.Debug("Query for local peer")
		closest = make([]*Entry, 1)
		closest[0] = &lp.Entry
	} else {
		log.Debug("Querying routing table")
		closest = lp.RoutingTable.FindClosest(address, BucketSize)
	}

	closest_json, err := json.Marshal(closest)
	log.Debug("Query results: ", string(closest_json))

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
		stream.Write(proto_ok)
		log.WithField("peer", entry.ZifAddress.Encode()).Info("Saved new peer")

	} else {
		stream.Write(proto_no)
		stream.Close()
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

			var p Peer
			err = p.Connect(i.PublicAddress+":"+strconv.Itoa(i.Port), lp)

			if err != nil {
				log.Warn("Failed to connect to peer: ", err.Error())
				continue
			}

			p.ConnectClient(lp)

			peer = &p
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

func (lp *LocalPeer) ListenStream(peer *Peer) {
	lp.Server.ListenStream(peer)
}
