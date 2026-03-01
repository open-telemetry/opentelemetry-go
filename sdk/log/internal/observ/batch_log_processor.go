// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ // import "go.opentelemetry.io/otel/sdk/log/internal/observ"

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/log/internal/x"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

const (
	// SchemaURL is the schema URL of the instrumentation.
	SchemaURL = semconv.SchemaURL
)

// ErrQueueFull is the attribute value for the "queue_full" error type.
var ErrQueueFull = otelconv.SDKProcessorLogProcessed{}.AttrErrorType("queue_full")

// BLPComponentName returns the component name attribute for a
// BatchLogProcessor with the given ID.
func BLPComponentName(id int64) attribute.KeyValue {
	t := otelconv.ComponentTypeBatchingLogProcessor
	name := fmt.Sprintf("%s/%d", t, id)
	return semconv.OTelComponentName(name)
}

// BLP is the instrumentation for an OTel SDK BatchLogProcessor.
type BLP struct {
	reg metric.Registration

	processed              metric.Int64Counter
	processedOpts          []metric.AddOption
	processedQueueFullOpts []metric.AddOption
}

// NewBLP creates a new BatchLogProcessor instrumentation.
// Returns nil if observability is not enabled.
func NewBLP(id int64, qLen func() int64, qMax int64) (*BLP, error) {
	if !x.Observability.Enabled() {
		return nil, nil
	}

	meter := otel.GetMeterProvider().Meter(
		ScopeName,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(SchemaURL),
	)

	var err error
	qCap, e := otelconv.NewSDKProcessorLogQueueCapacity(meter)
	if e != nil {
		e = fmt.Errorf("failed to create BLP queue capacity metric: %w", e)
		err = errors.Join(err, e)
	}
	qCapInst := qCap.Inst()

	qSize, e := otelconv.NewSDKProcessorLogQueueSize(meter)
	if e != nil {
		e = fmt.Errorf("failed to create BLP queue size metric: %w", e)
		err = errors.Join(err, e)
	}
	qSizeInst := qSize.Inst()

	cmpntT := semconv.OTelComponentTypeBatchingLogProcessor
	cmpnt := BLPComponentName(id)
	set := attribute.NewSet(cmpnt, cmpntT)

	// Register callback for async metrics
	obsOpts := []metric.ObserveOption{metric.WithAttributeSet(set)}
	reg, e := meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			o.ObserveInt64(qSizeInst, qLen(), obsOpts...)
			o.ObserveInt64(qCapInst, qMax, obsOpts...)
			return nil
		},
		qSizeInst,
		qCapInst,
	)
	if e != nil {
		e = fmt.Errorf("failed to register BLP queue size/capacity callback: %w", e)
		err = errors.Join(err, e)
	}

	processed, e := otelconv.NewSDKProcessorLogProcessed(meter)
	if e != nil {
		e = fmt.Errorf("failed to create BLP processed logs metric: %w", e)
		err = errors.Join(err, e)
	}

	processedOpts := []metric.AddOption{metric.WithAttributeSet(set)}
	setWithError := attribute.NewSet(cmpnt, cmpntT, ErrQueueFull)
	processedQueueFullOpts := []metric.AddOption{metric.WithAttributeSet(setWithError)}

	return &BLP{
		reg:                    reg,
		processed:              processed.Inst(),
		processedOpts:          processedOpts,
		processedQueueFullOpts: processedQueueFullOpts,
	}, err
}

func (b *BLP) Shutdown() error {
	return b.reg.Unregister()
}

func (b *BLP) Processed(ctx context.Context, n int64) {
	b.processed.Add(ctx, n, b.processedOpts...)
}

func (b *BLP) ProcessedQueueFull(ctx context.Context, n int64) {
	b.processed.Add(ctx, n, b.processedQueueFullOpts...)
}
