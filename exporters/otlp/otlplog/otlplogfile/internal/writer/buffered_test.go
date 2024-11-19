// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package writer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	msg = "hello, world!"

	sizeByte     = 1
	sizeKiloByte = 1 << (10 * iota)
	sizeMegaByte
)

type noopWriteCloser struct {
	w io.Writer
}

func (wc *noopWriteCloser) Write(p []byte) (int, error) { return wc.w.Write(p) }
func (wc *noopWriteCloser) Close() error                { return nil }

func TestBufferedWrites(t *testing.T) {
	t.Parallel()

	b := bytes.NewBuffer(nil)
	w := newBufferedWriteCloser(&noopWriteCloser{b})

	_, err := w.Write([]byte(msg))
	assert.NoError(t, err, "Must not error when writing data")
	assert.NoError(t, w.Close(), "Must not error when closing writer")
	assert.Equal(t, msg, b.String(), "Must match the expected string")
}

var errBenchmark error

func BenchmarkBufferedWriter(b *testing.B) {
	for _, payloadSize := range []int{
		10 * sizeKiloByte,
		100 * sizeKiloByte,
		sizeMegaByte,
		10 * sizeMegaByte,
	} {
		payload := make([]byte, payloadSize)
		for i := 0; i < payloadSize; i++ {
			payload[i] = 'a'
		}

		for name, w := range map[string]io.WriteCloser{
			"raw-file":      tempFile(b),
			"buffered-file": newBufferedWriteCloser(tempFile(b)),
		} {
			w := w
			b.Run(fmt.Sprintf("%s_%d_bytes", name, payloadSize), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				var err error
				for i := 0; i < b.N; i++ {
					_, err = w.Write(payload)
				}

				errBenchmark = errors.Join(err, w.Close())
			})
		}
	}
}
