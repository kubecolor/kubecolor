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

const (
	PagingAlways Paging = "always"
	PagingNever  Paging = "never"
)

// PagingOptions can be used for validation and schema generation (not implemented)
var (
	PagingOptions []Paging = []Paging{
		PagingAlways,
		PagingNever,
	}
)

type Paging string

type Config struct {
	Debug             bool          `jsonschema:"-"`
	Kubectl           string        `jsonschema:"default=kubectl,example=kubectl1.19,example=oc"` // Which kubectl executable to use
	ObjFreshThreshold time.Duration // Ages below this uses theme.data.durationfresh coloring
	Preset            Preset        // Color theme preset
	Theme             Theme         //
	Pager             string        // Command to use as pager
	Paging            Paging        `jsonschema:"default=always"` // Whether to enable paging: "always" or "never"
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
	v.MustBindEnv("pager", "KUBECOLOR_PAGER") // Don't bind PAGER yet as it should be overwritten by the config file

	v.SetDefault("kubectl", "kubectl")
	// mapstructure doesn't like "type X string" values, so we have to convert it via string(...)
	v.SetDefault(PresetKey, string(PresetDefault))
	v.SetDefault("paging", "always")

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
	theme := NewBaseTheme(preset)
	applyViperDefaults(theme, v)
	return nil
}
