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

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"bytes"
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// equalResourceMetrics returns reasons ResourceMetrics are not equal. If they
// are equal, the returned reasons will be empty.
//
// The ScopeMetrics each ResourceMetrics contains are compared based on
// containing the same ScopeMetrics, not the order they are stored in.
func equalResourceMetrics(a, b metricdata.ResourceMetrics, cfg config) (reasons []string) {
	if !a.Resource.Equal(b.Resource) {
		reasons = append(reasons, notEqualStr("Resources", a.Resource, b.Resource))
	}

	r := compareDiff(diffSlices(
		a.ScopeMetrics,
		b.ScopeMetrics,
		func(a, b metricdata.ScopeMetrics) bool {
			r := equalScopeMetrics(a, b, cfg)
			return len(r) == 0
		},
	))
	if r != "" {
		reasons = append(reasons, fmt.Sprintf("ResourceMetrics ScopeMetrics not equal:\n%s", r))
	}
	return reasons
}

// equalScopeMetrics returns reasons ScopeMetrics are not equal. If they are
// equal, the returned reasons will be empty.
//
// The Metrics each ScopeMetrics contains are compared based on containing the
// same Metrics, not the order they are stored in.
func equalScopeMetrics(a, b metricdata.ScopeMetrics, cfg config) (reasons []string) {
	if a.Scope != b.Scope {
		reasons = append(reasons, notEqualStr("Scope", a.Scope, b.Scope))
	}

	r := compareDiff(diffSlices(
		a.Metrics,
		b.Metrics,
		func(a, b metricdata.Metrics) bool {
			r := equalMetrics(a, b, cfg)
			return len(r) == 0
		},
	))
	if r != "" {
		reasons = append(reasons, fmt.Sprintf("ScopeMetrics Metrics not equal:\n%s", r))
	}
	return reasons
}

// equalMetrics returns reasons Metrics are not equal. If they are equal, the
// returned reasons will be empty.
func equalMetrics(a, b metricdata.Metrics, cfg config) (reasons []string) {
	if a.Name != b.Name {
		reasons = append(reasons, notEqualStr("Name", a.Name, b.Name))
	}
	if a.Description != b.Description {
		reasons = append(reasons, notEqualStr("Description", a.Description, b.Description))
	}
	if a.Unit != b.Unit {
		reasons = append(reasons, notEqualStr("Unit", a.Unit, b.Unit))
	}

	r := equalAggregations(a.Data, b.Data, cfg)
	if len(r) > 0 {
		reasons = append(reasons, "Metrics Data not equal:")
		reasons = append(reasons, r...)
	}
	return reasons
}

// equalAggregations returns reasons a and b are not equal. If they are equal,
// the returned reasons will be empty.
func equalAggregations(a, b metricdata.Aggregation, cfg config) (reasons []string) {
	if a == nil || b == nil {
		if a != b {
			return []string{notEqualStr("Aggregation", a, b)}
		}
		return reasons
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return []string{fmt.Sprintf("Aggregation types not equal:\nexpected: %T\nactual: %T", a, b)}
	}

	switch v := a.(type) {
	case metricdata.Gauge[int64]:
		r := equalGauges(v, b.(metricdata.Gauge[int64]), cfg)
		if len(r) > 0 {
			reasons = append(reasons, "Gauge[int64] not equal:")
			reasons = append(reasons, r...)
		}
	case metricdata.Gauge[float64]:
		r := equalGauges(v, b.(metricdata.Gauge[float64]), cfg)
		if len(r) > 0 {
			reasons = append(reasons, "Gauge[float64] not equal:")
			reasons = append(reasons, r...)
		}
	case metricdata.Sum[int64]:
		r := equalSums(v, b.(metricdata.Sum[int64]), cfg)
		if len(r) > 0 {
			reasons = append(reasons, "Sum[int64] not equal:")
			reasons = append(reasons, r...)
		}
	case metricdata.Sum[float64]:
		r := equalSums(v, b.(metricdata.Sum[float64]), cfg)
		if len(r) > 0 {
			reasons = append(reasons, "Sum[float64] not equal:")
			reasons = append(reasons, r...)
		}
	case metricdata.Histogram:
		r := equalHistograms(v, b.(metricdata.Histogram), cfg)
		if len(r) > 0 {
			reasons = append(reasons, "Histogram not equal:")
			reasons = append(reasons, r...)
		}
	default:
		reasons = append(reasons, fmt.Sprintf("Aggregation of unknown types %T", a))
	}
	return reasons
}

// equalGauges returns reasons Gauges are not equal. If they are equal, the
// returned reasons will be empty.
//
// The DataPoints each Gauge contains are compared based on containing the
// same DataPoints, not the order they are stored in.
func equalGauges[N int64 | float64](a, b metricdata.Gauge[N], cfg config) (reasons []string) {
	r := compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.DataPoint[N]) bool {
			r := equalDataPoints(a, b, cfg)
			return len(r) == 0
		},
	))
	if r != "" {
		reasons = append(reasons, fmt.Sprintf("Gauge DataPoints not equal:\n%s", r))
	}
	return reasons
}

