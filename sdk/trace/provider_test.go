// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
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
	injectExportError   error
}

func (t *basicSpanProcessor) Shutdown(context.Context) error {
	t.closed = true
	return t.injectShutdownError
}

func (*basicSpanProcessor) OnStart(context.Context, ReadWriteSpan) {}
func (*basicSpanProcessor) OnEnd(ReadOnlySpan)                     {}
func (t *basicSpanProcessor) ForceFlush(context.Context) error {
	t.flushed = true
	return t.injectExportError
}

type shutdownSpanProcessor struct {
	shutdown func(context.Context) error
}

func (t *shutdownSpanProcessor) Shutdown(ctx context.Context) error {
	return t.shutdown(ctx)
}

func (*shutdownSpanProcessor) OnStart(context.Context, ReadWriteSpan) {}
func (*shutdownSpanProcessor) OnEnd(ReadOnlySpan)                     {}
func (*shutdownSpanProcessor) ForceFlush(context.Context) error {
	return nil
}

type shutdownSampler struct {
	shutdown     func(context.Context) error
	shouldSample func(SamplingParameters) SamplingResult
}

func (s *shutdownSampler) ShouldSample(p SamplingParameters) SamplingResult {
	if s.shouldSample != nil {
		return s.shouldSample(p)
	}
	return SamplingResult{Decision: Drop}
}

func (*shutdownSampler) Description() string { return "shutdown sampler" }

func (s *shutdownSampler) Shutdown(ctx context.Context) error { return s.shutdown(ctx) }

type shutdownIDGenerator struct {
	randomIDGenerator
	shutdown func(context.Context) error
	newIDs   func(context.Context) (trace.TraceID, trace.SpanID)
}

func (g *shutdownIDGenerator) Shutdown(ctx context.Context) error { return g.shutdown(ctx) }

func (g *shutdownIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	if g.newIDs != nil {
		return g.newIDs(ctx)
	}
	return g.randomIDGenerator.NewIDs(ctx)
}

type forwardingSampler struct{ Sampler }

func (s forwardingSampler) Shutdown(ctx context.Context) error {
	return s.Sampler.(shutdowner).Shutdown(ctx)
}

const sensitiveExporterEndpoint = "user:pass@collector.internal:4318"

type marshalingSpanExporter struct{}

func (*marshalingSpanExporter) ExportSpans(context.Context, []ReadOnlySpan) error {
	return nil
}

func (*marshalingSpanExporter) Shutdown(context.Context) error {
	return nil
}

func (*marshalingSpanExporter) MarshalLog() any {
	return struct{ Endpoint string }{Endpoint: sensitiveExporterEndpoint}
}

func TestTracerProviderCreatedLogDoesNotIncludeExporterConfig(t *testing.T) {
	tests := []struct {
		name string
		opt  func(SpanExporter) TracerProviderOption
	}{
		{
			name: "batch",
			opt:  func(e SpanExporter) TracerProviderOption { return WithBatcher(e) },
		},
		{
			name: "simple",
			opt:  WithSyncer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orig := global.GetLogger()
			global.SetLogger(funcr.New(func(_, args string) {
				_, _ = buf.WriteString(args)
			}, funcr.Options{Verbosity: 4}))
			t.Cleanup(func() { global.SetLogger(orig) })

			tp := NewTracerProvider(tt.opt(&marshalingSpanExporter{}))
			require.NoError(t, tp.Shutdown(t.Context()))

			logged := buf.String()
			assert.Contains(t, logged, "TracerProvider created")
			assert.NotContains(t, logged, sensitiveExporterEndpoint)
		})
	}
}

func TestShutdownCallsTracerMethod(t *testing.T) {
	stp := NewTracerProvider()
	sp := &shutdownSpanProcessor{
		shutdown: func(context.Context) error {
			_ = stp.Tracer("abc") // must not deadlock
			return nil
		},
	}
	stp.RegisterSpanProcessor(sp)
	assert.NoError(t, stp.Shutdown(t.Context()))
	assert.True(t, stp.isShutdown.Load())
}

