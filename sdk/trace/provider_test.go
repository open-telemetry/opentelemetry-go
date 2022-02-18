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
	envStore.Record(env.SpanAttributesCountKey)
	envStore.Record(env.SpanEventCountKey)
	envStore.Record(env.SpanLinkCountKey)
	require.NoError(t, os.Setenv(env.SpanEventCountKey, "111"))
	require.NoError(t, os.Setenv(env.SpanAttributesCountKey, "222"))
	require.NoError(t, os.Setenv(env.SpanLinkCountKey, "333"))
	tp := NewTracerProvider()
	assert.Equal(t, 111, tp.spanLimits.EventCountLimit)
	assert.Equal(t, 222, tp.spanLimits.AttributeCountLimit)
	assert.Equal(t, 333, tp.spanLimits.LinkCountLimit)
}

func TestNewTraceProviderWithSpanLimitConfigurationFromOptsAndEnvironmentVariable(t *testing.T) {
	envStore := ottest.NewEnvStore()
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	envStore.Record(env.SpanAttributesCountKey)
	envStore.Record(env.SpanEventCountKey)
	envStore.Record(env.SpanLinkCountKey)
	require.NoError(t, os.Setenv(env.SpanEventCountKey, "111"))
	require.NoError(t, os.Setenv(env.SpanAttributesCountKey, "222"))
	require.NoError(t, os.Setenv(env.SpanLinkCountKey, "333"))
	tp := NewTracerProvider(WithSpanLimits(SpanLimits{
		EventCountLimit:     1,
		AttributeCountLimit: 2,
		LinkCountLimit:      3,
	}))
	assert.Equal(t, 1, tp.spanLimits.EventCountLimit)
	assert.Equal(t, 2, tp.spanLimits.AttributeCountLimit)
	assert.Equal(t, 3, tp.spanLimits.LinkCountLimit)
}

func TestNewTraceProviderWithInvalidSpanLimitConfigurationFromEnvironmentVariable(t *testing.T) {
	envStore := ottest.NewEnvStore()
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	envStore.Record(env.SpanAttributesCountKey)
	envStore.Record(env.SpanEventCountKey)
	envStore.Record(env.SpanLinkCountKey)
	require.NoError(t, os.Setenv(env.SpanEventCountKey, "-111"))
	require.NoError(t, os.Setenv(env.SpanAttributesCountKey, "-222"))
	require.NoError(t, os.Setenv(env.SpanLinkCountKey, "-333"))
	tp := NewTracerProvider()
	assert.Equal(t, 128, tp.spanLimits.EventCountLimit)
	assert.Equal(t, 128, tp.spanLimits.AttributeCountLimit)
	assert.Equal(t, 128, tp.spanLimits.LinkCountLimit)
}
