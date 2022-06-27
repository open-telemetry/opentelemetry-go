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

// Package metrictransform provides translations for opentelemetry-go concepts and
// structures to otlp structures.
package metrictransform // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/metrictransform"

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/resource"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

var (
	// ErrUnimplementedAgg is returned when a transformation of an unimplemented
	// aggregator is attempted.
	ErrUnimplementedAgg = errors.New("unimplemented aggregator")

	// ErrIncompatibleAgg is returned when
	// aggregation.Kind implies an interface conversion that has
	// failed.
	ErrIncompatibleAgg = errors.New("incompatible aggregation type")

	// ErrUnknownValueType is returned when a transformation of an unknown value
	// is attempted.
	ErrUnknownValueType = errors.New("invalid value type")

	// ErrContextCanceled is returned when a context cancellation halts a
	// transformation.
	ErrContextCanceled = errors.New("context canceled")

	// ErrTransforming is returned when an unexected error is encountered transforming.
	ErrTransforming = errors.New("transforming failed")
)

// result is the product of transforming Records into OTLP Metrics.
type result struct {
	Metric *metricpb.Metric
	Err    error
}

// toNanos returns the number of nanoseconds since the UNIX epoch.
func toNanos(t time.Time) uint64 {
	if t.IsZero() {
		return 0
	}
	return uint64(t.UnixNano())
}

// InstrumentationLibraryReader transforms all records contained in a checkpoint into
// batched OTLP ResourceMetrics.
func InstrumentationLibraryReader(ctx context.Context, temporalitySelector aggregation.TemporalitySelector, res *resource.Resource, ilmr export.InstrumentationLibraryReader, numWorkers uint) (*metricpb.ResourceMetrics, error) {
	var sms []*metricpb.ScopeMetrics

	err := ilmr.ForEach(func(lib instrumentation.Library, mr export.Reader) error {
		records, errc := source(ctx, temporalitySelector, mr)

		// Start a fixed number of goroutines to transform records.
		transformed := make(chan result)
		var wg sync.WaitGroup
		wg.Add(int(numWorkers))
		for i := uint(0); i < numWorkers; i++ {
			go func() {
				defer wg.Done()
				transformer(ctx, temporalitySelector, records, transformed)
			}()
		}
		go func() {
			wg.Wait()
			close(transformed)
		}()

		// Synchronously collect the transformed records and transmit.
		ms, err := sink(ctx, transformed)
		if err != nil {
			return nil
		}

		// source is complete, check for any errors.
		if err := <-errc; err != nil {
			return err
		}
		if len(ms) == 0 {
			return nil
		}

		sms = append(sms, &metricpb.ScopeMetrics{
			Metrics:   ms,
			SchemaUrl: lib.SchemaURL,
			Scope: &commonpb.InstrumentationScope{
				Name:    lib.Name,
				Version: lib.Version,
			},
		})
		return nil
	})
	if len(sms) == 0 {
		return nil, err
	}

	rms := &metricpb.ResourceMetrics{
		Resource:     Resource(res),
		SchemaUrl:    res.SchemaURL(),
		ScopeMetrics: sms,
	}

	return rms, err
}

// source starts a goroutine that sends each one of the Records yielded by
// the Reader on the returned chan. Any error encountered will be sent
// on the returned error chan after seeding is complete.
func source(ctx context.Context, temporalitySelector aggregation.TemporalitySelector, mr export.Reader) (<-chan export.Record, <-chan error) {
	errc := make(chan error, 1)
	out := make(chan export.Record)
	// Seed records into process.
	go func() {
		defer close(out)
		// No select is needed since errc is buffered.
		errc <- mr.ForEach(temporalitySelector, func(r export.Record) error {
			select {
			case <-ctx.Done():
				return ErrContextCanceled
			case out <- r:
			}
			return nil
		})
	}()
	return out, errc
}

// transformer transforms records read from the passed in chan into
// OTLP Metrics which are sent on the out chan.
func transformer(ctx context.Context, temporalitySelector aggregation.TemporalitySelector, in <-chan export.Record, out chan<- result) {
	for r := range in {
		m, err := Record(temporalitySelector, r)
		// Propagate errors, but do not send empty results.
		if err == nil && m == nil {
			continue
		}
		res := result{
			Metric: m,
			Err:    err,
		}
		select {
		case <-ctx.Done():
			return
		case out <- res:
		}
	}
}

