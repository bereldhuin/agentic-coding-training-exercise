package valueobject

import (
	"encoding/json"
)

// ItemImage represents an image associated with an item
type ItemImage struct {
	URL       string `json:"url" binding:"required,url"`
	Alt       string `json:"alt,omitempty"`
	SortOrder int    `json:"sortOrder,omitempty"`
}

// MarshalJSON implements json.Marshaler
func (ii ItemImage) MarshalJSON() ([]byte, error) {
	type alias ItemImage
	return json.Marshal(struct {
		URL       string `json:"url"`
		Alt       string `json:"alt,omitempty"`
		SortOrder int    `json:"sortOrder,omitempty"`
	}{
		URL:       ii.URL,
		Alt:       ii.Alt,
		SortOrder: ii.SortOrder,
	})
}

// UnmarshalJSON implements json.Unmarshaler
func (ii *ItemImage) UnmarshalJSON(data []byte) error {
	type alias ItemImage
	aux := &struct {
		URL       string `json:"url" binding:"required,url"`
		Alt       string `json:"alt,omitempty"`
		SortOrder *int   `json:"sortOrder,omitempty"`
	}{
		URL: ii.URL,
		Alt: ii.Alt,
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ii.URL = aux.URL
	ii.Alt = aux.Alt
	if aux.SortOrder != nil {
		ii.SortOrder = *aux.SortOrder
	}
	return nil
}
