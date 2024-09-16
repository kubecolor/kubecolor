package kubectl

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kubecolor/kubecolor/testutil"
)

// func TestInspectSubcommandInfo(args []string) (*SubcommandInfo, bool) {
func TestInspectSubcommandInfo(t *testing.T) {
	tests := []struct {
		args     string
		expected *SubcommandInfo
	}{
		{"get pods", &SubcommandInfo{Subcommand: Get}},
		{"get pod", &SubcommandInfo{Subcommand: Get}},
		{"get po", &SubcommandInfo{Subcommand: Get}},

		{"get pod -o wide", &SubcommandInfo{Subcommand: Get, Output: OutputWide}},
		{"get pod -o=wide", &SubcommandInfo{Subcommand: Get, Output: OutputWide}},
		{"get pod -owide", &SubcommandInfo{Subcommand: Get, Output: OutputWide}},

		{"get pod -o json", &SubcommandInfo{Subcommand: Get, Output: OutputJSON}},
		{"get pod -o=json", &SubcommandInfo{Subcommand: Get, Output: OutputJSON}},
		{"get pod -ojson", &SubcommandInfo{Subcommand: Get, Output: OutputJSON}},

		{"get pod -o yaml", &SubcommandInfo{Subcommand: Get, Output: OutputYAML}},
		{"get pod -o=yaml", &SubcommandInfo{Subcommand: Get, Output: OutputYAML}},
		{"get pod -oyaml", &SubcommandInfo{Subcommand: Get, Output: OutputYAML}},

		{"get pod --output json", &SubcommandInfo{Subcommand: Get, Output: OutputJSON}},
		{"get pod --output=json", &SubcommandInfo{Subcommand: Get, Output: OutputJSON}},
		{"get pod --output yaml", &SubcommandInfo{Subcommand: Get, Output: OutputYAML}},
		{"get pod --output=yaml", &SubcommandInfo{Subcommand: Get, Output: OutputYAML}},
		{"get pod --output wide", &SubcommandInfo{Subcommand: Get, Output: OutputWide}},
		{"get pod --output=wide", &SubcommandInfo{Subcommand: Get, Output: OutputWide}},

		{"get pod --no-headers", &SubcommandInfo{Subcommand: Get, NoHeader: true}},
		{"get pod -w", &SubcommandInfo{Subcommand: Get, Watch: true}},
		{"get pod --watch", &SubcommandInfo{Subcommand: Get, Watch: true}},
		{"get pod -h", &SubcommandInfo{Subcommand: Get, Help: true}},
		{"get pod --help", &SubcommandInfo{Subcommand: Get, Help: true}},

		{"get pod --output custom-columns=NAME:.metadata.name", &SubcommandInfo{Subcommand: Get, Output: OutputCustomColumns}},
		{"get pod --output=custom-columns=NAME:.metadata.name", &SubcommandInfo{Subcommand: Get, Output: OutputCustomColumns}},
		{"get pod --output custom-columns-file=./foo.txt", &SubcommandInfo{Subcommand: Get, Output: OutputCustomColumnsFile}},
		{"get pod --output=custom-columns-file=./foo.txt", &SubcommandInfo{Subcommand: Get, Output: OutputCustomColumnsFile}},

		{"get pod --output name", &SubcommandInfo{Subcommand: Get, Output: OutputOther}},
		{"get pod --output=name", &SubcommandInfo{Subcommand: Get, Output: OutputOther}},
		{"get pod --output jsonpath=...", &SubcommandInfo{Subcommand: Get, Output: OutputOther}},
		{"get pod --output=jsonpath=...", &SubcommandInfo{Subcommand: Get, Output: OutputOther}},

		{"describe pod pod-aaa", &SubcommandInfo{Subcommand: Describe}},
		{"top pod", &SubcommandInfo{Subcommand: Top}},
		{"top pods", &SubcommandInfo{Subcommand: Top}},

		{"api-versions", &SubcommandInfo{Subcommand: APIVersions}},

		{"explain pod", &SubcommandInfo{Subcommand: Explain}},
		{"explain pod --recursive=true", &SubcommandInfo{Subcommand: Explain, Recursive: true}},
		{"explain pod --recursive", &SubcommandInfo{Subcommand: Explain, Recursive: true}},

		{"version", &SubcommandInfo{Subcommand: Version}},
		{"version --client", &SubcommandInfo{Subcommand: Version, Client: true}},
		{"version -o json", &SubcommandInfo{Subcommand: Version, Output: OutputJSON}},
		{"version -o yaml", &SubcommandInfo{Subcommand: Version, Output: OutputYAML}},

		{"apply", &SubcommandInfo{Subcommand: Apply}},

		{"rsh", &SubcommandInfo{Subcommand: Rsh}},

		{"testplugin", &SubcommandInfo{Subcommand: KubectlPlugin}},
		{"testplugin with args", &SubcommandInfo{Subcommand: KubectlPlugin}},
		{"my-plugin with multiple words", &SubcommandInfo{Subcommand: KubectlPlugin}},
		// Args are not allowed in-between
		{"my-plugin --hello with multiple words", &SubcommandInfo{Subcommand: Unknown, Help: true}},
		// No plugin found, so assume it is help
		{"my-non-existing-plugin", &SubcommandInfo{Subcommand: Unknown, Help: true}},

		{"get pods -- --help", &SubcommandInfo{Subcommand: Get}},
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

			if diff := cmp.Diff(s, tc.expected); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestParseArgFlag(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFlag  string
		wantValue string
	}{
		{
			name:      "empty",
			args:      nil,
			wantFlag:  "",
			wantValue: "",
		},

		{
			name:      "single long",
			args:      []string{"--output"},
			wantFlag:  "--output",
			wantValue: "",
		},
		{
			name:      "long 2 args",
			args:      []string{"--output", "wide"},
			wantFlag:  "--output",
			wantValue: "wide",
		},
		{
			name:      "long equals",
			args:      []string{"--output=wide"},
			wantFlag:  "--output",
			wantValue: "wide",
		},
		{
			name:      "long double equals",
			args:      []string{"--output=wide"},
			wantFlag:  "--output",
			wantValue: "wide",
		},

		{
			name:      "single short",
			args:      []string{"-o"},
			wantFlag:  "-o",
			wantValue: "",
		},
		{
			name:      "short 2 args",
			args:      []string{"-o", "wide"},
			wantFlag:  "-o",
			wantValue: "wide",
		},
		{
			name:      "short smushed",
			args:      []string{"-owide"},
			wantFlag:  "-o",
			wantValue: "wide",
		},
		{
			name:      "short equals",
			args:      []string{"-o=wide"},
			wantFlag:  "-o",
			wantValue: "wide",
		},
		{
			name:      "short double equals",
			args:      []string{"-o==wide"},
			wantFlag:  "-o",
			wantValue: "=wide",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flag, value := parseArgFlag(tc.args)
			testutil.Equalf(t, tc.wantFlag, flag, "flag of %v", tc.args)
			testutil.Equalf(t, tc.wantValue, value, "value of %v", tc.args)
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
