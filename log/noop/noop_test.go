// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package noop // import "go.opentelemetry.io/otel/log/noop"

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

func TestImplementationNoPanics(t *testing.T) {
	// Check that if type has an embedded interface and that interface has
	// methods added to it than the No-Op implementation implements them.
	t.Run("LoggerProvider", assertAllExportedMethodNoPanic(
		reflect.ValueOf(LoggerProvider{}),
		reflect.TypeOf((*log.LoggerProvider)(nil)).Elem(),
	))
	t.Run("Logger", assertAllExportedMethodNoPanic(
		reflect.ValueOf(Logger{}),
		reflect.TypeOf((*log.Logger)(nil)).Elem(),
	))
}

func TestNewTracerProvider(t *testing.T) {
	provider := NewLoggerProvider()
	assert.Equal(t, provider, LoggerProvider{})
	logger := provider.Logger("")
	assert.Equal(t, logger, Logger{})
}
