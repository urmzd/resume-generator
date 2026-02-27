package generators

import (
	"reflect"
	"strings"
)

// skipEscapeTags lists JSON tag names whose values should not be escaped
// (e.g. URIs which must remain literal for \href / \url commands).
var skipEscapeTags = map[string]bool{
	"uri": true,
}

// escapeStructStrings returns a deep copy of the input struct with all string
// fields transformed by escapeFn. It walks the struct recursively, handling
// pointers, slices, maps, and nested structs. Fields whose JSON tag is in
// skipEscapeTags are left untouched. Unexported fields are skipped.
func escapeStructStrings(v any, escapeFn func(string) string) any {
	result := escapeValue(reflect.ValueOf(v), escapeFn, false)
	return result.Interface()
}

func escapeValue(v reflect.Value, escapeFn func(string) string, skip bool) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return v
		}
		elem := escapeValue(v.Elem(), escapeFn, skip)
		ptr := reflect.New(elem.Type())
		ptr.Elem().Set(elem)
		return ptr

	case reflect.Struct:
		// time.Time and similar stdlib structs — return as-is
		if v.Type().PkgPath() != "" && v.Type().PkgPath() != "github.com/urmzd/resume-generator/pkg/resume" {
			return v
		}
		cp := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !field.CanInterface() {
				continue
			}
			sf := v.Type().Field(i)
			shouldSkip := shouldSkipField(sf)
			cp.Field(i).Set(escapeValue(field, escapeFn, shouldSkip))
		}
		return cp

	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		cp := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			cp.Index(i).Set(escapeValue(v.Index(i), escapeFn, skip))
		}
		return cp

	case reflect.String:
		if skip {
			return v
		}
		return reflect.ValueOf(escapeFn(v.String()))

	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		return escapeValue(v.Elem(), escapeFn, skip)

	default:
		return v
	}
}

// shouldSkipField returns true if the struct field's JSON tag name is in skipEscapeTags.
func shouldSkipField(sf reflect.StructField) bool {
	tag := sf.Tag.Get("json")
	if tag == "" || tag == "-" {
		return false
	}
	name, _, _ := strings.Cut(tag, ",")
	return skipEscapeTags[name]
}
