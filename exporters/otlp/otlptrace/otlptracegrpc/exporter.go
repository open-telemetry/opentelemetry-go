// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracegrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/otlpconfig"
)

// New constructs a new Exporter and starts it.
func New(ctx context.Context, opts ...Option) (*otlptrace.Exporter, error) {
	client := newClient(opts...)
	if  otlpconfig.GetEnvConfigErr()!=nil{
		return nil,otlpconfig.GetEnvConfigErr()
	}
	return otlptrace.New(ctx, client)
}

// NewUnstarted constructs a new Exporter and does not start it.
func NewUnstarted(opts ...Option) *otlptrace.Exporter {
	return otlptrace.NewUnstarted(NewClient(opts...))
}
