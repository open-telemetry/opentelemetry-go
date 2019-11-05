// Copyright 2019, OpenTelemetry Authors
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

package simple // import "go.opentelemetry.io/otel/sdk/metric/selector/simpler"

import (
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/maxsumcount"
)

type selector struct{}

// New returns a simple aggregation selector that uses counter, gauge,
// and maxsumcount behavior for the three kinds of metric.
func New() export.AggregationSelector {
	return selector{}
}

func (s selector) AggregatorFor(record export.Record) export.Aggregator {
	switch record.Descriptor().MetricKind() {
	case export.GaugeKind:
		return gauge.New()
	case export.MeasureKind:
		return maxsumcount.New()
	default:
		return counter.New()
	}
}
