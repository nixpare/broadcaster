package broadcaster

import (
	"sync"
)

type broadcaster[T comparable] interface {
    getListeners() map[T]struct{}
    getMutex()     *sync.RWMutex
}

type Broadcaster[T any] struct {
    listeners map[*Channel[T]]struct{}
    mutex    *sync.RWMutex
    closed    bool
}

func NewBroadcaster[T any]() *Broadcaster[T] {
    return &Broadcaster[T]{
        listeners: make(map[*Channel[T]]struct{}),
        mutex:    new(sync.RWMutex),
    }
}

func (bc *Broadcaster[T]) getListeners() map[*Channel[T]]struct{} {
    return bc.listeners
}

func (bc *Broadcaster[T]) getMutex() *sync.RWMutex {
    return bc.mutex
}

func (bc *Broadcaster[T]) Register(bufSize int) *Channel[T] {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
   
    ch := &Channel[T]{ bc: bc, ch: make(chan T, bufSize) }

    if bc.closed {
        ch.closed = true
        close(ch.ch)
    } else {
        bc.listeners[ch] = struct{}{}
    }
    
    return ch
}

func (bc *Broadcaster[T]) Send(payload T) {
    bc.mutex.RLock()
    defer bc.mutex.RUnlock()

    for ch := range bc.listeners {
        ch.ch <- payload
    }
}

func (bc *Broadcaster[T]) Reset() {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    bc.resetNoLock()
}

func (bc *Broadcaster[T]) resetNoLock() {
    for ch := range bc.listeners {
        ch.unregisterNoLock()
    }
}

func (bc *Broadcaster[T]) Close() {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    if bc.closed {
        return
    }
    bc.closed = true

    bc.resetNoLock()
}

type BufBroadcaster[T any] struct {
    listeners map[*Channel[T]]struct{}
    mutex    *sync.RWMutex
    closed    bool

    data      []T
}

func NewBufBroadcaster[T any]() *BufBroadcaster[T] {
    return &BufBroadcaster[T]{
        listeners: make(map[*Channel[T]]struct{}),
        mutex:    new(sync.RWMutex),
    }
}

func (bc *BufBroadcaster[T]) getListeners() map[*Channel[T]]struct{} {
    return bc.listeners
}

func (bc *BufBroadcaster[T]) getMutex() *sync.RWMutex {
    return bc.mutex
}

func (bc *BufBroadcaster[T]) Register(bufSize int) *Channel[T] {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
   
    return bc.registerNoLock(bufSize)
}

func (bc *BufBroadcaster[T]) registerNoLock(bufSize int) *Channel[T] {
    ch := &Channel[T]{ bc: bc, ch: make(chan T, bufSize) }

    if bc.closed {
        ch.closed = true
        close(ch.ch)
    } else {
        bc.listeners[ch] = struct{}{}
    }
    
    return ch
}

func (bc *BufBroadcaster[T]) Send(payload T) {
    bc.mutex.RLock()
    defer bc.mutex.RUnlock()

    bc.data = append(bc.data, payload)
    for ch := range bc.listeners {
        ch.ch <- payload
    }
}

func (bc *BufBroadcaster[T]) Data() []T {
    return bc.data
}

func (bc *BufBroadcaster[T]) Connect(bufSize int) ([]T, *Channel[T]) {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    return bc.Data(), bc.registerNoLock(bufSize)
}

func (bc *BufBroadcaster[T]) Reset() {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    bc.resetNoLock()
}

func (bc *BufBroadcaster[T]) resetNoLock() {
    for ch := range bc.listeners {
        ch.unregisterNoLock()
    }
    bc.data = nil
}

func (bc *BufBroadcaster[T]) Close() {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    if bc.closed {
        return
    }
    bc.closed = true

    bc.resetNoLock()
}
