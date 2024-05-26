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

// when you add something here, it won't be colorized
var unsupported = map[kubectl.Subcommand]struct{}{
	kubectl.Attach:           {},
	kubectl.Completion:       {},
	kubectl.Create:           {},
	kubectl.Ctx:              {},
	kubectl.Debug:            {},
	kubectl.Delete:           {},
	kubectl.Edit:             {},
	kubectl.Exec:             {},
	kubectl.InternalComplete: {},
	kubectl.KubectlPlugin:    {},
	kubectl.Ns:               {},
	kubectl.Plugin:           {},
	kubectl.Proxy:            {},
	kubectl.Replace:          {},
	kubectl.Run:              {},
	kubectl.Wait:             {},

	// oc (OpenShift CLI) specific subcommands
	kubectl.Rsh: {},
}

func isColoringSupported(sc kubectl.Subcommand) bool {
	_, found := unsupported[sc]
	return !found
}
