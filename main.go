package main

import (
	"fmt"
	"sync"
	"time"
	"vk-test/subpub"
)

func Handler(msg interface{}) {
    msgString := msg.(string)
    fmt.Printf("Message received: %s\n", msgString)
}

func main() {
    const countOfTopics = 2

    sp := subpub.NewSubPub()
    for i := 0; i < countOfTopics; i++ {
        go func() {
            topicName := fmt.Sprintf("test%d", i)
            for {
                sp.Publish(topicName, topicName)
                time.Sleep(time.Millisecond * 100)
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
    
            dur := (time.Duration)(2 * (i + 1))
            time.Sleep(time.Second * dur) 
            
            sub.Unsubscribe()
        }()
    }
    
    wg.Wait()
}
