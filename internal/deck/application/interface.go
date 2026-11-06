package application

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
)

type DeckRepository interface {
	Create(ctx context.Context, deck domain.Deck) error
	Find(ctx context.Context) ([]domain.Deck, error)
	FindById(ctx context.Context, id string) (domain.Deck, error)
	Update(ctx context.Context, id string, deck domain.Deck) (domain.Deck, error)
	Delete(ctx context.Context, id string) error
}
