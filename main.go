package main

import (
	"time"
    "fmt"
	"vk-test/subpub"
)

func Handler(msg interface{}) {
    msgString := msg.(string)
    fmt.Printf("Message received: %s\n", msgString)
}

func main() {
    sp := subpub.NewSubPub()
    go func() {
        for {
            sp.Publish("test", "hello!")
            time.Sleep(time.Millisecond * 100)
        }
    }()

    sub, _ := sp.Subscribe("test", Handler)
    
    time.Sleep(time.Second * 10)

    sub.Unsubscribe()
}
