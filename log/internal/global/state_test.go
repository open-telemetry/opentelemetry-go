// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
)

func resetGlobalLoggerProvider() {
	globalLoggerProvider = defaultLoggerProvider()
	delegateLoggerOnce = sync.Once{}
}

type nonComparableLoggerProvider struct {
	log.LoggerProvider

	nonComparable [0]func() //nolint:structcheck,unused  // This is not called.
}

func TestSetLoggerProvider(t *testing.T) {
	t.Cleanup(resetGlobalLoggerProvider)

	t.Run("Set With default is a noop", func(t *testing.T) {
		t.Cleanup(resetGlobalLoggerProvider)
		SetLoggerProvider(GetLoggerProvider())

		provider, ok := GetLoggerProvider().(*loggerProvider)
		if !ok {
			t.Fatal("Global GetLoggerProvider should be the default logger provider")
		}

		if provider.delegate != nil {
			t.Fatal("logger provider should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		t.Cleanup(resetGlobalLoggerProvider)

		SetLoggerProvider(noop.NewLoggerProvider())

		_, ok := GetLoggerProvider().(*loggerProvider)
		if ok {
			t.Fatal("Global GetLoggerProvider was not changed")
		}
	})

	t.Run("Set() should delegate existing Logger Providers", func(t *testing.T) {
		t.Cleanup(resetGlobalLoggerProvider)

		provider := GetLoggerProvider()

		SetLoggerProvider(noop.NewLoggerProvider())

		dmp := provider.(*loggerProvider)

		if dmp.delegate == nil {
			t.Fatal("The delegated logger providers should have a delegate")
		}
	})

	t.Run("non-comparable types should not panic", func(t *testing.T) {
		t.Cleanup(resetGlobalLoggerProvider)

		provider := nonComparableLoggerProvider{}
		SetLoggerProvider(provider)
		assert.NotPanics(t, func() { SetLoggerProvider(provider) })
	})
}
