
package zif

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Command input types

type CommandPeer struct {
    Address string `json:"address"`
}

type CommandPing     CommandPeer
type CommandAnnounce CommandPeer
type CommandRSearch struct {
	CommandPeer
    Query   string `json:"query"`
    Page    int    `json:"page"`
}
type CommandPeerSearch CommandRSearch
type CommandPeerRecent struct {
	CommandPeer
    Page int `json:"page"`
}
type CommandPeerPopular CommandPeerRecent
type CommandMirror      CommandPeer
type CommandPeerIndex struct {
	CommandPeer
    Since int `json:"since"`
}

type CommandMeta struct {
    PId int    `json:"pid"`
    Key string `json:"key"`
}

type CommandAddPost Post
type CommandSelfIndex struct {
    Since int `json:"since"`
}
type CommandResolve   CommandPeer
type CommandBootstrap CommandPeer
type CommandSelfSearch struct {
    Query string `json:"query"`
    Page  int    `json:"page"`
}
type CommandSelfRecent struct {
    Page int `json:"page"`
}
type CommandSelfPopular CommandSelfRecent
type CommandAddMeta struct {
	CommandMeta
    Value string `json:"value"`
}
type CommandGetMeta CommandMeta
type CommandSaveCollection    interface{}
type CommandRebuildCollection interface{}
type CommandPeers             interface{}
type CommandSaveRoutingTable  interface{}

// Command output types

type CommandResult struct {
    IsOK   bool        `json:"status"`
    Result interface{} `json:"value"`
    Error  error       `json:"err"`
}

func (cr *CommandResult) WriteJSON(w io.Writer) {
	e := json.NewEncoder(w)

	if cr.IsOK {
		if cr.Result == nil {
			e.Encode(struct{
				Status string `json:"status"`
			}{"ok"})
		} else {
			e.Encode(struct{
				Status string      `json:"status"`
				Value  interface{} `json:"value"`
			}{"ok", cr.Result})
		}
	} else {
		if cr.Error == nil {
			cr.Error = errors.New("Something bad happened, but we don't know bad, which makes the fact much worse.")
		}

		e.Encode(struct{
			Status string `json:"status"`
			Error  string `json:"err"`
		}{"err", cr.Error.Error()})
	}
}

// Command server type

type CommandServer struct {
	LocalPeer *LocalPeer
}

// Command functions

