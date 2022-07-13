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

// TODO: NOTE this is a temporary space, it may be moved following the
// discussion of #2813, or #2841

package export // import "go.opentelemetry.io/otel/sdk/metric/export"

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

// ResourceMetrics is a collection of ScopeMetrics and the associated Resource
// that created them.
type ResourceMetrics struct {
	// Resource represents the entity that collected the metrics.
	Resource *resource.Resource
	// ScopeMetrics are the collection of metrics with unique Scopes.
	ScopeMetrics []ScopeMetrics
}

// compare returns true when an other ResourceMetrics is equivalent to this
// ResourceMetrics. It returns false when they differ, along with messages
// describing the difference.
//
// The ScopeMetrics each ResourceMetrics contains are compared based on
// containing the same ScopeMetrics, not the order they are stored in.
func (rm ResourceMetrics) compare(other ResourceMetrics) (equal bool, explination []string) {
	equal = true
	if !rm.Resource.Equal(other.Resource) {
		equal, explination = false, append(
			explination, notEqualStr("Resources", rm.Resource, other.Resource),
		)
	}

	var exp string
	equal, exp = compareDiff(diffSlices(
		rm.ScopeMetrics,
		other.ScopeMetrics,
		func(a, b ScopeMetrics) bool {
			equal, _ := a.compare(b)
			return equal
		},
	))
	if !equal {
		explination = append(explination, fmt.Sprintf(
			"ResourceMetrics ScopeMetrics not equal:\n%s", exp,
		))
	}
	return equal, explination
}

// ScopeMetrics is a collection of Metrics Produces by a Meter.
type ScopeMetrics struct {
	// Scope is the Scope that the Meter was created with.
	Scope instrumentation.Scope
	// Metrics are a list of aggregations created by the Meter.
	Metrics []Metrics
}

// compare returns true when an other ScopeMetrics is equivalent to this
// ScopeMetrics. It returns false when they differ, along with messages
// describing the difference.
//
// The Metrics each ScopeMetrics contains are compared based on containing the
// same Metrics, not the order they are stored in.
func (sm ScopeMetrics) compare(other ScopeMetrics) (equal bool, explination []string) {
	equal = true
	if sm.Scope != other.Scope {
		equal, explination = false, append(
			explination,
			notEqualStr("Scope", sm.Scope, other.Scope),
		)
	}

	var exp string
	equal, exp = compareDiff(diffSlices(
		sm.Metrics,
		other.Metrics,
		func(a, b Metrics) bool {
			equal, _ := a.compare(b)
			return equal
		},
	))
	if !equal {
		explination = append(explination, fmt.Sprintf(
			"ScopeMetrics Metrics not equal:\n%s", exp,
		))
	}
	return equal, explination
}

// Metrics is a collection of one or more aggregated timeseries from an Instrument.
type Metrics struct {
	// Name is the name of the Instrument that created this data.
	Name string
	// Description is the description of the Instrument, which can be used in documentation.
	Description string
	// Unit is the unit in which the Instrument reports.
	Unit unit.Unit
	// Data is the aggregated data from an Instrument.
	Data Aggregation
}

