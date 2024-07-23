// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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

	r10 := r
	r10.AddAttributes(
		log.String("k6", "str"),
		log.Float64("k7", 1.0),
		log.Int("k8", 2),
		log.Bool("k9", true),
		log.Bytes("k10", []byte{1}),
	)

	require.Equal(b, 5, r.AttributesLen())
	require.Equal(b, 10, r10.AttributesLen())

	b.Run("5 attributes", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.newRecord(context.Background(), r)
			}
		})
	})

	b.Run("10 attributes", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.newRecord(context.Background(), r10)
			}
		})
	})
}
