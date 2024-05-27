package config

import (
	"fmt"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCode string
	}{
		{
			name:     "fg/named color",
			input:    "yellow",
			wantCode: "33",
		},
		{
			name:     "bg/named color",
			input:    "bg=yellow",
			wantCode: "43",
		},
		{
			name:     "op",
			input:    "underline",
			wantCode: "4",
		},
		{
			name:     "fg/long hex without hash",
			input:    "ffff22",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "fg/long hex with hash",
			input:    "#ffff22",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "fg/short hex without hash",
			input:    "ff2",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "fg/short hex with hash",
			input:    "#ff2",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "bg/long hex without hash",
			input:    "bg=ffff22",
			wantCode: "48;2;255;255;34",
		},
		{
			name:     "bg/long hex with hash",
			input:    "bg=#ffff22",
			wantCode: "48;2;255;255;34",
		},
		{
			name:     "bg/short hex without hash",
			input:    "bg=ff2",
			wantCode: "48;2;255;255;34",
		},
		{
			name:     "bg/short hex with hash",
			input:    "bg=#ff2",
			wantCode: "48;2;255;255;34",
		},
		{
			name:     "fg/rgb without prefix",
			input:    "255, 255, 34",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "fg/rgb with prefix",
			input:    "rgb(255, 255, 34)",
			wantCode: "38;2;255;255;34",
		},
		{
			name:     "bg/rgb without prefix",
			input:    "bg=255, 255, 34",
			wantCode: "48;2;255;255;34",
		},
		{
			name:     "bg/rgb with prefix",
			input:    "bg=rgb(255, 255, 34)",
			wantCode: "48;2;255;255;34",
		},
		{
			name:     "fg/256bit",
			input:    "33",
			wantCode: "38;5;33",
		},
		{
			name:     "bg/256bit",
			input:    "bg=33",
			wantCode: "48;5;33",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseColor(tc.input)
			testutil.MustNoError(t, err)
			testutil.Equal(t, Color{Source: tc.input, Code: tc.wantCode}, got)
		})
	}
}

func TestRender_onColoredText(t *testing.T) {
	highlight := MustParseColor("cyan")

	s := fmt.Sprintf("prefix %s suffix", highlight.Render("highlighted"))
	testutil.Equal(t, "prefix \033[36mhighlighted\033[0m suffix", s, "only highlight")

	surrounding := MustParseColor("yellow")
	s2 := surrounding.Render(s)
	testutil.Equal(t, "\033[33mprefix \033[36mhighlighted\033[0m\033[33m suffix\033[0m", s2, "with surrounding color")
}
