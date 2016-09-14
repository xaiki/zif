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
	router.HandleFunc("/ping/{address}/", hs.Ping)
	router.HandleFunc("/announce/{address}/", hs.Announce)
	router.HandleFunc("/set_address/{address}/", hs.SetAddress)
	router.HandleFunc("/bootstrap/{address}/", hs.Bootstrap)

	router.HandleFunc("/{address}/resolve/", hs.Resolve)

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
	vars := mux.Vars(r)

	peer := NewPeer(hs.localPeer)
	err := peer.Connect(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}
	peer.ConnectClient()

	client := peer.Ping()
	defer client.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}

func (hs *HTTPServer) Announce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer := NewPeer(hs.localPeer)
	err := peer.Connect(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}
	peer.ConnectClient()

	peer.Announce()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (hs *HTTPServer) Bootstrap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer := NewPeer(hs.localPeer)
	err := peer.Connect(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}
	peer.ConnectClient()

	stream, err := peer.Bootstrap()
	defer stream.Close()

	if err != nil {
		log.Error("Failed to bootstrap: ", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (hs *HTTPServer) SetAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hs.localPeer.Entry.PublicAddress = vars["address"]

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
