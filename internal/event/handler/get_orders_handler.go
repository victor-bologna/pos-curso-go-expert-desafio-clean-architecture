package handler

import (
	"fmt"
	"sync"

	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/pkg/events"
)

type GetOrdersHandler struct {
	event string
}

func NewGetOrdersHandler() *GetOrdersHandler {
	return &GetOrdersHandler{event: "Get all orders handler"}
}

func (gaoh *GetOrdersHandler) Handle(EventInterface events.EventInterface, WaitGroup *sync.WaitGroup) {
	defer WaitGroup.Done()
	fmt.Printf("Get all orders: %v\n", EventInterface.GetPayload())
}
