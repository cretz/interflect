package interflect

import "reflect"

// TODO(cretz): Surely will change how this looks as we work on it
type Scheduler interface {
	StartCoroutine(func(Coroutine)) Coroutine
}

type Coroutine interface {
	Yielded() bool
	Select(cases []reflect.SelectCase) (chosen int, recv reflect.Value, recvOK bool)
	Send(ch reflect.Value, arg reflect.Value)
	Recv(ch reflect.Value) (x reflect.Value, ok bool)
}

var SystemScheduler Scheduler = systemScheduler{}

type systemScheduler struct{}

func (systemScheduler) StartCoroutine(fn func(Coroutine)) Coroutine {
	coro := &systemCoroutine{}
	go func() { fn(coro) }()
	return coro
}

type systemCoroutine struct{}

func (systemCoroutine) Yielded() bool {
	return false
}

func (systemCoroutine) Select(cases []reflect.SelectCase) (chosen int, recv reflect.Value, recvOK bool) {
	return reflect.Select(cases)
}

func (systemCoroutine) Send(ch reflect.Value, arg reflect.Value) {
	ch.Send(arg)
}

func (systemCoroutine) Recv(ch reflect.Value) (x reflect.Value, ok bool) {
	return ch.Recv()
}

// TODO(cretz): This
var CooperativeScheduler Scheduler = nil
