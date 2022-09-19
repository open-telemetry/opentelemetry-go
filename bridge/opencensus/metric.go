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

	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"

	"go.opentelemetry.io/otel/bridge/opencensus/opencensusmetric/internal"
	"go.opentelemetry.io/otel/sdk/metric"
)

// exporter implements the OpenCensus metric Exporter interface using an
// OpenTelemetry base exporter.
type exporter struct {
	base metric.Exporter
}

// NewMetricExporter returns an OpenCensus exporter that exports to an
// OpenTelemetry exporter.
func NewMetricExporter(base metric.Exporter) metricexport.Exporter {
	return &exporter{base: base}
}

// ExportMetrics implements the OpenCensus metric Exporter interface by sending
// to an OpenTelemetry exporter.
func (e *exporter) ExportMetrics(ctx context.Context, ocmetrics []*metricdata.Metric) error {
	otelmetrics, err := internal.ConvertMetrics(ocmetrics)
	if err != nil {
		return err
	}
	return e.base.Export(ctx, otelmetrics)
}
