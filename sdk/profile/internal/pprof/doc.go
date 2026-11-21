// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package pprof decodes Go runtime/pprof profiles, computes interval deltas for
// cumulative profiles, enriches binary mappings with OpenTelemetry build IDs,
// and converts pprof profiles into [go.opentelemetry.io/otel/sdk/profile/profiledata]
// structures.
package pprof
