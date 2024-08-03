package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/gookit/color"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/kubecolor/kubecolor/command"
	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/kubecolor/kubecolor/printer"
)

func RunTests(files []File) {
	var (
		testsPass int
		testsFail int
	)
	for _, file := range files {
		fmt.Printf("  %s:\n", colorHeader.Render(file.Name))
		if len(file.Tests) == 0 {
			fmt.Printf("    %s\n", colorMuted.Render("no tests found"))
		}
		for _, test := range file.Tests {
			if err := ExecuteTest(test); err != nil {
				testsFail++
				fmt.Println(indent(FormatTestError(test, err), "    "))
			} else {
				fmt.Printf("    ✅ %s\n", colorSuccess.Render(test.Name))
				testsPass++
			}
		}
	}
	fmt.Println()
	fmt.Printf("  %s\n", colorMuted.Render("---"))
	fmt.Println()
	fmt.Printf("  %s\n", colorHeader.Render("Results:"))
	if testsPass > 0 {
		fmt.Printf("    Passed: %s\n", colorSuccess.Render(strconv.Itoa(testsPass)))
	} else {
		fmt.Printf("    Passed: %s\n", colorMuted.Render(strconv.Itoa(testsPass)))
	}
	if testsFail > 0 {
		fmt.Printf("    Failed: %s\n", colorErrorText.Render(strconv.Itoa(testsFail)))
	} else {
		fmt.Printf("    Failed: %s\n", colorMuted.Render(strconv.Itoa(testsFail)))
	}
	fmt.Println()

	if testsFail > 0 {
		os.Exit(1)
	}
}

func ExecuteTest(test Test) error {
	args := strings.Fields(test.Command)
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}
	cmd := args[0]
	args = args[1:]

	if cmd != "kubectl" {
		return fmt.Errorf(`command must start with "kubectl", but got %q`, cmd)
	}

	gotOutput := printCommand(args, test.Input, test.Env)
	gotOutput = strings.TrimSpace(gotOutput)

	if test.Output != gotOutput {
		return fmt.Errorf("input does not match output:\n%s", createColoredDiff(test.Name, test.Output, gotOutput))
	}
	return nil
}

func FormatTestError(test Test, err error) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "❌ %s\n", colorErrorPrefix.Render(test.Name))
	for _, env := range test.Env {
		fmt.Fprintf(&buf, "%s %s\n", colorErrorPrefix.Render("│"), colorMuted.Render(fmt.Sprintf("(env %s=%q)", env.Key, env.Value)))
	}
	lines := strings.Split(err.Error(), "\n")
	for i, line := range lines {
		switch {
		case i == 0:
			fmt.Fprintf(&buf, "%s %s\n", colorErrorPrefix.Render("│"), colorErrorText.Render(line))
		case i == len(lines)-1:
			fmt.Fprintf(&buf, "%s%s\n", colorErrorPrefix.Render("└─"), line)
		default:
			fmt.Fprintf(&buf, "%s %s\n", colorErrorPrefix.Render("│"), line)
		}
	}
	fmt.Fprintln(&buf)
	return buf.String()
}

func indent(s, indent string) string {
	return indent + strings.ReplaceAll(s, "\n", "\n"+indent)
}

func printCommand(args []string, input string, env []EnvVar) string {
	os.Clearenv()
	for _, e := range env {
		if err := os.Setenv(e.Key, e.Value); err != nil {
			return fmt.Sprintf("config error: set env %s=%q: %s", e.Key, e.Value, err)
		}
	}

	v := config.NewViper()
	cfg, err := command.ResolveConfigViper(args, v)
	if err != nil {
		return fmt.Sprintf("config error: %s", err)
	}
	cfg.ForceColor = command.ColorLevelTrueColor

	subcommandInfo := kubectl.InspectSubcommandInfo(args, kubectl.NoopPluginHandler{})
	p := &printer.KubectlOutputColoredPrinter{
		SubcommandInfo:    subcommandInfo,
		Recursive:         subcommandInfo.Recursive,
		ObjFreshThreshold: cfg.ObjFreshThreshold,
		Theme:             &cfg.Theme,
		KubecolorVersion:  "dev",
	}
	var buf bytes.Buffer
	p.Print(strings.NewReader(input), &buf)
	return buf.String()
}

