// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/log/global"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
)

type testLoggerProvider struct{ embedded.LoggerProvider }

var _ log.LoggerProvider = &testLoggerProvider{}

func (*testLoggerProvider) Logger(_ string, _ ...log.LoggerOption) log.Logger {
	return noop.NewLoggerProvider().Logger("")
}

func TestMultipleGlobalLoggerProvider(t *testing.T) {
	p1, p2 := testLoggerProvider{}, noop.NewLoggerProvider()

	SetLoggerProvider(&p1)
	SetLoggerProvider(p2)

	assert.Equal(t, p2, GetLoggerProvider())
}
