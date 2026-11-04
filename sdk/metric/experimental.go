// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"errors"
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

// rejectDistinctExperimental returns an error if options contain more than one
// distinct option type of the same kind, or nil otherwise. Options that
// contribute views (experimentalViewsOption) are one kind, and all other
// experimental options are the other. Combinations within a kind are
// unsupported until a user requests one, while options of different kinds are
// independent: views do not affect how meters are built. Repeating the same
// type is allowed and counted once, and an option counts whether or not it
// enables its feature.
func rejectDistinctExperimental(options []experimentalOption) error {
	var meterTypes, viewTypes []reflect.Type
	for _, o := range options {
		types := &meterTypes
		if _, ok := o.(experimentalViewsOption); ok {
			types = &viewTypes
		}
		if t := reflect.TypeOf(o); !slices.Contains(*types, t) {
			*types = append(*types, t)
		}
	}
	return errors.Join(conflictError(meterTypes), conflictError(viewTypes))
}

func conflictError(types []reflect.Type) error {
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
