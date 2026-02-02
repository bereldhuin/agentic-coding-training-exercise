package cursor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		id   int64
	}{
		{"small id", 1},
		{"medium id", 42},
		{"large id", 999999},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cursor := Encode(tt.id)
			if cursor == "" {
				t.Error("Encode() returned empty string")
			}
			// Verify it's valid base64
			if _, err := Decode(cursor); err != nil {
				t.Errorf("Encode() produced invalid cursor: %v", err)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		cursor  string
		want    int64
		wantErr bool
	}{
		{
			name:    "valid cursor",
			cursor:  Encode(42),
			want:    42,
			wantErr: false,
		},
		{
			name:    "empty cursor",
			cursor:  "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid base64",
			cursor:  "not-base64!!!",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid json",
			cursor:  base64Encode("not json"),
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative id",
			cursor:  encodeJSON(map[string]int64{"id": -1}),
			want:    0,
			wantErr: true,
		},
		{
			name:    "zero id",
			cursor:  encodeJSON(map[string]int64{"id": 0}),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.cursor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	testIDs := []int64{1, 42, 100, 999999, 1234567890}
	for _, id := range testIDs {
		t.Run(fmt.Sprintf("id_%d", id), func(t *testing.T) {
			cursor := Encode(id)
			decoded, err := Decode(cursor)
			if err != nil {
				t.Errorf("Decode() error = %v", err)
				return
			}
			if decoded != id {
				t.Errorf("Decode() = %v, want %v", decoded, id)
			}
		})
	}
}

func TestMustEncode(t *testing.T) {
	tests := []struct {
		name string
		id   int64
	}{
		{"valid id", 42},
		{"large id", 999999},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cursor := MustEncode(tt.id)
			if cursor == "" {
				t.Error("MustEncode() returned empty string")
			}
		})
	}
}

// Helper functions for testing
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func encodeJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(jsonData)
}
