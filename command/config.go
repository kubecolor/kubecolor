package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/config"
	"github.com/spf13/viper"
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
	v, err := config.LoadViper()
	if err != nil {
		return nil, err
	}
	return ResolveConfigViper(inputArgs, v)
}

func ResolveConfigViper(inputArgs []string, v *viper.Viper) (*Config, error) {
	cfg := &Config{}

	if lightThemeEnv, ok, err := parseBoolEnv("KUBECOLOR_LIGHT_BACKGROUND"); err != nil {
		return nil, err
	} else if ok {
		if lightThemeEnv {
			v.Set(config.PresetKey, "dark")
		} else {
			v.Set(config.PresetKey, "light")
		}
	}

	if b, ok, err := parseBoolEnv("KUBECOLOR_FORCE_COLORS"); err != nil {
		return nil, err
	} else if ok {
		cfg.ForceColor = b
	}

	for _, s := range inputArgs {
		flag, value, _ := strings.Cut(s, "=")
		switch flag {
		case "--plain":
			b, err := parseBoolFlag(flag, value)
			if err != nil {
				return nil, err
			}
			cfg.Plain = b
		case "--light-background":
			b, err := parseBoolFlag(flag, value)
			if err != nil {
				return nil, err
			}
			if b {
				v.Set(config.PresetKey, "light")
			} else {
				v.Set(config.PresetKey, "dark")
			}
		case "--force-colors":
			b, err := parseBoolFlag(flag, value)
			if err != nil {
				return nil, err
			}
			cfg.ForceColor = b
		case "--kubecolor-version":
			b, err := parseBoolFlag(flag, value)
			if err != nil {
				return nil, err
			}
			cfg.ShowKubecolorVersion = b
		case "--kubecolor-theme":
			v.Set(config.PresetKey, value)
		default:
			cfg.ArgsPassthrough = append(cfg.ArgsPassthrough, s)
		}
	}

	newCfg, err := config.Unmarshal(v)
	if err != nil {
		return nil, err
	}
	cfg.KubectlCmd = newCfg.Kubectl
	cfg.ObjFreshThreshold = newCfg.ObjFreshThreshold
	cfg.Theme = &newCfg.Theme

	return cfg, nil
}

func parseBool(value string) (result bool, ok bool, err error) {
	switch strings.ToLower(value) {
	case "":
		return false, false, nil
	case "true":
		return true, true, nil
	case "false":
		return false, true, nil
	default:
		return false, false, fmt.Errorf(`must be either "true" or "false"`)
	}
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
