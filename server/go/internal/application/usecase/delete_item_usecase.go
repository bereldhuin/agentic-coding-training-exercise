package usecase

import (
	"context"

	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
)

// DeleteItemUseCase handles item deletion
type DeleteItemUseCase struct {
	repo repository.ItemRepository
}

// NewDeleteItemUseCase creates a new DeleteItemUseCase
func NewDeleteItemUseCase(repo repository.ItemRepository) *DeleteItemUseCase {
	return &DeleteItemUseCase{
		repo: repo,
	}
}

// Execute deletes an item by ID
func (uc *DeleteItemUseCase) Execute(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
