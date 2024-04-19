// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestRecordBuilder(t *testing.T) {
	now := time.Now()
	attrs := []log.KeyValue{
		log.Int("int", 1),
		log.String("str", "foo"),
		log.Float64("flt", 3.14),
	}
	dropped := 3
	scope := instrumentation.Scope{
		Name: t.Name(),
	}
	r := resource.NewSchemaless(attribute.Bool("works", true))

	b := RecordBuilder{}.
		SetTimestamp(now).
		SetAttributes(attrs...).
		SetDroppedAttributes(dropped).
		SetInstrumentationScope(scope).
		SetResource(r)

	got := b.Record()
	var gotAttrs []log.KeyValue
	got.WalkAttributes(func(kv log.KeyValue) bool {
		gotAttrs = append(gotAttrs, kv)
		return true
	})

	assert.Equal(t, now, got.Timestamp())
	assert.Equal(t, attrs, gotAttrs)
	assert.Equal(t, dropped, got.DroppedAttributes())
	assert.Equal(t, scope, got.InstrumentationScope())
	assert.Equal(t, *r, got.Resource())

	got = b.AddAttributes(log.Bool("added", true)).Record()
	gotAttrs = nil
	got.WalkAttributes(func(kv log.KeyValue) bool {
		gotAttrs = append(gotAttrs, kv)
		return true
	})

	assert.Equal(t, now, got.Timestamp())
	assert.Equal(t, append(attrs, log.Bool("added", true)), gotAttrs)
	assert.Equal(t, dropped, got.DroppedAttributes())
	assert.Equal(t, scope, got.InstrumentationScope())
	assert.Equal(t, *r, got.Resource())
}
