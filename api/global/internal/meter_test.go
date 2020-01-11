package internal_test

import (
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/exporter/metric/stdout"
	metrictest "go.opentelemetry.io/otel/internal/metric"
)

func TestDirect(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	meter1 := global.Scope().WithNamespace("test1").Meter()
	meter2 := global.Scope().WithNamespace("test2").Meter()
	lvals1 := key.String("A", "B")
	lvals2 := key.String("C", "D")
	lvals3 := key.String("E", "F")

	counter := meter1.NewInt64Counter("test.counter")
	counter.Add(ctx, 1, lvals1)
	counter.Add(ctx, 1, lvals1)

	gauge := meter1.NewInt64Gauge("test.gauge")
	gauge.Set(ctx, 1, lvals2)
	gauge.Set(ctx, 2, lvals2)

	measure := meter1.NewFloat64Measure("test.measure")
	measure.Record(ctx, 1, lvals1)
	measure.Record(ctx, 2, lvals1)

	second := meter2.NewFloat64Measure("test.second")
	second.Record(ctx, 1, lvals3)
	second.Record(ctx, 2, lvals3)

	sdk := metrictest.NewMeter()
	global.SetScope(scope.WithMeterSDK(sdk))

	counter.Add(ctx, 1, lvals1)
	gauge.Set(ctx, 3, lvals2)
	measure.Record(ctx, 3, lvals1)
	second.Record(ctx, 3, lvals3)

	require.Equal(t, 4, len(sdk.MeasurementBatches))

	require.Equal(t, []core.KeyValue{lvals1}, sdk.MeasurementBatches[0].Labels.Ordered())
	require.Equal(t, 1, len(sdk.MeasurementBatches[0].Measurements))
	require.Equal(t, core.NewInt64Number(1),
		sdk.MeasurementBatches[0].Measurements[0].Number)
	require.Equal(t, "test1/test.counter",
		sdk.MeasurementBatches[0].Measurements[0].Instrument.Name.String())

	require.Equal(t, []core.KeyValue{lvals2}, sdk.MeasurementBatches[1].Labels.Ordered())
	require.Equal(t, 1, len(sdk.MeasurementBatches[1].Measurements))
	require.Equal(t, core.NewInt64Number(3),
		sdk.MeasurementBatches[1].Measurements[0].Number)
	require.Equal(t, "test1/test.gauge",
		sdk.MeasurementBatches[1].Measurements[0].Instrument.Name.String())

	require.Equal(t, []core.KeyValue{lvals1}, sdk.MeasurementBatches[2].Labels.Ordered())
	require.Equal(t, 1, len(sdk.MeasurementBatches[2].Measurements))
	require.Equal(t, core.NewFloat64Number(3),
		sdk.MeasurementBatches[2].Measurements[0].Number)
	require.Equal(t, "test1/test.measure",
		sdk.MeasurementBatches[2].Measurements[0].Instrument.Name.String())

	require.Equal(t, []core.KeyValue{lvals3}, sdk.MeasurementBatches[3].Labels.Ordered())
	require.Equal(t, 1, len(sdk.MeasurementBatches[3].Measurements))
	require.Equal(t, core.NewFloat64Number(3),
		sdk.MeasurementBatches[3].Measurements[0].Number)
	require.Equal(t, "test2/test.second",
		sdk.MeasurementBatches[3].Measurements[0].Instrument.Name.String())
}

