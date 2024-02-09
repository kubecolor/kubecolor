package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_YamlPrinter_Print(t *testing.T) {
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
				apiVersion: v1
				kind: "Pod"
				num: 415
				unknown: <unknown>
				none: <none>
				bool: true`),
			expected: testutil.NewHereDoc(`
				[33mapiVersion[0m: [37mv1[0m
				[33mkind[0m: "[37mPod[0m"
				[33mnum[0m: [35m415[0m
				[33munknown[0m: [33m<unknown>[0m
				[33mnone[0m: [33m<none>[0m
				[33mbool[0m: [32mtrue[0m
			`),
		},
		{
			name:           "key color changes based on its indentation",
			darkBackground: true,
			input: testutil.NewHereDoc(`
				apiVersion: v1
				items:
				- apiVersion: v1
				  key:
				  - key2: 415
				    key3: true
				  key4:
				    key: val`),
			expected: testutil.NewHereDoc(`
				[33mapiVersion[0m: [37mv1[0m
				[33mitems[0m:
				- [37mapiVersion[0m: [37mv1[0m
				  [37mkey[0m:
				  - [33mkey2[0m: [35m415[0m
				    [33mkey3[0m: [32mtrue[0m
				  [37mkey4[0m:
				    [33mkey[0m: [37mval[0m
			`),
		},
		{
			name:           "elements in an array can be colored",
			darkBackground: true,
			input: testutil.NewHereDoc(`
				lifecycle:
				  preStop:
				    exec:
				      command:
				      - sh
				      - c
				      - sleep 30`),
			expected: testutil.NewHereDoc(`
				[33mlifecycle[0m:
				  [37mpreStop[0m:
				    [33mexec[0m:
				      [37mcommand[0m:
				      - [37msh[0m
				      - [37mc[0m
				      - [37msleep 30[0m
			`),
		},
		{
			name:           "a value contains dash",
			darkBackground: true,
			input: testutil.NewHereDoc(`
				apiVersion: v1
				items:
				- apiVersion: v1
				  key:
				  - key2: 415
				    key3: true
				  key4:
				    key: -val`),
			expected: testutil.NewHereDoc(`
				[33mapiVersion[0m: [37mv1[0m
				[33mitems[0m:
				- [37mapiVersion[0m: [37mv1[0m
				  [37mkey[0m:
				  - [33mkey2[0m: [35m415[0m
				    [33mkey3[0m: [32mtrue[0m
				  [37mkey4[0m:
				    [33mkey[0m: [37m-val[0m
			`),
		},
		{
			name:           "a long string which is broken into several lines can be colored",
			darkBackground: true,
			input: testutil.NewHereDoc(`
				- apiVersion: v1
				  kind: Pod
				  metadata:
				    annotations:
				      annotation.long.1: 'Sometimes, you may want to specify what to command to use as kubectl.
				        For example, when you want to use a versioned-kubectl kubectl.1.17, you can do that by an environment variable.'
				      annotation.long.2: kubecolor colorizes your kubectl command output and does nothing else.
				        kubecolor internally calls kubectl command and try to colorizes the output so you can use kubecolor as a
				        complete alternative of kubectl
				      annotation.short.1: normal length annotation`),
			expected: testutil.NewHereDoc(`
				- [37mapiVersion[0m: [37mv1[0m
				  [37mkind[0m: [37mPod[0m
				  [37mmetadata[0m:
				    [33mannotations[0m:
				      [37mannotation.long.1[0m: [37m'Sometimes, you may want to specify what to command to use as kubectl.[0m
				        [37mFor example, when you want to use a versioned-kubectl kubectl.1.17, you can do that by an environment variable.'[0m
				      [37mannotation.long.2[0m: [37mkubecolor colorizes your kubectl command output and does nothing else.[0m
				        [37mkubecolor internally calls kubectl command and try to colorizes the output so you can use kubecolor as a[0m
				        [37mcomplete alternative of kubectl[0m
				      [37mannotation.short.1[0m: [37mnormal length annotation[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := YamlPrinter{
				DarkBackground: tt.darkBackground,
				Theme:          color.NewTheme(tt.darkBackground),
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
