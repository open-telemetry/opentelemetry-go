// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var (
	_ metricx.Int64CounterBinder         = (*bindInt64Counter)(nil)
	_ metricx.Float64CounterBinder       = (*bindFloat64Counter)(nil)
	_ metricx.Int64UpDownCounterBinder   = (*bindInt64UpDownCounter)(nil)
	_ metricx.Float64UpDownCounterBinder = (*bindFloat64UpDownCounter)(nil)
)

type testBindOption bool

func (testBindOption) apply(config) config {
	panic("experimental options must not be applied as stable options")
}

func (testBindOption) Experimental() {}

func (o testBindOption) BindingEnabled() bool { return bool(o) }

type testObservedBindOption struct {
	testBindOption
	checked *int
}

func (o testObservedBindOption) BindingEnabled() bool {
	*o.checked++
	return o.testBindOption.BindingEnabled()
}

func TestBindOptionDisabled(t *testing.T) {
	factory, shutdown := newExperimentalMeterFactory([]experimentalOption{testBindOption(false)})

	assert.NotNil(t, factory)
	assert.Nil(t, shutdown)

	provider := NewMeterProvider(testBindOption(false))
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	_, ok := counter.(metricx.Int64CounterBinder)
	assert.False(t, ok)
}

func TestExperimentalMeterFactoryChecksAllBindOptions(t *testing.T) {
	for _, enabled := range [][]bool{{true, false}, {false, true}, {true, true}} {
		t.Run(fmt.Sprint(enabled), func(t *testing.T) {
			checked := 0
			options := []Option{testExperimentalOption{}}
			for _, value := range enabled {
				options = append(options, testObservedBindOption{
					testBindOption: testBindOption(value),
					checked:        &checked,
				})
			}
			provider := NewMeterProvider(options...)
			t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

			assert.Equal(t, len(enabled), checked)
			counter, err := provider.Meter("test").Int64Counter("requests")
			require.NoError(t, err)
			_, ok := counter.(metricx.Int64CounterBinder)
			assert.True(t, ok)
		})
	}
}

func TestBindWithFinishWarns(t *testing.T) {
	var msgs []string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) {
			msgs = append(msgs, args)
		}, funcr.Options{Verbosity: 20}))
		return func() { otel.SetLogger(orig) }
	}(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))))

	provider := NewMeterProvider(testFinishOption(true), testBindOption(true))
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

	assert.Contains(
		t,
		strings.Join(msgs, "\n"),
		"combining experimental Finish and Binding options is not yet supported",
	)
}

