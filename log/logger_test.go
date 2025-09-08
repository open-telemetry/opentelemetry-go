// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func TestNewLoggerConfig(t *testing.T) {
	version := "v1.1.1"
	schemaURL := "https://opentelemetry.io/schemas/1.37.0"
	attr := []attribute.KeyValue{
		attribute.String("user", "alice"),
		attribute.Bool("admin", true),
	}
	attrSet := attribute.NewSet(attr...)
	options := []log.LoggerOption{
		log.WithInstrumentationVersion(version),
		log.WithSchemaURL(schemaURL),
		log.WithInstrumentationAttributes(attr...),
	}

	// Modifications to attr should not affect the config.
	attr[0] = attribute.String("user", "bob")

	// Ensure that options can be used concurrently.
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := log.NewLoggerConfig(options...)
			assert.Equal(t, version, c.InstrumentationVersion(), "instrumentation version")
			assert.Equal(t, schemaURL, c.SchemaURL(), "schema URL")
			assert.Equal(t, attrSet, c.InstrumentationAttributes(), "instrumentation attributes")
		}()
	}
	wg.Wait()
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
}
