package main

import (
	"fmt"

	pb "vk-test/internal/pubsubpb"
	"vk-test/internal/server"

	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"

	"log"
	"net"
	"vk-test/pkg/config"
	"vk-test/pkg/logger"
)

func main() {
    l := logger.NewLogger()
    cfg := config.GetConfig(l)

    if cfg.IsDebug {
        l.SetLevel(logrus.DebugLevel)
    }

    lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.BindIp, cfg.Port))
    
    if err != nil {
        l.Fatal("Cannot start server: %s", err)
    }

    s := server.NewServer(l)

    rpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(logger.LoggingInterceptor(l)),
        grpc.StreamInterceptor(logger.StreamLoggingInterceptor(l)),
    )

    pb.RegisterPubSubServer(rpcServer, s)
    
    l.Info(fmt.Sprintf("Server is listening on: %s:%d", cfg.BindIp, cfg.Port))
    err = rpcServer.Serve(lis)

    if err != nil {
        log.Fatal("Impossible to serve: %s", err)
    }
}
