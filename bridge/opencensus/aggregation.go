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

package opencensus

import (
	"errors"
	"fmt"
	"time"

	"go.opencensus.io/metric/metricdata"

	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

var (
	errIncompatibleType = errors.New("incompatible type for aggregation")
	errEmpty            = errors.New("points may not be empty")
	errBadPoint         = errors.New("point cannot be converted")
)

// aggregationWithEndTime is an aggregation that can also provide the timestamp
// of the last recorded point.
type aggregationWithEndTime interface {
	aggregation.Aggregation
	end() time.Time
}

// newAggregationFromPoints creates an OpenTelemetry aggregation from
// OpenCensus points.  Points may not be empty and must be either
// all (int|float)64 or all *metricdata.Distribution.
func newAggregationFromPoints(points []metricdata.Point) (aggregationWithEndTime, error) {
	if len(points) == 0 {
		return nil, errEmpty
	}
	switch t := points[0].Value.(type) {
	case int64:
		return newExactAggregator(points)
	case float64:
		return newExactAggregator(points)
	case *metricdata.Distribution:
		return newDistributionAggregator(points)
	default:
		// TODO add *metricdata.Summary support
		return nil, fmt.Errorf("%w: %v", errIncompatibleType, t)
	}
}

var _ aggregation.Aggregation = &ocExactAggregator{}
var _ aggregation.LastValue = &ocExactAggregator{}
var _ aggregation.Points = &ocExactAggregator{}

// newExactAggregator creates an OpenTelemetry aggreation from OpenCensus points.
// Points may not be empty, and must only contain integers or floats.
func newExactAggregator(pts []metricdata.Point) (aggregationWithEndTime, error) {
	points := make([]aggregation.Point, len(pts))
	for i, pt := range pts {
		switch t := pt.Value.(type) {
		case int64:
			points[i] = aggregation.Point{
				Number: number.NewInt64Number(pt.Value.(int64)),
				Time:   pt.Time,
			}
		case float64:
			points[i] = aggregation.Point{
				Number: number.NewFloat64Number(pt.Value.(float64)),
				Time:   pt.Time,
			}
		default:
			return nil, fmt.Errorf("%w: %v", errIncompatibleType, t)
		}
	}
	return &ocExactAggregator{
		points: points,
	}, nil
}

type ocExactAggregator struct {
	points []aggregation.Point
}

// Kind returns the kind of aggregation this is.
func (o *ocExactAggregator) Kind() aggregation.Kind {
	return aggregation.ExactKind
}

// Points returns access to the raw data set.
func (o *ocExactAggregator) Points() ([]aggregation.Point, error) {
	return o.points, nil
}

// LastValue returns the last point.
func (o *ocExactAggregator) LastValue() (number.Number, time.Time, error) {
	last := o.points[len(o.points)-1]
	return last.Number, last.Time, nil
}

// end returns the timestamp of the last point
func (o *ocExactAggregator) end() time.Time {
	_, t, _ := o.LastValue()
	return t
}

var _ aggregation.Aggregation = &ocDistAggregator{}
var _ aggregation.Histogram = &ocDistAggregator{}

// newDistributionAggregator creates an OpenTelemetry aggreation from
// OpenCensus points. Points may not be empty, and must only contain
// Distributions.  The most recent disribution will be used in the aggregation.
func newDistributionAggregator(pts []metricdata.Point) (aggregationWithEndTime, error) {
	// only use the most recent datapoint for now.
	pt := pts[len(pts)-1]
	val, ok := pt.Value.(*metricdata.Distribution)
	if !ok {
		return nil, fmt.Errorf("%w: %v", errBadPoint, pt.Value)
	}
	bucketCounts := make([]uint64, len(val.Buckets))
	for i, bucket := range val.Buckets {
		if bucket.Count < 0 {
			return nil, fmt.Errorf("%w: bucket count may not be negative", errBadPoint)
		}
		bucketCounts[i] = uint64(bucket.Count)
	}
	if val.Count < 0 {
		return nil, fmt.Errorf("%w: count may not be negative", errBadPoint)
	}
	return &ocDistAggregator{
		sum:   number.NewFloat64Number(val.Sum),
		count: uint64(val.Count),
		buckets: aggregation.Buckets{
			Boundaries: val.BucketOptions.Bounds,
			Counts:     bucketCounts,
		},
		endTime: pts[len(pts)-1].Time,
	}, nil
}

type ocDistAggregator struct {
	sum     number.Number
	count   uint64
	buckets aggregation.Buckets
	endTime time.Time
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

// end returns the time the histogram was measured.
func (o *ocDistAggregator) end() time.Time {
	return o.endTime
}
