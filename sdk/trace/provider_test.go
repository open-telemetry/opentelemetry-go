// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	envTracesSampler    = "OTEL_TRACES_SAMPLER"
	envTracesSamplerArg = "OTEL_TRACES_SAMPLER_ARG"
)

type basicSpanProcessor struct {
	flushed             bool
	closed              bool
	injectShutdownError error
}

func (t *basicSpanProcessor) Shutdown(context.Context) error {
	t.closed = true
	return t.injectShutdownError
}

func (t *basicSpanProcessor) OnStart(context.Context, ReadWriteSpan) {}
func (t *basicSpanProcessor) OnEnd(ReadOnlySpan)                     {}
func (t *basicSpanProcessor) ForceFlush(context.Context) error {
	t.flushed = true
	return nil
}

type shutdownSpanProcessor struct {
	shutdown func(context.Context) error
}

func (t *shutdownSpanProcessor) Shutdown(ctx context.Context) error {
	return t.shutdown(ctx)
}

func (t *shutdownSpanProcessor) OnStart(context.Context, ReadWriteSpan) {}
func (t *shutdownSpanProcessor) OnEnd(ReadOnlySpan)                     {}
func (t *shutdownSpanProcessor) ForceFlush(context.Context) error {
	return nil
}

func TestShutdownCallsTracerMethod(t *testing.T) {
	stp := NewTracerProvider()
	sp := &shutdownSpanProcessor{
		shutdown: func(ctx context.Context) error {
			_ = stp.Tracer("abc") // must not deadlock
			return nil
		},
	}
	stp.RegisterSpanProcessor(sp)
	assert.NoError(t, stp.Shutdown(context.Background()))
	assert.True(t, stp.isShutdown.Load())
}

func TestForceFlushAndShutdownTraceProviderWithoutProcessor(t *testing.T) {
	stp := NewTracerProvider()
	assert.NoError(t, stp.ForceFlush(context.Background()))
	assert.NoError(t, stp.Shutdown(context.Background()))
	assert.True(t, stp.isShutdown.Load())
}

func TestUnregisterFirst(t *testing.T) {
	stp := NewTracerProvider()
	sp1 := &basicSpanProcessor{}
	sp2 := &basicSpanProcessor{}
	sp3 := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp1)
	stp.RegisterSpanProcessor(sp2)
	stp.RegisterSpanProcessor(sp3)

	stp.UnregisterSpanProcessor(sp1)

	sps := stp.getSpanProcessors()
	require.Len(t, sps, 2)
	assert.Same(t, sp2, sps[0].sp)
	assert.Same(t, sp3, sps[1].sp)
}

func TestUnregisterMiddle(t *testing.T) {
	stp := NewTracerProvider()
	sp1 := &basicSpanProcessor{}
	sp2 := &basicSpanProcessor{}
	sp3 := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp1)
	stp.RegisterSpanProcessor(sp2)
	stp.RegisterSpanProcessor(sp3)

	stp.UnregisterSpanProcessor(sp2)

	sps := stp.getSpanProcessors()
	require.Len(t, sps, 2)
	assert.Same(t, sp1, sps[0].sp)
	assert.Same(t, sp3, sps[1].sp)
}

func TestUnregisterLast(t *testing.T) {
	stp := NewTracerProvider()
	sp1 := &basicSpanProcessor{}
	sp2 := &basicSpanProcessor{}
	sp3 := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp1)
	stp.RegisterSpanProcessor(sp2)
	stp.RegisterSpanProcessor(sp3)

	stp.UnregisterSpanProcessor(sp3)

	sps := stp.getSpanProcessors()
	require.Len(t, sps, 2)
	assert.Same(t, sp1, sps[0].sp)
	assert.Same(t, sp2, sps[1].sp)
}

func TestShutdownTraceProvider(t *testing.T) {
	stp := NewTracerProvider()
	sp := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp)

	assert.NoError(t, stp.ForceFlush(context.Background()))
	assert.True(t, sp.flushed, "error ForceFlush basicSpanProcessor")
	assert.NoError(t, stp.Shutdown(context.Background()))
	assert.True(t, stp.isShutdown.Load())
	assert.True(t, sp.closed, "error Shutdown basicSpanProcessor")
}

