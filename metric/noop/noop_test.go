// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package noop // import "go.opentelemetry.io/otel/metric/noop"

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/metric"
)

func TestImplementationNoPanics(t *testing.T) {
	// Check that if type has an embedded interface and that interface has
	// methods added to it than the No-Op implementation implements them.
	t.Run("MeterProvider", assertAllExportedMethodNoPanic(
		reflect.ValueOf(MeterProvider{}),
		reflect.TypeFor[metric.MeterProvider](),
	))
	t.Run("Meter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Meter{}),
		reflect.TypeFor[metric.Meter](),
	))
	t.Run("Observer", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Observer{}),
		reflect.TypeFor[metric.Observer](),
	))
	t.Run("Registration", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Registration{}),
		reflect.TypeFor[metric.Registration](),
	))
	t.Run("Int64Counter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64Counter{}),
		reflect.TypeFor[metric.Int64Counter](),
	))
	t.Run("Float64Counter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64Counter{}),
		reflect.TypeFor[metric.Float64Counter](),
	))
	t.Run("Int64UpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64UpDownCounter{}),
		reflect.TypeFor[metric.Int64UpDownCounter](),
	))
	t.Run("Float64UpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64UpDownCounter{}),
		reflect.TypeFor[metric.Float64UpDownCounter](),
	))
	t.Run("Int64Histogram", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64Histogram{}),
		reflect.TypeFor[metric.Int64Histogram](),
	))
	t.Run("Float64Histogram", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64Histogram{}),
		reflect.TypeFor[metric.Float64Histogram](),
	))
	t.Run("Int64Gauge", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64Gauge{}),
		reflect.TypeFor[metric.Int64Gauge](),
	))
	t.Run("Float64Gauge", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64Gauge{}),
		reflect.TypeFor[metric.Float64Gauge](),
	))
	t.Run("Int64ObservableCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64ObservableCounter{}),
		reflect.TypeFor[metric.Int64ObservableCounter](),
	))
	t.Run("Float64ObservableCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64ObservableCounter{}),
		reflect.TypeFor[metric.Float64ObservableCounter](),
	))
	t.Run("Int64ObservableGauge", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64ObservableGauge{}),
		reflect.TypeFor[metric.Int64ObservableGauge](),
	))
	t.Run("Float64ObservableGauge", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64ObservableGauge{}),
		reflect.TypeFor[metric.Float64ObservableGauge](),
	))
	t.Run("Int64ObservableUpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64ObservableUpDownCounter{}),
		reflect.TypeFor[metric.Int64ObservableUpDownCounter](),
	))
	t.Run("Float64ObservableUpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64ObservableUpDownCounter{}),
		reflect.TypeFor[metric.Float64ObservableUpDownCounter](),
	))
	t.Run("Int64Observer", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64Observer{}),
		reflect.TypeFor[metric.Int64Observer](),
	))
	t.Run("Float64Observer", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64Observer{}),
		reflect.TypeFor[metric.Float64Observer](),
	))
}

func assertAllExportedMethodNoPanic(rVal reflect.Value, rType reflect.Type) func(*testing.T) {
	return func(t *testing.T) {
		for n := 0; n < rType.NumMethod(); n++ {
			mType := rType.Method(n)
			if !mType.IsExported() {
				t.Logf("ignoring unexported %s", mType.Name)
				continue
			}
			m := rVal.MethodByName(mType.Name)
			if !m.IsValid() {
				t.Errorf("unknown method for %s: %s", rVal.Type().Name(), mType.Name)
			}

			numIn := mType.Type.NumIn()
			if mType.Type.IsVariadic() {
				numIn--
			}
			args := make([]reflect.Value, numIn)
			for i := range args {
				aType := mType.Type.In(i)
				args[i] = reflect.New(aType).Elem()
			}

			assert.NotPanicsf(t, func() {
				_ = m.Call(args)
			}, "%s.%s", rVal.Type().Name(), mType.Name)
		}
	}
}

func TestNewMeterProvider(t *testing.T) {
	mp := NewMeterProvider()
	assert.Equal(t, MeterProvider{}, mp)
	meter := mp.Meter("")
	assert.Equal(t, Meter{}, meter)
}
