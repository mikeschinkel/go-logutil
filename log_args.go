package logutil

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"
)

func LogArgs(v any) (attrs []any) {
	var rv reflect.Value
	var rt reflect.Type

	if v == nil {
		goto end
	}

	rv = reflect.ValueOf(v)
	rt = rv.Type()

	// Handle pointer to struct
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			goto end
		}
		rv = rv.Elem()
		rt = rv.Type()
	}

	// If it is not a struct, just log it as "value"
	if rv.Kind() != reflect.Struct {
		attrs = []any{slog.Any("value", v)}
		goto end
	}

	attrs = make([]any, 0, rv.NumField())

	for i := 0; i < rv.NumField(); i++ {
		attr := logArg(rv, rt, i)
		if attr.Key == "" {
			continue
		}
		attrs = append(attrs, logArg(rv, rt, i))
	}
end:
	return attrs
}

func logArg(rv reflect.Value, rt reflect.Type, fldNo int) (attr slog.Attr) {
	var tag string
	var name string
	var omitEmpty bool

	field := rt.Field(fldNo)
	value := rv.Field(fldNo)

	// Skip unexported fields
	if !field.IsExported() {
		goto end
	}

	tag = field.Tag.Get("json")

	// Skip json:"-"
	if tag == "-" {
		goto end
	}

	name = field.Name
	omitEmpty = false

	if tag != "" {
		parts := strings.Split(tag, ",")
		if len(parts[0]) > 0 {
			name = parts[0]
		}
		for _, p := range parts[1:] {
			if p == "omitempty" {
				omitEmpty = true
				break
			}
		}
	}

	// Honor ,omitempty if present
	if omitEmpty {
		if value.IsZero() {
			goto end
		}
	}

	attr = formatAttr(name, value)

end:
	return attr
}

func formatAttr(name string, v reflect.Value) (attr slog.Attr) {
	var value any

	// Handle invalid values defensively
	if !v.IsValid() {
		goto end
	}

	// Deref pointers so *T works too
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			goto end
		}
		v = v.Elem()
	}

	if !v.CanInterface() {
		goto end
	}

	value = v.Interface()

	// 1) time.Time as RFC3339 string
	if t, ok := value.(time.Time); ok {
		// If you want to keep local time, drop the .UTC()
		attr = slog.String(name, t.UTC().Format(time.RFC3339Nano))
		goto end
	}

	// 2) error as error message
	if err, ok := value.(error); ok {
		attr = slog.String(name, err.Error())
		goto end
	}

	// 3) fmt.Stringer (enums, domain types, etc.)
	if s, ok := value.(fmt.Stringer); ok {
		attr = slog.String(name, s.String())
		goto end
	}

	// 4) nested structs as groups
	if v.Kind() == reflect.Struct {
		args := LogArgs(value)
		// If it has no exported/loggable fields, skip it
		if len(args) == 0 {
			goto end
		}
		attr = slog.Group(name, args...)
		goto end
	}

	// 5) fallback: let slog figure it out
	attr = slog.Any(name, value)

end:
	return attr
}
