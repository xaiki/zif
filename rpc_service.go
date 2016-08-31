package main

import "fmt"

type RPCService struct {
	localPeer *LocalPeer
}

func (rs *RPCService) Ping() string {
	return "pong"
}

func (rs *RPCService) Announce(entry *Entry) string {
	fmt.Println("Announced", entry.Name)
	return "ok"
}
