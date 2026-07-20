package command

import (
	"os"
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_ResolveConfig(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		env          map[string]string
		expectedConf *Config
	}{
		{
			name: "no config",
			args: []string{"get", "pods"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: time.Duration(0),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetDark,
				},
				ArgsPassthrough: []string{"get", "pods"},
				ForceColor:      ColorLevelUnset,
			},
		},
		{
			name: "plain, light, force",
			args: []string{"get", "pods", "--plain", "--light-background", "--force-colors"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: time.Duration(0),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.LightTheme,
					Preset:            config.PresetLight,
				},
				ForceColor:      ColorLevelAuto,
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "KUBECTL_COMMAND exists",
			args: []string{"get", "pods", "--plain"},
			env:  map[string]string{"KUBECTL_COMMAND": "kubectl.1.19"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl.1.19",
					ObjFreshThreshold: time.Duration(0),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetDark,
				},
				ForceColor:      ColorLevelNone,
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "KUBECOLOR_OBJ_FRESH exists",
			args: []string{"get", "pods"},
			env:  map[string]string{"KUBECOLOR_OBJ_FRESH": "1m"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: time.Minute,
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetDark,
				},
				ForceColor:      ColorLevelUnset,
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "KUBECOLOR_LIGHT_BACKGROUND via env",
			args: []string{"get", "pods"},
			env:  map[string]string{"KUBECOLOR_LIGHT_BACKGROUND": "true"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl: "kubectl",
					Paging:  config.PagingDefault,
					Theme:   *testconfig.LightTheme,
					Preset:  config.PresetLight,
				},
				ForceColor:      ColorLevelUnset,
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "KUBECOLOR_FORCE_COLORS env var bool",
			args: []string{"get", "pods"},
			env:  map[string]string{"KUBECOLOR_FORCE_COLORS": "true"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl: "kubectl",
					Paging:  config.PagingDefault,
					Theme:   *testconfig.DarkTheme,
					Preset:  config.PresetDark,
				},
				ForceColor:      ColorLevelAuto,
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "KUBECOLOR_FORCE_COLORS env var truecolor",
			args: []string{"get", "pods"},
			env:  map[string]string{"KUBECOLOR_FORCE_COLORS": "truecolor"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl: "kubectl",
					Paging:  config.PagingDefault,
					Theme:   *testconfig.DarkTheme,
					Preset:  config.PresetDark,
				},
				ForceColor:      ColorLevelTrueColor,
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "Pager flags overwrite (1)",
			args: []string{"--paging", "--pager=most", "get", "pods"},
			env: map[string]string{
				"KUBECOLOR_PAGING": "never",
				"KUBECOLOR_PAGER":  "more",
			},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl: "kubectl",
					Pager:   "most",
					Paging:  config.PagingAuto,
					Theme:   *testconfig.DarkTheme,
					Preset:  config.PresetDark,
				},
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "Pager flags overwrite (2)",
			args: []string{"--no-paging", "get", "pods"},
			env: map[string]string{
				"KUBECOLOR_PAGING": string(config.PagingAuto),
			},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl: "kubectl",
					Paging:  config.PagingNever,
					Theme:   *testconfig.DarkTheme,
					Preset:  config.PresetDark,
				},
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		// Issue #201: interactive kubectl commands must imply plain mode so that
		// kubecolor doesn't capture/colorize the prompt streams. The kubectl
		// arguments are preserved unchanged.
		{
			name: "interactive short flag -i implies plain mode",
			args: []string{"delete", "pod", "nginx", "-i"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: time.Duration(0),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetDark,
				},
				ForceColor:      ColorLevelNone,
				ArgsPassthrough: []string{"delete", "pod", "nginx", "-i"},
			},
		},
		{
			name: "interactive long flag --interactive implies plain mode",
			args: []string{"delete", "pod", "nginx", "--interactive"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: time.Duration(0),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetDark,
				},
				ForceColor:      ColorLevelNone,
				ArgsPassthrough: []string{"delete", "pod", "nginx", "--interactive"},
			},
		},
		{
			name: "interactive --interactive=true implies plain mode",
			args: []string{"delete", "pod", "nginx", "--interactive=true"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: time.Duration(0),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetDark,
				},
				ForceColor:      ColorLevelNone,
				ArgsPassthrough: []string{"delete", "pod", "nginx", "--interactive=true"},
			},
		},
		{
			name: "explicit --interactive=false does NOT imply plain mode",
			args: []string{"delete", "pod", "nginx", "--interactive=false"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: time.Duration(0),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetDark,
				},
				ForceColor:      ColorLevelUnset,
				ArgsPassthrough: []string{"delete", "pod", "nginx", "--interactive=false"},
			},
		},
		{
			name: "interactive implies plain mode even with explicit --force-colors",
			args: []string{"delete", "pod", "nginx", "-i", "--force-colors=truecolor"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: time.Duration(0),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetDark,
				},
				ForceColor:      ColorLevelNone,
				ArgsPassthrough: []string{"delete", "pod", "nginx", "-i"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				testutil.Setenv(t, k, v)
			}

			conf, err := ResolveConfig(tt.args)
			testutil.MustNoError(t, err)

			// Don't test flags field
			tt.expectedConf.Flags = conf.Flags

			testutil.MustEqual(t, tt.expectedConf, conf)
		})
	}
}

func TestHasInteractiveFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "short flag", args: []string{"delete", "-i"}, want: true},
		{name: "long flag", args: []string{"delete", "--interactive"}, want: true},
		{name: "short explicit false", args: []string{"delete", "-i=false"}, want: false},
		{name: "later flag wins", args: []string{"delete", "--interactive=false", "-i"}, want: true},
		{name: "later false wins", args: []string{"delete", "-i", "--interactive=false"}, want: false},
		{name: "separator stops parsing", args: []string{"get", "pod", "--", "--interactive"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasInteractiveFlag(tt.args); got != tt.want {
				t.Fatalf("hasInteractiveFlag(%q) = %t, want %t", tt.args, got, tt.want)
			}
		})
	}
}
