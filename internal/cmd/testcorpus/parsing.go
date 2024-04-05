package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

func ParseGlobFS(fsys fs.FS, glob string) ([]File, error) {
	matches, err := fs.Glob(fsys, glob)
	if err != nil {
		return nil, err
	}

	var files []File
	for _, path := range matches {
		stat, err := fs.Stat(fsys, path)
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			continue
		}
		file, err := ParseFileFS(fsys, path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		files = append(files, file)
	}
	return files, nil
}

func ParseGlob(glob string) ([]File, error) {
	fsys := os.DirFS(".")
	return ParseGlobFS(fsys, glob)
}

func ParseFileFS(fsys fs.FS, path string) (File, error) {
	file, err := fsys.Open(path)
	if err != nil {
		return File{}, err
	}
	defer file.Close()
	testScanner := NewTestScanner(file)

	var tests []Test
	for testScanner.Scan() {
		tests = append(tests, testScanner.Test())
	}
	if err := testScanner.Err(); err != nil {
		return File{}, err
	}

	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return File{
		Name:  name,
		Path:  path,
		Tests: tests,
	}, nil
}
