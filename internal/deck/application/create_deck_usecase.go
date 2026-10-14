package application

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
)

type CreateDeckInput struct {
	Name        string
	Format      string
	Description string
}

type CreateDeckUsecase struct {
	repository DeckRepository
}

func NewCreateDeckUsecase(repository DeckRepository) *CreateDeckUsecase {
	return &CreateDeckUsecase{repository: repository}
}

func (uc *CreateDeckUsecase) Execute(ctx context.Context, input CreateDeckInput) (domain.Deck, error) {
	deck, err := domain.NewDeck(input.Name, input.Format, input.Description)
	if err != nil {
		return domain.Deck{}, err
	}

	if err := uc.repository.Create(ctx, deck); err != nil {
		return domain.Deck{}, err
	}

	return deck, nil
}
