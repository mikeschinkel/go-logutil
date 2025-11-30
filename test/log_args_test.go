package test

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-logutil"
)

// LogEntryType is a simple enum that implements fmt.Stringer
// so we can verify Stringer-based formatting.
type LogEntryType int

const (
	EntryTypeUnknown LogEntryType = iota
	EntryTypeSearch
)

func (t LogEntryType) String() string {
	switch t {
	case EntryTypeSearch:
		return "search"
	case EntryTypeUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

type nested struct {
	ID   int       `json:"id"`
	When time.Time `json:"when"`
}

type testStruct struct {
	unexported string       // should be skipped
	ID         int          `json:"id"`
	Name       string       `json:"name,omitempty"`   // omit when empty
	When       time.Time    `json:"when"`             // time formatting
	Err        error        `json:"err,omitempty"`    // error formatting
	Type       LogEntryType `json:"type"`             // Stringer formatting
	Hidden     string       `json:"-"`                // explicit skip
	Nested     nested       `json:"nested,omitempty"` // group when non-zero
	Slice      []string     `json:"slice,omitempty"`  // omit when empty
	PtrNested  *nested      `json:"ptr_nested,omitempty"`
}

// attrsToMap is a helper to make assertions easier.
// It maps attr.Key -> attr.Value, assuming keys are unique.
func attrsToMap(attrs any) (m map[string]slog.Value, err error) {
	var errs []error
	switch t := attrs.(type) {
	case []slog.Attr:
		m = make(map[string]slog.Value, len(t))
		for _, attr := range t {
			m[attr.Key] = attr.Value
		}
	case []any:
		m = make(map[string]slog.Value, len(t))
		for _, attr := range t {
			a, ok := attr.(slog.Attr)
			if !ok {
				errs = append(errs, fmt.Errorf("expected `slog.Attr`, got `%T`", a))
				continue
			}
			m[a.Key] = a.Value
		}
		err = dt.CombineErrs(errs)
	default:
		goto end
	}
end:
	return m, err
}

func TestLogArgs_Nil(t *testing.T) {
	got := logutil.LogArgs(nil)
	if len(got) != 0 {
		t.Fatalf("expected 0 attrs for nil, got %d", len(got))
	}
}

func TestLogArgs_NonStruct(t *testing.T) {
	got := logutil.LogArgs("hello")

	if len(got) != 1 {
		t.Fatalf("expected 1 attr, got %d", len(got))
	}

	attr, ok := got[0].(slog.Attr)
	if !ok {
		t.Fatalf("expected `slog.Attr`, got `%T`", attr.Key)
	}
	if attr.Key != "value" {
		t.Fatalf("expected key %q, got %q", "value", attr.Key)
	}
	if attr.Value.Kind() != slog.KindString {
		t.Fatalf("expected KindString, got %v", attr.Value.Kind())
	}
	if attr.Value.String() != "hello" {
		t.Fatalf("expected value %q, got %q", "hello", attr.Value.String())
	}
}

func TestLogArgs_StructFormatting(t *testing.T) {
	var v slog.Value
	var ok bool

	now := time.Date(2025, 11, 22, 12, 0, 0, 123456789, time.UTC)
	err := errors.New("boom")

	ts := testStruct{
		unexported: "should be skipped",
		ID:         123,
		Name:       "", // omitempty -> skipped
		When:       now,
		Err:        err,
		Type:       EntryTypeSearch,
		Hidden:     "should be skipped",
		Nested: nested{
			ID:   1,
			When: now,
		},
		// Slice is nil -> omitempty -> skipped
		PtrNested: &nested{
			ID:   2,
			When: now,
		},
	}

	got := logutil.LogArgs(ts)
	if len(got) == 0 {
		t.Fatalf("expected some attrs, got none")
	}

	m, err := attrsToMap(got)
	if err != nil {
		t.Fatal(err.Error())
	}
	// Unexported and json:"-" fields must not appear
	if _, ok := m["unexported"]; ok {
		t.Fatalf("unexported field should not be logged")
	}
	if _, ok := m["Hidden"]; ok {
		t.Fatalf("field with json:\"-\" should not be logged")
	}

	// ID -> int64
	v, ok = m["id"]
	if !ok {
		t.Fatalf("missing \"id\" attr")
	}
	if v.Kind() != slog.KindInt64 {
		t.Fatalf("id: expected KindInt64, got %v", v.Kind())
	}
	if v.Int64() != 123 {
		t.Fatalf("id: expected 123, got %d", v.Int64())
	}

	// Name had omitempty and is empty -> should be omitted
	if _, ok := m["name"]; ok {
		t.Fatalf("name should be omitted due to omitempty and zero value")
	}

	// When: time.Time -> RFC3339Nano string
	v, ok = m["when"]
	if !ok {
		t.Fatalf("missing \"when\" attr")
	}
	if v.Kind() != slog.KindString {
		t.Fatalf("when: expected KindString, got %v", v.Kind())
	}
	expected := now.UTC().Format(time.RFC3339Nano)
	if v.String() != expected {
		t.Fatalf("when: expected %q, got %q", expected, v.String())
	}

	// Err: error -> string
	v, ok = m["err"]
	if !ok {
		t.Fatalf("missing \"err\" attr")
	}
	if v.Kind() != slog.KindString {
		t.Fatalf("err: expected KindString, got %v", v.Kind())
	}
	if v.String() != "boom" {
		t.Fatalf("err: expected %q, got %q", "boom", v.String())
	}

	// Type: LogEntryType (Stringer) -> string
	if v, ok = m["type"]; !ok {
		t.Fatalf("missing \"type\" attr")
	}
	if v.Kind() != slog.KindString {
		t.Fatalf("type: expected KindString, got %v", v.Kind())
	}
	if v.String() != "search" {
		t.Fatalf("type: expected %q, got %q", "search", v.String())
	}

	// Slice: omitempty + zero value -> should be omitted
	if _, ok := m["slice"]; ok {
		t.Fatalf("slice should be omitted due to omitempty and zero value")
	}

	// Nested: struct -> group
	nestedVal, ok := m["nested"]
	if !ok {
		t.Fatalf("missing \"nested\" attr")
	}
	if nestedVal.Kind() != slog.KindGroup {
		t.Fatalf("nested: expected KindGroup, got %v", nestedVal.Kind())
	}
	nestedAttrs := nestedVal.Group()
	nm, err := attrsToMap(nestedAttrs)
	if err != nil {
		t.Fatal(err.Error())
	}

	v, ok = nm["id"]
	if !ok {
		t.Fatalf("nested.id missing")
	}
	if v.Int64() != 1 {
		t.Fatalf("nested.id: expected 1, got %d", v.Int64())
	}

	v, ok = nm["when"]
	if !ok {
		t.Fatalf("nested.when missing")
	}

	expected = now.UTC().Format(time.RFC3339Nano)
	if v.String() != expected {
		t.Fatalf("nested.when: expected %q, got %q", expected, v.String())
	}

	// PtrNested: *nested -> should also be a group
	ptrVal, ok := m["ptr_nested"]
	if !ok {
		t.Fatalf("missing \"ptr_nested\" attr")
	}
	if ptrVal.Kind() != slog.KindGroup {
		t.Fatalf("ptr_nested: expected KindGroup, got %v", ptrVal.Kind())
	}
	ptrAttrs := ptrVal.Group()
	pm, err := attrsToMap(ptrAttrs)
	if err != nil {
		t.Fatal(err.Error())
	}

	v, ok = pm["id"]
	if !ok {
		t.Fatalf("ptr_nested.id missing")
	}
	if v.Int64() != 2 {
		t.Fatalf("ptr_nested.id: expected 2, got %d", v.Int64())
	}
}

func TestLogArgs_PointerToStruct(t *testing.T) {
	ts := &testStruct{
		ID: 42,
	}
	got := logutil.LogArgs(ts)
	if len(got) == 0 {
		t.Fatalf("expected attrs for pointer to struct, got none")
	}
	m, err := attrsToMap(got)
	if err != nil {
		t.Fatal(err.Error())
	}
	v, ok := m["id"]
	if !ok {
		t.Fatalf("missing id attr")
	}
	if v.Int64() != 42 {
		t.Fatalf("id: expected 42, got %d", v.Int64())
	}
}
