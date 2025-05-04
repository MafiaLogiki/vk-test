package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"vk-test/subpub"
    
    "net/http"
    _ "net/http/pprof"

    "runtime"
)

func Handler(msg interface{}) {
    msgString := msg.(string)
    fmt.Printf("Message received: %s\n", msgString)
}


func TestHandler(msg interface{}) {
    fmt.Printf("Handle %d\n", msg.(int))
}

func main() {
    const countOfTopics = 2

    fmt.Println(runtime.NumGoroutine())

    sp := subpub.NewSubPub()
    
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
  
    for i := 0; i < countOfTopics; i++ {
        go func() {
            topicName := fmt.Sprintf("test%d", i)
            for i := 0; i < 100000; i++ {
                sp.Publish(topicName, topicName)
            }
        }()
    }

    wg := sync.WaitGroup{}
    for i := 0; i < countOfTopics; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            topicName := fmt.Sprintf("test%d", i)
            sub, _ := sp.Subscribe(topicName, Handler)
    
            dur := (time.Duration)(i + 2)
            time.Sleep(time.Second * dur) 
            
            sub.Unsubscribe()
        }()
    }
    
    wg.Wait()
    
    fmt.Println(runtime.NumGoroutine())
    
    sp.Subscribe("test", TestHandler)
    for i := 0; i < 3; i++ {
        sp.Publish("test", i)
    }


    ctx, cancel := context.WithCancel(context.Background())

    fmt.Println(runtime.NumGoroutine())

    cancel()
    sp.Close(ctx)

    time.Sleep(time.Second)

    fmt.Println(runtime.NumGoroutine())

    for {}
}