func (cs *CommandServer) Ping(p CommandPing) CommandResult {
	log.Info("Command: Ping request")

	// TODO: implement
	return CommandResult{true,nil,nil}
}
func (cs *CommandServer) Announce(a CommandAnnounce) CommandResult {
	var err error

	log.Info("Command: Announce request")

	peer := cs.LocalPeer.GetPeer(a.Address)

	if peer == nil {
		peer, err = cs.LocalPeer.ConnectPeer(a.Address)

		if err != nil {
			return CommandResult{false,nil,err}
		}
	}

	_, err = peer.ConnectClient(cs.LocalPeer)
	if err != nil {
		return CommandResult{false,nil,err}
	}

	err = peer.Announce(cs.LocalPeer)

	return CommandResult{err != nil,nil,err}
}
func (cs *CommandServer) RSearch(rs CommandRSearch) CommandResult {
	var err error

	log.Info("Command: Peer Remote Search request");

    peer := cs.LocalPeer.GetPeer(rs.CommandPeer.Address)

	if peer == nil {
		peer, err = cs.LocalPeer.ConnectPeer(rs.CommandPeer.Address)
		if err != nil {
			return CommandResult{false,nil,err}
		}
	}

	posts, stream, err := peer.Search(rs.Query, rs.Page)

	if stream != nil {
		defer stream.Close()
	}

	return CommandResult{err != nil,posts,err}
}
func (cs *CommandServer) PeerSearch(ps CommandPeerSearch) CommandResult {
	var err error

	log.Info("Command: Peer Search request")

	if !cs.LocalPeer.Databases.Has(ps.CommandPeer.Address) {
		return cs.RSearch(CommandRSearch{ps.CommandPeer,ps.Query,ps.Page})
	}

	db, _ := cs.LocalPeer.Databases.Get(ps.CommandPeer.Address)

	posts, err := db.(*Database).Search(ps.Query, ps.Page)

	return CommandResult{err != nil,posts,err}
}
func (cs *CommandServer) PeerRecent(pr CommandPeerRecent) CommandResult {
	var err   error
	var posts []*Post

	log.Info("Command: Peer Recent request")

	if pr.CommandPeer.Address == cs.LocalPeer.Entry.ZifAddress.Encode() {
		posts, err = cs.LocalPeer.Database.QueryRecent(pr.Page)

		return CommandResult{err != nil,posts,err}
	}

	peer := cs.LocalPeer.GetPeer(pr.CommandPeer.Address)
	if peer == nil {
		peer, err = cs.LocalPeer.ConnectPeer(pr.CommandPeer.Address)
		if err != nil {
			return CommandResult{false,nil,err}
		}
	}

	posts, stream, err := peer.Recent(pr.Page)

	if stream != nil {
		defer stream.Close()
	}

	return CommandResult{err != nil,posts,err}
}
func (cs *CommandServer) PeerPopular(pp CommandPeerPopular) CommandResult {
	var err   error
	var posts []*Post

	log.Info("Command: Peer Popular request")

	if pp.CommandPeer.Address == cs.LocalPeer.Entry.ZifAddress.Encode() {
		posts, err = cs.LocalPeer.Database.QueryPopular(pp.Page)

		return CommandResult{err == nil,posts,err}
	}

	peer := cs.LocalPeer.GetPeer(pp.CommandPeer.Address)
	if peer == nil {
		peer, err = cs.LocalPeer.ConnectPeer(pp.CommandPeer.Address)
		if err != nil {
			return CommandResult{false,nil,err}
		}
	}

	posts, stream, err := peer.Popular(pp.Page)

	if stream != nil {
		defer stream.Close()
	}

	return CommandResult{err == nil,posts,err}
}
func (cs *CommandServer) Mirror(cm CommandMirror) CommandResult {
	var err error

	log.Info("Command: Peer Mirror request")

	peer := cs.LocalPeer.GetPeer(cm.Address)
	if peer == nil {
		peer, err = cs.LocalPeer.ConnectPeer(cm.Address)
		if err != nil {
			return CommandResult{false,nil,err}
		}
	}

	// TODO: make this configurable
	d := fmt.Sprintf("./data/%s", peer.ZifAddress.Encode())
	os.Mkdir(fmt.Sprintf("./data/%s", d), 0777)
	db := NewDatabase(d)
	db.Connect()

	cs.LocalPeer.Databases.Set(peer.ZifAddress.Encode(), db)

	_, err = peer.Mirror(db)
	if err != nil {
		return CommandResult{false,nil,err}
	}

	// TODO: wjh: is this needed? -poro
	cs.LocalPeer.Databases.Set(peer.ZifAddress.Encode(), db)

	return CommandResult{true,nil,nil}
}
func (cs *CommandServer) PeerIndex(ci CommandPeerIndex) CommandResult {
	var err error

	log.Info("Command: Peer Index request")

	if !cs.LocalPeer.Databases.Has(ci.CommandPeer.Address) {
		return CommandResult{false,nil,errors.New("Peer database not loaded.")}
	}

	db, _ := cs.LocalPeer.Databases.Get(ci.CommandPeer.Address)
	err = db.(*Database).GenerateFts(ci.Since)

	return CommandResult{err == nil,nil,err}
}

// self

