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

// exporter implements the OpenCensus metric Exporter interface using an
// OpenTelemetry base exporter.
type exporter struct {
	base metric.Exporter
	res  *resource.Resource
}

// NewMetricExporter returns an OpenCensus exporter that exports to an
// OpenTelemetry (push) exporter.
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

type reader struct {
	// This reader doesn't actually matter. It is just there to satisfy the Reader interface.
	metric.Reader
	res     *resource.Resource
	manager *metricproducer.Manager
}

// NewMetricExporter returns an OpenTelemetry metric.Reader that can be used as
// a source of metrics for OpenTelemetry (pull) exporters.
func NewMetricReader(res *resource.Resource) metric.Reader {
	return &reader{
		Reader:  metric.NewManualReader(),
		res:     res,
		manager: metricproducer.GlobalManager(),
	}
}

// Collect overrides the manual reader's collect method to retrieve metrics from OpenCensus.
func (r *reader) Collect(context.Context) (metricdata.ResourceMetrics, error) {
	producers := r.manager.GetAll()
	data := []*ocmetricdata.Metric{}
	for _, ocProducer := range producers {
		data = append(data, ocProducer.Read()...)
	}
	otelmetrics, err := internal.ConvertMetrics(data)
	if err != nil {
		return metricdata.ResourceMetrics{}, err
	}
	return metricdata.ResourceMetrics{
		Resource: r.res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name: "go.opentelemetry.io/otel/bridge/opencensus",
				},
				Metrics: otelmetrics,
			},
		}}, nil
}
