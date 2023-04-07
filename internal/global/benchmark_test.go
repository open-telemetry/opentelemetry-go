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

package global

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/noop"
)

func BenchmarkStartEndSpanNoSDK(b *testing.B) {
	// Compare with BenchmarkStartEndSpan() in
	// ../../sdk/trace/benchmark_test.go.
	ResetForTest(b)
	t := TracerProvider().Tracer("Benchmark StartEndSpan")
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := t.Start(ctx, "/foo")
		span.End()
	}
}

var benchMeter metric.Meter

func BenchmarkMetricMeter(b *testing.B) {
	reset := func() { globalMeterProvider = defaultMeterProvider() }

	b.Run("GetDefault", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			benchMeter = MeterProvider().Meter("")
		}
	})

	b.Run("UseDefault", useMeter(func() metric.Meter {
		globalMeterProvider = defaultMeterProvider()
		return MeterProvider().Meter("UseDefault")
	}))

	b.Run("GetDelegated", func(b *testing.B) {
		SetMeterProvider(noop.NewMeterProvider())
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			benchMeter = MeterProvider().Meter("")
		}
	})
	reset()

	b.Run("UseDelegated", useMeter(func() metric.Meter {
		globalMeterProvider = defaultMeterProvider()
		SetMeterProvider(noop.NewMeterProvider())
		return MeterProvider().Meter("UseDefault")
	}))
	reset()

	b.Run("SetDelegate", func(b *testing.B) {
		del := noop.NewMeterProvider()
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			SetMeterProvider(del)
			reset()
		}
	})
}

var (
	iCtr    instrument.Counter[int64]
	iUDCtr  instrument.UpDownCounter[int64]
	iHist   instrument.Histogram[int64]
	iOCtr   instrument.ObservableCounter[int64]
	iOUDCtr instrument.ObservableUpDownCounter[int64]
	iOGauge instrument.ObservableGauge[int64]

	fCtr    instrument.Counter[float64]
	fUDCtr  instrument.UpDownCounter[float64]
	fHist   instrument.Histogram[float64]
	fOCtr   instrument.ObservableCounter[float64]
	fOUDCtr instrument.ObservableUpDownCounter[float64]
	fOGauge instrument.ObservableGauge[float64]
)

func useMeter(newMeter func() metric.Meter) func(*testing.B) {
	return func(b *testing.B) {
		b.Run("CreateInt64Counter", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				iCtr, _ = m.Int64Counter(strconv.Itoa(n))
			}
		})
		b.Run("CreateInt64UpDownCounter", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				iUDCtr, _ = m.Int64UpDownCounter(strconv.Itoa(n))
			}
		})
		b.Run("CreateInt64Histogram", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				iHist, _ = m.Int64Histogram(strconv.Itoa(n))
			}
		})
		b.Run("CreateInt64ObservableCounter", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				iOCtr, _ = m.Int64ObservableCounter(strconv.Itoa(n))
			}
		})
		b.Run("CreateInt64ObservableUpDownCounter", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				iOUDCtr, _ = m.Int64ObservableUpDownCounter(strconv.Itoa(n))
			}
		})
		b.Run("CreateInt64ObservableGauge", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				iOGauge, _ = m.Int64ObservableGauge(strconv.Itoa(n))
			}
		})
		b.Run("CreateFloat64Counter", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				fCtr, _ = m.Float64Counter(strconv.Itoa(n))
			}
		})
		b.Run("CreateFloat64UpDownCounter", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				fUDCtr, _ = m.Float64UpDownCounter(strconv.Itoa(n))
			}
		})
		b.Run("CreateFloat64Histogram", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				fHist, _ = m.Float64Histogram(strconv.Itoa(n))
			}
		})
		b.Run("CreateFloat64ObservableCounter", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				fOCtr, _ = m.Float64ObservableCounter(strconv.Itoa(n))
			}
		})
		b.Run("CreateFloat64ObservableUpDownCounter", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				fOUDCtr, _ = m.Float64ObservableUpDownCounter(strconv.Itoa(n))
			}
		})
		b.Run("CreateFloat64ObservableGauge", func(b *testing.B) {
			m := newMeter()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				fOGauge, _ = m.Float64ObservableGauge(strconv.Itoa(n))
			}
		})

		b.Run("RegisterCallback", func(b *testing.B) {
			m := newMeter()
			iOCtr, _ := m.Int64ObservableCounter("")
			iOUDCtr, _ := m.Int64ObservableUpDownCounter("")
			iOGauge, _ := m.Int64ObservableGauge("")
			fOCtr, _ := m.Float64ObservableCounter("")
			fOUDCtr, _ := m.Float64ObservableUpDownCounter("")
			fOGauge, _ := m.Float64ObservableGauge("")
			fn := func(_ context.Context, o metric.Observer) error {
				return nil
			}
			b.ReportAllocs()
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				_, _ = m.RegisterCallback(fn, iOCtr, iOUDCtr, iOGauge, fOCtr, fOUDCtr, fOGauge)
			}
		})
	}
}

