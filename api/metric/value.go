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

//go:generate stringer -type=ValueKind

// ValueKind describes the data type of the measurement value the
// metric generates.
type ValueKind int8

const (
	// Int64ValueKind means that the metric generates values of
	// type int64.
	Int64ValueKind ValueKind = iota
	// Float64ValueKind means that the metric generates values of
	// type float64.
	Float64ValueKind
)

// MeasurementValue represents either an integral or a floating point
// value of a measurement. It needs to be accompanied with a value
// kind or some source that provides a value kind describing this
// measurement value.
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

// AsInt64 assumes that the measurement value contains an int64 and
// returns it as such. Make sure that the accompanying source of value
// kind indeed tells you its a 64 bit integral measurement value,
// otherwise the returned int64 will be wrong.
func (v MeasurementValue) AsInt64() int64 {
	return rawToInt64(v.AsRaw())
}

// AsFloat64 assumes that the measurement value contains a float64 and
// returns it as such. Make sure that the accompanying source of value
// kind indeed tells you its a 64 bit floating point measurement
// value, otherwise the returned float64 will be wrong.
func (v MeasurementValue) AsFloat64() float64 {
	return rawToFloat64(v.AsRaw())
}

// AsRaw gets the raw, uninterpreted value of the measurement. Might
// be useful for some atomic operations.
func (v MeasurementValue) AsRaw() uint64 {
	return uint64(v)
}

// AsRawPtr gets the pointer to the raw, uninterpreted value of the
// measurement. Might be useful for some atomic operations.
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

// Float64Compare assumes that the MeasurementValue contains a float64
// and performs a comparison between the value and the other value. It
// returns the typical result of the compare function: -1 if the value
// is less than the other, 0 if both are equal, 1 if the value is
// greater than the other.
func (v MeasurementValue) Float64Compare(other float64) int {
	this := v.AsFloat64()
	if this < other {
		return -1
	} else if this > other {
		return 1
	}
	return 0
}

// Int64Compare assumes that the MeasurementValue contains an int64
// and performs a comparison between the value and the other value. It
// returns the typical result of the compare function: -1 if the value
// is less than the other, 0 if both are equal, 1 if the value is
// greater than the other.
func (v MeasurementValue) Int64Compare(other int64) int {
	this := v.AsInt64()
	if this < other {
		return -1
	} else if this > other {
		return 1
	}
	return 0
}

// RawCompare calls either Float64Compare or Int64Compare, depending
// on the passed kind.
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

// IsPositive returns true if the actual value is greater than zero.
func (v MeasurementValue) IsPositive(kind ValueKind) bool {
	return v.compareWithZero(kind) > 0
}

// IsNegative returns true if the actual value is less than zero.
func (v MeasurementValue) IsNegative(kind ValueKind) bool {
	return v.compareWithZero(kind) < 0
}

// IsZero returns true if the actual value is equal to zero.
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
