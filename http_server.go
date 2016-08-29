// Used to control the Zif daemon

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type HTTPServer struct {
	localPeer *LocalPeer
}

func (hs *HTTPServer) ListenHTTP(addr string) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", hs.IndexHandler)
	router.HandleFunc("/bootstrap/{address}/{router}/{dht}/", hs.Bootstrap)
	router.HandleFunc("/ping/{address}/{dht}/", hs.Ping)
	router.HandleFunc("/query/{address}/{dht}/{target}/", hs.Query)
	router.HandleFunc("/announce/{address}/{dht}/", hs.Announce)

	fmt.Println("Starting HTTP server on", addr)

	err := http.ListenAndServe(addr, router)

	if err != nil {
		panic(err)
	}
}

func (hs *HTTPServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Zif"))
}

// TODO: Make this actually bootstrap rather than just ping.
// atm peers don't store peer data, there is no db, there is nothing to send.
func (hs *HTTPServer) Bootstrap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var peer Peer
	peer.RouterAddress = fmt.Sprintf("%s:%v", vars["address"], vars["router"])
	peer.DHTAddress = fmt.Sprintf("%s:%v", vars["address"], vars["dht"])

	peer.Connect()
	peer.Ping(hs.localPeer)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}

func (hs *HTTPServer) Ping(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var peer Peer
	peer.DHTAddress = fmt.Sprintf("%s:%v", vars["address"], vars["dht"])

	peer.ConnectUDP()
	peer.Ping(hs.localPeer)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}

func (hs *HTTPServer) Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var peer Peer
	peer.DHTAddress = fmt.Sprintf("%s:%v", vars["address"], vars["dht"])

	peer.ConnectUDP()
	peer.dht_client.Query(hs.localPeer, DecodeAddress(vars["target"]))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}

func (hs *HTTPServer) Announce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var peer Peer
	peer.DHTAddress = fmt.Sprintf("%s:%v", vars["address"], vars["dht"])

	peer.ConnectUDP()
	peer.dht_client.Announce(hs.localPeer, hs.localPeer.Entry)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}
