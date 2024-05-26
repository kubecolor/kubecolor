package command

import (
	"os"

	"github.com/gookit/color"
	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/testutil"
)

func init() {
	os.Clearenv()
	color.ForceColor()
	color.Enable = true

	testutil.DiffAllowUnexported(config.Color{})
}
