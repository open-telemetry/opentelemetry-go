// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/log"
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
		r.SetSeverityText("text")

		var enabled bool
		for i := 0; ; i++ {
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
