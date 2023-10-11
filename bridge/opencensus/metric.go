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
	"context"

	ocmetricdata "go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricproducer"

	internal "go.opentelemetry.io/otel/bridge/opencensus/internal/ocmetric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// MetricProducer implements the [go.opentelemetry.io/otel/sdk/metric.Producer] to provide metrics
// from OpenCensus to the OpenTelemetry SDK.
type MetricProducer struct {
	manager *metricproducer.Manager
}

// NewMetricProducer returns a metric.Producer that fetches metrics from
// OpenCensus.
func NewMetricProducer(opts ...MetricOption) *MetricProducer {
	return &MetricProducer{
		manager: metricproducer.GlobalManager(),
	}
}

var _ metric.Producer = (*MetricProducer)(nil)

// Produce fetches metrics from the OpenCensus manager,
// translates them to OpenTelemetry's data model, and returns them.
func (p *MetricProducer) Produce(context.Context) ([]metricdata.ScopeMetrics, error) {
	producers := p.manager.GetAll()
	data := []*ocmetricdata.Metric{}
	for _, ocProducer := range producers {
		data = append(data, ocProducer.Read()...)
	}
	otelmetrics, err := internal.ConvertMetrics(data)
	if len(otelmetrics) == 0 {
		return nil, err
	}
	return []metricdata.ScopeMetrics{{
		Scope: instrumentation.Scope{
			Name:    scopeName,
			Version: Version(),
		},
		Metrics: otelmetrics,
	}}, err
}
