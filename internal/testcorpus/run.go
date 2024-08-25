package testcorpus

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/kubecolor/kubecolor/command"
	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/kubecolor/kubecolor/printer"
)

var (
	ColorErrorPrefix = color.MustParse("hi-red:bold")
	ColorErrorText   = color.MustParse("red")
	ColorWarnPrefix  = color.MustParse("hi-yellow:bold")
	ColorWarnText    = color.MustParse("yellow")

	ColorHeader  = color.MustParse("bold")
	ColorMuted   = color.MustParse("gray")
	ColorSuccess = color.MustParse("green")

	ColorDiffAddPrefix      = color.MustParse("fg=green:bg=22:bold")
	ColorDiffAdd            = color.MustParse("bg=22") // dark green
	ColorDiffDelPrefix      = color.MustParse("fg=red:bg=52:bold")
	ColorDiffDel            = color.MustParse("bg=52") // dark red
	ColorDiffEqual          = color.MustParse("gray:italic")
	ColorDiffColorHighlight = color.MustParse(`magenta`)
)

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
	fmt.Fprintf(&buf, "❌ %s\n", ColorErrorPrefix.Render(test.Name))
	for _, env := range test.Env {
		fmt.Fprintf(&buf, "%s %s\n", ColorErrorPrefix.Render("│"), ColorMuted.Render(fmt.Sprintf("(env %s=%q)", env.Key, env.Value)))
	}
	lines := strings.Split(err.Error(), "\n")
	for i, line := range lines {
		switch {
		case i == 0:
			fmt.Fprintf(&buf, "%s %s\n", ColorErrorPrefix.Render("│"), ColorErrorText.Render(line))
		case i == len(lines)-1:
			fmt.Fprintf(&buf, "%s%s\n", ColorErrorPrefix.Render("└─"), line)
		default:
			fmt.Fprintf(&buf, "%s %s\n", ColorErrorPrefix.Render("│"), line)
		}
	}
	fmt.Fprintln(&buf)
	return buf.String()
}

func RunTests(files []File) {
	var (
		testsPass int
		testsFail int
	)
	for _, file := range files {
		fmt.Printf("  %s:\n", ColorHeader.Render(file.Name))
		if len(file.Tests) == 0 {
			fmt.Printf("    %s\n", ColorMuted.Render("no tests found"))
		}
		for _, test := range file.Tests {
			if err := ExecuteTest(test); err != nil {
				testsFail++
				fmt.Println(indent(FormatTestError(test, err), "    "))
			} else {
				fmt.Printf("    ✅ %s\n", ColorSuccess.Render(test.Name))
				testsPass++
			}
		}
	}
	fmt.Println()
	fmt.Printf("  %s\n", ColorMuted.Render("---"))
	fmt.Println()
	fmt.Printf("  %s\n", ColorHeader.Render("Results:"))
	if testsPass > 0 {
		fmt.Printf("    Passed: %s\n", ColorSuccess.Render(strconv.Itoa(testsPass)))
	} else {
		fmt.Printf("    Passed: %s\n", ColorMuted.Render(strconv.Itoa(testsPass)))
	}
	if testsFail > 0 {
		fmt.Printf("    Failed: %s\n", ColorErrorText.Render(strconv.Itoa(testsFail)))
	} else {
		fmt.Printf("    Failed: %s\n", ColorMuted.Render(strconv.Itoa(testsFail)))
	}
	fmt.Println()

	if testsFail > 0 {
		os.Exit(1)
	}
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

	if !subcommandInfo.SupportsColoring() {
		return input
	}

	var p printer.Printer = &printer.KubectlOutputColoredPrinter{
		SubcommandInfo:    subcommandInfo,
		Recursive:         subcommandInfo.Recursive,
		ObjFreshThreshold: cfg.ObjFreshThreshold,
		Theme:             &cfg.Theme,
		KubecolorVersion:  "dev",
	}

	if value, ok := os.LookupEnv("INPUT_IS_STDERR"); ok {
		if value != "true" {
			return fmt.Sprintf(`error: var INPUT_IS_STDERR can only be set to "true", but instead got: %q`, value)
		}
		p = &printer.StderrPrinter{
			Theme: &cfg.Theme,
		}
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
			fmt.Fprintf(&buf, "  %s\n", ColorDiffEqual.Render(color.ClearCode(line.Content)))
		case gotextdiff.Insert:
			fmt.Fprintf(&buf, "%s%s\n", ColorDiffAddPrefix.Render("+ "), injectColor(line.Content, ColorDiffAdd))
		case gotextdiff.Delete:
			fmt.Fprintf(&buf, "%s%s\n", ColorDiffDelPrefix.Render("- "), injectColor(line.Content, ColorDiffDel))
		}
	}

	fmt.Fprintf(&buf, "\n%s\n\n", ColorMuted.Render("-----"))

	tabbedLines := quoteAndTabWrite(lines)
	for _, diff := range tabbedLines {
		switch diff.Kind {
		case gotextdiff.Equal:
			fmt.Fprintf(&buf, "  %s\n", ColorDiffEqual.Render(diff.Content))
		case gotextdiff.Insert:
			text := injectColor(highlightEscapedColorCodes(diff.Content, ColorDiffColorHighlight), ColorDiffAdd)
			fmt.Fprintf(&buf, "%s%s\n", ColorDiffAddPrefix.Render("+ "), text)
		case gotextdiff.Delete:
			text := injectColor(highlightEscapedColorCodes(diff.Content, ColorDiffColorHighlight), ColorDiffDel)
			fmt.Fprintf(&buf, "%s%s\n", ColorDiffDelPrefix.Render("- "), text)
		}
	}

	fmt.Fprintf(&buf, "\n%s %s\n", ColorSuccess.Render("(+want)"), ColorErrorText.Render("(-got)"))

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

func injectColor(s string, color color.Color) string {
	newCode := strings.TrimSuffix(strings.TrimPrefix(color.ANSICode(), "\x1b["), "m")

	updatedColors := colorRegex.ReplaceAllStringFunc(s, func(s string) string {
		originalCode := strings.TrimSuffix(strings.TrimPrefix(s, "\x1b["), "m")

		return fmt.Sprintf("\x1b[%s;%sm", originalCode, newCode)
	})
	return color.Render(updatedColors)
}

var escapedColorRegex = regexp.MustCompile(`\\x1b\[[0-9;\.,]+m`)

func highlightEscapedColorCodes(s string, color color.Color) string {
	return escapedColorRegex.ReplaceAllStringFunc(s, func(s string) string {
		return color.Render(s)
	})
}
