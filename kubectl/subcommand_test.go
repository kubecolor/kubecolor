package kubectl

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func TestInspectSubcommandInfo(args []string) (*SubcommandInfo, bool) {
func TestInspectSubcommandInfo(t *testing.T) {
	tests := []struct {
		args       string
		expected   *SubcommandInfo
	}{
		{"get pods", &SubcommandInfo{Subcommand: Get}},
		{"get pod", &SubcommandInfo{Subcommand: Get}},
		{"get po", &SubcommandInfo{Subcommand: Get}},

		{"get pod -o wide", &SubcommandInfo{Subcommand: Get, FormatOption: Wide}},
		{"get pod -o=wide", &SubcommandInfo{Subcommand: Get, FormatOption: Wide}},
		{"get pod -owide", &SubcommandInfo{Subcommand: Get, FormatOption: Wide}},

		{"get pod -o json", &SubcommandInfo{Subcommand: Get, FormatOption: Json}},
		{"get pod -o=json", &SubcommandInfo{Subcommand: Get, FormatOption: Json}},
		{"get pod -ojson", &SubcommandInfo{Subcommand: Get, FormatOption: Json}},

		{"get pod -o yaml", &SubcommandInfo{Subcommand: Get, FormatOption: Yaml}},
		{"get pod -o=yaml", &SubcommandInfo{Subcommand: Get, FormatOption: Yaml}},
		{"get pod -oyaml", &SubcommandInfo{Subcommand: Get, FormatOption: Yaml}},

		{"get pod --output json", &SubcommandInfo{Subcommand: Get, FormatOption: Json}},
		{"get pod --output=json", &SubcommandInfo{Subcommand: Get, FormatOption: Json}},
		{"get pod --output yaml", &SubcommandInfo{Subcommand: Get, FormatOption: Yaml}},
		{"get pod --output=yaml", &SubcommandInfo{Subcommand: Get, FormatOption: Yaml}},
		{"get pod --output wide", &SubcommandInfo{Subcommand: Get, FormatOption: Wide}},
		{"get pod --output=wide", &SubcommandInfo{Subcommand: Get, FormatOption: Wide}},

		{"get pod --no-headers", &SubcommandInfo{Subcommand: Get, NoHeader: true}},
		{"get pod -w", &SubcommandInfo{Subcommand: Get, Watch: true}},
		{"get pod --watch", &SubcommandInfo{Subcommand: Get, Watch: true}},
		{"get pod -h", &SubcommandInfo{Subcommand: Get, Help: true}},
		{"get pod --help", &SubcommandInfo{Subcommand: Get, Help: true}},

		{"describe pod pod-aaa", &SubcommandInfo{Subcommand: Describe}},
		{"top pod", &SubcommandInfo{Subcommand: Top}},
		{"top pods", &SubcommandInfo{Subcommand: Top}},

		{"api-versions", &SubcommandInfo{Subcommand: APIVersions}},

		{"explain pod", &SubcommandInfo{Subcommand: Explain}},
		{"explain pod --recursive=true", &SubcommandInfo{Subcommand: Explain, Recursive: true}},
		{"explain pod --recursive", &SubcommandInfo{Subcommand: Explain, Recursive: true}},

		{"version", &SubcommandInfo{Subcommand: Version}},
		{"version --client", &SubcommandInfo{Subcommand: Version, Client: true}},
		{"version -o json", &SubcommandInfo{Subcommand: Version, FormatOption: Json}},
		{"version -o yaml", &SubcommandInfo{Subcommand: Version, FormatOption: Yaml}},

		{"apply", &SubcommandInfo{Subcommand: Apply}},

		{"rsh", &SubcommandInfo{Subcommand: Rsh}},

		{"testplugin", &SubcommandInfo{Subcommand: KubectlPlugin}},
		{"testplugin with args", &SubcommandInfo{Subcommand: KubectlPlugin}},
		{"my-plugin with multiple words", &SubcommandInfo{Subcommand: KubectlPlugin}},
		// Args are not allowed in-between
		{"my-plugin --hello with multiple words", &SubcommandInfo{Subcommand: Unknown, Help: true}},
		// No plugin found, so assume it is help
		{"my-non-existing-plugin", &SubcommandInfo{Subcommand: Unknown, Help: true}},

		{"", &SubcommandInfo{Subcommand: Unknown, Help: true}},
		{"--only-some-flag", &SubcommandInfo{Subcommand: Unknown, Help: true}},
	}

	pluginHandler := TestPluginHandler{LookupMap: map[string]string{
		"testplugin":                    "/bin/testplugin",
		"my_plugin-with-multiple-words": "/bin/my_plugin-with-multiple-words",
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.args, func(t *testing.T) {
			t.Parallel()
			s := InspectSubcommandInfo(strings.Fields(tt.args), pluginHandler)

			if diff := cmp.Diff(s, tt.expected); diff != "" {
				t.Errorf(diff)
			}
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
