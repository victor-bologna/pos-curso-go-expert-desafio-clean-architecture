// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from wire.go:

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreated, eventDispatcher)
	return createOrderUseCase
}

func NewGetOrdersUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.GetOrdersUseCase {
	orderRepository := database.NewOrderRepository(db)
	getOrders := event.NewGetOrders()
	getOrdersUseCase := usecase.NewGetOrderUseCase(orderRepository, getOrders, eventDispatcher)
	return getOrdersUseCase
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	getOrders := event.NewGetOrders()
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, orderCreated, getOrders)
	return webOrderHandler
}

// wire.go:

var setOrderRepositoryDependency = wire.NewSet(database.NewOrderRepository, wire.Bind(new(database.OrderRepositoryInterface), new(*database.OrderRepository)))

var setEventDispatcherDependency = wire.NewSet(events.NewEventDispatcher, event.NewOrderCreated, event.NewGetOrders, wire.Bind(new(event.OrderCreatedEvent), new(*event.OrderCreated)), wire.Bind(new(event.GetAllOrdersEvent), new(*event.GetOrders)), wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)))

var setOrderCreatedEvent = wire.NewSet(event.NewOrderCreated, wire.Bind(new(event.OrderCreatedEvent), new(*event.OrderCreated)))

var setGetOrdersEvent = wire.NewSet(event.NewGetOrders, wire.Bind(new(event.GetAllOrdersEvent), new(*event.GetOrders)))
