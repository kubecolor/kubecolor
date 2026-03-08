package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/kubecolor/kubecolor/internal/stringutil"
	"github.com/mitchellh/mapstructure"
)

// Duration holds the configurable duration thresholds for age-based coloring.
// Thresholds should be in ascending order (threshold1 < threshold2 < ... < threshold6).
// Each threshold defines the lower bound of a color bucket.
// Ages below threshold1 use [ThemeData.Duration] color.
type Duration struct {
	Threshold1 time.Duration // smallest threshold, e.g "1m"
	Threshold2 time.Duration // e.g "1h"
	Threshold3 time.Duration // e.g "24h"
	Threshold4 time.Duration // e.g "7d"
	Threshold5 time.Duration // e.g "30d"
	Threshold6 time.Duration // largest threshold, e.g "365d"
}

// humanDurationDecodeHook returns a mapstructure decode hook that parses
// human-readable durations like "2d", "7d", "1y" in addition to Go's
// standard duration format ("1h30m", "24h").
func humanDurationDecodeHook() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if to != reflect.TypeOf(time.Duration(0)) {
			return data, nil
		}
		if from.Kind() != reflect.String {
			return data, nil
		}
		s, ok := data.(string)
		if !ok || s == "" {
			return time.Duration(0), nil
		}
		if d, ok := stringutil.ParseHumanDuration(s); ok {
			return d, nil
		}
		return nil, fmt.Errorf("invalid duration %q: must be a human duration (e.g. \"1h\", \"7d\", \"1y\")", s)
	}
}
