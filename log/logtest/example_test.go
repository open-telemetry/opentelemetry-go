// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"
)

func Example() {
	// Create a recorder.
	rec := logtest.NewRecorder()

	// Emit a log record (code under test).
	l := rec.Logger("Example")
	ctx := context.Background()
	r := log.Record{}
	r.SetTimestamp(time.Now())
	r.SetSeverity(log.SeverityInfo)
	r.SetBody(log.StringValue("Hello there"))
	r.AddAttributes(log.String("foo", "bar"))
	r.AddAttributes(log.Int("n", 1))
	l.Emit(ctx, r)

	// Expected log records.
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
	opts := []cmpopts.Options{ 
		cmpopts.IgnoreFields(logtest.Record{}, "Context"), // Ignore Context.
		cmpopts.IgnoreTypes(time.Time{}), // Ignore Timestamps.
		cmpopts.SortSlices(func(a, b log.KeyValue) bool { return a.Key < b.Key }), // Unordered compare of the key values.
		cmpopts.EquateEmpty(), // Empty and nil collections are equal.
	}
	// Get the recorded log records.
	got := rec.Result()
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		fmt.Printf("recording mismatch (-want +got):\n%s", diff)
	}

	// Output:
	//
}