// sink collects transformed Records and batches them.
//
// Any errors encountered transforming input will be reported with an
// ErrTransforming as well as the completed ResourceMetrics. It is up to the
// caller to handle any incorrect data in these ResourceMetric.
func sink(ctx context.Context, in <-chan result) ([]*metricpb.Metric, error) {
	var errStrings []string

	// Group by the MetricDescriptor.
	grouped := map[string]*metricpb.Metric{}
	for res := range in {
		if res.Err != nil {
			errStrings = append(errStrings, res.Err.Error())
			continue
		}

		mID := res.Metric.GetName()
		m, ok := grouped[mID]
		if !ok {
			grouped[mID] = res.Metric
			continue
		}
		// Note: There is extra work happening in this code that can be
		// improved when the work described in #2119 is completed. The SDK has
		// a guarantee that no more than one point per period per attribute
		// set is produced, so this fallthrough should never happen. The final
		// step of #2119 is to remove all the grouping logic here.
		switch res.Metric.Data.(type) {
		case *metricpb.Metric_Gauge:
			m.GetGauge().DataPoints = append(m.GetGauge().DataPoints, res.Metric.GetGauge().DataPoints...)
		case *metricpb.Metric_Sum:
			m.GetSum().DataPoints = append(m.GetSum().DataPoints, res.Metric.GetSum().DataPoints...)
		case *metricpb.Metric_Histogram:
			m.GetHistogram().DataPoints = append(m.GetHistogram().DataPoints, res.Metric.GetHistogram().DataPoints...)
		case *metricpb.Metric_Summary:
			m.GetSummary().DataPoints = append(m.GetSummary().DataPoints, res.Metric.GetSummary().DataPoints...)
		default:
			err := fmt.Sprintf("unsupported metric type: %T", res.Metric.Data)
			errStrings = append(errStrings, err)
		}
	}

	if len(grouped) == 0 {
		return nil, nil
	}

	ms := make([]*metricpb.Metric, 0, len(grouped))
	for _, m := range grouped {
		ms = append(ms, m)
	}

	// Report any transform errors.
	if len(errStrings) > 0 {
		return ms, fmt.Errorf("%w:\n -%s", ErrTransforming, strings.Join(errStrings, "\n -"))
	}
	return ms, nil
}

// Record transforms a Record into an OTLP Metric. An ErrIncompatibleAgg
// error is returned if the Record Aggregator is not supported.
func Record(temporalitySelector aggregation.TemporalitySelector, r export.Record) (*metricpb.Metric, error) {
	agg := r.Aggregation()
	switch agg.Kind() {
	case aggregation.HistogramKind:
		h, ok := agg.(aggregation.Histogram)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrIncompatibleAgg, agg)
		}
		return histogramPoint(r, temporalitySelector.TemporalityFor(r.Descriptor(), aggregation.HistogramKind), h)

	case aggregation.SumKind:
		s, ok := agg.(aggregation.Sum)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrIncompatibleAgg, agg)
		}
		sum, err := s.Sum()
		if err != nil {
			return nil, err
		}
		return sumPoint(r, sum, r.StartTime(), r.EndTime(), temporalitySelector.TemporalityFor(r.Descriptor(), aggregation.SumKind), r.Descriptor().InstrumentKind().Monotonic())

	case aggregation.LastValueKind:
		lv, ok := agg.(aggregation.LastValue)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrIncompatibleAgg, agg)
		}
		value, tm, err := lv.LastValue()
		if err != nil {
			return nil, err
		}
		return gaugePoint(r, value, time.Time{}, tm)

	default:
		return nil, fmt.Errorf("%w: %T", ErrUnimplementedAgg, agg)
	}
}

func gaugePoint(record export.Record, num number.Number, start, end time.Time) (*metricpb.Metric, error) {
	desc := record.Descriptor()
	attrs := record.Attributes()

	m := &metricpb.Metric{
		Name:        desc.Name(),
		Description: desc.Description(),
		Unit:        string(desc.Unit()),
	}

	switch n := desc.NumberKind(); n {
	case number.Int64Kind:
		m.Data = &metricpb.Metric_Gauge{
			Gauge: &metricpb.Gauge{
				DataPoints: []*metricpb.NumberDataPoint{
					{
						Value: &metricpb.NumberDataPoint_AsInt{
							AsInt: num.CoerceToInt64(n),
						},
						Attributes:        Iterator(attrs.Iter()),
						StartTimeUnixNano: toNanos(start),
						TimeUnixNano:      toNanos(end),
					},
				},
			},
		}
	case number.Float64Kind:
		m.Data = &metricpb.Metric_Gauge{
			Gauge: &metricpb.Gauge{
				DataPoints: []*metricpb.NumberDataPoint{
					{
						Value: &metricpb.NumberDataPoint_AsDouble{
							AsDouble: num.CoerceToFloat64(n),
						},
						Attributes:        Iterator(attrs.Iter()),
						StartTimeUnixNano: toNanos(start),
						TimeUnixNano:      toNanos(end),
					},
				},
			},
		}
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownValueType, n)
	}

	return m, nil
}

