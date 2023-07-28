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
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const scopeName = "go.opentelemetry.io/otel/bridge/opencensus"

// NewMetricExporter returns an OpenTelemetry exporter that adds metrics from OpenCensus
// before exporting to the base OpenTelemetry exporter.
func NewMetricExporter(base metric.Exporter) *Exporter {
	return &Exporter{
		Exporter: base,
	}
}

// Exporter wraps an OpenTelemetry exporter and adds OpenCensus metrics to it.
type Exporter struct {
	metric.Exporter
}

// Export exports a batch of metrics from the OpenTelemetry SDK, and adds
// metrics from OpenCensus prior to exporting.
func (e *Exporter) Export(ctx context.Context, sdkMetrics *metricdata.ResourceMetrics) error {
	appendOpenCensusMetrics(sdkMetrics)
	return e.Exporter.Export(ctx, sdkMetrics)
}

// NewMetricReader wraps an existing metric reader, but overrides the Collect
// function to insert metrics from OpenCensus.
func NewMetricReader(opts ...metric.ManualReaderOption) *Reader {
	return &Reader{
		Reader: metric.NewManualReader(opts...),
	}
}

// Reader wraps a metric.Reader, and adds metrics from OpenCensus when Collect
// is invoked.
type Reader struct {
	metric.Reader
}

// Override the collect function with one that inserts metrics from OpenCensus.
func (r *Reader) Collect(ctx context.Context, sdkMetrics *metricdata.ResourceMetrics) error {
	// Collect metrics from the OTel SDK
	err := r.Reader.Collect(ctx, sdkMetrics)
	if err != nil {
		return err
	}
	appendOpenCensusMetrics(sdkMetrics)
	return nil
}

// appendOpenCensusMetrics gets metrics from the OpenCensus manager, and
// appends them to ResourceMetrics from the OpenTelemetry SDK.
func appendOpenCensusMetrics(sdkMetrics *metricdata.ResourceMetrics) {
	// Collect metrics from OpenCensus
	producers := metricproducer.GlobalManager().GetAll()
	data := []*ocmetricdata.Metric{}
	for _, ocProducer := range producers {
		data = append(data, ocProducer.Read()...)
	}
	convertedMetrics, err := internal.ConvertMetrics(data)
	if err != nil {
		otel.Handle(err)
	}
	// Insert metrics from OpenCensus into the metrics from the OTel SDK
	if len(convertedMetrics) > 0 {
		sdkMetrics.ScopeMetrics = append(
			sdkMetrics.ScopeMetrics,
			metricdata.ScopeMetrics{
				Scope: instrumentation.Scope{
					Name: scopeName,
				},
				Metrics: convertedMetrics,
			})
	}
}
