package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_YamlPrinter_Print(t *testing.T) {
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
				apiVersion: v1
				kind: "Pod"
				num: 415
				unknown: <unknown>
				none: <none>
				bool: true`),
			expected: testutil.NewHereDoc(`
				\e[33mapiVersion\e[0m: \e[37mv1\e[0m
				\e[33mkind\e[0m: "\e[37mPod\e[0m"
				\e[33mnum\e[0m: \e[35m415\e[0m
				\e[33munknown\e[0m: \e[33m<unknown>\e[0m
				\e[33mnone\e[0m: \e[33m<none>\e[0m
				\e[33mbool\e[0m: \e[32mtrue\e[0m
			`),
		},
		{
			name:  "key color changes based on its indentation",
			theme: testconfig.DarkTheme,
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
				\e[33mapiVersion\e[0m: \e[37mv1\e[0m
				\e[33mitems\e[0m:
				- \e[37mapiVersion\e[0m: \e[37mv1\e[0m
				  \e[37mkey\e[0m:
				  - \e[33mkey2\e[0m: \e[35m415\e[0m
				    \e[33mkey3\e[0m: \e[32mtrue\e[0m
				  \e[37mkey4\e[0m:
				    \e[33mkey\e[0m: \e[37mval\e[0m
			`),
		},
		{
			name:  "elements in an array can be colored",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				lifecycle:
				  preStop:
				    exec:
				      command:
				      - sh
				      - c
				      - sleep 30`),
			expected: testutil.NewHereDoc(`
				\e[33mlifecycle\e[0m:
				  \e[37mpreStop\e[0m:
				    \e[33mexec\e[0m:
				      \e[37mcommand\e[0m:
				      - \e[37msh\e[0m
				      - \e[37mc\e[0m
				      - \e[37msleep 30\e[0m
			`),
		},
		{
			name:  "a value contains dash",
			theme: testconfig.DarkTheme,
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
				\e[33mapiVersion\e[0m: \e[37mv1\e[0m
				\e[33mitems\e[0m:
				- \e[37mapiVersion\e[0m: \e[37mv1\e[0m
				  \e[37mkey\e[0m:
				  - \e[33mkey2\e[0m: \e[35m415\e[0m
				    \e[33mkey3\e[0m: \e[32mtrue\e[0m
				  \e[37mkey4\e[0m:
				    \e[33mkey\e[0m: \e[37m-val\e[0m
			`),
		},
		{
			name:  "a long string which is broken into several lines can be colored",
			theme: testconfig.DarkTheme,
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
				- \e[37mapiVersion\e[0m: \e[37mv1\e[0m
				  \e[37mkind\e[0m: \e[37mPod\e[0m
				  \e[37mmetadata\e[0m:
				    \e[33mannotations\e[0m:
				      \e[37mannotation.long.1\e[0m: \e[37m'Sometimes, you may want to specify what to command to use as kubectl.\e[0m
				        \e[37mFor example, when you want to use a versioned-kubectl kubectl.1.17, you can do that by an environment variable.'\e[0m
				      \e[37mannotation.long.2\e[0m: \e[37mkubecolor colorizes your kubectl command output and does nothing else.\e[0m
				        \e[37mkubecolor internally calls kubectl command and try to colorizes the output so you can use kubecolor as a\e[0m
				        \e[37mcomplete alternative of kubectl\e[0m
				      \e[37mannotation.short.1\e[0m: \e[37mnormal length annotation\e[0m
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
				Theme: tt.theme,
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
