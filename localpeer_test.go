package main

import (
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"
)

func CreateLocalPeer(name string, port int) LocalPeer {
	var lp LocalPeer

	lp.GenerateKey()
	lp.Setup()
	lp.RoutingTable.Setup(lp.ZifAddress)

	lp.Entry.Name = name
	lp.Entry.Port = port
	lp.Entry.PublicAddress = "127.0.0.1"
	lp.Entry.Desc = "Decentralize all the things!"

	lp.SignEntry()

	return lp
}

func TestLocalPeerAnnounce(t *testing.T) {
	const peer_count = 40

	log.SetLevel(log.InfoLevel)

	peers := make([]LocalPeer, 0, peer_count)

	for i := 0; i < peer_count; i++ {
		peers = append(peers, CreateLocalPeer(string(i), 5050+i))
		peers[i].Listen("0.0.0.0:" + strconv.Itoa(peers[i].Entry.Port))
	}

	// connect half of the peers to the first node
	for i := 1; i < peer_count/2; i++ {
		peer, err := peers[i].ConnectPeerDirect(peers[0].Entry.PublicAddress + ":" + strconv.Itoa(peers[0].Entry.Port))

		peer.ConnectClient()
		// TODO: THIS SUCKS GET RID OF IT
		go peers[i].ListenStream(peer)

		if err != nil {
			t.Fatal(err.Error())
			return
		}

		err = peer.Announce(&peers[i])

		if err != nil {
			t.Fatal(err.Error())
			return
		}
	}

}
