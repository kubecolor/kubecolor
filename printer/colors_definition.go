package printer

import "github.com/hidetatz/kubecolor/color"

var (
	// Preset of colors for background
	// Please use them when you just need random colors
	colorsForDarkBackground = []color.Color{
		color.White,
		color.Cyan,
	}

	colorsForLightBackground = []color.Color{
		color.Black,
		color.Blue,
	}

	// colors to be recommended to be used for some context
	// e.g. Json, Yaml, kubectl-describe format etc.

	// colors which look good in dark-backgrounded environment
	KeyColorForDark    = color.Cyan
	StringColorForDark = color.White
	BoolColorForDark   = color.Green
	NumberColorForDark = color.Magenta
	NullColorForDark   = color.Yellow
	HeaderColorForDark = color.White // for plain table

	// colors which look good in light-backgrounded environment
	KeyColorForLight    = color.Blue
	StringColorForLight = color.Black
	BoolColorForLight   = color.Green
	NumberColorForLight = color.Magenta
	NullColorForLight   = color.Yellow
	HeaderColorForLight = color.Black // for plain table
)