func BenchmarkMetricInstruments(b *testing.B) {
	b.Run("Default", benchmarkMetricInstruments(func() metric.Meter {
		globalMeterProvider = defaultMeterProvider()
		return MeterProvider().Meter("BenchmarkMetricInstruments")
	}))
	b.Run("Delegated", benchmarkMetricInstruments(func() metric.Meter {
		globalMeterProvider = defaultMeterProvider()
		SetMeterProvider(noop.NewMeterProvider())
		return MeterProvider().Meter("BenchmarkMetricInstruments")
	}))
	globalMeterProvider = defaultMeterProvider()
}

func benchmarkMetricInstruments(newMeter func() metric.Meter) func(b *testing.B) {
	return func(b *testing.B) {
		ctx := context.Background()

		b.Run("Int64Counter.Add", func(b *testing.B) {
			iCtr, err := newMeter().Int64Counter("")
			assert.NoError(b, err)
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				iCtr.Add(ctx, 1)
			}
		})

		b.Run("Int64UpDownCounter.Add", func(b *testing.B) {
			iUDCtr, err := newMeter().Int64UpDownCounter("")
			assert.NoError(b, err)
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				iUDCtr.Add(ctx, 1)
			}
		})

		b.Run("Int64Histogram.Record", func(b *testing.B) {
			iHist, err := newMeter().Int64Histogram("")
			assert.NoError(b, err)
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				iHist.Record(ctx, 1)
			}
		})

		b.Run("Int64ObservableCounter.Observe", func(b *testing.B) {
			iOCtr, err := newMeter().Int64ObservableCounter("")
			assert.NoError(b, err)
			obsv := noop.Observer{}
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				obsv.ObserveInt64(iOCtr, 1)
			}
		})

		b.Run("Int64ObservableUpDownCounter.Observe", func(b *testing.B) {
			iOUDCtr, err := newMeter().Int64ObservableUpDownCounter("")
			assert.NoError(b, err)
			obsv := noop.Observer{}
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				obsv.ObserveInt64(iOUDCtr, 1)
			}
		})

		b.Run("Int64ObservableGauge.Observe", func(b *testing.B) {
			iOGauge, err := newMeter().Int64ObservableGauge("")
			assert.NoError(b, err)
			obsv := noop.Observer{}
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				obsv.ObserveInt64(iOGauge, 1)
			}
		})

		b.Run("Float64Counter.Add", func(b *testing.B) {
			iCtr, err := newMeter().Float64Counter("")
			assert.NoError(b, err)
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				iCtr.Add(ctx, 1)
			}
		})

		b.Run("Float64UpDownCounter.Add", func(b *testing.B) {
			iUDCtr, err := newMeter().Float64UpDownCounter("")
			assert.NoError(b, err)
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				iUDCtr.Add(ctx, 1)
			}
		})

		b.Run("Float64Histogram.Record", func(b *testing.B) {
			iHist, err := newMeter().Float64Histogram("")
			assert.NoError(b, err)
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				iHist.Record(ctx, 1)
			}
		})

		b.Run("Float64ObservableCounter.Observe", func(b *testing.B) {
			iOCtr, err := newMeter().Float64ObservableCounter("")
			assert.NoError(b, err)
			obsv := noop.Observer{}
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				obsv.ObserveFloat64(iOCtr, 1)
			}
		})

		b.Run("Float64ObservableUpDownCounter.Observe", func(b *testing.B) {
			iOUDCtr, err := newMeter().Float64ObservableUpDownCounter("")
			assert.NoError(b, err)
			obsv := noop.Observer{}
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				obsv.ObserveFloat64(iOUDCtr, 1)
			}
		})

		b.Run("Float64ObservableGauge.Observe", func(b *testing.B) {
			iOGauge, err := newMeter().Float64ObservableGauge("")
			assert.NoError(b, err)
			obsv := noop.Observer{}
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				obsv.ObserveFloat64(iOGauge, 1)
			}
		})
	}
}
