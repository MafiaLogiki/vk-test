package subpub

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"container/list"

	"github.com/stretchr/testify/assert"
)

func TestCreateSubPub(t *testing.T) {
    subpub := NewSubPub()
    assert.NotEmpty(t, subpub)
}

func TestPublish(t *testing.T) {
    subpub := NewSubPub()
    err := subpub.Publish("test", "test")
    assert.Nil(t, err)
}

func TestSubscrbeAndUnsibscribe(t *testing.T) {
    subpub := NewSubPub()
    countBefore := runtime.NumGoroutine()

    sub, err := subpub.Subscribe("test", func(interface{}) {})

    assert.Nil(t, err)
    assert.NotEmpty(t, sub)
    
    sub.Unsubscribe()
    time.Sleep(time.Second) // Wait for all go-routine shut down (fighting with scheduler)

    countAfter := runtime.NumGoroutine()
    assert.Equal(t, countBefore, countAfter)
}

func TestPublish_AlreadySubscribed(t *testing.T) {
    subpub := NewSubPub()

    sub, _ := subpub.Subscribe("test", func(interface{}) {})
    
    err := subpub.Publish("test", "test")

    sub.Unsubscribe()
    time.Sleep(time.Second) // Fighting with scheduler

    assert.Nil(t, err) 
}

func getAllGoroutines() string {
    buf := make([]byte, 1 << 20)
    n := runtime.Stack(buf, true)
    return string(buf[:n])
}

func TestSubscribeAndUnsubscribe_MoreThanOneSub(t *testing.T) {
    subpub := NewSubPub()
    
    countBefore := runtime.NumGoroutine()

    sub1, err1 := subpub.Subscribe("test", func(interface{}) {})
    assert.Nil(t, err1)
    assert.NotEmpty(t, sub1)

    sub2, err2 := subpub.Subscribe("test", func(interface{}) {})
    assert.Nil(t, err2)
    assert.NotEmpty(t, sub2)
    
    sub3, err3 := subpub.Subscribe("test", func(interface{}) {})
    assert.Nil(t, err3)
    assert.NotEmpty(t, sub3)
    
    sub1.Unsubscribe()
    sub2.Unsubscribe()
    sub3.Unsubscribe()
    
    time.Sleep(time.Second) // Fighting with sheduler again

    countAfter := runtime.NumGoroutine()

    assert.Equal(t, countBefore, countAfter, "Not all goroutines done their work. Active goroutines:\n%s", getAllGoroutines())
}

func createHandlerForTest(mu *sync.Mutex, list *list.List) func(interface{}){
    return func(msg interface{}) {
        mu.Lock()
        list.PushBack(msg)
        mu.Unlock()
    }
}

func newSubscribtion(subpub SubPub) (Subscription, *list.List, *sync.Mutex){
    lst := list.New()
    mu := &sync.Mutex{}
    sub, _ := subpub.Subscribe("test", createHandlerForTest(mu, lst))

    return sub, lst, mu 
}

func TestPublish_FIFO(t *testing.T) {
    subpub := NewSubPub()
    const countOfPublishes = 100
    
    sub, msgs, mu := newSubscribtion(subpub)

    for i := 0; i < countOfPublishes; i++ {
        subpub.Publish("test", i)    
    }
    
    time.Sleep(time.Second)
   
    mu.Lock()
    elem := msgs.Front()
    for i := 0; i < countOfPublishes; i++ {
        assert.Equal(t, i, elem.Value)
        elem = elem.Next()
    }
    mu.Unlock()

    sub.Unsubscribe()
}


func TestPublish_FIFOWithMoreSubs(t *testing.T) {
    subpub := NewSubPub()

    const countOfPublishes = 10

    sub1, msgs1, mu1 := newSubscribtion(subpub)
    sub2, msgs2, mu2 := newSubscribtion(subpub)
    sub3, msgs3, mu3 := newSubscribtion(subpub)
    
    for i := 0; i < countOfPublishes; i++ {
        subpub.Publish("test", i)    
    }

    
    time.Sleep(time.Second * 2)
    
    mu1.Lock()
    mu2.Lock()
    mu3.Lock()

    elem1 := msgs1.Front()
    elem2 := msgs2.Front()
    elem3 := msgs3.Front()

    for i := 0; i < countOfPublishes; i++ {
        assert.Equal(t, i, elem1.Value)
        assert.Equal(t, i, elem2.Value)
        assert.Equal(t, i, elem3.Value)

        elem1 = elem1.Next()
        elem2 = elem2.Next()
        elem3 = elem3.Next()
    }

    mu1.Unlock()
    mu2.Unlock()
    mu3.Unlock()

    sub1.Unsubscribe()
    sub2.Unsubscribe()
    sub3.Unsubscribe()
}

func createLongFunc(lst list.List) func(interface{}){ 
    return func(msg interface{}) {
        time.Sleep(time.Millisecond * 600)
        lst.PushBack(msg)
    }
}

func TestPublish_FIFOWithLongFunc(t *testing.T) {
    subpub := NewSubPub()

    const countOfPublishes = 5

    sub, msgs, mu := newSubscribtion(subpub)

    for i := 0; i < countOfPublishes; i++ {
        subpub.Publish("test", i)
    }

    time.Sleep(time.Second)
    
    sub.Unsubscribe()
    time.Sleep(time.Second)
    
    mu.Lock()
    elem := msgs.Front()
    for i := 0; i < countOfPublishes; i++ {
        assert.Equal(t, i, elem.Value)
        elem = elem.Next()
    }

    mu.Unlock()

}
