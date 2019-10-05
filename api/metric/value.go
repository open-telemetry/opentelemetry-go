// Copyright 2019, OpenTelemetry Authors
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

package metric

import (
	"fmt"
	"math"
)

// MeasurementValue represents either an integral or a floating point
// value of a measurement. It needs to be accompanied with a
// descriptor of a metric that generated this value to decide what
// type of value it represents.
type MeasurementValue uint64

// NewInt64MeasurementValue creates an integral MeasurementValue.
func NewInt64MeasurementValue(i int64) MeasurementValue {
	return newFromRaw(int64ToRaw(i))
}

// NewFloat64MeasurementValue creates a floating point
// MeasurementValue.
func NewFloat64MeasurementValue(f float64) MeasurementValue {
	return newFromRaw(float64ToRaw(f))
}

func newFromRaw(raw uint64) MeasurementValue {
	return MeasurementValue(raw)
}

func (v MeasurementValue) AsInt64() int64 {
	return rawToInt64(v.AsRaw())
}

func (v MeasurementValue) AsFloat64() float64 {
	return rawToFloat64(v.AsRaw())
}

func (v MeasurementValue) AsRaw() uint64 {
	return uint64(v)
}

func (v *MeasurementValue) AsRawPtr() *uint64 {
	return (*uint64)(v)
}

// Emit returns a string representation of the actual value of the
// MeasurementValue. A %d is used for integral values, %f for floating
// point values.
func (v MeasurementValue) Emit(kind ValueKind) string {
	switch kind {
	case Int64ValueKind:
		return fmt.Sprintf("%d", v.AsInt64())
	case Float64ValueKind:
		return fmt.Sprintf("%f", v.AsFloat64())
	default:
		return ""
	}
}

func (v MeasurementValue) Float64Compare(other float64) int {
	this := v.AsFloat64()
	if this < other {
		return -1
	} else if this > other {
		return 1
	}
	return 0
}

func (v MeasurementValue) Int64Compare(other int64) int {
	this := v.AsInt64()
	if this < other {
		return -1
	} else if this > other {
		return 1
	}
	return 0
}

func (v MeasurementValue) RawCompare(other uint64, kind ValueKind) int {
	switch kind {
	case Int64ValueKind:
		return v.Int64Compare(rawToInt64(other))
	case Float64ValueKind:
		return v.Float64Compare(rawToFloat64(other))
	default:
		// you get what you deserve
		return 0
	}
}

func (v MeasurementValue) IsPositive(kind ValueKind) bool {
	return v.compareWithZero(kind) > 0
}

func (v MeasurementValue) IsNegative(kind ValueKind) bool {
	return v.compareWithZero(kind) < 0
}

func (v MeasurementValue) IsZero(kind ValueKind) bool {
	return v.compareWithZero(kind) == 0
}

func (v MeasurementValue) compareWithZero(kind ValueKind) int {
	switch kind {
	case Int64ValueKind:
		return v.Int64Compare(0)
	case Float64ValueKind:
		return v.Float64Compare(0.)
	default:
		// you get what you deserve
		return 0
	}
}

func rawToFloat64(r uint64) float64 {
	return math.Float64frombits(r)
}

func float64ToRaw(f float64) uint64 {
	return math.Float64bits(f)
}

func rawToInt64(r uint64) int64 {
	return int64(r)
}

func int64ToRaw(i int64) uint64 {
	return uint64(i)
}
