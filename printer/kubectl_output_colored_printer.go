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

	case kubectl.Logs:
		return &LogsPrinter{Theme: p.Theme}

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
		case kubectl.JSON:
			return &JSONPrinter{Theme: p.Theme}
		case kubectl.YAML:
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
		case p.SubcommandInfo.FormatOption == kubectl.JSON:
			return &VersionJSONInjectorPrinter{KubecolorVersion: p.KubecolorVersion, JsonPrinter: &JSONPrinter{Theme: p.Theme}}
		case p.SubcommandInfo.FormatOption == kubectl.YAML:
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

	case
		kubectl.Apply,
		kubectl.Create,
		kubectl.Delete,
		kubectl.Expose,
		kubectl.Patch,
		kubectl.Rollout:
		switch p.SubcommandInfo.FormatOption {
		case kubectl.JSON:
			return &JSONPrinter{Theme: p.Theme}
		case kubectl.YAML:
			return &YAMLPrinter{Theme: p.Theme}
		}
		switch p.SubcommandInfo.Subcommand {
		case kubectl.Apply:
			return &VerbPrinter{
				DryRunColor:   p.Theme.Apply.DryRun,
				FallbackColor: p.Theme.Apply.Fallback,
				VerbColor: map[string]config.Color{
					"created":    p.Theme.Apply.Created,
					"configured": p.Theme.Apply.Configured,
					"unchanged":  p.Theme.Apply.Unchanged,
				},
			}
		case kubectl.Create:
			return &VerbPrinter{
				DryRunColor:   p.Theme.Create.DryRun,
				FallbackColor: p.Theme.Create.Fallback,
				VerbColor: map[string]config.Color{
					"created": p.Theme.Create.Created,
				},
			}
		case kubectl.Delete:
			return &VerbPrinter{
				DryRunColor:   p.Theme.Delete.DryRun,
				FallbackColor: p.Theme.Delete.Fallback,
				VerbColor: map[string]config.Color{
					"deleted": p.Theme.Delete.Deleted,
				},
			}
		case kubectl.Expose:
			return &VerbPrinter{
				DryRunColor:   p.Theme.Expose.DryRun,
				FallbackColor: p.Theme.Expose.Fallback,
				VerbColor: map[string]config.Color{
					"exposed": p.Theme.Expose.Exposed,
				},
			}
		case kubectl.Patch:
			return &VerbPrinter{
				DryRunColor:   p.Theme.Patch.DryRun,
				FallbackColor: p.Theme.Patch.Fallback,
				VerbColor: map[string]config.Color{
					"patched": p.Theme.Patch.Patched,
				},
			}
		case kubectl.Rollout:
			return &VerbPrinter{
				DryRunColor:   p.Theme.Rollout.DryRun,
				FallbackColor: p.Theme.Rollout.Fallback,
				VerbColor: map[string]config.Color{
					"rolled back": p.Theme.Rollout.RolledBack,
					"paused":      p.Theme.Rollout.Paused,
					"resumed":     p.Theme.Rollout.Resumed,
					"restarted":   p.Theme.Rollout.Restarted,
				},
			}
		}
	}

	return &SingleColoredPrinter{Color: p.Theme.Default}
}
