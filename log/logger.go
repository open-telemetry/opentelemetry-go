// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import (
	"context"

	"go.opentelemetry.io/otel/log/embedded"
)

// Logger TODO: comment.
type Logger interface {
	embedded.Logger

	// Emit TODO: comment.
	Emit(ctx context.Context, record Record)
}
