package printer

import (
	"io"
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
	withHeader := !kp.SubcommandInfo.NoHeader

	var printer Printer = &SingleColoredPrinter{Color: kp.Theme.Default} // default in green

	switch kp.SubcommandInfo.Subcommand {
	case kubectl.Top, kubectl.APIResources:
		printer = NewTablePrinter(withHeader, kp.Theme, nil)

	case kubectl.APIVersions:
		printer = NewTablePrinter(false, kp.Theme, nil) // api-versions always doesn't have header

	case kubectl.Get, kubectl.Events:
		switch kp.SubcommandInfo.FormatOption {
		case kubectl.None, kubectl.Wide:
			printer = NewTablePrinter(
				withHeader,
				kp.Theme,
				func(_ int, column string) (config.Color, bool) {
					// first try to match a status
					col, matched := ColorStatus(column, kp.Theme)
					if matched {
						return col, true
					}

					// When Readiness is "n/m" then yellow
					if left, right, ok := stringutil.ParseRatio(column); ok {
						switch {
						case column == "0/0":
							return kp.Theme.Data.Ratio.Zero , true
						case left != right:
							return kp.Theme.Data.Ratio.Unequal, true
						default:
							return kp.Theme.Data.Ratio.Equal, true
						}
					}

					// Object age when fresh then green
					if age, ok := stringutil.ParseHumanDuration(column); ok {
						if age < kp.ObjFreshThreshold {
							return kp.Theme.DurationFresh, true
						}
					}

					return config.Color{}, false
				},
			)
		case kubectl.Json:
			printer = &JsonPrinter{Theme: kp.Theme}
		case kubectl.Yaml:
			printer = &YamlPrinter{Theme: kp.Theme}
		}

	case kubectl.Describe:
		printer = &DescribePrinter{
			TablePrinter: NewTablePrinter(false, kp.Theme, func(_ int, column string) (config.Color, bool) {
				return ColorStatus(column, kp.Theme)
			}),
		}
	case kubectl.Explain:
		printer = &ExplainPrinter{
			Theme:     kp.Theme,
			Recursive: kp.Recursive,
		}
	case kubectl.Version:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{Theme: kp.Theme}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			printer = &YamlPrinter{Theme: kp.Theme}
		case kp.SubcommandInfo.Client:
			printer = &VersionClientPrinter{
				Theme: kp.Theme,
			}
		default:
			printer = &VersionClientPrinter{
				Theme: kp.Theme,
			}
		}
	case kubectl.Options:
		printer = &OptionsPrinter{
			Theme: kp.Theme,
		}
	case kubectl.Apply:
		switch kp.SubcommandInfo.FormatOption {
		case kubectl.Json:
			printer = &JsonPrinter{Theme: kp.Theme}
		case kubectl.Yaml:
			printer = &YamlPrinter{Theme: kp.Theme}
		default:
			printer = &ApplyPrinter{Theme: kp.Theme}
		}
	}

	if kp.SubcommandInfo.Help {
		printer = &SingleColoredPrinter{Color: kp.Theme.Default}
	}

	printer.Print(r, w)
}
