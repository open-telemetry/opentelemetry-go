// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package transform provides transformation functionality from the
// sdk/log data-types into OTLP data-types.
package transform // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp/internal/transform"

import (
	lpb "go.opentelemetry.io/proto/otlp/logs/v1"

	"go.opentelemetry.io/otel/sdk/log"
)

func ResourceLogs([]log.Record) (*lpb.ResourceLogs, error) {
	// TODO: implement
	return nil, nil
}
