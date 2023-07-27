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

	"go.opentelemetry.io/otel"
	internal "go.opentelemetry.io/otel/bridge/opencensus/internal/ocmetric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const scopeName = "go.opentelemetry.io/otel/bridge/opencensus"

// exporter wraps an OpenTelemetry exporter and adds OpenCensus metrics to it.
type exporter struct {
	manager *metricproducer.Manager
	base    metric.Exporter
}

// NewMetricExporter returns an OpenTelemetry exporter that adds metrics from OpenCensus
// before exporting to the base OpenTelemetry exporter.
func NewMetricExporter(base metric.Exporter) metric.Exporter {
	return &exporter{
		base:    base,
		manager: metricproducer.GlobalManager(),
	}
}

func (e *exporter) Export(ctx context.Context, sdkMetrics *metricdata.ResourceMetrics) error {
	producers := e.manager.GetAll()
	data := []*ocmetricdata.Metric{}
	for _, ocProducer := range producers {
		data = append(data, ocProducer.Read()...)
	}
	otelmetrics, err := internal.ConvertMetrics(data)
	if err != nil {
		otel.Handle(err)
	}
	if len(otelmetrics) > 0 {
		// add metrics from OpenCensus to our exported batch of metrics under
		// its own scope.
		sdkMetrics.ScopeMetrics = append(
			sdkMetrics.ScopeMetrics,
			metricdata.ScopeMetrics{
				Scope: instrumentation.Scope{
					Name: scopeName,
				},
				Metrics: otelmetrics,
			})
	}
	return e.base.Export(ctx, sdkMetrics)
}

func (e *exporter) Temporality(kind metric.InstrumentKind) metricdata.Temporality {
	return e.base.Temporality(kind)

}

func (e *exporter) Aggregation(kind metric.InstrumentKind) aggregation.Aggregation {
	return e.base.Aggregation(kind)

}

func (e *exporter) ForceFlush(ctx context.Context) error {
	return e.base.ForceFlush(ctx)

}

func (e *exporter) Shutdown(ctx context.Context) error {
	return e.base.Shutdown(ctx)
}
