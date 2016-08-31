// The local peer. This runs on the current node, so we have access to its
// private key, database, etc.

package main

import (
	"errors"
	"io/ioutil"

	"golang.org/x/crypto/ed25519"
)

type LocalPeer struct {
	Peer
	Entry        Entry
	RPC          RPC
	RoutingTable RoutingTable
	privateKey   ed25519.PrivateKey
}

func (lp *LocalPeer) Setup() {
	lp.ZifAddress.Generate(lp.publicKey)
	lp.RPC.Setup()
}

// address, router (TCP) port, dht (udp) port
func (lp *LocalPeer) Listen(addr string) {
	go lp.RPC.Listen(addr)
}

func (lp *LocalPeer) GenerateKey() {
	var err error

	lp.publicKey, lp.privateKey, err = ed25519.GenerateKey(nil)

	if err != nil {
		panic(err)
	}
}

// Writes the private key to a file, in this way persisting your identity -
// all the other addresses can be generated from this, no need to save them.
func (lp *LocalPeer) WriteKey() error {
	if len(lp.privateKey) == 0 {
		return errors.
			New("LocalPeer does not have a private key, please generate")
	}

	err := ioutil.WriteFile("identity.dat", lp.privateKey, 0400)

	return err
}

func (lp *LocalPeer) ReadKey() error {
	pk, err := ioutil.ReadFile("identity.dat")

	if err != nil {
		return err
	}

	lp.privateKey = pk
	lp.publicKey = lp.privateKey.Public().(ed25519.PublicKey)

	return nil
}
