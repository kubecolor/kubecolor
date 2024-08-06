package printer

import (
	"io"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/internal/stringutil"
	"github.com/kubecolor/kubecolor/kubectl"
)

// KubectlOutputColoredPrinter is a printer to print data depending on
// which kubectl subcommand is executed.
//
// If given subcommand is not supported by the printer,
// it prints data using the theme's default command color.
type KubectlOutputColoredPrinter struct {
	SubcommandInfo    *kubectl.SubcommandInfo
	Recursive         bool
	ObjFreshThreshold time.Duration
	Theme             *config.Theme
	KubecolorVersion  string
}

// ensures it implements the interface
var _ Printer = &KubectlOutputColoredPrinter{}

// Print implements [Printer.Print]
func (p *KubectlOutputColoredPrinter) Print(r io.Reader, w io.Writer) {
	printer := p.getPrinter()
	printer.Print(r, w)
}

func (p *KubectlOutputColoredPrinter) getPrinter() Printer {
	withHeader := !p.SubcommandInfo.NoHeader

	if p.SubcommandInfo.Help {
		return &HelpPrinter{Theme: p.Theme}
	}

	switch p.SubcommandInfo.Subcommand {
	case kubectl.Top, kubectl.APIResources:
		return NewTablePrinter(withHeader, p.Theme, nil)

	case kubectl.APIVersions:
		return NewTablePrinter(false, p.Theme, nil) // api-versions always doesn't have header

	case kubectl.Get, kubectl.Events:
		switch p.SubcommandInfo.FormatOption {
		case kubectl.None, kubectl.Wide:
			return NewTablePrinter(
				withHeader,
				p.Theme,
				func(_ int, column string) string {
					// first try to match a status
					colored, matched := ColorStatus(column, p.Theme)
					if matched {
						return colored
					}

					// When Readiness is "n/m" then yellow
					if left, right, ok := stringutil.ParseRatio(strings.TrimPrefix(column, "Init:")); ok {
						switch {
						case left == "0" && right == "0":
							return p.Theme.Data.Ratio.Zero.Render(column)
						case left == right:
							return p.Theme.Data.Ratio.Equal.Render(column)
						default:
							return p.Theme.Data.Ratio.Unequal.Render(column)
						}
					}

					// Object age when fresh then green
					if age, ok := stringutil.ParseHumanDuration(column); ok {
						if age < p.ObjFreshThreshold {
							return p.Theme.Data.DurationFresh.Render(column)
						}
						return p.Theme.Data.Duration.Render(column)
					}

					return column
				},
			)
		case kubectl.Json:
			return &JSONPrinter{Theme: p.Theme}
		case kubectl.Yaml:
			return &YAMLPrinter{Theme: p.Theme}
		}

	case kubectl.Describe:
		return &DescribePrinter{
			TablePrinter: NewTablePrinter(false, p.Theme, func(_ int, column string) string {
				if colored, ok := ColorStatus(column, p.Theme); ok {
					return colored
				}
				return column
			}),
		}

	case kubectl.Explain:
		return &ExplainPrinter{
			Theme:     p.Theme,
			Recursive: p.Recursive,
		}
	case kubectl.Version:
		switch {
		case p.SubcommandInfo.FormatOption == kubectl.Json:
			return &VersionJSONInjectorPrinter{KubecolorVersion: p.KubecolorVersion, JsonPrinter: &JSONPrinter{Theme: p.Theme}}
		case p.SubcommandInfo.FormatOption == kubectl.Yaml:
			return &VersionYAMLInjectorPrinter{KubecolorVersion: p.KubecolorVersion, YamlPrinter: &YAMLPrinter{Theme: p.Theme}}
		default:
			return &VersionPrinter{
				Theme:            p.Theme,
				KubecolorVersion: p.KubecolorVersion,
			}
		}
	case kubectl.Options:
		return &OptionsPrinter{
			Theme: p.Theme,
		}
	case kubectl.Apply:
		switch p.SubcommandInfo.FormatOption {
		case kubectl.Json:
			return &JSONPrinter{Theme: p.Theme}
		case kubectl.Yaml:
			return &YAMLPrinter{Theme: p.Theme}
		default:
			return &ApplyPrinter{Theme: p.Theme}
		}
	}

	return &SingleColoredPrinter{Color: p.Theme.Default}
}
