package describe

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

var doubleSpace = []byte{' ', ' '}
var lineFeed = []byte{'\n'}
var lineFeedToken = Token{Kind: KindEOL, Bytes: lineFeed}

type Kind byte

const (
	KindUnknown Kind = iota
	KindKey
	KindValue
	KindWhitespace
	KindEOL
)

func (k Kind) String() string {
	switch k {
	case KindUnknown:
		return "unknown"
	case KindKey:
		return "key"
	case KindValue:
		return "value"
	case KindWhitespace:
		return "whitespace"
	case KindEOL:
		return "eol"
	default:
		return fmt.Sprintf("%[1]T(%[1]d)", k)
	}
}

type Token struct {
	Bytes []byte
	Kind  Kind
}

type State struct {
	KeyIndent   int
	ValueIndent int
}

type Scanner struct {
	lineScanner *bufio.Scanner
	tokens      []Token
	tokenIndex  int
	state       State
	lastState   State
}

func New(reader io.Reader) *Scanner {
	return &Scanner{
		lineScanner: bufio.NewScanner(reader),

		// Set fields to intentionally wrong values
		state:     State{KeyIndent: -1, ValueIndent: -1},
		lastState: State{KeyIndent: -1, ValueIndent: -1},
	}
}

func (s *Scanner) Token() Token {
	if s.tokenIndex >= len(s.tokens) {
		return Token{}
	}
	return s.tokens[s.tokenIndex]
}

func (s *Scanner) State() State {
	return s.state
}

func (s *Scanner) Err() error {
	return s.lineScanner.Err()
}

func (s *Scanner) Scan() bool {
	if s.tokenIndex < len(s.tokens)-1 {
		s.tokenIndex++
		return true
	}
	if !s.lineScanner.Scan() {
		return false
	}

	clear(s.tokens) // let byte slices get GC'd
	s.tokens = s.tokens[:0]
	s.tokenIndex = 0
	s.lastState = s.state
	s.state = State{}

	b := s.lineScanner.Bytes()

	// "  IP:           10.0.0.1"
	//    ^keyIndex
	keyIndex := indexOfNonSpace(b)
	if keyIndex < 0 {
		// No chars on this line. Must be empty line.
		if len(b) > 0 {
			s.tokens = append(s.tokens, Token{
				Kind:  KindWhitespace,
				Bytes: b,
			})
		}
		s.tokens = append(s.tokens, lineFeedToken)
		return true
	}

	// Add the indentation whitespace
	if keyIndex > 0 {
		// "  IP:           10.0.0.1"
		//  ^^
		s.state.KeyIndent = keyIndex
		s.tokens = append(s.tokens, Token{
			Kind:  KindWhitespace,
			Bytes: b[:keyIndex],
		})
	}

	// "  IP:           10.0.0.1"
	//    ^^^^^^^^^^^^^^^^^^^^^^
	leftTrimmed := b[keyIndex:]

	if keyIndex == s.lastState.ValueIndent {
		// Multiple values, so treat remainder just as value:
		// "Labels:           app.kubernetes.io/instance=traefik-traefik"
		// "                  app.kubernetes.io/managed-by=Helm"
		// "                  app.kubernetes.io/name=traefik"
		//                    ^lastValueIndent
		s.state = s.lastState
		s.tokens = append(s.tokens, Token{
			Kind:  KindValue,
			Bytes: leftTrimmed,
		})
		s.tokens = append(s.tokens, lineFeedToken)
		return true
	}

	// "IP:           10.0.0.1"
	//     ^endOfKey
	// Looking for double space, as some keys have spaces in them, e.g:
	// "QoS Class:                   Burstable"
	//            ^endOfKey
	endOfKey := bytes.Index(leftTrimmed, doubleSpace)
	if endOfKey < 0 {
		// No end of key, so there's no value here
		s.tokens = append(s.tokens, Token{
			Kind:  KindKey,
			Bytes: leftTrimmed,
		})
		s.tokens = append(s.tokens, lineFeedToken)
		return true
	}

	// "IP:           10.0.0.1"
	//  ^^^
	key := leftTrimmed[:endOfKey]
	s.tokens = append(s.tokens, Token{
		Kind:  KindKey,
		Bytes: key,
	})

	// "IP:           10.0.0.1"
	//     ^^^^^^^^^^^^^^^^^^^
	trailing := leftTrimmed[endOfKey:]

	// "IP:           10.0.0.1"
	//                ^valueIndex
	valueIndex := indexOfNonSpace(trailing)
	if valueIndex < 0 {
		// Maybe just some trailing whitespace on the line
		// "data:  " => "  "
		s.tokens = append(s.tokens, Token{
			Kind:  KindWhitespace,
			Bytes: trailing,
		})
		s.tokens = append(s.tokens, lineFeedToken)
		return true
	}
	s.state.ValueIndent = valueIndex + endOfKey + keyIndex

	// "           10.0.0.1"
	//  ^^^^^^^^^^^
	s.tokens = append(s.tokens, Token{
		Kind:  KindWhitespace,
		Bytes: trailing[:valueIndex],
	})
	// "           10.0.0.1"
	//             ^^^^^^^^
	s.tokens = append(s.tokens, Token{
		Kind:  KindValue,
		Bytes: trailing[valueIndex:],
	})

	s.tokens = append(s.tokens, lineFeedToken)
	return true
}

func indexOfNonSpace(b []byte) int {
	for i := 0; i < len(b); i++ {
		if b[i] != ' ' {
			return i
		}
	}
	return -1
}
