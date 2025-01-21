// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
)

func TestAssertRecord(t *testing.T) {
	r1 := log.Record{}
	r2 := log.Record{}
	AssertRecordEqual(t, r1, r2)

	now := time.Now()
	r1.SetEventName("my_event")
	r2.SetEventName("my_event")
	r1.SetTimestamp(now)
	r2.SetTimestamp(now)
	r1.SetObservedTimestamp(now)
	r2.SetObservedTimestamp(now)
	r1.SetSeverity(log.SeverityTrace1)
	r2.SetSeverity(log.SeverityTrace1)
	r1.SetSeverityText("trace")
	r2.SetSeverityText("trace")
	r1.SetBody(log.StringValue("log body"))
	r2.SetBody(log.StringValue("log body"))
	r1.AddAttributes(log.Bool("attr", true))
	r2.AddAttributes(log.Bool("attr", true))
	AssertRecordEqual(t, r1, r2)
}
