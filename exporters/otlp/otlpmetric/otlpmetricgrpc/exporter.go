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

package otlpmetricgrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"

import (
	"context"

	ominternal "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Exporter is a OpenTelemetry metric Exporter using gRPC.
type Exporter struct {
	wrapped *ominternal.Exporter
}

// Temporality returns the Temporality to use for an instrument kind.
func (e *Exporter) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return e.wrapped.Temporality(k)
}

// Aggregation returns the Aggregation to use for an instrument kind.
func (e *Exporter) Aggregation(k metric.InstrumentKind) aggregation.Aggregation {
	return e.wrapped.Aggregation(k)
}

// Export transforms and transmits metric data to an OTLP receiver.
//
// This method returns an error if called after Shutdown.
// This method returns an error if the method is canceled by the passed context.
func (e *Exporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	return e.wrapped.Export(ctx, rm)
}

// ForceFlush flushes any metric data held by an exporter.
//
// This method returns an error if called after Shutdown.
// This method returns an error if the method is canceled by the passed context.
//
// This method is safe to call concurrently.
func (e *Exporter) ForceFlush(ctx context.Context) error {
	return e.wrapped.ForceFlush(ctx)
}

// Shutdown flushes all metric data held by an exporter and releases any held
// computational resources.
//
// This method returns an error if called after Shutdown.
// This method returns an error if the method is canceled by the passed context.
//
// This method is safe to call concurrently.
func (e *Exporter) Shutdown(ctx context.Context) error {
	return e.wrapped.Shutdown(ctx)
}

// New returns an OpenTelemetry metric Exporter. The Exporter can be used with
// a PeriodicReader to export OpenTelemetry metric data to an OTLP receiving
// endpoint using gRPC.
//
// If an already established gRPC ClientConn is not passed in options using
// WithGRPCConn, a connection to the OTLP endpoint will be established based
// on options. If a connection cannot be establishes in the lifetime of ctx,
// an error will be returned.
func New(ctx context.Context, options ...Option) (*Exporter, error) {
	c, err := newClient(ctx, options...)
	if err != nil {
		return nil, err
	}
	exp := ominternal.New(c)
	return &Exporter{exp}, nil
}
