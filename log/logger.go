// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import (
	"context"

	"go.opentelemetry.io/otel/log/embedded"
)

// Logger emits log records.
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type Logger interface {
	// Users of the interface can ignore this. This embedded type is only used
	// by implementations of this interface. See the "API Implementations"
	// section of the package documentation for more information.
	embedded.Logger

	// Emit emits a log record.
	//
	// This method should:
	//   - be safe to call concurrently,
	//   - handle the trace context passed via ctx argument,
	//   - not modify the record's attributes,
	//   - copy the record's attributes in case of asynchronous processing,
	//   - use the current time as observed timestamp if the passed is empty.
	Emit(ctx context.Context, record Record)
}
