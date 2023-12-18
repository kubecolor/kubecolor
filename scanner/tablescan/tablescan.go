package tablescan

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/internal/bytesutil"
)

type Cell struct {
	Trimmed        string
	Full           string
	TrailingSpaces string
}

func (c Cell) String() string {
	return c.Full
}

var emptyCell = NewCell("")

func NewCell(full string) Cell {
	trimmed := strings.TrimRight(full, " \t")
	return Cell{
		Trimmed:        trimmed,
		Full:           full,
		TrailingSpaces: full[len(trimmed):],
	}
}

type Scanner struct {
	lineScanner     *bufio.Scanner
	headerIndices   []int
	currentCells    []Cell
	currentLine string
	leadingSpaces   string
	bufferedLines   []string
	reuseScanResult bool
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		lineScanner: bufio.NewScanner(reader),
	}
}

func (s *Scanner) Cells() []Cell {
	return s.currentCells
}

func (s *Scanner) LeadingSpaces() string {
	return s.leadingSpaces
}

func (s *Scanner) Bytes() []byte {
	return []byte(s.currentLine)
}

func (s *Scanner) Text() string {
	return s.currentLine
}

func (s *Scanner) Err() error {
	return s.lineScanner.Err()
}

func (s *Scanner) Scan() bool {
	if !s.bufferTableLines() {
		return false
	}

	s.currentLine = s.bufferedLines[0]
	s.bufferedLines = s.bufferedLines[1:]

	clear(s.currentCells) // let strings get GC'd
	s.currentCells = s.currentCells[:0]

	if len(s.currentLine) == 0 {
		// Empty line. Should reset header calculations then
		s.bufferedLines = nil
		return true
	}

	if len(s.headerIndices) > 0 {
		s.leadingSpaces = s.currentLine[:s.headerIndices[0]]
	}

	for i, columnIndex := range s.headerIndices {
		if columnIndex >= len(s.currentLine) {
			// empty cell at end of line
			s.currentCells = append(s.currentCells, emptyCell)
			continue
		}
		var str string
		if i+1 < len(s.headerIndices) {
			// there's more lines
			nextColumnIndex := s.headerIndices[i+1]
			str = s.currentLine[columnIndex:min(nextColumnIndex, len(s.currentLine))]
		} else {
			// last column
			str = s.currentLine[columnIndex:]
		}
		s.currentCells = append(s.currentCells, NewCell(str))
	}

	return true
}

func (s *Scanner) bufferTableLines() bool {
	if len(s.bufferedLines) > 0 {
		return true
	}

	s.headerIndices = nil
	var combinedLines []byte

	for {
		if s.reuseScanResult {
			s.reuseScanResult = false
		} else if !s.lineScanner.Scan() {
			return len(s.bufferedLines) > 0
		}

		b := s.lineScanner.Bytes()

		if len(b) == 0 {
			// Empty line. Keep it, but don't process further
			s.bufferedLines = append(s.bufferedLines, "")
			return len(s.bufferedLines) > 0
		}

		combinedLines = bytewiseAndNonSpace(combinedLines, b)
		newHeaderIndices := calcHeaderIndices(combinedLines)

		if isProbablyNewTable(s.headerIndices, combinedLines) {
			s.reuseScanResult = true
			return true
		}

		s.headerIndices = newHeaderIndices

		// Using [bufio.Scanner.Text], as using result from [bufio.Scanner.Bytes]
		// after another scan is undefined behavior,
		// and can lead to corrupted data.
		s.bufferedLines = append(s.bufferedLines, s.lineScanner.Text())
	}
}

func isProbablyNewTable(headerIndices []int, line []byte) bool {
	if len(headerIndices) == 0 {
		return false
	}
	// Checks if we have come across a new table,
	// via lossy logic on the heuristic that cells should have spacing.
	// There will be edge cases where this fails, but it'll be good enough.
	for _, columnIndex := range headerIndices[1:] {
		if columnIndex >= len(line) || columnIndex == 0 {
			continue
		}
		oneBefore := line[max(columnIndex-2, 0)]
		twoBefore := line[max(columnIndex-2, 0)]
		if oneBefore != ' ' || twoBefore != ' ' {
			// should have spacing between columns
			return true
		}
	}

	return false
}

func bytewiseAndNonSpace(oldBytes, newBytes []byte) []byte {
	if len(newBytes) > len(oldBytes) {
		oldBytes = append(oldBytes, bytes.Repeat([]byte(" "), len(newBytes) - len(oldBytes))...)
	}
	for i, b := range newBytes {
		if b == ' ' {
			continue
		}
		oldBytes[i] = b
	}
	return oldBytes
}

func calcHeaderIndices(line []byte) []int {
	var indices []int
	var indexOffset int
	for {
		index := bytesutil.IndexOfNonSpace(line, " \t")
		if index == -1 {
			break
		}
		line = line[index:]
		indexOffset += index
		indices = append(indices, indexOffset)

		index = bytesutil.IndexOfDoubleSpace(line)
		if index == -1 {
			break
		}
		line = line[index:]
		indexOffset += index
	}
	return indices
}
