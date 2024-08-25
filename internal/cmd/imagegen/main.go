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

	"github.com/kubecolor/kubecolor/command"
	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/kubecolor/kubecolor/printer"
)

var flags = struct {
	outputDir    string
	outputFile   string
	freezeConfig string
	preset       string
	width        int

	prompt       string
	promptColor  color.Color
	cmd          string
	cmdColor     color.Color
	argColor     color.Color
	valueColor   color.Color
	keywordColor color.Color
	flagColor    color.Color
}{
	outputDir:    "./docs",
	freezeConfig: "./docs/freeze-config.json",
	preset:       "dark",
	width:        100,

	prompt:       "❯",
	promptColor:  color.MustParseColor("green"),
	cmd:          "kubectl",
	cmdColor:     color.MustParseColor("green"),
	argColor:     color.MustParseColor("none"),
	valueColor:   color.MustParseColor("yellow"),
	keywordColor: color.MustParseColor("magenta"),
	flagColor:    color.MustParseColor("cyan"),
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

	printed, err := parseAndPrintCommand(string(input), &EnvStore{Vars: []EnvVar{
		{Key: "KUBECOLOR_THEME_PRESET", Value: flags.preset},
	}})
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

type EnvStore struct {
	Vars []EnvVar
}

func (e *EnvStore) Add(key, value string) {
	e.Vars = append(e.Vars, EnvVar{
		Key:   key,
		Value: value,
	})
}

type EnvVar struct {
	Key   string
	Value string
}

func parseAndPrintCommand(input string, env *EnvStore) (string, error) {
	commands, err := splitInputIntoCommands(input)
	if err != nil {
		return "", err
	}
	if len(commands) == 0 {
		return "", fmt.Errorf("no commands found")
	}

	var buf bytes.Buffer
	for i, cmd := range commands {
		if i > 0 {
			buf.WriteByte('\n')
		}
		commandOutput, err := runCommand(cmd, env)
		if err != nil {
			return "", err
		}
		commandOutput = setWidth(commandOutput, flags.width)
		buf.WriteString(commandOutput)
	}
	return buf.String(), nil
}

func runCommand(cmd Command, env *EnvStore) (string, error) {
	switch cmd.Exe {
	case "echo", "print":
		return runNoopCommand(cmd, env)
	case "export":
		return runExportCommand(cmd, env)
	case "kubectl", "kubecolor":
		return runKubecolorCommand(cmd, env)
	default:
		return "", fmt.Errorf("unsupported command: %q", cmd.Exe)
	}
}

func runNoopCommand(cmd Command, _ *EnvStore) (string, error) {
	var buf bytes.Buffer

	buf.WriteString(flags.promptColor.Render(flags.prompt))
	buf.WriteByte(' ')
	buf.WriteString(flags.cmdColor.Render(cmd.Exe))
	for _, arg := range cmd.Args {
		buf.WriteByte(' ')
		if strings.HasPrefix(arg, "-") {
			buf.WriteString(flags.flagColor.Render(arg))
		} else {
			buf.WriteString(flags.argColor.Render(arg))
		}
	}
	buf.WriteByte('\n')

	buf.WriteString(cmd.Input)
	buf.WriteByte('\n')

	return buf.String(), nil
}

func runExportCommand(cmd Command, env *EnvStore) (string, error) {
	var buf bytes.Buffer

	buf.WriteString(flags.promptColor.Render(flags.prompt))
	buf.WriteByte(' ')
	buf.WriteString(flags.keywordColor.Render(cmd.Exe))
	buf.WriteByte(' ')
	key, value := cmd.Args[0], cmd.Args[1]
	buf.WriteString(flags.argColor.Render(key))
	buf.WriteByte('=')
	buf.WriteString(flags.valueColor.Render(value))
	buf.WriteByte('\n')

	env.Add(key, value)

	return buf.String(), nil
}

func runKubecolorCommand(cmd Command, env *EnvStore) (string, error) {
	defer restoreEnv(os.Environ())
	os.Clearenv()

	var buf bytes.Buffer

	buf.WriteString(flags.promptColor.Render(flags.prompt))
	buf.WriteByte(' ')
	buf.WriteString(flags.cmdColor.Render(flags.cmd))
	for _, arg := range cmd.Args {
		buf.WriteByte(' ')
		if strings.HasPrefix(arg, "-") {
			buf.WriteString(flags.flagColor.Render(arg))
		} else {
			buf.WriteString(flags.argColor.Render(arg))
		}
	}
	buf.WriteByte('\n')

	for _, e := range env.Vars {
		if err := os.Setenv(e.Key, e.Value); err != nil {
			return "", fmt.Errorf("set env %s=%q: %s", e.Key, e.Value, err)
		}
	}

	v := config.NewViper()
	cfg, err := command.ResolveConfigViper(cmd.Args, v)
	if err != nil {
		return "", err
	}
	cfg.ForceColor = command.ColorLevelTrueColor

	subcommandInfo := kubectl.InspectSubcommandInfo(cmd.Args, kubectl.NoopPluginHandler{})
	p := &printer.KubectlOutputColoredPrinter{
		SubcommandInfo:    subcommandInfo,
		Recursive:         subcommandInfo.Recursive,
		ObjFreshThreshold: cfg.ObjFreshThreshold,
		Theme:             &cfg.Theme,
	}
	p.Print(strings.NewReader(cmd.Input), &buf)
	return buf.String(), nil
}

type Command struct {
	Exe   string
	Args  []string
	Input string
}

func splitInputIntoCommands(input string) ([]Command, error) {
	if !strings.HasPrefix(input, "$") {
		return nil, fmt.Errorf(`expected first line to start with '$ ...'`)
	}

	lines := strings.Split(input, "\n")

	var commands []Command
	var currentCommand Command

	for _, line := range lines {
		if cmdline, ok := strings.CutPrefix(line, "$"); ok {
			// new command
			cmdline = strings.TrimSpace(cmdline)
			if currentCommand.Exe != "" {
				currentCommand.Input = strings.TrimSpace(currentCommand.Input)
				commands = append(commands, currentCommand)
			}
			args := strings.Fields(cmdline)
			if len(args) == 0 {
				return nil, fmt.Errorf("missing command on line that starts with '$ ...'")
			}

			if args[0] == "export" {
				if len(args) != 2 {
					return nil, fmt.Errorf("export command must be in format '$ export KEY=VALUE'")
				}
				key, value, ok := strings.Cut(args[1], "=")
				if !ok {
					return nil, fmt.Errorf("export command must be in format '$ export KEY=VALUE'")
				}
				args = []string{"export", key, value}
			}

			currentCommand = Command{
				Exe:  args[0],
				Args: args[1:],
			}
			continue
		}

		if currentCommand.Input != "" {
			currentCommand.Input += "\n"
		}
		currentCommand.Input += line
	}

	if currentCommand.Exe != "" {
		currentCommand.Input = strings.TrimSpace(currentCommand.Input)
		commands = append(commands, currentCommand)
	}

	return commands, nil
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
