package main

import (
	"vk-test/subpub"
    "vk-test/server"
    rpc "vk-test/grpc"

    grpc "google.golang.org/grpc"

    "net"
    "log"
)

func main() {
    subpub := subpub.NewSubPub()

    lis, err := net.Listen("tcp", ":8089")

    if err != nil {
        log.Fatal("Cannot start server: %s", err)
    }

    s := &rpc.Server{
        Subpub: subpub,
    }

    rpcServer := grpc.NewServer()
    server.RegisterPubSubServer(rpcServer, s)
    err = rpcServer.Serve(lis)

    if err != nil {
        log.Fatal("Impossible to serve: %s", err)
    }
}