// equalSums returns reasons Sums are not equal. If they are equal, the
// returned reasons will be empty.
//
// The DataPoints each Sum contains are compared based on containing the same
// DataPoints, not the order they are stored in.
func equalSums[N int64 | float64](a, b metricdata.Sum[N], cfg config) (reasons []string) {
	if a.Temporality != b.Temporality {
		reasons = append(reasons, notEqualStr("Temporality", a.Temporality, b.Temporality))
	}
	if a.IsMonotonic != b.IsMonotonic {
		reasons = append(reasons, notEqualStr("IsMonotonic", a.IsMonotonic, b.IsMonotonic))
	}

	r := compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.DataPoint[N]) bool {
			r := equalDataPoints(a, b, cfg)
			return len(r) == 0
		},
	))
	if r != "" {
		reasons = append(reasons, fmt.Sprintf("Sum DataPoints not equal:\n%s", r))
	}
	return reasons
}

// equalHistograms returns reasons Histograms are not equal. If they are
// equal, the returned reasons will be empty.
//
// The DataPoints each Histogram contains are compared based on containing the
// same HistogramDataPoint, not the order they are stored in.
func equalHistograms(a, b metricdata.Histogram, cfg config) (reasons []string) {
	if a.Temporality != b.Temporality {
		reasons = append(reasons, notEqualStr("Temporality", a.Temporality, b.Temporality))
	}

	r := compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.HistogramDataPoint) bool {
			r := equalHistogramDataPoints(a, b, cfg)
			return len(r) == 0
		},
	))
	if r != "" {
		reasons = append(reasons, fmt.Sprintf("Histogram DataPoints not equal:\n%s", r))
	}
	return reasons
}

