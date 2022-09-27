package main

import (
	"errors"
	"os"

	"github.com/kubecolor/kubecolor/command"
)

// this is overridden on build time by GoReleaser
var Version string = "unset"

func main() {
	err := command.Run(os.Args[1:], Version)
	if err != nil {
		var ke *command.KubectlError
		if errors.As(err, &ke) {
			os.Exit(ke.ExitCode)
		}
		os.Exit(1)
	}
}
