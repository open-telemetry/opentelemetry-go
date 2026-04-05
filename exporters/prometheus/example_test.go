// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/model"
	"github.com/prometheus/otlptranslator"

	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func Example() {
	// Create a new Prometheus exporter. It automatically registers with the DefaultRegisterer.
	exporter, err := otelprom.New()
	if err != nil {
		log.Fatal(err)
	}

	// Register the exporter with a new MeterProvider.
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(
		"example-global",
		otelmetric.WithInstrumentationVersion("v1.0.0"),
		otelmetric.WithSchemaURL("https://opentelemetry.io/schemas/v1.0.0"),
	)

	// Create a counter instrument.
	counter, err := meter.Float64Counter("bar", otelmetric.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	counter.Add(ctx, 10)

	// Serve metrics using promhttp.Handler().
	// In production, you would use http.ListenAndServe(":8080", promhttp.Handler()).
	// For this testable example, we use httptest.NewServer.
	server := httptest.NewServer(promhttp.Handler())
	defer server.Close()

	// Make an HTTP request to the endpoint.
	resp, err := http.Get(server.URL) //nolint:noctx
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// The response contains all metrics registered with the default registry,
	// including Go runtime metrics. To make this example testable and deterministic,
	// we filter the output for lines containing our metric name.
	lines := strings.SplitSeq(string(body), "\n")
	for line := range lines {
		if strings.Contains(line, "bar") {
			fmt.Println(line)
		}
	}

	// Output:
	// # HELP bar_total a simple counter
	// # TYPE bar_total counter
	// bar_total{otel_scope_name="example-global",otel_scope_schema_url="https://opentelemetry.io/schemas/v1.0.0",otel_scope_version="v1.0.0"} 10
}

func Example_customRegistry() {
	// Create a custom Prometheus registry. This is often used to avoid global state.
	reg := prometheus.NewRegistry()

	// Create a new Prometheus exporter using the custom registry.
	exporter, err := otelprom.New(otelprom.WithRegisterer(reg))
	if err != nil {
		log.Fatal(err)
	}

	// Register the exporter with a new MeterProvider.
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(
		"example-custom",
		otelmetric.WithInstrumentationVersion("v1.0.0"),
		otelmetric.WithSchemaURL("https://opentelemetry.io/schemas/v1.0.0"),
	)

	// Create a counter instrument.
	counter, err := meter.Float64Counter("foo", otelmetric.WithDescription("another counter"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	counter.Add(ctx, 5)

	// Serve metrics using promhttp.HandlerFor.
	server := httptest.NewServer(promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	defer server.Close()

	// Make an HTTP request to the endpoint.
	resp, err := http.Get(server.URL) //nolint:noctx
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Filter the output for lines containing our metric name.
	lines := strings.SplitSeq(string(body), "\n")
	for line := range lines {
		if strings.Contains(line, "foo") {
			fmt.Println(line)
		}
	}

	// Output:
	// # HELP foo_total another counter
	// # TYPE foo_total counter
	// foo_total{otel_scope_name="example-custom",otel_scope_schema_url="https://opentelemetry.io/schemas/v1.0.0",otel_scope_version="v1.0.0"} 5
}

func Example_noTranslation() {
	// Set NameEscapingScheme to NoEscaping to prevent the prometheus client from escaping to underscores.
	model.NameEscapingScheme = model.NoEscaping

	// Create a new Prometheus exporter using NoTranslation strategy.
	// This keeps the original OpenTelemetry metric names.
	// It uses the global registry by default.
	exporter, err := otelprom.New(
		otelprom.WithTranslationStrategy(otlptranslator.NoTranslation),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register the exporter with a new MeterProvider.
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(
		"example-no-translation",
		otelmetric.WithInstrumentationVersion("v1.0.0"),
		otelmetric.WithSchemaURL("https://opentelemetry.io/schemas/v1.0.0"),
	)

	// Create a counter instrument.
	// We use a name with a dot ("my.metric").
	// With NoTranslation strategy, suffixes like _total are not added.
	counter, err := meter.Float64Counter("my.metric", otelmetric.WithDescription("a counter without translation"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	counter.Add(ctx, 5)

	// Serve metrics using promhttp.Handler (uses default gatherer).
	server := httptest.NewServer(promhttp.Handler())
	defer server.Close()

	// Make an HTTP request to the endpoint.
	resp, err := http.Get(server.URL) //nolint:noctx
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Filter the output for lines containing our metric name.
	lines := strings.SplitSeq(string(body), "\n")
	for line := range lines {
		if strings.Contains(line, "my.metric") {
			fmt.Println(line)
		}
	}

	// Output:
	// # HELP "my.metric" a counter without translation
	// # TYPE "my.metric" counter
	// {"my.metric",otel_scope_name="example-no-translation",otel_scope_schema_url="https://opentelemetry.io/schemas/v1.0.0",otel_scope_version="v1.0.0"} 5
}
