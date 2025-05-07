package subpub

import (
	"context"
	_ "fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type MessageHandler func (msg interface{})

type Subscription interface {
    Unsubscribe()
}

type SubPub interface {
    Subscribe(subject string, cb MessageHandler) (Subscription, error)
    Publish(subject string, msg interface{}) error
    Close (ctx context.Context) error
}

type subpub struct {
    ctx        context.Context
    cancelFunc context.CancelFunc

    mu         *sync.RWMutex
    subs       []subscr 
}

type subscr struct {
    uuid       uuid.UUID
    subject    string
    ch         chan interface{}

    parent     *subpub
    
    mu         *sync.RWMutex

    ctx        context.Context
    cancelFunc context.CancelFunc
}

func IsSubValid(sub Subscription) bool {
    select {
    case <-sub.(*subscr).ctx.Done():
        return false
    default:
        return true
    }
}

func messageProcessing(sub *subscr, cb MessageHandler) {
   
    prevDone := make(chan struct{}, 1)
    prevDone <- struct{}{}
    
    go func() {
        <-sub.ctx.Done()
        close(sub.ch)
    }()

    for msg := range sub.ch {
        
        curDone := make(chan struct{}, 1)
        started := make(chan struct{})

        msgCopy := msg

        go func(prevDone chan struct{}, done chan struct{}, msgCopy interface{}) {
            for {
                select {
                    case <-sub.ctx.Done():
                        return
                    default:
                        select {
                            case <-prevDone:
                                close(started)
                                cb(msgCopy)
                                
                                done <- struct{}{}
                                close(prevDone)

                                return
                        }
                }
            }
        }(prevDone, curDone, msgCopy)

        <-started

        prevDone = curDone
    }
}

func (s *subscr) Unsubscribe() {
    s.mu.Lock()
    
    s.cancelFunc()
    
    for i, sub := range s.parent.subs {
        if (strings.Compare(sub.uuid.String(), s.uuid.String()) == 0) {
            s.parent.subs[i] = s.parent.subs[len(s.parent.subs) - 1]
            s.parent.subs = s.parent.subs[:(len(s.parent.subs) - 1)]
            
            s.mu.Unlock()
            
            return
        }
    }
    
    s.mu.Unlock()
}

func (s *subpub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
    uuid, err := uuid.NewUUID()

    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithCancel(s.ctx)

    sub := &subscr {
        uuid:        uuid,
        subject:     subject,
        ch:          make(chan interface{}, 5),
        mu:          s.mu,
        parent:      s,
        ctx:         ctx,
        cancelFunc:  cancel,
    }
    
    s.mu.Lock()

    s.subs = append(s.subs, *sub)

    s.mu.Unlock()

    go messageProcessing(sub, cb)

    return sub, nil
}

func (s *subpub) Publish(subject string, msg interface{}) error {
    s.mu.RLock()

    for _, sub := range s.subs {
        if (sub.subject == subject) {
            sub.ch <- msg
        }
    }

    s.mu.RUnlock()
    return nil
}

func (s *subpub) Close (ctx context.Context) error {
    for _ = range ctx.Done() {
        return nil
    }

    s.cancelFunc()
    return nil
}

func NewSubPub() SubPub {
    ctx, cancel := context.WithCancel(context.Background())

    return &subpub {
        ctx:         ctx,
        cancelFunc:  cancel,
        mu:          &sync.RWMutex{},
        subs:        make([]subscr, 0),
    }
}
