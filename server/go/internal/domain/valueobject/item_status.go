package valueobject

import (
	"encoding/json"
	"errors"
)

// ItemStatus represents the status of an item
type ItemStatus string

const (
	ItemStatusDraft    ItemStatus = "draft"
	ItemStatusActive   ItemStatus = "active"
	ItemStatusReserved ItemStatus = "reserved"
	ItemStatusSold     ItemStatus = "sold"
	ItemStatusArchived ItemStatus = "archived"
)

var (
	ErrInvalidItemStatus = errors.New("invalid item status")
)

// AllItemStatuses returns all valid item statuses
func AllItemStatuses() []ItemStatus {
	return []ItemStatus{
		ItemStatusDraft,
		ItemStatusActive,
		ItemStatusReserved,
		ItemStatusSold,
		ItemStatusArchived,
	}
}

// ParseItemStatus parses a string to ItemStatus
func ParseItemStatus(s string) (ItemStatus, error) {
	is := ItemStatus(s)
	for _, valid := range AllItemStatuses() {
		if is == valid {
			return is, nil
		}
	}
	return "", ErrInvalidItemStatus
}

// String returns the string representation of ItemStatus
func (is ItemStatus) String() string {
	return string(is)
}

// Validate validates the ItemStatus
func (is ItemStatus) Validate() error {
	_, err := ParseItemStatus(string(is))
	return err
}

// MarshalJSON implements json.Marshaler
func (is ItemStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(is))
}

// UnmarshalJSON implements json.Unmarshaler
func (is *ItemStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseItemStatus(s)
	if err != nil {
		return err
	}
	*is = parsed
	return nil
}
