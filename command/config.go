package command

import (
	"cmp"
	"fmt"
	"os"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/spf13/viper"
)

type Config struct {
	*config.Config

	ForceColor           ColorLevel
	ShowKubecolorVersion bool
	StdinOverride        string

	ArgsPassthrough []string
	Flags           FlagSet
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
			v.Set(config.PresetKey, "light")
		} else {
			v.Set(config.PresetKey, "dark")
		}
	}

	if c, ok, err := parseColorLevelEnv("KUBECOLOR_FORCE_COLORS"); err != nil {
		return nil, err
	} else if ok {
		cfg.ForceColor = c
	}

	var (
		flagPlain = cfg.Flags.NewBool("--plain", "Disable colored output.")

		flagLightBg = cfg.Flags.NewBool("--light-background", "Switches to light theme, or dark when --light-background=false. Same as doing --kubecolor-theme=light or --kubecolor-theme=dark.")

		flagForceVal = ColorLevelAuto // value used when no flag value
		flagForce    = cfg.Flags.NewString("--force-colors", "Overrides the automatic color support detection. Overrides the KUBECOLOR_FORCE_COLORS env var.").
				WithUnmarshaller(&flagForceVal)

		flagVersion = cfg.Flags.NewBool("--kubecolor-version", "Print the kubecolor version and then exit.")

		flagStdin = cfg.Flags.NewString("--kubecolor-stdin", "Read command input from stdin or file instead of executing kubectl.")

		flagTheme = cfg.Flags.NewString("--kubecolor-theme", "Set kubecolor theme preset, e.g dark or light. Overrides the KUBECOLOR_PRESET env var.").
				WithRequiresValue()

		flagPager = cfg.Flags.NewString("--pager", `Set kubecolor pager, e.g "less -RF" or "more". Overrides the KUBECOLOR_PAGER and PAGER env vars.`).
				WithRequiresValue()

		flagPagingVal = config.PagingAuto // value used when no flag value
		flagPaging    = cfg.Flags.NewString("--paging", `Pipe kubecolor output into pager.`).
				WithUnmarshaller(&flagPagingVal)

		flagNoPaging = cfg.Flags.NewBool("--no-paging", `Disable paging. Alias to --paging=never.`)
	)

	for _, s := range inputArgs {
		f, err := cfg.Flags.ParseArg(s)
		if err != nil {
			return nil, err
		}
		switch f {
		case flagPlain:
			if f.BoolValue() {
				cfg.ForceColor = ColorLevelNone
			}
		case flagLightBg:
			if f.BoolValue() {
				v.Set(config.PresetKey, "light")
			} else {
				v.Set(config.PresetKey, "dark")
			}
		case flagForce:
			cfg.ForceColor = flagForceVal
		case flagVersion:
			cfg.ShowKubecolorVersion = f.BoolValue()
		case flagStdin:
			// Value means "read from file"
			// Dash "-" means "read from stdin"
			// Empty, as in just "--kubecolor-stdin", means "read from stdin"
			cfg.StdinOverride = cmp.Or(f.Value, "-")
		case flagTheme:
			v.Set(config.PresetKey, f.Value)
		case flagPager:
			v.Set("pager", f.Value)
		case flagPaging:
			// mapstructure doesn't like "type X string" values,
			// so we have to convert it via string(...)
			v.Set("paging", string(flagPagingVal))
		case flagNoPaging:
			if f.BoolValue() {
				v.Set("paging", string(config.PagingNever))
			}
		default:
			cfg.ArgsPassthrough = append(cfg.ArgsPassthrough, s)
		}
	}

	newCfg, err := config.Unmarshal(v)
	if err != nil {
		return nil, err
	}
	cfg.Config = newCfg

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

func parseBoolEnv(env string) (result bool, ok bool, err error) {
	result, ok, err = parseBool(os.Getenv(env))
	if err != nil {
		return false, false, fmt.Errorf("parse env %s: %w", env, err)
	}
	return result, ok, err
}
