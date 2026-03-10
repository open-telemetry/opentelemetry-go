// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package attribute provide several helper functions for some commonly used
logic of processing attributes.
*/
package attribute // import "go.opentelemetry.io/otel/attribute/internal"

import (
	"reflect"
)

// SliceValue converts a slice into an array with the same elements.
func SliceValue[T any](v []T) any {
	// Keep the common tiny-slice cases out of reflection. This matches the
	// short lengths that show up most often in local benchmarks and semantic
	// convention examples while leaving larger, less predictable slices on the
	// generic reflective path.
	switch len(v) {
	case 0:
		return [0]T{}
	case 1:
		return [1]T{v[0]}
	case 2:
		return [2]T{v[0], v[1]}
	case 3:
		return [3]T{v[0], v[1], v[2]}
	}

	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeFor[T]())).Elem()
	reflect.Copy(cp, reflect.ValueOf(v))
	return cp.Interface()
}

// AsSlice converts an array into a slice with the same elements.
func AsSlice[T any](v any) []T {
	// Mirror the small fixed-array fast path used by SliceValue.
	switch a := v.(type) {
	case [0]T:
		return []T{}
	case [1]T:
		return []T{a[0]}
	case [2]T:
		return []T{a[0], a[1]}
	case [3]T:
		return []T{a[0], a[1], a[2]}
	}

	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Array || rv.Type().Elem() != reflect.TypeFor[T]() {
		return nil
	}
	cpy := make([]T, rv.Len())
	if len(cpy) > 0 {
		_ = reflect.Copy(reflect.ValueOf(cpy), rv)
	}
	return cpy
}
