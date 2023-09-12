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

package cmp // import "go.opentelemetry.io/otel/sdk/resource/internal/schema/cmp"

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource/internal/schema"
)

// Result is the result of comparing two things.
type Result uint8

const (
	// invalidResult is an invalid result, it should not be used.
	invalidResult Result = iota

	// EqualTo is used to indicate two things are equal.
	EqualTo
	// GreaterThan is used to indicate one thing is greater than another.
	GreaterThan
	// LessThan is used to indicate one thing is less than another.
	LessThan
)

type errInvalidVer struct {
	ver string
	err error
}

func (e *errInvalidVer) Error() string {
	return fmt.Sprintf("invalid version for %q: %s", e.ver, e.err)
}

// Versions compares schema URL versions and returns the Result of a vs b (i.e.
// a is [result value] b).
func Versions(urlA, urlB string) (Result, error) {
	aVer, err := schema.Version(urlA)
	if err != nil {
		return invalidResult, &errInvalidVer{ver: urlA, err: err}
	}

	bVer, err := schema.Version(urlB)
	if err != nil {
		return invalidResult, &errInvalidVer{ver: urlB, err: err}
	}

	switch aVer.Compare(bVer) {
	case -1:
		return LessThan, nil
	case 0:
		return EqualTo, nil
	case 1:
		return GreaterThan, nil
	default:
		msg := fmt.Sprintf("unknown comparison: %q, %q", aVer, bVer)
		panic(msg)
	}
}
