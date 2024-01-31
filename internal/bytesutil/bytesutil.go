package bytesutil

import (
	"bytes"
	"strings"
)

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

var equalOrColon = [][]byte{[]byte(": "), []byte("=")}

func IndexOfEqualOrColon(b []byte) int {
	for _, mark := range equalOrColon {
		colonIndex := bytes.Index(b, mark)
		if colonIndex >= 0 {
			return colonIndex + 1
		}
	}
	return -1
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