func createColoredDiff(path, want, got string) string {
	edits := myers.ComputeEdits(span.URIFromPath(path), want, got)
	unified := gotextdiff.ToUnified("want", "got", want, edits)
	lines := splitHunksIntoLines(unified.Hunks)

	var buf bytes.Buffer
	buf.WriteByte('\n')

	for _, line := range lines {
		switch line.Kind {
		case gotextdiff.Equal:
			fmt.Fprintf(&buf, "  %s\n", colorDiffEqual.Render(color.ClearCode(line.Content)))
		case gotextdiff.Insert:
			fmt.Fprintf(&buf, "%s%s\n", colorDiffAddPrefix.Render("+ "), injectColor(line.Content, colorDiffAdd))
		case gotextdiff.Delete:
			fmt.Fprintf(&buf, "%s%s\n", colorDiffDelPrefix.Render("- "), injectColor(line.Content, colorDiffDel))
		}
	}

	fmt.Fprintf(&buf, "\n%s\n\n", colorMuted.Render("-----"))

	tabbedLines := quoteAndTabWrite(lines)
	for _, diff := range tabbedLines {
		switch diff.Kind {
		case gotextdiff.Equal:
			fmt.Fprintf(&buf, "  %s\n", colorDiffEqual.Render(diff.Content))
		case gotextdiff.Insert:
			text := injectColor(highlightEscapedColorCodes(diff.Content, colorDiffColorHighlight), colorDiffAdd)
			fmt.Fprintf(&buf, "%s%s\n", colorDiffAddPrefix.Render("+ "), text)
		case gotextdiff.Delete:
			text := injectColor(highlightEscapedColorCodes(diff.Content, colorDiffColorHighlight), colorDiffDel)
			fmt.Fprintf(&buf, "%s%s\n", colorDiffDelPrefix.Render("- "), text)
		}
	}

	fmt.Fprintf(&buf, "\n%s %s\n", colorSuccess.Render("(+want)"), colorErrorText.Render("(-got)"))

	return buf.String()
}

func splitHunksIntoLines(hunks []*gotextdiff.Hunk) []gotextdiff.Line {
	var result []gotextdiff.Line
	for _, hunk := range hunks {
		for _, line := range hunk.Lines {
			result = append(result, gotextdiff.Line{Kind: line.Kind, Content: strings.Trim(line.Content, "\n\r")})
		}
	}
	return result
}

var tripleSpaceRegex = regexp.MustCompile(` {3,}`)

func quoteAndTabWrite(diffs []gotextdiff.Line) []gotextdiff.Line {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	for _, diff := range diffs {
		quoted := fmt.Sprintf("%q", diff.Content)
		tabbed := tripleSpaceRegex.ReplaceAllLiteralString(quoted, "\t")
		fmt.Fprintln(tw, tabbed)
	}
	tw.Flush()
	lines := strings.Split(buf.String(), "\n")
	newDiffs := make([]gotextdiff.Line, len(diffs))
	for i, diff := range diffs {
		newDiffs[i] = gotextdiff.Line{Kind: diff.Kind, Content: lines[i]}
	}
	return newDiffs
}

var colorRegex = regexp.MustCompile(`\x1b\[[0-9;\.,]+m`)

func injectColor(s string, color config.Color) string {
	newCode := strings.TrimSuffix(strings.TrimPrefix(color.ANSICode(), "\x1b["), "m")

	updatedColors := colorRegex.ReplaceAllStringFunc(s, func(s string) string {
		originalCode := strings.TrimSuffix(strings.TrimPrefix(s, "\x1b["), "m")

		return fmt.Sprintf("\x1b[%s;%sm", originalCode, newCode)
	})
	return color.Render(updatedColors)
}

var escapedColorRegex = regexp.MustCompile(`\\x1b\[[0-9;\.,]+m`)

func highlightEscapedColorCodes(s string, color config.Color) string {
	return escapedColorRegex.ReplaceAllStringFunc(s, func(s string) string {
		return color.Render(s)
	})
}
