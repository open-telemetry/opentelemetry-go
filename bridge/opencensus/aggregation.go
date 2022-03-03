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

package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"

import (
	"errors"
	"fmt"
	"time"

	"go.opencensus.io/metric/metricdata"

	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
)

var (
	errIncompatibleType = errors.New("incompatible type for aggregation")
	errEmpty            = errors.New("points may not be empty")
	errBadPoint         = errors.New("point cannot be converted")
)

type recordFunc func(agg aggregation.Aggregation, end time.Time) error

// recordAggregationsFromPoints records one OpenTelemetry aggregation for
// each OpenCensus point.  Points may not be empty and must be either
// all (int|float)64 or all *metricdata.Distribution.
func recordAggregationsFromPoints(points []metricdata.Point, recorder recordFunc) error {
	if len(points) == 0 {
		return errEmpty
	}
	switch t := points[0].Value.(type) {
	case int64:
		return recordGaugePoints(points, recorder)
	case float64:
		return recordGaugePoints(points, recorder)
	case *metricdata.Distribution:
		return recordDistributionPoint(points, recorder)
	default:
		// TODO add *metricdata.Summary support
		return fmt.Errorf("%w: %v", errIncompatibleType, t)
	}
}

var _ aggregation.Aggregation = &ocRawAggregator{}
var _ aggregation.LastValue = &ocRawAggregator{}

// recordGaugePoints creates an OpenTelemetry aggregation from OpenCensus points.
// Points may not be empty, and must only contain integers or floats.
func recordGaugePoints(pts []metricdata.Point, recorder recordFunc) error {
	for _, pt := range pts {
		switch t := pt.Value.(type) {
		case int64:
			if err := recorder(&ocRawAggregator{
				value: number.NewInt64Number(pt.Value.(int64)),
				time:  pt.Time,
			}, pt.Time); err != nil {
				return err
			}
		case float64:
			if err := recorder(&ocRawAggregator{
				value: number.NewFloat64Number(pt.Value.(float64)),
				time:  pt.Time,
			}, pt.Time); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: %v", errIncompatibleType, t)
		}
	}
	return nil
}

type ocRawAggregator struct {
	value number.Number
	time  time.Time
}

// Kind returns the kind of aggregation this is.
func (o *ocRawAggregator) Kind() aggregation.Kind {
	return aggregation.LastValueKind
}

// LastValue returns the last point.
func (o *ocRawAggregator) LastValue() (number.Number, time.Time, error) {
	return o.value, o.time, nil
}

var _ aggregation.Aggregation = &ocDistAggregator{}
var _ aggregation.Histogram = &ocDistAggregator{}

// recordDistributionPoint creates an OpenTelemetry aggregation from
// OpenCensus points. Points may not be empty, and must only contain
// Distributions.  The most recent disribution will be used in the aggregation.
func recordDistributionPoint(pts []metricdata.Point, recorder recordFunc) error {
	// only use the most recent datapoint for now.
	pt := pts[len(pts)-1]
	val, ok := pt.Value.(*metricdata.Distribution)
	if !ok {
		return fmt.Errorf("%w: %v", errBadPoint, pt.Value)
	}
	bucketCounts := make([]uint64, len(val.Buckets))
	for i, bucket := range val.Buckets {
		if bucket.Count < 0 {
			return fmt.Errorf("%w: bucket count may not be negative", errBadPoint)
		}
		bucketCounts[i] = uint64(bucket.Count)
	}
	if val.Count < 0 {
		return fmt.Errorf("%w: count may not be negative", errBadPoint)
	}
	return recorder(&ocDistAggregator{
		sum:   number.NewFloat64Number(val.Sum),
		count: uint64(val.Count),
		buckets: aggregation.Buckets{
			Boundaries: val.BucketOptions.Bounds,
			Counts:     bucketCounts,
		},
	}, pts[len(pts)-1].Time)
}

type ocDistAggregator struct {
	sum     number.Number
	count   uint64
	buckets aggregation.Buckets
}

// Kind returns the kind of aggregation this is.
func (o *ocDistAggregator) Kind() aggregation.Kind {
	return aggregation.HistogramKind
}

// Sum returns the sum of values.
func (o *ocDistAggregator) Sum() (number.Number, error) {
	return o.sum, nil
}

// Count returns the number of values.
func (o *ocDistAggregator) Count() (uint64, error) {
	return o.count, nil
}

// Histogram returns the count of events in pre-determined buckets.
func (o *ocDistAggregator) Histogram() (aggregation.Buckets, error) {
	return o.buckets, nil
}
