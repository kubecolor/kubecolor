package command

import (
	"os"
	"testing"

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
		term         config.TermInfo
	}{
		{
			name: "no config",
			args: []string{"get", "pods"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:           "kubectl",
					ObjFreshThreshold: nil,
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetAuto,
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
					ObjFreshThreshold: nil,
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
					ObjFreshThreshold: nil,
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetAuto,
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
					ObjFreshThreshold: config.MustParseDurationSlice("1m"),
					Paging:            config.PagingDefault,
					Theme:             *testconfig.DarkTheme,
					Preset:            config.PresetAuto,
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
					Preset:  config.PresetAuto,
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
					Preset:  config.PresetAuto,
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
					Preset:  config.PresetAuto,
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
					Preset:  config.PresetAuto,
				},
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "Auto dark theme",
			args: []string{"get", "pods"},
			env: map[string]string{
				"KUBECOLOR_PRESET": string(config.PresetAuto),
			},
			term: testconfig.DarkTerm,
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl: "kubectl",
					Paging:  config.PagingNever,
					Theme:   *testconfig.DarkTheme,
					Preset:  config.PresetAuto,
				},
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "Auto light theme",
			args: []string{"get", "pods"},
			env: map[string]string{
				"KUBECOLOR_PRESET": string(config.PresetAuto),
			},
			term: testconfig.LightTerm,
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl: "kubectl",
					Paging:  config.PagingNever,
					Theme:   *testconfig.LightTheme,
					Preset:  config.PresetAuto,
				},
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range test.env {
				testutil.Setenv(t, k, v)
			}

			conf, err := ResolveConfig(test.args, test.term)
			testutil.MustNoError(t, err)

			// Don't test flags field
			test.expectedConf.Flags = conf.Flags

			testutil.MustEqual(t, test.expectedConf, conf)
		})
	}
}
