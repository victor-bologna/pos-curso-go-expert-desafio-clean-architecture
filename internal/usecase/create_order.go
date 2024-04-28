package usecase

import (
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/entity"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/event"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/database"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/pkg/events"
)

type InputOrderDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OutputOrderDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"finalPrice"`
}

type CreateOrderUseCase struct {
	OrderRepository database.OrderRepositoryInterface //Inversão de Controle
	OrderCreated    event.OrderCreatedEvent
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(orderRepository database.OrderRepositoryInterface, orderCreated event.OrderCreatedEvent,
	eventDispatcher events.EventDispatcherInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
		OrderCreated:    orderCreated,
		EventDispatcher: eventDispatcher,
	}
}

func (couc *CreateOrderUseCase) Execute(input InputOrderDTO) (OutputOrderDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice()
	if err := couc.OrderRepository.Save(&order); err != nil {
		return OutputOrderDTO{}, err
	}

	dto := OutputOrderDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
	couc.OrderCreated.SetPayload(dto)
	couc.EventDispatcher.Dispatch(couc.OrderCreated)
	return dto, nil
}