func TestForceFlushAndShutdownTraceProviderWithoutProcessor(t *testing.T) {
	stp := NewTracerProvider()
	assert.NoError(t, stp.ForceFlush(t.Context()))
	assert.NoError(t, stp.Shutdown(t.Context()))
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

func TestUnregisterUnknownSpanProcessor(t *testing.T) {
	stp := NewTracerProvider()
	sp1 := &basicSpanProcessor{}
	sp2 := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp1)
	stp.RegisterSpanProcessor(sp2)

	stp.UnregisterSpanProcessor(&basicSpanProcessor{})

	sps := stp.getSpanProcessors()
	require.Len(t, sps, 2)
	assert.Same(t, sp1, sps[0].sp)
	assert.Same(t, sp2, sps[1].sp)
}

func TestShutdownTraceProvider(t *testing.T) {
	stp := NewTracerProvider()
	sp := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp)

	assert.NoError(t, stp.ForceFlush(t.Context()))
	assert.True(t, sp.flushed, "error ForceFlush basicSpanProcessor")
	assert.NoError(t, stp.Shutdown(t.Context()))
	assert.True(t, stp.isShutdown.Load())
	assert.True(t, sp.closed, "error Shutdown basicSpanProcessor")
}

func TestShutdownComponentsAfterProcessorsOnce(t *testing.T) {
	var order []string
	sampler := &shutdownSampler{shutdown: func(context.Context) error {
		order = append(order, "sampler")
		return nil
	}}
	idGenerator := &shutdownIDGenerator{shutdown: func(context.Context) error {
		order = append(order, "ID generator")
		return nil
	}}
	processor := &shutdownSpanProcessor{shutdown: func(context.Context) error {
		order = append(order, "processor")
		return nil
	}}
	tp := NewTracerProvider(WithSampler(sampler), WithIDGenerator(idGenerator), WithSpanProcessor(processor))

	require.NoError(t, tp.Shutdown(t.Context()))
	assert.Equal(t, []string{"processor", "sampler", "ID generator"}, order)
	require.NoError(t, tp.Shutdown(t.Context()))
	assert.Equal(t, []string{"processor", "sampler", "ID generator"}, order)
}

func TestCachedTracerAfterShutdownDoesNotUseClosedComponents(t *testing.T) {
	var samplerUsed, idGeneratorUsed bool
	sampler := &shutdownSampler{
		shutdown: func(context.Context) error { return nil },
		shouldSample: func(SamplingParameters) SamplingResult {
			samplerUsed = true
			return SamplingResult{Decision: Drop}
		},
	}
	idGenerator := &shutdownIDGenerator{
		shutdown: func(context.Context) error { return nil },
		newIDs: func(context.Context) (trace.TraceID, trace.SpanID) {
			idGeneratorUsed = true
			return trace.TraceID{}, trace.SpanID{}
		},
	}
	tp := NewTracerProvider(WithSampler(sampler), WithIDGenerator(idGenerator))
	tracer := tp.Tracer("cached")
	require.NoError(t, tp.Shutdown(t.Context()))

	_, span := tracer.Start(t.Context(), "after shutdown")
	assert.False(t, span.IsRecording())
	assert.False(t, samplerUsed)
	assert.False(t, idGeneratorUsed)
}

func TestCompositeSamplerRepeatedDelegateShutdownPerSlot(t *testing.T) {
	var shutdownCount int
	delegate := &shutdownSampler{shutdown: func(context.Context) error {
		shutdownCount++
		return nil
	}}
	tp := NewTracerProvider(WithSampler(AlwaysRecord(ParentBased(
		delegate,
		WithRemoteParentSampled(delegate),
		WithLocalParentNotSampled(delegate),
	))))

	require.NoError(t, tp.Shutdown(t.Context()))
	assert.Equal(t, 3, shutdownCount)
}

