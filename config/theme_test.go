package config_test

import (
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
)

func TestResolveAutoTheme(t *testing.T) {
	t.Parallel()
	tests := []struct {
		preset config.Preset
		want   config.Preset
		term   config.TermInfo
	}{
		{preset: config.PresetAuto, want: config.PresetLight, term: testconfig.LightTerm},
		{preset: config.PresetAuto, want: config.PresetDark, term: testconfig.DarkTerm},
		{preset: config.PresetAuto, want: config.PresetDark, term: nil},
		{preset: config.PresetProtAuto, want: config.PresetProtLight, term: testconfig.LightTerm},
		{preset: config.PresetProtAuto, want: config.PresetProtDark, term: testconfig.DarkTerm},
		{preset: config.PresetProtAuto, want: config.PresetProtDark, term: nil},

		{preset: config.PresetDark, want: config.PresetDark, term: testconfig.LightTerm},
		{preset: config.PresetDark, want: config.PresetDark, term: testconfig.DarkTerm},
		{preset: config.PresetDark, want: config.PresetDark, term: nil},
		{preset: config.PresetLight, want: config.PresetLight, term: testconfig.LightTerm},
		{preset: config.PresetLight, want: config.PresetLight, term: testconfig.DarkTerm},
		{preset: config.PresetLight, want: config.PresetLight, term: nil},
	}

	for _, test := range tests {
		nameSuffix := ""
		switch {
		case test.term == nil:
			nameSuffix = "/nil term"
		case test.term.HasDarkBackground():
			nameSuffix = "/dark term"
		default:
			nameSuffix = "/light term"
		}
		t.Run(test.preset.String()+nameSuffix, func(t *testing.T) {
			got := config.ResolveAutoThemePreset(test.preset, test.term)
			if test.want != got {
				t.Errorf("want %q, got %q", test.want, got)
			}
		})
	}
}
