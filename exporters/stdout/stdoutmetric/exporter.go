// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/selfobservability"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/x"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

var (
	stdoutMetricExporterComponentType = otelconv.ComponentTypeAttr("stdout_metric_exporter")
	exporterIDCounter                 atomic.Int64
)

// nextExporterID returns an identifier for this stdoutmetric exporter,
// starting with 0 and incrementing by 1 each time it is called.
func nextExporterID() int64 {
	return exporterIDCounter.Add(1) - 1
}

// exporter is an OpenTelemetry metric exporter.
type exporter struct {
	encVal atomic.Value // encoderHolder

	shutdownOnce sync.Once

	temporalitySelector metric.TemporalitySelector
	aggregationSelector metric.AggregationSelector

	redactTimestamps bool

	selfObservabilityEnabled bool
	exporterMetric           *selfobservability.ExporterMetrics
}

// New returns a configured metric exporter.
//
// If no options are passed, the default exporter returned will use a JSON
// encoder with tab indentations that output to STDOUT.
func New(options ...Option) (metric.Exporter, error) {
	cfg := newConfig(options...)
	exp := &exporter{
		temporalitySelector: cfg.temporalitySelector,
		aggregationSelector: cfg.aggregationSelector,
		redactTimestamps:    cfg.redactTimestamps,
	}
	exp.initSelfObservability()
	exp.encVal.Store(*cfg.encoder)
	return exp, nil
}

func (e *exporter) initSelfObservability() {
	if !x.SelfObservability.Enabled() {
		return
	}

	e.selfObservabilityEnabled = true
	componentName := fmt.Sprintf("%s/%d", stdoutMetricExporterComponentType, nextExporterID())
	e.exporterMetric = selfobservability.NewExporterMetrics(
		"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric",
		semconv.OTelComponentName(componentName),
		semconv.OTelComponentTypeKey.String(string(stdoutMetricExporterComponentType)),
	)
}

func (e *exporter) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return e.temporalitySelector(k)
}

func (e *exporter) Aggregation(k metric.InstrumentKind) metric.Aggregation {
	return e.aggregationSelector(k)
}

func (e *exporter) Export(ctx context.Context, data *metricdata.ResourceMetrics) error {
	trackExportFunc := e.trackExport(context.Background(), countDataPoints(data))
	if err := ctx.Err(); err != nil {
		trackExportFunc(err)
		return err
	}
	if e.redactTimestamps {
		redactTimestamps(data)
	}

	global.Debug("STDOUT exporter export", "Data", data)

	err := e.encVal.Load().(encoderHolder).Encode(data)
	if err != nil {
		trackExportFunc(err)
		return err
	}
	trackExportFunc(nil)
	return nil
}

func (e *exporter) trackExport(ctx context.Context, count int64) func(err error) {
	if !e.selfObservabilityEnabled {
		return func(error) {}
	}
	return e.exporterMetric.TrackExport(ctx, count)
}

func (*exporter) ForceFlush(context.Context) error {
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

func (*exporter) MarshalLog() any {
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
