// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

const twoSec = time.Second * 2

func TestWithTimeout(t *testing.T) {
	test := func(d time.Duration) time.Duration {
		opts := []PeriodicReaderOption{WithTimeout(d)}
		return newPeriodicReaderConfig(opts).timeout
	}

	assert.Equal(t, twoSec, test(twoSec))
	assert.Equal(t, defaultTimeout, newPeriodicReaderConfig(nil).timeout)
	assert.Equal(t, defaultTimeout, test(time.Duration(0)), "invalid timeout should use default")
	assert.Equal(t, defaultTimeout, test(time.Duration(-1)), "invalid timeout should use default")
}

func TestWithInterval(t *testing.T) {
	test := func(d time.Duration) time.Duration {
		opts := []PeriodicReaderOption{WithInterval(d)}
		return newPeriodicReaderConfig(opts).interval
	}

	assert.Equal(t, twoSec, test(twoSec))
	assert.Equal(t, defaultInterval, newPeriodicReaderConfig(nil).interval)
	assert.Equal(t, defaultInterval, test(time.Duration(0)), "invalid interval should use default")
	assert.Equal(t, defaultInterval, test(time.Duration(-1)), "invalid interval should use default")
}

type fnExporter struct {
	exportFunc   func(context.Context, export.Metrics) error
	flushFunc    func(context.Context) error
	shutdownFunc func(context.Context) error
}

var _ Exporter = (*fnExporter)(nil)

func (e *fnExporter) Export(ctx context.Context, m export.Metrics) error {
	if e.exportFunc != nil {
		return e.exportFunc(ctx, m)
	}
	return nil
}

func (e *fnExporter) ForceFlush(ctx context.Context) error {
	if e.flushFunc != nil {
		return e.flushFunc(ctx)
	}
	return nil
}

func (e *fnExporter) Shutdown(ctx context.Context) error {
	if e.shutdownFunc != nil {
		return e.shutdownFunc(ctx)
	}
	return nil
}

func TestPeriodicReader(t *testing.T) {
	testReaderHarness(t, func() Reader {
		return NewPeriodicReader(new(fnExporter))
	})
}

func TestPeriodicReaderForceFlushPropagated(t *testing.T) {
	exp := &fnExporter{
		flushFunc: func(ctx context.Context) error { return assert.AnError },
	}
	r := NewPeriodicReader(exp)
	ctx := context.Background()
	assert.Equal(t, assert.AnError, r.ForceFlush(ctx))

	// Ensure Reader is allowed clean up attempt.
	_ = r.Shutdown(ctx)
}

func TestPeriodicReaderShutdownPropagated(t *testing.T) {
	exp := &fnExporter{
		shutdownFunc: func(ctx context.Context) error { return assert.AnError },
	}
	r := NewPeriodicReader(exp)
	ctx := context.Background()
	assert.Equal(t, assert.AnError, r.Shutdown(ctx))
}

type chErrorHandler struct {
	Err chan error
}

func newChErrorHandler() *chErrorHandler {
	return &chErrorHandler{
		Err: make(chan error, 1),
	}
}

func (eh chErrorHandler) Handle(err error) {
	eh.Err <- err
}

func TestPeriodicReaderRun(t *testing.T) {
	// Override the ticker C chan so tests are not flaky and rely on timing.
	defer func(orig func(time.Duration) *time.Ticker) {
		newTicker = orig
	}(newTicker)
	// Keep this at size zero so when triggered with a send it will hang until
	// the select case is selected and the collection loop is started.
	trigger := make(chan time.Time)
	newTicker = func(d time.Duration) *time.Ticker {
		ticker := time.NewTicker(d)
		ticker.C = trigger
		return ticker
	}

	// Register an error handler to validate export errors are passed to
	// otel.Handle.
	defer func(orig otel.ErrorHandler) {
		otel.SetErrorHandler(orig)
	}(otel.GetErrorHandler())
	eh := newChErrorHandler()
	otel.SetErrorHandler(eh)

	exp := &fnExporter{
		exportFunc: func(_ context.Context, m export.Metrics) error {
			// The testProducer produces testMetrics.
			assert.Equal(t, testMetrics, m)
			return assert.AnError
		},
	}

	r := NewPeriodicReader(exp)
	r.register(testProducer{})
	trigger <- time.Now()
	assert.Equal(t, assert.AnError, <-eh.Err)

	// Ensure Reader is allowed clean up attempt.
	_ = r.Shutdown(context.Background())
}
