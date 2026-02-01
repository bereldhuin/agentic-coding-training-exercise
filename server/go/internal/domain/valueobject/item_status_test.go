package valueobject

import (
	"encoding/json"
	"testing"
)

func TestItemStatusString(t *testing.T) {
	tests := []struct {
		name string
		is   ItemStatus
		want string
	}{
		{"draft", ItemStatusDraft, "draft"},
		{"active", ItemStatusActive, "active"},
		{"reserved", ItemStatusReserved, "reserved"},
		{"sold", ItemStatusSold, "sold"},
		{"archived", ItemStatusArchived, "archived"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.is.String(); got != tt.want {
				t.Errorf("ItemStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseItemStatus(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    ItemStatus
		wantErr bool
	}{
		{"valid draft", "draft", ItemStatusDraft, false},
		{"valid active", "active", ItemStatusActive, false},
		{"valid reserved", "reserved", ItemStatusReserved, false},
		{"valid sold", "sold", ItemStatusSold, false},
		{"valid archived", "archived", ItemStatusArchived, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseItemStatus(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseItemStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseItemStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemStatusValidate(t *testing.T) {
	tests := []struct {
		name    string
		is      ItemStatus
		wantErr bool
	}{
		{"valid draft", ItemStatusDraft, false},
		{"valid active", ItemStatusActive, false},
		{"valid reserved", ItemStatusReserved, false},
		{"valid sold", ItemStatusSold, false},
		{"valid archived", ItemStatusArchived, false},
		{"invalid", ItemStatus("invalid"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.is.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ItemStatus.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestItemStatusMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		is      ItemStatus
		want    string
		wantErr bool
	}{
		{"draft", ItemStatusDraft, `"draft"`, false},
		{"active", ItemStatusActive, `"active"`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.is)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemStatus.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("ItemStatus.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestItemStatusUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ItemStatus
		wantErr bool
	}{
		{"valid draft", `"draft"`, ItemStatusDraft, false},
		{"valid active", `"active"`, ItemStatusActive, false},
		{"invalid", `"invalid"`, "", true},
		{"invalid type", `123`, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var is ItemStatus
			err := json.Unmarshal([]byte(tt.data), &is)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemStatus.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && is != tt.want {
				t.Errorf("ItemStatus.UnmarshalJSON() = %v, want %v", is, tt.want)
			}
		})
	}
}
