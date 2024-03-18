// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
)

func TestLoggerProviderConcurrentSafe(t *testing.T) {
	p := &loggerProvider{}

	done := make(chan struct{})
	stop := make(chan struct{})

	go func() {
		defer close(done)
		var logger log.Logger
		for i := 0; ; i++ {
			logger = p.Logger(fmt.Sprintf("a%d", i))
			select {
			case <-stop:
				_ = logger
				return
			default:
			}
		}
	}()

	p.setDelegate(noop.NewLoggerProvider())
	close(stop)
	<-done
}

func TestLoggerConcurrentSafe(t *testing.T) {
	l := newLogger("", nil)

	done := make(chan struct{})
	stop := make(chan struct{})

	go func() {
		defer close(done)

		ctx := context.Background()
		var r log.Record

		var enabled bool
		for {
			l.Emit(ctx, r)
			enabled = l.Enabled(ctx, r)

			select {
			case <-stop:
				_ = enabled
				return
			default:
			}
		}
	}()

	l.setDelegate(noop.NewLoggerProvider())
	close(stop)
	<-done
}

type testLoggerProvider struct {
	embedded.LoggerProvider

	loggers map[string]*testLogger
	loggerN int
}

func (p *testLoggerProvider) Logger(name string, _ ...log.LoggerOption) log.Logger {
	if p.loggers == nil {
		l := &testLogger{name: name}
		p.loggers = map[string]*testLogger{name: l}
		p.loggerN++
		return l
	}

	if l, ok := p.loggers[name]; ok {
		return l
	}

	p.loggerN++
	l := &testLogger{name: name}
	p.loggers[name] = l
	return &testLogger{}
}

type testLogger struct {
	embedded.Logger

	name            string
	emitN, enabledN int
}

func (l *testLogger) Emit(context.Context, log.Record) { l.emitN++ }
func (l *testLogger) Enabled(context.Context, log.Record) bool {
	l.enabledN++
	return true
}

func emitRecord(l log.Logger) {
	ctx := context.Background()
	var r log.Record

	_ = l.Enabled(ctx, r)
	l.Emit(ctx, r)
}

func TestDelegation(t *testing.T) {
	delegate := &testLoggerProvider{}
	require.Equal(t, 0, delegate.loggerN, "invalid setup")

	provider := &loggerProvider{}

	const preName = "pre"
	pre0, pre1 := provider.Logger(preName), provider.Logger(preName)
	assert.Same(t, pre0, pre1, "same logger instance not returned")

	alt := provider.Logger("alt")
	assert.NotSame(t, pre0, alt)

	want := 0
	if !assert.Equal(t, want, delegate.loggerN, "delegate Logger created") {
		want = delegate.loggerN
	}

	provider.setDelegate(delegate)
	want += 2 // (pre0/pre1) and (alt)
	if !assert.Equal(t, want, delegate.loggerN, "previous Loggers not delegated") {
		want = delegate.loggerN
	}

	pre2 := provider.Logger(preName)
	if !assert.Equal(t, want, delegate.loggerN, "previous Logger recreated") {
		want = delegate.loggerN
	}

	post := provider.Logger("test")
	want++
	assert.Equal(t, want, delegate.loggerN, "new Logger not delegated")

	emitRecord(pre0)
	emitRecord(pre2)

	if assert.IsType(t, &testLogger{}, pre2, "wrong pre-delegation Logger type") {
		assert.Equal(t, 2, pre2.(*testLogger).emitN, "Emit not delegated")
		assert.Equal(t, 2, pre2.(*testLogger).enabledN, "Enabled not delegated")
	}

	emitRecord(post)

	if assert.IsType(t, &testLogger{}, post, "wrong post-delegation Logger type") {
		assert.Equal(t, 1, post.(*testLogger).emitN, "Emit not delegated")
		assert.Equal(t, 1, post.(*testLogger).enabledN, "Enabled not delegated")
	}
}
