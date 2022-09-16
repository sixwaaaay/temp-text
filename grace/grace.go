package grace

import "os"

type Endless struct {
	Running func()                // endless run
	OnEvent func()                // do when event
	Event   func() chan os.Signal // signal channel for event trigger
}

func NewEndless(running func(), onEvent func(), event func() chan os.Signal) *Endless {
	return &Endless{Running: running, OnEvent: onEvent, Event: event}
}

func (receiver *Endless) Run() {
	go receiver.Running()
	<-receiver.Event()
	receiver.OnEvent()
}
