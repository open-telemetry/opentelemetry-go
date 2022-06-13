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

package metric

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/resource"
)

type reader struct {
	producer        producer
	temporalityFunc func(InstrumentKind) Temporality
	collectFunc     func(context.Context) (export.Metrics, error)
	forceFlushFunc  func(context.Context) error
	shutdownFunc    func(context.Context) error
}

var _ Reader = (*reader)(nil)

func (r *reader) register(p producer)                                 { r.producer = p }
func (r *reader) temporality(kind InstrumentKind) Temporality         { return r.temporalityFunc(kind) }
func (r *reader) Collect(ctx context.Context) (export.Metrics, error) { return r.collectFunc(ctx) }
func (r *reader) ForceFlush(ctx context.Context) error                { return r.forceFlushFunc(ctx) }
func (r *reader) Shutdown(ctx context.Context) error                  { return r.shutdownFunc(ctx) }

func TestConfigReaderSignalsEmpty(t *testing.T) {
	f, s := config{}.readerSignals()

	require.NotNil(t, f)
	require.NotNil(t, s)

	ctx := context.Background()
	assert.Nil(t, f(ctx))
	assert.Nil(t, s(ctx))
	assert.ErrorIs(t, s(ctx), ErrReaderShutdown)
}

func TestConfigReaderSignalsForwarded(t *testing.T) {
	var flush, sdown int
	r := &reader{
		forceFlushFunc: func(ctx context.Context) error {
			flush++
			return nil
		},
		shutdownFunc: func(ctx context.Context) error {
			sdown++
			return nil
		},
	}
	c := newConfig([]Option{WithReader(r)})
	f, s := c.readerSignals()

	require.NotNil(t, f)
	require.NotNil(t, s)

	ctx := context.Background()
	assert.NoError(t, f(ctx))
	assert.NoError(t, f(ctx))
	assert.NoError(t, s(ctx))
	assert.ErrorIs(t, s(ctx), ErrReaderShutdown)

	assert.Equal(t, 2, flush, "flush not called 2 times")
	assert.Equal(t, 1, sdown, "shutdown not called 1 time")
}

func TestConfigReaderSignalsForwardedErrors(t *testing.T) {
	r := &reader{
		forceFlushFunc: func(ctx context.Context) error { return assert.AnError },
		shutdownFunc:   func(ctx context.Context) error { return assert.AnError },
	}
	c := newConfig([]Option{WithReader(r)})
	f, s := c.readerSignals()

	require.NotNil(t, f)
	require.NotNil(t, s)

	ctx := context.Background()
	assert.ErrorIs(t, f(ctx), assert.AnError)
	assert.ErrorIs(t, s(ctx), assert.AnError)
	assert.ErrorIs(t, s(ctx), ErrReaderShutdown)
}

func TestUnifyMultiError(t *testing.T) {
	f := func(context.Context) error { return assert.AnError }
	funcs := []func(context.Context) error{f, f, f}
	errs := []error{assert.AnError, assert.AnError, assert.AnError}
	target := fmt.Errorf("%v", errs)
	assert.Equal(t, unify(funcs)(context.Background()), target)
}

func TestWithResource(t *testing.T) {
	res := resource.NewSchemaless()
	c := newConfig([]Option{WithResource(res)})
	assert.Same(t, res, c.res)
}

func TestWithReader(t *testing.T) {
	r := &reader{}
	c := newConfig([]Option{WithReader(r)})
	assert.Contains(t, c.readers, r)
}
