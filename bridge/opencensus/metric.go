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

//go:build go1.18
// +build go1.18

package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"

import (
	"context"

	ocmetricdata "go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"
	"go.opencensus.io/metric/metricproducer"

	internal "go.opentelemetry.io/otel/bridge/opencensus/internal/ocmetric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	// instrumentationName is the name of this instrumentation package.
	instrumentationName = "go.opentelemetry.io/otel/bridge/opencensus"
)

// producer is a producer which provides metrics collected using OpenCensus
// instrumentation.
type producer struct {
	scope   instrumentation.Scope
	manager *metricproducer.Manager
}

// NewProducer returns a producer which can be invoked to collect metrics.
func NewProducer() metric.Producer {
	return &producer{
		scope:   instrumentation.Scope{Name: instrumentationName, Version: SemVersion()},
		manager: metricproducer.GlobalManager(),
	}
}

// Produce gathers all metrics from the OpenCensus in-memory state.
func (p *producer) Produce(context.Context) ([]metricdata.Metrics, error) {
	producers := p.manager.GetAll()
	data := []*ocmetricdata.Metric{}
	for _, ocProducer := range producers {
		data = append(data, ocProducer.Read()...)
	}
	return internal.ConvertMetrics(data)
}

// InstrumentationScope returns the instrumentation scope for the OpenCensus
// metrics bridge.
func (p *producer) InstrumentationScope() instrumentation.Scope {
	return p.scope
}

// exporter implements the OpenCensus metric Exporter interface using an
// OpenTelemetry base exporter.
type exporter struct {
	base metric.Exporter
	res  *resource.Resource
}

// NewMetricExporter returns an OpenCensus exporter that exports to an
// OpenTelemetry exporter.
// Deprecated: Pass metric.WithProducer(opencensus.NewProducer()) to NewMeterProvider instead.
func NewMetricExporter(base metric.Exporter, res *resource.Resource) metricexport.Exporter {
	return &exporter{base: base}
}

// ExportMetrics implements the OpenCensus metric Exporter interface by sending
// to an OpenTelemetry exporter.
func (e *exporter) ExportMetrics(ctx context.Context, ocmetrics []*ocmetricdata.Metric) error {
	otelmetrics, err := internal.ConvertMetrics(ocmetrics)
	if err != nil {
		return err
	}
	return e.base.Export(ctx, metricdata.ResourceMetrics{
		Resource: e.res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name: "go.opentelemetry.io/otel/bridge/opencensus",
				},
				Metrics: otelmetrics,
			},
		}})
}
