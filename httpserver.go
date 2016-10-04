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

	router.HandleFunc("/peer/{address}/ping/", hs.Ping)
	router.HandleFunc("/peer/{address}/announce/", hs.Announce)
	router.HandleFunc("/peer/{address}/search/{query}/", hs.PeerSearch)
	router.HandleFunc("/peer/{address}/recent/{page}/", hs.Recent)

	router.HandleFunc("/self/addpost/", hs.AddPost).Methods("POST")
	router.HandleFunc("/self/index/", hs.FtsIndex).Methods("POST")
	router.HandleFunc("/self/resolve/{address}", hs.Resolve)
	router.HandleFunc("/self/bootstrap/{address}/", hs.Bootstrap)

	log.Info("Starting HTTP server on ", addr)

	err := http.ListenAndServe(addr, router)

	if err != nil {
		panic(err)
	}
}

func http_error_check(w http.ResponseWriter, errCode int, err error) bool {
	if err != nil {
		w.WriteHeader(errCode)
		w.Write([]byte("{ \"status\": \"err\", \"err\": \"" + err.Error() + "\"}"))

		return true
	}

	return false
}

func http_write_ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"ok\" }"))
}

func http_write_posts(w http.ResponseWriter, posts []*Post) {
	json, err := json.Marshal(posts)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	post_length := strconv.Itoa(len(posts))

	// TODO: Use/write some sort of json building, based on a map.
	// This is kinda gross.
	w.Write([]byte("{ \"status\": \"ok\", \"count\": " + post_length + ", \"posts\": " + string(json) + "}"))
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

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	peer.ConnectClient(hs.LocalPeer)

	stream, err := peer.Bootstrap(&hs.LocalPeer.RoutingTable)
	defer stream.Close()

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}

func (hs *HTTPServer) Announce(w http.ResponseWriter, r *http.Request) {
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
	w.Write(entry_json)
}

func (hs *HTTPServer) PeerSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := vars["address"]
	query := vars["query"]

	log.Info("Searching ", addr, " for ", query)

	peer, err := hs.LocalPeer.ConnectPeer(addr)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	posts, stream, err := peer.Search(query)
	defer stream.Close()

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_posts(w, posts)
}

func (hs *HTTPServer) AddPost(w http.ResponseWriter, r *http.Request) {
	post_json := r.FormValue("data")

	log.Debug("Adding post, json: ", post_json)

	var post Post

	err := json.Unmarshal([]byte(post_json), &post)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	hs.LocalPeer.AddPost(post)

	http_write_ok(w)
}

func (hs *HTTPServer) Recent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	page_i, err := strconv.Atoi(page)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	peer, err := hs.LocalPeer.ConnectPeer(vars["address"])

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	posts, stream, err := peer.Recent(uint64(page_i))
	defer stream.Close()

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_posts(w, posts)
}

func (hs *HTTPServer) FtsIndex(w http.ResponseWriter, r *http.Request) {
	since := r.FormValue("since")

	since_i, err := strconv.Atoi(since)

	log.Info("Generating FTS index since ", since_i)

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	err = hs.LocalPeer.Database.GenerateFts(uint64(since_i))

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

	http_write_ok(w)
}
