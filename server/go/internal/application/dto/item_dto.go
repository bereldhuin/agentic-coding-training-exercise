package dto

import (
	"time"

	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
)

// ItemResponse represents an item response
type ItemResponse struct {
	ID                int64                     `json:"id"`
	Title             string                    `json:"title"`
	Description       string                    `json:"description,omitempty"`
	PriceCents        int                       `json:"priceCents"`
	Category          string                    `json:"category,omitempty"`
	Condition         valueobject.ItemCondition `json:"condition"`
	Status            valueobject.ItemStatus    `json:"status"`
	IsFeatured        bool                      `json:"isFeatured"`
	City              string                    `json:"city,omitempty"`
	PostalCode        string                    `json:"postalCode,omitempty"`
	Country           string                    `json:"country"`
	DeliveryAvailable bool                      `json:"deliveryAvailable"`
	CreatedAt         time.Time                 `json:"createdAt"`
	UpdatedAt         time.Time                 `json:"updatedAt"`
	PublishedAt       *time.Time                `json:"publishedAt,omitempty"`
	Images            []valueobject.ItemImage   `json:"images"`
}

// FromEntity creates an ItemResponse from a domain entity
func (ir *ItemResponse) FromEntity(item interface{}) {
	// This would be implemented if we have a common entity interface
	// For now, we'll convert directly in handlers
}

// CreateItemRequest represents a request to create an item
type CreateItemRequest struct {
	Title             string                    `json:"title" binding:"required,min=3,max=200"`
	Description       string                    `json:"description,omitempty"`
	PriceCents        int                       `json:"priceCents" binding:"required,gte=0"`
	Category          string                    `json:"category,omitempty"`
	Condition         valueobject.ItemCondition `json:"condition" binding:"required"`
	Status            valueobject.ItemStatus    `json:"status" binding:"required,oneof=draft active reserved sold archived"`
	IsFeatured        bool                      `json:"isFeatured"`
	City              string                    `json:"city,omitempty"`
	PostalCode        string                    `json:"postalCode,omitempty"`
	Country           string                    `json:"country" binding:"required"`
	DeliveryAvailable bool                      `json:"deliveryAvailable"`
	Images            []valueobject.ItemImage   `json:"images"`
}

// GetDefaultValues returns default values for optional fields
func (req *CreateItemRequest) GetDefaultValues() map[string]interface{} {
	defaults := make(map[string]interface{})

	if req.Status == "" {
		defaults["status"] = valueobject.ItemStatusDraft
	}
	if req.Country == "" {
		defaults["country"] = "FR"
	}
	if req.Images == nil {
		defaults["images"] = []valueobject.ItemImage{}
	}

	return defaults
}

// UpdateItemRequest represents a request to update an item (full replacement)
type UpdateItemRequest struct {
	Title             string                    `json:"title" binding:"required,min=3,max=200"`
	Description       string                    `json:"description,omitempty"`
	PriceCents        int                       `json:"priceCents" binding:"required,gte=0"`
	Category          string                    `json:"category,omitempty"`
	Condition         valueobject.ItemCondition `json:"condition" binding:"required"`
	Status            valueobject.ItemStatus    `json:"status" binding:"required,oneof=draft active reserved sold archived"`
	IsFeatured        bool                      `json:"isFeatured"`
	City              string                    `json:"city,omitempty"`
	PostalCode        string                    `json:"postalCode,omitempty"`
	Country           string                    `json:"country" binding:"required"`
	DeliveryAvailable bool                      `json:"deliveryAvailable"`
	Images            []valueobject.ItemImage   `json:"images"`
}

// PatchItemRequest represents a request to partially update an item
type PatchItemRequest struct {
	Title             *string                    `json:"title" binding:"omitempty,min=3,max=200"`
	Description       *string                    `json:"description,omitempty"`
	PriceCents        *int                       `json:"priceCents" binding:"omitempty,gte=0"`
	Category          *string                    `json:"category,omitempty"`
	Condition         *valueobject.ItemCondition `json:"condition,omitempty"`
	Status            *valueobject.ItemStatus    `json:"status" binding:"omitempty,oneof=draft active reserved sold archived"`
	IsFeatured        *bool                      `json:"isFeatured,omitempty"`
	City              *string                    `json:"city,omitempty"`
	PostalCode        *string                    `json:"postalCode,omitempty"`
	Country           *string                    `json:"country,omitempty"`
	DeliveryAvailable *bool                      `json:"deliveryAvailable,omitempty"`
	Images            *[]valueobject.ItemImage   `json:"images,omitempty"`
}

// HasChanges returns true if at least one field is set
func (req *PatchItemRequest) HasChanges() bool {
	return req.Title != nil ||
		req.Description != nil ||
		req.PriceCents != nil ||
		req.Category != nil ||
		req.Condition != nil ||
		req.Status != nil ||
		req.IsFeatured != nil ||
		req.City != nil ||
		req.PostalCode != nil ||
		req.Country != nil ||
		req.DeliveryAvailable != nil ||
		req.Images != nil
}

// ListItemsResponse represents a paginated list of items
type ListItemsResponse struct {
	Items      []*ItemResponse `json:"items"`
	NextCursor *string         `json:"nextCursor,omitempty"`
	TotalCount int             `json:"totalCount"`
	Limit      int             `json:"limit"`
}

// ErrorResponse represents the canonical error response format
type ErrorResponse struct {
	Error *ErrorDetail `json:"error"`
}

// ErrorDetail represents error details
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string, details map[string]interface{}) *ErrorResponse {
	return &ErrorResponse{
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// ValidationErrorDetail represents a single validation error
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse(message string, details []ValidationErrorDetail) *ErrorResponse {
	detailsMap := make(map[string]interface{})
	for _, d := range details {
		detailsMap[d.Field] = d.Message
	}
	return NewErrorResponse("validation_error", message, detailsMap)
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	Uptime    float64 `json:"uptime,omitempty"`
}

// QueryParams represents query parameters for listing items
type QueryParams struct {
	Status            *string
	Category          *string
	MinPriceCents     *int
	MaxPriceCents     *int
	City              *string
	PostalCode        *string
	IsFeatured        *bool
	DeliveryAvailable *bool
	Sort              *string
	Limit             *int
	Cursor            *string
	Query             *string // For full-text search
}

// GetLimit returns the limit with a default value
func (qp *QueryParams) GetLimit() int {
	if qp.Limit != nil && *qp.Limit > 0 {
		if *qp.Limit > 100 {
			return 100 // Max limit
		}
		return *qp.Limit
	}
	return 20 // Default limit
}

// GetCursor returns the cursor string
func (qp *QueryParams) GetCursor() string {
	if qp.Cursor != nil {
		return *qp.Cursor
	}
	return ""
}

// GetQuery returns the search query
func (qp *QueryParams) GetQuery() string {
	if qp.Query != nil {
		return *qp.Query
	}
	return ""
}
