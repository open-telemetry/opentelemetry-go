// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus // import "go.opentelemetry.io/otel/bridge/prometheus"

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	scopeName    = "go.opentelemetry.io/otel/bridge/prometheus"
	traceIDLabel = "trace_id"
	spanIDLabel  = "span_id"
)

var (
	errUnsupportedType = errors.New("unsupported metric type")
	processStartTime   = time.Now()
)

type producer struct {
	gatherers []prometheus.Gatherer
}

// NewMetricProducer returns a metric.Producer that fetches metrics from
// Prometheus. This can be used to allow Prometheus instrumentation to be
// added to an OpenTelemetry export pipeline.
func NewMetricProducer(opts ...Option) metric.Producer {
	cfg := newConfig(opts...)
	return &producer{
		gatherers: cfg.gatherers,
	}
}

func (p *producer) Produce(context.Context) ([]metricdata.ScopeMetrics, error) {
	now := time.Now()
	var errs multierr
	otelMetrics := make([]metricdata.Metrics, 0)
	for _, gatherer := range p.gatherers {
		promMetrics, err := gatherer.Gather()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		m, err := convertPrometheusMetricsInto(promMetrics, now)
		otelMetrics = append(otelMetrics, m...)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if errs.errOrNil() != nil {
		otel.Handle(errs.errOrNil())
	}
	if len(otelMetrics) == 0 {
		return nil, nil
	}
	return []metricdata.ScopeMetrics{{
		Scope: instrumentation.Scope{
			Name: scopeName,
		},
		Metrics: otelMetrics,
	}}, nil
}

func convertPrometheusMetricsInto(promMetrics []*dto.MetricFamily, now time.Time) ([]metricdata.Metrics, error) {
	var errs multierr
	otelMetrics := make([]metricdata.Metrics, 0)
	for _, pm := range promMetrics {
		newMetric := metricdata.Metrics{
			Name:        pm.GetName(),
			Description: pm.GetHelp(),
		}
		switch pm.GetType() {
		case dto.MetricType_GAUGE:
			newMetric.Data = convertGauge(pm.GetMetric(), now)
		case dto.MetricType_COUNTER:
			newMetric.Data = convertCounter(pm.GetMetric(), now)
		case dto.MetricType_HISTOGRAM:
			newMetric.Data = convertHistogram(pm.GetMetric(), now)
		default:
			// MetricType_GAUGE_HISTOGRAM, MetricType_SUMMARY, MetricType_UNTYPED
			errs = append(errs, fmt.Errorf("%w: %v for metric %v", errUnsupportedType, pm.GetType(), pm.GetName()))
			continue
		}
		otelMetrics = append(otelMetrics, newMetric)
	}
	return otelMetrics, errs.errOrNil()
}

func convertGauge(metrics []*dto.Metric, now time.Time) metricdata.Gauge[float64] {
	otelGauge := metricdata.Gauge[float64]{
		DataPoints: make([]metricdata.DataPoint[float64], len(metrics)),
	}
	for i, m := range metrics {
		dp := metricdata.DataPoint[float64]{
			Attributes: convertLabels(m.GetLabel()),
			Time:       now,
			Value:      m.GetGauge().GetValue(),
		}
		if m.GetTimestampMs() != 0 {
			dp.Time = time.UnixMilli(m.GetTimestampMs())
		}
		otelGauge.DataPoints[i] = dp
	}
	return otelGauge
}

func convertCounter(metrics []*dto.Metric, now time.Time) metricdata.Sum[float64] {
	otelCounter := metricdata.Sum[float64]{
		DataPoints:  make([]metricdata.DataPoint[float64], len(metrics)),
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
	}
	for i, m := range metrics {
		dp := metricdata.DataPoint[float64]{
			Attributes: convertLabels(m.GetLabel()),
			StartTime:  processStartTime,
			Time:       now,
			Value:      m.GetCounter().GetValue(),
			Exemplars:  []metricdata.Exemplar[float64]{convertExemplar(m.GetCounter().GetExemplar())},
		}
		createdTs := m.GetCounter().GetCreatedTimestamp()
		if createdTs.IsValid() {
			dp.StartTime = createdTs.AsTime()
		}
		if m.GetTimestampMs() != 0 {
			dp.Time = time.UnixMilli(m.GetTimestampMs())
		}
		otelCounter.DataPoints[i] = dp
	}
	return otelCounter
}

func convertHistogram(metrics []*dto.Metric, now time.Time) metricdata.Histogram[float64] {
	otelHistogram := metricdata.Histogram[float64]{
		DataPoints:  make([]metricdata.HistogramDataPoint[float64], len(metrics)),
		Temporality: metricdata.CumulativeTemporality,
	}
	for i, m := range metrics {
		bounds, bucketCounts, exemplars := convertBuckets(m.GetHistogram().GetBucket())
		dp := metricdata.HistogramDataPoint[float64]{
			Attributes:   convertLabels(m.GetLabel()),
			StartTime:    processStartTime,
			Time:         now,
			Count:        m.GetHistogram().GetSampleCount(),
			Sum:          m.GetHistogram().GetSampleSum(),
			Bounds:       bounds,
			BucketCounts: bucketCounts,
			Exemplars:    exemplars,
		}
		createdTs := m.GetHistogram().GetCreatedTimestamp()
		if createdTs.IsValid() {
			dp.StartTime = createdTs.AsTime()
		}
		if m.GetTimestampMs() != 0 {
			dp.Time = time.UnixMilli(m.GetTimestampMs())
		}
		otelHistogram.DataPoints[i] = dp
	}
	return otelHistogram
}

func convertBuckets(buckets []*dto.Bucket) ([]float64, []uint64, []metricdata.Exemplar[float64]) {
	bounds := make([]float64, len(buckets)-1)
	bucketCounts := make([]uint64, len(buckets))
	exemplars := make([]metricdata.Exemplar[float64], 0)
	for i, bucket := range buckets {
		// The last bound is the +Inf bound, which is implied in OTel, but is
		// explicit in Prometheus. Skip the last boundary, and assume it is the
		// +Inf bound.
		if i < len(bounds) {
			bounds[i] = bucket.GetUpperBound()
		}
		bucketCounts[i] = bucket.GetCumulativeCount()
		if bucket.GetExemplar() != nil {
			exemplars = append(exemplars, convertExemplar(bucket.GetExemplar()))
		}
	}
	return bounds, bucketCounts, exemplars
}

func convertLabels(labels []*dto.LabelPair) attribute.Set {
	kvs := make([]attribute.KeyValue, len(labels))
	for i, l := range labels {
		kvs[i] = attribute.String(l.GetName(), l.GetValue())
	}
	return attribute.NewSet(kvs...)
}

func convertExemplar(exemplar *dto.Exemplar) metricdata.Exemplar[float64] {
	attrs := make([]attribute.KeyValue, 0)
	var traceID, spanID []byte
	// find the trace ID and span ID in attributes, if it exists
	for _, label := range exemplar.GetLabel() {
		if label.GetName() == traceIDLabel {
			traceID = []byte(label.GetValue())
		} else if label.GetName() == spanIDLabel {
			spanID = []byte(label.GetValue())
		} else {
			attrs = append(attrs, attribute.String(label.GetName(), label.GetValue()))
		}
	}
	return metricdata.Exemplar[float64]{
		Value:              exemplar.GetValue(),
		Time:               exemplar.GetTimestamp().AsTime(),
		TraceID:            traceID,
		SpanID:             spanID,
		FilteredAttributes: attrs,
	}
}

type multierr []error

func (e multierr) errOrNil() error {
	if len(e) == 0 {
		return nil
	} else if len(e) == 1 {
		return e[0]
	}
	return e
}

func (e multierr) Error() string {
	es := make([]string, len(e))
	for i, err := range e {
		es[i] = fmt.Sprintf("* %s", err)
	}
	return strings.Join(es, "\n\t")
}
