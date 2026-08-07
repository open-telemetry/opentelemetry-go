// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/internal/x"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/semconv/v1.43.0/otelconv"
)

// ErrAlreadyShutdown is the attribute value for the "already_shutdown" error
// type.
var ErrAlreadyShutdown = otelconv.SDKProcessorSpanProcessed{}.AttrErrorType(
	otelconv.ErrorTypeAttr("already_shutdown"),
)

// SSP is the instrumentation for an OTel SDK SimpleSpanProcessor.
type SSP struct {
	spansProcessedCounter  metric.Int64Counter
	addOpts                []metric.AddOption
	alreadyShutdownAddOpts []metric.AddOption
}

// SSPComponentName returns the component name attribute for a
// SimpleSpanProcessor with the given ID.
func SSPComponentName(id int64) attribute.KeyValue {
	t := otelconv.ComponentTypeSimpleSpanProcessor
	name := fmt.Sprintf("%s/%d", t, id)
	return semconv.OTelComponentName(name)
}

// NewSSP returns instrumentation for an OTel SDK SimpleSpanProcessor with the
// provided ID.
//
// If the experimental observability is disabled, nil is returned.
func NewSSP(id int64) (*SSP, error) {
	if !x.Observability.Enabled() {
		return nil, nil
	}

	meter := otel.GetMeterProvider().Meter(
		ScopeName,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(SchemaURL),
	)
	spansProcessedCounter, err := otelconv.NewSDKProcessorSpanProcessed(meter)
	if err != nil {
		err = fmt.Errorf("failed to create SSP processed spans metric: %w", err)
	}

	componentName := SSPComponentName(id)
	componentType := spansProcessedCounter.AttrComponentType(otelconv.ComponentTypeSimpleSpanProcessor)
	attrs := []attribute.KeyValue{componentName, componentType}
	addOpts := []metric.AddOption{metric.WithAttributeSet(attribute.NewSet(attrs...))}

	shutdownAttrs := append(attrs, ErrAlreadyShutdown)
	alreadyShutdownAddOpts := []metric.AddOption{metric.WithAttributeSet(attribute.NewSet(shutdownAttrs...))}

	return &SSP{
		spansProcessedCounter:  spansProcessedCounter.Inst(),
		addOpts:                addOpts,
		alreadyShutdownAddOpts: alreadyShutdownAddOpts,
	}, err
}

// SpanProcessed records that a span has been submitted to the exporter by the
// SimpleSpanProcessor. Per the semantic conventions, this count is recorded at
// submission time and MUST NOT be affected by the export outcome.
func (ssp *SSP) SpanProcessed(ctx context.Context) {
	ssp.spansProcessedCounter.Add(ctx, 1, ssp.addOpts...)
}

// SpanProcessedAlreadyShutdown records that a span reached the
// SimpleSpanProcessor after it had already been shut down and therefore could
// not be submitted to the exporter. Per the semantic conventions, it is counted
// with error.type set to "already_shutdown".
func (ssp *SSP) SpanProcessedAlreadyShutdown(ctx context.Context) {
	ssp.spansProcessedCounter.Add(ctx, 1, ssp.alreadyShutdownAddOpts...)
}
