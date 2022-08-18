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

package otest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otest"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
)

// ClientFactory is a function that when called returns a
// otlpmetric.Client implementation that is connected to also returned
// Collector implementation. The Client is ready to upload metric data to the
// Collector which is ready to store that data.
type ClientFactory func() (otlpmetric.Client, Collector)

func RunClientTests(f ClientFactory) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Run("ClientHonorsContextErrors", func(t *testing.T) {
			t.Run("Shutdown", testCtxErrs(func() func(context.Context) error {
				c, _ := f()
				return c.Shutdown
			}))

			t.Run("ForceFlush", testCtxErrs(func() func(context.Context) error {
				c, _ := f()
				return c.ForceFlush
			}))

			t.Run("UploadMetrics", testCtxErrs(func() func(context.Context) error {
				c, _ := f()
				return func(ctx context.Context) error {
					return c.UploadMetrics(ctx, nil)
				}
			}))
		})

		t.Run("ForceFlushFlushes", func(t *testing.T) {
			ctx := context.Background()
			client, collector := f()
			require.NoError(t, client.UploadMetrics(ctx, resourceMetrics))

			require.NoError(t, client.ForceFlush(ctx))
			rm := collector.Collect().dump()
			// Data correctness is not important, just it was received.
			require.Greater(t, len(rm), 0, "no data uploaded")

			require.NoError(t, client.Shutdown(ctx))
			rm = collector.Collect().dump()
			assert.Len(t, rm, 0, "client did not flush all data")
		})

		t.Run("UploadMetrics", testUploadMetrics(f))
	}
}

func testCtxErrs(factory func() func(context.Context) error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		t.Run("DeadlineExceeded", func(t *testing.T) {
			innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
			t.Cleanup(innerCancel)
			<-innerCtx.Done()

			f := factory()
			assert.ErrorIs(t, f(innerCtx), context.DeadlineExceeded)
		})

		t.Run("Canceled", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			f := factory()
			assert.ErrorIs(t, f(innerCtx), context.Canceled)
		})
	}
}
