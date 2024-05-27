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
type KubectlOutputColoredPrinter struct {
	SubcommandInfo    *kubectl.SubcommandInfo
	Recursive         bool
	ObjFreshThreshold time.Duration
	Theme             *config.Theme
	KubecolorVersion  string
}

func ColorStatus(status string, theme *config.Theme) (config.Color, bool) {
	switch status {
	case
		// from https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/events/event.go
		// Container event reason list
		"Failed",
		"BackOff",
		"ExceededGracePeriod",
		// Pod event reason list
		"FailedKillPod",
		"FailedCreatePodContainer",
		// "Failed",
		"NetworkNotReady",
		// Image event reason list
		// "Failed",
		"InspectFailed",
		"ErrImageNeverPull",
		// "BackOff",
		// kubelet event reason list
		"NodeNotSchedulable",
		"KubeletSetupFailed",
		"FailedAttachVolume",
		"FailedMount",
		"VolumeResizeFailed",
		"FileSystemResizeFailed",
		"FailedMapVolume",
		"ContainerGCFailed",
		"ImageGCFailed",
		"FailedNodeAllocatableEnforcement",
		"FailedCreatePodSandBox",
		"FailedPodSandBoxStatus",
		"FailedMountOnFilesystemMismatch",
		// Image manager event reason list
		"InvalidDiskCapacity",
		"FreeDiskSpaceFailed",
		// Probe event reason list
		"Unhealthy",
		// Pod worker event reason list
		"FailedSync",
		// Config event reason list
		"FailedValidation",
		// Lifecycle hooks
		"FailedPostStartHook",
		"FailedPreStopHook",
		// Node status list
		"NotReady",
		"NetworkUnavailable",

		// some other status
		"ContainerStatusUnknown",
		"CrashLoopBackOff",
		"ImagePullBackOff",
		"Evicted",
		"FailedScheduling",
		"Error",
		"ErrImagePull",
		"OOMKilled",
		// PVC status
		"Lost":
		return theme.Status.Error, true
	case
		// from https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/events/event.go
		// Container event reason list
		"Killing",
		"Preempting",
		// Pod event reason list
		// Image event reason list
		"Pulling",
		// kubelet event reason list
		"NodeNotReady",
		"NodeSchedulable",
		"Starting",
		"AlreadyMountedVolume",
		"SuccessfulAttachVolume",
		"SuccessfulMountVolume",
		"NodeAllocatableEnforced",
		// Image manager event reason list
		// Probe event reason list
		"ProbeWarning",
		// Pod worker event reason list
		// Config event reason list
		// Lifecycle hooks

		// some other status
		"Pending",
		"ContainerCreating",
		"PodInitializing",
		"Terminating",
		"Terminated",
		"Warning",

		// PV reclaim policy
		"Delete",

		// PVC status
		"Available",
		"Released",

		"ScalingReplicaSet":
		return theme.Status.Warning, true
	case
		"Running",
		"Completed",
		"Pulled",
		"Created",
		"Rebooted",
		"NodeReady",
		"Started",
		"Normal",
		"VolumeResizeSuccessful",
		"FileSystemResizeSuccessful",
		"Ready",
		"Scheduled",
		"SuccessfulCreate",

		// PV reclaim policy
		"Retain",

		// PVC status
		"Bound":
		return theme.Status.Success, true
	}
	return config.Color{}, false
}

// Print reads r then write it to w, its format is based on kubectl subcommand.
// If given subcommand is not supported by the printer, it prints data in Green.
func (kp *KubectlOutputColoredPrinter) Print(r io.Reader, w io.Writer) {
	printer := kp.getPrinter()
	printer.Print(r, w)
}

func (kp *KubectlOutputColoredPrinter) getPrinter() Printer {
	withHeader := !kp.SubcommandInfo.NoHeader

	if kp.SubcommandInfo.Help {
		return &HelpPrinter{Theme: kp.Theme}
	}

	switch kp.SubcommandInfo.Subcommand {
	case kubectl.Top, kubectl.APIResources:
		return NewTablePrinter(withHeader, kp.Theme, nil)

	case kubectl.APIVersions:
		return NewTablePrinter(false, kp.Theme, nil) // api-versions always doesn't have header

	case kubectl.Get, kubectl.Events:
		switch kp.SubcommandInfo.FormatOption {
		case kubectl.None, kubectl.Wide:
			return NewTablePrinter(
				withHeader,
				kp.Theme,
				func(_ int, column string) (config.Color, bool) {
					column = strings.TrimPrefix(column, "Init:")

					// first try to match a status
					col, matched := ColorStatus(column, kp.Theme)
					if matched {
						return col, true
					}

					// When Readiness is "n/m" then yellow
					if left, right, ok := stringutil.ParseRatio(column); ok {
						switch {
						case column == "0/0":
							return kp.Theme.Data.Ratio.Zero, true
						case left == right:
							return kp.Theme.Data.Ratio.Equal, true
						default:
							return kp.Theme.Data.Ratio.Unequal, true
						}
					}

					// Object age when fresh then green
					if age, ok := stringutil.ParseHumanDuration(column); ok {
						if age < kp.ObjFreshThreshold {
							return kp.Theme.Data.DurationFresh, true
						}
						return kp.Theme.Data.Duration, true
					}

					return config.Color{}, false
				},
			)
		case kubectl.Json:
			return &JsonPrinter{Theme: kp.Theme}
		case kubectl.Yaml:
			return &YamlPrinter{Theme: kp.Theme}
		}

	case kubectl.Describe:
		return &DescribePrinter{
			TablePrinter: NewTablePrinter(false, kp.Theme, func(_ int, column string) (config.Color, bool) {
				return ColorStatus(column, kp.Theme)
			}),
		}
	case kubectl.Explain:
		return &ExplainPrinter{
			Theme:     kp.Theme,
			Recursive: kp.Recursive,
		}
	case kubectl.Version:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			return &VersionJSONInjectorPrinter{KubecolorVersion: kp.KubecolorVersion, JsonPrinter: &JsonPrinter{Theme: kp.Theme}}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			return &VersionYAMLInjectorPrinter{KubecolorVersion: kp.KubecolorVersion, YamlPrinter: &YamlPrinter{Theme: kp.Theme}}
		default:
			return &VersionPrinter{
				Theme:            kp.Theme,
				KubecolorVersion: kp.KubecolorVersion,
			}
		}
	case kubectl.Options:
		return &OptionsPrinter{
			Theme: kp.Theme,
		}
	case kubectl.Apply:
		switch kp.SubcommandInfo.FormatOption {
		case kubectl.Json:
			return &JsonPrinter{Theme: kp.Theme}
		case kubectl.Yaml:
			return &YamlPrinter{Theme: kp.Theme}
		default:
			return &ApplyPrinter{Theme: kp.Theme}
		}
	}

	return &SingleColoredPrinter{Color: kp.Theme.Default}
}
