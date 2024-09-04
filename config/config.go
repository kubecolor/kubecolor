package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/internal/slogutil"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// PresetKey is the Viper config key to use in [viper.Viper.Set].
const PresetKey = "preset"

type Config struct {
	Debug             bool          `jsonschema:"-"`
	Kubectl           string        `jsonschema:"default=kubectl,example=kubectl1.19,example=oc"` // Which kubectl executable to use
	ObjFreshThreshold time.Duration // Ages below this uses theme.data.durationfresh coloring
	Preset            Preset        // Color theme preset
	Theme             Theme         //
	Pager             string        `jsonschema:"example=less -RF,less --RAW-CONTROL-CHARS --quit-if-one-screen,example=more"` // Command to use as pager
	Paging            Paging        `jsonschema:"default=never"`                                                               // Whether to enable paging: "auto" or "never"
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
	// NOTE: Don't bind PAGER here as it should be overwritten by the config file

	v.SetDefault("kubectl", "kubectl")
	// mapstructure doesn't like "type X string" values, so we have to convert it via string(...)
	v.SetDefault(PresetKey, string(PresetDefault))
	v.SetDefault("paging", string(PagingDefault))
	v.SetDefault("pager", defaultPager())

	return v
}

func LoadViper() (*viper.Viper, error) {
	v := NewViper()

	if v.GetBool("debug") {
		if logger, ok := slog.Default().Handler().(*slogutil.SlogHandler); ok {
			logger.Level = slog.LevelDebug
		}
	}

	if path := os.Getenv("KUBECOLOR_CONFIG"); path != "" {
		v.SetConfigFile(path)
		slog.Debug("Overriding config path with environment variable", "KUBECOLOR_CONFIG", path)
	} else if homeDir, err := os.UserHomeDir(); err == nil {
		// ~/.kube/color.yaml
		v.AddConfigPath(filepath.Join(homeDir, ".kube"))
	}

	if err := v.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) || os.IsNotExist(err) {
			slog.Debug("No config file found. " + err.Error())
			// continue
		} else {
			return nil, err
		}
	} else {
		if fileUsed := v.ConfigFileUsed(); fileUsed != "" {
			slog.Debug("Using config", "file", fileUsed)
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
	slog.Debug("Applying theme", "preset", preset)
	theme := NewBaseTheme(preset)
	applyViperDefaults(theme, v)
	return nil
}

func defaultPager() string {
	if p := os.Getenv("PAGER"); p != "" {
		return p
	}
	if _, err := exec.LookPath("less"); err == nil {
		return "less -RF"
	}
	if _, err := exec.LookPath("more"); err == nil {
		return "more"
	}
	return ""
}
