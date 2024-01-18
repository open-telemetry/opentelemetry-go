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

	want := log.Record{}
	want.SetBody(testBody)
	want.SetSeverity(log.SeverityInfo)
	want.AddAttributes(attribute.String("string", testString))

	assert.Equal(t, testBody, spy.Record.Body())
	assert.Equal(t, log.SeverityInfo, spy.Record.Severity())
	assert.Equal(t, 1, spy.Record.AttributesLen())
	spy.Record.WalkAttributes(func(kv attribute.KeyValue) bool {
		assert.Equal(t, "string", string(kv.Key))
		assert.Equal(t, testString, kv.Value.AsString())
		return true
	})
}
