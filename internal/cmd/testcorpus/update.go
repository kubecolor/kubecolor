package main

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	testHeaderSeparator = strings.Repeat("=", 80)
	testOutputSeparator = strings.Repeat("-", 80)
)

func UpdateTests(files []File) {
	anyErr := false
	for _, file := range files {
		fmt.Printf("  %s:\n", colorHeader.Render(file.Name))
		if len(file.Tests) == 0 {
			fmt.Printf("    %s\n", colorMuted.Render("no tests found"))
			continue
		}
		changed, err := updateFile(file)
		if err != nil {
			fmt.Printf("    %s %s\n", colorErrorPrefix.Render("error:"), colorErrorText.Render(err.Error()))
			anyErr = true
		} else if changed {
			fmt.Printf("    %s\n", colorSuccess.Render("updated"))
		} else {
			fmt.Printf("    %s\n", colorMuted.Render("unchanged"))
		}
	}
	if anyErr {
		os.Exit(1)
	}
}

func updateFile(file File) (bool, error) {
	var buf bytes.Buffer
	beforeHash := hashFile(file.Path)

	for i, test := range file.Tests {
		if i > 0 {
			buf.WriteByte('\n')
		}
		if err := writeTest(&buf, test); err != nil {
			return false, fmt.Errorf("test %s: %w", test.Name, err)
		}
	}

	afterHash := hashReader(bytes.NewReader(buf.Bytes()))
	if beforeHash == afterHash {
		return false, nil
	}

	if err := os.WriteFile(file.Path, buf.Bytes(), 0644); err != nil {
		return false, err
	}
	return true, nil
}

func hashFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	return hashReader(f)
}

func hashReader(r io.Reader) string {
	h := sha512.New()
	io.Copy(h, r)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func writeTest(w io.Writer, test Test) error {
	args := strings.Fields(test.Command)
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}
	cmd := args[0]
	args = args[1:]

	if cmd != "kubectl" {
		return fmt.Errorf(`command must start with "kubectl", but got %q`, cmd)
	}

	fmt.Fprintln(w, testHeaderSeparator)
	if test.Name != test.Command {
		fmt.Fprintf(w, "# %s\n", test.Name)
	}
	fmt.Fprintf(w, "$ %s\n", test.Command)
	fmt.Fprintln(w, testHeaderSeparator)
	fmt.Fprintln(w)
	fmt.Fprintln(w, test.Input)
	fmt.Fprintln(w)
	fmt.Fprintln(w, testOutputSeparator)
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.TrimSpace(printCommand(args, test.Input)))
	return nil
}
