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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type req[V any] struct {
	t *testing.T
}

func (r req[V]) NoErr(v V, err error) V {
	require.NoError(r.t, err)
	return v
}

func TestImplementationNoPanics(t *testing.T) {
	meterProvider := NewMeterProvider()
	t.Run("MeterProvider", assertAllExportedMethodNoPanic(
		reflect.ValueOf(meterProvider),
		reflect.TypeOf((*metric.MeterProvider)(nil)).Elem(),
	))

	meter := meterProvider.Meter("")
	t.Run("Meter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(meter),
		reflect.TypeOf((*metric.Meter)(nil)).Elem(),
	))

	iC := req[instrument.Int64Counter]{t}.NoErr(meter.Int64Counter(""))
	t.Run("Int64Counter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(iC),
		reflect.TypeOf((*instrument.Int64Counter)(nil)).Elem(),
	))

	fC := req[instrument.Float64Counter]{t}.NoErr(meter.Float64Counter(""))
	t.Run("Float64Counter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(fC),
		reflect.TypeOf((*instrument.Float64Counter)(nil)).Elem(),
	))

	iUDC := req[instrument.Int64UpDownCounter]{t}.NoErr(meter.Int64UpDownCounter(""))
	t.Run("Int64UpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(iUDC),
		reflect.TypeOf((*instrument.Int64UpDownCounter)(nil)).Elem(),
	))

	fUDC := req[instrument.Float64UpDownCounter]{t}.NoErr(meter.Float64UpDownCounter(""))
	t.Run("Float64UpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(fUDC),
		reflect.TypeOf((*instrument.Float64UpDownCounter)(nil)).Elem(),
	))

	iH := req[instrument.Int64Histogram]{t}.NoErr(meter.Int64Histogram(""))
	t.Run("Int64Histogram", assertAllExportedMethodNoPanic(
		reflect.ValueOf(iH),
		reflect.TypeOf((*instrument.Int64Histogram)(nil)).Elem(),
	))

	fH := req[instrument.Float64Histogram]{t}.NoErr(meter.Float64Histogram(""))
	t.Run("Float64Histogram", assertAllExportedMethodNoPanic(
		reflect.ValueOf(fH),
		reflect.TypeOf((*instrument.Float64Histogram)(nil)).Elem(),
	))

	iOC := req[instrument.Int64ObservableCounter]{t}.NoErr(meter.Int64ObservableCounter(""))
	t.Run("Int64ObservableCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(iOC),
		reflect.TypeOf((*instrument.Int64ObservableCounter)(nil)).Elem(),
	))

	fOC := req[instrument.Float64ObservableCounter]{t}.NoErr(meter.Float64ObservableCounter(""))
	t.Run("Float64ObservableCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(fOC),
		reflect.TypeOf((*instrument.Float64ObservableCounter)(nil)).Elem(),
	))

	iOG := req[instrument.Int64ObservableGauge]{t}.NoErr(meter.Int64ObservableGauge(""))
	t.Run("Int64ObservableGauge", assertAllExportedMethodNoPanic(
		reflect.ValueOf(iOG),
		reflect.TypeOf((*instrument.Int64ObservableGauge)(nil)).Elem(),
	))

	fOG := req[instrument.Float64ObservableGauge]{t}.NoErr(meter.Float64ObservableGauge(""))
	t.Run("Float64ObservableGauge", assertAllExportedMethodNoPanic(
		reflect.ValueOf(fOG),
		reflect.TypeOf((*instrument.Float64ObservableGauge)(nil)).Elem(),
	))

	iOUDC := req[instrument.Int64ObservableUpDownCounter]{t}.NoErr(meter.Int64ObservableUpDownCounter(""))
	t.Run("Int64ObservableUpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(iOUDC),
		reflect.TypeOf((*instrument.Int64ObservableUpDownCounter)(nil)).Elem(),
	))

	fOUDC := req[instrument.Float64ObservableUpDownCounter]{t}.NoErr(meter.Float64ObservableUpDownCounter(""))
	t.Run("Float64ObservableUpDownCounter", assertAllExportedMethodNoPanic(
		reflect.ValueOf(fOUDC),
		reflect.TypeOf((*instrument.Float64ObservableUpDownCounter)(nil)).Elem(),
	))

	reg := req[metric.Registration]{t}.NoErr(meter.RegisterCallback(nil, fOC))
	t.Run("Registration", assertAllExportedMethodNoPanic(
		reflect.ValueOf(reg),
		reflect.TypeOf((*metric.Registration)(nil)).Elem(),
	))

	t.Run("Observer", assertAllExportedMethodNoPanic(
		reflect.ValueOf(newObserver()),
		reflect.TypeOf((*metric.Observer)(nil)).Elem(),
	))
}

var (
	ctxType       = reflect.TypeOf((*context.Context)(nil)).Elem()
	backgroundCtx = reflect.ValueOf(context.Background())
)

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
				if aType == ctxType {
					args[i] = backgroundCtx
					continue
				}
				args[i] = reflect.New(aType).Elem()
			}

			assert.NotPanicsf(t, func() {
				_ = m.Call(args)
			}, "%s.%s", rVal.Type().Name(), mType.Name)
		}
	}
}
