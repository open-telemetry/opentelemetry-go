// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
)

func TestLoggerProviderConcurrentSafe(*testing.T) {
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

func TestLoggerConcurrentSafe(*testing.T) {
	l := &logger{}

	done := make(chan struct{})
	stop := make(chan struct{})

	go func() {
		defer close(done)

		ctx := context.Background()
		var r log.Record
		var param log.EnabledParameters

		var enabled bool
		for {
			l.Emit(ctx, r)
			enabled = l.Enabled(ctx, param)

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
		l := &testLogger{}
		p.loggers = map[string]*testLogger{name: l}
		p.loggerN++
		return l
	}

	if l, ok := p.loggers[name]; ok {
		return l
	}

	p.loggerN++
	l := &testLogger{}
	p.loggers[name] = l
	return l
}

type testLogger struct {
	embedded.Logger

	emitN, enabledN int
}

func (l *testLogger) Emit(context.Context, log.Record) { l.emitN++ }
func (l *testLogger) Enabled(context.Context, log.EnabledParameters) bool {
	l.enabledN++
	return true
}

func emitRecord(l log.Logger) {
	ctx := context.Background()
	var param log.EnabledParameters
	var r log.Record

	_ = l.Enabled(ctx, param)
	l.Emit(ctx, r)
}

func TestDelegation(t *testing.T) {
	provider := &loggerProvider{}

	const preName = "pre"
	pre0, pre1 := provider.Logger(preName), provider.Logger(preName)
	assert.Same(t, pre0, pre1, "same logger instance not returned")

	alt := provider.Logger("alt")
	assert.NotSame(t, pre0, alt)

	alt2 := provider.Logger(preName, log.WithInstrumentationAttributes(attribute.String("k", "v")))
	assert.NotSame(t, pre0, alt2)

	delegate := &testLoggerProvider{}
	provider.setDelegate(delegate)

	want := 2 // (pre0/pre1) and (alt)
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

func TestLoggerIdentity(t *testing.T) {
	type id struct{ name, ver, url string }

	ids := []id{
		{"name-a", "version-a", "url-a"},
		{"name-a", "version-a", "url-b"},
		{"name-a", "version-b", "url-a"},
		{"name-a", "version-b", "url-b"},
		{"name-b", "version-a", "url-a"},
		{"name-b", "version-a", "url-b"},
		{"name-b", "version-b", "url-a"},
		{"name-b", "version-b", "url-b"},
	}

	provider := &loggerProvider{}
	newLogger := func(i id) log.Logger {
		return provider.Logger(
			i.name,
			log.WithInstrumentationVersion(i.ver),
			log.WithSchemaURL(i.url),
		)
	}

	for i, id0 := range ids {
		for j, id1 := range ids {
			l0, l1 := newLogger(id0), newLogger(id1)

			if i == j {
				assert.Samef(t, l0, l1, "logger(%v) != logger(%v)", id0, id1)
			} else {
				assert.NotSamef(t, l0, l1, "logger(%v) == logger(%v)", id0, id1)
			}
		}
	}
}
