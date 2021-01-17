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
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exact"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

type (
	selectorInexpensive struct{}
	selectorExact       struct{}
	selectorHistogram   struct {
		options []histogram.Option
	}
)

var (
	_ metric.AggregatorSelector = selectorInexpensive{}
	_ metric.AggregatorSelector = selectorExact{}
	_ metric.AggregatorSelector = selectorHistogram{}
)

// NewWithInexpensiveDistribution returns a simple aggregator selector
// that uses minmaxsumcount aggregators for `ValueRecorder`
// instruments.  This selector is faster and uses less memory than the
// others in this package because minmaxsumcount aggregators maintain
// the least information about the distribution among these choices.
func NewWithInexpensiveDistribution() metric.AggregatorSelector {
	return selectorInexpensive{}
}

// NewWithExactDistribution returns a simple aggregator selector that
// uses exact aggregators for `ValueRecorder` instruments.  This
// selector uses more memory than the others in this package because
// exact aggregators maintain the most information about the
// distribution among these choices.
func NewWithExactDistribution() metric.AggregatorSelector {
	return selectorExact{}
}

// NewWithHistogramDistribution returns a simple aggregator selector
// that uses histogram aggregators for `ValueRecorder` instruments.
// This selector is a good default choice for most metric exporters.
func NewWithHistogramDistribution(options ...histogram.Option) metric.AggregatorSelector {
	return selectorHistogram{options: options}
}

func sumAggs(aggPtrs []*metric.Aggregator) {
	aggs := sum.New(len(aggPtrs))
	for i := range aggPtrs {
		*aggPtrs[i] = &aggs[i]
	}
}

func lastValueAggs(aggPtrs []*metric.Aggregator) {
	aggs := lastvalue.New(len(aggPtrs))
	for i := range aggPtrs {
		*aggPtrs[i] = &aggs[i]
	}
}

func (selectorInexpensive) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*metric.Aggregator) {
	switch descriptor.InstrumentKind() {
	case metric.ValueObserverInstrumentKind:
		lastValueAggs(aggPtrs)
	case metric.ValueRecorderInstrumentKind:
		aggs := minmaxsumcount.New(len(aggPtrs), descriptor)
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		sumAggs(aggPtrs)
	}
}

func (selectorExact) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*metric.Aggregator) {
	switch descriptor.InstrumentKind() {
	case metric.ValueObserverInstrumentKind:
		lastValueAggs(aggPtrs)
	case metric.ValueRecorderInstrumentKind:
		aggs := exact.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		sumAggs(aggPtrs)
	}
}

func (s selectorHistogram) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*metric.Aggregator) {
	switch descriptor.InstrumentKind() {
	case metric.ValueObserverInstrumentKind:
		lastValueAggs(aggPtrs)
	case metric.ValueRecorderInstrumentKind:
		aggs := histogram.New(len(aggPtrs), descriptor, s.options...)
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		sumAggs(aggPtrs)
	}
}
