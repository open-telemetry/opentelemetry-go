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

package periodic

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/metric/reader"
)

type testExporter struct {
	exportCount   int64
	shutdownCount int64
}

func (e *testExporter) Export(_ context.Context, _ reader.Metrics) error {
	atomic.AddInt64(&e.exportCount, 1)
	return nil
}

func (e *testExporter) Shutdown(_ context.Context) error {
	atomic.AddInt64(&e.shutdownCount, 1)
	return nil
}

func (e *testExporter) Flush(_ context.Context) error {
	return nil
}

type testProducer struct{}

func (testProducer) Produce(_ context.Context, _ *reader.Metrics) reader.Metrics {
	return reader.Metrics{}
}

func Test_exporter_Periodic(t *testing.T) {
	texp := &testExporter{}
	exp := New(5*time.Millisecond, texp, WithTimeout(2*time.Millisecond))
	exp.Register(testProducer{})

	count := atomic.LoadInt64(&texp.exportCount)
	assert.Eventually(t, func() bool {
		newCount := atomic.LoadInt64(&texp.exportCount)
		return newCount >= count+2
	}, time.Second, time.Millisecond)

}

func Test_exporter_Shutdown(t *testing.T) {
	texp := &testExporter{}
	exp := New(5*time.Millisecond, texp, WithTimeout(2*time.Millisecond)).(*exporter)
	exp.Register(testProducer{})

	err := exp.Shutdown(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, int64(1), atomic.LoadInt64(&texp.shutdownCount))

	count := atomic.LoadInt64(&texp.exportCount)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, count, atomic.LoadInt64(&texp.exportCount))
}

func Test_exporter_NilProducer(t *testing.T) {
	texp := &testExporter{}
	_ = New(5*time.Millisecond, texp, WithTimeout(2*time.Millisecond))

	count := atomic.LoadInt64(&texp.exportCount)
	assert.Never(t, func() bool {
		newCount := atomic.LoadInt64(&texp.exportCount)
		return newCount >= count+2
	}, 500*time.Millisecond, time.Millisecond)

}

func Test_exporter_Flush(t *testing.T) {
	texp := &testExporter{}
	exp := New(time.Hour, texp, WithTimeout(2*time.Millisecond)).(*exporter)
	exp.Register(testProducer{})

	count := atomic.LoadInt64(&texp.exportCount)

	exp.collect(context.Background())

	assert.Equal(t, int64(count+1), atomic.LoadInt64(&texp.exportCount))

}
