package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	deckhttp "github.com/ropehapi/mtg-decklogger/internal/deck/adapter/inbound/http"
	"github.com/ropehapi/mtg-decklogger/internal/deck/adapter/outbound/mongodb"
	"github.com/ropehapi/mtg-decklogger/internal/deck/application"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	//TODO: Segregate this connection in another file
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://root:example@localhost:27017/?authSource=admin"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("disconnect MongoDB: %v", err)
		}
	}()

	pingContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(pingContext, nil); err != nil {
		log.Fatal(err)
	}

	collection := client.Database("deck").Collection("deckData")
	repository := mongodb.NewRepository(collection)

	createDeckUsecase := application.NewCreateDeckUsecase(repository)
	listDeckUsecase := application.NewListDecksUsecase(repository)

	handler := deckhttp.NewHandler(createDeckUsecase, listDeckUsecase)

	//TODO: Refactor to CHI or Gin, whatever fits best
	router := http.NewServeMux()
	router.HandleFunc("POST /decks", handler.Create)
	router.HandleFunc("GET /decks", handler.List)

	log.Println("API listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
