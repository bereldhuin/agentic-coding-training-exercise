package cursor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrInvalidCursor = errors.New("invalid cursor format")
)

// CursorData represents the data encoded in a cursor
type CursorData struct {
	ID int64 `json:"id"`
}

// Encode creates a cursor from the last item ID
func Encode(id int64) string {
	data := CursorData{ID: id}
	jsonData, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(jsonData)
}

// Decode parses a cursor and returns the ID
func Decode(cursor string) (int64, error) {
	if cursor == "" {
		return 0, nil
	}

	// Decode base64
	jsonData, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
	}

	// Parse JSON
	var data CursorData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
	}

	if data.ID <= 0 {
		return 0, fmt.Errorf("%w: invalid ID", ErrInvalidCursor)
	}

	return data.ID, nil
}

// MustEncode encodes a cursor or panics
func MustEncode(id int64) string {
	cursor := Encode(id)
	if cursor == "" {
		panic("failed to encode cursor")
	}
	return cursor
}
