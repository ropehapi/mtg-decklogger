package application

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck_version/domain.go"
)

type DeckVersionRepository interface {
	Create(ctx context.Context, deck domain.DeckVersion) error
	// Find(ctx context.Context) ([]domain.DeckVersion, error)
	// FindById(ctx context.Context, id string) (domain.DeckVersion, error)
	// Update(ctx context.Context, id string, deck domain.DeckVersion) (domain.DeckVersion, error)
	// Delete(ctx context.Context, id string) error
}

// DeckRepository is intentionally small: deck versions only need to verify
// that their parent deck exists.
type DeckRepository interface {
	ExistsByID(ctx context.Context, id string) (bool, error)
}
