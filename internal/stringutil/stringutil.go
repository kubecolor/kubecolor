package stringutil

import (
	"fmt"
	"strings"
	"time"
)

// ParseRatio attempts to parse a ratio delimited by slash, such as "1/2".
func ParseRatio(s string) (left, right string, ok bool) {
	if strings.Count(s, "/") != 1 {
		return "", "", false
	}
	left, right, ok = strings.Cut(s, "/")
	if !ok {
		return "", "", false
	}
	if left == "" || right == "" {
		return "", "", false
	}
	if !IsOnlyDigits(left) || !IsOnlyDigits(right) {
		return "", "", false
	}
	return left, right, ok
}

func IsOnlyDigits(s string) bool {
	for _, r := range s {
		if !IsDigit(r) {
			return false
		}
	}
	return true
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// ParseHumanDuration decodes HumanDuration from [k8s.io/apimachinery/pkg/util/duration]
func ParseHumanDuration(ageString string) (time.Duration, bool) {
	if ageString == "" {
		return 0, false
	}

	var objAgeDuration time.Duration
	rest := ageString

	for range 5 {
		var num int64
		var char rune
		varsRead, _ := fmt.Sscanf(rest, "%d%c%s", &num, &char, &rest)
		if varsRead < 2 {
			return 0, false
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
		default:
			return 0, false
		}
	}

	return objAgeDuration, true
}
