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

package tracetest

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type rwSpan struct {
	sdktrace.ReadWriteSpan
}

func TestSpanRecorderOnStartAppends(t *testing.T) {
	s0, s1 := new(rwSpan), new(rwSpan)
	ctx := context.Background()
	sr := new(SpanRecorder)

	assert.Len(t, sr.started, 0)
	sr.OnStart(ctx, s0)
	assert.Len(t, sr.started, 1)
	sr.OnStart(ctx, s1)
	assert.Len(t, sr.started, 2)

	// Ensure order correct.
	started := sr.Started()
	assert.Same(t, s0, started[0])
	assert.Same(t, s1, started[1])
}

type roSpan struct {
	sdktrace.ReadOnlySpan
}

func TestSpanRecorderOnEndAppends(t *testing.T) {
	s0, s1 := new(roSpan), new(roSpan)
	sr := new(SpanRecorder)

	assert.Len(t, sr.ended, 0)
	sr.OnEnd(s0)
	assert.Len(t, sr.ended, 1)
	sr.OnEnd(s1)
	assert.Len(t, sr.ended, 2)

	// Ensure order correct.
	ended := sr.Ended()
	assert.Same(t, s0, ended[0])
	assert.Same(t, s1, ended[1])
}

func TestSpanRecorderShutdownNoError(t *testing.T) {
	ctx := context.Background()
	assert.NoError(t, new(SpanRecorder).Shutdown(ctx))

	var c context.CancelFunc
	ctx, c = context.WithCancel(ctx)
	c()
	assert.NoError(t, new(SpanRecorder).Shutdown(ctx))
}

func TestSpanRecorderForceFlushNoError(t *testing.T) {
	ctx := context.Background()
	assert.NoError(t, new(SpanRecorder).ForceFlush(ctx))

	var c context.CancelFunc
	ctx, c = context.WithCancel(ctx)
	c()
	assert.NoError(t, new(SpanRecorder).ForceFlush(ctx))
}

func TestEndingConcurrency(t *testing.T) {
	sr := NewSpanRecorder()
	stop := make(chan struct{})

	var (
		wg    sync.WaitGroup
		count uint64
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		span := new(roSpan)
		for {
			sr.OnEnd(span)
			atomic.AddUint64(&count, 1)
			select {
			case <-stop:
				return
			default:
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			sr.Ended()
			select {
			case <-stop:
				return
			default:
			}
		}
	}()

	for atomic.LoadUint64(&count) < 10 {
		// Wait for things to flow.
		runtime.Gosched()
	}
	close(stop)
	wg.Wait()

	assert.Equal(t, uint64(len(sr.Ended())), atomic.LoadUint64(&count))
}

func TestStartingConcurrency(t *testing.T) {
	sr := NewSpanRecorder()
	stop := make(chan struct{})

	var (
		wg    sync.WaitGroup
		count uint64
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := context.Background()
		span := new(rwSpan)
		for {
			sr.OnStart(ctx, span)
			atomic.AddUint64(&count, 1)
			select {
			case <-stop:
				return
			default:
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			sr.Started()
			select {
			case <-stop:
				return
			default:
			}
		}
	}()

	for atomic.LoadUint64(&count) < 10 {
		// Wait for things to flow.
		runtime.Gosched()
	}
	close(stop)
	wg.Wait()

	assert.Equal(t, uint64(len(sr.Started())), atomic.LoadUint64(&count))
}
