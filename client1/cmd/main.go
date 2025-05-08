package main

import (
	"client/internal/pb"
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
    
    client := clientpb.NewPubSubClient(conn)
    
    var stream clientpb.PubSub_SubscribeClient

    go func() { stream, err = client.Subscribe(context.Background(), &clientpb.SubscribeRequest{Key: "test"}) }()

    go func() {
        for {
            client.Publish(context.Background(), &clientpb.PublishRequest{Key: "test", Data: "test"})
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
