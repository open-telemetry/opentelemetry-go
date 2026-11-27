// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profiledata_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	"go.opentelemetry.io/otel/trace"
)

func TestLinkIsValid(t *testing.T) {
	validTraceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	validSpanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}

	assert.False(t, (profiledata.Link{}).IsValid())
	assert.False(t, (profiledata.Link{TraceID: validTraceID}).IsValid())
	assert.False(t, (profiledata.Link{SpanID: validSpanID}).IsValid())
	assert.True(t, (profiledata.Link{TraceID: validTraceID, SpanID: validSpanID}).IsValid())
}
