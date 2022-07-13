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

package export // import "go.opentelemetry.io/otel/sdk/metric/export"

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

func assertCompare(equal bool, explination []string) func(*testing.T) bool {
	if equal {
		return func(*testing.T) bool { return true }
	}
	return func(t *testing.T) bool {
		return assert.Fail(t, strings.Join(explination, "\n"))
	}
}

// AssertResourceMetricsEqual asserts that two ResourceMetrics are equal.
func AssertResourceMetricsEqual(t *testing.T, expected, actual ResourceMetrics) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertScopeMetricsEqual asserts that two ScopeMetrics are equal.
func AssertScopeMetricsEqual(t *testing.T, expected, actual ScopeMetrics) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertMetricsEqual asserts that two Metrics are equal.
func AssertMetricsEqual(t *testing.T, expected, actual Metrics) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertGaugesEqual asserts that two Gauge are equal.
func AssertGaugesEqual(t *testing.T, expected, actual Gauge) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertSumsEqual asserts that two Sum are equal.
func AssertSumsEqual(t *testing.T, expected, actual Sum) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertHistogramsEqual asserts that two Histogram are equal.
func AssertHistogramsEqual(t *testing.T, expected, actual Histogram) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertDataPointsEqual asserts that two DataPoint are equal.
func AssertDataPointsEqual(t *testing.T, expected, actual DataPoint) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertHistogramDataPointsEqual asserts that two HistogramDataPoint are equal.
func AssertHistogramDataPointsEqual(t *testing.T, expected, actual HistogramDataPoint) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertInt64sEqual asserts that two Int64 are equal.
func AssertInt64sEqual(t *testing.T, expected, actual Int64) bool {
	return assertCompare(expected.compare(actual))(t)
}

// AssertFloat64sEqual asserts that two Float64 are equal.
func AssertFloat64sEqual(t *testing.T, expected, actual Float64) bool {
	return assertCompare(expected.compare(actual))(t)
}

func TestResourceMetricsComparison(t *testing.T) {
	a := ResourceMetrics{
		Resource: resource.NewSchemaless(attribute.String("resource", "a")),
	}

	b := ResourceMetrics{
		Resource: resource.NewSchemaless(attribute.String("resource", "b")),
		ScopeMetrics: []ScopeMetrics{
			{Scope: instrumentation.Scope{Name: "b"}},
		},
	}

	AssertResourceMetricsEqual(t, a, a)
	AssertResourceMetricsEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 2, "Resource and ScopeMetrics do not match")
}

func TestScopeMetricsComparison(t *testing.T) {
	a := ScopeMetrics{
		Scope: instrumentation.Scope{Name: "a"},
	}

	b := ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "b"},
		Metrics: []Metrics{{Name: "b"}},
	}

	AssertScopeMetricsEqual(t, a, a)
	AssertScopeMetricsEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 2, "Scope and Metrics do not match")
}

func TestMetricsComparison(t *testing.T) {
	a := Metrics{
		Name:        "a",
		Description: "a desc",
		Unit:        unit.Dimensionless,
	}

	b := Metrics{
		Name:        "b",
		Description: "b desc",
		Unit:        unit.Bytes,
		Data: Gauge{
			DataPoints: []DataPoint{
				{
					Attributes: attribute.NewSet(attribute.Bool("b", true)),
					StartTime:  time.Now(),
					Time:       time.Now(),
					Value:      Int64(1),
				},
			},
		},
	}

	AssertMetricsEqual(t, a, a)
	AssertMetricsEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 4, "Name, Description, Unit, and Data do not match")
}

func TestGaugesComparison(t *testing.T) {
	a := Gauge{
		DataPoints: []DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("a", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      Int64(2),
			},
		},
	}

	b := Gauge{
		DataPoints: []DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("b", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      Int64(1),
			},
		},
	}

	AssertGaugesEqual(t, a, a)
	AssertGaugesEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 1, "DataPoints do not match")
}

func TestSumsComparison(t *testing.T) {
	a := Sum{
		Temporality: CumulativeTemporality,
		IsMonotonic: true,
		DataPoints: []DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("a", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      Int64(2),
			},
		},
	}

	b := Sum{
		Temporality: DeltaTemporality,
		IsMonotonic: false,
		DataPoints: []DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("b", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      Int64(1),
			},
		},
	}

	AssertSumsEqual(t, a, a)
	AssertSumsEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 3, "Temporality, IsMonotonic, and DataPoints do not match")
}

func TestHistogramsComparison(t *testing.T) {
	a := Histogram{
		Temporality: CumulativeTemporality,
	}

	b := Histogram{
		Temporality: DeltaTemporality,
		DataPoints: []HistogramDataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("b", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
			},
		},
	}

	AssertHistogramsEqual(t, a, a)
	AssertHistogramsEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 2, "Temporality and DataPoints do not match")
}

func TestDataPointsComparison(t *testing.T) {
	a := DataPoint{
		Attributes: attribute.NewSet(attribute.Bool("a", true)),
		StartTime:  time.Now(),
		Time:       time.Now(),
		Value:      Int64(2),
	}

	b := DataPoint{
		Attributes: attribute.NewSet(attribute.Bool("b", true)),
		StartTime:  time.Now(),
		Time:       time.Now(),
		Value:      Float64(1),
	}

	AssertDataPointsEqual(t, a, a)
	AssertDataPointsEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 4, "Attributes, StartTime, Time and Value do not match")
}

func TestInt64sComparison(t *testing.T) {
	a := Int64(-1)

	b := Int64(2)

	AssertInt64sEqual(t, a, a)
	AssertInt64sEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 1, "Value does not match")
}

func TestFloat64sComparison(t *testing.T) {
	a := Float64(-1)

	b := Float64(2)

	AssertFloat64sEqual(t, a, a)
	AssertFloat64sEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 1, "Value does not match")
}

func TestHistogramDataPointsComparison(t *testing.T) {
	a := HistogramDataPoint{
		Attributes:   attribute.NewSet(attribute.Bool("a", true)),
		StartTime:    time.Now(),
		Time:         time.Now(),
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Sum:          2,
	}

	max, min := 99.0, 3.
	b := HistogramDataPoint{
		Attributes:   attribute.NewSet(attribute.Bool("b", true)),
		StartTime:    time.Now(),
		Time:         time.Now(),
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          &max,
		Min:          &min,
		Sum:          3,
	}

	AssertHistogramDataPointsEqual(t, a, a)
	AssertHistogramDataPointsEqual(t, b, b)

	equal, explination := a.compare(b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 9, "Attributes, StartTime, Time, Count, Bounds, BucketCounts, Max, Min, and Sum do not match")
}
