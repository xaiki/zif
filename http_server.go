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

	peer, err := hs.localPeer.ConnectPeerDirect(vars["address"])
	defer peer.Terminate()

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	peer.Ping()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}

func (hs *HTTPServer) Who(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer, err := hs.localPeer.ConnectPeerDirect(vars["address"])
	defer peer.CloseStreams()

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	entry, err := peer.Who()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ret, err := EntryToJson(&entry)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ret)
}

func (hs *HTTPServer) Announce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	peer, err := hs.localPeer.ConnectPeerDirect(vars["address"])
	defer peer.CloseStreams()

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

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
