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

//go:build go1.18
// +build go1.18

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"math"
)

// Atomic provides atomic access to a generic value type.
type Atomic[N int64 | float64] interface {
	// TODO: Add needed atomic methods.
}

// Int64 is an int64 implementation of an Atomic.
//
// An Int64 must not be copied.
type Int64 struct {
	value *int64
}

var _ Atomic[int64] = Int64{}

func NewInt64(v int64) Int64 {
	return Int64{value: &v}
}

// Float64 is a float64 implementation of an Atomic.
//
// An Float64 must not be copied.
type Float64 struct {
	value *uint64
}

var _ Atomic[float64] = Float64{}

func NewFloat64(v float64) Float64 {
	u := math.Float64bits(v)
	return Float64{value: &u}
}
