package broadcaster

import "sync"

type Broadcaster[T any] struct {
    listeners map[*Listener[T]]bool
    lMutex    *sync.RWMutex
}

func NewBroadcaster[T any]() *Broadcaster[T] {
    return &Broadcaster[T]{
        listeners: make(map[*Listener[T]]bool),
        lMutex:    new(sync.RWMutex),
    }
}

func (bc *Broadcaster[T]) Register(bufSize int) *Listener[T] {
    l := &Listener[T]{
        bc: bc,
        ch: make(chan Payload[T], bufSize),
    }

    bc.lMutex.Lock()
    bc.listeners[l] = true
    bc.lMutex.Unlock()
    
    return l
}

func (bc *Broadcaster[T]) Send(payload T) Waiter[T] {
    bc.lMutex.RLock()
    defer bc.lMutex.RUnlock()

    w := &payloadWaiter[T]{
        listeners: make([]*Listener[T], 0, len(bc.listeners)),
        payload:   payload,
    }
    w.wg.Add(len(bc.listeners))

    for l := range bc.listeners {
        w.listeners = append(w.listeners, l)
        l.ch <- w
    }

    return w
}

func (bc *Broadcaster[T]) Close() {
    bc.lMutex.Lock()
    defer bc.lMutex.Unlock()

    for l := range bc.listeners {
        l.unregisterNoLock()
    }
}
