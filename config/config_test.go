package config

import (
	"os"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestEnvVars_preset(t *testing.T) {
	os.Clearenv()
	testutil.Setenv(t, "KUBECOLOR_PRESET", "light")

	v := NewViper()
	cfg, err := Unmarshal(v)
	testutil.MustNoError(t, err)

	testutil.Equal(t, PresetLight, cfg.Preset, "Read from cfg.Preset")
	testutil.Equal(t, PresetLight, v.Get(PresetKey), "Read from v.Get(...)")
}

func TestEnvVars_theme(t *testing.T) {
	os.Clearenv()
	testutil.Setenv(t, "KUBECOLOR_THEME_TABLE_HEADER", "red")

	v := NewViper()
	cfg, err := Unmarshal(v)
	testutil.MustNoError(t, err)

	testutil.Equal(t, Color{Source: "red", Code: "31"}, cfg.Theme.Table.Header, "Read from cfg.Theme.Table.Header")
	testutil.Equal(t, "red", v.Get("theme.table.header"), "Read from v.Get(...)")
}
