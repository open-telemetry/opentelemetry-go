// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package observ provides self-observability metrics for OTLP metric exporters.
// This is an experimental feature controlled by the x.Observability feature flag.
package observ // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/observ"

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/semconv/v1.40.0/otelconv"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/x"
)

var (
	attrPool = sync.Pool{
		New: func() any {
			// Pre-allocate for common attributes + dynamic error attribute
			const n = 1 /* otel.component.type */ + 1 /* otel.component.name */ + 1 /* server.address */ + 1 /* server.port */ + 1 /* error.type */
			s := make([]attribute.KeyValue, 0, n)
			return &s
		},
	}

	recOptPool = sync.Pool{
		New: func() any {
			o := make([]metric.RecordOption, 0, 1)
			return &o
		},
	}
)

// Instrumentation holds the self-observability metric instruments for an OTLP metric exporter.
type Instrumentation struct {
	exported otelconv.SDKExporterMetricDataPointExported
	inflight otelconv.SDKExporterMetricDataPointInflight
	duration otelconv.SDKExporterOperationDuration
	attrs    []attribute.KeyValue
	addOpt   metric.AddOption
	recOpt   metric.RecordOption
	enabled  bool
}

// NewInstrumentation returns instrumentation for otlpmetric grpc exporter.
// If self-observability is disabled, returns nil, nil.
func NewInstrumentation(id int64, componentType, serverAddress string, serverPort int) (*Instrumentation, error) {
	if !x.Observability.Enabled() {
		return nil, nil
	}

	em := &Instrumentation{
		enabled: true,
	}

	meter := otel.GetMeterProvider().Meter(
		"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err error
	var instrumentErr error

	em.exported, instrumentErr = otelconv.NewSDKExporterMetricDataPointExported(meter)
	if instrumentErr != nil {
		err = errors.Join(err, fmt.Errorf("failed to create exported metric: %w", instrumentErr))
	}

	em.inflight, instrumentErr = otelconv.NewSDKExporterMetricDataPointInflight(meter)
	if instrumentErr != nil {
		err = errors.Join(err, fmt.Errorf("failed to create inflight metric: %w", instrumentErr))
	}

	em.duration, instrumentErr = otelconv.NewSDKExporterOperationDuration(meter)
	if instrumentErr != nil {
		err = errors.Join(err, fmt.Errorf("failed to create duration metric: %w", instrumentErr))
	}

	// Set up common attributes
	componentName := fmt.Sprintf("%s/%d", componentType, id)
	em.attrs = []attribute.KeyValue{
		semconv.OTelComponentTypeKey.String(componentType),
		semconv.OTelComponentName(componentName),
		semconv.ServerAddress(serverAddress),
		semconv.ServerPort(serverPort),
	}

	attrSet := attribute.NewSet(em.attrs...)
	em.addOpt = metric.WithAttributeSet(attrSet)
	em.recOpt = metric.WithAttributeSet(attrSet)

	return em, err
}

// TrackExport tracks an export operation and returns a function to complete the tracking.
// The returned function should be called when the export operation completes.
func (em *Instrumentation) TrackExport(ctx context.Context, rm *metricdata.ResourceMetrics) func(error) {
	if em == nil || !em.enabled {
		return func(error) {}
	}

	var dataPointCount int64
	inflightEnabled := em.inflight.Enabled(ctx)
	exportedEnabled := em.exported.Enabled(ctx)
	durationEnabled := em.duration.Enabled(ctx)

	if inflightEnabled || exportedEnabled {
		dataPointCount = countDataPoints(rm)
	}
	var startTime time.Time
	if durationEnabled {
		startTime = time.Now()
	}

	if inflightEnabled {
		em.inflight.Inst().Add(ctx, dataPointCount, em.addOpt)
	}

	return func(err error) {
		if inflightEnabled {
			em.inflight.Inst().Add(ctx, -dataPointCount, em.addOpt)
		}

		if durationEnabled {
			duration := time.Since(startTime).Seconds()
			if err != nil {
				attrsPtr := attrPool.Get().(*[]attribute.KeyValue)
				defer func() {
					*attrsPtr = (*attrsPtr)[:0]
					attrPool.Put(attrsPtr)
				}()

				*attrsPtr = append(*attrsPtr, em.attrs...)
				*attrsPtr = append(*attrsPtr, semconv.ErrorType(err))

				recOptPtr := recOptPool.Get().(*[]metric.RecordOption)
				defer func() {
					*recOptPtr = (*recOptPtr)[:0]
					recOptPool.Put(recOptPtr)
				}()

				set := attribute.NewSet(*attrsPtr...)
				*recOptPtr = append(*recOptPtr, metric.WithAttributeSet(set))

				em.duration.Inst().Record(ctx, duration, *recOptPtr...)
			} else {
				em.duration.Inst().Record(ctx, duration, em.recOpt)
			}
		}

		if exportedEnabled {
			if err != nil {
				attrsPtr := attrPool.Get().(*[]attribute.KeyValue)
				defer func() {
					*attrsPtr = (*attrsPtr)[:0]
					attrPool.Put(attrsPtr)
				}()

				*attrsPtr = append(*attrsPtr, em.attrs...)
				*attrsPtr = append(*attrsPtr, semconv.ErrorType(err))

				set := attribute.NewSet(*attrsPtr...)
				em.exported.Inst().Add(ctx, dataPointCount, metric.WithAttributeSet(set))
			} else {
				em.exported.Inst().Add(ctx, dataPointCount, em.addOpt)
			}
		}
	}
}

// countDataPoints counts the total number of data points in a ResourceMetrics.
func countDataPoints(rm *metricdata.ResourceMetrics) int64 {
	if rm == nil {
		return 0
	}

	var total int64
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch data := m.Data.(type) {
			case metricdata.Gauge[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Gauge[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Sum[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Sum[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Histogram[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Histogram[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.ExponentialHistogram[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.ExponentialHistogram[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Summary:
				total += int64(len(data.DataPoints))
			}
		}
	}
	return total
}

// ParseEndpoint extracts server address and port from an endpoint URL.
// Returns defaults if parsing fails or endpoint is empty.
func ParseEndpoint(endpoint string) (address string, port int) {
	address = "localhost"
	port = 4317

	if endpoint == "" {
		return address, port
	}

	// Handle endpoint without scheme
	if !strings.Contains(endpoint, "://") {
		endpoint = "http://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return address, port
	}

	if u.Hostname() != "" {
		address = u.Hostname()
	}

	if u.Port() != "" {
		if p, err := strconv.Atoi(u.Port()); err == nil {
			port = p
		}
	}

	return address, port
}
