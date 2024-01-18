// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func TestWriterLogger(t *testing.T) {
	sb := &strings.Builder{}
	l := &writerLogger{w: sb}

	r := log.Record{}
	r.SetTimestamp(testTimestamp)
	r.SetSeverity(testSeverity)
	r.SetBody(testBody)
	r.AddAttributes(
		attribute.String("string", testString),
		attribute.Float64("float", testFloat),
		attribute.Int("int", testInt),
		attribute.Bool("bool", testBool),
	)
	l.Emit(ctx, r)

	want := "timestamp=595728000 severity=9 body=log message string=7e3b3b2aaeff56a7108fe11e154200dd/7819479873059528190 float=1.2345 int=32768 bool=true traced=true\n"
	assert.Equal(t, want, sb.String())
}