// compare returns true when an other is equivalent to this Metrics. It
// returns false when they differ, along with messages describing the
// difference.
func (m Metrics) compare(other Metrics) (equal bool, explination []string) {
	equal = true
	if m.Name != other.Name {
		equal, explination = false, append(
			explination,
			notEqualStr("Name", m.Name, other.Name),
		)
	}
	if m.Description != other.Description {
		equal, explination = false, append(
			explination,
			notEqualStr("Description", m.Description, other.Description),
		)
	}
	if m.Unit != other.Unit {
		equal, explination = false, append(
			explination,
			notEqualStr("Unit", m.Unit, other.Unit),
		)
	}

	if m.Data == nil || other.Data == nil {
		if m.Data != other.Data {
			equal, explination = false, append(explination, notEqualStr(
				"Data", m.Data, other.Data,
			))
		}
		return equal, explination
	}

	if reflect.TypeOf(m.Data) != reflect.TypeOf(other.Data) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Data types not equal:\nexpected: %T\nactual: %T",
			m.Data,
			other.Data,
		))
		return equal, explination
	}

	switch v := m.Data.(type) {
	case Gauge:
		ok, exp := v.compare(other.Data.(Gauge))
		if !ok {
			equal, explination = false, append(explination, fmt.Sprintf("Gauge: %s", exp))
		}
	case Sum:
		ok, exp := v.compare(other.Data.(Sum))
		if !ok {
			equal, explination = false, append(explination, fmt.Sprintf("Sum: %s", exp))
		}
	case Histogram:
		ok, exp := v.compare(other.Data.(Histogram))
		if !ok {
			equal, explination = false, append(explination, fmt.Sprintf("Histogram: %s", exp))
		}
	default:
		equal, explination = false, append(
			explination,
			fmt.Sprintf("Data of unknown types %T", m.Data),
		)
	}
	return equal, explination
}

// Aggregation is the store of data reported by an Instrument.
// It will be one of: Gauge, Sum, Histogram.
type Aggregation interface {
	privateAggregation()
}

// Gauge represents a measurement of the current value of an instrument.
type Gauge struct {
	// DataPoints reprents individual aggregated measurements with unique Attributes.
	DataPoints []DataPoint
}

// compare returns true when an other is equivalent to this Gauge. It returns
// false when they differ, along with messages describing the difference.
//
// The DataPoints each Gauge contains are compared based on containing the
// same DataPoints, not the order they are stored in.
func (g Gauge) compare(other Gauge) (equal bool, explination []string) {
	var exp string
	equal, exp = compareDiff(diffSlices(
		g.DataPoints,
		other.DataPoints,
		func(a, b DataPoint) bool {
			equal, _ := a.compare(b)
			return equal
		},
	))
	if !equal {
		explination = append(explination, fmt.Sprintf(
			"Gauge DataPoints not equal:\n%s", exp,
		))
	}
	return equal, explination
}

func (Gauge) privateAggregation() {}

// Sum represents the sum of all measurements of values from an instrument.
type Sum struct {
	// DataPoints reprents individual aggregated measurements with unique Attributes.
	DataPoints []DataPoint
	// Temporality describes if the aggregation is reported as the change from the
	// last report time, or the cumulative changes since a fixed start time.
	Temporality Temporality
	// IsMonotonic represents if this aggregation only increases or decreases.
	IsMonotonic bool
}

// compare returns true when an other is equivalent to this Sum. It returns
// false when they differ, along with messages describing the difference.
//
// The DataPoints each Sum contains are compared based on containing the same
// DataPoints, not the order they are stored in.
func (s Sum) compare(other Sum) (equal bool, explination []string) {
	equal = true
	if s.Temporality != other.Temporality {
		equal, explination = false, append(
			explination,
			notEqualStr("Temporality", s.Temporality, other.Temporality),
		)
	}
	if s.IsMonotonic != other.IsMonotonic {
		equal, explination = false, append(
			explination,
			notEqualStr("IsMonotonic", s.IsMonotonic, other.IsMonotonic),
		)
	}

	var exp string
	equal, exp = compareDiff(diffSlices(
		s.DataPoints,
		other.DataPoints,
		func(a, b DataPoint) bool {
			equal, _ := a.compare(b)
			return equal
		},
	))
	if !equal {
		explination = append(explination, fmt.Sprintf(
			"Sum DataPoints not equal:\n%s", exp,
		))
	}
	return equal, explination
}

func (Sum) privateAggregation() {}

// DataPoint is a single data point in a timeseries.
type DataPoint struct {
	// Attributes is the set of key value pairs that uniquely identify the
	// timeseries.
	Attributes attribute.Set
	// StartTime is when the timeseries was started. (optional)
	StartTime time.Time
	// Time is the time when the timeseries was recorded. (optional)
	Time time.Time
	// Value is the value of this data point.
	Value Value
}

