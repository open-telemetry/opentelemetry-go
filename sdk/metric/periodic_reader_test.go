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

//go:build go1.18
// +build go1.18

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

const testDur = time.Second * 2

func TestWithTimeout(t *testing.T) {
	test := func(d time.Duration) time.Duration {
		opts := []PeriodicReaderOption{WithTimeout(d)}
		return newPeriodicReaderConfig(opts).timeout
	}

	assert.Equal(t, testDur, test(testDur))
	assert.Equal(t, defaultTimeout, newPeriodicReaderConfig(nil).timeout)
	assert.Equal(t, defaultTimeout, test(time.Duration(0)), "invalid timeout should use default")
	assert.Equal(t, defaultTimeout, test(time.Duration(-1)), "invalid timeout should use default")
}

func TestWithInterval(t *testing.T) {
	test := func(d time.Duration) time.Duration {
		opts := []PeriodicReaderOption{WithInterval(d)}
		return newPeriodicReaderConfig(opts).interval
	}

	assert.Equal(t, testDur, test(testDur))
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

type periodicReaderTestSuite struct {
	*readerTestSuite

	ErrReader Reader
}

func (ts *periodicReaderTestSuite) SetupTest() {
	ts.readerTestSuite.SetupTest()

	e := &fnExporter{
		exportFunc:   func(context.Context, export.Metrics) error { return assert.AnError },
		flushFunc:    func(context.Context) error { return assert.AnError },
		shutdownFunc: func(context.Context) error { return assert.AnError },
	}

	ts.ErrReader = NewPeriodicReader(e)
}

func (ts *periodicReaderTestSuite) TearDownTest() {
	ts.readerTestSuite.TearDownTest()

	_ = ts.ErrReader.Shutdown(context.Background())
}

func (ts *periodicReaderTestSuite) TestForceFlushPropagated() {
	ts.Equal(assert.AnError, ts.ErrReader.ForceFlush(context.Background()))
}

func (ts *periodicReaderTestSuite) TestShutdownPropagated() {
	ts.Equal(assert.AnError, ts.ErrReader.Shutdown(context.Background()))
}

func TestPeriodicReader(t *testing.T) {
	suite.Run(t, &periodicReaderTestSuite{
		readerTestSuite: &readerTestSuite{
			Factory: func() Reader {
				return NewPeriodicReader(new(fnExporter))
			},
		},
	})
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

func BenchmarkPeriodicReader(b *testing.B) {
	b.Run("Collect", benchReaderCollectFunc(
		NewPeriodicReader(new(fnExporter)),
	))
}

func TestPeriodiclReaderTemporality(t *testing.T) {
	tests := []struct {
		name    string
		options []PeriodicReaderOption
		// Currently only testing constant temporality. This should be expanded
		// if we put more advanced selection in the SDK
		wantTemporality Temporality
	}{
		{
			name:            "default",
			wantTemporality: CumulativeTemporality,
		},
		{
			name: "delta",
			options: []PeriodicReaderOption{
				WithTemporality(deltaTemporalitySelector),
			},
			wantTemporality: DeltaTemporality,
		},
		{
			name: "repeats overwrite",
			options: []PeriodicReaderOption{
				WithTemporality(deltaTemporalitySelector),
				WithTemporality(cumulativeTemporalitySelector),
			},
			wantTemporality: CumulativeTemporality,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdr := NewPeriodicReader(new(fnExporter), tt.options...)
			assert.Equal(t, tt.wantTemporality, rdr.temporality(undefinedInstrument))
		})
	}
}
