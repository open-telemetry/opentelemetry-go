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
