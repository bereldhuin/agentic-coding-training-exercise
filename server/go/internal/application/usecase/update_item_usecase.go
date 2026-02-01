package usecase

import (
	"context"

	"github.com/hymaia/lebonpoint-app/server/go/internal/application/dto"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/entity"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
	apperrors "github.com/hymaia/lebonpoint-app/server/go/internal/shared/errors"
)

// UpdateItemUseCase handles full replacement updates of items
type UpdateItemUseCase struct {
	repo repository.ItemRepository
}

// NewUpdateItemUseCase creates a new UpdateItemUseCase
func NewUpdateItemUseCase(repo repository.ItemRepository) *UpdateItemUseCase {
	return &UpdateItemUseCase{
		repo: repo,
	}
}

// Execute performs a full update of an item
func (uc *UpdateItemUseCase) Execute(ctx context.Context, id int64, req *dto.UpdateItemRequest) (*entity.Item, error) {
	// Validate request
	if req.Title == "" {
		return nil, apperrors.NewValidationError("title is required", map[string]interface{}{"title": "required field"})
	}
	if len(req.Title) < 3 || len(req.Title) > 200 {
		return nil, apperrors.NewValidationError("title must be between 3 and 200 characters", map[string]interface{}{"title": "invalid length"})
	}
	if req.PriceCents < 0 {
		return nil, apperrors.NewValidationError("price_cents must be >= 0", map[string]interface{}{"priceCents": "invalid value"})
	}
	if err := req.Condition.Validate(); err != nil {
		return nil, apperrors.NewValidationError("invalid condition", map[string]interface{}{"condition": err.Error()})
	}
	if req.Country == "" {
		return nil, apperrors.NewValidationError("country is required", map[string]interface{}{"country": "required field"})
	}

	// Check if item exists
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update all fields
	existing.Title = req.Title
	existing.Description = req.Description
	existing.PriceCents = req.PriceCents
	existing.Category = req.Category
	existing.Condition = req.Condition
	existing.Status = req.Status
	existing.IsFeatured = req.IsFeatured
	existing.City = req.City
	existing.PostalCode = req.PostalCode
	existing.Country = req.Country
	existing.DeliveryAvailable = req.DeliveryAvailable
	existing.Images = req.Images

	// Store updated item
	if err := uc.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}
