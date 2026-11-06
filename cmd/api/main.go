package main

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	deckhttp "github.com/ropehapi/mtg-decklogger/internal/deck/adapter/inbound/http"
	"github.com/ropehapi/mtg-decklogger/internal/deck/adapter/outbound/mongodb"
	"github.com/ropehapi/mtg-decklogger/internal/deck/application"
	"github.com/ropehapi/mtg-decklogger/pkg/config"
	"github.com/ropehapi/mtg-decklogger/pkg/database"
	"github.com/ropehapi/mtg-decklogger/pkg/logger"
)

// TODO: Add observability
// TODO: Add swagger docs
func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading configurations: %v", err))
	}

	mongoConnection, err := database.ConnectMongo(cfg.MongoUri)
	if err != nil {
		slog.Error(fmt.Sprintf("Error connecting to Mongo database: %v", err))
	}
	defer func() {
		if err := mongoConnection.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error disconnecting from Mongo database: %v", err))
		}
	}()

	collection := mongoConnection.Database("deck").Collection("deckData")
	repository := mongodb.NewRepository(collection)

	createDeckUsecase := application.NewCreateDeckUsecase(repository)
	listDeckUsecase := application.NewListDecksUsecase(repository)
	updateDeckUsecase := application.NewUpdateDeckUsecase(repository)
	deleteDeckUsecase := application.NewDeleteDeckUsecase(repository)

	handler := deckhttp.NewHandler(createDeckUsecase, listDeckUsecase, updateDeckUsecase, deleteDeckUsecase)

	r := gin.Default()
	r.POST("/decks", handler.Create)
	r.GET("/decks", handler.List)
	r.PUT("/decks/:id", handler.Update)
	r.DELETE("/decks/:id", handler.Delete)

	address := ":" + cfg.ServerPort
	r.Run(address)
}
