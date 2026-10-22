// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

// configuratorInt64Inst gates an int64 synchronous instrument on its
// configuratorMeter's config. Only configuratorMeter creates it, so instruments
// from a provider without a configurator record through int64Inst directly.
type configuratorInt64Inst struct {
	*int64Inst

	config *versionedMeterConfig
}

var (
	_ metric.Int64Counter       = (*configuratorInt64Inst)(nil)
	_ metric.Int64UpDownCounter = (*configuratorInt64Inst)(nil)
	_ metric.Int64Histogram     = (*configuratorInt64Inst)(nil)
	_ metric.Int64Gauge         = (*configuratorInt64Inst)(nil)
)

func (i *configuratorInt64Inst) Add(ctx context.Context, val int64, opts ...metric.AddOption) {
	if !i.config.Load() {
		return
	}
	i.int64Inst.Add(ctx, val, opts...)
}

func (i *configuratorInt64Inst) Record(ctx context.Context, val int64, opts ...metric.RecordOption) {
	if !i.config.Load() {
		return
	}
	i.int64Inst.Record(ctx, val, opts...)
}

func (i *configuratorInt64Inst) Enabled(ctx context.Context) bool {
	return i.config.Load() && i.int64Inst.Enabled(ctx)
}

// configuratorFloat64Inst gates a float64 synchronous instrument on its
// configuratorMeter's config. Only configuratorMeter creates it, so instruments
// from a provider without a configurator record through float64Inst directly.
type configuratorFloat64Inst struct {
	*float64Inst

	config *versionedMeterConfig
}

var (
	_ metric.Float64Counter       = (*configuratorFloat64Inst)(nil)
	_ metric.Float64UpDownCounter = (*configuratorFloat64Inst)(nil)
	_ metric.Float64Histogram     = (*configuratorFloat64Inst)(nil)
	_ metric.Float64Gauge         = (*configuratorFloat64Inst)(nil)
)

func (i *configuratorFloat64Inst) Add(ctx context.Context, val float64, opts ...metric.AddOption) {
	if !i.config.Load() {
		return
	}
	i.float64Inst.Add(ctx, val, opts...)
}

func (i *configuratorFloat64Inst) Record(ctx context.Context, val float64, opts ...metric.RecordOption) {
	if !i.config.Load() {
		return
	}
	i.float64Inst.Record(ctx, val, opts...)
}

func (i *configuratorFloat64Inst) Enabled(ctx context.Context) bool {
	return i.config.Load() && i.float64Inst.Enabled(ctx)
}
