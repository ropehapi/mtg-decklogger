package application

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck_version/domain.go"
)

type CreateDeckVersionInput struct {
	DeckId      string
	Version        int
	Description string
}

type CreateDeckVersionUsecase struct {
	repository DeckVersionRepository
}

func NewCreateDeckVersionUsecase(repository DeckVersionRepository) *CreateDeckVersionUsecase {
	return &CreateDeckVersionUsecase{repository: repository}
}

func (uc *CreateDeckVersionUsecase) Execute(ctx context.Context, input CreateDeckVersionInput) (domain.DeckVersion, error) {
	deckVersion, err := domain.NewDeckVersion(input.DeckId, input.Version, input.Description)
	if err != nil {
		return domain.DeckVersion{}, err
	}

	if err := uc.repository.Create(ctx, deckVersion); err != nil {
		return domain.DeckVersion{}, err
	}

	return deckVersion, nil
}
