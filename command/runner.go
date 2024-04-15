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

	"github.com/kubecolor/kubecolor/config"
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
var getPrinters = func(subcommandInfo *kubectl.SubcommandInfo, objFreshThreshold time.Duration, theme *config.Theme) *Printers {
	return &Printers{
		FullColoredPrinter: &printer.KubectlOutputColoredPrinter{
			SubcommandInfo:    subcommandInfo,
			Recursive:         subcommandInfo.Recursive,
			ObjFreshThreshold: objFreshThreshold,
			Theme:             theme,
		},
		ErrorPrinter: &printer.WithFuncPrinter{
			Fn: func(line string) config.Color {
				if strings.HasPrefix(strings.ToLower(line), "error") {
					return theme.Stderr.Error
				}

				return theme.Stderr.Default
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

	// when should not colorize, just run command and return
	if !shouldColorize {
		return execWithoutColors(config, args)
	}

	stdoutReader, stderrReader, err := execWithReaders(config, args)
	if err != nil {
		return err
	}
	// in case of panic, at least let kubectl finish executing
	defer stdoutReader.Close()
	defer stderrReader.Close()

	// make buffer to be used in defer recover()
	errBuf := new(bytes.Buffer)
	errBufReader := io.TeeReader(stderrReader, errBuf)

	printers := getPrinters(subcommandInfo, config.ObjFreshThreshold, config.Theme)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "[kubecolor] [ERROR] Recovered from panic: %v\n", r)
				fmt.Fprint(os.Stdout, errBuf.String())
			}
		}()

		// This can panic when kubecolor has bug, so recover in defer
		printers.FullColoredPrinter.Print(stdoutReader, Stdout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// This will unlikely panic
		printers.ErrorPrinter.Print(errBufReader, Stderr)
	}()

	wg.Wait()

	return stdoutReader.Close()
}

func execWithoutColors(config *Config, args []string) error {
	fmt.Println("in execWithoutColors")
	if config.StdinOverride != "" {
		r, err := getStdinOverrideReader(config.StdinOverride)
		if err != nil {
			return err
		}
		defer r.Close()
		_, err = io.Copy(Stdout, r)
		return err
	}

	cmd := exec.Command(config.KubectlCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr

	// when should not colorize, just run command and return
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w", &KubectlError{ExitCode: cmd.ProcessState.ExitCode()})
	}

	return nil
}

func execWithReaders(config *Config, args []string) (io.ReadCloser, io.ReadCloser, error) {
	if config.StdinOverride != "" {
		stdout, err := getStdinOverrideReader(config.StdinOverride)
		return stdout, nopReadCloser{}, err
	}

	cmd := exec.Command(config.KubectlCmd, args...)
	cmd.Stdin = os.Stdin

	// when colorize, capture stdout and err then colorize it
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return &cmdWaitReadCloser{cmd: cmd, stdout: cmdOut}, cmdErr, nil
}

func getStdinOverrideReader(stdinOverride string) (io.ReadCloser, error) {
	if stdinOverride == "-" {
		return os.Stdin, nil
	}
	file, err := os.Open(stdinOverride)
	if err != nil {
		return nil, fmt.Errorf("read file specified by --kubecolor-stdin: %w", err)
	}
	return file, nil
}

type nopReadCloser struct{}

func (nopReadCloser) Read(b []byte) (int, error) { return 0, io.EOF }

func (nopReadCloser) Close() error { return nil }

type cmdWaitReadCloser struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser

	closed  bool
	lastErr error
}

func (r *cmdWaitReadCloser) Read(b []byte) (int, error) {
	return r.stdout.Read(b)
}

func (r *cmdWaitReadCloser) Close() error {
	if r.closed {
		return r.lastErr
	}
	if err := r.stdout.Close(); err != nil {
		r.lastErr = err
		return err
	}
	if err := r.cmd.Wait(); err != nil {
		r.lastErr = err
		return &KubectlError{ExitCode: r.cmd.ProcessState.ExitCode()}
	}
	r.closed = true
	return nil
}
