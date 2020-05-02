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

package testtrace_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

type outOfThinAirPropagator struct {
	t *testing.T
}

var _ propagation.HTTPPropagator = outOfThinAirPropagator{}

func (p outOfThinAirPropagator) Extract(ctx context.Context, supplier propagation.HTTPSupplier) context.Context {
	traceID, err := trace.IDFromHex("938753245abe987f098c0987a9873987")
	require.NoError(p.t, err)
	spanID, err := trace.SpanIDFromHex("2345f98c0987a09d")
	require.NoError(p.t, err)
	sc := trace.SpanContext{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: 0,
	}
	require.True(p.t, sc.IsValid())
	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (outOfThinAirPropagator) Inject(context.Context, propagation.HTTPSupplier) {}

func (outOfThinAirPropagator) GetAllKeys() []string {
	return nil
}

type nilSupplier struct{}

var _ propagation.HTTPSupplier = nilSupplier{}

func (nilSupplier) Get(key string) string {
	return ""
}

func (nilSupplier) Set(key string, value string) {}

func TestMultiplePropagators(t *testing.T) {
	ootaProp := outOfThinAirPropagator{t: t}
	ns := nilSupplier{}
	testProps := []propagation.HTTPPropagator{
		trace.TraceContext{},
		trace.B3{SingleHeader: false},
		trace.B3{SingleHeader: true},
	}
	bg := context.Background()
	// sanity check of oota propagator, ensuring that it really
	// generates the valid span context out of thin air
	{
		props := propagation.New(propagation.WithExtractors(ootaProp))
		ctx := propagation.ExtractHTTP(bg, props, ns)
		sc := trace.RemoteSpanContextFromContext(ctx)
		require.True(t, sc.IsValid(), "oota prop failed sanity check")
	}
	// sanity check for real propagators, ensuring that they
	// really are not putting any valid span context into an empty
	// go context in absence of the HTTP headers.
	for _, prop := range testProps {
		props := propagation.New(propagation.WithExtractors(prop))
		ctx := propagation.ExtractHTTP(bg, props, ns)
		sc := trace.RemoteSpanContextFromContext(ctx)
		require.Falsef(t, sc.IsValid(), "%#v failed sanity check", prop)
	}
	for _, prop := range testProps {
		props := propagation.New(propagation.WithExtractors(ootaProp, prop))
		ctx := propagation.ExtractHTTP(bg, props, ns)
		sc := trace.RemoteSpanContextFromContext(ctx)
		assert.Truef(t, sc.IsValid(), "%#v clobbers span context", prop)
	}
}
