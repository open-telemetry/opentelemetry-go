package internal_test

import (
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/exporters/metric/stdout"
	metrictest "go.opentelemetry.io/otel/internal/metric"
)

func TestDirect(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	meter1 := global.MeterProvider().Meter("test1")
	meter2 := global.MeterProvider().Meter("test2")
	lvals1 := key.String("A", "B")
	labels1 := meter1.Labels(lvals1)
	lvals2 := key.String("C", "D")
	labels2 := meter1.Labels(lvals2)
	lvals3 := key.String("E", "F")
	labels3 := meter2.Labels(lvals3)

	counter := meter1.NewInt64Counter("test.counter")
	counter.Add(ctx, 1, labels1)
	counter.Add(ctx, 1, labels1)

	gauge := meter1.NewInt64Gauge("test.gauge")
	gauge.Set(ctx, 1, labels2)
	gauge.Set(ctx, 2, labels2)

	measure := meter1.NewFloat64Measure("test.measure")
	measure.Record(ctx, 1, labels1)
	measure.Record(ctx, 2, labels1)

	second := meter2.NewFloat64Measure("test.second")
	second.Record(ctx, 1, labels3)
	second.Record(ctx, 2, labels3)

	sdk := metrictest.NewProvider()
	global.SetMeterProvider(sdk)

	counter.Add(ctx, 1, labels1)
	gauge.Set(ctx, 3, labels2)
	measure.Record(ctx, 3, labels1)
	second.Record(ctx, 3, labels3)

	mock := sdk.Meter("test1").(*metrictest.Meter)
	require.Equal(t, 3, len(mock.MeasurementBatches))

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[0].LabelSet.Labels)
	require.Equal(t, 1, len(mock.MeasurementBatches[0].Measurements))
	require.Equal(t, core.NewInt64Number(1),
		mock.MeasurementBatches[0].Measurements[0].Number)
	require.Equal(t, "test.counter",
		mock.MeasurementBatches[0].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals2.Key: lvals2.Value,
	}, mock.MeasurementBatches[1].LabelSet.Labels)
	require.Equal(t, 1, len(mock.MeasurementBatches[1].Measurements))
	require.Equal(t, core.NewInt64Number(3),
		mock.MeasurementBatches[1].Measurements[0].Number)
	require.Equal(t, "test.gauge",
		mock.MeasurementBatches[1].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[2].LabelSet.Labels)
	require.Equal(t, 1, len(mock.MeasurementBatches[2].Measurements))
	require.Equal(t, core.NewFloat64Number(3),
		mock.MeasurementBatches[2].Measurements[0].Number)
	require.Equal(t, "test.measure",
		mock.MeasurementBatches[2].Measurements[0].Instrument.Name)

	// This tests the second Meter instance
	mock = sdk.Meter("test2").(*metrictest.Meter)
	require.Equal(t, 1, len(mock.MeasurementBatches))

	require.Equal(t, map[core.Key]core.Value{
		lvals3.Key: lvals3.Value,
	}, mock.MeasurementBatches[0].LabelSet.Labels)
	require.Equal(t, 1, len(mock.MeasurementBatches[0].Measurements))
	require.Equal(t, core.NewFloat64Number(3),
		mock.MeasurementBatches[0].Measurements[0].Number)
	require.Equal(t, "test.second",
		mock.MeasurementBatches[0].Measurements[0].Instrument.Name)
}

