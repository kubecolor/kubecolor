package command

import (
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/testutil"
)

func Test_ResolveConfig(t *testing.T) {
	tests := []struct {
		name                 string
		args                 []string
		kubectlCommand       string
		objFreshAgeThreshold string
		expectedArgs         []string
		expectedConf         *KubecolorConfig
	}{
		{
			name:         "no config",
			args:         []string{"get", "pods"},
			expectedArgs: []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:             false,
				DarkBackground:    true,
				ForceColor:        false,
				KubectlCmd:        "kubectl",
				ObjFreshThreshold: time.Duration(0),
			},
		},
		{
			name:         "plain, dark, force",
			args:         []string{"get", "pods", "--plain", "--light-background", "--force-colors"},
			expectedArgs: []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:             true,
				DarkBackground:    false,
				ForceColor:        true,
				KubectlCmd:        "kubectl",
				ObjFreshThreshold: time.Duration(0),
			},
		},
		{
			name:           "KUBECTL_COMMAND exists",
			args:           []string{"get", "pods", "--plain"},
			kubectlCommand: "kubectl.1.19",
			expectedArgs:   []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:             true,
				DarkBackground:    true,
				ForceColor:        false,
				KubectlCmd:        "kubectl.1.19",
				ObjFreshThreshold: time.Duration(0),
			},
		},
		{
			name:                 "KUBECOLOR_OBJ_FRESH exists",
			args:                 []string{"get", "pods"},
			expectedArgs:         []string{"get", "pods"},
			objFreshAgeThreshold: "1m",
			expectedConf: &KubecolorConfig{
				Plain:             false,
				DarkBackground:    true,
				ForceColor:        false,
				KubectlCmd:        "kubectl",
				ObjFreshThreshold: time.Minute,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			testutil.Setenv(t, "KUBECTL_COMMAND", tt.kubectlCommand)
			testutil.Setenv(t, "KUBECOLOR_OBJ_FRESH", tt.objFreshAgeThreshold)

			args, conf := ResolveConfig(tt.args)
			testutil.MustEqual(t, tt.expectedArgs, args)
			testutil.MustEqual(t, tt.expectedConf, conf)
		})
	}
}
