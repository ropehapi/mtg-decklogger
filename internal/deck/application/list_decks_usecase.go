package application

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
)

type ListDecksUsecase struct {
	repository DeckRepository
}

func NewListDecksUsecase(repository DeckRepository) *ListDecksUsecase {
	return &ListDecksUsecase{repository: repository}
}

func (uc *ListDecksUsecase) Execute(ctx context.Context) ([]domain.Deck, error) {
	results, err := uc.repository.Find(ctx)
	if err != nil {
		return nil, err
	}

	var decks []domain.Deck
	for _, result := range results {
		deck := domain.Deck{
			Name:        result.Name,
			Format:      result.Format,
			Description: result.Description,
		}
		decks = append(decks, deck)
	}

	return decks, nil
}
