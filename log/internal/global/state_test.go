// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
)

func TestSetLoggerProvider(t *testing.T) {
	reset := func() {
		globalLoggerProvider = defaultLoggerProvider()
		delegateLoggerOnce = sync.Once{}
	}

	t.Run("Set With default is a noop", func(t *testing.T) {
		t.Cleanup(reset)

		t.Cleanup(func(orig logr.Logger) func() {
			global.SetLogger(testr.New(t)) // Don't pollute output.
			return func() { global.SetLogger(orig) }
		}(global.GetLogger()))
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
		t.Cleanup(reset)

		SetLoggerProvider(noop.NewLoggerProvider())
		if _, ok := GetLoggerProvider().(*loggerProvider); ok {
			t.Fatal("Global GetLoggerProvider was not changed")
		}
	})

	t.Run("Set() should delegate existing Logger Providers", func(t *testing.T) {
		t.Cleanup(reset)

		provider := GetLoggerProvider()
		SetLoggerProvider(noop.NewLoggerProvider())

		if del := provider.(*loggerProvider); del.delegate == nil {
			t.Fatal("The delegated logger providers should have a delegate")
		}
	})

	t.Run("non-comparable types should not panic", func(t *testing.T) {
		t.Cleanup(reset)

		type nonComparableLoggerProvider struct {
			log.LoggerProvider
			noCmp [0]func() //nolint:unused  // This is indeed used.
		}

		provider := nonComparableLoggerProvider{}
		SetLoggerProvider(provider)
		assert.NotPanics(t, func() { SetLoggerProvider(provider) })
	})
}
