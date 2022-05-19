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

package metric // import "go.opentelemetry.io/otel/sdk/metric/reader"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/export"
)

type readerFactory func() Reader

func testReaderHarness(t *testing.T, f readerFactory) {
	t.Run("ErrorForNotRegistered", func(t *testing.T) {
		r := f()
		ctx := context.Background()

		_, err := r.Collect(ctx)
		require.ErrorIs(t, err, ErrReaderNotRegistered)

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(ctx)
	})

	t.Run("Producer", func(t *testing.T) {
		r := f()
		r.register(testProducer{})
		ctx := context.Background()

		m, err := r.Collect(ctx)
		assert.NoError(t, err)
		assert.Equal(t, testMetrics, m)

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(ctx)
	})

	t.Run("CollectAfterShutdown", func(t *testing.T) {
		r := f()
		r.register(testProducer{})
		require.NoError(t, r.Shutdown(context.Background()))

		m, err := r.Collect(context.Background())
		assert.ErrorIs(t, err, ErrReaderShutdown)
		assert.Equal(t, export.Metrics{}, m)
	})

	t.Run("ShutdownTwice", func(t *testing.T) {
		r := f()
		r.register(testProducer{})
		require.NoError(t, r.Shutdown(context.Background()))

		assert.ErrorIs(t, r.Shutdown(context.Background()), ErrReaderShutdown)
	})

	t.Run("MultipleForceFlush", func(t *testing.T) {
		r := f()
		r.register(testProducer{})
		ctx := context.Background()
		require.NoError(t, r.ForceFlush(ctx))
		assert.NoError(t, r.ForceFlush(ctx))

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(ctx)
	})

	t.Run("MultipleRegister", func(t *testing.T) {
		p0 := testProducer{
			produceFunc: func(ctx context.Context) (export.Metrics, error) {
				// Differentiate this producer from the second by returning an
				// error.
				return testMetrics, assert.AnError
			},
		}
		p1 := testProducer{}

		r := f()
		r.register(p0)
		// This should be ignored.
		r.register(p1)

		ctx := context.Background()
		_, err := r.Collect(ctx)
		assert.Equal(t, assert.AnError, err)

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(ctx)
	})
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
