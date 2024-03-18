// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/log/internal/global"

import (
	"errors"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/log"
)

var (
	globalLoggerProvider = defaultLoggerProvider()

	delegateLoggerOnce sync.Once
)

func defaultLoggerProvider() *atomic.Value {
	v := &atomic.Value{}
	v.Store(loggerProviderHolder{provider: &loggerProvider{}})
	return v
}

type loggerProviderHolder struct {
	provider log.LoggerProvider
}

// GetLoggerProvider returns the internal implementation for
// global.GetLoggerProvider.
func GetLoggerProvider() log.LoggerProvider {
	return globalLoggerProvider.Load().(loggerProviderHolder).provider
}

// SetLoggerProvider is the internal implementation for
// global.SetLoggerProvider.
func SetLoggerProvider(provider log.LoggerProvider) {
	current := GetLoggerProvider()
	if _, cOk := current.(*loggerProvider); cOk {
		if _, mpOk := provider.(*loggerProvider); mpOk && current == provider {
			// Do not assign the default delegating LoggerProvider to delegate
			// to itself.
			global.Error(
				errors.New("LoggerProvider delegate: self delegation"),
				"No delegate will be configured",
			)
			return
		}
	}

	delegateLoggerOnce.Do(func() {
		if def, ok := current.(*loggerProvider); ok {
			def.setDelegate(provider)
		}
	})
	globalLoggerProvider.Store(loggerProviderHolder{provider: provider})
}
