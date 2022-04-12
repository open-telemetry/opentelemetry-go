package controllertest

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type errorCatcher struct {
	lock   sync.Mutex
	errors []error
}

func (e *errorCatcher) Handle(err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.errors = append(e.errors, err)
}

func TestEndToEnd(t *testing.T) {
	h := &errorCatcher{}
	otel.SetErrorHandler(h)

	meter := global.Meter("go.opentelemetry.io/otel/sdk/metric/controller/controllertest_EndToEnd")
	gauge, err := meter.AsyncInt64().Gauge("test")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{gauge}, func(context.Context) {})
	require.NoError(t, err)

	c := controller.New(basic.NewFactory(simple.NewWithInexpensiveDistribution(), aggregation.CumulativeTemporalitySelector()))

	global.SetMeterProvider(c)

	gauge, err = meter.AsyncInt64().Gauge("test2")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{gauge}, func(context.Context) {})
	require.NoError(t, err)

	h.lock.Lock()
	require.Len(t, h.errors, 0)

}
