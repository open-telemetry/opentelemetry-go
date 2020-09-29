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

package propagators_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/propagators"
)

const (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanIDStr  = "00f067aa0ba902b7"
)

var (
	traceID = mustTraceIDFromHex(traceIDStr)
	spanID  = mustSpanIDFromHex(spanIDStr)
)

func mustTraceIDFromHex(s string) (t trace.ID) {
	var err error
	t, err = trace.IDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

type outOfThinAirPropagator struct {
	t *testing.T
}

var _ otel.HTTPPropagator = outOfThinAirPropagator{}

func (p outOfThinAirPropagator) Extract(ctx context.Context, supplier otel.HTTPSupplier) context.Context {
	sc := trace.SpanContext{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: 0,
	}
	require.True(p.t, sc.IsValid())
	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (outOfThinAirPropagator) Inject(context.Context, otel.HTTPSupplier) {}

func (outOfThinAirPropagator) GetAllKeys() []string {
	return nil
}

type nilSupplier struct{}

var _ otel.HTTPSupplier = nilSupplier{}

func (nilSupplier) Get(key string) string {
	return ""
}

func (nilSupplier) Set(key string, value string) {}

func TestMultiplePropagators(t *testing.T) {
	ootaProp := outOfThinAirPropagator{t: t}
	ns := nilSupplier{}
	testProps := []otel.HTTPPropagator{
		propagators.TraceContext{},
	}
	bg := context.Background()
	// sanity check of oota propagator, ensuring that it really
	// generates the valid span context out of thin air
	{
		props := otel.New(otel.WithExtractors(ootaProp))
		ctx := otel.ExtractHTTP(bg, props, ns)
		sc := trace.RemoteSpanContextFromContext(ctx)
		require.True(t, sc.IsValid(), "oota prop failed sanity check")
	}
	// sanity check for real propagators, ensuring that they
	// really are not putting any valid span context into an empty
	// go context in absence of the HTTP headers.
	for _, prop := range testProps {
		props := otel.New(otel.WithExtractors(prop))
		ctx := otel.ExtractHTTP(bg, props, ns)
		sc := trace.RemoteSpanContextFromContext(ctx)
		require.Falsef(t, sc.IsValid(), "%#v failed sanity check", prop)
	}
	for _, prop := range testProps {
		props := otel.New(otel.WithExtractors(ootaProp, prop))
		ctx := otel.ExtractHTTP(bg, props, ns)
		sc := trace.RemoteSpanContextFromContext(ctx)
		assert.Truef(t, sc.IsValid(), "%#v clobbers span context", prop)
	}
}

func TestDefaultHTTPPropagator(t *testing.T) {
	assert.IsType(t, propagators.TraceContext{}, propagators.DefaultHTTPPropagator())
}
