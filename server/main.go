package main

import (
    server "vk-test/grpcServer"
    rpc "vk-test/grpc"

    grpc "google.golang.org/grpc"

    "net"
    "log"
    "vk-test/logger"
)

func main() {
    logger.NewLogger()

    lis, err := net.Listen("tcp", ":8089")

    if err != nil {
        log.Fatal("Cannot start server: %s", err)
    }

    s := rpc.NewServer()

    rpcServer := grpc.NewServer()
    server.RegisterPubSubServer(rpcServer, s)

    err = rpcServer.Serve(lis)

    if err != nil {
        log.Fatal("Impossible to serve: %s", err)
    }
}
