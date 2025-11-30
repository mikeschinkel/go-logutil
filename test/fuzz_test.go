package test

import (
	"testing"

	"github.com/mikeschinkel/go-logutil"
)

// FuzzLogArgsString tests LogArgs with various string inputs
func FuzzLogArgsString(f *testing.F) {
	// Seed corpus with various string values
	seeds := []string{
		"",
		"simple",
		"with spaces",
		"Unicode: 你好世界",
		"Special!@#$%^&*()",
		"Multi\nLine\nString",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Create a simple struct with a string field
		type TestStruct struct {
			Value string
		}

		// Ensure doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LogArgs panicked with input %q: %v", input, r)
			}
		}()

		v := TestStruct{Value: input}
		_ = logutil.LogArgs(v)
	})
}

// FuzzLogArgsInt tests LogArgs with various integer inputs
func FuzzLogArgsInt(f *testing.F) {
	// Seed corpus with various integer values
	seeds := []int{
		0,
		1,
		-1,
		42,
		-42,
		2147483647,  // max int32
		-2147483648, // min int32
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input int) {
		// Create a simple struct with an int field
		type TestStruct struct {
			Value int
		}

		// Ensure doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LogArgs panicked with input %d: %v", input, r)
			}
		}()

		v := TestStruct{Value: input}
		_ = logutil.LogArgs(v)
	})
}

// FuzzLogArgsNil tests LogArgs with nil inputs
func FuzzLogArgsNil(f *testing.F) {
	// Seed with a simple value to start the fuzzer
	f.Add(true)

	f.Fuzz(func(t *testing.T, _ bool) {
		// Ensure doesn't panic with nil
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LogArgs panicked with nil: %v", r)
			}
		}()

		_ = logutil.LogArgs(nil)
	})
}
