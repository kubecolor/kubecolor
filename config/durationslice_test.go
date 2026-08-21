package config

import (
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestParseDurationSlice(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DurationSlice
		wantErr bool
	}{
		{
			name:  "single value",
			input: "5m",
			want:  DurationSlice{5 * time.Minute},
		},
		{
			name:  "multiple values",
			input: "5m/2h/1d",
			want:  DurationSlice{5 * time.Minute, 2 * time.Hour, 24 * time.Hour},
		},
		{
			name:  "human units days and years",
			input: "7d/1y",
			want:  DurationSlice{7 * 24 * time.Hour, 365 * 24 * time.Hour},
		},
		{
			name:  "go-style multi unit",
			input: "1h30m",
			want:  DurationSlice{90 * time.Minute},
		},
		{
			name:  "sub-second falls back to go parser",
			input: "500ms",
			want:  DurationSlice{500 * time.Millisecond},
		},
		{
			name:  "whitespace around separators is trimmed",
			input: " 5m / 2h ",
			want:  DurationSlice{5 * time.Minute, 2 * time.Hour},
		},
		{
			name:  "empty string is empty slice",
			input: "",
			want:  DurationSlice{},
		},
		{
			name:    "invalid value errors",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "invalid element among valid errors",
			input:   "5m/nope",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseDurationSlice(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (result %v)", got)
				}
				return
			}
			testutil.NoError(t, err)
			testutil.Equal(t, tt.want, got)
		})
	}
}

func TestDurationSliceUnmarshalText(t *testing.T) {
	var s DurationSlice
	testutil.NoError(t, s.UnmarshalText([]byte("5m/2h")))
	testutil.Equal(t, DurationSlice{5 * time.Minute, 2 * time.Hour}, s)
}

func TestDurationSliceString(t *testing.T) {
	s := DurationSlice{5 * time.Minute, 2 * time.Hour}
	// String output must round-trip back through ParseDurationSlice.
	got, err := ParseDurationSlice(s.String())
	testutil.NoError(t, err)
	testutil.Equal(t, s, got)
}