func TestFailedProcessorShutdown(t *testing.T) {
	stp := NewTracerProvider()
	spErr := errors.New("basic span processor shutdown failure")
	sp := &basicSpanProcessor{
		injectShutdownError: spErr,
	}
	stp.RegisterSpanProcessor(sp)

	err := stp.Shutdown(context.Background())
	assert.Error(t, err)
	assert.Equal(t, err, spErr)
	assert.True(t, stp.isShutdown.Load())
}

func TestFailedProcessorsShutdown(t *testing.T) {
	stp := NewTracerProvider()
	spErr1 := errors.New("basic span processor shutdown failure1")
	spErr2 := errors.New("basic span processor shutdown failure2")
	sp1 := &basicSpanProcessor{
		injectShutdownError: spErr1,
	}
	sp2 := &basicSpanProcessor{
		injectShutdownError: spErr2,
	}
	stp.RegisterSpanProcessor(sp1)
	stp.RegisterSpanProcessor(sp2)

	err := stp.Shutdown(context.Background())
	assert.Error(t, err)
	assert.EqualError(t, err, "basic span processor shutdown failure1; basic span processor shutdown failure2")
	assert.True(t, sp1.closed)
	assert.True(t, sp2.closed)
	assert.True(t, stp.isShutdown.Load())
}

func TestFailedProcessorShutdownInUnregister(t *testing.T) {
	handler.Reset()
	stp := NewTracerProvider()
	spErr := errors.New("basic span processor shutdown failure")
	sp := &basicSpanProcessor{
		injectShutdownError: spErr,
	}
	stp.RegisterSpanProcessor(sp)
	stp.UnregisterSpanProcessor(sp)

	assert.Contains(t, handler.errs, spErr)

	err := stp.Shutdown(context.Background())
	assert.NoError(t, err)
	assert.True(t, stp.isShutdown.Load())
}

func TestSchemaURL(t *testing.T) {
	stp := NewTracerProvider()
	schemaURL := "https://opentelemetry.io/schemas/1.2.0"
	tracerIface := stp.Tracer("tracername", trace.WithSchemaURL(schemaURL))

	// Verify that the SchemaURL of the constructed Tracer is correctly populated.
	tracerStruct := tracerIface.(*tracer)
	assert.Equal(t, schemaURL, tracerStruct.instrumentationScope.SchemaURL)
}

func TestRegisterAfterShutdownWithoutProcessors(t *testing.T) {
	stp := NewTracerProvider()
	err := stp.Shutdown(context.Background())
	assert.NoError(t, err)
	assert.True(t, stp.isShutdown.Load())

	sp := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp) // no-op
	assert.Empty(t, stp.getSpanProcessors())
}

func TestRegisterAfterShutdownWithProcessors(t *testing.T) {
	stp := NewTracerProvider()
	sp1 := &basicSpanProcessor{}

	stp.RegisterSpanProcessor(sp1)
	err := stp.Shutdown(context.Background())
	assert.NoError(t, err)
	assert.True(t, stp.isShutdown.Load())
	assert.Empty(t, stp.getSpanProcessors())

	sp2 := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp2) // no-op
	assert.Empty(t, stp.getSpanProcessors())
}

