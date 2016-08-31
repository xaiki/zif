// Used to control the Zif daemon

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	localPeer *LocalPeer
}

func (hs *HTTPServer) ListenHTTP(addr string) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", hs.IndexHandler)
	router.HandleFunc("/ping/{address}/", hs.Ping)
	/*router.HandleFunc("/query/{address}/{dht}/{target}/", hs.Query)*/
	router.HandleFunc("/announce/{address}/", hs.Announce)

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

func (hs *HTTPServer) Ping(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	c, ok := ConnectClient(vars["address"], hs.localPeer)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Connection failed"))
	}

	c.Handshake()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}

/*func (hs *HTTPServer) Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var peer Peer
	peer.DHTAddress = fmt.Sprintf("%s:%v", vars["address"], vars["dht"])

	peer.ConnectUDP()
	peer.dht_client.Query(hs.localPeer, DecodeAddress(vars["target"]))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Done."))
}*/

func (hs *HTTPServer) Announce(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var peer Peer
	peer.PublicAddress = vars["address"]

	peer.Connect()
	ok := peer.Announce(&hs.localPeer.Entry)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ok))
}
