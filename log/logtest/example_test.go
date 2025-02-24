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
	r.AddAttributes(log.String("fizz", "buzz"))
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
	// Ignore Context.
	cmpCtx := cmpopts.IgnoreFields(logtest.Record{}, "Context")
	// Ignore Timestamps.
	cmpStmps := cmpopts.IgnoreTypes(time.Time{})
	// Unordered compare of the key values.
	cmpKVs := cmpopts.SortSlices(func(a, b log.KeyValue) bool { return a.Key < b.Key })
	// Empty and nil collections are equal.
	cmpEpty := cmpopts.EquateEmpty()

	// Get the recorded log records.
	got := rec.Result()
	if diff := cmp.Diff(got, want, cmpCtx, cmpKVs, cmpStmps, cmpEpty); diff != "" {
		fmt.Printf("Recorded records mismatch (-want +got):\n%s", diff)
	}

	// Output:
	// Recorded records mismatch (-want +got):
	//   logtest.Recording{
	//   	{Name: "Example"}: {
	//   		{
	//   			... // 3 ignored and 2 identical fields
	//   			SeverityText: "",
	//   			Body:         s"Hello there",
	//   			Attributes: []log.KeyValue(Inverse(cmpopts.SortSlices, []log.KeyValue{
	// - 				s"fizz:buzz",
	// + 				s"foo:bar",
	//   				s"n:1",
	//   			})),
	//   		},
	//   	},
	//   }
}
