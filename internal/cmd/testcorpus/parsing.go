package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type File struct {
	Name  string
	Path  string
	Tests []Test
}

type Test struct {
	Name    string
	Command string
	Input   string
	Output  string
}

func ParseGlob(glob string) ([]File, error) {
	fsys := os.DirFS(".")
	matches, err := fs.Glob(fsys, glob)
	if err != nil {
		return nil, err
	}

	var files []File
	for _, path := range matches {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			continue
		}
		base := filepath.Base(path)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		tests, err := ParseFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		files = append(files, File{
			Name:  name,
			Path:  path,
			Tests: tests,
		})
	}
	return files, nil
}

func ParseFile(path string) ([]Test, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if _, err := scanLinesUntil(scanner, "==="); errors.Is(err, io.ErrUnexpectedEOF) {
		// no tests found
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var tests []Test

	for {
		var test Test
		headerLines, err := scanLinesUntil(scanner, "===")
		if err != nil {
			return nil, fmt.Errorf("scan until closing === line: %w", err)
		}
		for _, line := range headerLines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			firstChar, size := utf8.DecodeRuneInString(line)
			rest := strings.TrimSpace(line[size:])
			switch firstChar {
			case '#':
				test.Name = rest
			case '$':
				test.Command = rest
				if test.Name == "" {
					test.Name = rest
				}
			default:
				return nil, fmt.Errorf("test %q: invalid test header first char %q", test.Name, firstChar)
			}
		}

		inputLines, err := scanLinesUntil(scanner, "---")
		if err != nil {
			return nil, fmt.Errorf("scan test input: %w", err)
		}
		test.Input = strings.TrimSpace(strings.Join(inputLines, "\n"))

		outputLines, err := scanLinesUntil(scanner, "===")
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("scan test output: %w", err)
		}
		test.Output = strings.TrimSpace(strings.Join(outputLines, "\n"))
		tests = append(tests, test)

		if errors.Is(err, io.ErrUnexpectedEOF) {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tests, nil
}

func scanLinesUntil(scanner *bufio.Scanner, linePrefix string) ([]string, error) {
	var lines []string

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\r")

		if strings.HasPrefix(line, linePrefix) {
			return lines, nil
		}

		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, io.ErrUnexpectedEOF
}
