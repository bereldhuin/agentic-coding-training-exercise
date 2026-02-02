package usecase

import (
	"context"
	"strings"

	"github.com/hymaia/lebonpoint-app/server/go/internal/application/dto"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
	apperrors "github.com/hymaia/lebonpoint-app/server/go/internal/shared/errors"
)

// ListItemsUseCase handles listing items with filters and pagination
type ListItemsUseCase struct {
	repo repository.ItemRepository
}

// NewListItemsUseCase creates a new ListItemsUseCase
func NewListItemsUseCase(repo repository.ItemRepository) *ListItemsUseCase {
	return &ListItemsUseCase{
		repo: repo,
	}
}

// Execute lists items based on the provided query parameters
func (uc *ListItemsUseCase) Execute(ctx context.Context, params *dto.QueryParams) (*repository.ItemPage, error) {
	// Build filters
	filters := uc.buildFilters(params)

	// Parse sort options
	sort, err := uc.parseSortOptions(params)
	if err != nil {
		return nil, err
	}

	// Get pagination params
	limit := params.GetLimit()
	cursor := params.GetCursor()

	// Execute query
	page, err := uc.repo.FindAll(ctx, filters, sort, limit, cursor)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// buildFilters converts query params to ItemFilters
func (uc *ListItemsUseCase) buildFilters(params *dto.QueryParams) *repository.ItemFilters {
	filters := &repository.ItemFilters{}

	if params.Status != nil && *params.Status != "" {
		if status, err := valueobject.ParseItemStatus(*params.Status); err == nil {
			filters.Status = &status
		}
	}

	if params.Category != nil && *params.Category != "" {
		filters.Category = params.Category
	}

	if params.City != nil && *params.City != "" {
		filters.City = params.City
	}

	if params.PostalCode != nil && *params.PostalCode != "" {
		filters.PostalCode = params.PostalCode
	}

	if params.IsFeatured != nil {
		filters.IsFeatured = params.IsFeatured
	}

	if params.DeliveryAvailable != nil {
		filters.DeliveryAvailable = params.DeliveryAvailable
	}

	// Return nil if no filters are set
	if filters.Status == nil && filters.Category == nil &&
		filters.City == nil && filters.PostalCode == nil &&
		filters.IsFeatured == nil && filters.DeliveryAvailable == nil {
		return nil
	}

	return filters
}

// parseSortOptions parses the sort parameter
func (uc *ListItemsUseCase) parseSortOptions(params *dto.QueryParams) (*repository.SortOptions, error) {
	if params.Sort == nil || *params.Sort == "" {
		return repository.DefaultSortOptions(), nil
	}

	rawSort := strings.TrimSpace(*params.Sort)
	if rawSort == "" {
		return repository.DefaultSortOptions(), nil
	}

	parts := strings.Split(rawSort, ":")
	if len(parts) > 2 {
		return nil, apperrors.NewValidationError("Validation failed", map[string]interface{}{
			"sort": "sort must be in the form field:direction",
		})
	}

	fieldKey := parts[0]
	fieldMap := map[string]repository.SortField{
		"id":           repository.SortFieldID,
		"title":        repository.SortFieldTitle,
		"price_cents":  repository.SortFieldPriceCents,
		"priceCents":   repository.SortFieldPriceCents,
		"created_at":   repository.SortFieldCreatedAt,
		"createdAt":    repository.SortFieldCreatedAt,
		"updated_at":   repository.SortFieldUpdatedAt,
		"updatedAt":    repository.SortFieldUpdatedAt,
		"published_at": repository.SortFieldPublishedAt,
		"publishedAt":  repository.SortFieldPublishedAt,
	}

	sortField, ok := fieldMap[fieldKey]
	if !ok {
		return nil, apperrors.NewValidationError("Validation failed", map[string]interface{}{
			"sort": "sort field must be one of: id, title, price_cents, created_at, updated_at, published_at",
		})
	}

	sortDirection := repository.SortDirectionAsc
	if len(parts) == 2 {
		switch parts[1] {
		case "asc":
			sortDirection = repository.SortDirectionAsc
		case "desc":
			sortDirection = repository.SortDirectionDesc
		default:
			return nil, apperrors.NewValidationError("Validation failed", map[string]interface{}{
				"sort": "sort direction must be asc or desc",
			})
		}
	}

	return &repository.SortOptions{
		Field:     sortField,
		Direction: sortDirection,
	}, nil
}
