package main

import (
	"os"
	"testing"
)

func TestCorpus(t *testing.T) {
	const glob = "test/corpus/*.txt"
	repoRoot := os.DirFS("../../..")

	files, err := ParseGlobFS(repoRoot, glob)
	if err != nil {
		t.Fatalf("Error parsing corpus test file glob: %s", err)
	}

	if len(files) == 0 {
		cwd, _ := os.Getwd()
		t.Logf("Current directory: %q", cwd)
		t.Fatalf("Glob did not match any files: %s", glob)
	}

	for _, file := range files {
		t.Run(file.Name, func(t *testing.T) {
			if len(file.Tests) == 0 {
				t.Fatal("no tests found in file")
			}
			for _, test := range file.Tests {
				t.Run(test.Name, func(t *testing.T) {
					if err := ExecuteTest(test); err != nil {
						t.Error(FormatTestError(test, err))
					}
				})
			}
		})
	}
}
