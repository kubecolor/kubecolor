package command

import (
	"os"

	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/testutil"
)

func init() {
	os.Clearenv()
	color.ForceColor()
	testutil.DiffAllowUnexported(color.Color{})
}
