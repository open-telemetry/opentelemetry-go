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
	"go.opentelemetry.io/otel/metric/instrument"
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
	t.Run("Counter[int64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Counter[int64]{}),
		reflect.TypeOf((*instrument.Counter[int64])(nil)).Elem(),
	))
	t.Run("Counter[float64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Counter[float64]{}),
		reflect.TypeOf((*instrument.Counter[float64])(nil)).Elem(),
	))
	t.Run("UpDownCounter[int64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(UpDownCounter[int64]{}),
		reflect.TypeOf((*instrument.UpDownCounter[int64])(nil)).Elem(),
	))
	t.Run("UpDownCounter[float64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(UpDownCounter[float64]{}),
		reflect.TypeOf((*instrument.UpDownCounter[float64])(nil)).Elem(),
	))
	t.Run("Histogram[int64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Histogram[int64]{}),
		reflect.TypeOf((*instrument.Histogram[int64])(nil)).Elem(),
	))
	t.Run("Histogram[float64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Histogram[float64]{}),
		reflect.TypeOf((*instrument.Histogram[float64])(nil)).Elem(),
	))
	t.Run("ObservableCounter[int64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(ObservableCounter[int64]{}),
		reflect.TypeOf((*instrument.ObservableCounter[int64])(nil)).Elem(),
	))
	t.Run("ObservableCounter[float64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(ObservableCounter[float64]{}),
		reflect.TypeOf((*instrument.ObservableCounter[float64])(nil)).Elem(),
	))
	t.Run("ObservableGauge[int64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(ObservableGauge[int64]{}),
		reflect.TypeOf((*instrument.ObservableGauge[int64])(nil)).Elem(),
	))
	t.Run("ObservableGauge[float64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(ObservableGauge[float64]{}),
		reflect.TypeOf((*instrument.ObservableGauge[float64])(nil)).Elem(),
	))
	t.Run("ObservableUpDownCounter[int64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(ObservableUpDownCounter[int64]{}),
		reflect.TypeOf((*instrument.ObservableUpDownCounter[int64])(nil)).Elem(),
	))
	t.Run("ObservableUpDownCounter[float64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(ObservableUpDownCounter[float64]{}),
		reflect.TypeOf((*instrument.ObservableUpDownCounter[float64])(nil)).Elem(),
	))
	t.Run("Observer[int64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(ObserverT[int64]{}),
		reflect.TypeOf((*instrument.ObserverT[int64])(nil)).Elem(),
	))
	t.Run("Observer[float64]", assertAllExportedMethodNoPanic(
		reflect.ValueOf(ObserverT[float64]{}),
		reflect.TypeOf((*instrument.ObserverT[float64])(nil)).Elem(),
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
