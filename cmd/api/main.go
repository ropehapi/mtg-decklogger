package main

import (
	"log"

	deckhttp "github.com/ropehapi/mtg-decklogger/internal/deck/adapter/inbound/http"
	"github.com/ropehapi/mtg-decklogger/internal/deck/adapter/outbound/mongodb"
	"github.com/ropehapi/mtg-decklogger/internal/deck/application"
	"github.com/ropehapi/mtg-decklogger/pkg/config"
	"github.com/ropehapi/mtg-decklogger/pkg/database"
	"github.com/gin-gonic/gin"
)

//TODO: Add structured logs
//TODO: Add observability
//TODO: Add swagger docs
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	mongoConnection, err := database.ConnectMongo(cfg.MongoUri)
	if err != nil {
		log.Fatalf("Error connecting to Mongo database: %v", err)
	}
	defer func() {
		if err := mongoConnection.Close(); err != nil {
			log.Printf("Error disconnecting from Mongo database: %v", err)
		}
	}()

	collection := mongoConnection.Database("deck").Collection("deckData")
	repository := mongodb.NewRepository(collection)

	createDeckUsecase := application.NewCreateDeckUsecase(repository)
	listDeckUsecase := application.NewListDecksUsecase(repository)

	handler := deckhttp.NewHandler(createDeckUsecase, listDeckUsecase)

	r := gin.Default()
	r.POST("/decks", handler.Create)
	r.GET("/decks", handler.List)

	address := ":" + cfg.ServerPort
	r.Run(address)
}
