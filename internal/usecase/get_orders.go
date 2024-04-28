package usecase

import (
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/event"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/database"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/pkg/events"
)

type GetOrdersUseCase struct {
	OrderRepository database.OrderRepositoryInterface
	GetAllOrders    event.GetAllOrdersEvent
	EventDispatcher events.EventDispatcherInterface
}

func NewGetOrderUseCase(orderRepository database.OrderRepositoryInterface, getAllOrders event.GetAllOrdersEvent,
	eventDispatcher events.EventDispatcherInterface) *GetOrdersUseCase {
	return &GetOrdersUseCase{
		OrderRepository: orderRepository,
		GetAllOrders:    getAllOrders,
		EventDispatcher: eventDispatcher,
	}
}

func (gaous *GetOrdersUseCase) Execute() ([]OutputOrderDTO, error) {
	orders, err := gaous.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var outputOrderDTOs []OutputOrderDTO
	for _, order := range orders {
		outputOrderDTOs = append(outputOrderDTOs, OutputOrderDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	gaous.GetAllOrders.SetPayload(outputOrderDTOs)
	gaous.EventDispatcher.Dispatch(gaous.GetAllOrders)
	return outputOrderDTOs, nil
}
