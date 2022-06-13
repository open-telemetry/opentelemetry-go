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

package metric // import "go.opentelemetry.io/otel/sdk/metric/reader"

import (
	"context"
	"sync"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

type readerTestSuite struct {
	suite.Suite

	Factory func() Reader
	Reader  Reader
}

func (ts *readerTestSuite) SetupSuite() {
	otel.SetLogger(testr.New(ts.T()))
}

func (ts *readerTestSuite) SetupTest() {
	ts.Reader = ts.Factory()
}

func (ts *readerTestSuite) TearDownTest() {
	// Ensure Reader is allowed attempt to clean up.
	_ = ts.Reader.Shutdown(context.Background())
}

func (ts *readerTestSuite) TestErrorForNotRegistered() {
	_, err := ts.Reader.Collect(context.Background())
	ts.ErrorIs(err, ErrReaderNotRegistered)
}

func (ts *readerTestSuite) TestProducer() {
	ts.Reader.register(testProducer{})
	m, err := ts.Reader.Collect(context.Background())
	ts.NoError(err)
	ts.Equal(testMetrics, m)
}

func (ts *readerTestSuite) TestCollectAfterShutdown() {
	ctx := context.Background()
	ts.Reader.register(testProducer{})
	ts.Require().NoError(ts.Reader.Shutdown(ctx))

	m, err := ts.Reader.Collect(ctx)
	ts.ErrorIs(err, ErrReaderShutdown)
	ts.Equal(export.Metrics{}, m)
}

func (ts *readerTestSuite) TestShutdownTwice() {
	ctx := context.Background()
	ts.Reader.register(testProducer{})
	ts.Require().NoError(ts.Reader.Shutdown(ctx))
	ts.ErrorIs(ts.Reader.Shutdown(ctx), ErrReaderShutdown)
}

func (ts *readerTestSuite) TestMultipleForceFlush() {
	ctx := context.Background()
	ts.Reader.register(testProducer{})
	ts.Require().NoError(ts.Reader.ForceFlush(ctx))
	ts.NoError(ts.Reader.ForceFlush(ctx))
}

func (ts *readerTestSuite) TestMultipleRegister() {
	p0 := testProducer{
		produceFunc: func(ctx context.Context) (export.Metrics, error) {
			// Differentiate this producer from the second by returning an
			// error.
			return testMetrics, assert.AnError
		},
	}
	p1 := testProducer{}

	ts.Reader.register(p0)
	// This should be ignored.
	ts.Reader.register(p1)

	_, err := ts.Reader.Collect(context.Background())
	ts.Equal(assert.AnError, err)
}

func (ts *readerTestSuite) TestMethodConcurrency() {
	// Requires the race-detector (a default test option for the project).

	// All reader methods should be concurrent-safe.
	ts.Reader.register(testProducer{})
	ctx := context.Background()

	var wg sync.WaitGroup
	const threads = 2
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = ts.Reader.Collect(ctx)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Reader.ForceFlush(ctx)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Reader.Shutdown(ctx)
		}()
	}
	wg.Wait()
}

func (ts *readerTestSuite) TestShutdownBeforeRegister() {
	ctx := context.Background()
	ts.Require().NoError(ts.Reader.Shutdown(ctx))
	// Registering after shutdown should not revert the shutdown.
	ts.Reader.register(testProducer{})

	m, err := ts.Reader.Collect(ctx)
	ts.ErrorIs(err, ErrReaderShutdown)
	ts.Equal(export.Metrics{}, m)
}

var testMetrics = export.Metrics{
	// TODO: test with actual data.
}

type testProducer struct {
	produceFunc func(context.Context) (export.Metrics, error)
}

func (p testProducer) produce(ctx context.Context) (export.Metrics, error) {
	if p.produceFunc != nil {
		return p.produceFunc(ctx)
	}
	return testMetrics, nil
}

func benchReaderCollectFunc(r Reader) func(*testing.B) {
	ctx := context.Background()
	r.register(testProducer{})

	// Store bechmark results in a closure to prevent the compiler from
	// inlining and skipping the function.
	var (
		collectedMetrics export.Metrics
		err              error
	)

	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			collectedMetrics, err = r.Collect(ctx)
			if collectedMetrics != testMetrics || err != nil {
				b.Errorf("unexpected Collect response: (%#v, %v)", collectedMetrics, err)
			}
		}
	}
}
