package internal_test

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/stdout"
	metrictest "go.opentelemetry.io/otel/internal/metric"
)

func TestDirect(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	meter1 := global.Meter("test1")
	meter2 := global.Meter("test2")
	lvals1 := key.String("A", "B")
	labels1 := meter1.Labels(lvals1)
	lvals2 := key.String("C", "D")
	labels2 := meter1.Labels(lvals2)
	lvals3 := key.String("E", "F")
	labels3 := meter2.Labels(lvals3)

	counter := Must(meter1).NewInt64Counter("test.counter")
	counter.Add(ctx, 1, labels1)
	counter.Add(ctx, 1, labels1)

	measure := Must(meter1).NewFloat64Measure("test.measure")
	measure.Record(ctx, 1, labels1)
	measure.Record(ctx, 2, labels1)

	_ = Must(meter1).RegisterFloat64Observer("test.observer.float", func(result metric.Float64ObserverResult) {
		result.Observe(1., labels1)
		result.Observe(2., labels2)
	})

	_ = Must(meter1).RegisterInt64Observer("test.observer.int", func(result metric.Int64ObserverResult) {
		result.Observe(1, labels1)
		result.Observe(2, labels2)
	})

	second := Must(meter2).NewFloat64Measure("test.second")
	second.Record(ctx, 1, labels3)
	second.Record(ctx, 2, labels3)

	sdk := metrictest.NewProvider()
	global.SetMeterProvider(sdk)

	counter.Add(ctx, 1, labels1)
	measure.Record(ctx, 3, labels1)
	second.Record(ctx, 3, labels3)

	mock := sdk.Meter("test1").(*metrictest.Meter)
	mock.RunObservers()
	require.Len(t, mock.MeasurementBatches, 6)

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[0].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[0].Measurements, 1)
	require.Equal(t, int64(1),
		mock.MeasurementBatches[0].Measurements[0].Number.AsInt64())
	require.Equal(t, "test.counter",
		mock.MeasurementBatches[0].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[1].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[1].Measurements, 1)
	require.InDelta(t, float64(3),
		mock.MeasurementBatches[1].Measurements[0].Number.AsFloat64(),
		0.01)
	require.Equal(t, "test.measure",
		mock.MeasurementBatches[1].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[2].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[2].Measurements, 1)
	require.InDelta(t, float64(1),
		mock.MeasurementBatches[2].Measurements[0].Number.AsFloat64(),
		0.01)
	require.Equal(t, "test.observer.float",
		mock.MeasurementBatches[2].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals2.Key: lvals2.Value,
	}, mock.MeasurementBatches[3].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[3].Measurements, 1)
	require.InDelta(t, float64(2),
		mock.MeasurementBatches[3].Measurements[0].Number.AsFloat64(),
		0.01)
	require.Equal(t, "test.observer.float",
		mock.MeasurementBatches[3].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[4].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[4].Measurements, 1)
	require.Equal(t, int64(1),
		mock.MeasurementBatches[4].Measurements[0].Number.AsInt64())
	require.Equal(t, "test.observer.int",
		mock.MeasurementBatches[4].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals2.Key: lvals2.Value,
	}, mock.MeasurementBatches[5].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[5].Measurements, 1)
	require.Equal(t, int64(2),
		mock.MeasurementBatches[5].Measurements[0].Number.AsInt64())
	require.Equal(t, "test.observer.int",
		mock.MeasurementBatches[5].Measurements[0].Instrument.Name)

	// This tests the second Meter instance
	mock = sdk.Meter("test2").(*metrictest.Meter)
	require.Len(t, mock.MeasurementBatches, 1)

	require.Equal(t, map[core.Key]core.Value{
		lvals3.Key: lvals3.Value,
	}, mock.MeasurementBatches[0].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[0].Measurements, 1)
	require.InDelta(t, float64(3),
		mock.MeasurementBatches[0].Measurements[0].Number.AsFloat64(),
		0.01)
	require.Equal(t, "test.second",
		mock.MeasurementBatches[0].Measurements[0].Instrument.Name)
}