// compare returns true when an other is equivalent to this DataPoint. It
// returns false when they differ, along with messages describing the
// difference.
func (d DataPoint) compare(other DataPoint) (equal bool, explination []string) {
	equal = true
	if !d.Attributes.Equals(&other.Attributes) {
		equal, explination = false, append(explination, notEqualStr(
			"Attributes",
			d.Attributes.Encoded(attribute.DefaultEncoder()),
			other.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}
	if !d.StartTime.Equal(other.StartTime) {
		equal, explination = false, append(explination, notEqualStr(
			"StartTime",
			d.StartTime.UnixNano(),
			other.StartTime.UnixNano(),
		))
	}
	if !d.Time.Equal(other.Time) {
		equal, explination = false, append(explination, notEqualStr(
			"Time",
			d.Time.UnixNano(),
			other.Time.UnixNano(),
		))
	}
	if reflect.TypeOf(d.Value) != reflect.TypeOf(other.Value) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Value types not equal:\nexpected: %T\nactual: %T",
			d.Value,
			other.Value,
		))
		return equal, explination
	}

	switch v := d.Value.(type) {
	case Int64:
		ok, exp := v.compare(other.Value.(Int64))
		if !ok {
			equal, explination = false, append(explination, fmt.Sprintf("Int64: %s", exp))
		}
	case Float64:
		ok, exp := v.compare(other.Value.(Float64))
		if !ok {
			equal, explination = false, append(explination, fmt.Sprintf("Float64: %s", exp))
		}
	default:
		equal, explination = false, append(
			explination,
			fmt.Sprintf("Value of unknown types %T", d.Value),
		)
	}

	return equal, explination
}

// Value is a int64 or float64. All Values created by the sdk will be either
// Int64 or Float64.
type Value interface {
	privateValue()
}

// Int64 is a container for an int64 value.
type Int64 int64

func (Int64) privateValue() {}

// compare returns true when an other is equivalent to this Int64. It
// returns false when they differ, along with messages describing the
// difference.
func (i Int64) compare(other Int64) (equal bool, explination []string) {
	equal = true
	if i != other {
		equal, explination = false, append(
			explination, notEqualStr("Int64 value", i, other),
		)
	}
	return equal, explination
}

// Float64 is a container for a float64 value.
type Float64 float64

func (Float64) privateValue() {}

// compare returns true when an other is equivalent to this Float64. It
// returns false when they differ, along with messages describing the
// difference.
func (f Float64) compare(other Float64) (equal bool, explination []string) {
	equal = true
	if f != other {
		equal, explination = false, append(
			explination, notEqualStr("Float64 value", f, other),
		)
	}
	return equal, explination
}

// Histogram represents the histogram of all measurements of values from an instrument.
type Histogram struct {
	// DataPoints reprents individual aggregated measurements with unique Attributes.
	DataPoints []HistogramDataPoint
	// Temporality describes if the aggregation is reported as the change from the
	// last report time, or the cumulative changes since a fixed start time.
	Temporality Temporality
}

func (Histogram) privateAggregation() {}

// compare returns true when an other is equivalent to this Histogram. It
// returns false when they differ, along with messages describing the
// difference.
//
// The DataPoints each Histogram contains are compared based on containing the
// same HistogramDataPoint, not the order they are stored in.
func (h Histogram) compare(other Histogram) (equal bool, explination []string) {
	equal = true
	if h.Temporality != other.Temporality {
		equal, explination = false, append(
			explination,
			notEqualStr("Temporality", h.Temporality, other.Temporality),
		)
	}

	var exp string
	equal, exp = compareDiff(diffSlices(
		h.DataPoints,
		other.DataPoints,
		func(a, b HistogramDataPoint) bool {
			equal, _ := a.compare(b)
			return equal
		},
	))
	if !equal {
		explination = append(explination, fmt.Sprintf(
			"Histogram DataPoints not equal:\n%s", exp,
		))
	}
	return equal, explination
}

