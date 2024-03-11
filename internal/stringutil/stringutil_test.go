package stringutil

import (
	"testing"
	"time"
)

func TestParseRatio_success(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLeft  string
		wantRight string
	}{
		{
			name:      "zeros",
			input:     "0/0",
			wantLeft:  "0",
			wantRight: "0",
		},
		{
			name:      "ones",
			input:     "1/1",
			wantLeft:  "1",
			wantRight: "1",
		},
		{
			name:      "unequal",
			input:     "5/9",
			wantLeft:  "5",
			wantRight: "9",
		},
		{
			name:      "extra zeros",
			input:     "005/009",
			wantLeft:  "005",
			wantRight: "009",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotLeft, gotRight, ok := ParseRatio(tc.input)
			if !ok {
				t.Fatalf("failed to parse\ninput: %q", tc.input)
			}
			if tc.wantLeft != gotLeft {
				t.Errorf("want left  %q, got %q", tc.wantLeft, gotLeft)
			}
			if tc.wantRight != gotRight {
				t.Errorf("want right %q, got %q", tc.wantRight, gotRight)
			}
		})
	}
}

func TestParseRatio_fail(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty",
			input: "",
		},
		{
			name:  "missing left",
			input: "/1",
		},
		{
			name:  "missing right",
			input: "1/",
		},
		{
			name:  "spacing",
			input: " 1 / 1 ",
		},
		{
			name:  "decimals",
			input: "1.1/2.2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotLeft, gotRight, ok := ParseRatio(tc.input)
			if ok {
				t.Fatalf("should fail\ninput: %q\nunexpected result:\n  left:  %q\n  right: %q", tc.input, gotLeft, gotRight)
			}
		})
	}
}

func TestCutNumber_success(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNum   string
		wantAfter string
	}{
		{
			name:      "only numbers",
			input:     "12345",
			wantNum:   "12345",
			wantAfter: "",
		},
		{
			name:      "decimal",
			input:     "123.45",
			wantNum:   "123",
			wantAfter: ".45",
		},
		{
			name:      "duration",
			input:     "30m",
			wantNum:   "30",
			wantAfter: "m",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			num, after, ok := CutNumber(tc.input)
			if !ok {
				t.Fatalf("failed to parse\ninput: %q", tc.input)
			}
			if tc.wantNum != num || tc.wantAfter != after {
				t.Errorf("wrong value\ninput: %q\nwant:  %q, %q\ngot:   %q, %q",
					tc.input,
					tc.wantNum, tc.wantAfter,
					num, after)
			}
		})
	}
}

func TestCutNumber_fail(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty",
			input: "",
		},
		{
			name:  "spacing",
			input: "  12345  ",
		},
		{
			name:  "char before",
			input: "d20",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			num, after, ok := CutNumber(tc.input)
			if ok {
				t.Fatalf("should fail\ninput: %q\nunexpected result: %q, %q", tc.input, num, after)
			}
		})
	}
}

func TestParseHumanDuration_success(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Duration
	}{
		{
			name:  "zero",
			input: "0s",
			want:  0,
		},
		{
			name:  "minute",
			input: "1m",
			want:  1 * time.Minute,
		},
		{
			name:  "hour and minute",
			input: "2h30m",
			want:  2*time.Hour + 30*time.Minute,
		},
		{
			name:  "day",
			input: "1d",
			want:  24 * time.Hour,
		},
		{
			name:  "year",
			input: "1y",
			want:  365 * 24 * time.Hour,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ParseHumanDuration(tc.input)
			if !ok {
				t.Fatalf("failed to parse\ninput: %q", tc.input)
			}
			if tc.want != got {
				t.Errorf("wrong value\ninput: %q\nwant:  %s\ngot:   %s", tc.input, tc.want, got)
			}
		})
	}
}

func TestParseHumanDuration_fail(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "zero",
			input: "0",
		},
		{
			name:  "padded",
			input: " 0s ",
		},
		{
			name:  "invalid char",
			input: "14M",
		},
		{
			name:  "too many elements",
			input: "6m5m4m3m2m1m0m",
		},
		{
			name:  "pod name",
			input: "postgresql-0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ParseHumanDuration(tc.input)
			if ok {
				t.Fatalf("should fail\ninput: %q\nunexpected result: %s", tc.input, got)
			}
		})
	}
}
