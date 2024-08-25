package testcorpus

import (
	"os"

	"github.com/kubecolor/kubecolor/config/color"
)

func init() {
	os.Clearenv()
	color.ForceColor()
}
