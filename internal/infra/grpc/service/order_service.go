package service

import (
	"context"

	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/grpc/pb"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	GetOrdersUseCase   usecase.GetOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase,
	getOrdersUseCase usecase.GetOrdersUseCase) *OrderService {
	return &OrderService{CreateOrderUseCase: createOrderUseCase, GetOrdersUseCase: getOrdersUseCase}
}

func (os *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.InputOrderDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	output, err := os.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (os *OrderService) GetAllOrders(ctx context.Context, in *pb.Blank) (*pb.OrderResponseList, error) {
	output, err := os.GetOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var ordersResponse []*pb.CreateOrderResponse
	for _, order := range output {
		ordersResponse = append(ordersResponse, &pb.CreateOrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}
	return &pb.OrderResponseList{Response: ordersResponse}, nil
}
