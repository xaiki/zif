// Used to control the Zif daemon

package zif

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	router.HandleFunc("/bootstrap/{address}/", hs.Bootstrap)

	router.HandleFunc("/peer/{address}/resolve/", hs.Resolve)
	router.HandleFunc("/peer/{address}/ping/", hs.Ping)
	router.HandleFunc("/peer/{address}/announce/", hs.Announce)
	router.HandleFunc("/peer/{address}/search/{query}/", hs.PeerSearch)
	router.HandleFunc("/peer/{address}/recent/{page}/", hs.Recent)

	router.HandleFunc("/self/addpost/", hs.AddPost).Methods("POST")
	router.HandleFunc("/self/index/", hs.FtsIndex).Methods("POST")

	log.Info("Starting HTTP server on ", addr)

	err := http.ListenAndServe(addr, router)

	if err != nil {
		panic(err)
	}
}

func (hs *HTTPServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Zif"))
}

func (hs *HTTPServer) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}

func (hs *HTTPServer) Bootstrap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer, err := hs.LocalPeer.ConnectPeerDirect(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}
	peer.ConnectClient(hs.LocalPeer)

	stream, err := peer.Bootstrap(&hs.LocalPeer.RoutingTable)
	defer stream.Close()

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (hs *HTTPServer) Announce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer, err := hs.LocalPeer.ConnectPeer(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	peer.ConnectClient(hs.LocalPeer)

	peer.Announce(hs.LocalPeer)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (hs *HTTPServer) Resolve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := vars["address"]

	log.Info("Attempting to resolve address ", addr)

	entry, err := hs.LocalPeer.Resolve(addr)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))

		return
	}

	entry_json, err := EntryToJson(entry)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(entry_json)
}

func (hs *HTTPServer) PeerSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := vars["address"]
	query := vars["query"]

	log.Info("Searching ", addr, " for ", query)

	peer, err := hs.LocalPeer.ConnectPeer(addr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	posts, stream, err := peer.Search(query)
	defer stream.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	json, err := json.Marshal(posts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Write(json)
}

func (hs *HTTPServer) AddPost(w http.ResponseWriter, r *http.Request) {
	post_json := r.FormValue("data")

	log.Debug("Adding post, json: ", post_json)

	var post Post

	err := json.Unmarshal([]byte(post_json), &post)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	hs.LocalPeer.AddPost(post)
}

func (hs *HTTPServer) Recent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	page_i, err := strconv.Atoi(page)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return

	}

	posts, err := hs.LocalPeer.Database.QueryRecent(page_i)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	json, err := json.Marshal(posts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Write(json)
}

func (hs *HTTPServer) FtsIndex(w http.ResponseWriter, r *http.Request) {
	since := r.FormValue("since")

	since_i, err := strconv.Atoi(since)

	log.Info("Generating FTS index since ", since_i)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	hs.LocalPeer.Database.GenerateFts(uint64(since_i))

	w.Write([]byte("{ \"status\": \"ok\"}"))
}