func TestCompositeSamplerShutdownPropagatesThroughCustomWrapper(t *testing.T) {
	var shutdownCount int
	delegate := &shutdownSampler{shutdown: func(context.Context) error {
		shutdownCount++
		return nil
	}}
	composite := ParentBased(delegate)
	_, ok := composite.(shutdowner)
	require.True(t, ok)
	tp := NewTracerProvider(WithSampler(forwardingSampler{Sampler: composite}))

	require.NoError(t, tp.Shutdown(t.Context()))
	assert.Equal(t, 1, shutdownCount)
}

func TestCompositeSamplerShutdownErrorAndCancellation(t *testing.T) {
	t.Run("error does not stop other delegates", func(t *testing.T) {
		shutdownErr := errors.New("sampler shutdown failed")
		root := &shutdownSampler{shutdown: func(context.Context) error { return shutdownErr }}
		var nextCalled bool
		next := &shutdownSampler{shutdown: func(context.Context) error {
			nextCalled = true
			return nil
		}}
		tp := NewTracerProvider(WithSampler(ParentBased(root, WithRemoteParentSampled(next))))

		assert.ErrorIs(t, tp.Shutdown(t.Context()), shutdownErr)
		assert.True(t, nextCalled)
	})

	t.Run("cancellation skips later delegates", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		root := &shutdownSampler{shutdown: func(context.Context) error {
			cancel()
			return nil
		}}
		var nextCalled bool
		next := &shutdownSampler{shutdown: func(context.Context) error {
			nextCalled = true
			return nil
		}}
		tp := NewTracerProvider(WithSampler(ParentBased(root, WithRemoteParentSampled(next))))

		assert.ErrorIs(t, tp.Shutdown(ctx), context.Canceled)
		assert.False(t, nextCalled)
	})
}

func TestDefaultSamplerDecoratorsDoNotRequireShutdown(t *testing.T) {
	_, parentBasedShutdown := ParentBased(AlwaysSample()).(shutdowner)
	_, alwaysRecordShutdown := AlwaysRecord(AlwaysSample()).(shutdowner)
	assert.False(t, parentBasedShutdown)
	assert.False(t, alwaysRecordShutdown)
}

func TestShutdownComponentsJoinErrors(t *testing.T) {
	processorErr := errors.New("processor shutdown failed")
	samplerErr := errors.New("sampler shutdown failed")
	idGeneratorErr := errors.New("ID generator shutdown failed")
	sampler := &shutdownSampler{shutdown: func(context.Context) error { return samplerErr }}
	idGenerator := &shutdownIDGenerator{shutdown: func(context.Context) error { return idGeneratorErr }}
	processor := &shutdownSpanProcessor{shutdown: func(context.Context) error { return processorErr }}
	tp := NewTracerProvider(WithSampler(sampler), WithIDGenerator(idGenerator), WithSpanProcessor(processor))

	err := tp.Shutdown(t.Context())
	assert.ErrorIs(t, err, processorErr)
	assert.ErrorIs(t, err, samplerErr)
	assert.ErrorIs(t, err, idGeneratorErr)
}

func TestShutdownComponentsRespectCancellation(t *testing.T) {
	t.Run("default components", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		tp := NewTracerProvider()

		assert.NoError(t, tp.Shutdown(ctx))
	})

	t.Run("before shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		var samplerCalled, idGeneratorCalled bool
		sampler := &shutdownSampler{shutdown: func(context.Context) error {
			samplerCalled = true
			return nil
		}}
		idGenerator := &shutdownIDGenerator{shutdown: func(context.Context) error {
			idGeneratorCalled = true
			return nil
		}}
		tp := NewTracerProvider(WithSampler(sampler), WithIDGenerator(idGenerator))

		assert.ErrorIs(t, tp.Shutdown(ctx), context.Canceled)
		assert.False(t, samplerCalled)
		assert.False(t, idGeneratorCalled)
	})

	t.Run("during sampler shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		var idGeneratorCalled bool
		sampler := &shutdownSampler{shutdown: func(context.Context) error {
			cancel()
			return nil
		}}
		idGenerator := &shutdownIDGenerator{shutdown: func(context.Context) error {
			idGeneratorCalled = true
			return nil
		}}
		tp := NewTracerProvider(WithSampler(sampler), WithIDGenerator(idGenerator))

		assert.ErrorIs(t, tp.Shutdown(ctx), context.Canceled)
		assert.False(t, idGeneratorCalled)
	})

	t.Run("during ID generator shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		idGenerator := &shutdownIDGenerator{shutdown: func(context.Context) error {
			cancel()
			return nil
		}}
		tp := NewTracerProvider(WithIDGenerator(idGenerator))

		assert.ErrorIs(t, tp.Shutdown(ctx), context.Canceled)
	})
}

