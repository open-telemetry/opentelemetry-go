// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package noop // import "go.opentelemetry.io/otel/log/noop"

import (
	"context"
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

func assertAllExportedMethodNoPanic(rVal reflect.Value, rType reflect.Type) func(*testing.T) {
	return func(t *testing.T) {
		for n := 0; n < rType.NumMethod(); n++ {
			mType := rType.Method(n)
			if !mType.IsExported() {
				t.Logf("ignoring unexported %s", mType.Name)
				continue
			}
			m := rVal.MethodByName(mType.Name)
			if !m.IsValid() {
				t.Errorf("unknown method for %s: %s", rVal.Type().Name(), mType.Name)
			}

			numIn := mType.Type.NumIn()
			if mType.Type.IsVariadic() {
				numIn--
			}
			args := make([]reflect.Value, numIn)
			ctx := context.Background()
			for i := range args {
				aType := mType.Type.In(i)
				if aType.Name() == "Context" {
					// Do not panic on a nil context.
					args[i] = reflect.ValueOf(ctx)
				} else {
					args[i] = reflect.New(aType).Elem()
				}
			}

			assert.NotPanicsf(t, func() {
				_ = m.Call(args)
			}, "%s.%s", rVal.Type().Name(), mType.Name)
		}
	}
}

func TestNewTracerProvider(t *testing.T) {
	provider := NewLoggerProvider()
	assert.Equal(t, LoggerProvider{}, provider)
	logger := provider.Logger("")
	assert.Equal(t, Logger{}, logger)
}
