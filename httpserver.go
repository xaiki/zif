// Used to control the Zif daemon

package main

import (
	"net/http"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type HTTPServer struct {
	localPeer *LocalPeer
}

func (hs *HTTPServer) ListenHTTP(addr string) {
	router := mux.NewRouter().StrictSlash(true)

	// TODO: Bootstrap request
	router.HandleFunc("/", hs.IndexHandler)

	// This should be the ONLY route where the address is a non-Zif address
	router.HandleFunc("/bootstrap/{address}/", hs.Bootstrap)

	router.HandleFunc("/peer/{address}/resolve/", hs.Resolve)
	router.HandleFunc("/peer/{address}/ping/", hs.Ping)
	router.HandleFunc("/peer/{address}/announce/", hs.Announce)

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

func (hs *HTTPServer) Announce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer, err := hs.localPeer.ConnectPeer(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	peer.ConnectClient()

	// TODO: Not have to do this every time I connect a client... hmm.
	go hs.localPeer.ListenStream(peer)

	peer.Announce(hs.localPeer)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (hs *HTTPServer) Bootstrap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer, err := hs.localPeer.ConnectPeerDirect(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}
	peer.ConnectClient()
	go hs.localPeer.ListenStream(peer)

	stream, err := peer.Bootstrap(&hs.localPeer.RoutingTable)
	defer stream.Close()

	if err != nil {
		log.Error("Failed to bootstrap: ", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (hs *HTTPServer) Resolve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := vars["address"]

	log.Info("Attempting to resolve address ", addr)

	entry, err := hs.localPeer.Resolve(addr)

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
