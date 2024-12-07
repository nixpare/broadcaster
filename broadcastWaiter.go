package broadcaster

import "sync"

type BroadcastWaiter[T any] struct {
	listeners map[*Listener[T]]struct{}
	mutex    *sync.RWMutex
	closed    bool
}

func NewBroadcastWaiter[T any]() *BroadcastWaiter[T] {
	return &BroadcastWaiter[T]{
		listeners: make(map[*Listener[T]]struct{}),
		mutex:    new(sync.RWMutex),
	}
}

func (bc *BroadcastWaiter[T]) getListeners() map[*Listener[T]]struct{} {
	return bc.listeners
}

func (bc *BroadcastWaiter[T]) getMutex() *sync.RWMutex {
	return bc.mutex
}

func (bc *BroadcastWaiter[T]) Register(bufSize int) *Listener[T] {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	l := &Listener[T]{
		bc: bc,
		ch: make(chan Payload[T], bufSize),
	}

	if bc.closed {
		l.closed = true
		close(l.ch)
	} else {
		bc.listeners[l] = struct{}{}
	}

	return l
}

func (bc *BroadcastWaiter[T]) Send(payload T) Waiter[T] {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	w := &broadcastPayloadWaiter[T]{
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

func (bc *BroadcastWaiter[T]) Reset() {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    bc.resetNoLock()
}

func (bc *BroadcastWaiter[T]) resetNoLock() {
    for ch := range bc.listeners {
        ch.unregisterNoLock()
    }
}

func (bc *BroadcastWaiter[T]) Close() {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    if bc.closed {
        return
    }
    bc.closed = true

    bc.resetNoLock()
}

type BufBroadcastWaiter[T any] struct {
	listeners map[*Listener[T]]struct{}
	mutex    *sync.RWMutex
	closed    bool
	data      []T
}

func NewBufBroadcastWaiter[T any]() *BufBroadcastWaiter[T] {
	return &BufBroadcastWaiter[T]{
		listeners: make(map[*Listener[T]]struct{}),
		mutex:    new(sync.RWMutex),
	}
}

func (bc *BufBroadcastWaiter[T]) getListeners() map[*Listener[T]]struct{} {
	return bc.listeners
}

func (bc *BufBroadcastWaiter[T]) getMutex() *sync.RWMutex {
	return bc.mutex
}

func (bc *BufBroadcastWaiter[T]) registerNoLock(bufSize int) *Listener[T] {
	l := &Listener[T]{
		bc: bc,
		ch: make(chan Payload[T], bufSize),
	}

	if bc.closed {
		l.closed = true
		close(l.ch)
	} else {
		bc.listeners[l] = struct{}{}
	}

	bc.listeners[l] = struct{}{}
	return l
}

func (bc *BufBroadcastWaiter[T]) Register(bufSize int) *Listener[T] {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	return bc.registerNoLock(bufSize)
}

func (bc *BufBroadcastWaiter[T]) Send(payload T) Waiter[T] {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	bc.data = append(bc.data, payload)

	w := &broadcastPayloadWaiter[T]{
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

func (bc *BufBroadcastWaiter[T]) Data() []T {
	return bc.data
}

func (bc *BufBroadcastWaiter[T]) Connect(bufSize int) ([]T, *Listener[T]) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	return bc.Data(), bc.registerNoLock(bufSize)
}

func (bc *BufBroadcastWaiter[T]) Reset() {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    bc.resetNoLock()
}

func (bc *BufBroadcastWaiter[T]) resetNoLock() {
    for ch := range bc.listeners {
        ch.unregisterNoLock()
    }
    bc.data = nil
}

func (bc *BufBroadcastWaiter[T]) Close() {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()

    if bc.closed {
        return
    }
    bc.closed = true

    bc.resetNoLock()
}
