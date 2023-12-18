package tablescan

import (
	"bufio"
	"io"
	"strings"

	"github.com/kubecolor/kubecolor/internal/bytesutil"
)

type Cell struct {
	Trimmed        string
	Full           string
	TrailingSpaces string
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
	lineScanner   *bufio.Scanner
	headerIndices []int
	currentCells  []Cell
	leadingSpaces []byte
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		lineScanner: bufio.NewScanner(reader),
	}
}

func (s *Scanner) Cells() []Cell {
	return s.currentCells
}

func (s *Scanner) LeadingSpaces() []byte {
	return s.leadingSpaces
}

func (s *Scanner) Bytes() []byte {
	return s.lineScanner.Bytes()
}

func (s *Scanner) Text() string {
	return s.lineScanner.Text()
}

func (s *Scanner) Err() error {
	return s.lineScanner.Err()
}

func (s *Scanner) Scan() bool {
	if !s.lineScanner.Scan() {
		return false
	}

	b := s.lineScanner.Bytes()

	clear(s.currentCells) // let strings get GC'd
	s.currentCells = s.currentCells[:0]

	if len(b) == 0 {
		// Empty line. Should reset header calculations then
		s.headerIndices = nil
		return true
	} else if s.shouldCalcIndices(b) {
		s.headerIndices = calcHeaderIndices(b)
	}

	if len(s.headerIndices) > 0 {
		s.leadingSpaces = b[:s.headerIndices[0]]
	}

	for i, columnIndex := range s.headerIndices {
		if columnIndex >= len(b) {
			// empty cell at end of line
			s.currentCells = append(s.currentCells, emptyCell)
			continue
		}
		var str string
		if i+1 < len(s.headerIndices) {
			// there's more lines
			nextColumnIndex := s.headerIndices[i+1]
			str = string(b[columnIndex:min(nextColumnIndex, len(b))])
		} else {
			// last column
			str = string(b[columnIndex:])
		}
		s.currentCells = append(s.currentCells, NewCell(str))
	}

	return true
}

func (s *Scanner) shouldCalcIndices(line []byte) bool {
	if s.headerIndices == nil {
		return true
	}

	// Checks if we have come across a new table,
	// via lossy logic on the heuristic that cells should have spacing.
	// There will be edge cases where this fails, but it'll be good enough.
	for _, columnIndex := range s.headerIndices[1:] {
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
