package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_JsonPrinter_Print(t *testing.T) {
	tests := []struct {
		name     string
		theme    *config.Theme
		input    string
		expected string
	}{
		{
			name:  "values can be colored by its type",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				{
				    "apiVersion": "v1",
				    "kind": "Pod",
				    "num": 598,
				    "bool": true,
				    "null": null
				}`),
			expected: testutil.NewHereDoc(`
				{
				    "\e[37mapiVersion\e[0m": "\e[37mv1\e[0m",
				    "\e[37mkind\e[0m": "\e[37mPod\e[0m",
				    "\e[37mnum\e[0m": \e[35m598\e[0m,
				    "\e[37mbool\e[0m": \e[32mtrue\e[0m,
				    "\e[37mnull\e[0m": \e[33mnull\e[0m
				}
			`),
		},
		{
			name:  "keys can be colored by its indentation level",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				{
				    "k1": "v1",
				    "k2": {
				        "k3": "v3",
				        "k4": {
				            "k5": "v5"
				        },
				        "k6": "v6"
				    }
				}`),
			expected: testutil.NewHereDoc(`
				{
				    "\e[37mk1\e[0m": "\e[37mv1\e[0m",
				    "\e[37mk2\e[0m": {
				        "\e[33mk3\e[0m": "\e[37mv3\e[0m",
				        "\e[33mk4\e[0m": {
				            "\e[37mk5\e[0m": "\e[37mv5\e[0m"
				        },
				        "\e[33mk6\e[0m": "\e[37mv6\e[0m"
				    }
				}
			`),
		},
		{
			name:  "{} and [] are not colorized",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				{
				    "apiVersion": "v1",
				    "kind": {
				        "k2": [
				            "a",
				            "b",
				            "c"
				        ],
				        "k3": {
				            "k4": "val"
				        },
				        "k5": {}
				    }
				}`),
			expected: testutil.NewHereDoc(`
				{
				    "\e[37mapiVersion\e[0m": "\e[37mv1\e[0m",
				    "\e[37mkind\e[0m": {
				        "\e[33mk2\e[0m": [
				            "\e[37ma\e[0m",
				            "\e[37mb\e[0m",
				            "\e[37mc\e[0m"
				        ],
				        "\e[33mk3\e[0m": {
				            "\e[37mk4\e[0m": "\e[37mval\e[0m"
				        },
				        "\e[33mk5\e[0m": {}
				    }
				}
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := JsonPrinter{Theme: tt.theme}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