func TestBound(t *testing.T) {
	internal.ResetForTest()

	// Note: this test uses opposite Float64/Int64 number kinds
	// vs. the above, to cover all the instruments.
	ctx := context.Background()
	glob := global.MeterProvider().Meter("test")
	lvals1 := key.String("A", "B")
	labels1 := glob.Labels(lvals1)
	lvals2 := key.String("C", "D")
	labels2 := glob.Labels(lvals2)

	counter := glob.NewFloat64Counter("test.counter")
	boundC := counter.Bind(labels1)
	boundC.Add(ctx, 1)
	boundC.Add(ctx, 1)

	gauge := glob.NewFloat64Gauge("test.gauge")
	boundG := gauge.Bind(labels2)
	boundG.Set(ctx, 1)
	boundG.Set(ctx, 2)

	measure := glob.NewInt64Measure("test.measure")
	boundM := measure.Bind(labels1)
	boundM.Record(ctx, 1)
	boundM.Record(ctx, 2)

	sdk := metrictest.NewProvider()
	global.SetMeterProvider(sdk)

	boundC.Add(ctx, 1)
	boundG.Set(ctx, 3)
	boundM.Record(ctx, 3)

	mock := sdk.Meter("test").(*metrictest.Meter)
	require.Equal(t, 3, len(mock.MeasurementBatches))

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[0].LabelSet.Labels)
	require.Equal(t, 1, len(mock.MeasurementBatches[0].Measurements))
	require.Equal(t, core.NewFloat64Number(1),
		mock.MeasurementBatches[0].Measurements[0].Number)
	require.Equal(t, "test.counter",
		mock.MeasurementBatches[0].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals2.Key: lvals2.Value,
	}, mock.MeasurementBatches[1].LabelSet.Labels)
	require.Equal(t, 1, len(mock.MeasurementBatches[1].Measurements))
	require.Equal(t, core.NewFloat64Number(3),
		mock.MeasurementBatches[1].Measurements[0].Number)
	require.Equal(t, "test.gauge",
		mock.MeasurementBatches[1].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[2].LabelSet.Labels)
	require.Equal(t, 1, len(mock.MeasurementBatches[2].Measurements))
	require.Equal(t, core.NewInt64Number(3),
		mock.MeasurementBatches[2].Measurements[0].Number)
	require.Equal(t, "test.measure",
		mock.MeasurementBatches[2].Measurements[0].Instrument.Name)

	boundC.Unbind()
	boundG.Unbind()
	boundM.Unbind()
}

func TestUnbind(t *testing.T) {
	// Tests Unbind with SDK never installed.
	internal.ResetForTest()

	glob := global.MeterProvider().Meter("test")
	lvals1 := key.New("A").String("B")
	labels1 := glob.Labels(lvals1)
	lvals2 := key.New("C").String("D")
	labels2 := glob.Labels(lvals2)

	counter := glob.NewFloat64Counter("test.counter")
	boundC := counter.Bind(labels1)

	gauge := glob.NewFloat64Gauge("test.gauge")
	boundG := gauge.Bind(labels2)

	measure := glob.NewInt64Measure("test.measure")
	boundM := measure.Bind(labels1)

	boundC.Unbind()
	boundG.Unbind()
	boundM.Unbind()
}

func TestDefaultSDK(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	meter1 := global.MeterProvider().Meter("builtin")
	lvals1 := key.String("A", "B")
	labels1 := meter1.Labels(lvals1)

	counter := meter1.NewInt64Counter("test.builtin")
	counter.Add(ctx, 1, labels1)
	counter.Add(ctx, 1, labels1)

	in, out := io.Pipe()
	pusher, err := stdout.InstallNewPipeline(stdout.Config{
		Writer:         out,
		DoNotPrintTime: true,
	})
	if err != nil {
		panic(err)
	}

	counter.Add(ctx, 1, labels1)

	ch := make(chan string)
	go func() {
		data, _ := ioutil.ReadAll(in)
		ch <- string(data)
	}()

	pusher.Stop()
	out.Close()

	require.Equal(t, `{"updates":[{"name":"test.builtin{A=B}","sum":1}]}
`, <-ch)
}

func TestUnbindThenRecordOne(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	sdk := metrictest.NewProvider()
	meter := global.MeterProvider().Meter("test")
	counter := meter.NewInt64Counter("test.counter")
	boundC := counter.Bind(meter.Labels())
	global.SetMeterProvider(sdk)
	boundC.Unbind()

	require.NotPanics(t, func() {
		boundC.Add(ctx, 1)
	})
	mock := global.MeterProvider().Meter("test").(*metrictest.Meter)
	require.Equal(t, 0, len(mock.MeasurementBatches))
}
