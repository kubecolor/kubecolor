package kubectl

import (
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

// func TestInspectSubcommandInfo(args []string) (*SubcommandInfo, bool) {
func TestInspectSubcommandInfo(t *testing.T) {
	tests := []struct {
		args     string
		expected *SubcommandInfo
	}{
		{"get pods", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pods"}}},
		{"get pod", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}}},
		{"get po", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"po"}}},

		{"get pod -o wide", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputWide}},
		{"get pod -o=wide", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputWide}},
		{"get pod -owide", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputWide}},

		{"get pod -o json", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputJSON}},
		{"get pod -o=json", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputJSON}},
		{"get pod -ojson", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputJSON}},

		{"get pod -o yaml", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputYAML}},
		{"get pod -o=yaml", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputYAML}},
		{"get pod -oyaml", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputYAML}},

		{"get pod --output json", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputJSON}},
		{"get pod --output=json", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputJSON}},
		{"get pod --output yaml", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputYAML}},
		{"get pod --output=yaml", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputYAML}},
		{"get pod --output wide", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputWide}},
		{"get pod --output=wide", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputWide}},

		{"get pod --no-headers", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, NoHeader: true}},
		{"get pod -w", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Watch: true}},
		{"get pod --watch", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Watch: true}},
		{"get pod -h", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Help: true}},
		{"get pod --help", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Help: true}},

		{"get pod --output custom-columns=NAME:.metadata.name", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputCustomColumns}},
		{"get pod --output=custom-columns=NAME:.metadata.name", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputCustomColumns}},
		{"get pod --output custom-columns-file=./foo.txt", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputCustomColumnsFile}},
		{"get pod --output=custom-columns-file=./foo.txt", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputCustomColumnsFile}},

		{"get pod --output name", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputOther}},
		{"get pod --output=name", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputOther}},
		{"get pod --output jsonpath=...", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputOther}},
		{"get pod --output=jsonpath=...", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pod"}, Output: OutputOther}},

		{"describe pod pod-aaa", &SubcommandInfo{Subcommand: Describe, SubcommandArgs: []string{"pod", "pod-aaa"}}},
		{"top pod", &SubcommandInfo{Subcommand: Top, SubcommandArgs: []string{"pod"}}},
		{"top pods", &SubcommandInfo{Subcommand: Top, SubcommandArgs: []string{"pods"}}},

		{"auth can-i", &SubcommandInfo{Subcommand: Auth, SubcommandArgs: []string{"can-i"}}},
		{"auth whoami", &SubcommandInfo{Subcommand: Auth, SubcommandArgs: []string{"whoami"}}},
		{"config get-contexts", &SubcommandInfo{Subcommand: Config, SubcommandArgs: []string{"get-contexts"}}},
		{"config current-context", &SubcommandInfo{Subcommand: Config, SubcommandArgs: []string{"current-context"}}},
		{"config view --minify", &SubcommandInfo{Subcommand: Config, SubcommandArgs: []string{"view"}}},

		{"api-versions", &SubcommandInfo{Subcommand: APIVersions}},

		{"explain pod", &SubcommandInfo{Subcommand: Explain, SubcommandArgs: []string{"pod"}}},
		{"explain pod --recursive=true", &SubcommandInfo{Subcommand: Explain, SubcommandArgs: []string{"pod"}, Recursive: true}},
		{"explain pod --recursive", &SubcommandInfo{Subcommand: Explain, SubcommandArgs: []string{"pod"}, Recursive: true}},

		{"version", &SubcommandInfo{Subcommand: Version}},
		{"version --client", &SubcommandInfo{Subcommand: Version, Client: true}},
		{"version -o json", &SubcommandInfo{Subcommand: Version, Output: OutputJSON}},
		{"version -o yaml", &SubcommandInfo{Subcommand: Version, Output: OutputYAML}},
		{"auth can-i --list", &SubcommandInfo{Subcommand: Auth, SubcommandArgs: []string{"can-i"}, List: true}},

		{"apply", &SubcommandInfo{Subcommand: Apply}},

		{"rsh", &SubcommandInfo{Subcommand: Rsh}},

		{"testplugin", &SubcommandInfo{Subcommand: KubectlPlugin}},
		{"testplugin with args", &SubcommandInfo{Subcommand: KubectlPlugin, SubcommandArgs: []string{"with", "args"}}},
		{"my-plugin with multiple words", &SubcommandInfo{Subcommand: KubectlPlugin, SubcommandArgs: []string{"with", "multiple", "words"}}},
		// Args are not allowed in-between
		{"my-plugin --hello with multiple words", &SubcommandInfo{Subcommand: Unknown, Help: true}},
		// No plugin found, so assume it is help
		{"my-non-existing-plugin", &SubcommandInfo{Subcommand: Unknown, Help: true}},

		{"get pods -- --help", &SubcommandInfo{Subcommand: Get, SubcommandArgs: []string{"pods"}}},
		{"-- --help", &SubcommandInfo{Subcommand: Unknown, Help: true}},

		{"", &SubcommandInfo{Subcommand: Unknown, Help: true}},
		{"--only-some-flag", &SubcommandInfo{Subcommand: Unknown, Help: true}},
	}

	pluginHandler := TestPluginHandler{LookupMap: map[string]string{
		"testplugin":                    "/bin/testplugin",
		"my_plugin-with-multiple-words": "/bin/my_plugin-with-multiple-words",
	}}

	for _, tc := range tests {
		t.Run(tc.args, func(t *testing.T) {
			t.Parallel()
			s := InspectSubcommandInfo(strings.Fields(tc.args), pluginHandler)
			testutil.Equal(t, s, tc.expected)
		})
	}
}

