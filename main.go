package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/segment3d-app/segment3d-be/api"
	db "github.com/segment3d-app/segment3d-be/db/sqlc"
	"github.com/segment3d-app/segment3d-be/util"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
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

	// rabbitmq
	rconn, err := amqp091.Dial(config.RabbitSource)
	if err != nil {
		log.Fatal("can't connect to rabbitmq")
	}
	ch, err := rconn.Channel()
	if err != nil {
		log.Fatal("can't create channel")
	}

	// server
	server, err := api.NewServer(&config, store, ch)
	if err != nil {
		log.Fatal("can't create server", err)
	}

	// start server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server", err)
	}
}
