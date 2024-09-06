// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package writer

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tempFile creates a temporary file for the given test case and returns its path on disk.
// The file is automatically cleaned up when the test ends.
func tempFile(tb testing.TB) *os.File {
	f, err := os.CreateTemp(tb.TempDir(), tb.Name())
	require.NoError(tb, err, "must not error when creating temporary file")
	tb.Cleanup(func() {
		assert.NoError(tb, os.RemoveAll(path.Dir(f.Name())), "must clean up files after being written")
	})
	return f
}

func TestNewFileWriter(t *testing.T) {
	f := tempFile(t)

	writer, err := New(f, 0)
	// nolint: errcheck
	defer writer.Shutdown()

	assert.NoError(t, err, "must not error when creating the file writer")

	// Ensure file was created
	_, err = os.Stat(f.Name())
	assert.NoError(t, err, "must not error when trying to retrieve file stats")
}

func TestFileWriterExport(t *testing.T) {
	f := tempFile(t)

	writer, err := New(f, 0)
	// nolint: errcheck
	defer writer.Shutdown()
	require.NoError(t, err, "must not error when creating the file writer")

	data := []byte("helloworld")
	assert.NoError(t, writer.Export(data))

	// Force data to be written to disk.
	_ = writer.Flush()

	// Read file and verify content
	content, err := os.ReadFile(f.Name())
	require.NoError(t, err, "must not error when reading file content")
	assert.Equal(t, "helloworld\n", string(content))
}

func TestFileWriterShutdown(t *testing.T) {
	f := tempFile(t)

	writer, err := New(f, 0)
	require.NoError(t, err, "must not error when creating the file writer")
	assert.NoError(t, writer.Shutdown(), "must not error when calling Shutdown()")
}

func TestFileWriterConcurrentSafe(t *testing.T) {
	f := tempFile(t)

	writer, err := New(f, 0)
	require.NoError(t, err, "must not error when creating the file writer")

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
					_ = writer.Export([]byte(fmt.Sprintf("data from goroutine %d", i)))
					_ = writer.Flush()
					atomic.AddUint64(runs, 1)
				}
			}
		}()
	}

	for atomic.LoadUint64(runs) == 0 {
		runtime.Gosched()
	}

	assert.NoError(t, writer.Shutdown(), "must not error when shutting down")
	cancel()
	wg.Wait()
}
