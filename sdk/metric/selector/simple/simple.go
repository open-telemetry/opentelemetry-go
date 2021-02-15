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
	"sync"

	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
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

	// SelectorDelegator implements export.AggregatorSelector and provides a way to
	// register a specific AggregatorSelector for a given metric name.
	SelectorDelegator struct {
		// defaultSelector is the default AggregatorSelector for an unregistered metric.
		defaultSelector export.AggregatorSelector

		// namedSelector stores the mapping from metric name to AggregatorSelector.
		namedSelector sync.Map
	}
)

var (
	_ export.AggregatorSelector = selectorInexpensive{}
	_ export.AggregatorSelector = selectorExact{}
	_ export.AggregatorSelector = selectorHistogram{}
	_ export.AggregatorSelector = &SelectorDelegator{}
)

// NewWithDelegate returns a SelectorDelegator with given default selector.
func NewWithDelegate(defaultSelector export.AggregatorSelector) *SelectorDelegator {
	return &SelectorDelegator{
		defaultSelector: defaultSelector,
	}
}

// NewWithInexpensiveDistribution returns a simple aggregator selector
// that uses minmaxsumcount aggregators for `ValueRecorder`
// instruments.  This selector is faster and uses less memory than the
// others in this package because minmaxsumcount aggregators maintain
// the least information about the distribution among these choices.
func NewWithInexpensiveDistribution() export.AggregatorSelector {
	return selectorInexpensive{}
}

// NewWithExactDistribution returns a simple aggregator selector that
// uses exact aggregators for `ValueRecorder` instruments.  This
// selector uses more memory than the others in this package because
// exact aggregators maintain the most information about the
// distribution among these choices.
func NewWithExactDistribution() export.AggregatorSelector {
	return selectorExact{}
}

// NewWithHistogramDistribution returns a simple aggregator selector
// that uses histogram aggregators for `ValueRecorder` instruments.
// This selector is a good default choice for most metric exporters.
func NewWithHistogramDistribution(options ...histogram.Option) export.AggregatorSelector {
	return selectorHistogram{options: options}
}

func sumAggs(aggPtrs []*export.Aggregator) {
	aggs := sum.New(len(aggPtrs))
	for i := range aggPtrs {
		*aggPtrs[i] = &aggs[i]
	}
}

func lastValueAggs(aggPtrs []*export.Aggregator) {
	aggs := lastvalue.New(len(aggPtrs))
	for i := range aggPtrs {
		*aggPtrs[i] = &aggs[i]
	}
}

func (selectorInexpensive) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*export.Aggregator) {
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

func (selectorExact) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*export.Aggregator) {
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

func (s selectorHistogram) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*export.Aggregator) {
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

// Register provides a way to assosicate a metric name with a specific AggregatorSelector.
// If the name is registered before, the call has no effect. Returns true if the name was
// not registered before.
func (s *SelectorDelegator) Register(instrumentName string, aselector export.AggregatorSelector) bool {
	_, loaded := s.namedSelector.LoadOrStore(instrumentName, aselector)
	return !loaded
}

// AggregatorFor initializes required AggregatorSelector based on the metric name. It always
// returns the same kind of AggregatorSelector for a given metric name.
func (s *SelectorDelegator) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*export.Aggregator) {
	actual, _ := s.namedSelector.LoadOrStore(descriptor.Name(), s.defaultSelector)
	actual.(export.AggregatorSelector).AggregatorFor(descriptor, aggPtrs...)
}
