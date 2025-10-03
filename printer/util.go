package printer

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/color"
	"github.com/kubecolor/kubecolor/internal/stringutil"
)

// ColorDataKey returns a color based on the given indent.
// When you want to change key color based on indent depth (e.g. Json, Yaml), use this function
func ColorDataKey(indent, basicIndentWidth int, colors color.Slice) color.Color {
	if len(colors) == 0 {
		return color.Color{}
	}
	return colors[indent/basicIndentWidth%len(colors)]
}

var isQuantityRegex = regexp.MustCompile(`^[\+\-]?(?:\d+|\.\d+|\d+\.|\d+\.\d+)?(?:m|[kMGTPE]i?)$`)

// ColorDataValue returns a color by value.
// This is intended to be used to colorize any structured data e.g. JSON, YAML.
//
// Logic:
//
// - "null" return "data.null" theme color
// - "true" return "data.true" theme color
// - "false" return "data.false" theme color
// - "123" return "data.number" theme color
// - "5Gi" return "data.quantity" theme color
// - "15m10s" return "data.duration" theme color
// - otherwise, return "data.string" theme color
func ColorDataValue(val string, theme *config.Theme) color.Color {
	c, _ := TryColorDataValue(val, theme)
	return c
}

func TryColorDataValue(val string, theme *config.Theme) (color.Color, bool) {
	switch val {
	case "null", "<none>", "<unknown>", "<unset>", "<nil>", "<invalid>":
		return theme.Data.Null, true
	case "true", "True":
		return theme.Data.True, true
	case "false", "False":
		return theme.Data.False, true
	}

	// Ints: 123
	if stringutil.IsOnlyDigits(val) {
		return theme.Data.Number, true
	}

	// Floats: 123.456
	if left, right, ok := strings.Cut(val, "."); ok {
		if stringutil.IsOnlyDigits(left) && stringutil.IsOnlyDigits(right) {
			return theme.Data.Number, true
		}
	}

	// Quantity: 100m, 5Gi
	if isQuantityRegex.MatchString(val) {
		return theme.Data.Quantity, true
	}

	// Duration: 15m10s
	if _, ok := stringutil.ParseHumanDuration(val); ok {
		return theme.Data.Duration, true
	}

	return theme.Data.String, false
}

// ColorStatus returns the color that should be used for a given status text.
func ColorStatus(status string, theme *config.Theme) (string, bool) {
	if strings.ContainsRune(status, ',') {
		statuses := strings.Split(status, ",")
		any := false
		for i, s := range statuses {
			if colored, ok := colorSingleStatus(s, theme); ok {
				statuses[i] = colored
				any = true
			}
		}
		if !any {
			return status, false
		}
		return strings.Join(statuses, ","), true
	}

	return colorSingleStatus(status, theme)
}

func colorSingleStatus(status string, theme *config.Theme) (string, bool) {
	switch strings.TrimPrefix(status, "Init:") {
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
		"FailedCreate",
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
		"PreCreateHookError",
		"PreStartHookError",
		"PostStartHookError",
		"FailedPostStartHook",
		"FailedPreStopHook",
		// Node status list
		"NotReady",
		"NetworkUnavailable",

		// some other status
		"ContainerStatusUnknown",
		"CreateContainerConfigError",
		"CreateContainerError",
		"ContainerCannotRun",
		"CrashLoopBackOff",
		"DeadlineExceeded",
		"ImagePullBackOff",
		"Evicted",
		"FailedScheduling",
		"Error",
		"ErrImagePull",
		"OOMKilled",
		"RunContainerError",
		"StartError",
		// PVC status
		"Lost":
		return theme.Status.Error.Render(status), true
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
		"SchedulingDisabled",
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
		return theme.Status.Warning.Render(status), true
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
		return theme.Status.Success.Render(status), true
	}
	return status, false
}

// findIndent returns a length of indent (spaces at left) in the given line
func findIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

// isAllUpper is used to identity header lines like this:
//
//	NAME  READY  STATUS  RESTARTS  AGE
//
// Commonly found in "kubectl get" output
func isAllUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// isOnlySymbols is used to identity header underline like this:
//
//	Resources  Non-Resource URLs  Resource Names  Verbs
//	---------  -----------------  --------------  -----
//
// Commonly found in "kubectl describe" output
func isOnlySymbols(s string) bool {
	anyPuncts := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return false
		}
		if unicode.IsPunct(r) {
			anyPuncts = true
		}
	}
	return anyPuncts
}
