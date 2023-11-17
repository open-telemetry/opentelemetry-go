// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log_test

import "testing"

// These benchmarks are based on slog/internal/benchmarks.
// They have the following desirable properties:
//
//   - They test a complete log record, from the user's call to its return.
//
//   - The benchmarked code is run concurrently in multiple goroutines, to
//     better simulate a real server (the most common environment for structured
//     logs).
//
//   - Some handlers are optimistic versions of real handlers, doing real-world
//     tasks as fast as possible (and sometimes faster, in that an
//     implementation may not be concurrency-safe). This gives us an upper bound
//     on handler performance, so we can evaluate the (handler-independent) core
//     activity of the package in an end-to-end context without concern that a
//     slow handler implementation is skewing the results.
func BenchmarkEndToEnd(b *testing.B) {
	// TODO: Replicate https://github.com/golang/go/blob/master/src/log/slog/internal/benchmarks/benchmarks_test.go
	// Run benchmarks against a "noop.Logger" and "fastTextLogger" (based on fastTextHandler)
}
