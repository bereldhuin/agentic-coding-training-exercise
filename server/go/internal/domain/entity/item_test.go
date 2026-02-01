package entity

import (
	"testing"
	"time"

	"github.com/hymaia/lebonpoint-app/server/go/internal/domain/valueobject"
)

func TestNewItem(t *testing.T) {
	item := NewItem("Test Item", 1000, valueobject.ItemConditionNew)

	if item.Title != "Test Item" {
		t.Errorf("Expected title 'Test Item', got '%s'", item.Title)
	}
	if item.PriceCents != 1000 {
		t.Errorf("Expected priceCents 1000, got %d", item.PriceCents)
	}
	if item.Condition != valueobject.ItemConditionNew {
		t.Errorf("Expected condition 'new', got '%s'", item.Condition)
	}
	if item.Status != valueobject.ItemStatusDraft {
		t.Errorf("Expected status 'draft', got '%s'", item.Status)
	}
	if item.Country != "FR" {
		t.Errorf("Expected country 'FR', got '%s'", item.Country)
	}
	if item.IsFeatured {
		t.Error("Expected IsFeatured to be false")
	}
	if item.DeliveryAvailable {
		t.Error("Expected DeliveryAvailable to be false")
	}
	if len(item.Images) != 0 {
		t.Errorf("Expected empty images slice, got %d images", len(item.Images))
	}
	if item.PublishedAt != nil {
		t.Error("Expected PublishedAt to be nil")
	}
	if item.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if item.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestItem_IsPublished(t *testing.T) {
	tests := []struct {
		name   string
		status valueobject.ItemStatus
		want   bool
	}{
		{"draft", valueobject.ItemStatusDraft, false},
		{"active", valueobject.ItemStatusActive, true},
		{"reserved", valueobject.ItemStatusReserved, false},
		{"sold", valueobject.ItemStatusSold, false},
		{"archived", valueobject.ItemStatusArchived, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &Item{Status: tt.status}
			if got := item.IsPublished(); got != tt.want {
				t.Errorf("Item.IsPublished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_Publish(t *testing.T) {
	item := NewItem("Test", 100, valueobject.ItemConditionGood)
	item.Publish()

	if item.Status != valueobject.ItemStatusActive {
		t.Errorf("Expected status 'active', got '%s'", item.Status)
	}
	if item.PublishedAt == nil {
		t.Error("Expected PublishedAt to be set")
	}
	if !item.UpdatedAt.After(item.CreatedAt) {
		t.Error("Expected UpdatedAt to be updated after Publish")
	}
}

func TestItem_Archive(t *testing.T) {
	item := NewItem("Test", 100, valueobject.ItemConditionGood)
	updatedAt := item.UpdatedAt
	item.Archive()

	if item.Status != valueobject.ItemStatusArchived {
		t.Errorf("Expected status 'archived', got '%s'", item.Status)
	}
	if !item.UpdatedAt.After(updatedAt) {
		t.Error("Expected UpdatedAt to be updated after Archive")
	}
}

func TestItem_MarkAsSold(t *testing.T) {
	item := NewItem("Test", 100, valueobject.ItemConditionGood)
	updatedAt := item.UpdatedAt
	item.MarkAsSold()

	if item.Status != valueobject.ItemStatusSold {
		t.Errorf("Expected status 'sold', got '%s'", item.Status)
	}
	if !item.UpdatedAt.After(updatedAt) {
		t.Error("Expected UpdatedAt to be updated after MarkAsSold")
	}
}

func TestItem_MarkAsReserved(t *testing.T) {
	item := NewItem("Test", 100, valueobject.ItemConditionGood)
	updatedAt := item.UpdatedAt
	item.MarkAsReserved()

	if item.Status != valueobject.ItemStatusReserved {
		t.Errorf("Expected status 'reserved', got '%s'", item.Status)
	}
	if !item.UpdatedAt.After(updatedAt) {
		t.Error("Expected UpdatedAt to be updated after MarkAsReserved")
	}
}

func TestItem_SetDraft(t *testing.T) {
	item := NewItem("Test", 100, valueobject.ItemConditionGood)
	now := time.Now()
	item.Publish()
	if item.PublishedAt == nil {
		t.Fatal("Expected PublishedAt to be set after Publish")
	}

	item.SetDraft()

	if item.Status != valueobject.ItemStatusDraft {
		t.Errorf("Expected status 'draft', got '%s'", item.Status)
	}
	if item.PublishedAt != nil {
		t.Error("Expected PublishedAt to be nil after SetDraft")
	}
	if !item.UpdatedAt.After(now) {
		t.Error("Expected UpdatedAt to be updated after SetDraft")
	}
}
