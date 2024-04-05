package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/kubecolor/kubecolor/command"
	"github.com/kubecolor/kubecolor/config"
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
				fmt.Printf("    ❌ %s\n", colorErrorPrefix.Render(test.Name))
				lines := strings.Split(err.Error(), "\n")
				for i, line := range lines {
					if i == 0 {
						fmt.Printf("     %s %s\n", colorErrorPrefix.Render("│"), colorErrorText.Render(line))
					} else if i == len(lines)-1 {
						fmt.Printf("     %s%s\n", colorErrorPrefix.Render("└─"), line)
					} else {
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
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")
	diff := cmp.Diff(wantLines, gotLines)
	diffLines := strings.Split(diff, "\n")

	var buf bytes.Buffer
	for _, line := range diffLines {
		if strings.HasPrefix(line, "+") {
			fmt.Fprintln(&buf, colorSuccess.Render(line))
		} else if strings.HasPrefix(line, "-") {
			fmt.Fprintln(&buf, colorErrorText.Render(line))
		} else {
			fmt.Fprintln(&buf, colorMuted.Render(line))
		}
	}
	return buf.String()
}