func TestParseArgFlag(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFlag  string
		wantValue string
		wantSkip  int
	}{
		{
			name:      "empty",
			args:      nil,
			wantFlag:  "",
			wantValue: "",
			wantSkip:  0,
		},

		{
			name:      "single long",
			args:      []string{"--output"},
			wantFlag:  "--output",
			wantValue: "",
			wantSkip:  1,
		},
		{
			name:      "long 2 args",
			args:      []string{"--output", "wide"},
			wantFlag:  "--output",
			wantValue: "wide",
			wantSkip:  2,
		},
		{
			name:      "long equals",
			args:      []string{"--output=wide"},
			wantFlag:  "--output",
			wantValue: "wide",
			wantSkip:  1,
		},
		{
			name:      "long double equals",
			args:      []string{"--output=wide"},
			wantFlag:  "--output",
			wantValue: "wide",
			wantSkip:  1,
		},

		{
			name:      "single short",
			args:      []string{"-o"},
			wantFlag:  "-o",
			wantValue: "",
			wantSkip:  1,
		},
		{
			name:      "short 2 args",
			args:      []string{"-o", "wide"},
			wantFlag:  "-o",
			wantValue: "wide",
			wantSkip:  2,
		},
		{
			name:      "short smushed",
			args:      []string{"-owide"},
			wantFlag:  "-o",
			wantValue: "wide",
			wantSkip:  1,
		},
		{
			name:      "short equals",
			args:      []string{"-o=wide"},
			wantFlag:  "-o",
			wantValue: "wide",
			wantSkip:  1,
		},
		{
			name:      "short double equals",
			args:      []string{"-o==wide"},
			wantFlag:  "-o",
			wantValue: "=wide",
			wantSkip:  1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flag, value, skip := parseArgFlag(tc.args)
			testutil.Equalf(t, tc.wantFlag, flag, "flag of %v", tc.args)
			testutil.Equalf(t, tc.wantValue, value, "value of %v", tc.args)
			testutil.Equalf(t, tc.wantSkip, skip, "skip of %v", tc.args)
		})
	}
}

type TestPluginHandler struct {
	LookupMap map[string]string
}

// Ensure it implements the interface
var _ PluginHandler = TestPluginHandler{}

// Lookup implements PluginHandler
func (t TestPluginHandler) Lookup(filename string) (string, bool) {
	path, found := t.LookupMap[filename]
	return path, found
}
