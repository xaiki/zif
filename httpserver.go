// Used to control the Zif daemon

package zif

import (
	"net/http"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type HttpServer struct {
	CommandServer *CommandServer
}

func (hs *HttpServer) ListenHttp(addr string) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", hs.IndexHandler)

	// This should be the ONLY route where the address is a non-Zif address

	router.HandleFunc("/peer/{address}/ping/", hs.Ping)
	/*router.HandleFunc("/peer/{address}/announce/", hs.Announce)
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
	router.HandleFunc("/self/saveroutingtable/", hs.SaveRoutingTable)*/

	log.Info("Starting HTTP server on ", addr)

	err := http.ListenAndServe(addr, router)

	if err != nil {
		panic(err)
	}
}

func write_http_response(w http.ResponseWriter, cr CommandResult) {
	var err int
	if cr.IsOK {
		err = 200
	} else {
		err = 500
	}

	w.WriteHeader(err)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cr.WriteJSON(w)
}

func (hs *HttpServer) Ping(w http.ResponseWriter, r *http.Request) {
	// TODO
	write_http_response(w, CommandResult{true,nil,nil})
}

func (hs *HttpServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Zif"))
}

