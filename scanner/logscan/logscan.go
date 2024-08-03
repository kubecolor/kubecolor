package logscan

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/kubecolor/kubecolor/internal/bytesutil"
)

type Token struct {
	Kind Kind
	Text string
}

type Kind byte

const (
	KindUnknown Kind = iota

	KindNewline      // end of the line
	KindPreformatted // e.g "\e[33mAlready with some colors\e[0m"

	KindDate        // e.g "2024-08-03T12:38:44.049832713Z"
	KindSourceRef   // e.g "reconciler.go:142]" or "[main.py:10]"
	KindQuote       // double-quoted or single-quoted string, e.g `"Updated object"`
	KindParenthases // e.g "(" + some other token + ")"

	KindKey   // e.g `controller` in `controller="apiservice"`
	KindValue // e.g `"apiservice"` in `controller="apiservice"`

	KindSeverityTrace
	KindSeverityDebug
	KindSeverityInfo
	KindSeverityWarn
	KindSeverityError
	KindSeverityFatal
	KindSeverityPanic
)

type Scanner struct {
	tokenIndex  int
	tokenBuffer []Token
	lineBuffer  []byte
	lineScanner *bufio.Scanner

	newlineBeforeScan bool
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		lineScanner: bufio.NewScanner(reader),
		tokenBuffer: make([]Token, 0, 10),
	}
}

func (s *Scanner) Token() Token {
	if s.tokenIndex < 0 || s.tokenIndex >= len(s.tokenBuffer) {
		return Token{}
	}
	return s.tokenBuffer[s.tokenIndex]
}

func (s *Scanner) Err() error {
	return s.lineScanner.Err()
}

func (s *Scanner) Scan() bool {
	for {
		if !s.tryScan() {
			return false
		}
		if len(s.tokenBuffer) > 0 {
			return true
		}
	}
}

func (s *Scanner) tryScan() bool {
	if s.tokenIndex < len(s.tokenBuffer)-1 {
		s.tokenIndex++
		return true
	}
	s.tokenIndex = 0
	s.tokenBuffer = s.tokenBuffer[:0]
	if len(s.lineBuffer) == 0 {
		if s.newlineBeforeScan {
			s.pushToken(KindNewline, "\n")
			s.newlineBeforeScan = false
			return true
		}
		if !s.lineScanner.Scan() {
			return false
		}
		s.newlineBeforeScan = true
		s.lineBuffer = s.lineScanner.Bytes()

		if bytes.Contains(s.lineBuffer, []byte("\033[")) {
			s.pushToken(KindPreformatted, string(s.lineBuffer))
			s.lineBuffer = nil
			return true
		}
	}

	length := s.scan(s.lineBuffer)
	s.lineBuffer = s.lineBuffer[length:]
	return true
}

func (s *Scanner) scan(rest []byte) int {
	word := readWord(rest)
	if len(word) == 0 {
		// just return the rest as-is
		return len(rest)
	}

	firstRune, _ := utf8.DecodeRune(word)
	switch firstRune {
	case '[':
		group := readParenthases(rest, '[', ']')
		if group != nil {
			s.scanParenthases(group)
			return len(group)
		}
	case '(':
		group := readParenthases(rest, '(', ')')
		if group != nil {
			s.scanParenthases(group)
			return len(group)
		}
	case '"', '\'', '`':
		if quoted := readQuoted(rest); len(quoted) != 0 {
			return s.pushToken(KindQuote, string(quoted))
		}
	case '{':
		// TODO: Try read JSON
	}

	if key, _, ok := bytes.Cut(word, []byte("=")); ok {
		return s.scanKeyValue(key, rest[len(key)+1:]) // +1 to skip the "=" sign
	}

	// Kubernetes "klog" style source mapping
	if word[len(word)-1] == ']' && bytes.ContainsRune(word, '.') && bytes.ContainsRune(word, ':') {
		return s.pushToken(KindSourceRef, string(word))
	}

	severity := bytes.TrimRight(word, ":")
	if bytesutil.IsOnlyLetters(severity) {
		switch string(severity) {
		case "TRACE", "TRC", "trace", "trc":
			return s.pushToken(KindSeverityTrace, string(severity))
		case "DEBUG", "DBG", "debug", "dbg":
			return s.pushToken(KindSeverityDebug, string(severity))
		case "INFORMATION", "INFO", "INF", "info", "inf":
			return s.pushToken(KindSeverityInfo, string(severity))
		case "WARNING", "WARN", "WRN", "warning", "warn", "wrn":
			return s.pushToken(KindSeverityWarn, string(severity))
		case "ERROR", "ERRO", "ERR", "error", "erro", "err":
			return s.pushToken(KindSeverityError, string(severity))
		case "FATAL", "fatal":
			return s.pushToken(KindSeverityFatal, string(severity))
		case "PANIC", "panic":
			return s.pushToken(KindSeverityPanic, string(severity))
		}
	}

	return s.pushToken(KindUnknown, string(word))
}

