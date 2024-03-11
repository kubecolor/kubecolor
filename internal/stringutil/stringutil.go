package stringutil

import (
	"strconv"
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

func CutNumber(s string) (num string, after string, found bool) {
	if s == "" {
		return "", s, false
	}
	for i, r := range s {
		if IsDigit(r) {
			continue
		}
		if i == 0 {
			return "", s, false
		}
		return s[:i], s[i:], true
	}
	return s, "", true
}

// ParseHumanDuration decodes HumanDuration from [k8s.io/apimachinery/pkg/util/duration]
func ParseHumanDuration(ageString string) (time.Duration, bool) {
	if ageString == "" {
		return 0, false
	}

	var objAgeDuration time.Duration
	rest := ageString

	for range 5 {
		numStr, after, ok := CutNumber(rest)
		if !ok || after == "" {
			return 0, false
		}
		char := after[0]
		rest = after[1:]

		num, err := strconv.ParseUint(numStr, 10, 64)
		if err != nil {
			return 0, false
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

		if rest == "" {
			return objAgeDuration, true
		}
	}

	// should've returned from the end of the loop
	// this means the string contains too many duration elements
	return 0, false
}
