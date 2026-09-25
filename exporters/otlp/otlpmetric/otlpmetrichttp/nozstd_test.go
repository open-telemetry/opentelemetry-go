// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build nozstd

package otlpmetrichttp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

// TestWithCompressionZstdFailsFastUnderNoZstdTag checks that a build
// excluding the zstd codec (via the nozstd build tag) fails at
// construction, rather than at first export with an opaque error.
func TestWithCompressionZstdFailsFastUnderNoZstdTag(t *testing.T) {
	_, err := otlpmetrichttp.New(
		context.Background(), //nolint:usetesting // no cleanup needed; New fails before dialing.
		otlpmetrichttp.WithEndpoint("localhost:0"),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithCompression(otlpmetrichttp.ZstdCompression),
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nozstd")
}
