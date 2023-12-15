package describe

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var doubleSpace = []byte{' ', ' '}

type Line struct {
	Indent   []byte
	Key      []byte
	Spacing  []byte
	Value    []byte
	Trailing []byte
}

func (line Line) KeyIndent() int {
	return len(line.Indent)
}

func (line Line) ValueIndent() int {
	return len(line.Indent) + len(line.Key) + len(line.Spacing)
}

func (line Line) String() string {
	return string(bytes.Join([][]byte{line.Indent, line.Key, line.Spacing, line.Value, line.Trailing}, nil))
}

func (line Line) GoString() string {
	return string(bytes.Join([][]byte{line.Indent, line.Key, line.Spacing, line.Value, line.Trailing}, []byte("~")))
}

type Scanner struct {
	lineScanner *bufio.Scanner
	prevLine    Line
	path        []string
}

func New(reader io.Reader) *Scanner {
	return &Scanner{
		lineScanner: bufio.NewScanner(reader),
	}
}

func (s *Scanner) Line() Line {
	return s.prevLine
}

func (s *Scanner) Path() string {
	return strings.Join(s.path, "/")
}

func (s *Scanner) Err() error {
	return s.lineScanner.Err()
}

func (s *Scanner) Scan() bool {
	if !s.lineScanner.Scan() {
		return false
	}

	b := s.lineScanner.Bytes()

	line := s.parseLine(b)
	if len(line.Key) > 0 {
		s.path = []string{newPathSegment(line.Key)}
	}
	s.prevLine = line
	return true
}

func (s *Scanner) parseLine(b []byte) Line {
	var line Line

	// "  IP:           10.0.0.1"
	//    ^keyIndex
	keyIndex := indexOfNonSpace(b)
	if keyIndex < 0 {
		// No chars on this line. Must be empty line.
		if len(b) > 0 {
			line.Trailing = b
		}
		return line
	}

	// Add the indentation whitespace
	if keyIndex > 0 {
		// "  IP:           10.0.0.1"
		//  ^^
		line.Indent = b[:keyIndex]
	}

	// "  IP:           10.0.0.1"
	//    ^^^^^^^^^^^^^^^^^^^^^^
	leftTrimmed := b[keyIndex:]

	if len(s.prevLine.Value) > 0 && keyIndex == s.prevLine.ValueIndent() {
		// Multiple values, so treat remainder just as value:
		// "Labels:           app.kubernetes.io/instance=traefik-traefik"
		// "                  app.kubernetes.io/managed-by=Helm"
		// "                  app.kubernetes.io/name=traefik"
		//                    ^lastValueIndent
		line.Value = leftTrimmed
		return line
	}

	// "IP:           10.0.0.1"
	//     ^endOfKey
	// Looking for double space, as some keys have spaces in them, e.g:
	// "QoS Class:                   Burstable"
	//            ^endOfKey
	endOfKey := bytes.Index(leftTrimmed, doubleSpace)
	if endOfKey < 0 {
		// No end of key, so there's no value here

		if len(s.prevLine.Key) > 0 && len(s.prevLine.Value) == 0 && keyIndex > s.prevLine.KeyIndent() {
			// "Args:"
			// "  --this-flag"
			//    ^keyIndex
			line.Value = leftTrimmed
			return line
		}
		if len(s.prevLine.Key) == 0 && len(s.prevLine.Value) > 0 && keyIndex == s.prevLine.ValueIndent() {
			// Previous was array element, so keep being array element
			// "Args:"
			// "  --first-flag"
			// "  --this-flag"
			//    ^keyIndex
			line.Value = leftTrimmed
			return line
		}

		line.Key = leftTrimmed
		return line
	}

	// "IP:           10.0.0.1"
	//  ^^^
	key := leftTrimmed[:endOfKey]
	line.Key = key

	// "IP:           10.0.0.1"
	//     ^^^^^^^^^^^^^^^^^^^
	pastKey := leftTrimmed[endOfKey:]

	// "IP:           10.0.0.1"
	//                ^valueIndex
	valueIndex := indexOfNonSpace(pastKey)
	if valueIndex < 0 {
		// Maybe just some trailing whitespace on the line
		// "data:  " => "  "
		line.Trailing = pastKey
		return line
	}

	// "           10.0.0.1"
	//  ^^^^^^^^^^^
	line.Spacing = pastKey[:valueIndex]
	// "           10.0.0.1"
	//             ^^^^^^^^
	line.Value = pastKey[valueIndex:]
	return line
}

func indexOfNonSpace(b []byte) int {
	for i := 0; i < len(b); i++ {
		if b[i] != ' ' {
			return i
		}
	}
	return -1
}

func newPathSegment(key []byte) string {
	if len(key) == 0 {
		return ""
	}
	if key[len(key)-1] == ':' {
		return string(key[:len(key)-1])
	}
	return string(key)
}
