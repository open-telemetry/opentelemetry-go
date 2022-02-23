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
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/sdk/internal/env"

	"github.com/stretchr/testify/assert"

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

func TestNewTraceProviderWithoutSpanLimitConfiguration(t *testing.T) {
	envStore := ottest.NewEnvStore()
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	envStore.Record(env.SpanAttributeValueLengthKey)
	envStore.Record(env.SpanAttributeCountKey)
	envStore.Record(env.SpanEventCountKey)
	envStore.Record(env.SpanLinkCountKey)
	envStore.Record(env.SpanEventAttributeCountKey)
	envStore.Record(env.SpanLinkAttributeCountKey)
	require.NoError(t, os.Setenv(env.SpanAttributeValueLengthKey, "42"))
	require.NoError(t, os.Setenv(env.SpanEventCountKey, "111"))
	require.NoError(t, os.Setenv(env.SpanAttributeCountKey, "222"))
	require.NoError(t, os.Setenv(env.SpanLinkCountKey, "333"))
	require.NoError(t, os.Setenv(env.SpanEventAttributeCountKey, "42"))
	require.NoError(t, os.Setenv(env.SpanLinkAttributeCountKey, "42"))
	assert.Equal(t, NewTracerProvider().spanLimits, SpanLimits{
		AttributeValueLengthLimit:   42,
		AttributeCountLimit:         222,
		EventCountLimit:             111,
		LinkCountLimit:              333,
		AttributePerEventCountLimit: 42,
		AttributePerLinkCountLimit:  42,
	})
}

func TestNewTraceProviderWithSpanLimitConfigurationFromOptsAndEnvironmentVariable(t *testing.T) {
	envStore := ottest.NewEnvStore()
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	envStore.Record(env.SpanAttributeValueLengthKey)
	envStore.Record(env.SpanAttributeCountKey)
	envStore.Record(env.SpanEventCountKey)
	envStore.Record(env.SpanLinkCountKey)
	envStore.Record(env.SpanEventAttributeCountKey)
	envStore.Record(env.SpanLinkAttributeCountKey)
	require.NoError(t, os.Setenv(env.SpanAttributeValueLengthKey, "-2"))
	require.NoError(t, os.Setenv(env.SpanEventCountKey, "-2"))
	require.NoError(t, os.Setenv(env.SpanAttributeCountKey, "-2"))
	require.NoError(t, os.Setenv(env.SpanLinkCountKey, "-2"))
	require.NoError(t, os.Setenv(env.SpanEventAttributeCountKey, "-2"))
	require.NoError(t, os.Setenv(env.SpanLinkAttributeCountKey, "-2"))
	sl := SpanLimits{
		AttributeValueLengthLimit:   42,
		AttributeCountLimit:         222,
		EventCountLimit:             111,
		LinkCountLimit:              333,
		AttributePerEventCountLimit: 42,
		AttributePerLinkCountLimit:  42,
	}
	assert.Equal(t, NewTracerProvider(WithSpanLimits(sl)).spanLimits, sl)
}

func TestNewTraceProviderSpanLimitUnlimitedFromEnv(t *testing.T) {
	envStore := ottest.NewEnvStore()
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	envStore.Record(env.SpanAttributeValueLengthKey)
	envStore.Record(env.SpanAttributeCountKey)
	envStore.Record(env.SpanEventCountKey)
	envStore.Record(env.SpanLinkCountKey)
	envStore.Record(env.SpanEventAttributeCountKey)
	envStore.Record(env.SpanLinkAttributeCountKey)
	// OTel spec says this is invalid (negative) for
	// SpanLinkAttributeCountKey, but since we will revert to the default
	// (unlimited) which uses negative values to signal this than this value
	// is expected to pass through.
	require.NoError(t, os.Setenv(env.SpanAttributeValueLengthKey, "-1"))
	require.NoError(t, os.Setenv(env.SpanEventCountKey, "-1"))
	require.NoError(t, os.Setenv(env.SpanAttributeCountKey, "-1"))
	require.NoError(t, os.Setenv(env.SpanLinkCountKey, "-1"))
	require.NoError(t, os.Setenv(env.SpanEventAttributeCountKey, "-1"))
	require.NoError(t, os.Setenv(env.SpanLinkAttributeCountKey, "-1"))
	assert.Equal(t, NewTracerProvider().spanLimits, SpanLimits{
		AttributeValueLengthLimit:   -1,
		AttributeCountLimit:         -1,
		EventCountLimit:             -1,
		LinkCountLimit:              -1,
		AttributePerEventCountLimit: -1,
		AttributePerLinkCountLimit:  -1,
	})
}

func TestNewTraceProviderSpanLimitDefaults(t *testing.T) {
	assert.Equal(t, NewTracerProvider().spanLimits, NewSpanLimits())
}