// HistogramDataPoint is a single histogram data point in a timeseries.
type HistogramDataPoint struct {
	// Attributes is the set of key value pairs that uniquely identify the
	// timeseries.
	Attributes attribute.Set
	// StartTime is when the timeseries was started.
	StartTime time.Time
	// Time is the time when the timeseries was recorded.
	Time time.Time

	// Count is the number of updates this histogram has been calculated with.
	Count uint64
	// Bounds are the upper bounds of the buckets of the histogram. Because the
	// last boundary is +infinity this one is implied.
	Bounds []float64
	// BucketCounts is the count of each of the buckets.
	BucketCounts []uint64

	// Min is the minimum value recorded. (optional)
	Min *float64
	// Max is the maximum value recorded. (optional)
	Max *float64
	// Sum is the sum of the values recorded.
	Sum float64
}

// compare returns true when an other is equivalent to this
// HistogramDataPoint. It returns false when they differ, along with messages
// describing the difference.
func (hd HistogramDataPoint) compare(other HistogramDataPoint) (equal bool, explination []string) {
	equal = true
	if !hd.Attributes.Equals(&other.Attributes) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Attributes not equal:\nexpected: %s\nactual: %s",
			hd.Attributes.Encoded(attribute.DefaultEncoder()),
			other.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}
	if !hd.StartTime.Equal(other.StartTime) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"StartTime not equal:\nexpected: %d\nactual: %d",
			hd.StartTime.UnixNano(),
			other.StartTime.UnixNano(),
		))
	}
	if !hd.Time.Equal(other.Time) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Time not equal:\nexpected: %d\nactual: %d",
			hd.Time.UnixNano(),
			other.Time.UnixNano(),
		))
	}
	if hd.Count != other.Count {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Count not equal:\nexpected: %d\nactual: %d",
			hd.Count,
			other.Count,
		))
	}
	if !equalSlices(hd.Bounds, other.Bounds) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Bounds not equal:\nexpected: %v\nactual: %v",
			hd.Bounds,
			other.Bounds,
		))
	}
	if !equalSlices(hd.BucketCounts, other.BucketCounts) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"BucketCounts not equal:\nexpected: %v\nactual: %v",
			hd.BucketCounts,
			other.BucketCounts,
		))
	}
	if !equalPtrValues(hd.Min, other.Min) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Min not equal:\nexpected: %v\nactual: %v",
			hd.Min,
			other.Min,
		))
	}
	if !equalPtrValues(hd.Max, other.Max) {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Max not equal:\nexpected: %v\nactual: %v",
			hd.Max,
			other.Max,
		))
	}
	if hd.Sum != other.Sum {
		equal, explination = false, append(explination, fmt.Sprintf(
			"Sum not equal:\nexpected: %g\nactual: %g",
			hd.Sum,
			other.Sum,
		))
	}
	return equal, explination
}

func notEqualStr(prefix string, expected, actual interface{}) string {
	return fmt.Sprintf("%s not equal:\nexpected: %v\nactual: %v", prefix, expected, actual)
}

func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func equalPtrValues[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b

}

func diffSlices[T any](a, b []T, equal func(T, T) bool) (extraA, extraB []T) {
	visited := make([]bool, len(b))
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if visited[j] {
				continue
			}
			if equal(a[i], b[j]) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, a[i])
		}
	}

	for j := 0; j < len(b); j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, b[j])
	}

	return extraA, extraB
}

func compareDiff[T any](extraExpected, extraActual []T) (equal bool, explination string) {
	if len(extraExpected) == 0 && len(extraActual) == 0 {
		return true, explination
	}

	formater := func(v T) string {
		return fmt.Sprintf("%#v", v)
	}

	var msg bytes.Buffer
	if len(extraExpected) > 0 {
		msg.WriteString("missing expected values:\n")
		for _, v := range extraExpected {
			msg.WriteString(formater(v) + "\n")
		}
	}

	if len(extraActual) > 0 {
		msg.WriteString("unexpected additional values:\n")
		for _, v := range extraActual {
			msg.WriteString(formater(v) + "\n")
		}
	}

	return false, msg.String()
}
