package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rabbitmq/amqp091-go"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/configs"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/event/handler"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/graph"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/grpc/pb"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/grpc/service"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/infra/web/webserver"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/internal/usecase"
	"github.com/victor-bologna/pos-curso-go-expert-clean-architecture/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := configs.LoadConfig(".")

	db := loadDB(config)
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(config)

	eventDispatcher := loadDispatcher(rabbitMQChannel)

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	getOrdersUseCase := NewGetOrdersUseCase(db, eventDispatcher)

	executeWebServer(config, db, eventDispatcher)

	executeGRPCServer(createOrderUseCase, getOrdersUseCase, config)

	executeGraphQLServer(createOrderUseCase, getOrdersUseCase, config)
}

func executeGraphQLServer(createOrderUseCase *usecase.CreateOrderUseCase, getOrdersUseCase *usecase.GetOrdersUseCase, config *configs.Conf) {
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		GetOrdersUseCase:   *getOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", config.API.GraphQLServerPort)
	http.ListenAndServe(":"+config.API.GraphQLServerPort, nil)
}

func executeGRPCServer(createOrderUseCase *usecase.CreateOrderUseCase, getOrdersUseCase *usecase.GetOrdersUseCase, config *configs.Conf) {
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *getOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", config.API.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.API.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)
}

func executeWebServer(config *configs.Conf, db *sql.DB, eventDispatcher *events.EventDispatcher) {
	webserver := webserver.NewWebServer(config.API.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/orders", webOrderHandler.GetAll)
	fmt.Println("Starting web server on port", config.API.WebServerPort)
	go webserver.Start()
}

func loadDispatcher(rabbitMQChannel *amqp091.Channel) *events.EventDispatcher {
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.Register("Get Orders", &handler.GetOrdersHandler{})
	return eventDispatcher
}

func loadDB(config *configs.Conf) *sql.DB {
	db, err := sql.Open(config.DB.Driver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name))
	if err != nil {
		panic(err)
	}
	return db
}

func getRabbitMQChannel(config *configs.Conf) *amqp091.Channel {
	conn, err := amqp091.Dial(config.API.Rabbit)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
