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
	"sync"
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

	// Requires the race-detector (a default test option for the project).
	t.Run("MethodConcurrency", func(t *testing.T) {
		// All reader methods should be concurrent-safe.
		r := f()
		r.register(testProducer{})
		ctx := context.Background()

		var wg sync.WaitGroup
		const threads = 2
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = r.Collect(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = r.ForceFlush(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = r.Shutdown(ctx)
			}()
		}
		wg.Wait()
	})

	t.Run("ShutdownBeforeRegister", func(t *testing.T) {
		r := f()

		err := r.Shutdown(context.Background())
		require.NoError(t, err)

		// Registering after the reader is shutdown, while not expected user
		// behavior, needs to be not reset the reader to not be shutdown.
		r.register(testProducer{})

		m, err := r.Collect(context.Background())
		assert.ErrorIs(t, err, ErrReaderShutdown)
		assert.Equal(t, export.Metrics{}, m)
	})
}

var testMetrics = export.Metrics{
	// TODO: test with actual data.
}

type testProducer struct{}

func (p testProducer) produce(context.Context) (export.Metrics, error) {
	return testMetrics, nil
}
