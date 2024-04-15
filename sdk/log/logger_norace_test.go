// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !race

package log

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

func TestAllocationLimits(t *testing.T) {
	// This test is not run with a race detector. The sync.Pool used by parts
	// of the SDK has memory optimizations removed for the race detector. Do
	// not test performance of the SDK in that state.

	const runs = 10

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

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		logger.newRecord(context.Background(), r)
	}), "newRecord")
}