func TestExperimentalBindConstruction(t *testing.T) {
	type instKindCase struct {
		name      string
		isFloat   bool
		monotonic bool
		create    func(m api.Meter, name string, opts ...api.InstrumentOption) (any, error)
		bind      func(t *testing.T, inst any, attrs ...attribute.KeyValue) any
		add       func(ctx context.Context, inst any, opts ...api.AddOption)
		enabled   func(ctx context.Context, inst any) bool
	}

	cases := []instKindCase{
		{
			name:      "Int64Counter",
			monotonic: true,
			create: func(m api.Meter, name string, opts ...api.InstrumentOption) (any, error) {
				o := make([]api.Int64CounterOption, len(opts))
				for i, opt := range opts {
					o[i] = opt
				}
				return m.Int64Counter(name, o...)
			},
			bind: func(t *testing.T, inst any, attrs ...attribute.KeyValue) any {
				t.Helper()
				b, ok := inst.(metricx.Int64CounterBinder)
				require.True(t, ok)
				bound := b.Bind(attrs...)
				_, rebindable := bound.(metricx.Int64CounterBinder)
				assert.False(t, rebindable)
				return bound
			},
			add: func(ctx context.Context, inst any, opts ...api.AddOption) {
				inst.(api.Int64Counter).Add(ctx, 3, opts...)
			},
			enabled: func(ctx context.Context, inst any) bool {
				return inst.(api.Int64Counter).Enabled(ctx)
			},
		},
		{
			name:      "Float64Counter",
			isFloat:   true,
			monotonic: true,
			create: func(m api.Meter, name string, opts ...api.InstrumentOption) (any, error) {
				o := make([]api.Float64CounterOption, len(opts))
				for i, opt := range opts {
					o[i] = opt
				}
				return m.Float64Counter(name, o...)
			},
			bind: func(t *testing.T, inst any, attrs ...attribute.KeyValue) any {
				t.Helper()
				b, ok := inst.(metricx.Float64CounterBinder)
				require.True(t, ok)
				bound := b.Bind(attrs...)
				_, rebindable := bound.(metricx.Float64CounterBinder)
				assert.False(t, rebindable)
				return bound
			},
			add: func(ctx context.Context, inst any, opts ...api.AddOption) {
				inst.(api.Float64Counter).Add(ctx, 2.5, opts...)
			},
			enabled: func(ctx context.Context, inst any) bool {
				return inst.(api.Float64Counter).Enabled(ctx)
			},
		},
		{
			name: "Int64UpDownCounter",
			create: func(m api.Meter, name string, opts ...api.InstrumentOption) (any, error) {
				o := make([]api.Int64UpDownCounterOption, len(opts))
				for i, opt := range opts {
					o[i] = opt
				}
				return m.Int64UpDownCounter(name, o...)
			},
			bind: func(t *testing.T, inst any, attrs ...attribute.KeyValue) any {
				t.Helper()
				b, ok := inst.(metricx.Int64UpDownCounterBinder)
				require.True(t, ok)
				bound := b.Bind(attrs...)
				_, rebindable := bound.(metricx.Int64UpDownCounterBinder)
				assert.False(t, rebindable)
				return bound
			},
			add: func(ctx context.Context, inst any, opts ...api.AddOption) {
				inst.(api.Int64UpDownCounter).Add(ctx, -3, opts...)
			},
			enabled: func(ctx context.Context, inst any) bool {
				return inst.(api.Int64UpDownCounter).Enabled(ctx)
			},
		},
		{
			name:    "Float64UpDownCounter",
			isFloat: true,
			create: func(m api.Meter, name string, opts ...api.InstrumentOption) (any, error) {
				o := make([]api.Float64UpDownCounterOption, len(opts))
				for i, opt := range opts {
					o[i] = opt
				}
				return m.Float64UpDownCounter(name, o...)
			},
			bind: func(t *testing.T, inst any, attrs ...attribute.KeyValue) any {
				t.Helper()
				b, ok := inst.(metricx.Float64UpDownCounterBinder)
				require.True(t, ok)
				bound := b.Bind(attrs...)
				_, rebindable := bound.(metricx.Float64UpDownCounterBinder)
				assert.False(t, rebindable)
				return bound
			},
			add: func(ctx context.Context, inst any, opts ...api.AddOption) {
				inst.(api.Float64UpDownCounter).Add(ctx, -2.5, opts...)
			},
			enabled: func(ctx context.Context, inst any) bool {
				return inst.(api.Float64UpDownCounter).Enabled(ctx)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("OptionsCachingAndCaseInsensitiveSharing", func(t *testing.T) {
				reader := NewManualReader()
				provider := NewMeterProvider(WithReader(reader), testBindOption(true))
				t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
				meter := provider.Meter("test")

				opts := []api.InstrumentOption{
					api.WithDescription("Active items"),
					api.WithUnit("{item}"),
					metricx.WithDefaultAttributes("route"),
				}
				inst, err := tc.create(meter, "requests", opts...)
				require.NoError(t, err)
				again, err := tc.create(meter, "requests", opts...)
				require.NoError(t, err)
				assert.Same(t, inst, again)
				assert.True(t, tc.enabled(t.Context(), inst))

				// Case-insensitive duplicate shares the underlying bound aggregation.
				upper, err := tc.create(meter, "Requests", opts...)
				require.NoError(t, err)

				attrs := []attribute.KeyValue{
					attribute.String("zone", "east"),
					attribute.String("route", "/orders"),
				}
				bound1 := tc.bind(t, inst, attrs...)
				bound2 := tc.bind(t, upper, attrs...)
				assert.True(t, tc.enabled(t.Context(), bound1))

				tc.add(t.Context(), bound1)
				tc.add(t.Context(), bound2)

				var rm metricdata.ResourceMetrics
				require.NoError(t, reader.Collect(t.Context(), &rm))
				require.Len(t, rm.ScopeMetrics, 1)
				require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
				m := rm.ScopeMetrics[0].Metrics[0]
				assert.Equal(t, "requests", m.Name)
				assert.Equal(t, "Active items", m.Description)
				assert.Equal(t, "{item}", m.Unit)
				wantAttrs := attribute.NewSet(attribute.String("route", "/orders"))
				if tc.isFloat {
					sum, ok := m.Data.(metricdata.Sum[float64])
					require.True(t, ok)
					assert.Equal(t, tc.monotonic, sum.IsMonotonic)
					require.Len(t, sum.DataPoints, 1)
					assert.Equal(t, wantAttrs, sum.DataPoints[0].Attributes)
					if tc.monotonic {
						assert.Equal(t, float64(5), sum.DataPoints[0].Value)
					} else {
						assert.Equal(t, float64(-5), sum.DataPoints[0].Value)
					}
				} else {
					sum, ok := m.Data.(metricdata.Sum[int64])
					require.True(t, ok)
					assert.Equal(t, tc.monotonic, sum.IsMonotonic)
					require.Len(t, sum.DataPoints, 1)
					assert.Equal(t, wantAttrs, sum.DataPoints[0].Attributes)
					if tc.monotonic {
						assert.Equal(t, int64(6), sum.DataPoints[0].Value)
					} else {
						assert.Equal(t, int64(-6), sum.DataPoints[0].Value)
					}
				}
			})

			t.Run("DroppedAndNoReaderDisabled", func(t *testing.T) {
				noReaderProvider := NewMeterProvider(testBindOption(true))
				t.Cleanup(func() { assert.NoError(t, noReaderProvider.Shutdown(context.WithoutCancel(t.Context()))) })

				inst, err := tc.create(noReaderProvider.Meter("test"), "no_reader")
				require.NoError(t, err)
				assert.False(t, tc.enabled(t.Context(), inst))
				bound := tc.bind(t, inst, attribute.String("k", "v"))
				assert.False(t, tc.enabled(t.Context(), bound))
				assert.NotPanics(t, func() {
					tc.add(t.Context(), inst)
					tc.add(t.Context(), bound)
					tc.add(t.Context(), bound, api.WithAttributes(attribute.String("extra", "1")))
				})

				reader := NewManualReader()
				dropProvider := NewMeterProvider(
					WithReader(reader),
					WithView(NewView(Instrument{Name: "dropped"}, Stream{Aggregation: AggregationDrop{}})),
					testBindOption(true),
				)
				t.Cleanup(func() { assert.NoError(t, dropProvider.Shutdown(context.WithoutCancel(t.Context()))) })

				dropInst, err := tc.create(dropProvider.Meter("test"), "dropped")
				require.NoError(t, err)
				assert.False(t, tc.enabled(t.Context(), dropInst))
				dropBound := tc.bind(t, dropInst, attribute.String("k", "v"))
				assert.False(t, tc.enabled(t.Context(), dropBound))
				assert.NotPanics(t, func() {
					tc.add(t.Context(), dropInst)
					tc.add(t.Context(), dropBound)
					tc.add(t.Context(), dropBound, api.WithAttributes(attribute.String("extra", "1")))
				})
				var rm metricdata.ResourceMetrics
				require.NoError(t, reader.Collect(t.Context(), &rm))
				assert.Empty(t, rm.ScopeMetrics)
			})

			t.Run("Errors", func(t *testing.T) {
				provider := NewMeterProvider(testBindOption(true))
				t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
				inst, err := tc.create(provider.Meter("test"), "invalid name")
				assert.ErrorIs(t, err, ErrInstrumentName)
				assert.NotNil(t, inst)

				reader := NewManualReader()
				failedProvider := NewMeterProvider(
					WithReader(reader),
					testBindOption(true),
					WithView(NewView(Instrument{Name: "active"}, Stream{Aggregation: testUnknownAggregation{}})),
				)
				t.Cleanup(func() { assert.NoError(t, failedProvider.Shutdown(context.WithoutCancel(t.Context()))) })
				failedInst, err := tc.create(failedProvider.Meter("test"), "active")
				assert.ErrorIs(t, err, errUnknownAggregation)
				bound := tc.bind(t, failedInst, attribute.String("k", "v"))
				assert.False(t, tc.enabled(t.Context(), bound))
				assert.NotPanics(t, func() {
					tc.add(t.Context(), failedInst)
					tc.add(t.Context(), bound)
				})
			})
		})
	}
}
