package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ropehapi/mtg-decklogger/internal/deck/application"
	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
)

type CreateDeckUsecase interface {
	Execute(ctx context.Context, input application.CreateDeckInput) (domain.Deck, error)
}

type ListDecksUsecase interface {
	Execute(ctx context.Context) ([]domain.Deck, error)
}

type Handler struct {
	createDeckUsecase CreateDeckUsecase
	listDecksUsecase ListDecksUsecase
}

type createDeckRequest struct {
	Name        string `json:"name"`
	Format      string `json:"format"`
	Description string `json:"description"`
}

//TODO: Create default Http response struct
type deckResponse struct {
	Name        string `json:"name"`
	Format      string `json:"format"`
	Description string `json:"description"`
}

//TODO: Create default http error response struct
type errorResponse struct {
	Message string `json:"message"`
}

func NewHandler(createDeckUsecase CreateDeckUsecase, listDecksUsecase ListDecksUsecase) *Handler {
	return &Handler{createDeckUsecase: createDeckUsecase, listDecksUsecase: listDecksUsecase}
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body createDeckRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid request body"})
		return
	}

	deck, err := handler.createDeckUsecase.Execute(r.Context(), application.CreateDeckInput{
		Name:        body.Name,
		Format:      body.Format,
		Description: body.Description,
	})
	if errors.Is(err, domain.ErrDeckNameRequired) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Message: fmt.Sprintf("Error creating deck: %v", err.Error())})
		return
	}

	writeJSON(w, http.StatusCreated, deckResponse{
		Name:        deck.Name,
		Format:      deck.Format,
		Description: deck.Description,
	})
}

func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	decks, err := handler.listDecksUsecase.Execute(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Message: fmt.Sprintf("Error listing decks: %v", err.Error())})
		return
	}

	writeJSON(w, http.StatusCreated, decks)
}

func writeJSON(response http.ResponseWriter, status int, body any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(body)
}
