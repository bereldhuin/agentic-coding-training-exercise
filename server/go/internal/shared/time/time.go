package timeutil

import (
	"encoding/json"
	"fmt"
	"time"
)

// CustomTime wraps time.Time for custom JSON marshaling
type CustomTime struct {
	time.Time
}

// NewCustomTime creates a new CustomTime from time.Time
func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}

// Now returns the current time as CustomTime
func Now() CustomTime {
	return CustomTime{Time: time.Now()}
}

// MarshalJSON implements json.Marshaler
// Uses RFC3339 format for consistency with other implementations
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := ct.Time.Format(time.RFC3339Nano)
	return json.Marshal(formatted)
}

// UnmarshalJSON implements json.Unmarshaler
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ct.Time = time.Time{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		// Try RFC3339 without nanoseconds
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}

	ct.Time = t
	return nil
}

// ToTime returns the underlying time.Time
func (ct CustomTime) ToTime() time.Time {
	return ct.Time
}

// String returns the RFC3339 representation
func (ct CustomTime) String() string {
	return ct.Time.Format(time.RFC3339)
}

// IsZero returns true if the time is zero
func (ct CustomTime) IsZero() bool {
	return ct.Time.IsZero()
}

// ParseCustomTime parses a string in RFC3339 format to CustomTime
func ParseCustomTime(s string) (CustomTime, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return CustomTime{}, err
		}
	}
	return CustomTime{Time: t}, nil
}

// MustParseCustomTime parses a string in RFC3339 format to CustomTime or panics
func MustParseCustomTime(s string) CustomTime {
	ct, err := ParseCustomTime(s)
	if err != nil {
		panic(err)
	}
	return ct
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return d.String()
	}
	if d < time.Minute {
		return d.Truncate(time.Millisecond).String()
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
