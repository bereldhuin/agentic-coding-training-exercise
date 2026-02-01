package timeutil

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCustomTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		ct      CustomTime
		want    string
		wantErr bool
	}{
		{
			name:    "valid time",
			ct:      CustomTime{Time: mustParseTime("2026-01-31T12:00:00Z")},
			want:    `"2026-01-31T12:00:00Z"`,
			wantErr: false,
		},
		{
			name:    "zero time",
			ct:      CustomTime{Time: time.Time{}},
			want:    "null",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.ct)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomTime.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("CustomTime.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestCustomTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid RFC3339",
			data:    `"2026-01-31T12:00:00Z"`,
			want:    mustParseTime("2026-01-31T12:00:00Z"),
			wantErr: false,
		},
		{
			name:    "valid RFC3339 with nanoseconds",
			data:    `"2026-01-31T12:00:00.123456789Z"`,
			want:    mustParseTime("2026-01-31T12:00:00.123456789Z"),
			wantErr: false,
		},
		{
			name:    "null",
			data:    "null",
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "invalid format",
			data:    `"not-a-time"`,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "invalid type",
			data:    `123`,
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct CustomTime
			err := json.Unmarshal([]byte(tt.data), &ct)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !ct.Time.Equal(tt.want) {
				t.Errorf("CustomTime.UnmarshalJSON() = %v, want %v", ct.Time, tt.want)
			}
		})
	}
}

func TestCustomTime_RoundTrip(t *testing.T) {
	originalTime := time.Now().UTC().Truncate(time.Microsecond)
	original := CustomTime{Time: originalTime}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}

	// Unmarshal
	var unmarshaled CustomTime
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	// Compare
	if !unmarshaled.Time.Equal(original.Time) {
		t.Errorf("RoundTrip = %v, want %v", unmarshaled.Time, original.Time)
	}
}

func TestNow(t *testing.T) {
	ct := Now()
	if ct.IsZero() {
		t.Error("Now() returned zero time")
	}
	// Check it's within last second
	if time.Since(ct.Time) > time.Second {
		t.Error("Now() returned time too far in the past")
	}
}

func TestParseCustomTime(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantErr bool
	}{
		{
			name:    "valid RFC3339",
			s:       "2026-01-31T12:00:00Z",
			wantErr: false,
		},
		{
			name:    "valid RFC3339 with nanoseconds",
			s:       "2026-01-31T12:00:00.123456789Z",
			wantErr: false,
		},
		{
			name:    "invalid format",
			s:       "not-a-time",
			wantErr: true,
		},
		{
			name:    "empty string",
			s:       "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCustomTime(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCustomTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMustParseCustomTime(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool // should panic
	}{
		{
			name: "valid",
			s:    "2026-01-31T12:00:00Z",
			want: false,
		},
		{
			name: "invalid",
			s:    "not-a-time",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.want {
					t.Errorf("MustParseCustomTime() panic = %v, want panic %v", r, tt.want)
				}
			}()
			MustParseCustomTime(tt.s)
		})
	}
}

func TestCustomTime_String(t *testing.T) {
	ct := CustomTime{Time: mustParseTime("2026-01-31T12:00:00Z")}
	want := "2026-01-31T12:00:00Z"
	if got := ct.String(); got != want {
		t.Errorf("CustomTime.String() = %v, want %v", got, want)
	}
}

func TestCustomTime_IsZero(t *testing.T) {
	tests := []struct {
		name string
		ct   CustomTime
		want bool
	}{
		{
			name: "zero time",
			ct:   CustomTime{Time: time.Time{}},
			want: true,
		},
		{
			name: "non-zero time",
			ct:   CustomTime{Time: mustParseTime("2026-01-31T12:00:00Z")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.IsZero(); got != tt.want {
				t.Errorf("CustomTime.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function
func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			panic(err)
		}
	}
	return t
}
