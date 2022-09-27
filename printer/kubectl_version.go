package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/color"
)

type VersionShortPrinter struct {
	DarkBackground bool
}

// kubectl version --short format
// Client Version: v1.19.3
// Server Version: v1.19.2
func (vsp *VersionShortPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ": ")
		key, val := splitted[0], splitted[1]
		fmt.Fprintf(w, "%s: %s\n",
			color.Apply(key, getColorByKeyIndent(0, 2, vsp.DarkBackground)),
			color.Apply(val, getColorByValueType(val, vsp.DarkBackground)),
		)
	}
}

type VersionPrinter struct {
	DarkBackground bool
}

func (vp *VersionPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.SplitN(line, ": ", 2)
		key, val := splitted[0], splitted[1]
		key = color.Apply(key, getColorByKeyIndent(0, 2, vp.DarkBackground))

		// val is go struct like
		// version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.2", GitCommit:"f5743093fd1c663cb0cbc89748f730662345d44d", GitTreeState:"clean", BuildDate:"2020-09-16T13:32:58Z", GoVersion:"go1.15", Compiler:"gc", Platform:"linux/amd64"}
		val = strings.TrimRight(val, "}")

		pkgAndValues := strings.SplitN(val, "{", 2)
		packageName := pkgAndValues[0]

		values := strings.Split(pkgAndValues[1], ", ")
		coloredValues := make([]string, len(values))

		fmt.Fprintf(w, "%s: %s{", key, color.Apply(packageName, getColorByKeyIndent(2, 2, vp.DarkBackground)))
		for i, value := range values {
			kv := strings.SplitN(value, ":", 2)
			coloredKey := color.Apply(kv[0], getColorByKeyIndent(0, 2, vp.DarkBackground))

			isValDoubleQuotationSurrounded := strings.HasPrefix(kv[1], `"`) && strings.HasSuffix(kv[1], `"`)
			val := strings.TrimRight(strings.TrimLeft(kv[1], `"`), `"`)

			coloredVal := color.Apply(val, getColorByValueType(kv[1], vp.DarkBackground))

			if isValDoubleQuotationSurrounded {
				coloredValues[i] = fmt.Sprintf(`%s:"%s"`, coloredKey, coloredVal)
			} else {
				coloredValues[i] = fmt.Sprintf(`%s:%s`, coloredKey, coloredVal)
			}
		}

		fmt.Fprintf(w, "%s}\n", strings.Join(coloredValues, ", "))
	}
}
