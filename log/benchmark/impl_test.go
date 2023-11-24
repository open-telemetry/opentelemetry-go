// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

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

	r := log.Record{Timestamp: testTimestamp, Severity: testSeverity, Body: testBody}
	r.AddAttributes(
		attribute.String("string", testString),
		attribute.Float64("float", testFloat),
		attribute.Int("int", testInt),
		attribute.Bool("bool", testBool),
	)
	l.Emit(ctx, r)

	want := "timestamp=595728000 severity=9 body=log message string=7e3b3b2aaeff56a7108fe11e154200dd/7819479873059528190 float=1.2345 int=32768 bool=true\n"
	assert.Equal(t, want, sb.String())
}
