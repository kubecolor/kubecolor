package command

import (
	"os"

	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/mattn/go-isatty"
)

// mocked in unit tests
var isOutputTerminal = func() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

var pluginHandler kubectl.PluginHandler = kubectl.DefaultPluginHandler{}

func ResolveSubcommand(args []string, config *Config) (bool, *kubectl.SubcommandInfo) {
	// if --plain found, it does not colorize
	if config.Plain {
		return false, nil
	}

	// subcommandFound becomes false when subcommand is not found; e.g. "kubecolor --help"
	subcommandInfo, subcommandFound := kubectl.InspectSubcommandInfo(args, pluginHandler)

	// if subcommand is not found (e.g. kubecolor --help or just "kubecolor"),
	// it is treated as help because kubectl shows help for such input
	if !subcommandFound {
		subcommandInfo.Help = true
		return true, subcommandInfo
	}

	if !subcommandInfo.SupportsColoring() {
		return false, subcommandInfo
	}

	// when the command output tty is not standard output, shouldColorize depends on --force-colors flag.
	// For example, if the command is run in a shellscript, it should not colorize. (e.g. in "kubectl completion bash")
	// However, if user wants colored output even if the out is not tty (e.g. kubecolor get xx | grep yy)
	// it colorizes the output based on --force-colors.
	if !isOutputTerminal() {
		return config.ForceColor, subcommandInfo
	}

	// else, when the given subcommand is supported, then we colorize it
	return subcommandFound, subcommandInfo
}
