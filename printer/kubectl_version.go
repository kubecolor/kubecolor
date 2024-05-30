package printer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"gopkg.in/yaml.v3"
)

type VersionPrinter struct {
	Theme            *config.Theme
	KubecolorVersion string
}

// kubectl version --client=false
// Client Version: v1.29.0
// Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
// Server Version: v1.27.5-gke.200
func (vp *VersionPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	any := false
	for scanner.Scan() {
		line := scanner.Text()
		key, val, ok := strings.Cut(line, ": ")
		if !ok {
			fmt.Fprintln(w, line)
			continue
		}
		fmt.Fprintf(w, "%s: %s\n",
			getColorByKeyIndent(0, 2, vp.Theme.Version.Key).Render(key),
			getColorByValueType(val, vp.Theme).Render(val),
		)
		any = true
	}

	// Check if any got printed, so we don't print version after an error like
	// 	error: invalid argument "foo" for "--client" flag: strconv.ParseBool: parsing "foo": invalid syntax
	if any {
		key := "Kubecolor Version"
		val := vp.KubecolorVersion
		fmt.Fprintf(w, "%s: %s\n",
			getColorByKeyIndent(0, 2, vp.Theme.Version.Key).Render(key),
			getColorByValueType(val, vp.Theme).Render(val),
		)
	}
}

type VersionJSONInjectorPrinter struct {
	JsonPrinter      *JsonPrinter
	KubecolorVersion string
}

func (vp *VersionJSONInjectorPrinter) Print(r io.Reader, w io.Writer) {
	b, err := io.ReadAll(r)
	if err != nil {
		return
	}
	var output map[string]any
	if err := json.Unmarshal(b, &output); err != nil {
		w.Write(b)
		return
	}
	output["kubecolorVersion"] = vp.KubecolorVersion
	result, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		w.Write(b)
		return
	}
	vp.JsonPrinter.Print(bytes.NewReader(result), w)
}

type VersionYAMLInjectorPrinter struct {
	YamlPrinter      *YamlPrinter
	KubecolorVersion string
}

func (vp *VersionYAMLInjectorPrinter) Print(r io.Reader, w io.Writer) {
	b, err := io.ReadAll(r)
	if err != nil {
		return
	}
	var output map[string]any
	if err := yaml.Unmarshal(b, &output); err != nil {
		w.Write(b)
		return
	}
	output["kubecolorVersion"] = vp.KubecolorVersion
	result, err := yaml.Marshal(output)
	if err != nil {
		w.Write(b)
		return
	}
	vp.YamlPrinter.Print(bytes.NewReader(result), w)
}
