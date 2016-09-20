package main

import (
	"strconv"
	"testing"
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
	lp.Entry.PublicKey = lp.publicKey
	lp.Entry.ZifAddress = lp.ZifAddress

	lp.SignEntry()

	lp.Listen("0.0.0.0:" + strconv.Itoa(port))

	return lp
}

func BootstrapLocalPeer(lp *LocalPeer, peer *LocalPeer, t *testing.T) {
	p, err := lp.ConnectPeerDirect(peer.Entry.PublicAddress + ":" +
		strconv.Itoa(peer.Entry.Port))

	if err != nil {
		t.Fatal(err.Error())
	}

	p.ConnectClient(lp)

	stream, err := p.Bootstrap(&lp.RoutingTable)

	if err != nil {
		t.Fatal(err.Error())
	}

	stream.Close()

}

// Does a *simple* test of announcing.
// Further testig with a much larger number of networked peers (over the internet)
// is definitely needed.
func TestLocalPeerAnnounce(t *testing.T) {
	// both peers know this one
	lp_initial := CreateLocalPeer("initial", 5050)

	// this peers announces itself
	lp_announcer := CreateLocalPeer("announcer", 5051)

	// This peer should be able to resolve the announcers address after it announces
	lp_test := CreateLocalPeer("test", 5052)

	BootstrapLocalPeer(&lp_announcer, &lp_initial, t)
	BootstrapLocalPeer(&lp_test, &lp_initial, t)

	// announce the test to the initial peer
	peer, err := lp_test.ConnectPeer(lp_initial.ZifAddress.Encode())

	if err != nil {
		t.Fatal(err.Error())
	}

	peer.Announce(&lp_test)
	////////////////////////////////////////////

	if lp_initial.RoutingTable.NumPeers() != 1 {
		t.Fatal("Failed to store announcements properly")
	}

	// announce the announcer to the initial peer
	peer, err = lp_announcer.ConnectPeer(lp_initial.ZifAddress.Encode())

	if err != nil {
		t.Fatal(err.Error())
	}

	peer.Announce(&lp_announcer)
	////////////////////////////////////////////

	if lp_initial.RoutingTable.NumPeers() != 2 {
		t.Fatal("Failed to store announcements properly")
	}

	// block until the test peer recieved a msg (should be an announce forward)
	<-lp_test.msg_chan

	if !lp_test.RoutingTable.FindClosest(lp_initial.ZifAddress, 1)[0].
		ZifAddress.Equals(&lp_initial.ZifAddress) {
		t.Fatal("Announce forwarding failed")
	}
}
