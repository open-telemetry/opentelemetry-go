// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func TestLogrSink(t *testing.T) {
	spy := &spyLogger{}

	l := logr.New(&logrSink{spy})

	l.Info(testBody, "string", testString, "ctx", ctx)

	want := log.Record{
		Body:     testBody,
		Severity: log.SeverityInfo,
		Attributes: []attribute.KeyValue{
			attribute.String("string", testString),
		},
	}

	assert.Equal(t, want, spy.Record)
	assert.Equal(t, ctx, spy.Context)
}
