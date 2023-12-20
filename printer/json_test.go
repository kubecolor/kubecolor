package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func Test_JsonPrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		darkBackground bool
		input          string
		expected       string
	}{
		{
			name:           "values can be colored by its type",
			darkBackground: true,
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
				    "[37mapiVersion[0m": "[37mv1[0m",
				    "[37mkind[0m": "[37mPod[0m",
				    "[37mnum[0m": [35m598[0m,
				    "[37mbool[0m": [32mtrue[0m,
				    "[37mnull[0m": [33mnull[0m
				}
			`),
		},
		{
			name:           "keys can be colored by its indentation level",
			darkBackground: true,
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
				    "[37mk1[0m": "[37mv1[0m",
				    "[37mk2[0m": {
				        "[33mk3[0m": "[37mv3[0m",
				        "[33mk4[0m": {
				            "[37mk5[0m": "[37mv5[0m"
				        },
				        "[33mk6[0m": "[37mv6[0m"
				    }
				}
			`),
		},
		{
			name:           "{} and [] are not colorized",
			darkBackground: true,
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
				    "[37mapiVersion[0m": "[37mv1[0m",
				    "[37mkind[0m": {
				        "[33mk2[0m": [
				            "[37ma[0m",
				            "[37mb[0m",
				            "[37mc[0m"
				        ],
				        "[33mk3[0m": {
				            "[37mk4[0m": "[37mval[0m"
				        },
				        "[33mk5[0m": {}
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
			printer := JsonPrinter{DarkBackground: tt.darkBackground, ColorSchema: NewColorSchema(tt.darkBackground)}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
