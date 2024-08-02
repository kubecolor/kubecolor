package printer

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/internal/stringutil"
)

// ColorDataKey returns a color based on the given indent.
// When you want to change key color based on indent depth (e.g. Json, Yaml), use this function
func ColorDataKey(indent int, basicIndentWidth int, colors config.ColorSlice) config.Color {
	if len(colors) == 0 {
		return config.Color{}
	}
	return colors[indent/basicIndentWidth%len(colors)]
}

var (
	isQuantityRegex = regexp.MustCompile(`^[\+\-]?(?:\d+|\.\d+|\d+\.|\d+\.\d+)?(?:m|[kMGTPE]i?)$`)
)

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
func ColorDataValue(val string, theme *config.Theme) config.Color {
	switch val {
	case "null", "<none>", "<unknown>", "<unset>", "<nil>", "<invalid>":
		return theme.Data.Null
	case "true", "True":
		return theme.Data.True
	case "false", "False":
		return theme.Data.False
	}

	// Ints: 123
	if stringutil.IsOnlyDigits(val) {
		return theme.Data.Number
	}

	// Floats: 123.456
	if left, right, ok := strings.Cut(val, "."); ok {
		if stringutil.IsOnlyDigits(left) && stringutil.IsOnlyDigits(right) {
			return theme.Data.Number
		}
	}

	// Quantity: 100m, 5Gi
	if isQuantityRegex.MatchString(val) {
		return theme.Data.Quantity
	}

	// Duration: 15m10s
	if _, ok := stringutil.ParseHumanDuration(val); ok {
		return theme.Data.Duration
	}

	return theme.Data.String
}

// ColorStatus returns the color that should be used for a given status text.
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

// toSpaces returns repeated spaces whose length is n.
func toSpaces(n int) string {
	return strings.Repeat(" ", n)
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
