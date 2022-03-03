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

package metric // import "go.opentelemetry.io/otel/sdk/export/metric"

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type Accumulation = export.Accumulation

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/aggregator"
type Aggregator = aggregator.Aggregator

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type AggregatorSelector = export.AggregatorSelector

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type Checkpointer = export.Checkpointer

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type CheckpointerFactory = export.CheckpointerFactory

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type Exporter = export.Exporter

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type InstrumentationLibraryReader = export.Exporter

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type Metadata = export.Metadata

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type Processor = export.Processor

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type Reader = export.Reader

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
type Record = export.Record

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
func NewAccumulation(descriptor *sdkapi.Descriptor, labels *attribute.Set, aggregator Aggregator) Accumulation {
	return export.NewAccumulation(descriptor, labels, aggregator)
}

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export"
func NewRecord(descriptor *sdkapi.Descriptor, labels *attribute.Set, aggregation aggregation.Aggregation, start, end time.Time) Record {
	return export.NewRecord(descriptor, labels, aggregation, start, end)
}
