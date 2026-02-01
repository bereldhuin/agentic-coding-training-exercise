package valueobject

import (
	"encoding/json"
	"testing"
)

func TestItemConditionString(t *testing.T) {
	tests := []struct {
		name string
		ic   ItemCondition
		want string
	}{
		{"new", ItemConditionNew, "new"},
		{"like_new", ItemConditionLikeNew, "like_new"},
		{"good", ItemConditionGood, "good"},
		{"fair", ItemConditionFair, "fair"},
		{"parts", ItemConditionParts, "parts"},
		{"unknown", ItemConditionUnknown, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ic.String(); got != tt.want {
				t.Errorf("ItemCondition.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseItemCondition(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    ItemCondition
		wantErr bool
	}{
		{"valid new", "new", ItemConditionNew, false},
		{"valid like_new", "like_new", ItemConditionLikeNew, false},
		{"valid good", "good", ItemConditionGood, false},
		{"valid fair", "fair", ItemConditionFair, false},
		{"valid parts", "parts", ItemConditionParts, false},
		{"valid unknown", "unknown", ItemConditionUnknown, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseItemCondition(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseItemCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseItemCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemConditionValidate(t *testing.T) {
	tests := []struct {
		name    string
		ic      ItemCondition
		wantErr bool
	}{
		{"valid new", ItemConditionNew, false},
		{"valid like_new", ItemConditionLikeNew, false},
		{"valid good", ItemConditionGood, false},
		{"valid fair", ItemConditionFair, false},
		{"valid parts", ItemConditionParts, false},
		{"valid unknown", ItemConditionUnknown, false},
		{"invalid", ItemCondition("invalid"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ic.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ItemCondition.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestItemConditionMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		ic      ItemCondition
		want    string
		wantErr bool
	}{
		{"new", ItemConditionNew, `"new"`, false},
		{"like_new", ItemConditionLikeNew, `"like_new"`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.ic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemCondition.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("ItemCondition.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestItemConditionUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ItemCondition
		wantErr bool
	}{
		{"valid new", `"new"`, ItemConditionNew, false},
		{"valid like_new", `"like_new"`, ItemConditionLikeNew, false},
		{"invalid", `"invalid"`, "", true},
		{"invalid type", `123`, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ic ItemCondition
			err := json.Unmarshal([]byte(tt.data), &ic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemCondition.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ic != tt.want {
				t.Errorf("ItemCondition.UnmarshalJSON() = %v, want %v", ic, tt.want)
			}
		})
	}
}
