// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewConfigEnvEndpointGRPCTargets asserts parity with the sibling
// exporters otlptracegrpc (#8852) and otlpmetricgrpc (#8862): non-HTTP(S)
// gRPC targets from OTEL_EXPORTER_OTLP_ENDPOINT must be preserved.
func TestNewConfigEnvEndpointGRPCTargets(t *testing.T) {
	testcases := []struct {
		name         string
		env          string
		wantEndpoint string
		wantInsecure bool
	}{
		{
			name:         "unix scheme preserved",
			env:          "unix:///tmp/grpc.sock",
			wantEndpoint: "unix:///tmp/grpc.sock",
			wantInsecure: true,
		},
		{
			name:         "unix-abstract scheme preserved",
			env:          "unix-abstract:///grpc.sock",
			wantEndpoint: "unix-abstract:///grpc.sock",
			wantInsecure: true,
		},
		{
			// url.Parse reads "localhost" as the scheme; the endpoint must
			// survive verbatim. Scheme-derived security inference for
			// non-URL values is unchanged by this fix.
			name:         "bare host:port preserved",
			env:          "localhost:4317",
			wantEndpoint: "localhost:4317",
			wantInsecure: true,
		},
		{
			name:         "https URL uses host only",
			env:          "https://env.endpoint:8080/prefix",
			wantEndpoint: "env.endpoint:8080",
			wantInsecure: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", tc.env)
			c := newConfig(nil)
			assert.Equal(t, tc.wantEndpoint, c.endpoint.Value, "endpoint")
			assert.Equal(t, tc.wantInsecure, c.insecure.Value, "insecure")
		})
	}
}
