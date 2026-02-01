package valueobject

import (
	"encoding/json"
	"errors"
)

// ItemCondition represents the condition of an item
type ItemCondition string

const (
	ItemConditionNew     ItemCondition = "new"
	ItemConditionLikeNew ItemCondition = "like_new"
	ItemConditionGood    ItemCondition = "good"
	ItemConditionFair    ItemCondition = "fair"
	ItemConditionParts   ItemCondition = "parts"
	ItemConditionUnknown ItemCondition = "unknown"
)

var (
	ErrInvalidItemCondition = errors.New("invalid item condition")
)

// AllItemConditions returns all valid item conditions
func AllItemConditions() []ItemCondition {
	return []ItemCondition{
		ItemConditionNew,
		ItemConditionLikeNew,
		ItemConditionGood,
		ItemConditionFair,
		ItemConditionParts,
		ItemConditionUnknown,
	}
}

// ParseItemCondition parses a string to ItemCondition
func ParseItemCondition(s string) (ItemCondition, error) {
	ic := ItemCondition(s)
	for _, valid := range AllItemConditions() {
		if ic == valid {
			return ic, nil
		}
	}
	return "", ErrInvalidItemCondition
}

// String returns the string representation of ItemCondition
func (ic ItemCondition) String() string {
	return string(ic)
}

// Validate validates the ItemCondition
func (ic ItemCondition) Validate() error {
	_, err := ParseItemCondition(string(ic))
	return err
}

// MarshalJSON implements json.Marshaler
func (ic ItemCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ic))
}

// UnmarshalJSON implements json.Unmarshaler
func (ic *ItemCondition) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseItemCondition(s)
	if err != nil {
		return err
	}
	*ic = parsed
	return nil
}
