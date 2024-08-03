package testcorpus

import (
	"errors"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestExecuteTest_success(t *testing.T) {
	err := ExecuteTest(Test{
		// Using a command that doesn't use colors,
		// as we're not testing the coloring feature here
		Command: "kubectl __complete",
		Input:   "1.2.3",
		Output:  "1.2.3",
	})
	testutil.MustNoError(t, err)
}

func TestExecuteTest_error(t *testing.T) {
	tests := []struct {
		name string
		test Test
	}{
		{
			name: "missing command",
			test: Test{},
		},
		{
			name: "invalid command",
			test: Test{Command: "invalid-command"},
		},
		{
			name: "wrong output",
			test: Test{
				Command: "kubectl get pods",
				Input:   "some sample input",
				Output:  "not the same as the input",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := ExecuteTest(tc.test)

			if gotErr == nil {
				t.Fatalf("expected error to invalid test: %#v", tc.test)
			}
		})
	}
}

func TestPrintCommand_invalidConfig(t *testing.T) {
	got := printCommand([]string{}, "", []EnvVar{
		{Key: "KUBECOLOR_PRESET", Value: "non-existing-preset"},
	})

	if !strings.HasPrefix(got, "config error: ") {
		t.Fatalf("Expected 'config error' output, but got:\n%s", got)
	}
}

func TestIndent(t *testing.T) {
	got := indent("foo\n"+
		"  bar\n"+
		"\t aoe",
		"__")

	want := "__foo\n" +
		"__  bar\n" +
		"__\t aoe"

	testutil.MustEqual(t, want, got)
}

func TestFormatTestError(t *testing.T) {
	got := FormatTestError(Test{
		Name: "example.txt",
		Env: []EnvVar{
			{Key: "FOO", Value: "bar"},
			{Key: "MOO", Value: "doo"},
		},
	}, errors.New("some sample error\nwith multiple\nlines"))

	want := "❌ \033[91;1mexample.txt\033[0m\n" +
		"\033[91;1m│\033[0m \033[90m(env FOO=\"bar\")\033[0m\n" +
		"\033[91;1m│\033[0m \033[90m(env MOO=\"doo\")\033[0m\n" +
		"\033[91;1m│\033[0m \033[31msome sample error\033[0m\n" +
		"\033[91;1m│\033[0m with multiple\n" +
		"\033[91;1m└─\033[0mlines\n" +
		"\n"
	testutil.MustEqual(t, want, got)
}

func TestCreateColoredDiff(t *testing.T) {
	got := createColoredDiff("./example/file.txt",
		"equal\ndifferent \033[33mcolor\033[0m string",
		"equal\ndifferent \033[31mcolor\033[0m string")

	want := []string{"",
		"  \033[90;3mequal\033[0m",
		"\033[31;48;5;52;1m- \033[0m\033[48;5;52mdifferent \x1b[33;48;5;52mcolor\x1b[0;48;5;52m string\033[0m",
		"\033[32;48;5;22;1m+ \033[0m\033[48;5;22mdifferent \x1b[31;48;5;22mcolor\x1b[0;48;5;22m string\033[0m",
		"",
		"\033[90m-----\033[0m",
		"",
		"  \033[90;3m\"equal\"\033[0m",
		"\033[31;48;5;52;1m- \033[0m\033[48;5;52m\"different \x1b[35;48;5;52m\\x1b[33m\x1b[0;48;5;52mcolor\x1b[35;48;5;52m\\x1b[0m\x1b[0;48;5;52m string\"\033[0m",
		"\033[32;48;5;22;1m+ \033[0m\033[48;5;22m\"different \x1b[35;48;5;22m\\x1b[31m\x1b[0;48;5;22mcolor\x1b[35;48;5;22m\\x1b[0m\x1b[0;48;5;22m string\"\033[0m",
		"",
		"\033[32m(+want)\033[0m \033[31m(-got)\033[0m",
		""}
	testutil.MustEqual(t, want, strings.Split(got, "\n"))
}
