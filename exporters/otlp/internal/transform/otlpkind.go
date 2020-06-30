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

package transform

import (
	"fmt"

	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

// OTLPKind returns the OTLP Kind value given the SDK's descriptor.
// This assumes that export.ExportKind == export.PassThroughExporter.
func OTLPKind(mdesc *metric.Descriptor, akind aggregation.Kind) (metricpb.MetricDescriptor_Kind, error) {

	var pkind metricpb.MetricDescriptor_KindMask

	mkind := mdesc.MetricKind()

	if mkind.Adding() {
		pkind |= metricpb.MetricDescriptor_ADDING
	} else {
		pkind |= metricpb.MetricDescriptor_GROUPING
	}

	if mkind.Synchronous() {
		pkind |= metricpb.MetricDescriptor_SYNCHRONOUS
	}

	if mkind.Monotonic() {
		pkind |= metricpb.MetricDescriptor_MONOTONIC
	}

	switch akind {
	case aggregation.SumKind, aggregation.MinMaxSumCountKind, aggregation.HistogramKind, aggregation.ExactKind:

		if mkind.PrecomputedSum() {
			pkind |= metricpb.MetricDescriptor_CUMULATIVE
		} else {
			pkind |= metricpb.MetricDescriptor_DELTA
		}
		return metricpb.MetricDescriptor_Kind(pkind), nil

	case aggregation.LastValueKind:
		pkind |= metricpb.MetricDescriptor_INSTANTANEOUS
		return metricpb.MetricDescriptor_Kind(pkind), nil
	}

	fmt.Println("HERE", akind, mkind)
	// Note: this includes aggregation.SketchKind.
	return metricpb.MetricDescriptor_INVALID_KIND, aggregation.ErrNotSupported
}
