// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace_test

import (
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func basicTracerProvider(_ *testing.T) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	return tp
}
