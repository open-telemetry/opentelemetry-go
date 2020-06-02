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

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	controllerTest "go.opentelemetry.io/otel/sdk/metric/controller/test"
	"go.opentelemetry.io/otel/sdk/resource"
)

var testResource = resource.New(kv.String("R", "V"))

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

func newFixture(t *testing.T) testFixture {
	checkpointSet := test.NewCheckpointSet(testResource)

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

func (e *testExporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
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

	p := push.New(
		testSelector{},
		fix.exporter,
		push.WithPeriod(time.Second),
		push.WithResource(testResource),
	)
	meter := p.Provider().Meter("name")

	mock := controllerTest.NewMockClock()
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

	records, exports = fix.exporter.resetRecords()
	require.Equal(t, 1, exports)
	require.Equal(t, 1, len(records))
	require.Equal(t, "counter", records[0].Descriptor().Name())
	require.Equal(t, "R=V", records[0].Resource().Encoded(label.DefaultEncoder()))

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
	require.Equal(t, "R=V", records[0].Resource().Encoded(label.DefaultEncoder()))

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
		{"errNone", nil, []string{"counter1{R=V,X=Y}", "counter2{R=V,}"}, nil},
		{"errNoData", aggregator.ErrNoData, []string{"counter2{R=V,}"}, nil},
		{"errUnexpected", errAggregator, []string{}, errAggregator},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fix := newFixture(t)
			fix.exporter.injectErr = injector("counter1", tt.injectedError)

			p := push.New(
				testSelector{},
				fix.exporter,
				push.WithPeriod(time.Second),
				push.WithResource(testResource),
			)

			var err error
			var lock sync.Mutex
			p.SetErrorHandler(func(sdkErr error) {
				lock.Lock()
				defer lock.Unlock()
				err = sdkErr
			})

			mock := controllerTest.NewMockClock()
			p.SetClock(mock)

			ctx := context.Background()

			meter := p.Provider().Meter("name")
			counter1 := metric.Must(meter).NewInt64Counter("counter1")
			counter2 := metric.Must(meter).NewInt64Counter("counter2")

			p.Start()
			runtime.Gosched()

			counter1.Add(ctx, 3, kv.String("X", "Y"))
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
				require.Contains(t, tt.expectedDescriptors,
					fmt.Sprintf("%s{%s,%s}",
						r.Descriptor().Name(),
						r.Resource().Encoded(label.DefaultEncoder()),
						r.Labels().Encoded(label.DefaultEncoder()),
					),
				)
			}

			p.Stop()
		})
	}
}
