package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/metric/registry"
)

type (
	newFunc func(name, libraryName string) (metric.InstrumentImpl, error)
)

var (
	allNew = map[string]newFunc{
		"counter.int64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewInt64Counter(name))
		},
		"counter.float64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewFloat64Counter(name))
		},
		"measure.int64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewInt64Measure(name))
		},
		"measure.float64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewFloat64Measure(name))
		},
		"observer.int64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).RegisterInt64Observer(name, func(metric.Int64ObserverResult) {}))
		},
		"observer.float64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).RegisterFloat64Observer(name, func(metric.Float64ObserverResult) {}))
		},
	}
)

func unwrap(impl interface{}, err error) (metric.InstrumentImpl, error) {
	if impl == nil {
		return nil, err
	}
	if s, ok := impl.(interface {
		SyncImpl() metric.SyncImpl
	}); ok {
		return s.SyncImpl(), err
	}
	if a, ok := impl.(interface {
		AsyncImpl() metric.AsyncImpl
	}); ok {
		return a.AsyncImpl(), err
	}
	return nil, err
}

func TestRegistrySameInstruments(t *testing.T) {
	for _, nf := range allNew {
		ResetForTest()
		inst1, err1 := nf("this", "meter")
		inst2, err2 := nf("this", "meter")

		require.Nil(t, err1)
		require.Nil(t, err2)
		require.Equal(t, inst1, inst2)
	}
}

func TestRegistryDifferentNamespace(t *testing.T) {
	for _, nf := range allNew {
		ResetForTest()
		inst1, err1 := nf("this", "meter1")
		inst2, err2 := nf("this", "meter2")

		require.Nil(t, err1)
		require.Nil(t, err2)
		require.NotEqual(t, inst1, inst2)
	}
}

func TestRegistryDiffInstruments(t *testing.T) {
	for origName, origf := range allNew {
		ResetForTest()

		_, err := origf("this", "super")
		require.Nil(t, err)

		for newName, nf := range allNew {
			if newName == origName {
				continue
			}

			other, err := nf("this", "super")
			require.NotNil(t, err)
			require.NotNil(t, other)
			require.True(t, errors.Is(err, registry.ErrMetricKindMismatch))
			require.Contains(t, err.Error(), "super")
		}
	}
}
