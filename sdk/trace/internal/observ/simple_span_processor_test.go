// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package observ_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/internal/observ"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

const sspComponentId = 0

func TestSSPComponentName(t *testing.T) {
	got := observ.SSPComponentName(10)
	want := semconv.OTelComponentName("simple_span_processor/10")
	assert.Equal(t, want, got)
}

func TestNewSSPError(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	_, err := observ.NewSSP(sspComponentId)
	require.ErrorIs(t, err, assert.AnError, "new instrument errors")
	assert.ErrorContains(t, err, "create SSP processed spans metric")
}

func TestNewSSPDisabled(t *testing.T) {
	ssp, err := observ.NewSSP(sspComponentId)
	assert.NoError(t, err)
	assert.Nil(t, ssp)
}

func TestSSPSpanProcessed(t *testing.T) {
	ctx := t.Context()
	collect := setup(t)
	ssp, err := observ.NewSSP(sspComponentId)
	assert.NoError(t, err)

	ssp.SpanProcessed(ctx, 1, nil)
	check(t, collect(), processed(dPt(sspSet(), 1)))
	ssp.SpanProcessed(ctx, 2, nil)
	check(t, collect(), processed(dPt(sspSet(), 3)))

	processErr := errors.New("error processing span")
	ssp.SpanProcessed(ctx, 1, processErr)
	check(t, collect(), processed(
		dPt(sspSet(), 3),
		dPt(sspSet(semconv.ErrorType(processErr)), 1),
	))
}

func BenchmarkSSP(b *testing.B) {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
}

func sspSet(attrs ...attribute.KeyValue) attribute.Set {
	return attribute.NewSet(append([]attribute.KeyValue{
		semconv.OTelComponentTypeSimpleSpanProcessor,
		observ.SSPComponentName(sspComponentId),
	}, attrs...)...)
}
