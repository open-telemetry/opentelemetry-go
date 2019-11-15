// Copyright 2019, OpenTelemetry Authors
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
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporter/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

type testBatcher struct {
	t             *testing.T
	checkpointSet *test.CheckpointSet
	checkpoints   int
	finishes      int
}

type testExporter struct {
	t       *testing.T
	exports int
	records []export.Record
	retErr  error
}

type testFixture struct {
	checkpointSet *test.CheckpointSet
	batcher       *testBatcher
	exporter      *testExporter
}

type mockClock struct {
	mock *clock.Mock
}

type mockTicker struct {
	ticker *clock.Ticker
}

var _ push.Clock = mockClock{}
var _ push.Ticker = mockTicker{}

func newFixture(t *testing.T) testFixture {
	checkpointSet := test.NewCheckpointSet(sdk.DefaultLabelEncoder())

	batcher := &testBatcher{
		t:             t,
		checkpointSet: checkpointSet,
	}
	exporter := &testExporter{
		t: t,
	}
	return testFixture{
		checkpointSet: checkpointSet,
		batcher:       batcher,
		exporter:      exporter,
	}
}

func (b *testBatcher) AggregatorFor(*export.Descriptor) export.Aggregator {
	return counter.New()
}

func (b *testBatcher) CheckpointSet() export.CheckpointSet {
	b.checkpoints++
	return b.checkpointSet
}

func (b *testBatcher) FinishedCollection() {
	b.finishes++
}

func (b *testBatcher) Process(_ context.Context, record export.Record) error {
	b.checkpointSet.Add(record.Descriptor(), record.Aggregator(), record.Labels().Ordered()...)
	return nil
}

func (e *testExporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	e.exports++
	checkpointSet.ForEach(func(r export.Record) {
		e.records = append(e.records, r)
	})
	return e.retErr
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
	p := push.New(fix.batcher, fix.exporter, time.Second)
	p.Start()
	p.Stop()
	p.Stop()
}

func TestPushDoubleStart(t *testing.T) {
	fix := newFixture(t)
	p := push.New(fix.batcher, fix.exporter, time.Second)
	p.Start()
	p.Start()
	p.Stop()
}

func TestPushTicker(t *testing.T) {
	fix := newFixture(t)

	p := push.New(fix.batcher, fix.exporter, time.Second)
	meter := p.GetMeter("name")

	mock := mockClock{clock.NewMock()}
	p.SetClock(mock)

	ctx := context.Background()

	counter := meter.NewInt64Counter("counter")

	p.Start()

	counter.Add(ctx, 3, meter.Labels())

	require.Equal(t, 0, fix.batcher.checkpoints)
	require.Equal(t, 0, fix.batcher.finishes)
	require.Equal(t, 0, fix.exporter.exports)
	require.Equal(t, 0, len(fix.exporter.records))

	mock.Add(time.Second)
	runtime.Gosched()

	require.Equal(t, 1, fix.batcher.checkpoints)
	require.Equal(t, 1, fix.exporter.exports)
	require.Equal(t, 1, fix.batcher.finishes)
	require.Equal(t, 1, len(fix.exporter.records))
	require.Equal(t, "counter", fix.exporter.records[0].Descriptor().Name())

	sum, err := fix.exporter.records[0].Aggregator().(aggregator.Sum).Sum()
	require.Equal(t, int64(3), sum.AsInt64())
	require.Nil(t, err)

	fix.checkpointSet.Reset()
	fix.exporter.records = nil

	counter.Add(ctx, 7, meter.Labels())

	mock.Add(time.Second)
	runtime.Gosched()

	require.Equal(t, 2, fix.batcher.checkpoints)
	require.Equal(t, 2, fix.batcher.finishes)
	require.Equal(t, 2, fix.exporter.exports)
	require.Equal(t, 1, len(fix.exporter.records))
	require.Equal(t, "counter", fix.exporter.records[0].Descriptor().Name())

	sum, err = fix.exporter.records[0].Aggregator().(aggregator.Sum).Sum()
	require.Equal(t, int64(7), sum.AsInt64())
	require.Nil(t, err)

	p.Stop()
}

func TestPushExportError(t *testing.T) {
	fix := newFixture(t)
	fix.exporter.retErr = fmt.Errorf("Test export error")

	p := push.New(fix.batcher, fix.exporter, time.Second)

	var err error
	p.SetErrorHandler(func(sdkErr error) {
		err = sdkErr
	})

	mock := mockClock{clock.NewMock()}
	p.SetClock(mock)

	p.Start()
	runtime.Gosched()

	require.Equal(t, 0, fix.exporter.exports)
	require.Nil(t, err)

	mock.Add(time.Second)
	runtime.Gosched()

	require.Equal(t, 1, fix.exporter.exports)
	require.Error(t, err)
	require.Equal(t, fix.exporter.retErr, err)

	p.Stop()
}
