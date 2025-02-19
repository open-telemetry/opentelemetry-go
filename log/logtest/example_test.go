// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest_test

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"
)

func Example() {
	// Create a recorder.
	rec := logtest.NewRecorder()

	// Emit a record).
	l := rec.Logger("Example")
	ctx := context.Background()
	r := log.Record{}
	r.SetTimestamp(time.Now())
	r.SetSeverity(log.SeverityInfo)
	r.SetBody(log.StringValue("Hello there"))
	l.Emit(ctx, r)

	// Check what was recorded.
	got := rec.Result()
	for _, records := range got {
		for _, record := range records {
			fmt.Printf("%s: %s\n", record.Severity, record.Body)
		}
	}

	// Output:
	// INFO: Hello there
}

func Example_ignoreTimestamp() {
	// Create a recorder.
	rec := logtest.NewRecorder()

	// Emit a record).
	l := rec.Logger("Example")
	ctx := context.Background()
	r := log.Record{}
	r.SetTimestamp(time.Now())
	r.SetSeverity(log.SeverityInfo)
	r.SetBody(log.StringValue("Hello there"))
	l.Emit(ctx, r)

	// Ignore Timestamp.
	got := rec.Result()
	for _, recs := range got {
		for i, r := range recs {
			r.Timestamp = time.Time{}
			recs[i] = r
		}
	}

	for _, records := range got {
		for _, record := range records {
			fmt.Printf("%s: %s: %s\n", record.Timestamp, record.Severity, record.Body)
		}
	}

	// Output:
	// 0001-01-01 00:00:00 +0000 UTC: INFO: Hello there
}
