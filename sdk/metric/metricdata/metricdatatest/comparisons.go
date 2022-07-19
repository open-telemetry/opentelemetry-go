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

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"bytes"
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
) // equalResourceMetrics returns true when ResourceMetrics are equal. It
// returns false when they differ, along with the reasons why they differ.
//
// The ScopeMetrics each ResourceMetrics contains are compared based on
// containing the same ScopeMetrics, not the order they are stored in.
func equalResourceMetrics(a, b metricdata.ResourceMetrics) (equal bool, reasons []string) {
	equal = true
	if !a.Resource.Equal(b.Resource) {
		equal, reasons = false, append(
			reasons, notEqualStr("Resources", a.Resource, b.Resource),
		)
	}

	var r string
	equal, r = compareDiff(diffSlices(
		a.ScopeMetrics,
		b.ScopeMetrics,
		func(a, b metricdata.ScopeMetrics) bool {
			equal, _ := equalScopeMetrics(a, b)
			return equal
		},
	))
	if !equal {
		reasons = append(reasons, fmt.Sprintf(
			"ResourceMetrics ScopeMetrics not equal:\n%s", r,
		))
	}
	return equal, reasons
}

// equalScopeMetrics returns true when ScopeMetrics are equal. It returns
// false when they differ, along with the reasons why they differ.
//
// The Metrics each ScopeMetrics contains are compared based on containing the
// same Metrics, not the order they are stored in.
func equalScopeMetrics(a, b metricdata.ScopeMetrics) (equal bool, reasons []string) {
	equal = true
	if a.Scope != b.Scope {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Scope", a.Scope, b.Scope),
		)
	}

	var r string
	equal, r = compareDiff(diffSlices(
		a.Metrics,
		b.Metrics,
		func(a, b metricdata.Metrics) bool {
			equal, _ := equalMetrics(a, b)
			return equal
		},
	))
	if !equal {
		reasons = append(reasons, fmt.Sprintf(
			"ScopeMetrics Metrics not equal:\n%s", r,
		))
	}
	return equal, reasons
}

// equalMetrics returns true when Metrics are equal. It returns false when
// they differ, along with the reasons why they differ.
func equalMetrics(a, b metricdata.Metrics) (equal bool, reasons []string) {
	equal = true
	if a.Name != b.Name {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Name", a.Name, b.Name),
		)
	}
	if a.Description != b.Description {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Description", a.Description, b.Description),
		)
	}
	if a.Unit != b.Unit {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Unit", a.Unit, b.Unit),
		)
	}

	var r []string
	equal, r = equalAggregations(a.Data, b.Data)
	if !equal {
		reasons = append(reasons, "Metrics Data not equal:")
		reasons = append(reasons, r...)
	}
	return equal, reasons
}

// equalAggregations returns true when a and b are equal. It returns false
// when they differ, along with the reasons why they differ.
func equalAggregations(a, b metricdata.Aggregation) (equal bool, reasons []string) {
	equal = true
	if a == nil || b == nil {
		if a != b {
			equal, reasons = false, []string{notEqualStr("Aggregation", a, b)}
		}
		return equal, reasons
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, []string{fmt.Sprintf(
			"Aggregation types not equal:\nexpected: %T\nactual: %T", a, b,
		)}
	}

	switch v := a.(type) {
	case metricdata.Gauge:
		var r []string
		equal, r = equalGauges(v, b.(metricdata.Gauge))
		if !equal {
			reasons = append(reasons, "Gauge not equal:")
			reasons = append(reasons, r...)
		}
	case metricdata.Sum:
		var r []string
		equal, r = equalSums(v, b.(metricdata.Sum))
		if !equal {
			reasons = append(reasons, "Sum not equal:")
			reasons = append(reasons, r...)
		}
	case metricdata.Histogram:
		var r []string
		equal, r = equalHistograms(v, b.(metricdata.Histogram))
		if !equal {
			reasons = append(reasons, "Histogram not equal:")
			reasons = append(reasons, r...)
		}
	default:
		equal = false
		reasons = append(reasons, fmt.Sprintf("Aggregation of unknown types %T", a))
	}
	return equal, reasons
}

