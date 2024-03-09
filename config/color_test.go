package config

import "testing"

func TestParseColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCode string
	}{
		{
			name:     "named color",
			input:    "yellow",
			wantCode: "33",
		},
		{
			name:     "op",
			input:    "underline",
			wantCode: "4",
		},
		{
			name:     "hex without hash",
			input:    "ff2",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "hex with hash",
			input:    "#ff2",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "rgb without prefix",
			input:    "255, 255, 34",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "rgb with prefix",
			input:    "rgb(255, 255, 34)",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "hsl",
			input:    "hsl(0.1, 0.8, 0.4)",
			wantCode: "38;2;184;118;20",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseColor(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got.Source != tc.input {
				t.Errorf("wrong source\nwant %q\ngot  %q", tc.input, got.Source)
			}
			if got.Code != tc.wantCode {
				t.Errorf("wrong code\nwant %q\ngot  %q", tc.wantCode, got.Code)
			}
		})
	}
}
