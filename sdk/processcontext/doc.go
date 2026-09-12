// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package processcontext implements the SDK publishing side of the OTEL_CTX
// process context mechanism (OTEP 4719).
//
// A [Publisher] serializes a [resource.Resource] into a protobuf payload and
// exposes it through a named memory-mapped region that external readers (such
// as OBI or eBPF profiler) can discover via /proc/<pid>/maps.
//
// The package is supported on Linux only. [NewPublisher] returns an error on
// other platforms.
package processcontext
