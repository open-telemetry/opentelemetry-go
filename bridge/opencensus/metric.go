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
	"go.opencensus.io/metric/metricexport"
	"go.opencensus.io/metric/metricproducer"

	"go.opentelemetry.io/otel"
	internal "go.opentelemetry.io/otel/bridge/opencensus/internal/ocmetric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

const scopeName = "go.opentelemetry.io/otel/bridge/opencensus"

type producer struct {
	manager *metricproducer.Manager
}

// NewMetricProducer returns a metric.Producer that fetches metrics from
// OpenCensus.
func NewMetricProducer() metric.Producer {
	return &producer{
		manager: metricproducer.GlobalManager(),
	}
}

func (p *producer) Produce(context.Context) ([]metricdata.ScopeMetrics, error) {
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
			Name: scopeName,
		},
		Metrics: otelmetrics,
	}}, err
}

// exporter implements the OpenCensus metric Exporter interface using an
// OpenTelemetry base exporter.
type exporter struct {
	base metric.Exporter
	res  *resource.Resource
}

// NewMetricExporter returns an OpenCensus exporter that exports to an
// OpenTelemetry (push) exporter.
// Deprecated: Use NewMetricProducer instead.
func NewMetricExporter(base metric.Exporter, res *resource.Resource) metricexport.Exporter {
	return &exporter{base: base, res: res}
}

// ExportMetrics implements the OpenCensus metric Exporter interface by sending
// to an OpenTelemetry exporter.
func (e *exporter) ExportMetrics(ctx context.Context, ocmetrics []*ocmetricdata.Metric) error {
	otelmetrics, err := internal.ConvertMetrics(ocmetrics)
	if err != nil {
		otel.Handle(err)
	}
	if len(otelmetrics) == 0 {
		return nil
	}
	return e.base.Export(ctx, &metricdata.ResourceMetrics{
		Resource: e.res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name: scopeName,
				},
				Metrics: otelmetrics,
			},
		}})
}
