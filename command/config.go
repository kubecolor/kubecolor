package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/color"
)

type Config struct {
	Plain                bool
	ForceColor           bool
	ShowKubecolorVersion bool
	KubectlCmd           string
	ObjFreshThreshold    time.Duration
	Theme                *color.Theme

	ArgsPassthrough []string
}

func ResolveConfig(inputArgs []string) (*Config, error) {
	config := &Config{
		KubectlCmd:        "kubectl",
		ObjFreshThreshold: time.Duration(0),
	}

	themePreset := color.PresetDefault

	if themePresetEnv, ok, err := parseThemePresetEnv("KUBECOLOR_THEME"); err != nil {
		return nil, err
	} else if ok {
		themePreset = themePresetEnv
	}

	if lightThemeEnv, ok, err := parseBoolEnv("KUBECOLOR_LIGHT_BACKGROUND"); err != nil {
		return nil, err
	} else if ok {
		if lightThemeEnv {
			themePreset = color.PresetDark
		} else {
			themePreset = color.PresetLight
		}
	}

	for _, s := range inputArgs {
		flag, value, _ := strings.Cut(s, "=")
		switch flag {
		case "--plain":
			if b, err := parseBoolFlag(flag, value); err != nil {
				return nil, err
			} else {
				config.Plain = b
			}
		case "--light-background":
			if b, err := parseBoolFlag(flag, value); err != nil {
				return nil, err
			} else {
				if b {
					themePreset = color.PresetLight
				} else {
					themePreset = color.PresetDark
				}
			}
		case "--force-colors":
			if b, err := parseBoolFlag(flag, value); err != nil {
				return nil, err
			} else {
				config.ForceColor = b
			}
		case "--kubecolor-version":
			if b, err := parseBoolFlag(flag, value); err != nil {
				return nil, err
			} else {
				config.ShowKubecolorVersion = b
			}
		case "--kubecolor-theme":
			if result, err := parseThemePresetFlag(flag, value); err != nil {
				return nil, err
			} else {
				themePreset = result
			}
		default:
			config.ArgsPassthrough = append(config.ArgsPassthrough, s)
		}
	}

	config.Theme = color.NewTheme(themePreset)

	if err := config.Theme.OverrideFromEnv(); err != nil {
		return nil, fmt.Errorf("read theme from env: %w", err)
	}

	if b, ok, err := parseBoolEnv("KUBECOLOR_FORCE_COLORS"); err != nil {
		return nil, err
	} else if ok {
		config.ForceColor = b
	}

	if dur, ok, err := parseDurationEnv("KUBECOLOR_OBJ_FRESH"); err != nil {
		return nil, err
	} else if ok {
		config.ObjFreshThreshold = dur
	}

	if cmd := os.Getenv("KUBECTL_COMMAND"); cmd != "" {
		config.KubectlCmd = cmd
	}

	return config, nil
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

func parseThemePresetEnv(env string) (color.Preset, bool, error) {
	value := os.Getenv(env)
	if value == "" {
		return 0, false, nil
	}
	t, err := color.ParsePreset(value)
	if err != nil {
		return 0, false, fmt.Errorf("parse env %s: %w", env, err)
	}
	return t, true, nil
}

func parseThemePresetFlag(flag, value string) (color.Preset, error) {
	if value == "" {
		return 0, fmt.Errorf("flag %s: must be in format %s=value", flag, flag)
	}
	t, err := color.ParsePreset(value)
	if err != nil {
		return 0, fmt.Errorf("flag %s: %w", flag, err)
	}
	return t, nil
}
