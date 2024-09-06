// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package writer // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile/internal/writer"

import (
	"io"
	"sync"
	"time"
)

// Writer writes data to the configured io.WriteCloser.
// It is buffered to reduce I/O operations to improve performance.
type Writer struct {
	out io.WriteCloser
	mu  sync.Mutex

	flushInterval time.Duration
	flushTicker   *time.Ticker
	stopTicker    chan struct{}
}

var _ flusher = (*Writer)(nil)

// New initializes a writer for the given io.WriteCloser.
func New(w io.WriteCloser, flushInterval time.Duration) (*Writer, error) {
	fw := &Writer{
		flushInterval: flushInterval,
		out:           newBufferedWriteCloser(w),
	}

	if fw.flushInterval > 0 {
		fw.startFlusher()
	}

	return fw, nil
}

// Export writes the given data in the file.
func (w *Writer) Export(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := w.out.Write(data); err != nil {
		return err
	}

	// As stated in the specification, line separator is \n.
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/file-exporter.md#json-lines-file
	if _, err := io.WriteString(w.out, "\n"); err != nil {
		return err
	}

	return nil
}

// Shutdown stops the flusher. It also stops the flush ticker if set.
func (w *Writer) Shutdown() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.flushTicker != nil {
		close(w.stopTicker)
	}
	return w.out.Close()
}

// Flush writes buffered data to disk.
func (w *Writer) Flush() error {
	ff, ok := w.out.(flusher)
	if !ok {
		return nil
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	return ff.Flush()
}

// startFlusher starts the flusher to periodically flush the buffer.
func (w *Writer) startFlusher() {
	w.mu.Lock()
	defer w.mu.Unlock()

	ff, ok := w.out.(flusher)
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
