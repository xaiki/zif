package libzif

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"
	data "github.com/wjh/zif/libzif/data"
	"github.com/wjh/zif/libzif/dht"
	"github.com/wjh/zif/libzif/proto"
)

const MaxSearchLength = 256

// TODO: Move this into some sort of handler object, can handle general requests.

// TODO: While I think about it, move all these TODOs to issues or a separate
// file/issue tracker or something.

// Querying peer sends a Zif address
// This peer will respond with a list of the k closest peers, ordered by distance.
// The top peer may well be the one that is being queried for :)
func (lp *LocalPeer) HandleQuery(msg *proto.Message) error {
	log.Info("Handling query")
	cl := Client{msg.Stream, nil, nil}

	//msg.From.limiter.queryLimiter.Wait()

	address := dht.DecodeAddress(string(msg.Content))
	log.WithField("target", address.Encode()).Info("Recieved query")

	ok := &Message{Header: ProtoOk}
	err := cl.WriteMessage(ok)

	if err != nil {
		return err
	}

	closest_json := bytes.Buffer{}
	encoder := json.NewEncoder(&closest_json)

	if address.Equals(&lp.Address) {
		log.WithField("name", lp.Entry.Name).Debug("Query for local peer")

		json, err := EntryToJson(&lp.Entry)

		if err != nil {
			return err
		}
		kv := dht.NewKeyValue(lp.Entry.Address, json)

		encoder.Encode(kv)
	} else {
		log.Debug("Querying routing table")

		for _, i := range lp.RoutingTable.FindClosest(address, dht.MaxBucketSize) {
			encoder.Encode(i)
		}
	}

	log.Debug("Query results: ", string(closest_json.String()))

	if err != nil {
		return errors.New("Failed to encode closest peers to json")
	}

	results := &Message{
		Header:  ProtoEntry,
		Content: closest_json.Bytes(),
	}

	err = cl.WriteMessage(results)

	return err
}

