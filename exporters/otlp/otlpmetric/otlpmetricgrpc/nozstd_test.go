// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build nozstd

package otlpmetricgrpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
)

// TestWithCompressorZstdFailsFastUnderNoZstdTag checks that a build
// excluding the zstd codec (via the nozstd build tag) fails at construction,
// rather than at first export with an opaque grpc-internal error.
func TestWithCompressorZstdFailsFastUnderNoZstdTag(t *testing.T) {
	_, err := otlpmetricgrpc.New(
		context.Background(), //nolint:usetesting // no cleanup needed; New fails before dialing.
		otlpmetricgrpc.WithCompressor("zstd"),
		otlpmetricgrpc.WithEndpoint("localhost:0"),
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nozstd")
}
