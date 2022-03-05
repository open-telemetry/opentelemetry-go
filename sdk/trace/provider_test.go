// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/trace"
)

type basicSpanProcesor struct {
	running             bool
	injectShutdownError error
}

func (t *basicSpanProcesor) Shutdown(context.Context) error {
	t.running = false
	return t.injectShutdownError
}

func (t *basicSpanProcesor) OnStart(context.Context, ReadWriteSpan) {}
func (t *basicSpanProcesor) OnEnd(ReadOnlySpan)                     {}
func (t *basicSpanProcesor) ForceFlush(context.Context) error {
	return nil
}

func TestShutdownTraceProvider(t *testing.T) {
	stp := NewTracerProvider()
	sp := &basicSpanProcesor{}
	stp.RegisterSpanProcessor(sp)

	sp.running = true

	_ = stp.Shutdown(context.Background())

	if sp.running != false {
		t.Errorf("Error shutdown basicSpanProcesor\n")
	}
}

func TestFailedProcessorShutdown(t *testing.T) {
	stp := NewTracerProvider()
	spErr := errors.New("basic span processor shutdown failure")
	sp := &basicSpanProcesor{
		running:             true,
		injectShutdownError: spErr,
	}
	stp.RegisterSpanProcessor(sp)

	err := stp.Shutdown(context.Background())
	assert.Error(t, err)
	assert.Equal(t, err, spErr)
}

func TestFailedProcessorShutdownInUnregister(t *testing.T) {
	handler.Reset()
	stp := NewTracerProvider()
	spErr := errors.New("basic span processor shutdown failure")
	sp := &basicSpanProcesor{
		running:             true,
		injectShutdownError: spErr,
	}
	stp.RegisterSpanProcessor(sp)
	stp.UnregisterSpanProcessor(sp)

	assert.Contains(t, handler.errs, spErr)

	err := stp.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestSchemaURL(t *testing.T) {
	stp := NewTracerProvider()
	schemaURL := "https://opentelemetry.io/schemas/1.2.0"
	tracerIface := stp.Tracer("tracername", trace.WithSchemaURL(schemaURL))

	// Verify that the SchemaURL of the constructed Tracer is correctly populated.
	tracerStruct := tracerIface.(*tracer)
	assert.EqualValues(t, schemaURL, tracerStruct.instrumentationLibrary.SchemaURL)
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
			envVars := map[string]string{
				"OTEL_TRACES_SAMPLER": test.sampler,
			}

			if test.samplerArg != "" {
				envVars["OTEL_TRACES_SAMPLER_ARG"] = test.samplerArg
			}
			envStore, err := ottest.SetEnvVariables(envVars)
			require.NoError(t, err)
			t.Cleanup(func() {
				handler.Reset()
				require.NoError(t, envStore.Restore())
			})

			stp := NewTracerProvider(WithSyncer(NewTestExporter()))
			assert.Equal(t, test.description, stp.sampler.Description())
			if test.errorType != nil {
				testStoredError(t, test.errorType)
			} else {
				assert.Empty(t, handler.errs)
			}

			if test.argOptional {
				t.Run("invalid sampler arg", func(t *testing.T) {
					envStore, err := ottest.SetEnvVariables(map[string]string{
						"OTEL_TRACES_SAMPLER":     test.sampler,
						"OTEL_TRACES_SAMPLER_ARG": "invalid-ignored-string",
					})
					require.NoError(t, err)
					t.Cleanup(func() {
						handler.Reset()
						require.NoError(t, envStore.Restore())
					})

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
		require.NotNil(t, target.(error))

		defer handler.Reset()
		if errors.Is(err, target.(error)) {
			return
		}

		assert.ErrorAs(t, err, target)
	}
}