func TestFailedProcessorShutdown(t *testing.T) {
	stp := NewTracerProvider()
	spErr := errors.New("basic span processor shutdown failure")
	sp := &basicSpanProcessor{
		injectShutdownError: spErr,
	}
	stp.RegisterSpanProcessor(sp)

	err := stp.Shutdown(t.Context())
	assert.Error(t, err)
	assert.ErrorIs(t, err, spErr)
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

	err := stp.Shutdown(t.Context())
	assert.Error(t, err)
	assert.ErrorIs(t, err, spErr1)
	assert.ErrorIs(t, err, spErr2)
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

	err := stp.Shutdown(t.Context())
	assert.NoError(t, err)
	assert.True(t, stp.isShutdown.Load())
}

func TestSchemaURL(t *testing.T) {
	stp := NewTracerProvider()
	schemaURL := "https://opentelemetry.io/schemas/1.21.0"
	tracerIface := stp.Tracer("tracername", trace.WithSchemaURL(schemaURL))

	// Verify that the SchemaURL of the constructed Tracer is correctly populated.
	tracerStruct := tracerIface.(*tracer)
	assert.Equal(t, schemaURL, tracerStruct.instrumentationScope.SchemaURL)
}

func TestRegisterAfterShutdownWithoutProcessors(t *testing.T) {
	stp := NewTracerProvider()
	err := stp.Shutdown(t.Context())
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
	err := stp.Shutdown(t.Context())
	assert.NoError(t, err)
	assert.True(t, stp.isShutdown.Load())
	assert.Empty(t, stp.getSpanProcessors())

	sp2 := &basicSpanProcessor{}
	stp.RegisterSpanProcessor(sp2) // no-op
	assert.Empty(t, stp.getSpanProcessors())
}

func TestTracerProviderForceFlush(t *testing.T) {
	t.Run("AfterShutdown", func(t *testing.T) {
		stp := NewTracerProvider()
		sp1 := &basicSpanProcessor{}
		stp.RegisterSpanProcessor(sp1)
		ctx := t.Context()

		require.NoError(t, stp.ForceFlush(ctx))
		require.True(t, sp1.flushed, "SpanProcessor ForceFlush not called")

		sp1.flushed = false
		require.NoError(t, stp.Shutdown(ctx))

		require.NoError(t, stp.ForceFlush(ctx))
		assert.False(t, sp1.flushed, "SpanProcessor ForceFlush called after Shutdown")
	})

	t.Run("Multi", func(t *testing.T) {
		stp := NewTracerProvider()
		sp1 := &basicSpanProcessor{}
		sp2 := &basicSpanProcessor{}
		stp.RegisterSpanProcessor(sp1)
		stp.RegisterSpanProcessor(sp2)
		ctx := t.Context()

		require.NoError(t, stp.ForceFlush(ctx))
		require.True(t, sp1.flushed, "SpanProcessor ForceFlush not called")
		require.True(t, sp2.flushed, "SpanProcessor ForceFlush not called")
	})

	t.Run("Error", func(t *testing.T) {
		stp := NewTracerProvider()
		sp1 := &basicSpanProcessor{injectExportError: assert.AnError}
		sp2 := &basicSpanProcessor{}
		stp.RegisterSpanProcessor(sp1)
		stp.RegisterSpanProcessor(sp2)
		ctx := t.Context()

		assert.ErrorIs(t, stp.ForceFlush(ctx), assert.AnError, "span processor error not returned")
		require.True(t, sp1.flushed, "SpanProcessor ForceFlush not called")
		require.True(t, sp2.flushed, "SpanProcessor ForceFlush not called")
	})

	t.Run("WithCancel", func(t *testing.T) {
		stp := NewTracerProvider()
		sp1 := &basicSpanProcessor{}
		stp.RegisterSpanProcessor(sp1)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		assert.ErrorIs(t, stp.ForceFlush(ctx), context.Canceled)
	})
}

