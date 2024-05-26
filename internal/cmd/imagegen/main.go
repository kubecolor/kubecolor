package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/kubecolor/kubecolor/command"
	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/kubecolor/kubecolor/printer"
)

var flags = struct {
	outputDir    string
	outputFile   string
	freezeConfig string
	preset       string
	width        int

	prompt      string
	promptColor config.Color
	cmd         string
	cmdColor    config.Color
	argColor    config.Color
	flagColor   config.Color
}{
	outputDir:    "./docs",
	freezeConfig: "./docs/freeze-config.json",
	preset:       "dark",
	width:        100,

	prompt:      "❯",
	promptColor: config.MustParseColor("green"),
	cmd:         "kubectl",
	cmdColor:    config.MustParseColor("green"),
	argColor:    config.MustParseColor("none"),
	flagColor:   config.MustParseColor("cyan"),
}

func init() {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), `Usage: go run ./cmd/imagegen <path to test file>

Examples:
  go run ./internal/cmd/imagegen ./docs/kubectl-get-pods.txt
  go run ./internal/cmd/imagegen ./docs/kubectl-describe-pod.txt

Flags:
`)
		flag.PrintDefaults()
	}

	flag.StringVar(&flags.outputDir, "output-dir", flags.outputDir, "Default directory to output to")
	flag.StringVar(&flags.outputFile, "output", flags.outputFile, `Path to output to (default "${-output-dir flag}")/${filename of test file}.svg")`)
	flag.StringVar(&flags.freezeConfig, "freeze-config", flags.freezeConfig, "Path to the charmbracelet freeze config file")
	flag.StringVar(&flags.preset, "preset", flags.preset, "Kubecolor theme preset")
	flag.IntVar(&flags.width, "width", flags.width, "Terminal width in output image")
	flag.StringVar(&flags.prompt, "prompt", flags.prompt, "Shell prompt used in output image")
	flag.Var(&flags.promptColor, "prompt-color", "Color of the shell prompt used in output image")
	flag.StringVar(&flags.cmd, "cmd", flags.cmd, "Shell command used in output image")
	flag.Var(&flags.cmdColor, "cmd-color", "Color of the shell command used in output image")
	flag.Var(&flags.argColor, "arg-color", "Color of shell command arguments used in output image")
	flag.Var(&flags.flagColor, "flag-color", "Color of shell command flags used in output image")
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputPath := flag.Arg(0)
	input, err := os.ReadFile(inputPath)
	if err != nil {
		slog.Error("Failed to read input file", "error", err)
		os.Exit(1)
	}

	printed, err := parseAndPrintCommand(string(input), []EnvVar{
		{Key: "KUBECOLOR_THEME_PRESET", Value: flags.preset},
	})
	if err != nil {
		slog.Error("Failed to print command via kubecolor", "error", err)
		os.Exit(1)
	}

	outputPath := flags.outputFile
	if outputPath == "" {
		_, inputFilename := filepath.Split(inputPath)
		inputExt := filepath.Ext(inputFilename)
		outputFilename := strings.TrimSuffix(inputFilename, inputExt) + ".svg"
		outputPath = filepath.Join(flags.outputDir, outputFilename)
	}

	if err := runFreeze(printed, outputPath, flags.freezeConfig); err != nil {
		os.Exit(1)
	}
}

func runFreeze(inputText, outputPath, configPath string) error {
	cmd := freezeCmd("--language=ansi", "--output="+outputPath, "--config="+configPath)
	cmd.Stdin = strings.NewReader(inputText)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) == 0 {
			out = []byte(err.Error())
		}
		slog.Error("Freeze failed", "error", out)
		return err
	}
	slog.Info(string(out))
	return nil
}

func freezeCmd(args ...string) *exec.Cmd {
	slog.Debug("Running freeze", "args", args)
	freezePath, err := exec.LookPath("freeze")
	if err == nil {
		return exec.Command(freezePath, args...)
	}

	return exec.Command("go", append([]string{"run", "github.com/charmbracelet/freeze@latest"}, args...)...)
}

type EnvVar struct {
	Key   string
	Value string
}

func parseAndPrintCommand(input string, env []EnvVar) (string, error) {
	args, commandInput, err := splitInputIntoArgsAndOutput(input)
	if err != nil {
		return "", err
	}

	commandOutput, err := printCommand(args, commandInput, env)
	if err != nil {
		return "", err
	}
	commandOutput = setWidth(commandOutput, flags.width)

	var buf bytes.Buffer

	buf.WriteString(flags.promptColor.Render(flags.prompt))
	buf.WriteByte(' ')
	buf.WriteString(flags.cmdColor.Render(flags.cmd))
	for _, arg := range args {
		buf.WriteByte(' ')
		if strings.HasPrefix(arg, "-") {
			buf.WriteString(flags.flagColor.Render(arg))
		} else {
			buf.WriteString(flags.argColor.Render(arg))
		}
	}
	buf.WriteByte('\n')

	buf.WriteString(commandOutput)
	return buf.String(), nil
}

func printCommand(args []string, commandInput string, env []EnvVar) (string, error) {
	defer restoreEnv(os.Environ())
	os.Clearenv()

	for _, e := range env {
		if err := os.Setenv(e.Key, e.Value); err != nil {
			return "", fmt.Errorf("set env %s=%q: %s", e.Key, e.Value, err)
		}
	}

	v := config.NewViper()
	cfg, err := command.ResolveConfigViper(args, v)
	if err != nil {
		return "", err
	}
	cfg.ForceColor = command.ColorLevelTrueColor

	subcommandInfo := kubectl.InspectSubcommandInfo(args, kubectl.NoopPluginHandler{})
	p := &printer.KubectlOutputColoredPrinter{
		SubcommandInfo:    subcommandInfo,
		Recursive:         subcommandInfo.Recursive,
		ObjFreshThreshold: cfg.ObjFreshThreshold,
		Theme:             cfg.Theme,
	}
	var buf bytes.Buffer
	p.Print(strings.NewReader(commandInput), &buf)
	return buf.String(), nil
}

func splitInputIntoArgsAndOutput(input string) ([]string, string, error) {
	cmdline, rest, _ := strings.Cut(input, "\n")
	args := strings.Split(cmdline, " ")

	if len(args) < 2 || args[0] != "$" || args[1] != "kubectl" {
		return nil, "", fmt.Errorf(`expected first line to start with '$ kubectl ...', but got: '%s'`, cmdline)
	}

	return args[2:], strings.TrimSpace(rest), nil
}

func setWidth(s string, width int) string {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		lineLen := len(color.ClearCode(line))
		switch {
		case lineLen > width:
			buf.WriteString(line[:width])
		case lineLen < width:
			buf.WriteString(line)
			buf.WriteString(strings.Repeat(" ", width-lineLen))
		default:
			buf.WriteString(line)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func restoreEnv(oldEnv []string) {
	os.Clearenv()
	for _, e := range oldEnv {
		k, v, _ := strings.Cut(e, "=")
		os.Setenv(k, v)
	}
}
