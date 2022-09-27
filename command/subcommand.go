package command

import (
	"os"
	"strings"

	"github.com/hidetatz/kubecolor/kubectl"
	"github.com/mattn/go-isatty"
)

// mocked in unit tests
var isOutputTerminal = func() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

func ResolveSubcommand(args []string, config *KubecolorConfig) (bool, *kubectl.SubcommandInfo) {
	// subcommandFound becomes false when subcommand is not found; e.g. "kubecolor --help"
	subcommandInfo, subcommandFound := kubectl.InspectSubcommandInfo(args)

	// if --plain found, it does not colorize
	if config.Plain {
		return false, subcommandInfo
	}

	// if subcommand is not found (e.g. kubecolor --help or just "kubecolor"),
	// it is treated as help because kubectl shows help for such input
	if !subcommandFound {
		// if there is an argument starting with __,
		// the subcommand is probably an internal subcommand (like __completeNoDesc)
		// and should probably not be colorized
		for i := range args {
			if strings.HasPrefix(args[i], "__") {
				return false, subcommandInfo
			}
		}

		subcommandInfo.Help = true
		return true, subcommandInfo
	}

	// when the command output tty is not standard output, shouldColorize depends on --force-colors flag.
	// For example, if the command is run in a shellscript, it should not colorize. (e.g. in "kubectl completion bash")
	// However, if user wants colored output even if the out is not tty (e.g. kubecolor get xx | grep yy)
	// it colorizes the output based on --force-colors.
	if !isOutputTerminal() {
		return config.ForceColor, subcommandInfo
	}

	// else, when the given subcommand is supported, then we colorize it
	return subcommandFound && isColoringSupported(subcommandInfo.Subcommand), subcommandInfo
}

func isColoringSupported(sc kubectl.Subcommand) bool {
	// when you add something here, it won't be colorized
	unsupported := []kubectl.Subcommand{
		kubectl.Create,
		kubectl.Delete,
		kubectl.Edit,
		kubectl.Attach,
		kubectl.Replace,
		kubectl.Completion,
		kubectl.Exec,
		kubectl.Proxy,
		kubectl.Plugin,
		kubectl.Wait,
		kubectl.Run,
		kubectl.Ctx,
		kubectl.Ns,
	}

	for _, u := range unsupported {
		if sc == u {
			return false
		}
	}

	return true
}
