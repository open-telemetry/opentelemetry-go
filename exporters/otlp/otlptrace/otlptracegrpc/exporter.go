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

package otlptracegrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// NewExporter constructs a new Exporter and starts it.
func NewExporter(ctx context.Context, opts ...Option) (*otlptrace.Exporter, error) {
	return otlptrace.NewExporter(ctx, NewClient(opts...))
}

// NewUnstartedExporter constructs a new Exporter and does not start it.
func NewUnstartedExporter(opts ...Option) *otlptrace.Exporter {
	return otlptrace.NewUnstartedExporter(NewClient(opts...))
}

// NewExportPipeline sets up a complete export pipeline
// with the recommended TracerProvider setup.
func NewExportPipeline(ctx context.Context, opts ...Option) (*otlptrace.Exporter, *tracesdk.TracerProvider, error) {
	return otlptrace.NewExportPipeline(ctx, NewClient(opts...))
}

// InstallNewPipeline instantiates a NewExportPipeline with the
// recommended configuration and registers it globally.
func InstallNewPipeline(ctx context.Context, opts ...Option) (*otlptrace.Exporter, *tracesdk.TracerProvider, error) {
	return otlptrace.InstallNewPipeline(ctx, NewClient(opts...))
}
