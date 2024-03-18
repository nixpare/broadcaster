package broadcaster

type Listener[T any] struct {
	bc broadcaster[*Listener[T]]
	ch chan Payload[T]
}

func (l *Listener[T]) Ch() <-chan Payload[T] {
	return l.ch
}

func (l *Listener[T]) Get() T {
	payload := <-l.Ch()
	return payload.Get()
}

func (l *Listener[T]) unregisterNoLock() {
	delete(l.bc.getListeners(), l)
	close(l.ch)
}

func (l *Listener[T]) Unregister() {
	mutex := l.bc.getMutex()
    mutex.Lock()
    defer mutex.Unlock()
	
	l.unregisterNoLock()
}
