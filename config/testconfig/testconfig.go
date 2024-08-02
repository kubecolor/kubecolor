package testconfig

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/kubecolor/kubecolor/config"
)

var (
	DarkTheme  *config.Theme
	LightTheme *config.Theme
)

func init() {
	os.Clearenv()
	color.ForceColor()
	color.Enable = true

	DarkTheme = NewTheme(config.PresetDark)
	LightTheme = NewTheme(config.PresetLight)
}

// NewTheme returns a theme from a preset that's meant to be used in testing.
func NewTheme(preset config.Preset) *config.Theme {
	v := config.NewViper()
	// mapstructure doesn't like "type X string" values, so we have to convert it via string(...)
	v.Set(config.PresetKey, string(preset))
	cfg, err := config.Unmarshal(v)
	if err != nil {
		panic(fmt.Errorf("unmarshal config: %w", err))
	}
	return &cfg.Theme
}
