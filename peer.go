// This represents a peer in the network.
// the minimum that a peer requires to be "valid" is just an address.
// everything else can be discovered via the network.

package main

import (
	"fmt"

	"github.com/valyala/gorpc"
	"golang.org/x/crypto/ed25519"
)

type Peer struct {
	ZifAddress    Address
	PublicAddress string
	Client        *gorpc.Client
	Dispatch      *gorpc.DispatcherClient
	publicKey     ed25519.PublicKey

	dispatcher *gorpc.Dispatcher
}

func CreatePeer(pub ed25519.PublicKey, addr string, port int) Peer {
	var ret Peer

	ret.publicKey = pub
	ret.ZifAddress.Generate(ret.publicKey)

	return ret
}

func (p *Peer) Connect() {
	p.dispatcher = gorpc.NewDispatcher()
	p.dispatcher.AddService("service", &RPCService{})

	p.Client = gorpc.NewTCPClient(p.PublicAddress)
	p.Dispatch = p.dispatcher.NewServiceClient("service", p.Client)
	p.Client.Start()
}

func (p *Peer) Ping() {
	fmt.Println(p.Dispatch.Call("Ping", nil))
}

func (p *Peer) Announce(entry Entry) {
	p.Dispatch.Call("Announce", entry)
}
