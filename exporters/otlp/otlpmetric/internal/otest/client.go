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
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
)

func ClientContextErrorTests(factory func() otlpmetric.Client) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("Shutdown", testCtxErrs(func() func(context.Context) error {
			return factory().Shutdown
		}))

		t.Run("ForceFlush", testCtxErrs(func() func(context.Context) error {
			return factory().ForceFlush
		}))

		t.Run("UploadMetrics", testCtxErrs(func() func(context.Context) error {
			return func(ctx context.Context) error {
				return factory().UploadMetrics(ctx, nil)
			}
		}))
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

		t.Run("NoError", func(t *testing.T) {
			f := factory()
			assert.NoError(t, f(ctx))
		})
	}
}
