package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/kubecolor/kubecolor/command"
	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/printer"
)

var (
	colorErrorPrefix = config.MustParseColor("hi-red:bold")
	colorErrorText   = config.MustParseColor("red")
	colorWarnPrefix  = config.MustParseColor("hi-yellow:bold")
	colorWarnText    = config.MustParseColor("yellow")

	colorHeader  = config.MustParseColor("bold")
	colorMuted   = config.MustParseColor("gray")
	colorSuccess = config.MustParseColor("green")
)

var flags = struct {
	glob   string
	update bool
}{
	glob: "test/corpus/*.txt",
}

func init() {
	flag.StringVar(&flags.glob, "glob", flags.glob, "Glob pattern to find test files")
	flag.BoolVar(&flags.update, "update", flags.update, "Update all outputs in corpus files with current kubecolor output")

	color.ForceColor()
	color.Enable = true
}

func main() {
	flag.Parse()

	fmt.Println()
	files, err := ParseGlob(flags.glob)
	if err != nil {
		logErrorf("%s", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		logWarnf("glob did not match any files: %s", flags.glob)
		os.Exit(0)
	}

	if flags.update {
		UpdateTests(files)
	} else {
		RunTests(files)
	}
}

func printCommand(args []string, input string) string {
	cfg := &command.Config{ForceColor: true, Theme: testconfig.DarkTheme}
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

func logErrorf(format string, args ...any) {
	fmt.Printf("  %s %s\n", colorErrorPrefix.Render("error:"), colorErrorText.Render(fmt.Sprintf(format, args...)))
}

func logWarnf(format string, args ...any) {
	fmt.Printf("  %s  %s\n", colorWarnPrefix.Render("warn:"), colorWarnText.Render(fmt.Sprintf(format, args...)))
}
