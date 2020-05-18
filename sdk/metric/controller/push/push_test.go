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

package push_test

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/resource"
)

type testExporter struct {
	t         *testing.T
	lock      sync.Mutex
	exports   int
	records   []export.Record
	injectErr func(r export.Record) error
}

type testFixture struct {
	checkpointSet *test.CheckpointSet
	exporter      *testExporter
}

type testSelector struct{}

type mockClock struct {
	mock *clock.Mock
}

type mockTicker struct {
	ticker *clock.Ticker
}

var _ push.Clock = mockClock{}
var _ push.Ticker = mockTicker{}

func newFixture(t *testing.T) testFixture {
	checkpointSet := test.NewCheckpointSet()

	exporter := &testExporter{
		t: t,
	}
	return testFixture{
		checkpointSet: checkpointSet,
		exporter:      exporter,
	}
}

func (testSelector) AggregatorFor(*metric.Descriptor) export.Aggregator {
	return sum.New()
}

// func (b *testIntegrator) Process(_ context.Context, record export.Record) error {
// 	b.lock.Lock()
// 	defer b.lock.Unlock()
// 	labels := record.Labels().ToSlice()
// 	b.checkpointSet.Add(record.Descriptor(), record.Aggregator(), labels...)
// 	return nil
// }

// func (b *testIntegrator) getCounts() (checkpoints, finishes int) {
// 	b.lock.Lock()
// 	defer b.lock.Unlock()
// 	return b.checkpoints, b.finishes
// }

func (e *testExporter) Export(_ context.Context, _ *resource.Resource, checkpointSet export.CheckpointSet) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.exports++
	var records []export.Record
	if err := checkpointSet.ForEach(func(r export.Record) error {
		if e.injectErr != nil {
			if err := e.injectErr(r); err != nil {
				return err
			}
		}
		records = append(records, r)
		return nil
	}); err != nil {
		return err
	}
	e.records = records
	return nil
}

func (e *testExporter) resetRecords() ([]export.Record, int) {
	e.lock.Lock()
	defer e.lock.Unlock()
	r := e.records
	e.records = nil
	return r, e.exports
}

func (c mockClock) Now() time.Time {
	return c.mock.Now()
}

func (c mockClock) Ticker(period time.Duration) push.Ticker {
	return mockTicker{c.mock.Ticker(period)}
}

func (c mockClock) Add(d time.Duration) {
	c.mock.Add(d)
}

func (t mockTicker) Stop() {
	t.ticker.Stop()
}

func (t mockTicker) C() <-chan time.Time {
	return t.ticker.C
}

func TestPushDoubleStop(t *testing.T) {
	fix := newFixture(t)
	p := push.New(testSelector{}, fix.exporter)
	p.Start()
	p.Stop()
	p.Stop()
}

func TestPushDoubleStart(t *testing.T) {
	fix := newFixture(t)
	p := push.New(testSelector{}, fix.exporter)
	p.Start()
	p.Start()
	p.Stop()
}

func TestPushTicker(t *testing.T) {
	fix := newFixture(t)

	p := push.New(testSelector{}, fix.exporter, push.WithPeriod(time.Second))
	meter := p.Meter("name")

	mock := mockClock{clock.NewMock()}
	p.SetClock(mock)

	ctx := context.Background()

	counter := metric.Must(meter).NewInt64Counter("counter")

	p.Start()

	counter.Add(ctx, 3)

	records, exports := fix.exporter.resetRecords()
	require.Equal(t, 0, exports)
	require.Equal(t, 0, len(records))

	mock.Add(time.Second)
	runtime.Gosched()

	fmt.Println("NOW")

	records, exports = fix.exporter.resetRecords()
	require.Equal(t, 1, exports)
	require.Equal(t, 1, len(records))
	require.Equal(t, "counter", records[0].Descriptor().Name())

	sum, err := records[0].Aggregator().(aggregator.Sum).Sum()
	require.Equal(t, int64(3), sum.AsInt64())
	require.Nil(t, err)

	fix.checkpointSet.Reset()

	counter.Add(ctx, 7)

	mock.Add(time.Second)
	runtime.Gosched()

	records, exports = fix.exporter.resetRecords()
	require.Equal(t, 2, exports)
	require.Equal(t, 1, len(records))
	require.Equal(t, "counter", records[0].Descriptor().Name())

	sum, err = records[0].Aggregator().(aggregator.Sum).Sum()
	require.Equal(t, int64(7), sum.AsInt64())
	require.Nil(t, err)

	p.Stop()
}

func TestPushExportError(t *testing.T) {
	injector := func(name string, e error) func(r export.Record) error {
		return func(r export.Record) error {
			if r.Descriptor().Name() == name {
				return e
			}
			return nil
		}
	}
	var errAggregator = fmt.Errorf("unexpected error")
	var tests = []struct {
		name                string
		injectedError       error
		expectedDescriptors []string
		expectedError       error
	}{
		{"errNone", nil, []string{"counter1", "counter2"}, nil},
		{"errNoData", aggregator.ErrNoData, []string{"counter2"}, nil},
		{"errUnexpected", errAggregator, []string{}, errAggregator},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fix := newFixture(t)
			fix.exporter.injectErr = injector("counter1", tt.injectedError)

			p := push.New(testSelector{}, fix.exporter, push.WithPeriod(time.Second))

			var err error
			var lock sync.Mutex
			p.SetErrorHandler(func(sdkErr error) {
				lock.Lock()
				defer lock.Unlock()
				err = sdkErr
			})

			mock := mockClock{clock.NewMock()}
			p.SetClock(mock)

			ctx := context.Background()

			meter := p.Meter("name")
			counter1 := metric.Must(meter).NewInt64Counter("counter1")
			counter2 := metric.Must(meter).NewInt64Counter("counter2")

			p.Start()
			runtime.Gosched()

			counter1.Add(ctx, 3)
			counter2.Add(ctx, 5)

			require.Equal(t, 0, fix.exporter.exports)
			require.Nil(t, err)

			mock.Add(time.Second)
			runtime.Gosched()

			records, exports := fix.exporter.resetRecords()
			require.Equal(t, 1, exports)
			lock.Lock()
			if tt.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
			}
			lock.Unlock()
			require.Equal(t, len(tt.expectedDescriptors), len(records))
			for _, r := range records {
				require.Contains(t, tt.expectedDescriptors, r.Descriptor().Name())
			}

			p.Stop()

		})
	}
}
