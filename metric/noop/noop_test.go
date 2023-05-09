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
		reflect.TypeOf((*metric.MeterProvider)(nil)).Elem(),
	))
	t.Run("Meter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Meter{}),
		reflect.TypeOf((*metric.Meter)(nil)).Elem(),
	))
	t.Run("Observer", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Observer{}),
		reflect.TypeOf((*metric.Observer)(nil)).Elem(),
	))
	t.Run("Registration", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Registration{}),
		reflect.TypeOf((*metric.Registration)(nil)).Elem(),
	))
	t.Run("Int64Counter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64Counter{}),
		reflect.TypeOf((*metric.Int64Counter)(nil)).Elem(),
	))
	t.Run("Float64Counter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64Counter{}),
		reflect.TypeOf((*metric.Float64Counter)(nil)).Elem(),
	))
	t.Run("Int64UpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64UpDownCounter{}),
		reflect.TypeOf((*metric.Int64UpDownCounter)(nil)).Elem(),
	))
	t.Run("Float64UpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64UpDownCounter{}),
		reflect.TypeOf((*metric.Float64UpDownCounter)(nil)).Elem(),
	))
	t.Run("Int64Histogram", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64Histogram{}),
		reflect.TypeOf((*metric.Int64Histogram)(nil)).Elem(),
	))
	t.Run("Float64Histogram", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64Histogram{}),
		reflect.TypeOf((*metric.Float64Histogram)(nil)).Elem(),
	))
	t.Run("Int64ObservableCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64ObservableCounter{}),
		reflect.TypeOf((*metric.Int64ObservableCounter)(nil)).Elem(),
	))
	t.Run("Float64ObservableCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64ObservableCounter{}),
		reflect.TypeOf((*metric.Float64ObservableCounter)(nil)).Elem(),
	))
	t.Run("Int64ObservableGauge", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64ObservableGauge{}),
		reflect.TypeOf((*metric.Int64ObservableGauge)(nil)).Elem(),
	))
	t.Run("Float64ObservableGauge", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64ObservableGauge{}),
		reflect.TypeOf((*metric.Float64ObservableGauge)(nil)).Elem(),
	))
	t.Run("Int64ObservableUpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64ObservableUpDownCounter{}),
		reflect.TypeOf((*metric.Int64ObservableUpDownCounter)(nil)).Elem(),
	))
	t.Run("Float64ObservableUpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64ObservableUpDownCounter{}),
		reflect.TypeOf((*metric.Float64ObservableUpDownCounter)(nil)).Elem(),
	))
	t.Run("Int64Observer", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Int64Observer{}),
		reflect.TypeOf((*metric.Int64Observer)(nil)).Elem(),
	))
	t.Run("Float64Observer", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Float64Observer{}),
		reflect.TypeOf((*metric.Float64Observer)(nil)).Elem(),
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
	assert.Equal(t, mp, MeterProvider{})
	meter := mp.Meter("")
	assert.Equal(t, meter, Meter{})
}
