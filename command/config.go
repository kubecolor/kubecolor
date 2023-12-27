package command

import (
	"fmt"
	"os"
	"time"
)

type KubecolorConfig struct {
	Plain                bool
	DarkBackground       bool
	ForceColor           bool
	ShowKubecolorVersion bool
	KubectlCmd           string
	ObjFreshThreshold    time.Duration
}

func ResolveConfig(args []string) ([]string, *KubecolorConfig) {
	args, plainFlagFound := findAndRemoveBoolFlagIfExists(args, "--plain")
	args, lightBackgroundFlagFound := findAndRemoveBoolFlagIfExists(args, "--light-background")
	args, forceColorFlagFound := findAndRemoveBoolFlagIfExists(args, "--force-colors")
	args, kubecolorVersionFlagFound := findAndRemoveBoolFlagIfExists(args, "--kubecolor-version")

	colorsForcedViaEnv := os.Getenv("KUBECOLOR_FORCE_COLORS") == "true"
	lightBackgroundViaEnv := os.Getenv("KUBECOLOR_LIGHT_BACKGROUND") == "true"

	kubectlCmd := "kubectl"
	if kc := os.Getenv("KUBECTL_COMMAND"); kc != "" {
		kubectlCmd = kc
	}

	objFreshAgeThresholdDuration := time.Duration(0)
	objFreshAgeThresholdEnv := "KUBECOLOR_OBJ_FRESH"
	if objFreshAgeThreshold := os.Getenv(objFreshAgeThresholdEnv); objFreshAgeThreshold != "" {
		var err error
		objFreshAgeThresholdDuration, err = time.ParseDuration(objFreshAgeThreshold)
		if err != nil {
			fmt.Printf("[WARN] [kubecolor] cannot parse duration taken from env %s. See kubecolor document. %v\n", objFreshAgeThresholdEnv, err)
		}
	}

	return args, &KubecolorConfig{
		Plain:                plainFlagFound,
		DarkBackground:       !lightBackgroundFlagFound && !lightBackgroundViaEnv,
		ForceColor:           forceColorFlagFound || colorsForcedViaEnv,
		ShowKubecolorVersion: kubecolorVersionFlagFound,
		KubectlCmd:           kubectlCmd,
		ObjFreshThreshold:    objFreshAgeThresholdDuration,
	}
}

func findAndRemoveBoolFlagIfExists(args []string, key string) ([]string, bool) {
	for i, arg := range args {
		if arg == key {
			return append(args[:i], args[i+1:]...), true
		}
	}

	return args, false
}
