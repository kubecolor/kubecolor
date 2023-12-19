package printer

import "github.com/kubecolor/kubecolor/color"

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
	KeyColorForDark      = color.Cyan
	StringColorForDark   = color.White
	TrueColorForDark     = color.Green
	FalseColorForDark    = color.Red
	NumberColorForDark   = color.Magenta
	NullColorForDark     = color.Yellow
	HeaderColorForDark   = color.White // for plain table
	RequiredColorForDark = color.Red   // for `kubectl explain`

	// colors which look good in light-backgrounded environment
	KeyColorForLight      = color.Blue
	StringColorForLight   = color.Black
	TrueColorForLight     = color.Green
	FalseColorForLight    = color.Red
	NumberColorForLight   = color.Magenta
	NullColorForLight     = color.Yellow
	HeaderColorForLight   = color.Black // for plain table
	RequiredColorForLight = color.Red   // for `kubectl explain`
)

type ColorSchema struct {
	DefaultColor, // defaut when no specific mapping is found for the command
	KeyColor, // used to color keys (where ?)
	StringColor, // default color for strings
	TrueColor, // used when value is true
	FalseColor, // used when value is false or an error
	NumberColor, // used when the value is a number
	NullColor, // used when the value is null or a warning
	HeaderColor, // used to print headers
	FreshColor, // color used when the time value is under a certain delay
	RequiredColor color.Color // used when the value is required or is an error
	RandomColor []color.Color // used to display multiple colons, cycle between colors
}

// newColorSchema returns the color schema depending on the chosen theme
func NewColorSchema(dark bool) ColorSchema {

	if dark {
		return ColorSchema{
			DefaultColor:  color.Yellow,
			KeyColor:      color.Cyan,
			StringColor:   color.White,
			TrueColor:     color.Green,
			FalseColor:    color.Red,
			NumberColor:   color.Magenta,
			NullColor:     color.Yellow,
			HeaderColor:   color.White,
			FreshColor:    color.Green,
			RequiredColor: color.Red,
			RandomColor:   []color.Color{color.White, color.Cyan},
		}
	}

	return ColorSchema{
		DefaultColor:  color.Yellow,
		KeyColor:      color.Blue,
		StringColor:   color.Black,
		TrueColor:     color.Green,
		FalseColor:    color.Red,
		NumberColor:   color.Magenta,
		NullColor:     color.Yellow,
		HeaderColor:   color.Black,
		FreshColor:    color.Green,
		RequiredColor: color.Red,
		RandomColor:   []color.Color{color.Black, color.Blue},
	}
}
