package zif

import (
	"encoding/json"
	"errors"
	"net"

	log "github.com/sirupsen/logrus"
)

const MaxSearchLength = 256

// TODO: Move this into some sort of handler object, can handle general requests.

// TODO: While I think about it, move all these TODOs to issues or a separate
// file/issue tracker or something.

// Querying peer sends a Zif address
// This peer will respond with a list of the k closest peers, ordered by distance.
// The top peer may well be the one that is being queried for :)
func (lp *LocalPeer) HandleQuery(stream net.Conn, from *Peer) error {
	from.limiter.queryLimiter.Wait()

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

func (lp *LocalPeer) HandleAnnounce(stream net.Conn, from *Peer) {
	/*from.limiter.announceLimiter.Wait()
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
	}*/

}

func (lp *LocalPeer) HandleSearch(conn net.Conn, from *Peer) {
	length, err := net_recvlength(conn)

	if err != nil {
		log.Debug(err.Error())
		return
	}

	if length > MaxSearchLength {
		log.Debug("Query too long")
		return
	}

	//conn.Write(proto_ok)

	buf := make([]byte, length)
	net_recvall(buf, conn)

	query := string(buf)
	log.Info("Post query for ", query)

	posts, err := lp.Database.Search(query, 0)

	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info(len(posts), " results")

	net_sendlength(conn, uint64(len(posts)))

	for _, p := range posts {
		net_sendpost(conn, p)
	}
}

func (lp *LocalPeer) HandleRecent(conn net.Conn, from *Peer) {
	log.Info("Recieved query for recent posts")
	page, err := net_recvlength(conn)

	if err != nil {
		log.Debug(err.Error())
		return
	}

	posts, err := lp.Database.QueryRecent(int(page))
	net_sendlength(conn, uint64(len(posts)))

	for _, p := range posts {
		net_sendpost(conn, p)
	}
}

func (lp *LocalPeer) ListenStream(peer *Peer) {
	lp.Server.ListenStream(peer)
}
