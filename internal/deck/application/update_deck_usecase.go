package application

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck/domain"
)

type UpdateDeckInput struct {
	Name        string
	Format      string
	Description string
}

type UpdateDeckUsecase struct {
	repository DeckRepository
}

func NewUpdateDeckUsecase(repository DeckRepository) *UpdateDeckUsecase {
	return &UpdateDeckUsecase{repository: repository}
}

func (uc *UpdateDeckUsecase) Execute(ctx context.Context, id string, input UpdateDeckInput) (domain.Deck, error) {
	oldDeck, err := uc.repository.FindById(ctx, id)
	if err != nil {
		return domain.Deck{}, err
	}

	newDeck := &domain.Deck{Name: input.Name, Format: input.Format, Description: input.Description}
	deck, err := domain.UpdateDeck(oldDeck, *newDeck)
	if err != nil {
		return domain.Deck{}, err
	}

	updatedDeck, err := uc.repository.Update(ctx, id, deck)
	if err != nil {
		return domain.Deck{}, err
	}

	return updatedDeck, nil
}
