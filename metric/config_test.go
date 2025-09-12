// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func TestConfig(t *testing.T) {
	version := "v1.1.1"
	schemaURL := "https://opentelemetry.io/schemas/1.37.0"
	attr := attribute.NewSet(
		attribute.String("user", "alice"),
		attribute.Bool("admin", true),
	)

	c := metric.NewMeterConfig(
		metric.WithInstrumentationVersion(version),
		metric.WithSchemaURL(schemaURL),
		metric.WithInstrumentationAttributes(attr.ToSlice()...),
	)

	assert.Equal(t, version, c.InstrumentationVersion(), "instrumentation version")
	assert.Equal(t, schemaURL, c.SchemaURL(), "schema URL")
	assert.Equal(t, attr, c.InstrumentationAttributes(), "instrumentation attributes")
}

func TestWithInstrumentationAttributesNotLazy(t *testing.T) {
	attrs := []attribute.KeyValue{
		attribute.String("service", "test"),
		attribute.Int("three", 3),
	}
	want := attribute.NewSet(attrs...)

	// WithInstrumentationAttributes is expected to immediately
	// create an immutable set from the attributes, so later changes
	// to attrs should not affect the config.
	opt := metric.WithInstrumentationAttributes(attrs...)
	attrs[0] = attribute.String("service", "changed")

	c := metric.NewMeterConfig(opt)
	assert.Equal(t, want, c.InstrumentationAttributes(), "instrumentation attributes")
}

func TestWithInstrumentationAttributeSet(t *testing.T) {
	attrs := attribute.NewSet(
		attribute.String("service", "test"),
		attribute.Int("three", 3),
	)

	c := metric.NewMeterConfig(
		metric.WithInstrumentationAttributeSet(attrs),
	)

	assert.Equal(t, attrs, c.InstrumentationAttributes(), "instrumentation attributes")
}

func TestWithInstrumentationAttributesMerge(t *testing.T) {
	aliceAttr := attribute.String("user", "Alice")
	bobAttr := attribute.String("user", "Bob")
	adminAttr := attribute.Bool("admin", true)

	alice := attribute.NewSet(aliceAttr)
	bob := attribute.NewSet(bobAttr)
	aliceAdmin := attribute.NewSet(aliceAttr, adminAttr)
	bobAdmin := attribute.NewSet(bobAttr, adminAttr)

	t.Run("SameKey", func(t *testing.T) {
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributes(aliceAttr),
			metric.WithInstrumentationAttributes(bobAttr),
		)
		assert.Equal(t, bob, c.InstrumentationAttributes(),
			"Later values for the same key should overwrite earlier ones.")
	})

	t.Run("DifferentKeys", func(t *testing.T) {
		// Different keys should be merged
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributes(aliceAttr),
			metric.WithInstrumentationAttributes(adminAttr),
		)
		assert.Equal(t, aliceAdmin, c.InstrumentationAttributes(),
			"Different keys should be merged.")
	})

	t.Run("Mixed", func(t *testing.T) {
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributes(aliceAttr, adminAttr),
			metric.WithInstrumentationAttributes(bobAttr),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Combination of same and different keys should be merged.")
	})

	t.Run("MergedEmpty", func(t *testing.T) {
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributes(aliceAttr),
			metric.WithInstrumentationAttributes(),
		)
		assert.Equal(t, alice, c.InstrumentationAttributes(),
			"Empty attributes should not affect existing ones.")
	})

	t.Run("SameKeyWithSet", func(t *testing.T) {
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributeSet(alice),
			metric.WithInstrumentationAttributeSet(bob),
		)
		assert.Equal(t, bob, c.InstrumentationAttributes(),
			"Later values for the same key should overwrite earlier ones.")
	})

	t.Run("DifferentKeysWithSet", func(t *testing.T) {
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributeSet(alice),
			metric.WithInstrumentationAttributeSet(attribute.NewSet(adminAttr)),
		)
		assert.Equal(t, aliceAdmin, c.InstrumentationAttributes(),
			"Different keys should be merged.")
	})

	t.Run("MixedWithSet", func(t *testing.T) {
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributeSet(aliceAdmin),
			metric.WithInstrumentationAttributeSet(bob),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Combination of same and different keys should be merged.")
	})

	t.Run("MergedEmptyWithSet", func(t *testing.T) {
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributeSet(alice),
			metric.WithInstrumentationAttributeSet(attribute.NewSet()),
		)
		assert.Equal(t, alice, c.InstrumentationAttributes(),
			"Empty attribute set should not affect existing ones.")
	})

	t.Run("MixedAttributesAndSet", func(t *testing.T) {
		c := metric.NewMeterConfig(
			metric.WithInstrumentationAttributes(aliceAttr),
			metric.WithInstrumentationAttributeSet(attribute.NewSet(bobAttr, adminAttr)),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Attributes and attribute sets should be merged together.")
	})
}

func BenchmarkNewMeterConfig(b *testing.B) {
	for _, bb := range []struct {
		name    string
		options []metric.MeterOption
	}{
		{
			name: "with no options",
		},
		{
			name: "with an instrumentation version",
			options: []metric.MeterOption{
				metric.WithInstrumentationVersion("testing version"),
			},
		},
		{
			name: "with a schema url",
			options: []metric.MeterOption{
				metric.WithSchemaURL("testing URL"),
			},
		},
		{
			name: "with instrumentation attribute",
			options: []metric.MeterOption{
				metric.WithInstrumentationAttributes(attribute.String("key", "value")),
			},
		},
		{
			name: "with instrumentation attribute set",
			options: []metric.MeterOption{
				metric.WithInstrumentationAttributeSet(attribute.NewSet(attribute.String("key", "value"))),
			},
		},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				metric.NewMeterConfig(bb.options...)
			}
		})
	}
}
