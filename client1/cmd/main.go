package main

import (
	"client/internal/pb"
    "client/internal/config"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func main() {
    cfg := config.GetConfig()
    conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.ServerIp, cfg.ServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal("Error in dialing: %v", err)
    }
    
    client := clientpb.NewPubSubClient(conn)
    
    var stream clientpb.PubSub_SubscribeClient

    go func() { stream, err = client.Subscribe(context.Background(), &clientpb.SubscribeRequest{Key: "test"}) }()

    go func() {
        for {
            client.Publish(context.Background(), &clientpb.PublishRequest{Key: "test", Data: "test"})
        }
    }()

    
    time.Sleep(time.Second)
    if err != nil {
        fmt.Println(err)
        return
    }
    for {
        event, err := stream.Recv()

        if err == io.EOF {
            fmt.Printf("%v\n", err)
            break
        }

        fmt.Println(event.GetData())
    }
}
