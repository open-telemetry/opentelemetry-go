// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
)

type finishOption interface {
	FinishEnabled() bool
}

// experimentalConflict is a pair of experimental features that cannot be used
// together. Each feature is identified by an interface its options implement,
// such as finishOption.
type experimentalConflict [2]reflect.Type

// validateExperimental returns an error for every conflict in conflicts whose
// features are both implemented by options, or nil if there is none.
func validateExperimental(options []experimentalOption, conflicts []experimentalConflict) error {
	if len(options) == 0 || len(conflicts) == 0 {
		return nil
	}

	// present maps each feature named by conflicts to the first option
	// implementing it.
	present := make(map[reflect.Type]experimentalOption)
	for _, c := range conflicts {
		for _, feature := range c {
			if _, ok := present[feature]; ok || feature == nil || feature.Kind() != reflect.Interface {
				continue
			}
			for _, o := range options {
				if reflect.TypeOf(o).Implements(feature) {
					present[feature] = o
					break
				}
			}
		}
	}

	var errs []error
	for _, c := range conflicts {
		optA, okA := present[c[0]]
		optB, okB := present[c[1]]
		if okA && okB {
			errs = append(errs, fmt.Errorf("experimental option %T cannot be used with %T", optA, optB))
		}
	}
	return errors.Join(errs...)
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
