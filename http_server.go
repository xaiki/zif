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
	router.HandleFunc("/query/{address}/", hs.Query)
	router.HandleFunc("/announce/{address}/", hs.Announce)
	router.HandleFunc("/set_address/{address}/", hs.SetAddress)
	router.HandleFunc("/bootstrap/{address}/", hs.Bootstrap)

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

// TODO: Query needs to work without specifying a direct connection. Other requests
// need to be able to as well. So, make a bootstrap request which will kickstart
// a dht table. Then, all further requests can use this :D
func (hs *HTTPServer) Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer := NewPeer(hs.localPeer)
	err := peer.Connect(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}
	peer.ConnectClient()

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
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
