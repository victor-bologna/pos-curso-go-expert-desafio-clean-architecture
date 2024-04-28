package event

import (
	"time"

	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/pkg/events"
)

type OrderCreatedEvent events.EventInterface

type GetOrders struct {
	Name    string
	Payload interface{}
}

func NewGetOrders() *GetOrders {
	return &GetOrders{Name: "Get all orders"}
}

func (g *GetOrders) GetName() string {
	return g.Name
}
func (g *GetOrders) GetDateTime() time.Time {
	return time.Now()
}
func (g *GetOrders) GetPayload() interface{} {
	return g.Payload
}
func (g *GetOrders) SetPayload(payload interface{}) {
	g.Payload = payload
}
