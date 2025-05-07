package main

import (
    "fmt"

    "vk-test/internal/server"
    pb "vk-test/internal/pubsubpb"

    grpc "google.golang.org/grpc"

    "net"
    "log"
    "vk-test/pkg/logger"
    "vk-test/pkg/config"
)

func main() {
    l := logger.NewLogger()
    cfg := config.GetConfig(l)

    l.Info(cfg.Port)

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))

    if err != nil {
        l.Fatal("Cannot start server: %s", err)
    }

    s := server.NewServer(l)

    rpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(logger.LoggingInterceptor(l)),
        grpc.StreamInterceptor(logger.StreamLoggingInterceptor(l)),
    )

    pb.RegisterPubSubServer(rpcServer, s)

    err = rpcServer.Serve(lis)

    if err != nil {
        log.Fatal("Impossible to serve: %s", err)
    }
}
