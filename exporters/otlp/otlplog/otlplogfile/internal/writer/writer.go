// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package writer // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile/internal/writer"

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// FileWriter writes data to a configured file.
// It is buffered to reduce I/O operations to improve performance.
type FileWriter struct {
	path string
	file io.WriteCloser
	mu   sync.Mutex

	flushInterval time.Duration
	flushTicker   *time.Ticker
	stopTicker    chan struct{}
}

var _ flusher = (*FileWriter)(nil)

// NewFileWriter initializes a file writer for the file at the given path.
func NewFileWriter(path string, flushInterval time.Duration) (*FileWriter, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	fw := &FileWriter{
		path:          path,
		flushInterval: flushInterval,
		file:          newBufferedWriteCloser(file),
	}

	if fw.flushInterval > 0 {
		fw.startFlusher()
	}

	return fw, nil
}

// Export writes the given data in the file.
func (w *FileWriter) Export(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := w.file.Write(data); err != nil {
		return err
	}

	// As stated in the specification, line separator is \n.
	// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.36.0/specification/protocol/file-exporter.md#json-lines-file
	if _, err := io.WriteString(w.file, "\n"); err != nil {
		return err
	}

	return nil
}

// Shutdown stops the flusher. It also stops the flush ticker if set.
func (w *FileWriter) Shutdown() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.flushTicker != nil {
		close(w.stopTicker)
	}
	return w.file.Close()
}

// Flush writes buffered data to disk.
func (w *FileWriter) Flush() error {
	ff, ok := w.file.(flusher)
	if !ok {
		return nil
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	return ff.Flush()
}

// startFlusher starts the flusher to periodically flush the buffer.
func (w *FileWriter) startFlusher() {
	w.mu.Lock()
	defer w.mu.Unlock()

	ff, ok := w.file.(flusher)
	if !ok {
		return
	}

	w.stopTicker = make(chan struct{})
	w.flushTicker = time.NewTicker(w.flushInterval)
	go func() {
		for {
			select {
			case <-w.flushTicker.C:
				_ = ff.Flush()
			case <-w.stopTicker:
				w.flushTicker.Stop()
				w.flushTicker = nil
				return
			}
		}
	}()
}