func TestBound(t *testing.T) {
	internal.ResetForTest()

	// Note: this test uses oppsite Float64/Int64 number kinds
	// vs. the above, to cover all the instruments.
	ctx := context.Background()
	glob := global.Scope().WithNamespace("test").Meter()
	lvals1 := key.String("A", "B")
	lvals2 := key.String("C", "D")

	counter := glob.NewFloat64Counter("test.counter")
	boundC := counter.Bind(ctx, lvals1)
	boundC.Add(ctx, 1)
	boundC.Add(ctx, 1)

	gauge := glob.NewFloat64Gauge("test.gauge")
	boundG := gauge.Bind(ctx, lvals2)
	boundG.Set(ctx, 1)
	boundG.Set(ctx, 2)

	measure := glob.NewInt64Measure("test.measure")
	boundM := measure.Bind(ctx, lvals1)
	boundM.Record(ctx, 1)
	boundM.Record(ctx, 2)

	sdk := metrictest.NewMeter()
	global.SetScope(scope.WithMeterSDK(sdk))

	boundC.Add(ctx, 1)
	boundG.Set(ctx, 3)
	boundM.Record(ctx, 3)

	require.Equal(t, 3, len(sdk.MeasurementBatches))

	require.Equal(t, []core.KeyValue{lvals1}, sdk.MeasurementBatches[0].Labels.Ordered())
	require.Equal(t, 1, len(sdk.MeasurementBatches[0].Measurements))
	require.Equal(t, core.NewFloat64Number(1),
		sdk.MeasurementBatches[0].Measurements[0].Number)
	require.Equal(t, "test/test.counter",
		sdk.MeasurementBatches[0].Measurements[0].Instrument.Name.String())

	require.Equal(t, []core.KeyValue{lvals2}, sdk.MeasurementBatches[1].Labels.Ordered())
	require.Equal(t, 1, len(sdk.MeasurementBatches[1].Measurements))
	require.Equal(t, core.NewFloat64Number(3),
		sdk.MeasurementBatches[1].Measurements[0].Number)
	require.Equal(t, "test/test.gauge",
		sdk.MeasurementBatches[1].Measurements[0].Instrument.Name.String())

	require.Equal(t, []core.KeyValue{lvals1}, sdk.MeasurementBatches[2].Labels.Ordered())
	require.Equal(t, 1, len(sdk.MeasurementBatches[2].Measurements))
	require.Equal(t, core.NewInt64Number(3),
		sdk.MeasurementBatches[2].Measurements[0].Number)
	require.Equal(t, "test/test.measure",
		sdk.MeasurementBatches[2].Measurements[0].Instrument.Name.String())

	boundC.Unbind()
	boundG.Unbind()
	boundM.Unbind()
}

func TestUnbind(t *testing.T) {
	// Tests Unbind with SDK never installed.
	internal.ResetForTest()

	ctx := context.Background()

	glob := global.Scope().WithNamespace("test").Meter()
	lvals1 := key.New("A").String("B")
	lvals2 := key.New("C").String("D")

	counter := glob.NewFloat64Counter("test/counter")
	boundC := counter.Bind(ctx, lvals1)

	gauge := glob.NewFloat64Gauge("test/gauge")
	boundG := gauge.Bind(ctx, lvals2)

	measure := glob.NewInt64Measure("test/measure")
	boundM := measure.Bind(ctx, lvals1)

	boundC.Unbind()
	boundG.Unbind()
	boundM.Unbind()
}

func TestDefaultSDK(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	meter1 := global.Scope().WithNamespace("builtin").Meter()
	lvals1 := key.String("A", "B")

	counter := meter1.NewInt64Counter("count.b")
	counter.Add(ctx, 1, lvals1)
	counter.Add(ctx, 1, lvals1)

	in, out := io.Pipe()
	pusher, err := stdout.NewExportPipeline(stdout.Config{
		Writer:         out,
		DoNotPrintTime: true,
	})
	if err != nil {
		panic(err)
	}
	global.SetScope(scope.WithMeterSDK(pusher.Meter()))

	counter.Add(ctx, 1, lvals1)

	ch := make(chan string)
	go func() {
		data, _ := ioutil.ReadAll(in)
		ch <- string(data)
	}()

	pusher.Stop()
	out.Close()

	require.Equal(t, `{"updates":[{"name":"builtin/count.b{A=B}","sum":1}]}
`, <-ch)
}
