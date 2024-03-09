package printer

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/kubecolor/kubecolor/config"
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
		return theme.False, true
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
		return theme.Null, true
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
		return theme.True, true
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
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.None, kp.SubcommandInfo.FormatOption == kubectl.Wide:
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
					if strings.Count(column, "/") == 1 {
						if arr := strings.Split(column, "/"); arr[0] != arr[1] {
							_, e1 := strconv.Atoi(arr[0])
							_, e2 := strconv.Atoi(arr[1])
							if e1 == nil && e2 == nil { // check both is number
								// TODO: Replace with theme color
								return config.MustParseColor("yellow"), true
							}
						}

					}

					// Object age when fresh then green
					if checkIfObjFresh(column, kp.ObjFreshThreshold) {
						return kp.Theme.DurationFresh, true
					}

					return config.Color{}, false
				},
			)
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{Theme: kp.Theme}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
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
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{Theme: kp.Theme}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
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

func checkIfObjFresh(ageString string, threshold time.Duration) bool {
	// decode HumanDuration from [k8s.io/apimachinery/pkg/util/duration]
	var objAgeDuration time.Duration
	rest := ageString

	for range 5 {
		var num int64
		var char rune
		varsRead, _ := fmt.Sscanf(rest, "%d%c%s", &num, &char, &rest)
		if varsRead < 2 {
			break
		} else if varsRead < 3 {
			rest = ""
		}

		switch char {
		case 'y':
			objAgeDuration += 365 * 24 * time.Hour * time.Duration(num)
		case 'd':
			objAgeDuration += 24 * time.Hour * time.Duration(num)
		case 'h':
			objAgeDuration += time.Hour * time.Duration(num)
		case 'm':
			objAgeDuration += time.Minute * time.Duration(num)
		case 's':
			objAgeDuration += time.Second * time.Duration(num)
		}
	}

	return objAgeDuration > 0 && objAgeDuration < threshold
}
