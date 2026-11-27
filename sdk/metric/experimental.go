// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
)

type finishOption interface {
	FinishEnabled() bool
}

type bindOption interface {
	BindingEnabled() bool
}

type experimentalMeterConfig struct {
	finishRegistry *finish.Registry
	bindEnabled    bool
}

func newExperimentalMeterFactory(
	options []experimentalOption,
) (meterFactoryFn meterFactory, shutdown func(context.Context) error) {
	var conf experimentalMeterConfig
	for _, option := range options {
		if opt, ok := option.(finishOption); ok && opt.FinishEnabled() {
			if conf.finishRegistry == nil {
				conf.finishRegistry = &finish.Registry{}
			}
		}
		if opt, ok := option.(bindOption); ok && opt.BindingEnabled() {
			conf.bindEnabled = true
		}
	}
	if conf.finishRegistry == nil && !conf.bindEnabled {
		return defaultMeterFactory, nil
	}
	if conf.finishRegistry != nil && conf.bindEnabled {
		global.Warn(
			"combining experimental Finish and Binding options is not yet supported; " +
				"Binding overrides Finish on supported instruments",
		)
	}
	var shutdownFn func(context.Context) error
	if conf.finishRegistry != nil {
		shutdownFn = conf.finishRegistry.Shutdown
	}
	return conf.newMeter, shutdownFn
}

func (c experimentalMeterConfig) newMeter(s instrumentation.Scope, p pipelines) metric.Meter {
	m := newMeter(s, p)
	var decorated metric.Meter = m
	// Apply enabled extensions to one SDK meter so they share stream resolution.
	// Combining Finish and Binding is not yet supported; bindMeter overrides
	// finishMeter for supported Counter and UpDownCounter instruments.
	if c.finishRegistry != nil {
		decorated = finish.NewMeter(decorated, &factory{meter: m, registry: c.finishRegistry})
	}
	if c.bindEnabled {
		decorated = newBindMeter(decorated, m)
	}
	return decorated
}
