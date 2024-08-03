package logscan

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/kubecolor/kubecolor/internal/bytesutil"
)

// klogLevelAndDateRegex is for parsing Kubernetes klog line: https://github.com/kubernetes/klog/blob/75663bb798999a49e3e4c0f2375ed5cca8164194/klog.go#L637-L650
//
//	Lmmdd hh:mm:ss.uuuuuu threadid file:line] msg...
var klogLevelAndDateRegex = regexp.MustCompile(`^([IWEF])(\d{4} \d\d:\d\d:\d\d\.\d+)(\s*\d+\s*)([\w\._]+:\d+)\]`)

// dateRegex is for parsing dates. E.g:
//
//	2024-08-03T19:57:19.446242
//	2024-08-03 20:04:28.614 GMT
var dateRegex = regexp.MustCompile(`^\d{4}-\d\d-\d\dT\d\d:\d\d(:\d\d(\.\d+)?)?(Z|\+\d\d:\d\d|\+\d{4})?\b|^(\d{4}-\d\d-\d\d|\d\d ([a-zA-Z][a-z]+) \d{4}|\d\d/([a-zA-Z][a-z]+)/\d{4})[ :]\d\d:\d\d(:\d\d(\.\d+)?)?( ?(GMT|UTC|\+\d\d:\d\d|\+\d\d\d\d))?\b`)

// guidRegex is for matching on GUIDs and UUIDs. E.g:
//
//	70d5707e-b07b-41c3-9411-cad84c6db764
//	70d5707eb07b41c39411cad84c6db764
var guidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$|^[0-9a-fA-F]{32}$`)

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
	KindGUID        // e.g "70d5707e-b07b-41c3-9411-cad84c6db764"
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
	hasFoundSeverity  bool
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
		s.hasFoundSeverity = false
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

	// Kubernetes "klog" style source mapping, e.g "dynamic_source.go:290]"
	if word[len(word)-1] == ']' && bytes.ContainsRune(word, '.') && bytes.ContainsRune(word, ':') {
		return s.pushToken(KindSourceRef, string(word))
	}

	// Kubernetes "klog" header, which follows the format: https://github.com/kubernetes/klog/blob/75663bb798999a49e3e4c0f2375ed5cca8164194/klog.go#L637-L650
	//
	//	Lmmdd hh:mm:ss.uuuuuu threadid file:line] msg...
	if klogMatches := klogLevelAndDateRegex.FindSubmatch(rest); klogMatches != nil {
		fullMatch := klogMatches[0]
		severity := klogMatches[1]
		date := klogMatches[2]
		threadIDWithPadding := klogMatches[3]
		sourceRef := klogMatches[4]

		var severityKind Kind
		switch firstRune {
		case 'I':
			severityKind = KindSeverityInfo
		case 'W':
			severityKind = KindSeverityWarn
		case 'E':
			severityKind = KindSeverityError
		case 'F':
			severityKind = KindSeverityFatal
		}

		s.hasFoundSeverity = true
		s.pushToken(severityKind, string(severity))
		s.pushToken(KindDate, string(date))
		s.pushToken(KindUnknown, string(threadIDWithPadding))
		s.pushToken(KindSourceRef, string(sourceRef))
		s.pushToken(KindParenthases, "]")
		return len(fullMatch)
	}

	if dateMatch := dateRegex.Find(rest); dateMatch != nil {
		return s.pushToken(KindDate, string(dateMatch))
	}

	if guidRegex.Match(word) {
		return s.pushToken(KindGUID, string(word))
	}

	if !s.hasFoundSeverity {
		severity := bytes.TrimRight(word, ":!,")
		if bytesutil.IsOnlyLetters(severity) {
			severityKind := severityKindFromName(string(severity))
			if severityKind != KindUnknown {
				s.hasFoundSeverity = true
				return s.pushToken(severityKind, string(severity))
			}
		}
	}

	return s.pushToken(KindUnknown, string(word))
}

func severityKindFromName(severity string) Kind {
	switch severity {
	case "TRACE", "TRC",
		"Trace", "Trc",
		"trace", "trc":
		return KindSeverityTrace
	case "DEBUG", "DBG",
		"Debug", "Dbg",
		"debug", "dbg":
		return KindSeverityDebug
	case "INFORMATION", "INFO", "INF",
		"Information", "Info", "Inf",
		"information", "info", "inf",
		"NOTE", "Note", "note",
		"SUCCESSFULLY", "Successfully", "successfully",
		"SUCCESS", "Success", "success":
		return KindSeverityInfo
	case "WARNING", "WARN", "WRN",
		"Warning", "Warn", "Wrn",
		"warning", "warn", "wrn":
		return KindSeverityWarn
	case "ERROR", "ERRO", "ERR",
		"Error", "Erro", "Err",
		"error", "erro", "err",
		"FAILED", "Failed", "failed":
		return KindSeverityError
	case "FATAL", "Fatal", "fatal":
		return KindSeverityFatal
	case "PANIC", "Panic", "panic":
		return KindSeverityPanic
	default:
		return KindUnknown
	}
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
		if group := readParenthases(valueAndRest, '(', ')'); len(group) != 0 && len(group) >= len(word) {
			s.pushToken(KindValue, string(group))
			return len(key) + 1 + len(group)
		}
	case '[':
		if group := readParenthases(valueAndRest, '[', ']'); len(group) != 0 && len(group) >= len(word) {
			s.pushToken(KindValue, string(group))
			return len(key) + 1 + len(group)
		}
	case '"', '\'', '`':
		if quoted := readQuoted(valueAndRest); len(quoted) != 0 && len(quoted) >= len(word) {
			s.pushToken(KindValue, string(quoted))
			return len(key) + 1 + len(quoted)
		}

	case '{':
		// TODO: parse JSON
	}

	if dateMatch := dateRegex.Find(word); dateMatch != nil {
		s.pushToken(KindDate, string(word))
		return len(key) + 1 + len(word)
	}

	switch string(key) {
	case "level", "lvl", "severity", "l", "s":
		severityKind := severityKindFromName(string(word))
		if severityKind != KindUnknown {
			s.hasFoundSeverity = true
			s.pushToken(severityKind, string(word))
			return len(key) + 1 + len(word)
		}

	case "caller", "source":
		s.pushToken(KindSourceRef, string(word))
		return len(key) + 1 + len(word)
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