func TestBound(t *testing.T) {
	internal.ResetForTest()

	// Note: this test uses opposite Float64/Int64 number kinds
	// vs. the above, to cover all the instruments.
	ctx := context.Background()
	glob := global.Meter("test")
	lvals1 := key.String("A", "B")
	labels1 := glob.Labels(lvals1)

	counter := Must(glob).NewFloat64Counter("test.counter")
	boundC := counter.Bind(labels1)
	boundC.Add(ctx, 1)
	boundC.Add(ctx, 1)

	measure := Must(glob).NewInt64Measure("test.measure")
	boundM := measure.Bind(labels1)
	boundM.Record(ctx, 1)
	boundM.Record(ctx, 2)

	sdk := metrictest.NewProvider()
	global.SetMeterProvider(sdk)

	boundC.Add(ctx, 1)
	boundM.Record(ctx, 3)

	mock := sdk.Meter("test").(*metrictest.Meter)
	require.Len(t, mock.MeasurementBatches, 2)

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[0].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[0].Measurements, 1)
	require.InDelta(t, float64(1),
		mock.MeasurementBatches[0].Measurements[0].Number.AsFloat64(),
		0.01)
	require.Equal(t, "test.counter",
		mock.MeasurementBatches[0].Measurements[0].Instrument.Name)

	require.Equal(t, map[core.Key]core.Value{
		lvals1.Key: lvals1.Value,
	}, mock.MeasurementBatches[1].LabelSet.Labels)
	require.Len(t, mock.MeasurementBatches[1].Measurements, 1)
	require.Equal(t, int64(3),
		mock.MeasurementBatches[1].Measurements[0].Number.AsInt64())
	require.Equal(t, "test.measure",
		mock.MeasurementBatches[1].Measurements[0].Instrument.Name)

	boundC.Unbind()
	boundM.Unbind()
}

func TestUnbind(t *testing.T) {
	// Tests Unbind with SDK never installed.
	internal.ResetForTest()

	glob := global.Meter("test")
	lvals1 := key.New("A").String("B")
	labels1 := glob.Labels(lvals1)

	counter := Must(glob).NewFloat64Counter("test.counter")
	boundC := counter.Bind(labels1)

	measure := Must(glob).NewInt64Measure("test.measure")
	boundM := measure.Bind(labels1)

	observerInt := Must(glob).RegisterInt64Observer("test.observer.int", nil)
	observerFloat := Must(glob).RegisterFloat64Observer("test.observer.float", nil)

	boundC.Unbind()
	boundM.Unbind()
	observerInt.Unregister()
	observerFloat.Unregister()
}

func TestDefaultSDK(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	meter1 := global.Meter("builtin")
	lvals1 := key.String("A", "B")
	labels1 := meter1.Labels(lvals1)

	counter := Must(meter1).NewInt64Counter("test.builtin")
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

	meter := global.Meter("test")
	counter := Must(meter).NewInt64Counter("test.counter")
	boundC := counter.Bind(meter.Labels())
	global.SetMeterProvider(sdk)
	boundC.Unbind()

	require.NotPanics(t, func() {
		boundC.Add(ctx, 1)
	})
	mock := global.Meter("test").(*metrictest.Meter)
	require.Equal(t, 0, len(mock.MeasurementBatches))
}

type meterProviderWithConstructorError struct {
	metric.Provider
}

type meterWithConstructorError struct {
	metric.Meter
}

func (m *meterProviderWithConstructorError) Meter(name string) metric.Meter {
	return &meterWithConstructorError{m.Provider.Meter(name)}
}

func (m *meterWithConstructorError) NewInt64Counter(name string, cos ...metric.CounterOptionApplier) (metric.Int64Counter, error) {
	return metric.Int64Counter{}, errors.New("constructor error")
}

func TestErrorInDeferredConstructor(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	meter := global.MeterProvider().Meter("builtin")

	c1 := Must(meter).NewInt64Counter("test")
	c2 := Must(meter).NewInt64Counter("test")

	sdk := &meterProviderWithConstructorError{metrictest.NewProvider()}

	require.Panics(t, func() {
		global.SetMeterProvider(sdk)
	})

	c1.Add(ctx, 1, meter.Labels())
	c2.Add(ctx, 2, meter.Labels())
}
