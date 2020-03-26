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

// Package transform provides translations for opentelemetry-go concepts and
// structures to otlp structures.
package transform

import (
	"errors"

	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

// ErrUnimplementedAgg is returned when a transformation of an unimplemented
// aggregator is attempted.
var ErrUnimplementedAgg = errors.New("unimplemented aggregator")

// Record transforms a Record into an OTLP Metric. An ErrUnimplementedAgg
// error is returned if the Record Aggregator is not supported.
func Record(r export.Record) (*metricpb.Metric, error) {
	d := r.Descriptor()
	l := r.Labels()
	switch a := r.Aggregator().(type) {
	case aggregator.MinMaxSumCount:
		return minMaxSumCount(d, l, a)
	case aggregator.Sum:
		return sum(d, l, a)
	}
	return nil, ErrUnimplementedAgg
}

// sum transforms a Sum Aggregator into an OTLP Metric.
func sum(desc *metric.Descriptor, labels export.Labels, a aggregator.Sum) (*metricpb.Metric, error) {
	sum, err := a.Sum()
	if err != nil {
		return nil, err
	}

	m := &metricpb.Metric{
		MetricDescriptor: &metricpb.MetricDescriptor{
			Name:        desc.Name(),
			Description: desc.Description(),
			Unit:        string(desc.Unit()),
			Labels:      stringKeyValues(labels.Iter()),
		},
	}

	switch n := desc.NumberKind(); n {
	case core.Int64NumberKind, core.Uint64NumberKind:
		m.MetricDescriptor.Type = metricpb.MetricDescriptor_COUNTER_INT64
		m.Int64DataPoints = []*metricpb.Int64DataPoint{
			{Value: sum.CoerceToInt64(n)},
		}
	case core.Float64NumberKind:
		m.MetricDescriptor.Type = metricpb.MetricDescriptor_COUNTER_DOUBLE
		m.DoubleDataPoints = []*metricpb.DoubleDataPoint{
			{Value: sum.CoerceToFloat64(n)},
		}
	}

	return m, nil
}

// minMaxSumCountValue returns the values of the MinMaxSumCount Aggregator
// as discret values.
func minMaxSumCountValues(a aggregator.MinMaxSumCount) (min, max, sum core.Number, count int64, err error) {
	if min, err = a.Min(); err != nil {
		return
	}
	if max, err = a.Max(); err != nil {
		return
	}
	if sum, err = a.Sum(); err != nil {
		return
	}
	if count, err = a.Count(); err != nil {
		return
	}
	return
}

// minMaxSumCount transforms a MinMaxSumCount Aggregator into an OTLP Metric.
func minMaxSumCount(desc *metric.Descriptor, labels export.Labels, a aggregator.MinMaxSumCount) (*metricpb.Metric, error) {
	min, max, sum, count, err := minMaxSumCountValues(a)
	if err != nil {
		return nil, err
	}

	numKind := desc.NumberKind()
	return &metricpb.Metric{
		MetricDescriptor: &metricpb.MetricDescriptor{
			Name:        desc.Name(),
			Description: desc.Description(),
			Unit:        string(desc.Unit()),
			Type:        metricpb.MetricDescriptor_SUMMARY,
			Labels:      stringKeyValues(labels.Iter()),
		},
		SummaryDataPoints: []*metricpb.SummaryDataPoint{
			{
				Count: uint64(count),
				Sum:   sum.CoerceToFloat64(numKind),
				PercentileValues: []*metricpb.SummaryDataPoint_ValueAtPercentile{
					{
						Percentile: 0.0,
						Value:      min.CoerceToFloat64(numKind),
					},
					{
						Percentile: 100.0,
						Value:      max.CoerceToFloat64(numKind),
					},
				},
			},
		},
	}, nil
}

// stringKeyValues transforms a label iterator into an OTLP StringKeyValues.
func stringKeyValues(iter export.LabelIterator) []*commonpb.StringKeyValue {
	l := iter.Len()
	if l == 0 {
		return nil
	}
	result := make([]*commonpb.StringKeyValue, 0, l)
	for iter.Next() {
		kv := iter.Label()
		result = append(result, &commonpb.StringKeyValue{
			Key:   string(kv.Key),
			Value: kv.Value.Emit(),
		})
	}
	return result
}
