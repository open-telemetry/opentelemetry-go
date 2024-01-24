// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

func TestLogrSink(t *testing.T) {
	spy := &spyLogger{}

	l := logr.New(&logrSink{spy})

	l.Info(testBodyString, "string", testString, "ctx", ctx)

	want := log.Record{}
	want.SetBody(testBody)
	want.SetSeverity(log.SeverityInfo)
	want.AddAttributes(log.String("string", testString))

	assert.Equal(t, testBody, spy.Record.Body())
	assert.Equal(t, log.SeverityInfo, spy.Record.Severity())
	assert.Equal(t, 1, spy.Record.AttributesLen())
	spy.Record.WalkAttributes(func(kv log.KeyValue) bool {
		assert.Equal(t, "string", string(kv.Key))
		assert.Equal(t, testString, kv.Value.String())
		return true
	})
}
