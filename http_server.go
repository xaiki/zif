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

	router.HandleFunc("/", hs.IndexHandler)
	router.HandleFunc("/ping/{address}/", hs.Ping)
	router.HandleFunc("/who/{address}/", hs.Who)
	/*router.HandleFunc("/query/{address}/{dht}/{target}/", hs.Query)*/
	router.HandleFunc("/announce/{address}/", hs.Announce)
	router.HandleFunc("/set_address/{address}/", hs.SetAddress)

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

func (hs *HTTPServer) Who(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer := NewPeer(hs.localPeer)
	err := peer.Connect(vars["address"])

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}
	peer.ConnectClient()

	client, entry, err := peer.Who()
	defer client.Close()

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	json, err := EntryToJson(&entry)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
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

func (hs *HTTPServer) SetAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hs.localPeer.Entry.PublicAddress = vars["address"]

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