func sdkTemporalityToTemporality(temporality aggregation.Temporality) metricpb.AggregationTemporality {
	switch temporality {
	case aggregation.DeltaTemporality:
		return metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA
	case aggregation.CumulativeTemporality:
		return metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE
	}
	return metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_UNSPECIFIED
}

func sumPoint(record export.Record, num number.Number, start, end time.Time, temporality aggregation.Temporality, monotonic bool) (*metricpb.Metric, error) {
	desc := record.Descriptor()
	attrs := record.Attributes()

	m := &metricpb.Metric{
		Name:        desc.Name(),
		Description: desc.Description(),
		Unit:        string(desc.Unit()),
	}

	switch n := desc.NumberKind(); n {
	case number.Int64Kind:
		m.Data = &metricpb.Metric_Sum{
			Sum: &metricpb.Sum{
				IsMonotonic:            monotonic,
				AggregationTemporality: sdkTemporalityToTemporality(temporality),
				DataPoints: []*metricpb.NumberDataPoint{
					{
						Value: &metricpb.NumberDataPoint_AsInt{
							AsInt: num.CoerceToInt64(n),
						},
						Attributes:        Iterator(attrs.Iter()),
						StartTimeUnixNano: toNanos(start),
						TimeUnixNano:      toNanos(end),
					},
				},
			},
		}
	case number.Float64Kind:
		m.Data = &metricpb.Metric_Sum{
			Sum: &metricpb.Sum{
				IsMonotonic:            monotonic,
				AggregationTemporality: sdkTemporalityToTemporality(temporality),
				DataPoints: []*metricpb.NumberDataPoint{
					{
						Value: &metricpb.NumberDataPoint_AsDouble{
							AsDouble: num.CoerceToFloat64(n),
						},
						Attributes:        Iterator(attrs.Iter()),
						StartTimeUnixNano: toNanos(start),
						TimeUnixNano:      toNanos(end),
					},
				},
			},
		}
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownValueType, n)
	}

	return m, nil
}

func histogramValues(a aggregation.Histogram) (boundaries []float64, counts []uint64, err error) {
	var buckets aggregation.Buckets
	if buckets, err = a.Histogram(); err != nil {
		return
	}
	boundaries, counts = buckets.Boundaries, buckets.Counts
	if len(counts) != len(boundaries)+1 {
		err = ErrTransforming
		return
	}
	return
}

// histogram transforms a Histogram Aggregator into an OTLP Metric.
func histogramPoint(record export.Record, temporality aggregation.Temporality, a aggregation.Histogram) (*metricpb.Metric, error) {
	desc := record.Descriptor()
	attrs := record.Attributes()
	boundaries, counts, err := histogramValues(a)
	if err != nil {
		return nil, err
	}

	count, err := a.Count()
	if err != nil {
		return nil, err
	}

	sum, err := a.Sum()
	if err != nil {
		return nil, err
	}

	sumFloat64 := sum.CoerceToFloat64(desc.NumberKind())
	m := &metricpb.Metric{
		Name:        desc.Name(),
		Description: desc.Description(),
		Unit:        string(desc.Unit()),
		Data: &metricpb.Metric_Histogram{
			Histogram: &metricpb.Histogram{
				AggregationTemporality: sdkTemporalityToTemporality(temporality),
				DataPoints: []*metricpb.HistogramDataPoint{
					{
						Sum:               &sumFloat64,
						Attributes:        Iterator(attrs.Iter()),
						StartTimeUnixNano: toNanos(record.StartTime()),
						TimeUnixNano:      toNanos(record.EndTime()),
						Count:             uint64(count),
						BucketCounts:      counts,
						ExplicitBounds:    boundaries,
					},
				},
			},
		},
	}
	return m, nil
}
