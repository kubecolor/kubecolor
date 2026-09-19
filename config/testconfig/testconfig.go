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
	NullTheme  *config.Theme

	DarkTerm  = NewTermInfoWithDarkBackground(true)
	LightTerm = NewTermInfoWithDarkBackground(false)
)

func init() {
	os.Clearenv()
	color.ForceColor()
	color.Enable = true

	DarkTheme = NewTheme(config.PresetDark, nil)
	LightTheme = NewTheme(config.PresetLight, nil)
	NullTheme = &config.Theme{}
}

// NewTheme returns a theme from a preset that's meant to be used in testing.
func NewTheme(preset config.Preset, term config.TermInfo) *config.Theme {
	v := config.NewViper()
	// mapstructure doesn't like "type X string" values, so we have to convert it via string(...)
	v.Set(config.PresetKey, string(preset))
	cfg, err := config.Unmarshal(v, term)
	if err != nil {
		panic(fmt.Errorf("unmarshal config: %w", err))
	}
	return &cfg.Theme
}

func NewTermInfoWithDarkBackground(hasDark bool) config.TermInfo {
	return terminfo{hasDarkBackground: hasDark}
}

type terminfo struct {
	hasDarkBackground bool
}

var _ config.TermInfo = terminfo{}

// HasDarkBackground implements [config.TermInfo].
func (t terminfo) HasDarkBackground() bool {
	return t.hasDarkBackground
}
