package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
	"github.com/segment3d-app/segment3d-be/api"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
	"github.com/segment3d-app/segment3d-be/util"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	"github.com/wagslane/go-rabbitmq"
)

type Message struct {
	URL string `json:"url"`
}

// @title Segment3d App API Documentation
// @version 1.0
// @description This is a documentation for Segment3d App API

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config")
	}

	// postgresql
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	store := db.NewStore(conn)

	// rabbit_mq
	rConn, err := rabbitmq.NewConn(
		config.RabbitSource,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal("can't connect to message broker queue", err)
	}
	defer conn.Close()

	publisher, err := rabbitmq.NewPublisher(
		rConn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("direct"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		// rabbitmq.WithPublisherOptionsExchangeDurable(&rabbitmq.PublisherOptions{
		// 	ExchangeOptions: rabbitmq.ExchangeOptions{Name: "direct",
		// 		Kind:       "direct",
		// 		Durable:    true,
		// 		AutoDelete: false,
		// 		Internal:   false,
		// 		NoWait:     false,
		// 		Passive:    false,
		// 		Args:       nil,
		// 		Declare:    false,
		// 	},
		// 	Logger:      rabbitmq.Logger{},
		// 	ConfirmMode: false,
		// }),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	// consumer, err := rabbitmq.NewConsumer(
	// 	rConn,
	// 	func(d rabbitmq.Delivery) rabbitmq.Action {
	// 		var msg Message
	// 		err := json.Unmarshal(d.Body, &msg)
	// 		if err != nil {
	// 			log.Printf("error unmarshalling message: %v", err)
	// 			return rabbitmq.NackDiscard
	// 		}

	// 		log.Printf("consumed URL: %v", msg.URL)
	// 		return rabbitmq.Ack
	// 	},
	// 	"splat_generation_queue",
	// 	rabbitmq.WithConsumerOptionsRoutingKey("splat_generation"),
	// 	rabbitmq.WithConsumerOptionsExchangeName("direct"),
	// 	rabbitmq.WithConsumerOptionsExchangeDeclare,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer consumer.Close()

	msg := Message{URL: "http://example.com"}
	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	err = publisher.Publish(
		[]byte(jsonMessage),
		[]string{"splat_generation"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("direct"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
	if err != nil {
		log.Println(err)
	}

	// server
	server, err := api.NewServer(&config, store)
	if err != nil {
		log.Fatal("can't create server", err)
	}

	// start server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server", err)
	}
}
