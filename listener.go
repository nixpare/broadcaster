package broadcaster

type Listener[T any] struct {
    bc *Broadcaster[T]
    ch chan Payload[T]
}

func (l *Listener[T]) unregisterNoLock() {
    delete(l.bc.listeners, l)
    close(l.ch)
}

func (l *Listener[T]) Unregister() {
    l.bc.lMutex.Lock()
    l.unregisterNoLock()
    l.bc.lMutex.Unlock()
}

func (l *Listener[T]) Ch() <-chan Payload[T] {
    return l.ch
}

func (l *Listener[T]) Get() T {
    payload := <- l.Ch()
    return payload.Get()
}
