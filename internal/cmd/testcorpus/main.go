package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/kubecolor/kubecolor/internal/testcorpus"
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
	files, err := testcorpus.ParseGlob(flags.glob)
	if err != nil {
		logErrorf("%s", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		logWarnf("glob did not match any files: %s", flags.glob)
		os.Exit(0)
	}

	if flags.update {
		testcorpus.UpdateTests(files)
	} else {
		testcorpus.RunTests(files)
	}
}

func logErrorf(format string, args ...any) {
	fmt.Printf("  %s %s\n",
		testcorpus.ColorErrorPrefix.Render("error:"),
		testcorpus.ColorErrorText.Render(fmt.Sprintf(format, args...)))
}

func logWarnf(format string, args ...any) {
	fmt.Printf("  %s  %s\n",
		testcorpus.ColorWarnPrefix.Render("warn:"),
		testcorpus.ColorWarnText.Render(fmt.Sprintf(format, args...)))
}
