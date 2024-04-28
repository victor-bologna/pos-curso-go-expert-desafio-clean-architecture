package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/pkg/events"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp091.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp091.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{RabbitMQChannel: rabbitMQChannel}
}

func (och *OrderCreatedHandler) Handle(EventInterface events.EventInterface, WaitGroup *sync.WaitGroup) {
	defer WaitGroup.Done()
	fmt.Printf("Order created: %v", EventInterface.GetPayload())
	jsonOutput, _ := json.Marshal(EventInterface.GetPayload())

	msgRabbitMq := amqp091.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	och.RabbitMQChannel.PublishWithContext(
		context.TODO(),
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitMq,  // message to publish
	)
}
