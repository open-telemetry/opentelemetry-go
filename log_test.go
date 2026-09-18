// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
)

type testLoggerProvider struct {
	embedded.LoggerProvider

	logger log.Logger
}

var _ log.LoggerProvider = &testLoggerProvider{}

func (p *testLoggerProvider) Logger(string, ...log.LoggerOption) log.Logger {
	if p.logger != nil {
		return p.logger
	}
	return noop.NewLoggerProvider().Logger("")
}

type testLogger struct{ embedded.Logger }

func (*testLogger) Emit(context.Context, log.Record) {}

func (*testLogger) Enabled(context.Context, log.EnabledParameters) bool { return false }

func TestMultipleGlobalLoggerProvider(t *testing.T) {
	p1 := testLoggerProvider{}
	p2 := noop.NewLoggerProvider()
	SetLoggerProvider(&p1)
	SetLoggerProvider(p2)

	assert.Equal(t, p2, GetLoggerProvider())
}

func TestLogger(t *testing.T) {
	provider := &testLoggerProvider{logger: &testLogger{}}
	SetLoggerProvider(provider)

	assert.Same(t, provider.logger, Logger("test"))
}
