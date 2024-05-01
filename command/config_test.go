package command

import (
	"os"
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

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
				Plain:             false,
				ForceColor:        false,
				KubectlCmd:        "kubectl",
				ObjFreshThreshold: time.Duration(0),
				Theme:             testconfig.DarkTheme,
				ArgsPassthrough:   []string{"get", "pods"},
			},
		},
		{
			name: "plain, light, force",
			args: []string{"get", "pods", "--plain", "--light-background", "--force-colors"},
			expectedConf: &Config{
				Plain:             true,
				ForceColor:        true,
				KubectlCmd:        "kubectl",
				ObjFreshThreshold: time.Duration(0),
				Theme:             testconfig.LightTheme,
				ArgsPassthrough:   []string{"get", "pods"},
			},
		},
		{
			name: "KUBECTL_COMMAND exists",
			args: []string{"get", "pods", "--plain"},
			env:  map[string]string{"KUBECTL_COMMAND": "kubectl.1.19"},
			expectedConf: &Config{
				Plain:             true,
				ForceColor:        false,
				KubectlCmd:        "kubectl.1.19",
				ObjFreshThreshold: time.Duration(0),
				Theme:             testconfig.DarkTheme,
				ArgsPassthrough:   []string{"get", "pods"},
			},
		},
		{
			name: "KUBECOLOR_OBJ_FRESH exists",
			args: []string{"get", "pods"},
			env:  map[string]string{"KUBECOLOR_OBJ_FRESH": "1m"},
			expectedConf: &Config{
				Plain:             false,
				ForceColor:        false,
				KubectlCmd:        "kubectl",
				ObjFreshThreshold: time.Minute,
				Theme:             testconfig.DarkTheme,
				ArgsPassthrough:   []string{"get", "pods"},
			},
		},
		{
			name: "KUBECOLOR_LIGHT_BACKGROUND via env",
			args: []string{"get", "pods"},
			env:  map[string]string{"KUBECOLOR_LIGHT_BACKGROUND": "true"},
			expectedConf: &Config{
				Plain:           false,
				ForceColor:      false,
				KubectlCmd:      "kubectl",
				Theme:           testconfig.LightTheme,
				ArgsPassthrough: []string{"get", "pods"},
			},
		},
		{
			name: "KUBECOLOR_FORCE_COLORS env var",
			args: []string{"get", "pods"},
			env:  map[string]string{"KUBECOLOR_FORCE_COLORS": "true"},
			expectedConf: &Config{
				Plain:           false,
				ForceColor:      true,
				KubectlCmd:      "kubectl",
				Theme:           testconfig.DarkTheme,
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
			testutil.MustEqual(t, tt.expectedConf, conf)
		})
	}
}
