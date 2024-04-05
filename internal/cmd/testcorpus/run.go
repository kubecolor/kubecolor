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
	"github.com/kubecolor/kubecolor/command"
	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/printer"
	"github.com/sergi/go-diff/diffmatchpatch"
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
				fmt.Printf("    ❌ %s\n", colorErrorPrefix.Render(test.Name))
				lines := strings.Split(err.Error(), "\n")
				for i, line := range lines {
					switch {
					case i == 0:
						fmt.Printf("     %s %s\n", colorErrorPrefix.Render("│"), colorErrorText.Render(line))
					case i == len(lines)-1:
						fmt.Printf("     %s%s\n", colorErrorPrefix.Render("└─"), line)
					default:
						fmt.Printf("     %s %s\n", colorErrorPrefix.Render("│"), line)
					}
				}
				testsFail++
				fmt.Println()
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
		return fmt.Errorf("input does not match output:\n%s", createColoredDiff(test.Output, gotOutput))
	}
	return nil
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
	cfg.ForceColor = true

	shouldColorize, subcommandInfo := command.ResolveSubcommand(args, cfg)
	if !shouldColorize {
		return input
	}
	p := &printer.KubectlOutputColoredPrinter{
		SubcommandInfo:    subcommandInfo,
		Recursive:         subcommandInfo.Recursive,
		ObjFreshThreshold: cfg.ObjFreshThreshold,
		Theme:             cfg.Theme,
	}
	var buf bytes.Buffer
	p.Print(strings.NewReader(input), &buf)
	return buf.String()
}

func createColoredDiff(want, got string) string {
	dmp := diffmatchpatch.New()

	wSrc, wDst, linesArray := dmp.DiffLinesToRunes(want, got)
	diffs := dmp.DiffMainRunes(wSrc, wDst, false)
	diffs = dmp.DiffCharsToLines(diffs, linesArray)
	diffs = splitDiffsIntoLines(diffs)

	var buf bytes.Buffer
	buf.WriteByte('\n')

	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			fmt.Fprintf(&buf, "  %s\n", colorDiffEqual.Render(color.ClearCode(diff.Text)))
		case diffmatchpatch.DiffInsert:
			fmt.Fprintf(&buf, "%s%s\n", colorDiffAddPrefix.Render("+ "), injectColor(diff.Text, colorDiffAdd))
		case diffmatchpatch.DiffDelete:
			fmt.Fprintf(&buf, "%s%s\n", colorDiffDelPrefix.Render("- "), injectColor(diff.Text, colorDiffDel))
		}
	}

	fmt.Fprintf(&buf, "\n%s\n\n", colorMuted.Render("-----"))

	tabbedDiffs := quoteAndTabWrite(diffs)
	for _, diff := range tabbedDiffs {
		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			fmt.Fprintf(&buf, "  %s\n", colorDiffEqual.Render(diff.Text))
		case diffmatchpatch.DiffInsert:
			text := injectColor(highlightEscapedColorCodes(diff.Text, colorDiffColorHighlight), colorDiffAdd)
			fmt.Fprintf(&buf, "%s%s\n", colorDiffAddPrefix.Render("+ "), text)
		case diffmatchpatch.DiffDelete:
			text := injectColor(highlightEscapedColorCodes(diff.Text, colorDiffColorHighlight), colorDiffDel)
			fmt.Fprintf(&buf, "%s%s\n", colorDiffDelPrefix.Render("- "), text)
		}
	}

	fmt.Fprintf(&buf, "\n%s %s\n", colorSuccess.Render("(+want)"), colorErrorText.Render("(-got)"))

	return buf.String()
}

func splitDiffsIntoLines(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	var result []diffmatchpatch.Diff
	for _, diff := range diffs {
		lines := strings.Split(strings.Trim(diff.Text, "\n\r"), "\n")
		for _, line := range lines {
			match := strings.Trim(line, "\n\r")
			result = append(result, diffmatchpatch.Diff{Type: diff.Type, Text: match})
		}
	}
	return result
}

var tripleSpaceRegex = regexp.MustCompile(` {3,}`)

func quoteAndTabWrite(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	for _, diff := range diffs {
		quoted := fmt.Sprintf("%q", diff.Text)
		tabbed := tripleSpaceRegex.ReplaceAllLiteralString(quoted, "\t")
		fmt.Fprintln(tw, tabbed)
	}
	tw.Flush()
	lines := strings.Split(buf.String(), "\n")
	newDiffs := make([]diffmatchpatch.Diff, len(diffs))
	for i, diff := range diffs {
		newDiffs[i] = diffmatchpatch.Diff{Type: diff.Type, Text: lines[i]}
	}
	return newDiffs
}

var colorRegex = regexp.MustCompile(`\x1b\[[0-9;\.,]+m`)

func injectColor(s string, color config.Color) string {
	newCode := strings.TrimSuffix(strings.TrimPrefix(color.Code, "\x1b["), "m")

	updatedColors := colorRegex.ReplaceAllStringFunc(s, func(s string) string {
		originalCode := strings.TrimSuffix(strings.TrimPrefix(s, "\x1b["), "m")

		return fmt.Sprintf("\x1b[%s;%sm", originalCode, newCode)
	})
	return color.Render(updatedColors)
}

var escapedColorRegex = regexp.MustCompile(`\\x1b\[[0-9;\.,]+m`)

func highlightEscapedColorCodes(s string, color config.Color) string {
	return escapedColorRegex.ReplaceAllStringFunc(s, func(s string) string {
		return colorDiffColorHighlight.Render(s)
	})
}
