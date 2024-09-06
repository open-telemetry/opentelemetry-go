// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlplogfile // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile"
import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"

	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// tempFile creates a temporary file for the given test case and returns its path on disk.
// The file is automatically cleaned up when the test ends.
func tempFile(tb testing.TB) *os.File {
	f, err := os.CreateTemp(tb.TempDir(), tb.Name())
	assert.NoError(tb, err, "must not error when creating temporary file")
	tb.Cleanup(func() {
		assert.NoError(tb, os.RemoveAll(path.Dir(f.Name())), "must clean up files after being written")
	})
	return f
}

// makeRecords is a helper function to generate an array of log record with the desired size.
func makeRecords(count int, message string) []sdklog.Record {
	var records []sdklog.Record
	for i := 0; i < count; i++ {
		r := sdklog.Record{}
		r.SetSeverityText("INFO")
		r.SetSeverity(log.SeverityInfo)
		r.SetBody(log.StringValue(message))
		r.SetTimestamp(time.Now())
		r.SetObservedTimestamp(time.Now())
		records = append(records, r)
	}
	return records
}

func TestExporter(t *testing.T) {
	file := tempFile(t)
	records := makeRecords(1, "hello, world!")

	exporter, err := New(WithWriter(file))
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, exporter.Shutdown(context.TODO()))
	})

	err = exporter.Export(context.TODO(), records)
	assert.NoError(t, err)
	err = exporter.ForceFlush(context.TODO())
	assert.NoError(t, err)
}

func TestExporterConcurrentSafe(t *testing.T) {
	file := tempFile(t)
	exporter, err := New(WithWriter(file))
	require.NoError(t, err, "New()")

	const goroutines = 10

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	runs := new(uint64)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					_ = exporter.Export(ctx, makeRecords(1, fmt.Sprintf("log from goroutine %d", i)))
					_ = exporter.ForceFlush(ctx)
					atomic.AddUint64(runs, 1)
				}
			}
		}()
	}

	for atomic.LoadUint64(runs) == 0 {
		runtime.Gosched()
	}

	assert.NoError(t, exporter.Shutdown(ctx), "must not error when shutting down")
	cancel()
	wg.Wait()
}
