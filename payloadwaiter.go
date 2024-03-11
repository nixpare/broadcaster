package broadcaster

import "sync"

type Payload[T any] interface {
    Done()
	Data() T
	Get() T
}

type Waiter[T any] interface {
	Wait()
}

type payloadWaiter[T any] struct {
	listeners []*Listener[T]
	payload   T
	wg        sync.WaitGroup
}

func (pw *payloadWaiter[T]) Wait() {
	pw.wg.Wait()
}

func (pw *payloadWaiter[T]) Done() {
	pw.wg.Done()
}

func (pw *payloadWaiter[T]) Data() T {
	return pw.payload
}

func (pw *payloadWaiter[T]) Get() T {
    defer pw.Done()
	return pw.Data()
}
