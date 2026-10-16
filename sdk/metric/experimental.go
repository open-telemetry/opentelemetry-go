// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
)

type finishOption interface {
	FinishEnabled() bool
}

// validateExperimental returns an error if options contain more than one
// distinct experimental option type, or nil otherwise. Combinations of
// experimental features are unsupported until a user requests one, so any
// two distinct types are rejected. Repeating the same type is allowed, and an
// option counts whether or not it enables its feature.
func validateExperimental(options []experimentalOption) error {
	var types []reflect.Type
	for _, o := range options {
		if t := reflect.TypeOf(o); !slices.Contains(types, t) {
			types = append(types, t)
		}
	}
	if len(types) <= 1 {
		return nil
	}

	names := make([]string, len(types))
	for i, t := range types {
		names[i] = t.String()
	}
	return fmt.Errorf("experimental options cannot be used together: %s", strings.Join(names, ", "))
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
