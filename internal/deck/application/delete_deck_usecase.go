package application

import (
	"context"
)

type DeleteDeckUsecase struct {
	repository DeckRepository
}

func NewDeleteDeckUsecase(repository DeckRepository) *DeleteDeckUsecase {
	return &DeleteDeckUsecase{repository: repository}
}

func (uc *DeleteDeckUsecase) Execute(ctx context.Context, id string) error {
	return uc.repository.Delete(ctx, id)
}
