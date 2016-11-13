// Used to control the Zif daemon

package zif

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type HTTPServer struct {
	LocalPeer *LocalPeer
}

func (hs *HTTPServer) ListenHTTP(addr string) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", hs.IndexHandler)

	// This should be the ONLY route where the address is a non-Zif address

	router.HandleFunc("/peer/{address}/ping/", hs.Ping)
	router.HandleFunc("/peer/{address}/announce/", hs.Announce)
	router.HandleFunc("/peer/{address}/rsearch/", hs.PeerRSearch).Methods("POST")
	router.HandleFunc("/peer/{address}/search/", hs.PeerSearch).Methods("POST")
	router.HandleFunc("/peer/{address}/recent/{page}/", hs.Recent)
	router.HandleFunc("/peer/{address}/popular/{page}/", hs.Popular)
	router.HandleFunc("/peer/{address}/mirror/", hs.Mirror)
	router.HandleFunc("/peer/{address}/index/{since}/", hs.PeerFtsIndex)

	router.HandleFunc("/self/addpost/", hs.AddPost).Methods("POST")
	router.HandleFunc("/self/index/{since}/", hs.FtsIndex)
	router.HandleFunc("/self/resolve/{address}/", hs.Resolve)
	router.HandleFunc("/self/bootstrap/{address}/", hs.Bootstrap)
	router.HandleFunc("/self/search/", hs.SelfSearch).Methods("POST")
	router.HandleFunc("/self/recent/{page}/", hs.SelfRecent)
	router.HandleFunc("/self/popular/{page}/", hs.SelfPopular)
	router.HandleFunc("/self/addmeta/{pid}/{key}/{value}/", hs.AddMeta)
	router.HandleFunc("/self/getmeta/{pid}/{key}/", hs.GetMeta)
	router.HandleFunc("/self/savecollection/", hs.SaveCollection)
	router.HandleFunc("/self/rebuildcollection/", hs.RebuildCollection)
	router.HandleFunc("/self/peers/", hs.Peers)

	log.Info("Starting HTTP server on ", addr)

	err := http.ListenAndServe(addr, router)

	if err != nil {
		panic(err)
	}
}

func http_error_check(w http.ResponseWriter, errCode int, err error) bool {
	if err != nil {
		http_write_error(w, errCode, err)

		return true
	}

	return false
}

func http_write_error(w http.ResponseWriter, errCode int, err error) {
	w.WriteHeader(errCode)
	w.Write([]byte("{ \"status\": \"err\", \"err\": \"" + err.Error() + "\"}"))
}

func http_write_ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"ok\" }"))
}

// writes a single string value (eg, for metadata gets)
func http_write_value(w http.ResponseWriter, val string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"ok\", \"value\":\"" + val + "\" }"))
}

func http_write_posts(w http.ResponseWriter, posts []*Post) {
	json, err := json.Marshal(posts)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	// TODO: Use/write some sort of json building, based on a map.
	// This is kinda gross.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{ \"status\": \"ok\", \"value\": " + string(json) + "}"))
}

// TODO: wjh: actually implement these (and no, no "shhh") -poro
func (hs *HTTPServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http_write_ok(w)
}

func (hs *HTTPServer) Ping(w http.ResponseWriter, r *http.Request) {
	http_write_ok(w)
}

func (hs *HTTPServer) Bootstrap(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Bootstrap request")
	vars := mux.Vars(r)

	address := strings.Split(vars["address"], ":")
	host := address[0]
	var port string
	if len(address) == 1 {
		port = "5050"
	} else {
		port = address[1]
	}

	peer, err := hs.LocalPeer.ConnectPeerDirect(host + ":" + port)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	peer.ConnectClient(hs.LocalPeer)

	_, err = peer.Bootstrap(&hs.LocalPeer.RoutingTable)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}

func (hs *HTTPServer) Announce(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Announce request")
	vars := mux.Vars(r)

	peer, err := hs.LocalPeer.ConnectPeer(vars["address"])

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	_, err = peer.ConnectClient(hs.LocalPeer)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	err = peer.Announce(hs.LocalPeer)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}

func (hs *HTTPServer) Resolve(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Resolve request")

	vars := mux.Vars(r)
	addr := vars["address"]

	log.Info("Attempting to resolve address ", addr)

	entry, err := hs.LocalPeer.Resolve(addr)

	if http_error_check(w, http.StatusNotFound, err) {
		return
	}

	entry_json, err := EntryToJson(entry)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("{\"status\": \"ok\", \"value\":"))
	w.Write(entry_json)
	w.Write([]byte(" }"))
}