// equalGauges returns true when Gauges are equal. It returns false when they
// differ, along with the reasons why they differ.
//
// The DataPoints each Gauge contains are compared based on containing the
// same DataPoints, not the order they are stored in.
func equalGauges(a, b metricdata.Gauge) (equal bool, reasons []string) {
	var r string
	equal, r = compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.DataPoint) bool {
			equal, _ := equalDataPoints(a, b)
			return equal
		},
	))
	if !equal {
		reasons = append(reasons, fmt.Sprintf(
			"Gauge DataPoints not equal:\n%s", r,
		))
	}
	return equal, reasons
}

// equalSums returns true when Sums are equal. It returns false when they
// differ, along with the reasons why they differ.
//
// The DataPoints each Sum contains are compared based on containing the same
// DataPoints, not the order they are stored in.
func equalSums(a, b metricdata.Sum) (equal bool, reasons []string) {
	equal = true
	if a.Temporality != b.Temporality {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Temporality", a.Temporality, b.Temporality),
		)
	}
	if a.IsMonotonic != b.IsMonotonic {
		equal, reasons = false, append(
			reasons,
			notEqualStr("IsMonotonic", a.IsMonotonic, b.IsMonotonic),
		)
	}

	var r string
	equal, r = compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.DataPoint) bool {
			equal, _ := equalDataPoints(a, b)
			return equal
		},
	))
	if !equal {
		reasons = append(reasons, fmt.Sprintf(
			"Sum DataPoints not equal:\n%s", r,
		))
	}
	return equal, reasons
}

// equalHistograms returns true when Histograms are equal. It returns false
// when they differ, along with the reasons why they differ.
//
// The DataPoints each Histogram contains are compared based on containing the
// same HistogramDataPoint, not the order they are stored in.
func equalHistograms(a, b metricdata.Histogram) (equal bool, reasons []string) {
	equal = true
	if a.Temporality != b.Temporality {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Temporality", a.Temporality, b.Temporality),
		)
	}

	var r string
	equal, r = compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.HistogramDataPoint) bool {
			equal, _ := equalHistogramDataPoints(a, b)
			return equal
		},
	))
	if !equal {
		reasons = append(reasons, fmt.Sprintf(
			"Histogram DataPoints not equal:\n%s", r,
		))
	}
	return equal, reasons
}

// equalDataPoints returns true when DataPoints are equal. It returns false
// when they differ, along with the reasons why they differ.
func equalDataPoints(a, b metricdata.DataPoint) (equal bool, reasons []string) {
	equal = true
	if !a.Attributes.Equals(&b.Attributes) {
		equal, reasons = false, append(reasons, notEqualStr(
			"Attributes",
			a.Attributes.Encoded(attribute.DefaultEncoder()),
			b.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}
	if !a.StartTime.Equal(b.StartTime) {
		equal, reasons = false, append(reasons, notEqualStr(
			"StartTime",
			a.StartTime.UnixNano(),
			b.StartTime.UnixNano(),
		))
	}
	if !a.Time.Equal(b.Time) {
		equal, reasons = false, append(reasons, notEqualStr(
			"Time",
			a.Time.UnixNano(),
			b.Time.UnixNano(),
		))
	}

	var r []string
	equal, r = equalValues(a.Value, b.Value)
	if !equal {
		reasons = append(reasons, "DataPoint Value not equal:")
		reasons = append(reasons, r...)
	}
	return equal, reasons
}

// equalHistogramDataPoints returns true when HistogramDataPoints are equal.
// It returns false when they differ, along with the reasons why they differ.
func equalHistogramDataPoints(a, b metricdata.HistogramDataPoint) (equal bool, reasons []string) {
	equal = true
	if !a.Attributes.Equals(&b.Attributes) {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"Attributes not equal:\nexpected: %s\nactual: %s",
			a.Attributes.Encoded(attribute.DefaultEncoder()),
			b.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}
	if !a.StartTime.Equal(b.StartTime) {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"StartTime not equal:\nexpected: %d\nactual: %d",
			a.StartTime.UnixNano(),
			b.StartTime.UnixNano(),
		))
	}
	if !a.Time.Equal(b.Time) {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"Time not equal:\nexpected: %d\nactual: %d",
			a.Time.UnixNano(),
			b.Time.UnixNano(),
		))
	}
	if a.Count != b.Count {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"Count not equal:\nexpected: %d\nactual: %d",
			a.Count,
			b.Count,
		))
	}
	if !equalSlices(a.Bounds, b.Bounds) {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"Bounds not equal:\nexpected: %v\nactual: %v",
			a.Bounds,
			b.Bounds,
		))
	}
	if !equalSlices(a.BucketCounts, b.BucketCounts) {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"BucketCounts not equal:\nexpected: %v\nactual: %v",
			a.BucketCounts,
			b.BucketCounts,
		))
	}
	if !equalPtrValues(a.Min, b.Min) {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"Min not equal:\nexpected: %v\nactual: %v",
			a.Min,
			b.Min,
		))
	}
	if !equalPtrValues(a.Max, b.Max) {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"Max not equal:\nexpected: %v\nactual: %v",
			a.Max,
			b.Max,
		))
	}
	if a.Sum != b.Sum {
		equal, reasons = false, append(reasons, fmt.Sprintf(
			"Sum not equal:\nexpected: %g\nactual: %g",
			a.Sum,
			b.Sum,
		))
	}
	return equal, reasons
}

