package subpub

import (
	"context"
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

    subs       []subscr
}

type subscr struct {
    uuid       uuid.UUID

    subject    string
    ch         chan interface{}

    subs       []subscr

    cancelFunc context.CancelFunc
}

func (s *subscr) Unsubscribe() {
    s.cancelFunc()
    
    for i, sub := range s.subs {
        if (sub.uuid.String() == s.uuid.String()) {
            s.subs[i] = s.subs[len(s.subs) - 1]
            s.subs = s.subs[:len(s.subs) - 1]
            return
        }
    }
}

func (s *subpub) Subscribe(subject string, cb MessageHandler) (Subscription, error){
    uuid, err := uuid.NewUUID()

    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithCancel(s.ctx)

    sub := &subscr {
        uuid:        uuid,
        subject:     subject,
        ch:          make(chan interface{}),
        subs:        s.subs,
        cancelFunc:  cancel,

    }
    
    go func(ch chan interface{}, mh MessageHandler) {
       for {
            select {
            case <-ctx.Done():
                return
            default:
                select {
                case msg := <-ch:
                    mh(msg)
                }
            }
       }
    }(sub.ch, cb)
    
    s.subs = append(s.subs, *sub)

    return sub, nil
}

func (s subpub) Publish(subject string, msg interface{}) error {
    for _, sub := range s.subs {
        if (sub.subject == subject) {
            go func() {
                sub.ch <- msg       
            }()
        }
    }
    return nil
}

func (s subpub) Close (ctx context.Context) error {
    return nil
}

func NewSubPub() SubPub {
    ctx, cancel := context.WithCancel(context.Background())
    return &subpub {
        ctx: ctx,
        cancelFunc: cancel,
        subs: make([]subscr, 0),
    }
}