// equalDataPoints returns reasons DataPoints are not equal. If they are
// equal, the returned reasons will be empty.
func equalDataPoints[N int64 | float64](a, b metricdata.DataPoint[N], cfg config) (reasons []string) { // nolint: revive // Intentional internal control flag
	if !a.Attributes.Equals(&b.Attributes) {
		reasons = append(reasons, notEqualStr(
			"Attributes",
			a.Attributes.Encoded(attribute.DefaultEncoder()),
			b.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}

	if !cfg.ignoreTimestamp {
		if !a.StartTime.Equal(b.StartTime) {
			reasons = append(reasons, notEqualStr("StartTime", a.StartTime.UnixNano(), b.StartTime.UnixNano()))
		}
		if !a.Time.Equal(b.Time) {
			reasons = append(reasons, notEqualStr("Time", a.Time.UnixNano(), b.Time.UnixNano()))
		}
	}

	if a.Value != b.Value {
		reasons = append(reasons, notEqualStr("Value", a.Value, b.Value))
	}
	return reasons
}

// equalHistogramDataPoints returns reasons HistogramDataPoints are not equal.
// If they are equal, the returned reasons will be empty.
func equalHistogramDataPoints(a, b metricdata.HistogramDataPoint, cfg config) (reasons []string) { // nolint: revive // Intentional internal control flag
	if !a.Attributes.Equals(&b.Attributes) {
		reasons = append(reasons, notEqualStr(
			"Attributes",
			a.Attributes.Encoded(attribute.DefaultEncoder()),
			b.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}
	if !cfg.ignoreTimestamp {
		if !a.StartTime.Equal(b.StartTime) {
			reasons = append(reasons, notEqualStr("StartTime", a.StartTime.UnixNano(), b.StartTime.UnixNano()))
		}
		if !a.Time.Equal(b.Time) {
			reasons = append(reasons, notEqualStr("Time", a.Time.UnixNano(), b.Time.UnixNano()))
		}
	}
	if a.Count != b.Count {
		reasons = append(reasons, notEqualStr("Count", a.Count, b.Count))
	}
	if !equalSlices(a.Bounds, b.Bounds) {
		reasons = append(reasons, notEqualStr("Bounds", a.Bounds, b.Bounds))
	}
	if !equalSlices(a.BucketCounts, b.BucketCounts) {
		reasons = append(reasons, notEqualStr("BucketCounts", a.BucketCounts, b.BucketCounts))
	}
	if !eqExtrema(a.Min, b.Min) {
		reasons = append(reasons, notEqualStr("Min", a.Min, b.Min))
	}
	if !eqExtrema(a.Max, b.Max) {
		reasons = append(reasons, notEqualStr("Max", a.Max, b.Max))
	}
	if a.Sum != b.Sum {
		reasons = append(reasons, notEqualStr("Sum", a.Sum, b.Sum))
	}
	return reasons
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

func equalExtrema(a, b metricdata.Extrema, _ config) (reasons []string) {
	if !eqExtrema(a, b) {
		reasons = append(reasons, notEqualStr("Extrema", a, b))
	}
	return reasons
}

func eqExtrema(a, b metricdata.Extrema) bool {
	aV, aOk := a.Value()
	bV, bOk := b.Value()

	if !aOk || !bOk {
		return aOk == bOk
	}
	return aV == bV
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

func compareDiff[T any](extraExpected, extraActual []T) string {
	if len(extraExpected) == 0 && len(extraActual) == 0 {
		return ""
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

	return msg.String()
}

func missingAttrStr(name string) string {
	return fmt.Sprintf("missing attribute %s", name)
}

func hasAttributesDataPoints[T int64 | float64](dp metricdata.DataPoint[T], attrs ...attribute.KeyValue) (reasons []string) {
	for _, attr := range attrs {
		val, ok := dp.Attributes.Value(attr.Key)
		if !ok {
			reasons = append(reasons, missingAttrStr(string(attr.Key)))
			continue
		}
		if val != attr.Value {
			reasons = append(reasons, notEqualStr(string(attr.Key), attr.Value.Emit(), val.Emit()))
		}
	}
	return reasons
}

func hasAttributesGauge[T int64 | float64](gauge metricdata.Gauge[T], attrs ...attribute.KeyValue) (reasons []string) {
	for n, dp := range gauge.DataPoints {
		reas := hasAttributesDataPoints(dp, attrs...)
		if len(reas) > 0 {
			reasons = append(reasons, fmt.Sprintf("gauge datapoint %d attributes:\n", n))
			reasons = append(reasons, reas...)
		}
	}
	return reasons
}

func hasAttributesSum[T int64 | float64](sum metricdata.Sum[T], attrs ...attribute.KeyValue) (reasons []string) {
	for n, dp := range sum.DataPoints {
		reas := hasAttributesDataPoints(dp, attrs...)
		if len(reas) > 0 {
			reasons = append(reasons, fmt.Sprintf("sum datapoint %d attributes:\n", n))
			reasons = append(reasons, reas...)
		}
	}
	return reasons
}

func hasAttributesHistogramDataPoints(dp metricdata.HistogramDataPoint, attrs ...attribute.KeyValue) (reasons []string) {
	for _, attr := range attrs {
		val, ok := dp.Attributes.Value(attr.Key)
		if !ok {
			reasons = append(reasons, missingAttrStr(string(attr.Key)))
			continue
		}
		if val != attr.Value {
			reasons = append(reasons, notEqualStr(string(attr.Key), attr.Value.Emit(), val.Emit()))
		}
	}
	return reasons
}

func hasAttributesHistogram(histogram metricdata.Histogram, attrs ...attribute.KeyValue) (reasons []string) {
	for n, dp := range histogram.DataPoints {
		reas := hasAttributesHistogramDataPoints(dp, attrs...)
		if len(reas) > 0 {
			reasons = append(reasons, fmt.Sprintf("histogram datapoint %d attributes:\n", n))
			reasons = append(reasons, reas...)
		}
	}
	return reasons
}

func hasAttributesAggregation(agg metricdata.Aggregation, attrs ...attribute.KeyValue) (reasons []string) {
	switch agg := agg.(type) {
	case metricdata.Gauge[int64]:
		reasons = hasAttributesGauge(agg, attrs...)
	case metricdata.Gauge[float64]:
		reasons = hasAttributesGauge(agg, attrs...)
	case metricdata.Sum[int64]:
		reasons = hasAttributesSum(agg, attrs...)
	case metricdata.Sum[float64]:
		reasons = hasAttributesSum(agg, attrs...)
	case metricdata.Histogram:
		reasons = hasAttributesHistogram(agg, attrs...)
	default:
		reasons = []string{fmt.Sprintf("unknown aggregation %T", agg)}
	}
	return reasons
}

func hasAttributesMetrics(metrics metricdata.Metrics, attrs ...attribute.KeyValue) (reasons []string) {
	reas := hasAttributesAggregation(metrics.Data, attrs...)
	if len(reas) > 0 {
		reasons = append(reasons, fmt.Sprintf("Metric %s:\n", metrics.Name))
		reasons = append(reasons, reas...)
	}
	return reasons
}

func hasAttributesScopeMetrics(sm metricdata.ScopeMetrics, attrs ...attribute.KeyValue) (reasons []string) {
	for n, metrics := range sm.Metrics {
		reas := hasAttributesMetrics(metrics, attrs...)
		if len(reas) > 0 {
			reasons = append(reasons, fmt.Sprintf("ScopeMetrics %s Metrics %d:\n", sm.Scope.Name, n))
			reasons = append(reasons, reas...)
		}
	}
	return reasons
}
func hasAttributesResourceMetrics(rm metricdata.ResourceMetrics, attrs ...attribute.KeyValue) (reasons []string) {
	for n, sm := range rm.ScopeMetrics {
		reas := hasAttributesScopeMetrics(sm, attrs...)
		if len(reas) > 0 {
			reasons = append(reasons, fmt.Sprintf("ResourceMetrics ScopeMetrics %d:\n", n))
			reasons = append(reasons, reas...)
		}
	}
	return reasons
}
