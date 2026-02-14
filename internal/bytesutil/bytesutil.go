package bytesutil

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MaxLineLength is the maximum size the line buffer will grow to.
//
// Value is based on the maximum file size in Secrets and ConfigMaps (1 MiB),
// but base64-encoded which increases the size by ~33% (rounded up to 1.5MB)
// [https://kubernetes.io/docs/concepts/configuration/secret/#restriction-data-size]
const MaxLineLength = 1500000

func IndexOfNonSpace(b []byte, spaceCharset string) int {
	for i := 0; i < len(b); i++ {
		if !strings.ContainsRune(spaceCharset, rune(b[i])) {
			return i
		}
	}
	return -1
}

var doubleSpace = []byte("  ")

func IndexOfDoubleSpace(b []byte) int {
	spaceIndex := bytes.Index(b, doubleSpace)
	tabIndex := bytes.IndexByte(b, '\t')
	if spaceIndex >= 0 && tabIndex >= 0 {
		return min(spaceIndex, tabIndex)
	}
	if spaceIndex >= 0 {
		return spaceIndex
	}
	return tabIndex
}

func CountColumns(b []byte, spaceCharset string) int {
	var count int
	for {
		index := IndexOfNonSpace(b, spaceCharset)
		if index == -1 {
			break
		}
		b = b[index:]
		count++

		index = IndexOfDoubleSpace(b)
		if index == -1 {
			break
		}
		b = b[index:]
	}
	return count
}

func IsOnlyDigits(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	var index int
	for {
		r, size := utf8.DecodeRune(b[index:])
		if size == 0 {
			// EOF
			return true
		}
		if r == utf8.RuneError {
			return false
		}
		index += size
		if !unicode.IsDigit(r) {
			return false
		}
	}
}

func IsOnlyLetters(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	var index int
	for {
		r, size := utf8.DecodeRune(b[index:])
		if size == 0 {
			// EOF
			return true
		}
		if r == utf8.RuneError {
			return false
		}
		index += size
		if !unicode.IsLetter(r) {
			return false
		}
	}
}
