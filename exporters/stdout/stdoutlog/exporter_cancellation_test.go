// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func TestExporterBlockedWriteContext(t *testing.T) {
	for _, deadline := range []bool{false, true} {
		name := "canceled"
		if deadline {
			name = "deadline"
		}
		t.Run(name, func(t *testing.T) {
			for _, count := range []int{1, 2} {
				t.Run(strconv.Itoa(count)+" records", func(t *testing.T) {
					w, release := newBlockedWriter(t)
					e, err := New(WithWriter(w))
					require.NoError(t, err)
					ctx, cancel := context.WithCancel(t.Context())
					want := context.Canceled
					if deadline {
						cancel()
						ctx, cancel = context.WithTimeout(t.Context(), 100*time.Millisecond)
						want = context.DeadlineExceeded
					}
					defer cancel()
					result := make(chan error, 1)
					go func() { result <- e.Export(ctx, make([]sdklog.Record, count)) }()
					waitForWrite(t, w.entered)
					if !deadline {
						cancel()
					}
					<-ctx.Done()
					select {
					case err := <-result:
						assert.ErrorIs(t, err, want)
					case <-time.After(time.Second):
						t.Error("Export did not return when its context ended")
						release()
						assert.ErrorIs(t, waitForExport(t, result), want)
					}
					release()
					require.NoError(t, e.ForceFlush(t.Context()))
					assert.Len(t, w.data(), 1)
				})
			}
		})
	}
}

func TestExporterContextEndsDuringLastWrite(t *testing.T) {
	for _, writeErr := range []error{nil, errors.New("write failed")} {
		name := "success"
		if writeErr != nil {
			name = "error"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			e, err := New(WithWriter(writerFunc(func(p []byte) (int, error) {
				cancel()
				return len(p), writeErr
			})))
			require.NoError(t, err)
			assert.ErrorIs(t, e.Export(ctx, []sdklog.Record{{}}), context.Canceled)
		})
	}
}

func TestExporterWriterPanic(t *testing.T) {
	var calls int
	e, err := New(WithWriter(writerFunc(func(p []byte) (int, error) {
		calls++
		if calls == 1 {
			panic("writer panic")
		}
		return len(p), nil
	})))
	require.NoError(t, err)
	assert.ErrorContains(t, e.Export(t.Context(), []sdklog.Record{{}}), "writer panic")
	assert.NoError(t, e.Export(t.Context(), []sdklog.Record{{}}))
	assert.Equal(t, 2, calls)
}

func TestExporterEncodingErrorRecovery(t *testing.T) {
	var buf bytes.Buffer
	e, err := New(WithWriter(&buf))
	require.NoError(t, err)
	var record sdklog.Record
	// time.Time cannot be JSON encoded when its year is outside [0, 9999].
	record.SetTimestamp(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.Error(t, e.Export(t.Context(), []sdklog.Record{record}))
	assert.Empty(t, buf.String())

	record.SetTimestamp(time.Now())
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	assert.NoError(t, e.Export(ctx, []sdklog.Record{record}))
	assert.NotEmpty(t, buf.String())
}

func TestExporterShutdownWhileWaitingForWrite(t *testing.T) {
	w, release := newBlockedWriter(t)
	e, err := New(WithWriter(w))
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	go func() { result <- e.Export(ctx, []sdklog.Record{{}}) }()
	waitForWrite(t, w.entered)
	cancel()
	require.ErrorIs(t, waitForExport(t, result), context.Canceled)

	// Wait until the next export is trying to acquire the still-blocked writer.
	waiting := &waitingContext{Context: t.Context(), waiting: make(chan struct{})}
	go func() { result <- e.Export(waiting, []sdklog.Record{{}}) }()
	waitForWrite(t, waiting.waiting)
	ctx, cancel = context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	require.ErrorIs(t, e.Shutdown(ctx), context.DeadlineExceeded)
	release()
	assert.ErrorIs(t, waitForExport(t, result), sdklog.ErrExporterShutdown)
	assert.Empty(t, w.entered, "an export waiting before shutdown must not start another write")
}

type waitingContext struct {
	context.Context
	waiting chan struct{}
	once    sync.Once
}

func (c *waitingContext) Done() <-chan struct{} {
	c.once.Do(func() { close(c.waiting) })
	return c.Context.Done()
}

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) {
	return f(p)
}

