package usecase

import (
	"context"

	"github.com/hymaia/lebonpoint-app/server/go/internal/application/dto"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/entity"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
	apperrors "github.com/hymaia/lebonpoint-app/server/go/internal/shared/errors"
)

// CreateItemUseCase handles the creation of new items
type CreateItemUseCase struct {
	repo repository.ItemRepository
}

// NewCreateItemUseCase creates a new CreateItemUseCase
func NewCreateItemUseCase(repo repository.ItemRepository) *CreateItemUseCase {
	return &CreateItemUseCase{
		repo: repo,
	}
}

// Execute creates a new item
func (uc *CreateItemUseCase) Execute(ctx context.Context, req *dto.CreateItemRequest) (*entity.Item, error) {
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

	// Set defaults
	if req.Status == "" {
		req.Status = valueobject.ItemStatusDraft
	}
	if req.Images == nil {
		req.Images = []valueobject.ItemImage{}
	}

	// Create entity
	item := &entity.Item{
		Title:             req.Title,
		Description:       req.Description,
		PriceCents:        req.PriceCents,
		Category:          req.Category,
		Condition:         req.Condition,
		Status:            req.Status,
		IsFeatured:        req.IsFeatured,
		City:              req.City,
		PostalCode:        req.PostalCode,
		Country:           req.Country,
		DeliveryAvailable: req.DeliveryAvailable,
		Images:            req.Images,
	}

	// Store item
	created, err := uc.repo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return created, nil
}
