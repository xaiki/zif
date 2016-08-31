package main

import (
	"fmt"

	"github.com/valyala/gorpc"
)

type RPC struct {
	dispatcher *gorpc.Dispatcher
	server     *gorpc.Server
}

func (rpc *RPC) Setup() {
	rpc.dispatcher = gorpc.NewDispatcher()

	rpc.dispatcher.AddService("service", &RPCService{})
}

func (rpc *RPC) Stop() {
	rpc.server.Stop()
}

func (rpc *RPC) Listen(address string) {
	rpc.server = gorpc.NewTCPServer(address, rpc.dispatcher.NewHandlerFunc())

	if err := rpc.server.Start(); err != nil {
		panic(err)
	}

	fmt.Println("RPC server listening on", address)
}
