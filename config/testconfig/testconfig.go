package testconfig

import (
	"fmt"

	"github.com/kubecolor/kubecolor/config"
)

var DarkTheme *config.Theme = NewTheme(config.PresetDark)
var LightTheme *config.Theme = NewTheme(config.PresetLight)

// NewTheme returns a theme from a preset that's meant to be used in testing.
func NewTheme(preset config.Preset) *config.Theme {
	v := config.NewViper()
	v.Set(config.PresetKey, preset)
	cfg, err := config.Unmarshal(v)
	if err != nil {
		panic(fmt.Errorf("unmarshal config: %w", err))
	}
	return &cfg.Theme
}
