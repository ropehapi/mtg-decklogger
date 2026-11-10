package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	deckhttp "github.com/ropehapi/mtg-decklogger/internal/deck/adapter/inbound/http"
	deckmongodb "github.com/ropehapi/mtg-decklogger/internal/deck/adapter/outbound/mongodb"
	deckapplication "github.com/ropehapi/mtg-decklogger/internal/deck/application"
	deckversionhttp "github.com/ropehapi/mtg-decklogger/internal/deck_version/adapter/inbound/http"
	deckversionmongodb "github.com/ropehapi/mtg-decklogger/internal/deck_version/adapter/outbound/mongodb"
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

	db := mongoConnection.Database("mtg_decklogger")
	deckCollection := db.Collection("decks")
	deckVersionCollection := db.Collection("deck_versions")

	deckRepository := deckmongodb.NewRepository(deckCollection)
	deckVersionRepository := deckversionmongodb.NewRepository(deckVersionCollection)
	if err := deckVersionRepository.EnsureIndexes(context.Background()); err != nil {
		slog.Error(fmt.Sprintf("Error creating deck version indexes: %v", err))
		return
	}

	createDeckUsecase := deckapplication.NewCreateDeckUsecase(deckRepository)
	listDeckUsecase := deckapplication.NewListDecksUsecase(deckRepository)
	updateDeckUsecase := deckapplication.NewUpdateDeckUsecase(deckRepository)
	deleteDeckUsecase := deckapplication.NewDeleteDeckUsecase(deckRepository)

	createDeckVersionUsecase := deckversionapplication.NewCreateDeckVersionUsecase(deckVersionRepository, deckRepository)

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
