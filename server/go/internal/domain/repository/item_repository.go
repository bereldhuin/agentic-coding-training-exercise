package repository

import (
	"context"

	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/entity"
	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
)

// ItemRepository defines the interface for item persistence operations
type ItemRepository interface {
	// Create creates a new item and returns it with the generated ID
	Create(ctx context.Context, item *entity.Item) (*entity.Item, error)

	// FindByID retrieves an item by its ID
	// Returns NotFoundError if the item doesn't exist
	FindByID(ctx context.Context, id int64) (*entity.Item, error)

	// FindAll retrieves items with optional filtering, sorting, and pagination
	FindAll(ctx context.Context, filters *ItemFilters, sort *SortOptions, limit int, cursor string) (*ItemPage, error)

	// Update updates an existing item
	// Returns NotFoundError if the item doesn't exist
	Update(ctx context.Context, item *entity.Item) error

	// Delete deletes an item by its ID
	// Returns NotFoundError if the item doesn't exist
	Delete(ctx context.Context, id int64) error

	// Search performs full-text search on items with optional filters
	Search(ctx context.Context, query string, filters *ItemFilters, sort *SortOptions, limit int, cursor string) (*ItemPage, error)
}

// ItemFilters represents filter options for listing items
type ItemFilters struct {
	Status            *valueobject.ItemStatus
	Category          *string
	City              *string
	PostalCode        *string
	IsFeatured        *bool
	DeliveryAvailable *bool
}

// SortField represents a sortable field
type SortField string

const (
	SortFieldID          SortField = "id"
	SortFieldTitle       SortField = "title"
	SortFieldPriceCents  SortField = "priceCents"
	SortFieldCreatedAt   SortField = "createdAt"
	SortFieldUpdatedAt   SortField = "updatedAt"
	SortFieldPublishedAt SortField = "publishedAt"
)

// SortDirection represents the sort direction
type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

// SortOptions represents sorting options
type SortOptions struct {
	Field     SortField
	Direction SortDirection
}

// ItemPage represents a paginated page of items
type ItemPage struct {
	Items      []*entity.Item
	NextCursor *string
	TotalCount int // Total count of items matching the query (for debugging/info)
}

// NewItemPage creates a new ItemPage
func NewItemPage(items []*entity.Item, nextCursor *string, totalCount int) *ItemPage {
	return &ItemPage{
		Items:      items,
		NextCursor: nextCursor,
		TotalCount: totalCount,
	}
}

// HasMore returns true if there are more items
func (p *ItemPage) HasMore() bool {
	return p.NextCursor != nil
}

// DefaultSortOptions returns the default sort options (created_at DESC)
func DefaultSortOptions() *SortOptions {
	return &SortOptions{
		Field:     SortFieldCreatedAt,
		Direction: SortDirectionDesc,
	}
}
