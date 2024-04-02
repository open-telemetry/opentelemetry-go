// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
)

func TestNewBatchingConfig(t *testing.T) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		t.Log(err)
	}))

	testcases := []struct {
		name    string
		envars  map[string]string
		options []BatchingOption
		want    batchingConfig
	}{
		{
			name: "Defaults",
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
			},
		},
		{
			name: "Options",
			options: []BatchingOption{
				WithMaxQueueSize(1),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
			},
			want: batchingConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
			},
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(1),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(10),
			},
			want: batchingConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(100 * time.Millisecond),
				expTimeout:      newSetting(1000 * time.Millisecond),
				expMaxBatchSize: newSetting(10),
			},
		},
		{
			name: "InvalidOptions",
			options: []BatchingOption{
				WithMaxQueueSize(-11),
				WithExportInterval(-1 * time.Microsecond),
				WithExportTimeout(-1 * time.Hour),
				WithExportMaxBatchSize(-2),
			},
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
			},
		},
		{
			name: "InvalidEnvironment",
			envars: map[string]string{
				envarMaxQSize:        "-1",
				envarExpInterval:     "-1",
				envarExpTimeout:      "-1",
				envarExpMaxBatchSize: "-1",
			},
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
			},
		},
		{
			name: "Precedence",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(1),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(10),
			},
			options: []BatchingOption{
				// These override the environment variables.
				WithMaxQueueSize(3),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
			},
			want: batchingConfig{
				maxQSize:        newSetting(3),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, newBatchingConfig(tc.options))
		})
	}
}

func TestQueue(t *testing.T) {
	var r Record
	r.SetBody(log.BoolValue(true))

	t.Run("newQueue", func(t *testing.T) {
		const size = 1
		q := newQueue(size)
		assert.Equal(t, q.len, 0)
		assert.Equal(t, size, q.cap, "capacity")
		assert.Equal(t, size, q.read.Len(), "read ring")
		assert.Same(t, q.read, q.write, "different rings")
	})

	t.Run("Enqueue", func(t *testing.T) {
		const size = 2
		q := newQueue(size)

		var notR Record
		notR.SetBody(log.IntValue(10))

		assert.Equal(t, 1, q.Enqueue(notR), "incomplete batch")
		assert.Equal(t, 1, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, q.Enqueue(r), "complete batch")
		assert.Equal(t, 2, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, q.Enqueue(r), "overflow batch")
		assert.Equal(t, 2, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, []Record{r, r}, q.Flush(), "flushed Records")
	})

	t.Run("Flush", func(t *testing.T) {
		const size = 2
		q := newQueue(size)
		q.write.Value = r
		q.write = q.write.Next()
		q.len = 1

		assert.Equal(t, []Record{r}, q.Flush(), "flushed")
	})

	t.Run("TryFlush", func(t *testing.T) {
		const size = 3
		q := newQueue(size)
		for i := 0; i < size-1; i++ {
			q.write.Value = r
			q.write = q.write.Next()
			q.len++
		}

		buf := make([]Record, 1)
		f := func([]Record) bool { return false }
		assert.Equal(t, size-1, q.TryDequeue(buf, f), "not flushed")
		require.Equal(t, size-1, q.len, "length")
		require.NotSame(t, q.read, q.write, "read ring advanced")

		var flushed []Record
		f = func(r []Record) bool {
			flushed = append(flushed, r...)
			return true
		}
		if assert.Equal(t, size-2, q.TryDequeue(buf, f), "did not flush len(buf)") {
			assert.Equal(t, []Record{r}, flushed, "Records")
		}

		buf = slices.Grow(buf, size)
		flushed = flushed[:0]
		if assert.Equal(t, 0, q.TryDequeue(buf, f), "did not flush len(queue)") {
			assert.Equal(t, []Record{r}, flushed, "Records")
		}
	})

	t.Run("ConcurrentSafe", func(t *testing.T) {
		const goRoutines = 10

		flushed := make(chan []Record, goRoutines)
		out := make([]Record, 0, goRoutines)
		done := make(chan struct{})
		go func() {
			defer close(done)
			for recs := range flushed {
				out = append(out, recs...)
			}
		}()

		var wg sync.WaitGroup
		wg.Add(goRoutines)

		b := newQueue(goRoutines)
		for i := 0; i < goRoutines; i++ {
			go func() {
				defer wg.Done()
				b.Enqueue(Record{})
				flushed <- b.Flush()
			}()
		}

		wg.Wait()
		close(flushed)
		<-done

		assert.Len(t, out, goRoutines, "flushed Records")
	})
}
