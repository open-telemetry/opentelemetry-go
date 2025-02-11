// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package xlog

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func BenchmarkLoggerEnabled(b *testing.B) {
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(newFltrProcessor(false)),
		sdklog.WithProcessor(newFltrProcessor(true)),
	)
	logger := provider.Logger(b.Name())
	ctx := context.Background()
	param := log.EnabledParameters{Severity: log.SeverityDebug}
	var enabled bool

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		enabled = logger.Enabled(ctx, param)
	}

	_ = enabled
}
