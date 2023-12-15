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

type PathSegment struct {
	Segment   string
	KeyIndent int
}

func newPathSegment(line Line) PathSegment {
	if len(line.Key) == 0 {
		return PathSegment{}
	}

	// Converting it to string here does copy the slice, but this is intended.
	// The [[]byte] returned by [bufio.Scanner.Bytes] references a mutating
	// slices, which means if we keep the value for multiple [bufio.Scanner.Scan]
	// calls, then the value might get corrupted.

	if cut, ok := bytes.CutSuffix(line.Key, []byte{':'}); ok {
		return PathSegment{Segment: string(cut), KeyIndent: line.KeyIndent()}
	}
	return PathSegment{Segment: string(line.Key), KeyIndent: line.KeyIndent()}
}

type Scanner struct {
	lineScanner  *bufio.Scanner
	prevLine     Line
	pathSegments []PathSegment
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		lineScanner: bufio.NewScanner(reader),
	}
}

func (s *Scanner) Line() Line {
	return s.prevLine
}

func (s *Scanner) Path() string {
	if len(s.pathSegments) == 0 {
		return ""
	}
	if len(s.pathSegments) == 1 {
		return s.pathSegments[0].Segment
	}
	var sb strings.Builder
	for i, p := range s.pathSegments {
		if i > 0 {
			sb.WriteByte('/')
		}
		sb.WriteString(p.Segment)
	}
	return sb.String()
}

func (s *Scanner) PathSegments() []PathSegment {
	return s.pathSegments
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
		segment := newPathSegment(line)
		for len(s.pathSegments) > 0 && s.pathSegments[len(s.pathSegments)-1].KeyIndent >= segment.KeyIndent {
			s.pathSegments = s.pathSegments[:len(s.pathSegments)-1]
		}
		s.pathSegments = append(s.pathSegments, segment)
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
	endOfKey := indexOfSpace(leftTrimmed)
	if endOfKey < 0 {
		// No end of key, so there's no value here

		if leftTrimmed[len(leftTrimmed)-1] == ':' {
			// Ending with ":" always means it's a key
			line.Key = leftTrimmed
			return line
		}

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
		if b[i] != ' ' && b[i] != '\t' {
			return i
		}
	}
	return -1
}

func indexOfSpace(b []byte) int {
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
