package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(name string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[name]; ok {
		for _, h := range ed.handlers[name] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[name] = append(ed.handlers[name], handler)
	return nil
}

func (ed *EventDispatcher) Clear() {
	ed.handlers = make(map[string][]EventHandlerInterface)
}

func (ed *EventDispatcher) Has(name string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[name]; ok {
		for _, h := range ed.handlers[name] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, h := range handlers {
			wg.Add(1)
			go h.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (ed *EventDispatcher) Remove(name string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[name]; ok {
		for i, h := range ed.handlers[name] {
			if h == handler {
				ed.handlers[name] = append(ed.handlers[name][:i], ed.handlers[name][i+1:]...)
				return nil
			}
		}
	}
	return nil
}
