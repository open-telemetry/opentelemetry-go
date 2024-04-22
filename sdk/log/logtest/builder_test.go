// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
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

	assert.Equal(t, now, got.Timestamp())
	assertAttributes(t, attrs, got)
	assert.Equal(t, dropped, got.DroppedAttributes())
	assert.Equal(t, scope, got.InstrumentationScope())
	assert.Equal(t, *r, got.Resource())

	got = b.AddAttributes(log.Bool("added", true)).Record()

	assert.Equal(t, now, got.Timestamp())
	assertAttributes(t, append(attrs, log.Bool("added", true)), got)
	assert.Equal(t, dropped, got.DroppedAttributes())
	assert.Equal(t, scope, got.InstrumentationScope())
	assert.Equal(t, *r, got.Resource())
}

func assertAttributes(t *testing.T, want []log.KeyValue, r sdklog.Record) {
	t.Helper()
	var got []log.KeyValue
	r.WalkAttributes(func(kv log.KeyValue) bool {
		got = append(got, kv)
		return true
	})
	if !slices.EqualFunc(want, got, log.KeyValue.Equal) {
		t.Errorf("Attributes are not equal:\nwant: %v\ngot:  %v", want, got)
	}
}
