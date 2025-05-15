// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package propagation_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestExtractValidTraceContextEnvCarrier(t *testing.T) {
	stateStr := "key1=value1,key2=value2"
	state, err := trace.ParseTraceState(stateStr)
	require.NoError(t, err)

	tests := []struct {
		name string
		envs map[string]string
		want trace.SpanContext
	}{
		{
			name: "sampled",
			envs: map[string]string{
				"TRACEPARENT": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			},
			want: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name: "valid tracestate",
			envs: map[string]string{
				"TRACEPARENT": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
				"TRACESTATE":  stateStr,
			},
			want: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceState: state,
				Remote:     true,
			}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}
			ctx = prop.Extract(ctx, propagation.EnvCarrier{})
			assert.Equal(t, tc.want, trace.SpanContextFromContext(ctx))
		})
	}
}

func TestInjectTraceContextEnvCarrier(t *testing.T) {
	stateStr := "key1=value1,key2=value2"
	state, err := trace.ParseTraceState(stateStr)
	require.NoError(t, err)

	tests := []struct {
		name string
		want map[string]string
		sc   trace.SpanContext
	}{
		{
			name: "sampled",
			want: map[string]string{
				"TRACEPARENT": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			},
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name: "with tracestate",
			want: map[string]string{
				"TRACEPARENT": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
				"TRACESTATE":  stateStr,
			},
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceState: state,
				Remote:     true,
			}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = trace.ContextWithRemoteSpanContext(ctx, tc.sc)
			c := propagation.EnvCarrier{
				SetEnvFunc: func(key, value string) error {
					t.Setenv(key, value)
					return nil
				},
			}

			prop.Inject(ctx, c)

			for k, v := range tc.want {
				if got := os.Getenv(k); got != v {
					t.Errorf("got %s=%s, want %s=%s", k, got, k, v)
				}

			}
		})
	}
}
