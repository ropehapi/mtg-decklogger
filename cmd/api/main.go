package main

import (
	"log"
	"net/http"

	deckhttp "github.com/ropehapi/mtg-decklogger/internal/deck/adapter/inbound/http"
	"github.com/ropehapi/mtg-decklogger/internal/deck/adapter/outbound/mongodb"
	"github.com/ropehapi/mtg-decklogger/internal/deck/application"
	"github.com/ropehapi/mtg-decklogger/pkg/config"
	"github.com/ropehapi/mtg-decklogger/pkg/database"
)

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

	//TODO: Refactor to CHI or Gin, whatever fits best
	router := http.NewServeMux()
	router.HandleFunc("POST /decks", handler.Create)
	router.HandleFunc("GET /decks", handler.List)

	address := ":" + cfg.ServerPort
	log.Printf("API listening on %s", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Printf("HTTP server stopped: %v", err)
	}
}
