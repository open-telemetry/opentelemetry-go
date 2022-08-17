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

package otlpmetrichttp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientHonorsContextErrors(t *testing.T) {
	t.Run("Shutdown", testCtxErr(func(t *testing.T) func(context.Context) error {
		c, err := newClient()
		require.NoError(t, err)
		return c.Shutdown
	}))

	t.Run("ForceFlush", testCtxErr(func(t *testing.T) func(context.Context) error {
		c, err := newClient()
		require.NoError(t, err)
		return c.ForceFlush
	}))

	t.Run("UploadMetrics", testCtxErr(func(t *testing.T) func(context.Context) error {
		c, err := newClient()
		require.NoError(t, err)
		return func(ctx context.Context) error {
			return c.UploadMetrics(ctx, nil)
		}
	}))
}

func testCtxErr(factory func(*testing.T) func(context.Context) error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		t.Run("DeadlineExceeded", func(t *testing.T) {
			innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
			t.Cleanup(innerCancel)
			<-innerCtx.Done()

			f := factory(t)
			assert.ErrorIs(t, f(innerCtx), context.DeadlineExceeded)
		})

		t.Run("Canceled", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			f := factory(t)
			assert.ErrorIs(t, f(innerCtx), context.Canceled)
		})
	}
}
