// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package reservoir // import "go.opentelemetry.io/otel/sdk/metric/internal/reservoir"

// DeferTimestamp is an interface that can be embedded in an
// exemplar.Reservoir to indicate to the SDK that it would like to take control
// over measuring the current timestamp for performance reasons. The SDK will
// provide a zero timestamp to reservoirs that embed this interface.
type DeferTimestamp interface{ deferTimestamp() }
