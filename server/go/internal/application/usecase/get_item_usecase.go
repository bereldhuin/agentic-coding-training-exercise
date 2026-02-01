package usecase

import (
	"context"

	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/entity"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
)

// GetItemUseCase handles retrieving a single item by ID
type GetItemUseCase struct {
	repo repository.ItemRepository
}

// NewGetItemUseCase creates a new GetItemUseCase
func NewGetItemUseCase(repo repository.ItemRepository) *GetItemUseCase {
	return &GetItemUseCase{
		repo: repo,
	}
}

// Execute retrieves an item by ID
func (uc *GetItemUseCase) Execute(ctx context.Context, id int64) (*entity.Item, error) {
	return uc.repo.FindByID(ctx, id)
}
