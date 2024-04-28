package web

import (
	"encoding/json"
	"net/http"

	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/event"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/database"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/usecase"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/pkg/events"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   database.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
	GetAllOrdersEvent events.EventInterface
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository database.OrderRepositoryInterface,
	OrderCreatedEvent event.OrderCreatedEvent,
	GetAllOrdersEvent event.GetAllOrdersEvent,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
		GetAllOrdersEvent: GetAllOrdersEvent,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.InputOrderDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	getAllOrder := usecase.NewGetOrderUseCase(h.OrderRepository, h.GetAllOrdersEvent, h.EventDispatcher)

	output, err := getAllOrder.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
