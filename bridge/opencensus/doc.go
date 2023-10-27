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

// Package opencensus provides a migration bridge from OpenCensus to
// OpenTelemetry for metrics and traces. The bridge incorporates metrics and
// traces from OpenCensus into the OpenTelemetry SDK, combining them with
// metrics and traces from OpenTelemetry instrumentation.
//
// # Migration Guide
//
// For most applications, it would be difficult to migrate an application
// from OpenCensus to OpenTelemetry all-at-once. Libraries used by the
// application may still be using OpenCensus, and the application itself may
// have many lines of instrumentation.
//
// Bridges help in this situation by allowing your application to have "mixed"
// instrumentation, while incorporating all instrumentation into a single
// export path. To migrate with bridges, a user would:
//
//  1. Configure the OpenTelemetry SDK for metrics and traces, with the OpenTelemetry exporters matching to your current OpenCensus exporters.
//  2. Install this OpenCensus bridge, which sends OpenCensus telemetry to your new OpenTelemetry exporters.
//  3. Over time, migrate your instrumentation from OpenCensus to OpenTelemetry.
//  4. Once all instrumentation is migrated, remove the OpenCensus bridge.
//
// With this approach, you can migrate your telemetry, including in dependent
// libraries over time without disruption.
//
// # Warnings
//
// Installing a metric or tracing bridge will cause OpenCensus telemetry to be
// exported by OpenTelemetry exporters. Since OpenCensus telemetry uses globals,
// installing a bridge will result in telemetry collection from _all_ libraries
// that use OpenCensus, including some you may not expect, such as the
// telemetry exporter itself.
//
// # Limitations
//
// There are known limitations to the trace bridge:
//
//   - The AddLink method for OpenCensus Spans is ignored, and an error is sent
//     to the OpenTelemetry ErrorHandler.
//   - The NewContext method of the OpenCensus Tracer cannot embed an OpenCensus
//     Span in a context unless that Span was created by that Tracer.
//   - Conversion of custom OpenCensus Samplers to OpenTelemetry is not
//     implemented, and An error will be sent to the OpenTelemetry ErrorHandler.
//
// There are known limitations to the metric bridge:
//   - GaugeDistribution-typed metrics are dropped
//   - Histogram's SumOfSquaredDeviation field is dropped
package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"
