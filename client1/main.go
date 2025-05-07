package main

import (
	"client/grpcClient"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func main() {
    conn, err := grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal("Error in dialing: %v", err)
    }
    
    client := grpcClient.NewPubSubClient(conn)
    
    var stream grpcClient.PubSub_SubscribeClient

    go func() { stream, err = client.Subscribe(context.Background(), &grpcClient.SubscribeRequest{Key: "test"}) }()

    go func() {
        for {
            client.Publish(context.Background(), &grpcClient.PublishRequest{Key: "test", Data: "test"})
        }
    }()

    
    time.Sleep(time.Second)
    
    for {
        event, err := stream.Recv()

        if err == io.EOF {
            fmt.Printf("%v\n", err)
            break
        }

        fmt.Println(event.GetData())
    }
}
