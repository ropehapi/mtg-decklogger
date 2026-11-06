package main

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	deckhttp "github.com/ropehapi/mtg-decklogger/internal/deck/adapter/inbound/http"
	deckversionhttp "github.com/ropehapi/mtg-decklogger/internal/deck_version/adapter/inbound/http"
	"github.com/ropehapi/mtg-decklogger/internal/deck/adapter/outbound/mongodb"
	deckversionmongodb "github.com/ropehapi/mtg-decklogger/internal/deck_version/adapter/outbound/mongodb"
	"github.com/ropehapi/mtg-decklogger/internal/deck/application"
	deckversionapplication "github.com/ropehapi/mtg-decklogger/internal/deck_version/application"
	"github.com/ropehapi/mtg-decklogger/pkg/config"
	"github.com/ropehapi/mtg-decklogger/pkg/database"
	"github.com/ropehapi/mtg-decklogger/pkg/logger"
)

// TODO: Add observability
// TODO: Add swagger docs
// TODO: Add unit tests
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

	deckCollection := mongoConnection.Database("deck").Collection("deckData")
	deckVersionCollection := mongoConnection.Database("deck_version").Collection("deckData")

	deckRepository := mongodb.NewRepository(deckCollection)
	deckVersionRepository := deckversionmongodb.NewRepository(deckVersionCollection)

	createDeckUsecase := application.NewCreateDeckUsecase(deckRepository)
	listDeckUsecase := application.NewListDecksUsecase(deckRepository)
	updateDeckUsecase := application.NewUpdateDeckUsecase(deckRepository)
	deleteDeckUsecase := application.NewDeleteDeckUsecase(deckRepository)

	createDeckVersionUsecase := deckversionapplication.NewCreateDeckVersionUsecase(deckVersionRepository)

	deckHandler := deckhttp.NewHandler(createDeckUsecase, listDeckUsecase, updateDeckUsecase, deleteDeckUsecase)
	deckVersionHandler := deckversionhttp.NewHandler(createDeckVersionUsecase)

	r := gin.Default()
	r.POST("/decks", deckHandler.Create)
	r.GET("/decks", deckHandler.List)
	r.PUT("/decks/:id", deckHandler.Update)
	r.DELETE("/decks/:id", deckHandler.Delete)

	r.POST("/deck_version", deckVersionHandler.Create)
	address := ":" + cfg.ServerPort
	r.Run(address)
}
