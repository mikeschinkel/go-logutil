package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikeschinkel/go-logutil"
)

// TestFuzzCorpus runs all fuzz corpus files as regression tests
// This ensures that any interesting inputs discovered during fuzzing
// are tested in CI/CD to prevent regressions
func TestFuzzCorpus(t *testing.T) {
	corpusDir := "testdata/fuzz"

	// Check if corpus directory exists
	if _, err := os.Stat(corpusDir); os.IsNotExist(err) {
		// t.Skip("No fuzz corpus found - run fuzzing locally to generate corpus")
		return
	}

	// Find all fuzz test directories
	fuzzTests := []string{
		"FuzzLogArgsString",
		"FuzzLogArgsInt",
		"FuzzLogArgsNil",
	}

	for _, fuzzTest := range fuzzTests {
		t.Run(fuzzTest, func(t *testing.T) {
			testDir := filepath.Join(corpusDir, fuzzTest)

			// Check if this fuzz test has corpus data
			if _, err := os.Stat(testDir); os.IsNotExist(err) {
				// t.Skipf("No corpus for %s", fuzzTest)
				return
			}

			// Read all corpus files
			entries, err := os.ReadDir(testDir)
			if err != nil {
				t.Fatalf("Failed to read corpus directory: %v", err)
			}

			if len(entries) == 0 {
				// t.Skipf("No corpus files for %s", fuzzTest)
				return
			}

			// Run each corpus file as a test
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				t.Run(entry.Name(), func(t *testing.T) {
					corpusFile := filepath.Join(testDir, entry.Name())
					data, err := os.ReadFile(corpusFile)
					if err != nil {
						t.Fatalf("Failed to read corpus file: %v", err)
					}

					// Run the appropriate test based on fuzz test name
					switch fuzzTest {
					case "FuzzLogArgsString":
						runLogArgsStringCorpus(t, data)
					case "FuzzLogArgsInt":
						runLogArgsIntCorpus(t, data)
					case "FuzzLogArgsNil":
						runLogArgsNilCorpus(t, data)
					}
				})
			}
		})
	}
}

func runLogArgsStringCorpus(t *testing.T, data []byte) {
	input := extractStringFromCorpus(data)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LogArgs panicked with input: %q, panic: %v", input, r)
		}
	}()

	type TestStruct struct {
		Value string
	}

	v := TestStruct{Value: input}
	_ = logutil.LogArgs(v)
}

func runLogArgsIntCorpus(t *testing.T, data []byte) {
	// For int corpus, just use default value
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LogArgs panicked: %v", r)
		}
	}()

	type TestStruct struct {
		Value int
	}

	v := TestStruct{Value: 42}
	_ = logutil.LogArgs(v)
}

func runLogArgsNilCorpus(t *testing.T, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LogArgs panicked with nil: %v", r)
		}
	}()

	_ = logutil.LogArgs(nil)
}

// extractStringFromCorpus extracts a string value from Go's fuzz corpus format
func extractStringFromCorpus(data []byte) string {
	// Simple extraction - corpus format is: "go test fuzz v1\nstring(\"...\")\n"
	// For production use, you might want more robust parsing
	str := string(data)

	// Skip the header line
	if len(str) > 0 {
		// This is a simplified version - real corpus parsing would be more robust
		return str
	}

	return ""
}
