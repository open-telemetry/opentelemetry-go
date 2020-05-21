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

package simple // import "go.opentelemetry.io/otel/sdk/metric/selector/simple"

import (
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

type (
	selectorInexpensive struct{}
	selectorExact       struct{}
	selectorSketch      struct {
		config *ddsketch.Config
	}
	selectorHistogram struct {
		boundaries []float64
	}
)

var (
	_ export.AggregationSelector = selectorInexpensive{}
	_ export.AggregationSelector = selectorSketch{}
	_ export.AggregationSelector = selectorExact{}
	_ export.AggregationSelector = selectorHistogram{}
)

// NewWithInexpensiveDistribution returns a simple aggregation selector
// that uses counter, minmaxsumcount and minmaxsumcount aggregators
// for the three kinds of metric.  This selector is faster and uses
// less memory than the others because minmaxsumcount does not
// aggregate quantile information.
func NewWithInexpensiveDistribution() export.AggregationSelector {
	return selectorInexpensive{}
}

// NewWithSketchDistribution returns a simple aggregation selector that
// uses counter, ddsketch, and ddsketch aggregators for the three
// kinds of metric.  This selector uses more cpu and memory than the
// NewWithInexpensiveDistribution because it uses one DDSketch per distinct
// instrument and label set.
func NewWithSketchDistribution(config *ddsketch.Config) export.AggregationSelector {
	return selectorSketch{
		config: config,
	}
}

// NewWithExactDistribution returns a simple aggregation selector that uses
// counter, array, and array aggregators for the three kinds of metric.
// This selector uses more memory than the NewWithSketchDistribution
// because it aggregates an array of all values, therefore is able to
// compute exact quantiles.
func NewWithExactDistribution() export.AggregationSelector {
	return selectorExact{}
}

// NewWithHistogramDistribution returns a simple aggregation selector that uses counter,
// histogram, and histogram aggregators for the three kinds of metric. This
// selector uses more memory than the NewWithInexpensiveDistribution because it
// uses a counter per bucket.
func NewWithHistogramDistribution(boundaries []float64) export.AggregationSelector {
	return selectorHistogram{boundaries: boundaries}
}

func (selectorInexpensive) AggregatorFor(descriptor *metric.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case metric.ValueObserverKind, metric.ValueRecorderKind:
		return minmaxsumcount.New(descriptor)
	default:
		return sum.New()
	}
}

func (s selectorSketch) AggregatorFor(descriptor *metric.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case metric.ValueObserverKind, metric.ValueRecorderKind:
		return ddsketch.New(s.config, descriptor)
	default:
		return sum.New()
	}
}

func (selectorExact) AggregatorFor(descriptor *metric.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case metric.ValueObserverKind, metric.ValueRecorderKind:
		return array.New()
	default:
		return sum.New()
	}
}

func (s selectorHistogram) AggregatorFor(descriptor *metric.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case metric.ValueObserverKind, metric.ValueRecorderKind:
		return histogram.New(descriptor, s.boundaries)
	default:
		return sum.New()
	}
}
