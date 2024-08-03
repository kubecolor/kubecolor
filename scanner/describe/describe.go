package describe

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/internal/bytesutil"
)

var (
	spaceCharset = " \t"
	doubleSpace  = []byte{' ', ' '}
	tabBytes     = []byte{'\t'}
)

type Line struct {
	Indent   []byte
	Key      []byte
	Spacing  []byte
	Value    []byte
	Trailing []byte
}

func (line Line) IsZero() bool {
	return len(line.Indent) == 0 &&
		len(line.Key) == 0 &&
		len(line.Spacing) == 0 &&
		len(line.Value) == 0 &&
		len(line.Trailing) == 0
}

func (line Line) KeyIndent() int {
	// Treat tabs as 8 characters long
	return len(line.Indent) + bytes.Count(line.Indent, tabBytes)*7
}

func (line Line) ValueIndent() int {
	return line.KeyIndent() + len(line.Key) + len(line.Spacing)
}

func (line Line) String() string {
	return string(bytes.Join([][]byte{line.Indent, line.Key, line.Spacing, line.Value, line.Trailing}, nil))
}

func (line Line) GoString() string {
	return string(bytes.Join([][]byte{line.Indent, line.Key, line.Spacing, line.Value, line.Trailing}, []byte("~")))
}

type Path []PathSegment

func (p Path) String() string {
	if len(p) == 0 {
		return ""
	}
	if len(p) == 1 {
		return p[0].Segment
	}
	var sb strings.Builder
	for i, p := range p {
		if i > 0 {
			sb.WriteByte('/')
		}
		sb.WriteString(p.Segment)
	}
	return sb.String()
}

func (p Path) HasPrefix(segments ...string) bool {
	for i, s := range segments {
		if i >= len(p) {
			return false
		}
		if p[i].Segment != s {
			return false
		}
	}
	return true
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
	lineScanner *bufio.Scanner
	// currentLine is used in [Scanner.Line]
	currentLine Line
	// prevNonZeroLine is used when code wants to reference previous line
	prevNonZeroLine Line
	pathSegments    []PathSegment
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		lineScanner: bufio.NewScanner(reader),
	}
}

func (s *Scanner) Line() Line {
	return s.currentLine
}

func (s *Scanner) Path() Path {
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
	s.currentLine = line
	if len(line.Key) > 0 || len(line.Value) > 0 {
		s.prevNonZeroLine = line
	}
	return true
}

func (s *Scanner) parseLine(b []byte) Line {
	var line Line

	// "  IP:           10.0.0.1"
	//    ^keyIndex
	keyIndex := bytesutil.IndexOfNonSpace(b, spaceCharset)
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

	if len(s.prevNonZeroLine.Value) > 0 && keyIndex == s.prevNonZeroLine.ValueIndent() {
		// Multiple values, so treat remainder just as value:
		// "Labels:           app.kubernetes.io/instance=traefik-traefik"
		// "                  app.kubernetes.io/managed-by=Helm"
		// "                  app.kubernetes.io/name=traefik"
		//                    ^lastValueIndent
		line.Value = leftTrimmed
		return line
	}

	// "--flag='':"
	//        ^endOfKey
	endOfKey := bytes.IndexRune(leftTrimmed, '=')
	if endOfKey > 0 && leftTrimmed[len(leftTrimmed)-1] == ':' {
		// "--flag='':"
		//  ^^^^^^
		line.Key = leftTrimmed[:endOfKey]
		// "--flag='':"
		//        ^
		line.Spacing = leftTrimmed[endOfKey : endOfKey+1]

		// "--flag='':"
		//         ^^
		line.Value = leftTrimmed[endOfKey+1 : len(leftTrimmed)-1]

		// "--flag='':"
		//           ^
		line.Trailing = leftTrimmed[len(leftTrimmed)-1:]

		return line
	}

	// "IP:           10.0.0.1"
	//     ^endOfKey
	// Looking for double space, as some keys have spaces in them, e.g:
	// "QoS Class:                   Burstable"
	//            ^endOfKey
	endOfKey = bytesutil.IndexOfDoubleSpace(leftTrimmed)
	if endOfKey < 0 {
		// No end of key, so there's no value here

		if leftTrimmed[len(leftTrimmed)-1] == ':' {
			// Ending with ":" always means it's a key.
			// Exception: in `kubectl explain` lots of lines ends with "More info:"
			if !bytes.HasSuffix(leftTrimmed, []byte(" info:")) {
				line.Key = leftTrimmed
				return line
			}
		}

		if len(s.prevNonZeroLine.Key) > 0 && keyIndex > s.prevNonZeroLine.KeyIndent() {
			if len(s.prevNonZeroLine.Value) == 0 {
				// "Args:"
				// "  --this-flag"
				//    ^keyIndex
				line.Value = leftTrimmed
				return line
			}

			lastValueWithoutRequired := bytes.TrimSuffix(s.prevNonZeroLine.Value, []byte(" -required-"))
			if s.prevNonZeroLine.Value[0] == '<' && lastValueWithoutRequired[len(lastValueWithoutRequired)-1] == '>' {
				// Special case for type declarations (e.g in `kubectl explain`)
				// "apiVersion    <string>"
				// "  --this-flag"
				//    ^keyIndex
				line.Value = leftTrimmed
				return line
			}
		}

		if len(s.prevNonZeroLine.Key) == 0 && len(s.prevNonZeroLine.Value) > 0 && keyIndex == s.prevNonZeroLine.ValueIndent() {
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
	valueIndex := bytesutil.IndexOfNonSpace(pastKey, spaceCharset)
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
