package main

import (
	"database/sql"
	"fmt"
	graphqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"github.com/winstonjr/goexpert-clean-arch/configs"
	"github.com/winstonjr/goexpert-clean-arch/internal/event/handler"
	"github.com/winstonjr/goexpert-clean-arch/internal/infra/graph"
	"github.com/winstonjr/goexpert-clean-arch/internal/infra/grpc/pb"
	"github.com/winstonjr/goexpert-clean-arch/internal/infra/grpc/service"
	"github.com/winstonjr/goexpert-clean-arch/internal/infra/web/webserver"
	"github.com/winstonjr/goexpert-clean-arch/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(conf.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(conf)

	eventDispatcher := events.NewEventDispatcher()
	err = eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	if err != nil {
		panic(err)
	}

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	ws := webserver.NewWebServer(conf.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	ws.AddHandler(http.MethodPost, "/order", webOrderHandler.Create)
	ws.AddHandler(http.MethodGet, "/order", webOrderHandler.ListAll)
	fmt.Println("Starting web server on port", conf.WebServerPort)
	go ws.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", conf.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphqlHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", conf.GraphQLServerPort)
	http.ListenAndServe(":"+conf.GraphQLServerPort, nil)
}

func getRabbitMQChannel(conf *configs.Conf) *amqp.Channel {
	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/", conf.RabbitMQUser, conf.RabbitMQPassword, conf.RabbitMQHost, conf.RabbitMQPort))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
