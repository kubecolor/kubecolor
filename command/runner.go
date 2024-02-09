package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/kubecolor/kubecolor/printer"
	"github.com/mattn/go-colorable"
)

var (
	Stdout = colorable.NewColorableStdout()
	Stderr = colorable.NewColorableStderr()
)

type Printers struct {
	FullColoredPrinter printer.Printer
	ErrorPrinter       printer.Printer
}

// This is defined here to be replaced in test
var getPrinters = func(subcommandInfo *kubectl.SubcommandInfo, objFreshThreshold time.Duration, theme *color.Theme) *Printers {
	return &Printers{
		FullColoredPrinter: &printer.KubectlOutputColoredPrinter{
			SubcommandInfo:    subcommandInfo,
			Recursive:         subcommandInfo.Recursive,
			ObjFreshThreshold: objFreshThreshold,
			Theme:             theme,
		},
		ErrorPrinter: &printer.WithFuncPrinter{
			Fn: func(line string) color.Color {
				if strings.HasPrefix(strings.ToLower(line), "error") {
					return theme.ErrorColor
				}

				return theme.StringColor
			},
		},
	}
}

func Run(args []string, version string) error {
	config, err := ResolveConfig(args)
	if err != nil {
		return fmt.Errorf("resolve config: %w", err)
	}
	args = config.ArgsPassthrough

	if config.ShowKubecolorVersion {
		fmt.Fprintf(os.Stdout, "%s\n", version)
		return nil
	}

	shouldColorize, subcommandInfo := ResolveSubcommand(args, config)

	cmd := exec.Command(config.KubectlCmd, args...)
	cmd.Stdin = os.Stdin

	// when should not colorize, just run command and return
	if !shouldColorize {
		cmd.Stdout = Stdout
		cmd.Stderr = Stderr
		if err := cmd.Start(); err != nil {
			return err
		}

		// inherit the kubectl exit code
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("%w", &KubectlError{ExitCode: cmd.ProcessState.ExitCode()})
		}
		return nil
	}

	// when colorize, capture stdout and err then colorize it
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// make buffer to be used in defer recover()
	buff := new(bytes.Buffer)
	outReader := io.TeeReader(cmdOut, buff)
	errReader := io.TeeReader(cmdErr, buff)

	if err := cmd.Start(); err != nil {
		return err
	}

	printers := getPrinters(subcommandInfo, config.ObjFreshThreshold, config.Theme)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stdout, buff.String())
			}
		}()

		// This can panic when kubecolor has bug, so recover in defer
		printers.FullColoredPrinter.Print(outReader, Stdout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// This will unlikely panic
		printers.ErrorPrinter.Print(errReader, Stderr)
	}()

	wg.Wait()

	// inherit the kubectl exit code
	if err := cmd.Wait(); err != nil {
		return &KubectlError{ExitCode: cmd.ProcessState.ExitCode()}
	}

	return nil
}