func TestExporterBlockedWriteOwnershipAndRecovery(t *testing.T) {
	w, release := newBlockedWriter(t)
	e, err := New(WithWriter(w))
	require.NoError(t, err)
	now := time.Now()
	records := []sdklog.Record{getRecord(now)}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	go func() { result <- e.Export(ctx, records) }()
	waitForWrite(t, w.entered)
	cancel()
	require.ErrorIs(t, waitForExport(t, result), context.Canceled)

	// The SDK can immediately reuse records after Export returns.
	records[0].SetBody(attribute.StringValue("reused"))
	records[0].AddAttributes(attribute.String("key", "reused"))
	records[0] = sdklog.Record{}
	for range 10 {
		ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
		err := e.Export(ctx, records)
		cancel()
		require.ErrorIs(t, err, context.DeadlineExceeded)
	}
	assert.Empty(t, w.entered, "a blocked write must prevent further writer calls")
	release()
	require.NoError(t, e.ForceFlush(t.Context()))
	require.Equal(t, []string{getJSON(&now)}, w.data())

	records[0] = getRecord(now)
	require.NoError(t, e.Export(t.Context(), records))
	assert.Equal(t, []string{getJSON(&now), getJSON(&now)}, w.data())
}

func TestExporterBlockedWriteLifecycle(t *testing.T) {
	w, release := newBlockedWriter(t)
	e, err := New(WithWriter(w))
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	go func() { result <- e.Export(ctx, []sdklog.Record{{}}) }()
	waitForWrite(t, w.entered)
	cancel()
	require.ErrorIs(t, waitForExport(t, result), context.Canceled)

	ctx, cancel = context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	assert.ErrorIs(t, e.ForceFlush(ctx), context.DeadlineExceeded)
	ctx, cancel = context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	assert.ErrorIs(t, e.Shutdown(ctx), context.DeadlineExceeded)
	assert.ErrorIs(t, e.Export(t.Context(), []sdklog.Record{{}}), sdklog.ErrExporterShutdown)
	assert.NoError(t, e.ForceFlush(t.Context()))
	assert.NoError(t, e.Shutdown(t.Context()))
	assert.False(t, w.closed, "the application owns the writer")
	release()
	waitForWrite(t, w.finished)
}

func waitForExport(t *testing.T, result <-chan error) error {
	t.Helper()
	select {
	case err := <-result:
		return err
	case <-time.After(time.Second):
		t.Fatal("Export did not return when its context ended")
		return nil
	}
}

func waitForWrite(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("writer did not reach the expected state")
	}
}

type blockedWriter struct {
	entered  chan struct{}
	release  chan struct{}
	finished chan struct{}
	mu       sync.Mutex
	writes   []string
	closed   bool
}

func newBlockedWriter(t *testing.T) (*blockedWriter, func()) {
	t.Helper()
	w := &blockedWriter{
		entered:  make(chan struct{}, 16),
		release:  make(chan struct{}),
		finished: make(chan struct{}, 16),
	}
	var once sync.Once
	release := func() { once.Do(func() { close(w.release) }) }
	t.Cleanup(release)
	return w, release
}

func (w *blockedWriter) Write(p []byte) (int, error) {
	w.entered <- struct{}{}
	<-w.release
	w.mu.Lock()
	w.writes = append(w.writes, string(p))
	w.mu.Unlock()
	w.finished <- struct{}{}
	return len(p), nil
}

func (w *blockedWriter) Close() error {
	w.closed = true
	return nil
}

func (w *blockedWriter) data() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string(nil), w.writes...)
}
