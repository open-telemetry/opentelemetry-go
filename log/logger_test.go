// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func TestNewLoggerConfig(t *testing.T) {
	version := "v1.1.1"
	schemaURL := "https://opentelemetry.io/schemas/1.37.0"
	attr := attribute.NewSet(
		attribute.String("user", "alice"),
		attribute.Bool("admin", true),
	)

	c := log.NewLoggerConfig(
		log.WithInstrumentationVersion(version),
		log.WithSchemaURL(schemaURL),
		log.WithInstrumentationAttributes(attr.ToSlice()...),
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
	opt := log.WithInstrumentationAttributes(attrs...)
	attrs[0] = attribute.String("service", "changed")

	c := log.NewLoggerConfig(opt)
	assert.Equal(t, want, c.InstrumentationAttributes(), "instrumentation attributes")
}

func TestWithInstrumentationAttributeSet(t *testing.T) {
	attrs := attribute.NewSet(
		attribute.String("service", "test"),
		attribute.Int("three", 3),
	)

	c := log.NewLoggerConfig(
		log.WithInstrumentationAttributeSet(attrs),
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
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributes(aliceAttr),
			log.WithInstrumentationAttributes(bobAttr),
		)
		assert.Equal(t, bob, c.InstrumentationAttributes(),
			"Later values for the same key should overwrite earlier ones.")
	})

	t.Run("DifferentKeys", func(t *testing.T) {
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributes(aliceAttr),
			log.WithInstrumentationAttributes(adminAttr),
		)
		assert.Equal(t, aliceAdmin, c.InstrumentationAttributes(),
			"Different keys should be merged.")
	})

	t.Run("Mixed", func(t *testing.T) {
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributes(aliceAttr, adminAttr),
			log.WithInstrumentationAttributes(bobAttr),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Combination of same and different keys should be merged.")
	})

	t.Run("MergedEmpty", func(t *testing.T) {
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributes(aliceAttr),
			log.WithInstrumentationAttributes(),
		)
		assert.Equal(t, alice, c.InstrumentationAttributes(),
			"Empty attributes should not affect existing ones.")
	})

	t.Run("SameKeyWithSet", func(t *testing.T) {
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributeSet(alice),
			log.WithInstrumentationAttributeSet(bob),
		)
		assert.Equal(t, bob, c.InstrumentationAttributes(),
			"Later values for the same key should overwrite earlier ones.")
	})

	t.Run("DifferentKeysWithSet", func(t *testing.T) {
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributeSet(alice),
			log.WithInstrumentationAttributeSet(attribute.NewSet(adminAttr)),
		)
		assert.Equal(t, aliceAdmin, c.InstrumentationAttributes(),
			"Different keys should be merged.")
	})

	t.Run("MixedWithSet", func(t *testing.T) {
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributeSet(aliceAdmin),
			log.WithInstrumentationAttributeSet(bob),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Combination of same and different keys should be merged.")
	})

	t.Run("MergedEmptyWithSet", func(t *testing.T) {
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributeSet(alice),
			log.WithInstrumentationAttributeSet(attribute.NewSet()),
		)
		assert.Equal(t, alice, c.InstrumentationAttributes(),
			"Empty attribute set should not affect existing ones.")
	})

	t.Run("MixedAttributesAndSet", func(t *testing.T) {
		c := log.NewLoggerConfig(
			log.WithInstrumentationAttributes(aliceAttr),
			log.WithInstrumentationAttributeSet(attribute.NewSet(bobAttr, adminAttr)),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Attributes and attribute sets should be merged together.")
	})
}

func BenchmarkNewLoggerConfig(b *testing.B) {
	for _, bb := range []struct {
		name    string
		options []log.LoggerOption
	}{
		{
			name: "with no options",
		},
		{
			name: "with an instrumentation version",
			options: []log.LoggerOption{
				log.WithInstrumentationVersion("testing version"),
			},
		},
		{
			name: "with a schema url",
			options: []log.LoggerOption{
				log.WithSchemaURL("testing URL"),
			},
		},
		{
			name: "with instrumentation attribute",
			options: []log.LoggerOption{
				log.WithInstrumentationAttributes(attribute.String("foo", "value")),
			},
		},
		{
			name: "with instrumentation attribute set",
			options: []log.LoggerOption{
				log.WithInstrumentationAttributeSet(attribute.NewSet(attribute.String("bar", "value"))),
			},
		},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				log.NewLoggerConfig(bb.options...)
			}
		})
	}
}
