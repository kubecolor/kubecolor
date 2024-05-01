package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
)

type VersionClientPrinter struct {
	Theme *config.Theme
}

// kubectl version --client=true
// Client Version: v1.29.0
// Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
func (vsp *VersionClientPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		key, val, ok := strings.Cut(line, ": ")
		if !ok {
			fmt.Fprintln(w, line)
			continue
		}
		fmt.Fprintf(w, "%s: %s\n",
			getColorByKeyIndent(0, 2, vsp.Theme.Version.Key).Render(key),
			getColorByValueType(val, vsp.Theme).Render(val),
		)
	}
}

type VersionPrinter struct {
	Theme *config.Theme
}

// kubectl version --client=false
// Client Version: v1.29.0
// Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
// Server Version: v1.27.5-gke.200
func (vp *VersionPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, ": ", 2)
		key, val := split[0], split[1]
		key = getColorByKeyIndent(0, 2, vp.Theme.Version.Key).Render(key)

		// val is go struct like
		// version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.2", GitCommit:"f5743093fd1c663cb0cbc89748f730662345d44d", GitTreeState:"clean", BuildDate:"2020-09-16T13:32:58Z", GoVersion:"go1.15", Compiler:"gc", Platform:"linux/amd64"}
		val = strings.TrimRight(val, "}")

		pkgAndValues := strings.SplitN(val, "{", 2)
		packageName := pkgAndValues[0]

		values := strings.Split(pkgAndValues[1], ", ")
		coloredValues := make([]string, len(values))

		fmt.Fprintf(w, "%s: %s{", key, getColorByKeyIndent(2, 2, vp.Theme.Version.Key).Render(packageName))
		for i, value := range values {
			kv := strings.SplitN(value, ":", 2)
			coloredKey := getColorByKeyIndent(0, 2, vp.Theme.Version.Key).Render(kv[0])

			isValDoubleQuotationSurrounded := strings.HasPrefix(kv[1], `"`) && strings.HasSuffix(kv[1], `"`)
			val := strings.TrimRight(strings.TrimLeft(kv[1], `"`), `"`)

			coloredVal := getColorByValueType(kv[1], vp.Theme).Render(val)

			if isValDoubleQuotationSurrounded {
				coloredValues[i] = fmt.Sprintf(`%s:"%s"`, coloredKey, coloredVal)
			} else {
				coloredValues[i] = fmt.Sprintf(`%s:%s`, coloredKey, coloredVal)
			}
		}

		fmt.Fprintf(w, "%s}\n", strings.Join(coloredValues, ", "))
	}
}
