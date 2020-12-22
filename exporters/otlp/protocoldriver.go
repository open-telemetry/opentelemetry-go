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

package otlp // import "go.opentelemetry.io/otel/exporters/otlp"

import (
	"context"

	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
)

// ProtocolDriver is an interface used by OTLP exporter. It's
// responsible for connecting to and disconnecting from the collector,
// and for transforming traces and metrics into wire format and
// transmitting them to the collector.
type ProtocolDriver interface {
	// Start should establish connection(s) to endpoint(s). It is
	// called just once by the exporter, so the implementation
	// does not need to worry about idempotence and locking.
	Start(ctx context.Context) error
	// Stop should close the connections. The function is called
	// only once by the exporter, so the implementation does not
	// need to worry about idempotence, but it may be called
	// concurrently with ExportMetrics or ExportTraces, so proper
	// locking is required. The function serves as a
	// synchronization point - after the function returns, the
	// process of closing connections is assumed to be finished.
	Stop(ctx context.Context) error
	// ExportMetrics should transform the passed metrics to the
	// wire format and send it to the collector. May be called
	// concurrently with ExportTraces, so the manager needs to
	// take this into account by doing proper locking.
	ExportMetrics(ctx context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error
	// ExportTraces should transform the passed traces to the wire
	// format and send it to the collector. May be called
	// concurrently with ExportMetrics, so the manager needs to
	// take this into account by doing proper locking.
	ExportTraces(ctx context.Context, ss []*tracesdk.SpanSnapshot) error
}