func TestTracerProviderSamplerConfigFromEnv(t *testing.T) {
	type testCase struct {
		sampler             string
		samplerArg          string
		argOptional         bool
		description         string
		errorType           error
		invalidArgErrorType any
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
						//nolint:usetesting // required to avoid getting a canceled context at cleanup.
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

func TestTracerProviderSamplerConfigFromEnvEmptyValues(t *testing.T) {
	tests := []struct {
		name          string
		sampler       string
		samplerArg    string
		setSamplerArg bool
		description   string
	}{
		{
			name:        "empty sampler",
			sampler:     "",
			description: ParentBased(AlwaysSample()).Description(),
		},
		{
			name:          "empty traceidratio sampler arg",
			sampler:       "traceidratio",
			samplerArg:    "",
			setSamplerArg: true,
			description:   TraceIDRatioBased(1.0).Description(),
		},
		{
			name:          "empty parentbased traceidratio sampler arg",
			sampler:       "parentbased_traceidratio",
			samplerArg:    "",
			setSamplerArg: true,
			description:   ParentBased(TraceIDRatioBased(1.0)).Description(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler.Reset()
			t.Cleanup(handler.Reset)

			t.Setenv(envTracesSampler, test.sampler)
			if test.setSamplerArg {
				t.Setenv(envTracesSamplerArg, test.samplerArg)
			}

			stp := NewTracerProvider(WithSyncer(NewTestExporter()))
			t.Cleanup(func() {
				//nolint:usetesting // required to avoid getting a canceled context at cleanup.
				require.NoError(t, stp.Shutdown(context.Background()))
			})

			assert.Equal(t, test.description, stp.sampler.Description())
			assert.Empty(t, handler.errs)
		})
	}
}

func testStoredError(t *testing.T, target any) {
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

func TestTracerProviderObservability(t *testing.T) {
	handler.Reset()
	p := NewTracerProvider()

	// Enable observability
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	tr := p.Tracer("test-tracer")
	require.IsType(t, &tracer{}, tr)

	tStruct := tr.(*tracer)
	assert.True(t, tStruct.inst.Enabled(), "observability should be enabled")

	// Verify errors are passed to the otel handler
	handlerErrs := handler.errs
	assert.Empty(t, handlerErrs, "No errors should occur during instrument creation")
}

func TestTracerProviderObservabilityErrorsHandled(t *testing.T) {
	handler.Reset()

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })
	otel.SetMeterProvider(&errMeterProvider{err: assert.AnError})

	p := NewTracerProvider()

	// Enable observability
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Create a tracer to trigger instrument creation.
	tr := p.Tracer("test-tracer")
	_ = tr

	require.Len(t, handler.errs, 1)
	assert.ErrorIs(t, handler.errs[0], assert.AnError)
}

type errMeterProvider struct {
	metric.MeterProvider

	err error
}

func (mp *errMeterProvider) Meter(string, ...metric.MeterOption) metric.Meter {
	return &errMeter{err: mp.err}
}

type errMeter struct {
	metric.Meter

	err error
}

func (m *errMeter) Int64Counter(string, ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return nil, m.err
}

func (m *errMeter) Int64UpDownCounter(string, ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return nil, m.err
}

type testExperimentalOption struct {
	TracerProviderOption
}

func (testExperimentalOption) Experimental() {}

func TestExperimentalOptionSafe(t *testing.T) {
	var opt testExperimentalOption

	assert.NotPanics(t, func() { _ = NewTracerProvider(opt) })
}
