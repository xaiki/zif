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
	router.HandleFunc("/peer/{address}/popular/{page}/", hs.Popular)

	router.HandleFunc("/self/addpost/", hs.AddPost).Methods("POST")
	router.HandleFunc("/self/index/", hs.FtsIndex).Methods("POST")
	router.HandleFunc("/self/resolve/{address}", hs.Resolve)
	router.HandleFunc("/self/bootstrap/{address}/", hs.Bootstrap)
	router.HandleFunc("/self/search/", hs.SelfSearch).Methods("POST")
	router.HandleFunc("/self/recent/{page}/", hs.SelfRecent)
	router.HandleFunc("/self/popular/{page}/", hs.SelfPopular)
	router.HandleFunc("/self/addmeta/{pid}/{key}/{value}/", hs.AddMeta)
	router.HandleFunc("/self/getmeta/{pid}/{key}/", hs.GetMeta)

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
	log.Info("HTTP: Bootstrap request")
	vars := mux.Vars(r)

	peer, err := hs.LocalPeer.ConnectPeerDirect(vars["address"])

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
	w.Write(entry_json)
}

func (hs *HTTPServer) PeerSearch(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Peer Search request")
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
	log.Info("HTTP: Add Post request")

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
	defer stream.Close()

	if http_error_check(w, http.StatusInternalServerError, err) {
		return
	}

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

func (hs *HTTPServer) SelfSearch(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP: Self Search request")

	query := r.FormValue("query")
	page := r.FormValue("page")

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
