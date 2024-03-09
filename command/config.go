package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/config"
)

type Config struct {
	Plain                bool
	ForceColor           bool
	ShowKubecolorVersion bool
	KubectlCmd           string
	ObjFreshThreshold    time.Duration
	Theme                *config.Theme

	ArgsPassthrough []string
}

func ResolveConfig(inputArgs []string) (*Config, error) {
	cfg := &Config{}

	v, err := config.LoadViper()
	if err != nil {
		return nil, err
	}

	if lightThemeEnv, ok, err := parseBoolEnv("KUBECOLOR_LIGHT_BACKGROUND"); err != nil {
		return nil, err
	} else if ok {
		if lightThemeEnv {
			v.Set("preset", "dark")
		} else {
			v.Set("preset", "light")
		}
	}

	for _, s := range inputArgs {
		flag, value, _ := strings.Cut(s, "=")
		switch flag {
		case "--plain":
			if b, err := parseBoolFlag(flag, value); err != nil {
				return nil, err
			} else {
				cfg.Plain = b
			}
		case "--light-background":
			if b, err := parseBoolFlag(flag, value); err != nil {
				return nil, err
			} else {
				if b {
					v.Set("preset", "light")
				} else {
					v.Set("preset", "dark")
				}
			}
		case "--force-colors":
			if b, err := parseBoolFlag(flag, value); err != nil {
				return nil, err
			} else {
				cfg.ForceColor = b
			}
		case "--kubecolor-version":
			if b, err := parseBoolFlag(flag, value); err != nil {
				return nil, err
			} else {
				cfg.ShowKubecolorVersion = b
			}
		case "--kubecolor-theme":
			v.Set("preset", value)
		default:
			cfg.ArgsPassthrough = append(cfg.ArgsPassthrough, s)
		}
	}

	if err := config.ApplyThemePreset(v); err != nil {
		return nil, err
	}

	newCfg, err := config.Unmarshal(v)
	if err != nil {
		return nil, err
	}
	cfg.KubectlCmd = newCfg.Kubectl
	cfg.ObjFreshThreshold = newCfg.ObjFreshThreshold
	cfg.Theme = &newCfg.Theme

	if b, ok, err := parseBoolEnv("KUBECOLOR_FORCE_COLORS"); err != nil {
		return nil, err
	} else if ok {
		cfg.ForceColor = b
	}

	return cfg, nil
}

func parseBool(value string) (result bool, ok bool, err error) {
	if value == "" {
		return false, false, nil
	}
	ok = true
	switch strings.ToLower(value) {
	case "true":
		result = true
	case "false":
		result = false
	default:
		return false, false, fmt.Errorf(`must be either "true" or "false"`)
	}
	return result, ok, err
}

func parseBoolFlag(flag, value string) (result bool, err error) {
	result, ok, err := parseBool(value)
	if err != nil {
		return false, fmt.Errorf("flag %s: %w", flag, err)
	}
	if !ok {
		// bool flags treat no value as true (e.g "--plain" is same as "--plain=true")
		return true, nil
	}
	return result, nil
}

func parseBoolEnv(env string) (result bool, ok bool, err error) {
	result, ok, err = parseBool(os.Getenv(env))
	if err != nil {
		return false, false, fmt.Errorf("parse env %s: %w", env, err)
	}
	return result, ok, err
}

func parseDuration(value string) (time.Duration, bool, error) {
	if value == "" {
		return 0, false, nil
	}
	result, err := time.ParseDuration(value)
	if err != nil {
		return 0, false, err
	}
	return result, true, nil
}

func parseDurationEnv(env string) (time.Duration, bool, error) {
	result, ok, err := parseDuration(os.Getenv(env))
	if err != nil {
		return 0, false, fmt.Errorf("parse env %s: %w", env, err)
	}
	return result, ok, nil
}