// Runs a remote search on a peer, ie, a search performed over a network connection.
func (hs *HTTPServer) PeerRSearch(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Peer Remote Search request")
	vars := mux.Vars(r)
	addr := vars["address"]

	query := r.FormValue("query")
	page := r.FormValue("page")

	page_i, err := strconv.Atoi(page)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	log.Info("Searching ", addr, " for ", query)

	peer, err := hs.LocalPeer.ConnectPeer(addr)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	posts, stream, err := peer.Search(query, page_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	defer stream.Close()

	http_write_posts(w, posts)
}

func (hs *HTTPServer) PeerSearch(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Peer Search request")
	vars := mux.Vars(r)
	addr := vars["address"]

	if !hs.LocalPeer.Databases.Has(addr) {
		hs.PeerRSearch(w, r)
		return
	}

	query := r.FormValue("query")
	page := r.FormValue("page")

	page_i, err := strconv.Atoi(page)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	log.Info("Searching ", addr, " for ", query)

	db, _ := hs.LocalPeer.Databases.Get(addr)

	posts, err := db.(*Database).Search(query, page_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_posts(w, posts)
}

func (hs *HTTPServer) AddPost(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Add Post request")

	post_json := r.FormValue("data")

	log.Debug("Adding post, json: ", post_json)

	var post Post

	err := json.Unmarshal([]byte(post_json), &post)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	hs.LocalPeer.AddPost(post, false)

	http_write_ok(w)
}

func (hs *HTTPServer) Recent(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Recent request")

	vars := mux.Vars(r)
	page := vars["page"]

	page_i, err := strconv.Atoi(page)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	var posts []*Post
	if vars["address"] == hs.LocalPeer.Entry.ZifAddress.Encode() {
		posts, err = hs.LocalPeer.Database.QueryRecent(page_i)

		if http_error_check(w, http.StatusInternalServerError, err) {
			return
		}

		http_write_posts(w, posts)
	}

	peer, err := hs.LocalPeer.ConnectPeer(vars["address"])

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	posts, stream, err := peer.Recent(page_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	defer stream.Close()

	http_write_posts(w, posts)
}

func (hs *HTTPServer) Popular(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Popular request")

	vars := mux.Vars(r)
	page := vars["page"]

	page_i, err := strconv.Atoi(page)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	var posts []*Post
	if vars["address"] == hs.LocalPeer.Entry.ZifAddress.Encode() {
		posts, err = hs.LocalPeer.Database.QueryPopular(page_i)

		if http_error_check(w, http.StatusInternalServerError, err) {
			return
		}

		http_write_posts(w, posts)
	}

	peer, err := hs.LocalPeer.ConnectPeer(vars["address"])

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	posts, stream, err := peer.Popular(page_i)
	defer stream.Close()

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_posts(w, posts)
}

func (hs *HTTPServer) FtsIndex(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: FTS Index request")

	vars := mux.Vars(r)
	since := vars["since"]

	since_i, err := strconv.Atoi(since)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	err = hs.LocalPeer.Database.GenerateFts(since_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}

func (hs *HTTPServer) PeerFtsIndex(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: FTS Index request")

	vars := mux.Vars(r)
	since := vars["since"]
	addr := vars["address"]

	since_i, err := strconv.Atoi(since)

	if !hs.LocalPeer.Databases.Has(addr) {
		err = errors.New("Peer database not loaded")
	}

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	db, _ := hs.LocalPeer.Databases.Get(addr)

	err = db.(*Database).GenerateFts(since_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}

func (hs *HTTPServer) SelfSearch(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Self Search request")

	query := r.FormValue("query")
	page := r.FormValue("page")

	log.Debug(query)
	log.Debug(page)

	page_i, err := strconv.Atoi(page)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	log.Info("Searching for ", query)

	posts, err := hs.LocalPeer.Database.Search(query, page_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_posts(w, posts)
}

func (hs *HTTPServer) SelfRecent(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Self Recent request")

	vars := mux.Vars(r)
	page := vars["page"]

	page_i, err := strconv.Atoi(page)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	posts, err := hs.LocalPeer.Database.QueryRecent(page_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_posts(w, posts)
}

func (hs *HTTPServer) SelfPopular(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Self Popular request")

	vars := mux.Vars(r)
	page := vars["page"]

	page_i, err := strconv.Atoi(page)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	posts, err := hs.LocalPeer.Database.QueryPopular(page_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_posts(w, posts)
}

func (hs *HTTPServer) AddMeta(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Add Meta request")

	vars := mux.Vars(r)
	pid := vars["pid"]
	key := vars["key"]
	value := vars["value"]

	log.WithFields(log.Fields{
		"pid":   pid,
		"key":   key,
		"value": value,
	}).Info("Adding meta")

	pid_i, err := strconv.Atoi(pid)
	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	err = hs.LocalPeer.Database.AddMeta(pid_i, key, value)
	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}

func (hs *HTTPServer) GetMeta(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Get Meta request")

	vars := mux.Vars(r)
	pid := vars["pid"]
	key := vars["key"]

	log.WithFields(log.Fields{
		"pid": pid,
		"key": key,
	}).Info("Getting meta")

	pid_i, err := strconv.Atoi(pid)
	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	value, err := hs.LocalPeer.Database.GetMeta(pid_i, key)
	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_value(w, value)
}

func (hs *HTTPServer) Mirror(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Mirror request")

	vars := mux.Vars(r)

	peer, err := hs.LocalPeer.ConnectPeer(vars["address"])

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	// Open a database for the peer
	os.Mkdir(fmt.Sprintf("%s/%s", "./data", peer.ZifAddress.Encode()), 0777)
	db := NewDatabase(fmt.Sprintf("%s/%s/posts.db", "./data", peer.ZifAddress.Encode()))
	db.Connect()

	hs.LocalPeer.Databases.Set(peer.ZifAddress.Encode(), db)

	_, err = peer.Mirror(db)

	hs.LocalPeer.Databases.Set(peer.ZifAddress.Encode(), db)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}

func (hs *HTTPServer) Peers(w http.ResponseWriter, r *http.Request) {
	log.Info("Peers request")

	peers := make([]*Peer, 0, hs.LocalPeer.peers.Count())

	for i := range hs.LocalPeer.peers.IterBuffered() {
		peers = append(peers, i.Val.(*Peer))
	}

	data, err := json.Marshal(peers)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_value(w, string(data))
}

func (hs *HTTPServer) SaveCollection(w http.ResponseWriter, r *http.Request) {
	hs.LocalPeer.Collection.Save("./data/collection.dat")

	http_write_ok(w)
}

func (hs *HTTPServer) RebuildCollection(w http.ResponseWriter, r *http.Request) {
	var err error
	hs.LocalPeer.Collection, err = CreateCollection(hs.LocalPeer.Database, 0, PieceSize)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}