func (lp *LocalPeer) HandleAnnounce(msg *proto.Message) error {
	cl := Client{msg.Stream, nil, nil}
	msg.From.limiter.announceLimiter.Wait()

	defer msg.Stream.Close()

	entry := Entry{}
	err := msg.Decode(&entry)

	log.WithField("address", entry.Address.Encode()).Info("Announce")

	if err != nil {
		return err
	}

	json := EntryToBytes(&entry)
	saved := lp.RoutingTable.Update(dht.NewKeyValue(entry.Address, json))

	if saved {
		cl.WriteMessage(&Message{Header: ProtoOk})
		log.WithField("peer", entry.Address.Encode()).Info("Saved new peer")

	} else {
		cl.WriteMessage(&Message{Header: ProtoNo})
		return errors.New("Failed to save entry")
	}

	// next up, tell other people!
	closest := lp.RoutingTable.FindClosest(entry.Address, dht.MaxBucketSize)

	// TODO: Parallize this
	for _, i := range closest {
		if i.Key.Equals(&entry.Address) || i.Key.Equals(&msg.From.Address) {
			continue
		}

		peer := lp.GetPeer(entry.Address.Encode())

		if peer == nil {
			log.Debug("Connecting to new peer")

			decoded, err := JsonToEntry(i.Value)

			if err != nil {
				return err
			}

			var p Peer
			err = p.Connect(decoded.PublicAddress+":"+strconv.Itoa(decoded.Port), lp)

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

		peer_announce := &Message{
			Header:  ProtoDhtAnnounce,
			Content: msg.Content,
		}
		peer_stream.WriteMessage(peer_announce)
	}
	return nil

}

func (lp *LocalPeer) HandleSearch(msg *proto.Message) error {
	if len(msg.Content) > MaxSearchLength {
		return errors.New("Search query too long")
	}

	sq := MessageSearchQuery{}
	err := msg.Decode(&sq)

	if err != nil {
		return err
	}

	log.WithField("query", sq.Query).Info("Search recieved")

	posts, err := lp.Database.Search(sq.Query, sq.Page, 25)

	if err != nil {
		return err
	}
	log.Info("Posts loaded")

	json, err := PostsToJson(posts)

	if err != nil {
		return err
	}

	post_msg := &Message{
		Header:  ProtoPosts,
		Content: json,
	}

	msg.Client.WriteMessage(post_msg)

	return nil
}

func (lp *LocalPeer) HandleRecent(msg *proto.Message) error {
	log.Info("Recieved query for recent posts")

	page, err := strconv.Atoi(string(msg.Content))

	if err != nil {
		return err
	}

	recent, err := lp.Database.QueryRecent(page)

	if err != nil {
		return err
	}

	recent_json, err := PostsToJson(recent)

	if err != nil {
		return err
	}

	resp := &Message{
		Header:  ProtoPosts,
		Content: recent_json,
	}

	msg.Client.WriteMessage(resp)

	return nil
}

func (lp *LocalPeer) HandlePopular(msg *proto.Message) error {
	log.Info("Recieved query for popular posts")

	page, err := strconv.Atoi(string(msg.Content))

	if err != nil {
		return err
	}

	recent, err := lp.Database.QueryPopular(page)

	if err != nil {
		return err
	}

	recent_json, err := PostsToJson(recent)

	if err != nil {
		return err
	}

	resp := &Message{
		Header:  ProtoPosts,
		Content: recent_json,
	}

	msg.Client.WriteMessage(resp)

	return nil
}

func (lp *LocalPeer) HandleHashList(msg *proto.Message) error {
	address := dht.Address{msg.Content}

	log.WithField("address", address.Encode()).Info("Collection request recieved")

	var sig []byte

	if address.Equals(&lp.Address) {
		sig = lp.Sign(lp.Collection.HashList)
	} else {
		// this means that the hash list wanted does not belong to this peer
		// TODO: sort out getting a hash list for a peer that has been mirrored
	}

	mhl := MessageCollection{lp.Collection.Hash(), lp.Collection.HashList, len(lp.Collection.HashList) / 32, sig}
	data, err := mhl.Encode()

	if err != nil {
		return err
	}

	resp := &Message{
		Header:  ProtoHashList,
		Content: data,
	}

	msg.Client.WriteMessage(resp)

	return nil
}

func (lp *LocalPeer) HandlePiece(msg *proto.Message) error {

	mrp := MessageRequestPiece{}
	err := msg.Decode(&mrp)

	log.WithFields(log.Fields{
		"id":     mrp.Id,
		"length": mrp.Length,
	}).Info("Recieved piece request")

	if err != nil {
		return err
	}

	var posts chan *data.Post

	if mrp.Address == lp.Entry.Address.Encode() {
		posts = lp.Database.QueryPiecePosts(mrp.Id, mrp.Length, true)

	} else if lp.Databases.Has(mrp.Address) {
		db, _ := lp.Databases.Get(mrp.Address)
		posts = db.(*data.Database).QueryPiecePosts(mrp.Id, mrp.Length, true)

	} else {
		return errors.New("Piece not found")
	}

	// Buffered writer -> gzip -> net
	// or
	// gzip -> buffered writer -> net
	// I'm guessing the latter allows for gzip to maybe run a little faster?
	// The former may allow for database reads to occur a little faster though.
	// buffer both?
	bw := bufio.NewWriter(msg.Stream)
	gzw := gzip.NewWriter(bw)

	for i := range posts {
		data.WritePost(i, "|", "", gzw)
	}
	data.WritePost(&data.Post{Id: -1}, "|", "", gzw)

	gzw.Flush()
	bw.Flush()

	log.Info("Sent all")

	return nil
}

func (lp *LocalPeer) HandleAddPeer(msg *proto.Message) error {
	// The AddPeer message contains the address of the peer that the client
	// wishes to be registered for.

	peerFor := string(msg.Content)

	log.Info("Handling add peer request for ", peerFor)

	// First up, we need the address in binary form
	address := dht.DecodeAddress(peerFor)

	if len(address.Bytes) != dht.AddressBinarySize {
		msg.Client.WriteMessage(&Message{Header: ProtoNo})
		return errors.New("Invalid binary address size")
	}

	if address.Equals(&lp.Address) {
		log.WithField("peer", address.Encode()).Info("New seed peer")

		add := true

		for _, i := range lp.Entry.Seeds {
			if address.Equals(&dht.Address{i}) {
				add = false
			}
		}

		if add {
			lp.Entry.Seeds = append(lp.Entry.Seeds, address.Bytes)
		}

	} else {
		// then we need to see if we have the entry for that address
		results := lp.RoutingTable.FindClosest(address, dht.MaxBucketSize)

		if len(results) == 0 {
			msg.Client.WriteMessage(&Message{Header: ProtoNo})
			return errors.New("Could not resolve address")
		}

		decoded, err := JsonToEntry(results[0].Value)

		if err != nil {
			return err
		}

		// if the routing table contains the address we are looking for,
		// register a new seed.
		if decoded.Address.Equals(&address) {
			decoded.Seeds = append(decoded.Seeds, address.Bytes)
		}
	}

	msg.Client.WriteMessage(&Message{Header: ProtoOk})

	return nil
}

func (lp *LocalPeer) HandlePing(msg *proto.Message) error {
	log.WithField("peer", msg.From.Address.Encode()).Info("Ping")

	return msg.Client.WriteMessage(&Message{Header: ProtoPong})
}

func (lp *LocalPeer) ListenStream(peer *Peer) {
	lp.Server.ListenStream(peer)
}
