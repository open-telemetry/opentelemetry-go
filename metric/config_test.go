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
