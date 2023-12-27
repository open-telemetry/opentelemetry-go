// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestSlogHandler(t *testing.T) {
	spy := &spyLogger{}
	l := slog.New(&slogHandler{spy})

	l.InfoContext(ctx, testBody, "string", testString)

	want := log.Record{
		Body:     testBody,
		Severity: log.SeverityInfo,
		Attributes: []attribute.KeyValue{
			attribute.String("string", testString),
		},
	}

	assert.NotZero(t, spy.Record.Timestamp, "should set a timestamp")
	spy.Record.Timestamp = time.Time{}
	assert.Equal(t, want, spy.Record)
	assert.Equal(t, ctx, spy.Context)
}
