package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/kubecolor/kubecolor/command"
)

// this is overridden on build time by GoReleaser
var Version string

func main() {
	err := command.Run(os.Args[1:], getVersion())
	if err != nil {
		var ke *command.KubectlError
		if errors.As(err, &ke) {
			os.Exit(ke.ExitCode)
		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] [kubecolor] %s\n", err)
		}
		os.Exit(1)
	}
}

func getVersion() string {
	if Version != "" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "(devel)" {
			return info.Main.Version
		}
		if v, ok := getVCSBuildVersion(info); ok {
			return v
		}
	}
	return "(unset)"
}

func getVCSBuildVersion(info *debug.BuildInfo) (string, bool) {
	var (
		revision string
		dirty    string
	)
	for _, v := range info.Settings {
		switch v.Key {
		case "vcs.revision":
			revision = v.Value
		case "vcs.modified":
			dirty = " (dirty)"
		}
	}
	if revision == "" {
		return "", false
	}
	return revision + dirty, true
}
