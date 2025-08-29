// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func TestNewMeterConfig(t *testing.T) {
	version := "v1.1.1"
	schemaURL := "https://opentelemetry.io/schemas/1.0.0"
	attr := attribute.NewSet(
		attribute.String("user", "alice"),
		attribute.Bool("admin", true),
	)
	options := []metric.MeterOption{
		metric.WithInstrumentationVersion(version),
		metric.WithSchemaURL(schemaURL),
		metric.WithInstrumentationAttributes(attr.ToSlice()...),
	}

	// Ensure that options can be used concurrently.
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := metric.NewMeterConfig(options...)
			assert.Equal(t, version, c.InstrumentationVersion(), "instrumentation version")
			assert.Equal(t, schemaURL, c.SchemaURL(), "schema URL")
			assert.Equal(t, attr, c.InstrumentationAttributes(), "instrumentation attributes")
		}()
	}
	wg.Wait()
}
