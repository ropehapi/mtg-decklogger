package application

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
)

type DeckRepository interface {
	Create(ctx context.Context, deck domain.Deck) error
	Find(ctx context.Context) ([]domain.Deck, error)
}
