// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlphttp_test

import (
	"testing"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

func TestHTTPExporters(t *testing.T) {
	metrics, err := otlpmetrichttp.New(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if err := metrics.Shutdown(t.Context()); err != nil {
		t.Fatal(err)
	}
	logs, err := otlploghttp.New(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if err := logs.Shutdown(t.Context()); err != nil {
		t.Fatal(err)
	}
}
