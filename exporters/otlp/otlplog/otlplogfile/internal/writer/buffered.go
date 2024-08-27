// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package writer // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile/internal/writer"

import (
	"bufio"
	"errors"
	"io"
)

// flusher implementations are responsible for ensuring that any buffered or pending data
// is written out or processed.
type flusher interface {
	// Flush writes any buffered data to the underlying storage.
	// It returns an error if the data could not be flushed.
	Flush() error
}

// bufferedWriter is intended to use more memory
// in order to optimize writing to disk to help improve performance.
type bufferedWriter struct {
	wrapped io.Closer
	buffer  *bufio.Writer
}

// Ensure that the implementation satisfies the interface at compile-time.
var (
	_ io.WriteCloser = (*bufferedWriter)(nil)
	_ flusher        = (*bufferedWriter)(nil)
)

func newBufferedWriteCloser(f io.WriteCloser) io.WriteCloser {
	return &bufferedWriter{
		wrapped: f,
		buffer:  bufio.NewWriter(f),
	}
}

func (bw *bufferedWriter) Write(p []byte) (n int, err error) {
	return bw.buffer.Write(p)
}

func (bw *bufferedWriter) Close() error {
	return errors.Join(
		bw.buffer.Flush(),
		bw.wrapped.Close(),
	)
}

func (bw *bufferedWriter) Flush() error {
	return bw.buffer.Flush()
}