// equalValues returns true when Values are equal. It returns false when they
// differ, along with the reasons why they differ.
func equalValues(a, b metricdata.Value) (equal bool, reasons []string) {
	equal = true
	if a == nil || b == nil {
		if a != b {
			equal, reasons = false, []string{notEqualStr("Values", a, b)}
		}
		return equal, reasons
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, []string{fmt.Sprintf(
			"Value types not equal:\nexpected: %T\nactual: %T", a, b,
		)}
	}

	switch v := a.(type) {
	case metricdata.Int64:
		var r []string
		equal, r = equalInt64(v, b.(metricdata.Int64))
		if !equal {
			reasons = append(reasons, "Int64 not equal:")
			reasons = append(reasons, r...)
		}
	case metricdata.Float64:
		var r []string
		equal, r = equalFloat64(v, b.(metricdata.Float64))
		if !equal {
			reasons = append(reasons, "Int64 not equal:")
			reasons = append(reasons, r...)
		}
	default:
		equal = false
		reasons = append(reasons, fmt.Sprintf("Value of unknown types %T", a))
	}

	return equal, reasons
}

// equalFloat64 returns true when Float64s are equal. It returns false when
// they differ, along with the reasons why they differ.
func equalFloat64(a, b metricdata.Float64) (equal bool, reasons []string) {
	equal = a == b
	if !equal {
		reasons = append(reasons, notEqualStr("Float64 value", a, b))
	}
	return equal, reasons
}

// equalInt64 returns true when Int64s are equal. It returns false when they
// differ, along with the reasons why they differ.
func equalInt64(a, b metricdata.Int64) (equal bool, reasons []string) {
	equal = a == b
	if !equal {
		reasons = append(reasons, notEqualStr("Int64 value", a, b))
	}
	return equal, reasons
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

func compareDiff[T any](extraExpected, extraActual []T) (equal bool, reasons string) {
	if len(extraExpected) == 0 && len(extraActual) == 0 {
		return true, reasons
	}

	formater := func(v T) string {
		return fmt.Sprintf("%#v", v)
	}

	var msg bytes.Buffer
	if len(extraExpected) > 0 {
		_, _ = msg.WriteString("missing expected values:\n")
		for _, v := range extraExpected {
			_, _ = msg.WriteString(formater(v) + "\n")
		}
	}

	if len(extraActual) > 0 {
		_, _ = msg.WriteString("unexpected additional values:\n")
		for _, v := range extraActual {
			_, _ = msg.WriteString(formater(v) + "\n")
		}
	}

	return false, msg.String()
}
