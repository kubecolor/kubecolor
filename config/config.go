package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// PresetKey is the Viper config key to use in [viper.Viper.Set].
const PresetKey = "preset"

type Config struct {
	Debug             bool
	Kubectl           string
	ObjFreshThreshold time.Duration
	Preset            Preset
	Theme             Theme
}

func NewViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("color")
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvPrefix("KUBECOLOR")
	v.SetEnvKeyReplacer(strings.NewReplacer(
		".", "_",
	))

	v.MustBindEnv("kubectl", "KUBECTL_COMMAND")
	v.MustBindEnv("objfreshthreshold", "KUBECOLOR_OBJ_FRESH")

	v.SetDefault("kubectl", "kubectl")
	v.SetDefault(PresetKey, "dark")

	return v
}

func LoadViper() (*viper.Viper, error) {
	v := NewViper()

	if path := os.Getenv("KUBECOLOR_CONFIG"); path != "" {
		v.AddConfigPath(path)
	}
	if homeDir, err := os.UserHomeDir(); err == nil {
		// ~/.kube/color.yaml
		v.AddConfigPath(filepath.Join(homeDir, ".kube"))
	}

	if err := v.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			// continue
		} else {
			return nil, err
		}
	}

	if v.GetBool("debug") {
		if fileUsed := v.ConfigFileUsed(); fileUsed != "" {
			fmt.Fprintf(os.Stderr, "[kubecolor] [debug] using config: %s\n", fileUsed)
		}
	}

	return v, nil
}

func Unmarshal(v *viper.Viper) (*Config, error) {
	if err := ApplyThemePreset(v); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.TextUnmarshallerHookFunc(),
		))); err != nil {
		return nil, err
	}
	return cfg, nil
}

func ApplyThemePreset(v *viper.Viper) error {
	preset, err := ParsePreset(v.GetString(PresetKey))
	if err != nil {
		return fmt.Errorf("parse preset: %w", err)
	}
	if v.GetBool("debug") {
		fmt.Fprintf(os.Stderr, "[kubecolor] [debug] applying preset: %s\n", preset)
	}
	v.Set(PresetKey, preset) // to skip parsing it twice
	theme := NewBaseTheme(preset)
	applyViperDefaults(theme, v)
	return nil
}
