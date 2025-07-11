// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest_test

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"
)

func Example() {
	t := &testing.T{} // Provided by testing framework.

	// Create a recorder.
	rec := logtest.NewRecorder()

	// Emit a log record (code under test).
	l := rec.Logger("Example")
	r := log.Record{}
	r.SetTimestamp(time.Now())
	r.SetSeverity(log.SeverityInfo)
	r.SetBody(log.StringValue("Hello there"))
	r.AddAttributes(log.String("foo", "bar"))
	r.AddAttributes(log.Int("n", 1))
	l.Emit(context.Background(), r)

	// Verify that the expected and actual log records match.
	want := logtest.Recording{
		logtest.Scope{Name: "Example"}: []logtest.Record{
			{
				Severity: log.SeverityInfo,
				Body:     log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.Int("n", 1),
					log.String("foo", "bar"),
				},
			},
		},
	}
	got := rec.Result()
	logtest.AssertEqual(t, want, got,
		logtest.Transform(func(r logtest.Record) logtest.Record {
			r = r.Clone()
			r.Context = nil           // Ignore context.
			r.Timestamp = time.Time{} // Ignore timestamp.
			return r
		}),
	)
	// Output:
	//
}
