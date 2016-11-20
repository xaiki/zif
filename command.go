
package zif

import (
	"encoding/json"
	"fmt"
	"os"

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

type CommandAddPost struct {
	// TODO
}
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
	if err != nil
		return {false,nil,err}

	return {true,nil,nil}
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

	if err != nil
		return {false,nil,err}

	return {true,posts,nil}
}
func (cs *CommandServer) PeerSearch(ps CommandPeerSearch) CommandResult {
	var err error

	log.Info("Command: Peer Search request")

	if !cs.localPeer.Databases.Has(ps.CommandPeer.address)
		return cs.RSearch(ps.(CommandRSearch)) // this is completely fine

	db, _ := cs.localPeer.Databases.Get(ps.CommandPeer.address)

	posts, err := db.(*Database).Search(ps.query, ps.page)
	if err != nil
		return {false,nil,err}

	return {true,posts,nil}
}
func (cs *CommandServer) PeerRecent(pr CommandPeerRecent) CommandResult {
	var err   error
	var posts []*Post

	log.Info("Command: Peer Recent request")

	if pr.CommandPeer.address == cs.localPeer.Entry.ZifAddress.Encode() {
		posts, err = cs.localPeer.Database.Query
		if err != nil
			return {false,nil,err}

		return {true,posts,nil}
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

	if err != nil
		return {false,nil,err}

	return {true,posts,nil}
}
func (cs *CommandServer) PeerPopular(pp CommandPeerPopular) CommandResult {
	var err   error
	var posts []*Post

	log.Info("Command: Peer Popular request")

	if pp.CommandPeer.address == cs.localPeer.Entry.ZifAddress.Encode() {
		posts, err = cs.localPeer.Database.Query
		if err != nil
			return {false,nil,err}

		return {true,posts,nil}
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

	if err != nil
		return {false,nil,err}

	return {true,posts,nil}
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
	if err != nil
		return {false,nil,err}

	return {true,nil,nil}
}

