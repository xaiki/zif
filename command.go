
package zif

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Command input types

type CommandPeer struct {
	address string
}

type CommandPing     CommandPeer
type CommandAnnounce CommandPeer
type CommandRSearch struct {
	CommandPeer
	query   string
	page    int
}
type CommandPeerSearch CommandRSearch
type CommandPeerRecent struct {
	CommandPeer
	page int
}
type CommandPeerPopular CommandPeerRecent
type CommandMirror      CommandPeer
type CommandPeerIndex struct {
	CommandPeer
	since int
}

type CommandMeta struct {
	pid int
	key string
}

type CommandAddPost Post
type CommandSelfIndex struct {
	since int
}
type CommandResolve   CommandPeer
type CommandBootstrap CommandPeer
type CommandSelfSearch struct {
	query string
	page  int
}
type CommandSelfRecent struct {
	page int
}
type CommandSelfPopular CommandSelfRecent
type CommandAddMeta struct {
	CommandMeta
	value string
}
type CommandGetMeta CommandMeta
type CommandSaveCollection    interface{}
type CommandRebuildCollection interface{}
type CommandPeers             interface{}
type CommandSaveRoutingTable  interface{}

// Command output types

type CommandResult struct {
	isOk   bool
	result interface{}
	errmsg error
}

func (cr *CommandResult) WriteJSON(w io.Writer) {
	e := json.NewEncoder(w)

	if cr.isOk {
		if cr.result == nil
			e.Encode(struct{
				Status string `json:"status"`
			}{"ok"})
		else
			e.Encode(struct{
				Status string	  `json:"status"`
				Value  interface{} `json:"value"`
			}{"ok", cr.result})
	}
	else {
		if cr.errmsg == nil
			cr.errmsg = errors.New("An unspecified error occured.")

		e.Encode(struct{
			Status string `json:"status"`
			Error  string `json:"err"`
		}{"err", cr.errmsg.Error()})
	}
}

// Command server type

type CommandServer struct {
	localPeer *LocalPeer
}

// Command functions

