// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
)

func TestMultipleGlobalLoggerProvider(t *testing.T) {
	type provider struct{ log.LoggerProvider }

	p1, p2 := provider{}, noop.NewLoggerProvider()

	SetLoggerProvider(&p1)
	SetLoggerProvider(p2)

	assert.Equal(t, p2, GetLoggerProvider())
}

func TestLogger(t *testing.T) {
	provider := &testLoggerProvider{logger: &testLogger{}}
	SetLoggerProvider(provider)

	assert.Same(t, provider.logger, Logger("test"))
}

type testLoggerProvider struct {
	embedded.LoggerProvider

	logger log.Logger
}

func (p *testLoggerProvider) Logger(string, ...log.LoggerOption) log.Logger {
	return p.logger
}

type testLogger struct{ embedded.Logger }

func (*testLogger) Emit(context.Context, log.Record) {}

func (*testLogger) Enabled(context.Context, log.EnabledParameters) bool { return false }
