// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

func TestInstrumentationScope(t *testing.T) {
	want := &commonpb.InstrumentationScope{
		Name:    "name",
		Version: "1.0.0",
		Attributes: []*commonpb.KeyValue{
			{
				Key: "foo",
				Value: &commonpb.AnyValue{
					Value: &commonpb.AnyValue_StringValue{StringValue: "bar"},
				},
			},
		},
	}

	in := instrumentation.Scope{
		Name:       "name",
		Version:    "1.0.0",
		Attributes: attribute.NewSet(attribute.String("foo", "bar")),
	}

	got := InstrumentationScope(in)

	assert.Equal(t, want, got)
}
