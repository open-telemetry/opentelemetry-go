// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package attribute provide several helper functions for some commonly used
logic of processing attributes.
*/
package attribute // import "go.opentelemetry.io/otel/internal/attribute"

import (
	"reflect"
)

// BoolSliceValue converts a bool slice into an array with same elements as slice.
func BoolSliceValue(v []bool) any {
	var zero bool
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero)))
	copy(cp.Elem().Slice(0, len(v)).Interface().([]bool), v)
	return cp.Elem().Interface()
}

// SliceValue convert a int64 slice into an array with same elements as slice.
func Int64SliceValue(v []int64) any {
	var zero int64
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero)))
	copy(cp.Elem().Slice(0, len(v)).Interface().([]int64), v)
	return cp.Elem().Interface()
}

// SliceValue convert a float64 slice into an array with same elements as slice.
func Float64SliceValue(v []float64) any {
	var zero float64
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero)))
	copy(cp.Elem().Slice(0, len(v)).Interface().([]float64), v)
	return cp.Elem().Interface()
}

// SliceValue convert a string slice into an array with same elements as slice.
func StringSliceValue(v []string) any {
	var zero string
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero)))
	copy(cp.Elem().Slice(0, len(v)).Interface().([]string), v)
	return cp.Elem().Interface()
}

// AsSlice convert an bool array into a slice into with same elements as array.
func AsBoolSlice(v any) []bool {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	var zero bool
	correctLen := rv.Len()
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(zero))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]bool)
}

// AsSlice convert an int64 array into a slice into with same elements as array.
func AsInt64Slice(v any) []int64 {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	var zero int64
	correctLen := rv.Len()
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(zero))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]int64)
}

// AsSlice convert an float64 array into a slice into with same elements as array.
func AsFloat64Slice(v any) []float64 {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	var zero float64
	correctLen := rv.Len()
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(zero))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]float64)
}

// AsSlice convert an string array into a slice into with same elements as array.
func AsStringSlice(v any) []string {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	var zero string
	correctLen := rv.Len()
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(zero))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]string)
}