func (cs *CommandServer) AddPost(ap CommandAddPost) CommandResult {
	log.Info("Command: Add Post request")

	cs.LocalPeer.AddPost(
		Post{ap.Id,ap.InfoHash,ap.Title,ap.Size,ap.FileCount,ap.Seeders,ap.Leechers,ap.UploadDate,ap.Tags},
		false)

	return CommandResult{true,nil,nil}
}
func (cs *CommandServer) SelfIndex(ci CommandSelfIndex) CommandResult {
	log.Info("Command: FTS Index request")

	err := cs.LocalPeer.Database.GenerateFts(ci.Since)

	return CommandResult{err == nil, nil, err}
}
func (cs *CommandServer) Bootstrap(cb CommandBootstrap) CommandResult {
	log.Info("Command: Bootstrap request")

	addrnport := strings.Split(cb.Address, ":")

	host := addrnport[0]
	var port string
	if len(addrnport) == 1 {
		port = "5050" // TODO: make this configurable
	} else {
		port = addrnport[1]
	}

	peer, err := cs.LocalPeer.ConnectPeerDirect(host + ":" + port)
	if err != nil {
		return CommandResult{false,nil,err}
	}

	peer.ConnectClient(cs.LocalPeer)

	_, err = peer.Bootstrap(cs.LocalPeer.RoutingTable)

	return CommandResult{err == nil,nil,err}
}
func (cs *CommandServer) SelfSearch(css CommandSelfSearch) CommandResult {
	log.Info("Command: Search request")

	posts, err := cs.LocalPeer.Database.Search(css.Query, css.Page)

	return CommandResult{err == nil,posts,err}
}
func (cs *CommandServer) SelfRecent(cr CommandSelfRecent) CommandResult {
	log.Info("Command: Recent request")

	posts, err := cs.LocalPeer.Database.QueryRecent(cr.Page)

	return CommandResult{err == nil,posts,err}
}
func (cs *CommandServer) SelfPopular(cp CommandSelfPopular) CommandResult {
	log.Info("Command: Popular request")

	posts, err := cs.LocalPeer.Database.QueryPopular(cp.Page)

	return CommandResult{err == nil,posts,err}
}
func (cs *CommandServer) AddMeta(cam CommandAddMeta) CommandResult {
	log.Info("Command: Add Meta request")

	err := cs.LocalPeer.Database.AddMeta(cam.CommandMeta.PId, cam.CommandMeta.Key, cam.Value)

	return CommandResult{err == nil,nil,err}
}
func (cs *CommandServer) GetMeta(cgm CommandGetMeta) CommandResult {
	log.Info("Command: Get Meta request")

	val, err := cs.LocalPeer.Database.GetMeta(cgm.PId, cgm.Key)

	return CommandResult{err == nil,val,err}
}
func (cs *CommandServer) SaveCollection(csc CommandSaveCollection) CommandResult {
	log.Info("Command: Save Collection request")

	// TODO: make this configurable
	cs.LocalPeer.Collection.Save("./data/collection.dat")

	return CommandResult{true,nil,nil}
}
func (cs *CommandServer) RebuildCollection(crc CommandRebuildCollection) CommandResult {
	var err error

	log.Info("Command: Rebuild Collection request")

	cs.LocalPeer.Collection, err = CreateCollection(cs.LocalPeer.Database, 0, PieceSize)
	return CommandResult{err != nil,nil,err}
}
func (cs *CommandServer) Peers(cp CommandPeers) CommandResult {
	log.Info("Command: Peers request")

	ps := make([]*Peer, cs.LocalPeer.Peers.Count() + 1)

	ps[0] = &cs.LocalPeer.Peer

	i := 1
	for _, p := range cs.LocalPeer.Peers.Items() {
		ps[i] = p.(*Peer)
		i = i + 1
	}

	return CommandResult{true,ps,nil}
}
func (cs *CommandServer) SaveRoutingTable(csrt CommandSaveRoutingTable) CommandResult {
	log.Info("Command: Save Routing Table request")

	// TODO: make this configurable
	err := cs.LocalPeer.RoutingTable.Save("dht")

	return CommandResult{err != nil,nil,err}
}

