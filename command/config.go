package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/printer"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
)

type KubecolorConfig struct {
	Plain                bool
	DarkBackground       bool
	ForceColor           bool
	ShowKubecolorVersion bool
	KubectlCmd           string
	ObjFreshThreshold    time.Duration
	ColorSchema          printer.ColorSchema
}

func ResolveConfig(args []string) ([]string, *KubecolorConfig) {
	args, plainFlagFound := findAndRemoveBoolFlagIfExists(args, "--plain")
	args, lightBackgroundFlagFound := findAndRemoveBoolFlagIfExists(args, "--light-background")
	args, forceColorFlagFound := findAndRemoveBoolFlagIfExists(args, "--force-colors")
	args, kubecolorVersionFlagFound := findAndRemoveBoolFlagIfExists(args, "--kubecolor-version")

	darkBackground := !lightBackgroundFlagFound

	colorsForcedViaEnv := os.Getenv("KUBECOLOR_FORCE_COLORS") == "true"

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
			fmt.Printf("[WARN] [kubecolor] cannot parse duration taken from env '%s'. See kubecolor document. err: %v\n", objFreshAgeThresholdEnv, err)
		}
	}

	// parse the config file, deafult to ~/.kubecolor
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configFile := fmt.Sprintf("%s/.kubecolor.yml", homeDir)
	configFileEnv := "KUBECOLOR_CONFIG_FILE"
	if customConfigFile := os.Getenv(configFileEnv); customConfigFile != "" {
		if strings.HasPrefix(customConfigFile, "~/") {
			configFile = filepath.Join(homeDir, customConfigFile[2:])
		} else {
			configFile = customConfigFile
		}
	}
	config.WithOptions(config.ParseEnv)
	// config.AddDriver(json.Driver)
	config.AddDriver(yamlv3.Driver)

	err = config.LoadFiles(configFile)
	if err != nil {
		fmt.Printf("[ERROR] [kubecolor] cannot parse config file from '%s'. See kubecolor docs. err: %v\n", configFile, err)
	}

	// Parse the color schema if provided
	// format is like "default:32;key:36;string:37;true:32;false:31;number:35;null:33;header:37;fresh:32;required:31;random:36,37"
	ColorSchema := printer.NewColorSchema(darkBackground)

	// Set the color schema based on the config file
	for colorKey, colorName := range config.StringMap("theme.all") {

		switch colorKey {
		case "default":
			ColorSchema.DefaultColor = color.StringToColor(colorName)
		case "key":
			ColorSchema.KeyColor = color.StringToColor(colorName)
		case "string":
			ColorSchema.StringColor = color.StringToColor(colorName)
		case "true":
			ColorSchema.TrueColor = color.StringToColor(colorName)
		case "false":
			ColorSchema.FalseColor = color.StringToColor(colorName)
		case "number":
			ColorSchema.NumberColor = color.StringToColor(colorName)
		case "null":
			ColorSchema.NumberColor = color.StringToColor(colorName)
		case "header":
			ColorSchema.HeaderColor = color.StringToColor(colorName)
		case "fresh":
			ColorSchema.NumberColor = color.StringToColor(colorName)
		case "required":
			ColorSchema.RequiredColor = color.StringToColor(colorName)
		case "tabs":
			for _, val := range config.Strings("theme.all.tabs") {
				ColorSchema.RandomColor = append(ColorSchema.RandomColor, color.StringToColor(val))
			}
		default:
		}
	}

	customColorEnv := "KUBECOLOR_CUSTOM_COLOR"
	if customColor := os.Getenv(customColorEnv); customColor != "" {
		customSelections := strings.Split(customColor, ";")
		for _, customSelection := range customSelections {
			keyVal := strings.Split(customSelection, ":")
			if len(keyVal) != 2 {
				fmt.Printf("[WARN] [kubecolor] cannot parse custom color selection taken from env '%s'. See kubecolor docs.\n", customSelection)
				break
			}

			// Check the random case as the value if a list
			if keyVal[0] == "random" {
				rndColors := strings.Split(keyVal[1], ",")
				for _, newColor := range rndColors {
					val, err := strconv.Atoi(newColor)
					if err != nil {
						fmt.Printf("[WARN] [kubecolor] cannot parse custom color taken from env %s. See kubecolor document. %v\n", customSelection, err)
						break
					}
					ColorSchema.RandomColor = append(ColorSchema.RandomColor, color.Color(val))
				}
				continue
			}

			key := keyVal[0]
			val, err := strconv.Atoi(keyVal[1])
			if err != nil {
				fmt.Printf("[WARN] [kubecolor] cannot parse custom color taken from env %s. See kubecolor document. %v\n", customSelection, err)
				break
			}
			colorName := color.Color(val)

			switch key {
			case "default":
				ColorSchema.DefaultColor = colorName
			case "key":
				ColorSchema.KeyColor = colorName
			case "string":
				ColorSchema.StringColor = colorName
			case "true":
				ColorSchema.TrueColor = colorName
			case "false":
				ColorSchema.FalseColor = colorName
			case "number":
				ColorSchema.NumberColor = colorName
			case "null":
				ColorSchema.NumberColor = colorName
			case "header":
				ColorSchema.HeaderColor = colorName
			case "fresh":
				ColorSchema.NumberColor = colorName
			case "required":
				ColorSchema.RequiredColor = colorName
			default:
			}
		}
	}

	return args, &KubecolorConfig{
		Plain:                plainFlagFound,
		DarkBackground:       darkBackground,
		ForceColor:           forceColorFlagFound || colorsForcedViaEnv,
		ShowKubecolorVersion: kubecolorVersionFlagFound,
		KubectlCmd:           kubectlCmd,
		ObjFreshThreshold:    objFreshAgeThresholdDuration,
		ColorSchema:          ColorSchema,
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
