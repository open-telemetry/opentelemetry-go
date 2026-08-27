// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// MeterProvider is the experimental metric SDK facade. Its stable delegate
// owns all readers, views, pipelines, and aggregation state.
type MeterProvider struct {
	embedded.MeterProvider
	delegate *sdkmetric.MeterProvider
}

var _ metric.MeterProvider = (*MeterProvider)(nil)

// NewMeterProvider returns a stable-backed provider whose Int64Counter values
// implement the experimental binding and finishing contracts.
func NewMeterProvider(options ...Option) *MeterProvider {
	options = append(options, boundCounterOption{Option: sdkmetric.WithView()})
	return &MeterProvider{delegate: sdkmetric.NewMeterProvider(options...)}
}

// Meter returns a Meter backed by the stable SDK runtime.
func (p *MeterProvider) Meter(name string, options ...metric.MeterOption) metric.Meter {
	return p.delegate.Meter(name, options...)
}

// ForceFlush flushes all configured Readers.
func (p *MeterProvider) ForceFlush(ctx context.Context) error {
	return p.delegate.ForceFlush(ctx)
}

// Shutdown shuts down the underlying stable SDK provider.
func (p *MeterProvider) Shutdown(ctx context.Context) error {
	return p.delegate.Shutdown(ctx)
}
