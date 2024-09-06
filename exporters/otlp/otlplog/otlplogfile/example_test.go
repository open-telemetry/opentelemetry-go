// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlplogfile_test

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

func Example() {
	ctx := context.Background()
	exp, err := otlplogfile.New(
		otlplogfile.WithFile("/tmp/otlp-logs.jsonl"),
		otlplogfile.WithFlushInterval(time.Second),
	)
	if err != nil {
		panic(err)
	}

	processor := log.NewBatchProcessor(exp)
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	global.SetLoggerProvider(provider)

	// From here, the provider can be used by instrumentation to collect
	// telemetry.
}
