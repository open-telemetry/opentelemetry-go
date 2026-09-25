// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build nozstd

package otlpmetrichttp

// zstdSupported reports whether this build includes the zstd codec. It is
// excluded by the nozstd build tag.
const zstdSupported = false

// compressZstd is unreachable: callers must check zstdSupported first.
func compressZstd([]byte) ([]byte, error) {
	panic("otlpmetrichttp: compressZstd called in a nozstd build")
}
