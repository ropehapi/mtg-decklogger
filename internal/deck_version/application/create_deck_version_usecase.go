package application

import (
	"context"

	"github.com/ropehapi/mtg-decklogger/internal/deck_version/domain.go"
)

type CreateDeckVersionInput struct {
	DeckId      string
	Version     int
	Description string
}

type CreateDeckVersionUsecase struct {
	repository     DeckVersionRepository
	deckRepository DeckRepository
}

func NewCreateDeckVersionUsecase(repository DeckVersionRepository, deckRepository DeckRepository) *CreateDeckVersionUsecase {
	return &CreateDeckVersionUsecase{
		repository:     repository,
		deckRepository: deckRepository,
	}
}

func (uc *CreateDeckVersionUsecase) Execute(ctx context.Context, input CreateDeckVersionInput) (domain.DeckVersion, error) {
	deckVersion, err := domain.NewDeckVersion(input.DeckId, input.Version, input.Description)
	if err != nil {
		return domain.DeckVersion{}, err
	}

	exists, err := uc.deckRepository.ExistsByID(ctx, deckVersion.DeckId)
	if err != nil {
		return domain.DeckVersion{}, err
	}
	if !exists {
		return domain.DeckVersion{}, domain.ErrDeckNotFound
	}

	if err := uc.repository.Create(ctx, deckVersion); err != nil {
		return domain.DeckVersion{}, err
	}

	return deckVersion, nil
}