func (s *Scanner) pushToken(kind Kind, text string) int {
	s.tokenBuffer = append(s.tokenBuffer, Token{Kind: kind, Text: text})
	return len(text)
}

func (s *Scanner) scanParenthases(group []byte) {
	s.pushToken(KindParenthases, string(group[:1]))
	inner := group[1 : len(group)-1]
	s.scan(inner)
	s.pushToken(KindParenthases, string(group[len(group)-1:]))
}

func (s *Scanner) scanKeyValue(key, valueAndRest []byte) int {
	s.pushToken(KindKey, string(key))
	s.pushToken(KindUnknown, "=")

	word := readWord(valueAndRest)

	firstRune, _ := utf8.DecodeRune(word)
	switch firstRune {
	case '(':
		if group := readParenthases(valueAndRest, '(', ')'); len(group) != 0 {
			s.pushToken(KindValue, string(group))
			return len(key) + 1 + len(group)
		}
	case '[':
		if group := readParenthases(valueAndRest, '[', ']'); len(group) != 0 {
			s.pushToken(KindValue, string(group))
			return len(key) + 1 + len(group)
		}
	case '"', '\'', '`':
		if quoted := readQuoted(valueAndRest); len(quoted) != 0 {
			s.pushToken(KindValue, string(quoted))
			return len(key) + 1 + len(quoted)
		}

	case '{':
		// TODO: parse JSON

	}
	s.pushToken(KindValue, string(word))
	return len(key) + 1 + len(word)
}

func readParenthases(lineBuffer []byte, open, close rune) []byte {
	var openCount int = -1 // start at -1 because the first symbol is the open
	var index int
	for {
		nextRune, size := utf8.DecodeRune(lineBuffer[index:])
		if nextRune == utf8.RuneError {
			return nil
		}
		index += size

		switch nextRune {
		case open:
			openCount++
		case close:
			if openCount == 0 {
				return lineBuffer[:index]
			}
			openCount--
		}
	}
}

func readQuoted(lineBuffer []byte) []byte {
	firstRune, size := utf8.DecodeRune(lineBuffer)
	if firstRune == utf8.RuneError {
		return nil
	}
	index := size
	for {
		nextRune, size := utf8.DecodeRune(lineBuffer[index:])
		if nextRune == utf8.RuneError {
			return nil
		}
		index += size
		if nextRune == '\\' {
			index++ // skip 1 more
			continue
		}
		if nextRune == firstRune {
			return lineBuffer[:index]
		}
	}
}

func readWord(lineBuffer []byte) []byte {
	firstRune, size := utf8.DecodeRune(lineBuffer)
	if firstRune == utf8.RuneError {
		// Invalid utf8. Just return rest of the line as-is
		return lineBuffer
	}

	if unicode.IsSpace(firstRune) {
		// Read spaces
		index := size
		for {
			nextRune, size := utf8.DecodeRune(lineBuffer[index:])
			if nextRune == utf8.RuneError || !unicode.IsSpace(nextRune) {
				return lineBuffer[:index]
			}
			index += size
		}
	}

	// read non-spaces
	index := size
	for {
		nextRune, size := utf8.DecodeRune(lineBuffer[index:])
		if nextRune == utf8.RuneError || unicode.IsSpace(nextRune) {
			return lineBuffer[:index]
		}
		index += size
	}
}
