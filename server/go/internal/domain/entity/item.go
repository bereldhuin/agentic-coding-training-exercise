package entity

import (
	"time"

	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
)

// Item represents a marketplace item
type Item struct {
	ID                int64                     `json:"id"`
	Title             string                    `json:"title" binding:"required,min=3,max=200"`
	Description       string                    `json:"description,omitempty"`
	PriceCents        int                       `json:"priceCents" binding:"required,gte=0"`
	Category          string                    `json:"category,omitempty"`
	Condition         valueobject.ItemCondition `json:"condition" binding:"required"`
	Status            valueobject.ItemStatus    `json:"status" binding:"required"`
	IsFeatured        bool                      `json:"isFeatured"`
	City              string                    `json:"city,omitempty"`
	PostalCode        string                    `json:"postalCode,omitempty"`
	Country           string                    `json:"country" binding:"required"`
	DeliveryAvailable bool                      `json:"deliveryAvailable"`
	CreatedAt         time.Time                 `json:"createdAt"`
	UpdatedAt         time.Time                 `json:"updatedAt"`
	PublishedAt       *time.Time                `json:"publishedAt,omitempty"`
	Images            []valueobject.ItemImage   `json:"images" binding:"required"`
}

// NewItem creates a new Item with default values
func NewItem(title string, priceCents int, condition valueobject.ItemCondition) *Item {
	now := time.Now()
	return &Item{
		Title:             title,
		PriceCents:        priceCents,
		Condition:         condition,
		Status:            valueobject.ItemStatusDraft,
		IsFeatured:        false,
		Country:           "FR",
		DeliveryAvailable: false,
		CreatedAt:         now,
		UpdatedAt:         now,
		Images:            []valueobject.ItemImage{},
	}
}

// IsPublished returns true if the item is published (status is active)
func (i *Item) IsPublished() bool {
	return i.Status == valueobject.ItemStatusActive
}

// Publish publishes the item
func (i *Item) Publish() {
	i.Status = valueobject.ItemStatusActive
	now := time.Now()
	i.PublishedAt = &now
	i.UpdatedAt = now
}

// Archive archives the item
func (i *Item) Archive() {
	i.Status = valueobject.ItemStatusArchived
	i.UpdatedAt = time.Now()
}

// MarkAsSold marks the item as sold
func (i *Item) MarkAsSold() {
	i.Status = valueobject.ItemStatusSold
	i.UpdatedAt = time.Now()
}

// MarkAsReserved marks the item as reserved
func (i *Item) MarkAsReserved() {
	i.Status = valueobject.ItemStatusReserved
	i.UpdatedAt = time.Now()
}

// SetDraft sets the item status back to draft
func (i *Item) SetDraft() {
	i.Status = valueobject.ItemStatusDraft
	i.PublishedAt = nil
	i.UpdatedAt = time.Now()
}
