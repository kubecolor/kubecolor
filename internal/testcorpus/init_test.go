package testcorpus

import (
	"os"

	"github.com/gookit/color"
)

func init() {
	os.Clearenv()
	color.ForceColor()
	color.Enable = true
}
