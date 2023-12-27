// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

type spyLogger struct {
	embedded.Logger
	Context context.Context
	Record  log.Record
}

func (l *spyLogger) Emit(ctx context.Context, r log.Record) {
	l.Context = ctx
	l.Record = r
}
