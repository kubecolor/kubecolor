package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode/utf8"
)

func NewTestScanner(reader io.Reader) *TestScanner {
	return &TestScanner{
		scanner: bufio.NewScanner(reader),
		first:   true,
	}
}

type TestScanner struct {
	scanner *bufio.Scanner

	test  Test
	err   error
	done  bool
	first bool
}

func (s *TestScanner) Scan() bool {
	s.err = s.scanErr()
	return s.err == nil && !s.done
}

func (s *TestScanner) Test() Test {
	return s.test
}

func (s *TestScanner) Err() error {
	return s.err
}

func (s *TestScanner) scanErr() error {
	if s.first {
		s.first = false
		if _, err := scanLinesUntil(s.scanner, "==="); errors.Is(err, io.ErrUnexpectedEOF) {
			// no tests found
			s.done = true
			return nil
		} else if err != nil {
			return err
		}
	}

	s.test = Test{}
	if err := s.scanHeader(); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			s.done = true
			return nil
		}
		return err
	}
	if err := s.scanInput(); err != nil {
		return err
	}
	if err := s.scanOutput(); err != nil {
		return err
	}
	return nil
}

func (s *TestScanner) scanHeader() error {
	headerLines, err := scanLinesUntil(s.scanner, "===")
	if err != nil {
		return fmt.Errorf("scan until closing === line: %w", err)
	}
	for _, line := range headerLines {
		if err := s.readHeaderLine(line); err != nil {
			return err
		}
	}
	return nil
}

func (s *TestScanner) readHeaderLine(line string) error {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return nil
	}
	firstChar, size := utf8.DecodeRuneInString(line)
	rest := strings.TrimSpace(trimmed[size:])
	switch firstChar {
	case '#':
		s.test.Name = rest
		return nil
	case '$':
		if s.test.Command != "" {
			return fmt.Errorf("test %q: found multiple command lines (starting with '$') in header", s.test.Name)
		}
		s.test.Command = rest
		if s.test.Name == "" {
			s.test.Name = rest
		}
		return nil
	}

	if key, value, ok := strings.Cut(trimmed, "="); ok {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			if _, err := fmt.Sscanf(value, "%q", &value); err != nil {
				return fmt.Errorf("test %q: env %q: failed to parse value as quoted string: %w", s.test.Name, key, err)
			}
		}
		if slices.ContainsFunc(s.test.Env, func(e EnvVar) bool {
			return e.Key == key
		}) {
			return fmt.Errorf("test %q: env %q: key specified multiple times", s.test.Name, key)
		}

		s.test.Env = append(s.test.Env, EnvVar{
			Key:   key,
			Value: value,
		})
		return nil
	}

	return fmt.Errorf("test %q: invalid test header line %q", s.test.Name, trimmed)
}

func (s *TestScanner) scanInput() error {
	inputLines, err := scanLinesUntil(s.scanner, "---")
	if err != nil {
		return fmt.Errorf("scan test input: %w", err)
	}
	s.test.Input = strings.TrimSpace(strings.Join(inputLines, "\n"))
	return nil
}

func (s *TestScanner) scanOutput() error {
	outputLines, err := scanLinesUntil(s.scanner, "===")
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("scan test output: %w", err)
	}
	s.test.Output = strings.TrimSpace(strings.Join(outputLines, "\n"))
	return nil
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
