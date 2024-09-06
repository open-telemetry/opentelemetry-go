// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	ottest "go.opentelemetry.io/otel/sdk/internal/internaltest"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type reader struct {
	producer          sdkProducer
	externalProducers []Producer
	temporalityFunc   TemporalitySelector
	aggregationFunc   AggregationSelector
	collectFunc       func(context.Context, *metricdata.ResourceMetrics) error
	forceFlushFunc    func(context.Context) error
	shutdownFunc      func(context.Context) error
}

const envVarResourceAttributes = "OTEL_RESOURCE_ATTRIBUTES"

var _ Reader = (*reader)(nil)

func (r *reader) aggregation(kind InstrumentKind) Aggregation { // nolint:revive  // import-shadow for method scoped by type.
	return r.aggregationFunc(kind)
}

func (r *reader) register(p sdkProducer)      { r.producer = p }
func (r *reader) RegisterProducer(p Producer) { r.externalProducers = append(r.externalProducers, p) }
func (r *reader) temporality(kind InstrumentKind) metricdata.Temporality {
	return r.temporalityFunc(kind)
}

func (r *reader) Collect(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	return r.collectFunc(ctx, rm)
}
func (r *reader) ForceFlush(ctx context.Context) error { return r.forceFlushFunc(ctx) }
func (r *reader) Shutdown(ctx context.Context) error   { return r.shutdownFunc(ctx) }

func TestConfigReaderSignalsEmpty(t *testing.T) {
	f, s := config{}.readerSignals()

	require.NotNil(t, f)
	require.NotNil(t, s)

	ctx := context.Background()
	assert.Nil(t, f(ctx))
	assert.Nil(t, s(ctx))
	assert.ErrorIs(t, s(ctx), ErrReaderShutdown)
}

func TestConfigReaderSignalsForwarded(t *testing.T) {
	var flush, sdown int
	r := &reader{
		forceFlushFunc: func(ctx context.Context) error {
			flush++
			return nil
		},
		shutdownFunc: func(ctx context.Context) error {
			sdown++
			return nil
		},
	}
	c := newConfig([]Option{WithReader(r)})
	f, s := c.readerSignals()

	require.NotNil(t, f)
	require.NotNil(t, s)

	ctx := context.Background()
	assert.NoError(t, f(ctx))
	assert.NoError(t, f(ctx))
	assert.NoError(t, s(ctx))
	assert.ErrorIs(t, s(ctx), ErrReaderShutdown)

	assert.Equal(t, 2, flush, "flush not called 2 times")
	assert.Equal(t, 1, sdown, "shutdown not called 1 time")
}

func TestConfigReaderSignalsForwardedErrors(t *testing.T) {
	r := &reader{
		forceFlushFunc: func(ctx context.Context) error { return assert.AnError },
		shutdownFunc:   func(ctx context.Context) error { return assert.AnError },
	}
	c := newConfig([]Option{WithReader(r)})
	f, s := c.readerSignals()

	require.NotNil(t, f)
	require.NotNil(t, s)

	ctx := context.Background()
	assert.ErrorIs(t, f(ctx), assert.AnError)
	assert.ErrorIs(t, s(ctx), assert.AnError)
	assert.ErrorIs(t, s(ctx), ErrReaderShutdown)
}

func TestUnifyMultiError(t *testing.T) {
	f := func(context.Context) error { return assert.AnError }
	funcs := []func(context.Context) error{f, f, f}
	errs := []error{assert.AnError, assert.AnError, assert.AnError}
	target := fmt.Errorf("%v", errs)
	assert.Equal(t, unify(funcs)(context.Background()), target)
}

func mergeResource(t *testing.T, r1, r2 *resource.Resource) *resource.Resource {
	r, err := resource.Merge(r1, r2)
	assert.NoError(t, err)
	return r
}

func TestWithResource(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		envVarResourceAttributes: "key=value,rk5=7",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	cases := []struct {
		name    string
		options []Option
		want    *resource.Resource
		msg     string
	}{
		{
			name:    "explicitly empty resource",
			options: []Option{WithResource(resource.Empty())},
			want:    resource.Environment(),
		},
		{
			name:    "uses default if no resource option",
			options: []Option{},
			want:    resource.Default(),
		},
		{
			name:    "explicit resource",
			options: []Option{WithResource(resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk2", 5)))},
			want:    mergeResource(t, resource.Environment(), resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk2", 5))),
		},
		{
			name: "last resource wins",
			options: []Option{
				WithResource(resource.NewSchemaless(attribute.String("rk1", "vk1"), attribute.Int64("rk2", 5))),
				WithResource(resource.NewSchemaless(attribute.String("rk3", "rv3"), attribute.Int64("rk4", 10))),
			},
			want: mergeResource(t, resource.Environment(), resource.NewSchemaless(attribute.String("rk3", "rv3"), attribute.Int64("rk4", 10))),
		},
		{
			name:    "overlapping attributes with environment resource",
			options: []Option{WithResource(resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk5", 10)))},
			want:    mergeResource(t, resource.Environment(), resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk5", 10))),
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := newConfig(tc.options).res
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("WithResource:\n  -got +want %s", diff)
			}
		})
	}
}

func TestWithReader(t *testing.T) {
	r := &reader{}
	c := newConfig([]Option{WithReader(r)})
	require.Len(t, c.readers, 1)
	assert.Same(t, r, c.readers[0])
}

func TestWithView(t *testing.T) {
	c := newConfig([]Option{WithView(
		NewView(
			Instrument{Kind: InstrumentKindObservableCounter},
			Stream{Name: "a"},
		),
		NewView(
			Instrument{Kind: InstrumentKindCounter},
			Stream{Name: "b"},
		),
	)})
	assert.Len(t, c.views, 2)
}
