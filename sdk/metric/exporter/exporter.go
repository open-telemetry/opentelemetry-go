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
package exporter

import (
	"context"

	"go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Exporter handles presentation of the checkpoint of aggregate
// metrics.  This is the final stage of a metrics export pipeline,
// where metric data are formatted for a specific system
type Exporter interface {
	// Export is called immediately after completing a collection
	// pass in the SDK.
	//
	// The Context comes from the controller that initiated
	// collection.
	//
	// The InstrumentationLibraryReader interface refers to the
	// Processor that just completed collection.
	Export(ctx context.Context, resource *resource.Resource, reader metric.InstrumentationLibraryReader) error

	// TemporalitySelector is an interface used by the Processor
	// in deciding whether to compute Delta or Cumulative
	// Aggregations when passing Records to this Exporter.
	aggregation.TemporalitySelector
}
