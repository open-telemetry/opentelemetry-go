// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build nozstd

package otlptracehttp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

// TestWithCompressionZstdFailsFastUnderNoZstdTag checks that a build
// excluding the zstd codec (via the nozstd build tag) fails at Start,
// rather than at first export with an opaque error.
func TestWithCompressionZstdFailsFastUnderNoZstdTag(t *testing.T) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("localhost:0"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(otlptracehttp.ZstdCompression),
	)
	err := client.Start(context.Background()) //nolint:usetesting // no cleanup needed; Start fails before dialing.
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nozstd")
}
