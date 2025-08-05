// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/x"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

// otelComponentType is a name identifying the type of the OpenTelemetry component.
const otelComponentType = "stdout_metric_exporter"

// exporter is an OpenTelemetry metric exporter.
type exporter struct {
	encVal atomic.Value // encoderHolder

	shutdownOnce sync.Once

	temporalitySelector sdkmetric.TemporalitySelector
	aggregationSelector sdkmetric.AggregationSelector

	redactTimestamps bool

	selfObservabilityEnabled bool
	selfObservabilityAttrs   []attribute.KeyValue // selfObservability common attributes
	spanInflightMetric       otelconv.SDKExporterSpanInflight
	spanExportedMetric       otelconv.SDKExporterSpanExported
	operationDurationMetric  otelconv.SDKExporterOperationDuration
}

// New returns a configured metric exporter.
//
// If no options are passed, the default exporter returned will use a JSON
// encoder with tab indentations that output to STDOUT.
func New(options ...Option) (sdkmetric.Exporter, error) {
	cfg := newConfig(options...)
	exp := &exporter{
		temporalitySelector: cfg.temporalitySelector,
		aggregationSelector: cfg.aggregationSelector,
		redactTimestamps:    cfg.redactTimestamps,
	}
	exp.encVal.Store(*cfg.encoder)
	exp.initSelfObservability()
	return exp, nil
}

func (e *exporter) initSelfObservability() {
	if !x.SelfObservability.Enabled() {
		return
	}

	e.selfObservabilityEnabled = true

	// common attributes
	e.selfObservabilityAttrs = []attribute.KeyValue{
		semconv.OTelComponentName(fmt.Sprintf("%s/%d", otelComponentType, newExporterID())),
		semconv.OTelComponentTypeKey.String(otelComponentType),
	}

	mp := otel.GetMeterProvider()
	m := mp.Meter("go.opentelemetry.io/otel/exporters/stdout/stdoutmetric",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err error
	if e.spanInflightMetric, err = otelconv.NewSDKExporterSpanInflight(m); err != nil {
		otel.Handle(err)
	}
	if e.spanExportedMetric, err = otelconv.NewSDKExporterSpanExported(m); err != nil {
		otel.Handle(err)
	}
	if e.operationDurationMetric, err = otelconv.NewSDKExporterOperationDuration(m); err != nil {
		otel.Handle(err)
	}
}

func (e *exporter) Temporality(k sdkmetric.InstrumentKind) metricdata.Temporality {
	return e.temporalitySelector(k)
}

func (e *exporter) Aggregation(k sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return e.aggregationSelector(k)
}

func (e *exporter) Export(ctx context.Context, data *metricdata.ResourceMetrics) (err error) {
	if e.selfObservabilityEnabled {
		e.spanInflightMetric.Add(context.Background(), 1, e.selfObservabilityAttrs...)

		defer func(starting time.Time) {
			// additional attributes for self-observability,
			// only spanExportedMetric and operationDurationMetric are supported
			addAttrs := e.selfObservabilityAttrs
			if err != nil {
				addAttrs = append(addAttrs, semconv.ErrorType(err))
			}

			e.spanInflightMetric.Add(context.Background(), -1, e.selfObservabilityAttrs...)
			e.spanExportedMetric.Add(context.Background(), 1, addAttrs...)
			e.operationDurationMetric.Record(context.Background(), time.Since(starting).Seconds(), addAttrs...)
		}(time.Now())
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	if e.redactTimestamps {
		redactTimestamps(data)
	}

	global.Debug("STDOUT exporter export", "Data", data)

	return e.encVal.Load().(encoderHolder).Encode(data)
}

func (e *exporter) ForceFlush(context.Context) error {
	// exporter holds no state, nothing to flush.
	return nil
}

func (e *exporter) Shutdown(context.Context) error {
	e.shutdownOnce.Do(func() {
		e.encVal.Store(encoderHolder{
			encoder: shutdownEncoder{},
		})
	})
	return nil
}

func (e *exporter) MarshalLog() any {
	return struct{ Type string }{Type: "STDOUT"}
}

func redactTimestamps(orig *metricdata.ResourceMetrics) {
	for i, sm := range orig.ScopeMetrics {
		metrics := sm.Metrics
		for j, m := range metrics {
			data := m.Data
			orig.ScopeMetrics[i].Metrics[j].Data = redactAggregationTimestamps(data)
		}
	}
}

var errUnknownAggType = errors.New("unknown aggregation type")

func redactAggregationTimestamps(orig metricdata.Aggregation) metricdata.Aggregation {
	switch a := orig.(type) {
	case metricdata.Sum[float64]:
		return metricdata.Sum[float64]{
			Temporality: a.Temporality,
			DataPoints:  redactDataPointTimestamps(a.DataPoints),
			IsMonotonic: a.IsMonotonic,
		}
	case metricdata.Sum[int64]:
		return metricdata.Sum[int64]{
			Temporality: a.Temporality,
			DataPoints:  redactDataPointTimestamps(a.DataPoints),
			IsMonotonic: a.IsMonotonic,
		}
	case metricdata.Gauge[float64]:
		return metricdata.Gauge[float64]{
			DataPoints: redactDataPointTimestamps(a.DataPoints),
		}
	case metricdata.Gauge[int64]:
		return metricdata.Gauge[int64]{
			DataPoints: redactDataPointTimestamps(a.DataPoints),
		}
	case metricdata.Histogram[int64]:
		return metricdata.Histogram[int64]{
			Temporality: a.Temporality,
			DataPoints:  redactHistogramTimestamps(a.DataPoints),
		}
	case metricdata.Histogram[float64]:
		return metricdata.Histogram[float64]{
			Temporality: a.Temporality,
			DataPoints:  redactHistogramTimestamps(a.DataPoints),
		}
	default:
		global.Error(errUnknownAggType, fmt.Sprintf("%T", a))
		return orig
	}
}

func redactHistogramTimestamps[T int64 | float64](
	hdp []metricdata.HistogramDataPoint[T],
) []metricdata.HistogramDataPoint[T] {
	out := make([]metricdata.HistogramDataPoint[T], len(hdp))
	for i, dp := range hdp {
		out[i] = metricdata.HistogramDataPoint[T]{
			Attributes:   dp.Attributes,
			Count:        dp.Count,
			Sum:          dp.Sum,
			Bounds:       dp.Bounds,
			BucketCounts: dp.BucketCounts,
			Min:          dp.Min,
			Max:          dp.Max,
		}
	}
	return out
}

func redactDataPointTimestamps[T int64 | float64](sdp []metricdata.DataPoint[T]) []metricdata.DataPoint[T] {
	out := make([]metricdata.DataPoint[T], len(sdp))
	for i, dp := range sdp {
		out[i] = metricdata.DataPoint[T]{
			Attributes: dp.Attributes,
			Value:      dp.Value,
		}
	}
	return out
}

var exporterIDCounter atomic.Int64

// newExporterID returns a new unique ID for an exporter.
// the starting value is 0, and it increments by 1 for each call.
func newExporterID() int64 {
	return exporterIDCounter.Add(1) - 1
}
