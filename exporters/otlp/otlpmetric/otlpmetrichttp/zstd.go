// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !nozstd

package otlpmetrichttp

import (
	"bytes"
	"io"
	"sync"

	"github.com/klauspost/compress/zstd"
)

// zstdSupported reports whether this build includes the zstd codec. It is
// excluded by the nozstd build tag.
const zstdSupported = true

// zstdMaxWindowSize caps the zstd window per RFC 9659 §3's HTTP
// interoperability recommendation (see the zstd OTEP): encoders should not
// exceed an 8 MB window.
const zstdMaxWindowSize = 8 << 20 // 8 MB

var zstdPool = sync.Pool{
	New: func() any {
		w, err := zstd.NewWriter(io.Discard, zstd.WithEncoderConcurrency(1), zstd.WithWindowSize(zstdMaxWindowSize))
		if err != nil {
			panic(err)
		}
		return w
	},
}

// compressZstd compresses body as a zstd frame, mirroring the pooled
// streaming gzip.Writer pattern this package already uses for gzip.
func compressZstd(body []byte) ([]byte, error) {
	zw := zstdPool.Get().(*zstd.Encoder)
	defer func() {
		zw.Reset(io.Discard)
		zstdPool.Put(zw)
	}()

	var b bytes.Buffer
	zw.Reset(&b)

	if _, err := zw.Write(body); err != nil {
		return nil, err
	}
	// Close needs to be called to ensure body is fully written.
	if err := zw.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
