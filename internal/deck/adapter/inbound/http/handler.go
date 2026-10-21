package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ropehapi/mtg-decklogger/internal/deck/application"
	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
	"github.com/ropehapi/mtg-decklogger/pkg/response"
)

type CreateDeckUsecase interface {
	Execute(ctx context.Context, input application.CreateDeckInput) (domain.Deck, error)
}

type ListDecksUsecase interface {
	Execute(ctx context.Context) ([]domain.Deck, error)
}

type Handler struct {
	createDeckUsecase CreateDeckUsecase
	listDecksUsecase  ListDecksUsecase
}

type createDeckRequest struct {
	Name        string `json:"name"`
	Format      string `json:"format"`
	Description string `json:"description"`
}

func NewHandler(createDeckUsecase CreateDeckUsecase, listDecksUsecase ListDecksUsecase) *Handler {
	return &Handler{createDeckUsecase: createDeckUsecase, listDecksUsecase: listDecksUsecase}
}

func (handler *Handler) Create(c *gin.Context) {
	var body createDeckRequest
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	deck, err := handler.createDeckUsecase.Execute(c, application.CreateDeckInput{
		Name:        body.Name,
		Format:      body.Format,
		Description: body.Description,
	})
	if errors.Is(err, domain.ErrDeckNameRequired) {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Error creating deck: %v", err.Error()))
		return
	}
	
	response.Success(c, http.StatusCreated, deck)
}

func (handler *Handler) List(c *gin.Context) {
	decks, err := handler.listDecksUsecase.Execute(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Error listing decks: %v", err.Error()))
	}

	response.Success(c, http.StatusOK, decks)
}
