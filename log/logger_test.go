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
	schemaURL := "https://opentelemetry.io/schemas/1.0.0"
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
