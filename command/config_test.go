package command

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

var defaultDuration = config.Duration{
	Threshold1: 5 * time.Minute,
	Threshold2: 2 * time.Hour,
	Threshold3: 24 * time.Hour,
	Threshold4: 30 * 24 * time.Hour,
	Threshold5: 365 * 24 * time.Hour,
}

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
					Kubectl:  "kubectl",
					Duration: defaultDuration,
					Paging:   config.PagingDefault,
					Theme:    *testconfig.DarkTheme,
					Preset:   config.PresetDark,
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
					Kubectl:  "kubectl",
					Duration: defaultDuration,
					Paging:   config.PagingDefault,
					Theme:    *testconfig.LightTheme,
					Preset:   config.PresetLight,
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
					Kubectl:  "kubectl.1.19",
					Duration: defaultDuration,
					Paging:   config.PagingDefault,
					Theme:    *testconfig.DarkTheme,
					Preset:   config.PresetDark,
				},
				ForceColor:      ColorLevelNone,
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "KUBECOLOR_LIGHT_BACKGROUND via env",
			args: []string{"get", "pods"},
			env:  map[string]string{"KUBECOLOR_LIGHT_BACKGROUND": "true"},
			expectedConf: &Config{
				Config: &config.Config{
					Kubectl:  "kubectl",
					Duration: defaultDuration,
					Paging:   config.PagingDefault,
					Theme:    *testconfig.LightTheme,
					Preset:   config.PresetLight,
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
					Kubectl:  "kubectl",
					Duration: defaultDuration,
					Paging:   config.PagingDefault,
					Theme:    *testconfig.DarkTheme,
					Preset:   config.PresetDark,
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
					Kubectl:  "kubectl",
					Duration: defaultDuration,
					Paging:   config.PagingDefault,
					Theme:    *testconfig.DarkTheme,
					Preset:   config.PresetDark,
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
					Kubectl:  "kubectl",
					Duration: defaultDuration,
					Pager:    "most",
					Paging:   config.PagingAuto,
					Theme:    *testconfig.DarkTheme,
					Preset:   config.PresetDark,
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
					Kubectl:  "kubectl",
					Duration: defaultDuration,
					Paging:   config.PagingNever,
					Theme:    *testconfig.DarkTheme,
					Preset:   config.PresetDark,
				},
				ArgsPassthrough: []string{"get", "pods"},
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

func Test_ResolveConfig_LegacyDuration_ConfigFile(t *testing.T) {
	writeConfig := func(t *testing.T, content string) string {
		t.Helper()
		f := filepath.Join(t.TempDir(), "color.yaml")
		if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return f
	}

	t.Run("legacy objFreshThreshold in config file", func(t *testing.T) {
		os.Clearenv()
		testutil.Setenv(t, "KUBECOLOR_CONFIG", writeConfig(t, `
objFreshThreshold: 2m
`))

		conf, err := ResolveConfig([]string{"get", "pods"})
		testutil.MustNoError(t, err)

		testutil.Equal(t, 2*time.Minute, conf.Duration.Threshold1, "threshold1")
		testutil.Equal(t, time.Duration(0), conf.Duration.Threshold2, "threshold2")
		testutil.Equal(t, "green", conf.Theme.Data.Duration.Default.Source, "fresh color")
		testutil.Equal(t, "", conf.Theme.Data.Duration.Threshold1.Source, "non-fresh color")
	})

	t.Run("legacy objFreshThreshold + durationFresh in config file", func(t *testing.T) {
		os.Clearenv()
		testutil.Setenv(t, "KUBECOLOR_CONFIG", writeConfig(t, `
objFreshThreshold: 3m
theme:
  data:
    durationFresh: cyan
`))

		conf, err := ResolveConfig([]string{"get", "pods"})
		testutil.MustNoError(t, err)

		testutil.Equal(t, 3*time.Minute, conf.Duration.Threshold1, "threshold1")
		testutil.Equal(t, "cyan", conf.Theme.Data.Duration.Default.Source, "fresh color")
		testutil.Equal(t, "", conf.Theme.Data.Duration.Threshold1.Source, "non-fresh color")
	})

	t.Run("legacy objFreshThreshold + durationFresh + duration in config file", func(t *testing.T) {
		os.Clearenv()
		testutil.Setenv(t, "KUBECOLOR_CONFIG", writeConfig(t, `
objFreshThreshold: 3m
theme:
  data:
    durationFresh: cyan
    duration: yellow
`))

		conf, err := ResolveConfig([]string{"get", "pods"})
		testutil.MustNoError(t, err)

		testutil.Equal(t, 3*time.Minute, conf.Duration.Threshold1, "threshold1")
		testutil.Equal(t, "cyan", conf.Theme.Data.Duration.Default.Source, "fresh color")
		testutil.Equal(t, "yellow", conf.Theme.Data.Duration.Threshold1.Source, "non-fresh color")
	})

	t.Run("new duration thresholds and colors in config file", func(t *testing.T) {
		os.Clearenv()
		testutil.Setenv(t, "KUBECOLOR_CONFIG", writeConfig(t, `
duration:
  threshold1: 10m
  threshold2: 1h
  threshold3: 7d
  threshold4: 90d
  threshold5: 365d
theme:
  data:
    duration:
      default: "hi-green"
      threshold1: "green"
      threshold2: "hi-yellow"
      threshold3: "yellow"
      threshold4: "hi-red"
      threshold5: "red"
`))

		conf, err := ResolveConfig([]string{"get", "pods"})
		testutil.MustNoError(t, err)

		testutil.Equal(t, 10*time.Minute, conf.Duration.Threshold1, "threshold1")
		testutil.Equal(t, time.Hour, conf.Duration.Threshold2, "threshold2")
		testutil.Equal(t, 7*24*time.Hour, conf.Duration.Threshold3, "threshold3")
		testutil.Equal(t, 90*24*time.Hour, conf.Duration.Threshold4, "threshold4")
		testutil.Equal(t, 365*24*time.Hour, conf.Duration.Threshold5, "threshold5")

		colors := conf.Theme.Data.Duration
		testutil.Equal(t, "hi-green", colors.Default.Source, "duration.default")
		testutil.Equal(t, "green", colors.Threshold1.Source, "duration.threshold1")
		testutil.Equal(t, "hi-yellow", colors.Threshold2.Source, "duration.threshold2")
		testutil.Equal(t, "yellow", colors.Threshold3.Source, "duration.threshold3")
		testutil.Equal(t, "hi-red", colors.Threshold4.Source, "duration.threshold4")
		testutil.Equal(t, "red", colors.Threshold5.Source, "duration.threshold5")
		testutil.Equal(t, "", colors.Threshold6.Source, "duration.threshold6")
	})
}

func Test_ResolveConfig_LegacyDuration_Env(t *testing.T) {
	type wantColors struct {
		Default    string
		Threshold1 string
		Threshold2 string
		Threshold3 string
		Threshold4 string
		Threshold5 string
		Threshold6 string
	}

	// Dark preset base colors used in defaultFrom tags:
	//   Default    → theme.base.primary   = "magenta"
	//   Threshold1 → theme.base.secondary = "cyan"
	//   Threshold2 → theme.base.info      = "white"
	//   Threshold3 → theme.base.success   = "green"
	//   Threshold4 → theme.base.warning   = "yellow"
	//   Threshold5 → theme.base.danger    = "red"
	//   Threshold6 → (no default)         = ""
	darkDefaults := wantColors{
		Default:    "magenta",
		Threshold1: "cyan",
		Threshold2: "white",
		Threshold3: "green",
		Threshold4: "yellow",
		Threshold5: "red",
		Threshold6: "",
	}

	tests := []struct {
		name         string
		env          map[string]string
		wantDuration config.Duration
		wantColors   wantColors
	}{
		{
			name: "defaults",
			wantDuration: config.Duration{
				Threshold1: 5 * time.Minute,
				Threshold2: 2 * time.Hour,
				Threshold3: 24 * time.Hour,
				Threshold4: 30 * 24 * time.Hour,
				Threshold5: 365 * 24 * time.Hour,
			},
			wantColors: darkDefaults,
		},
		{
			name: "custom thresholds via env",
			env: map[string]string{
				"KUBECOLOR_DURATION_THRESHOLD1": "10m",
				"KUBECOLOR_DURATION_THRESHOLD2": "6h",
				"KUBECOLOR_DURATION_THRESHOLD3": "7d",
				"KUBECOLOR_DURATION_THRESHOLD4": "60d",
				"KUBECOLOR_DURATION_THRESHOLD5": "180d",
			},
			wantDuration: config.Duration{
				Threshold1: 10 * time.Minute,
				Threshold2: 6 * time.Hour,
				Threshold3: 7 * 24 * time.Hour,
				Threshold4: 60 * 24 * time.Hour,
				Threshold5: 180 * 24 * time.Hour,
			},
			wantColors: darkDefaults,
		},
		{
			name: "legacy KUBECOLOR_OBJ_FRESH sets threshold1 and zeroes rest",
			env: map[string]string{
				"KUBECOLOR_OBJ_FRESH": "1m",
			},
			wantDuration: config.Duration{
				Threshold1: time.Minute,
			},
			wantColors: wantColors{
				Default:    "green", // legacy DurationFresh defaulted from base.success
				Threshold1: "",      // originally Duration had no default
			},
		},
		{
			name: "legacy KUBECOLOR_OBJ_FRESH + KUBECOLOR_THEME_DATA_DURATIONFRESH",
			env: map[string]string{
				"KUBECOLOR_OBJ_FRESH":                "2m",
				"KUBECOLOR_THEME_DATA_DURATIONFRESH": "yellow",
			},
			wantDuration: config.Duration{
				Threshold1: 2 * time.Minute,
			},
			wantColors: wantColors{
				Default:    "yellow",
				Threshold1: "", // originally Duration had no default
			},
		},
		{
			name: "legacy KUBECOLOR_OBJ_FRESH + KUBECOLOR_THEME_DATA_DURATION",
			env: map[string]string{
				"KUBECOLOR_OBJ_FRESH":           "3m",
				"KUBECOLOR_THEME_DATA_DURATION": "red",
			},
			wantDuration: config.Duration{
				Threshold1: 3 * time.Minute,
			},
			wantColors: wantColors{
				Default:    "green",
				Threshold1: "red",
			},
		},
		{
			name: "legacy KUBECOLOR_OBJ_FRESH + KUBECOLOR_THEME_DATA_DURATIONFRESH + old KUBECOLOR_THEME_DATA_DURATION",
			env: map[string]string{
				"KUBECOLOR_OBJ_FRESH":                "3m",
				"KUBECOLOR_THEME_DATA_DURATIONFRESH": "yellow",
				"KUBECOLOR_THEME_DATA_DURATION":      "red",
			},
			wantDuration: config.Duration{
				Threshold1: 3 * time.Minute,
			},
			wantColors: wantColors{
				Default:    "yellow",
				Threshold1: "red",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				testutil.Setenv(t, k, v)
			}

			conf, err := ResolveConfig([]string{"get", "pods"})
			testutil.MustNoError(t, err)

			testutil.Equal(t, tt.wantDuration, conf.Duration, "duration thresholds")

			colors := conf.Theme.Data.Duration
			testutil.Equal(t, tt.wantColors.Default, colors.Default.Source, "duration.default")
			testutil.Equal(t, tt.wantColors.Threshold1, colors.Threshold1.Source, "duration.threshold1")
			testutil.Equal(t, tt.wantColors.Threshold2, colors.Threshold2.Source, "duration.threshold2")
			testutil.Equal(t, tt.wantColors.Threshold3, colors.Threshold3.Source, "duration.threshold3")
			testutil.Equal(t, tt.wantColors.Threshold4, colors.Threshold4.Source, "duration.threshold4")
			testutil.Equal(t, tt.wantColors.Threshold5, colors.Threshold5.Source, "duration.threshold5")
			testutil.Equal(t, tt.wantColors.Threshold6, colors.Threshold6.Source, "duration.threshold6")
		})
	}
}
