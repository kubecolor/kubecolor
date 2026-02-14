package printer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/internal/bytesutil"
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
func (p *VersionPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, bytesutil.MaxLineLength)
	any := false
	for scanner.Scan() {
		line := scanner.Text()
		key, val, ok := strings.Cut(line, ": ")
		if !ok {
			fmt.Fprintln(w, line)
			continue
		}
		fmt.Fprintf(w, "%s: %s\n",
			ColorDataKey(0, 2, p.Theme.Version.Key).Render(key),
			ColorDataValue(val, p.Theme).Render(val),
		)
		any = true
	}
	if err := scanner.Err(); err != nil {
		slog.Error("Failed to print version output.", "error", err)
	}

	// Check if any got printed, so we don't print version after an error like
	// 	error: invalid argument "foo" for "--client" flag: strconv.ParseBool: parsing "foo": invalid syntax
	if any {
		key := "Kubecolor Version"
		val := p.KubecolorVersion
		fmt.Fprintf(w, "%s: %s\n",
			ColorDataKey(0, 2, p.Theme.Version.Key).Render(key),
			ColorDataValue(val, p.Theme).Render(val),
		)
	}
}

type VersionJSONInjectorPrinter struct {
	JsonPrinter      *JSONPrinter
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
	YamlPrinter      *YAMLPrinter
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
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(output); err != nil {
		w.Write(b)
		return
	}
	vp.YamlPrinter.Print(&buf, w)
}
