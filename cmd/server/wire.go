//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/event"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/database"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/web"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/usecase"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(database.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	event.NewGetOrders,
	wire.Bind(new(event.OrderCreatedEvent), new(*event.OrderCreated)),
	wire.Bind(new(event.GetAllOrdersEvent), new(*event.GetOrders)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(event.OrderCreatedEvent), new(*event.OrderCreated)),
)

var setGetOrdersEvent = wire.NewSet(
	event.NewGetOrders,
	wire.Bind(new(event.GetAllOrdersEvent), new(*event.GetOrders)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewGetOrdersUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.GetOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setGetOrdersEvent,
		usecase.NewGetOrderUseCase,
	)
	return &usecase.GetOrdersUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		setGetOrdersEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
