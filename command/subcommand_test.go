package command

import (
	"testing"

	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_ResolveSubcommand(t *testing.T) {
	tests := []struct {
		name                   string
		args                   []string
		conf                   *Config
		isOutputTerminal       func() bool
		expectedShouldColorize bool
		expectedInfo           *kubectl.SubcommandInfo
	}{
		{
			name:             "basic case",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &Config{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
				Theme:      testconfig.DarkTheme,
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when plain, it won't colorize",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &Config{
				Plain:      true,
				ForceColor: false,
				KubectlCmd: "kubectl",
				Theme:      testconfig.DarkTheme,
			},
			expectedShouldColorize: false,
			expectedInfo:           nil,
		},
		{
			name:             "when help, it will colorize",
			args:             []string{"get", "pods", "-h"},
			isOutputTerminal: func() bool { return true },
			conf: &Config{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
				Theme:      testconfig.DarkTheme,
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get, Help: true},
		},
		{
			name:             "when both plain and force, plain is chosen",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &Config{
				Plain:      true,
				ForceColor: true,
				KubectlCmd: "kubectl",
				Theme:      testconfig.DarkTheme,
			},
			expectedShouldColorize: false,
			expectedInfo:           nil,
		},
		{
			name:             "when no subcommand is found, it becomes help",
			args:             []string{},
			isOutputTerminal: func() bool { return true },
			conf: &Config{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
				Theme:      testconfig.DarkTheme,
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Help: true},
		},
		{
			name:             "when not tty, it won't colorize",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return false },
			conf: &Config{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
				Theme:      testconfig.DarkTheme,
			},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "even if not tty, if force, it colorizes",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return false },
			conf: &Config{
				Plain:      false,
				ForceColor: true,
				KubectlCmd: "kubectl",
				Theme:      testconfig.DarkTheme,
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when the subcommand is unsupported, it won't colorize",
			args:             []string{"-h"},
			isOutputTerminal: func() bool { return true },
			conf: &Config{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
				Theme:      testconfig.DarkTheme,
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Help: true},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			isOutputTerminal = tt.isOutputTerminal
			shouldColorize, info := ResolveSubcommand(tt.args, tt.conf)
			testutil.MustEqual(t, tt.expectedShouldColorize, shouldColorize)
			testutil.MustEqual(t, tt.expectedInfo, info)
		})
	}
}
