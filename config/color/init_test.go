package color

import (
	"os"

	"github.com/kubecolor/kubecolor/testutil"
)

func init() {
	os.Clearenv()
	ForceColor()
	testutil.DiffAllowUnexported(Color{})
}
