package usecase

import (
	"context"

	"github.com/hymaia/lebonpoint-app/server/go/internal/application/dto"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/repository"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
)

// SearchItemsUseCase handles full-text search of items
type SearchItemsUseCase struct {
	repo   repository.ItemRepository
	listUC *ListItemsUseCase
}

// NewSearchItemsUseCase creates a new SearchItemsUseCase
func NewSearchItemsUseCase(repo repository.ItemRepository) *SearchItemsUseCase {
	return &SearchItemsUseCase{
		repo:   repo,
		listUC: NewListItemsUseCase(repo),
	}
}

// Execute performs a full-text search
func (uc *SearchItemsUseCase) Execute(ctx context.Context, params *dto.QueryParams) (*repository.ItemPage, error) {
	query := params.GetQuery()

	// If no query, delegate to list use case
	if query == "" {
		return uc.listUC.Execute(ctx, params)
	}

	// Build filters (same as list)
	filters := uc.buildFilters(params)

	// Parse sort options
	sort, err := uc.listUC.parseSortOptions(params)
	if err != nil {
		return nil, err
	}

	// Get pagination params
	limit := params.GetLimit()
	cursor := params.GetCursor()

	// Execute search
	page, err := uc.repo.Search(ctx, query, filters, sort, limit, cursor)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// buildFilters converts query params to ItemFilters
func (uc *SearchItemsUseCase) buildFilters(params *dto.QueryParams) *repository.ItemFilters {
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
