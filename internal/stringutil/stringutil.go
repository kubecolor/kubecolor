package stringutil

import (
	"strconv"
	"strings"
	"time"
)

// ParseRatio attempts to parse a ratio delimited by slash, such as "1/2".
func ParseRatio(s string) (left, right string, ok bool) {
	left, right, ok = strings.Cut(s, "/")
	if !ok {
		return "", "", false
	}
	if strings.ContainsRune(right, '/') {
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

func IsAllSameRune(s string, r rune) bool {
	for _, elem := range s {
		if elem != r {
			return false
		}
	}
	return true
}

func CutNumber(s string) (num, after string, found bool) {
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

func SplitAndTrimSpace(s, sep string) []string {
	split := strings.Split(s, sep)
	for i := range split {
		split[i] = strings.TrimSpace(split[i])
	}
	return split
}

// CutSurrounding removes byte at beginning and at end around a string
func CutSurrounding(line string, surrounding byte) (inner string, ok bool) {
	if len(line) < 2 {
		return line, false
	}
	if line[0] == surrounding && line[len(line)-1] == surrounding {
		return line[1 : len(line)-1], true
	}
	return line, false
}

// CutSurroundingAny removes any of the cutset bytes at beginning and at end around a string
//
// NOTE: This function does not support multi-byte runes (e.g emojies), and is deemed undefined behavior.
func CutSurroundingAny(line, cutset string) (quote byte, inner string, ok bool) {
	for _, r := range cutset {
		if inner, ok := CutSurrounding(line, byte(r)); ok {
			return byte(r), inner, true
		}
	}
	return 0, line, false
}

func CutPrefixAny(s string, prefixes ...string) (prefix, after string, ok bool) {
	for _, prefix := range prefixes {
		if after, ok := strings.CutPrefix(s, prefix); ok {
			return prefix, after, true
		}
	}
	return "", s, false
}

func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}
