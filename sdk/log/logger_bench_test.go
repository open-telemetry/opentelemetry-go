// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

func BenchmarkLoggerNewRecord(b *testing.B) {
	logger := newLogger(NewLoggerProvider(), instrumentation.Scope{})

	r := log.Record{}
	r.SetTimestamp(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	r.SetObservedTimestamp(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	r.SetBody(log.StringValue("testing body value"))
	r.SetSeverity(log.SeverityInfo)
	r.SetSeverityText("testing text")
	r.AddAttributes(
		log.String("k1", "str"),
		log.Float64("k2", 1.0),
		log.Int("k3", 2),
		log.Bool("k4", true),
		log.Bytes("k5", []byte{1}),
	)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.newRecord(context.Background(), r)
		}
	})
}
