// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/internaltest/alignment.go.tmpl

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

package internaltest // import "go.opentelemetry.io/otel/internal/internaltest"

/*
This file contains common utilities and objects to validate memory alignment
of Go types. The primary use of this functionality is intended to ensure
`struct` fields that need to be 64-bit aligned so they can be passed as
arguments to 64-bit atomic operations.

The common workflow is to define a slice of `FieldOffset` and pass them to the
`Aligned8Byte` function from within a `TestMain` function from a package's
tests. It is important to make this call from the `TestMain` function prior
to running the rest of the test suit as it can provide useful diagnostics
about field alignment instead of ambiguous nil pointer dereference and runtime
panic.

For more information:
https://github.com/open-telemetry/opentelemetry-go/issues/341
*/

import (
	"fmt"
	"io"
)

// FieldOffset is a preprocessor representation of a struct field alignment.
type FieldOffset struct {
	// Name of the field.
	Name string

	// Offset of the field in bytes.
	//
	// To compute this at compile time use unsafe.Offsetof.
	Offset uintptr
}

// Aligned8Byte returns if all fields are aligned modulo 8-bytes.
//
// Error messaging is printed to out for any field determined misaligned.
func Aligned8Byte(fields []FieldOffset, out io.Writer) bool {
	misaligned := make([]FieldOffset, 0)
	for _, f := range fields {
		if f.Offset%8 != 0 {
			misaligned = append(misaligned, f)
		}
	}

	if len(misaligned) == 0 {
		return true
	}

	fmt.Fprintln(out, "struct fields not aligned for 64-bit atomic operations:")
	for _, f := range misaligned {
		fmt.Fprintf(out, "  %s: %d-byte offset\n", f.Name, f.Offset)
	}

	return false
}
