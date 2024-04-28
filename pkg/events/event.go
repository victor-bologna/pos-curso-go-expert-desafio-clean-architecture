package events

import (
	"sync"
	"time"
)

type EventInterface interface { //Evento
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

type EventHandlerInterface interface { //Operação quando evento for chamado (Tem que ter um handle).
	Handle(EventInterface EventInterface, WaitGroup *sync.WaitGroup)
}

type EventDispatcherInterface interface { // Cria o evento com um nome.
	Register(name string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error // Dispara o(s) evento(s) registrados acima
	Remove(name string, handler EventHandlerInterface) error
	Has(name string, handler EventHandlerInterface) bool
	Clear()
}
