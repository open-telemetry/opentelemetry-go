// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
)

type finishOption interface {
	FinishEnabled() bool
}

type experimentalMeterConfig struct {
	finishRegistry *finish.Registry
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
	}
	if conf.finishRegistry == nil {
		return defaultMeterFactory, nil
	}
	return conf.newMeter, conf.finishRegistry.Shutdown
}

func (c experimentalMeterConfig) newMeter(s instrumentation.Scope, p pipelines) metric.Meter {
	m := newMeter(s, p)
	var decorated metric.Meter = m
	// Apply enabled extensions to one SDK meter so they share stream resolution.
	if c.finishRegistry != nil {
		decorated = finish.NewMeter(decorated, &factory{meter: m, registry: c.finishRegistry})
	}
	return decorated
}
