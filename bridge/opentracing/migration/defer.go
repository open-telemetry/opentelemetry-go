// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package migration // import "go.opentelemetry.io/otel/bridge/opentracing/migration"

import (
	"context"
)

type doDeferredContextSetupType struct{}

var (
	doDeferredContextSetupTypeKey   = doDeferredContextSetupType{}
	doDeferredContextSetupTypeValue = doDeferredContextSetupType{}
)

// WithDeferredSetup returns a context that can tell the OpenTelemetry
// tracer to skip the context setup in the Start() function.
func WithDeferredSetup(ctx context.Context) context.Context {
	return context.WithValue(ctx, doDeferredContextSetupTypeKey, doDeferredContextSetupTypeValue)
}

// SkipContextSetup can tell the OpenTelemetry tracer to skip the
// context setup during the span creation in the Start() function.
func SkipContextSetup(ctx context.Context) bool {
	_, ok := ctx.Value(doDeferredContextSetupTypeKey).(doDeferredContextSetupType)
	return ok
}
