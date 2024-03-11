package stringutil

import "testing"

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
