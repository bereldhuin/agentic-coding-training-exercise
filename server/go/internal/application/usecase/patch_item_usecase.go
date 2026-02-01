package usecase

import (
	"context"

	"github.com/hymaia/lebonpoint-app/server/go/internal/application/dto"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/entity"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
	apperrors "github.com/hymaia/lebonpoint-app/server/go/internal/shared/errors"
)

// PatchItemUseCase handles partial updates of items
type PatchItemUseCase struct {
	repo repository.ItemRepository
}

// NewPatchItemUseCase creates a new PatchItemUseCase
func NewPatchItemUseCase(repo repository.ItemRepository) *PatchItemUseCase {
	return &PatchItemUseCase{
		repo: repo,
	}
}

// Execute performs a partial update of an item
func (uc *PatchItemUseCase) Execute(ctx context.Context, id int64, req *dto.PatchItemRequest) (*entity.Item, error) {
	// Validate that at least one field is being updated
	if !req.HasChanges() {
		return nil, apperrors.NewBadRequestError("no fields to update")
	}

	// Get existing item
	item, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Title != nil {
		if len(*req.Title) < 3 || len(*req.Title) > 200 {
			return nil, apperrors.NewValidationError("title must be between 3 and 200 characters", map[string]interface{}{"title": "invalid length"})
		}
		item.Title = *req.Title
	}

	if req.Description != nil {
		item.Description = *req.Description
	}

	if req.PriceCents != nil {
		if *req.PriceCents < 0 {
			return nil, apperrors.NewValidationError("price_cents must be >= 0", map[string]interface{}{"priceCents": "invalid value"})
		}
		item.PriceCents = *req.PriceCents
	}

	if req.Category != nil {
		item.Category = *req.Category
	}

	if req.Condition != nil {
		if err := req.Condition.Validate(); err != nil {
			return nil, apperrors.NewValidationError("invalid condition", map[string]interface{}{"condition": err.Error()})
		}
		item.Condition = *req.Condition
	}

	if req.Status != nil {
		item.Status = *req.Status
	}

	if req.IsFeatured != nil {
		item.IsFeatured = *req.IsFeatured
	}

	if req.City != nil {
		item.City = *req.City
	}

	if req.PostalCode != nil {
		item.PostalCode = *req.PostalCode
	}

	if req.Country != nil {
		item.Country = *req.Country
	}

	if req.DeliveryAvailable != nil {
		item.DeliveryAvailable = *req.DeliveryAvailable
	}

	if req.Images != nil {
		item.Images = *req.Images
	}

	// Store updated item
	if err := uc.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}