func TestTracerProviderSamplerConfigFromEnv(t *testing.T) {
	type testCase struct {
		sampler             string
		samplerArg          string
		argOptional         bool
		description         string
		errorType           error
		invalidArgErrorType interface{}
	}

	randFloat := rand.Float64()

	tests := []testCase{
		{
			sampler:             "invalid-sampler",
			argOptional:         true,
			description:         ParentBased(AlwaysSample()).Description(),
			errorType:           errUnsupportedSampler("invalid-sampler"),
			invalidArgErrorType: func() *errUnsupportedSampler { e := errUnsupportedSampler("invalid-sampler"); return &e }(),
		},
		{
			sampler:     "always_on",
			argOptional: true,
			description: AlwaysSample().Description(),
		},
		{
			sampler:     "always_off",
			argOptional: true,
			description: NeverSample().Description(),
		},
		{
			sampler:     "traceidratio",
			samplerArg:  fmt.Sprintf("%g", randFloat),
			description: TraceIDRatioBased(randFloat).Description(),
		},
		{
			sampler:     "traceidratio",
			samplerArg:  fmt.Sprintf("%g", -randFloat),
			description: TraceIDRatioBased(1.0).Description(),
			errorType:   errNegativeTraceIDRatio,
		},
		{
			sampler:     "traceidratio",
			samplerArg:  fmt.Sprintf("%g", 1+randFloat),
			description: TraceIDRatioBased(1.0).Description(),
			errorType:   errGreaterThanOneTraceIDRatio,
		},
		{
			sampler:             "traceidratio",
			argOptional:         true,
			description:         TraceIDRatioBased(1.0).Description(),
			invalidArgErrorType: new(samplerArgParseError),
		},
		{
			sampler:     "parentbased_always_on",
			argOptional: true,
			description: ParentBased(AlwaysSample()).Description(),
		},
		{
			sampler:     "parentbased_always_off",
			argOptional: true,
			description: ParentBased(NeverSample()).Description(),
		},
		{
			sampler:     "parentbased_traceidratio",
			samplerArg:  fmt.Sprintf("%g", randFloat),
			description: ParentBased(TraceIDRatioBased(randFloat)).Description(),
		},
		{
			sampler:     "parentbased_traceidratio",
			samplerArg:  fmt.Sprintf("%g", -randFloat),
			description: ParentBased(TraceIDRatioBased(1.0)).Description(),
			errorType:   errNegativeTraceIDRatio,
		},
		{
			sampler:     "parentbased_traceidratio",
			samplerArg:  fmt.Sprintf("%g", 1+randFloat),
			description: ParentBased(TraceIDRatioBased(1.0)).Description(),
			errorType:   errGreaterThanOneTraceIDRatio,
		},
		{
			sampler:             "parentbased_traceidratio",
			argOptional:         true,
			description:         ParentBased(TraceIDRatioBased(1.0)).Description(),
			invalidArgErrorType: new(samplerArgParseError),
		},
	}

	handler.Reset()

	for _, test := range tests {
		t.Run(test.sampler, func(t *testing.T) {
			t.Setenv(envTracesSampler, test.sampler)

			if test.samplerArg != "" {
				t.Setenv(envTracesSamplerArg, test.samplerArg)
			}

			stp := NewTracerProvider(WithSyncer(NewTestExporter()))
			assert.Equal(t, test.description, stp.sampler.Description())
			if test.errorType != nil {
				testStoredError(t, test.errorType)
			} else {
				assert.Empty(t, handler.errs)
			}

			if test.argOptional {
				t.Run("invalid sampler arg", func(t *testing.T) {
					t.Setenv(envTracesSampler, test.sampler)
					t.Setenv(envTracesSamplerArg, "invalid-ignored-string")

					stp := NewTracerProvider(WithSyncer(NewTestExporter()))
					t.Cleanup(func() {
						require.NoError(t, stp.Shutdown(context.Background()))
					})
					assert.Equal(t, test.description, stp.sampler.Description())

					if test.invalidArgErrorType != nil {
						testStoredError(t, test.invalidArgErrorType)
					} else {
						assert.Empty(t, handler.errs)
					}
				})
			}
		})
	}
}

func testStoredError(t *testing.T, target interface{}) {
	t.Helper()

	if assert.Len(t, handler.errs, 1) && assert.Error(t, handler.errs[0]) {
		err := handler.errs[0]

		require.Implements(t, (*error)(nil), target)
		require.Error(t, target.(error))

		defer handler.Reset()
		if errors.Is(err, target.(error)) {
			return
		}

		assert.ErrorAs(t, err, target)
	}
}

func TestTracerProviderReturnsSameTracer(t *testing.T) {
	p := NewTracerProvider()

	t0, t1, t2 := p.Tracer(
		"t0",
	), p.Tracer(
		"t1",
	), p.Tracer(
		"t0",
		trace.WithInstrumentationAttributes(attribute.String("foo", "bar")),
	)
	assert.NotSame(t, t0, t1)
	assert.NotSame(t, t0, t2)
	assert.NotSame(t, t1, t2)

	t3, t4, t5 := p.Tracer(
		"t0",
	), p.Tracer(
		"t1",
	), p.Tracer(
		"t0",
		trace.WithInstrumentationAttributes(attribute.String("foo", "bar")),
	)
	assert.Same(t, t0, t3)
	assert.Same(t, t1, t4)
	assert.Same(t, t2, t5)
}