func (cs *CommandServer) Ping(p CommandPing) CommandResult {
	log.Info("Command: Ping request")

	// TODO: implement
	return {true,nil,nil}
}
func (cs *CommandServer) Announce(a CommandAnnounce) CommandResult {
	var err error

	log.Info("Command: Announce request")

	peer := cs.localPeer.GetPeer(a.address)

	if peer == nil {
		peer, err = cs.localPeer.ConnectPeer(a.address)

		if err != nil
			return {false,nil,err}
	}

	_, err = peer.ConnectClient(cs.localPeer)
	if err != nil
		return {false,nil,err}

	err = peer.Announce(cs.localPeer)

	return {err != nil,nil,err}
}
func (cs *CommandServer) RSearch(rs CommandRSearch) CommandResult {
	var err error

	log.Info("Command: Peer Remote Search request");

	peer = cs.localPeer.GetPeer(rs.CommandPeer.address)

	if peer == nil {
		peer, err = cs.localPeer.ConnectPeer(rs.CommandPeer.address)
		if err != nil
			return {false,nil,err}
	}

	posts, stream, err := peer.Search(rs.query, rs.page)

	if stream != nil
		defer stream.Close()

	return {err != nil,posts,err}
}
func (cs *CommandServer) PeerSearch(ps CommandPeerSearch) CommandResult {
	var err error

	log.Info("Command: Peer Search request")

	if !cs.localPeer.Databases.Has(ps.CommandPeer.address)
		return cs.RSearch(ps.(CommandRSearch)) // this is completely fine

	db, _ := cs.localPeer.Databases.Get(ps.CommandPeer.address)

	posts, err := db.(*Database).Search(ps.query, ps.page)

	return {err != nil,posts,err}
}
func (cs *CommandServer) PeerRecent(pr CommandPeerRecent) CommandResult {
	var err   error
	var posts []*Post

	log.Info("Command: Peer Recent request")

	if pr.CommandPeer.address == cs.localPeer.Entry.ZifAddress.Encode() {
		posts, err = cs.localPeer.Database.Query

		return {err != nil,posts,err}
	}

	peer := cs.localPeer.GetPeer(pr.CommandPeer.address)
	if peer == null {
		peer, err = cs.localPeer.ConnectPeer(pr.CommandPeer.address)
		if err != nil
			return {false,nil,err}
	}

	posts, stream, err := peer.Recent(pr.page)

	if stream != nil
		defer stream.Close()

	return {err != nil,posts,err}
}
func (cs *CommandServer) PeerPopular(pp CommandPeerPopular) CommandResult {
	var err   error
	var posts []*Post

	log.Info("Command: Peer Popular request")

	if pp.CommandPeer.address == cs.localPeer.Entry.ZifAddress.Encode() {
		posts, err = cs.localPeer.Database.Query

		return {err == nil,posts,err}
	}

	peer := cs.localPeer.GetPeer(pp.CommandPeer.address)
	if peer == null {
		peer, err = cs.localPeer.ConnectPeer(pp.CommandPeer.address)
		if err != nil
			return {false,nil,err}
	}

	posts, stream, err := peer.Popular(pp.page)

	if stream != nil
		defer stream.Close()

	return {err == nil,posts,err}
}
func (cs *CommandServer) Mirror(cm CommandMirror) CommandResult {
	var err error

	log.Info("Command: Peer Mirror request")

	peer := cs.localPeer.GetPeer(cm.address)
	if peer == nil {
		peer, err = cs.localPeer.ConnectPeer(cm.address)
		if err != nil
			return {false,nil,err}
	}

	// TODO: make this configurable
	d := fmt.Sprintf("./data/%s", peer.ZifAddress.Encode())
	os.Mkdir(fmt.Sprintf("./data/%s", d, 0777)
	db := NewDatabase(d)
	db.Connect()

	cs.localPeer.Databases.Set(peer.ZifAddress.Encode(), db)

	_, err = peer.Mirror(db)
	if err != nil
		return {false,nil,err}

	// TODO: wjh: is this needed? -poro
	cs.localPeer.Databases.Set(peer.ZifAddress.Encode(), db)

	return {true,nil,nil}
}
func (cs *CommandServer) PeerIndex(ci CommandPeerIndex) CommandResult {
	var err error

	log.Info("Command: Peer Index request")

	if !cs.localPeer.Database.Has(ci.CommandPeer.address)
		return {false,nil,errors.New("Peer database not loaded.")}

	db, _ := cs.localPeer.Databases.Get(ci.CommandPeer.address)
	err = db.(*Database).GenerateFts(ci.since)

	return {err == nil,nil,err}
}

// self

func (cs *CommandServer) AddPost(cp CommandAddPost) CommandResult {
	log.Info("Command: Add Post request")

	cs.localPeer.AddPost(cp, false)

	return {true,nil,nil}
}
func (cs *CommandServer) SelfIndex(ci CommandSelfIndex) CommandResult {
	log.Info("Command: FTS Index request")

	err := cs.LocalPeer.Database.GenerateFts(ci.since)

	return {err == nil, nil, err}
}
func (cs *CommandServer) Bootstrap(cb CommandBootstrap) CommandResult {
	log.Info("Command: Bootstrap request")

	addrnport = strings.Split(cb.address, ":")

	address := addrnport[0]
	var port string
	if len(address) == 1
		port = "5050" // TODO: make this configurable
	else
		port = addrnport[1]

	peer, err := cs.localPeer.ConnectPeerDirect(host + ":" + port)
	if err != nil
		return {false,nil,err}

	peer.ConnectClient(cs.localPeer)

	_, err = peer.Bootstrap(cs.localPeer.RoutingTable)

	return {err == nil,nil,err}
}
func (cs *CommandServer) SelfSearch(cs CommandSelfSearch) CommandResult {
	log.Info("Command: Search request")

	posts, err := cs.localPeer.Database.Search(cs.query, cs.page)

	return {err == nil,posts,err}
}
func (cs *CommandServer) SelfRecent(cr CommandSelfRecent) CommandResult {
	log.Info("Command: Recent request")

	posts, err := cs.localPeer.Database.QueryRecent(cs.page)

	return {err == nil,posts,err}
}
func (cs *CommandServer) SelfPopular(cp CommandSelfPopular) CommandResult {
	log.Info("Command: Popular request")

	posts, err := cs.localPeer.Database.QueryPopular(cp.page)

	return {err == nil,posts,err}
}
func (cs *CommandServer) AddMeta(cam CommandAddMeta) CommandResult {
	log.Info("Command: Add Meta request")

	err := cs.localPeer.Database.AddMeta(cam.CommandMeta.pid, cam.CommandMeta.key, cam.value)

	return {err == nil,nil,err}
}
func (cs *CommandServer) GetMeta(cgm CommandGetMeta) CommandResult {
	log.Info("Command: Get Meta request")

	val, err := cs.localPeer.Database.GetMeta(cgm.pid, cgm.key)

	return {err == nil,nil,err}
}
func (cs *CommandServer) SaveCollection(csc CommandSaveCollection) CommandResult {
	log.Info("Command: Save Collection request")

	// TODO: make this configurable
	cs.localPeer.Collection.Save("./data/collection.dat")

	return {true,nil,nil}
}
func (cs *CommandServer) RebuildCollection(crc CommandRebuildCollection) CommandResult {
	var err error

	log.Info("Command: Rebuild Collection request")

	cs.localPeer.Collection, err = CreateCollection(cs.localPeer.Database, 0, PieceSize)
	return {err =! nil,nil,err}
}
func (cs *CommandServer) Peers(cp CommandPeers) CommandResult {
	log.Info("Command: Peers request")

	ps := make([]*Peer, cs.localPeer.Peers.Count() + 1)

	ps[0] = &cs.localPeer.Peer

	i := 1
	for _, p := range cs.localPeer.Peers.Items() {
		ps[i] = p.(*Peer)
		i = i + 1
	}

	return {true,ps,nil}
}
func (cs *CommandServer) SaveRoutingTable(csrt CommandSaveRoutingTable) CommandResult {
	log.Info("Command: Save Routing Table request")

	// TODO: make this configurable
	err := cs.localPeer.RoutingTable.Save("dht")

	return {err != nil,nil,err}
}

