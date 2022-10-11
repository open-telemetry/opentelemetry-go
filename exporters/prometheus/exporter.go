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

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

import (
	"context"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/prometheus/client_golang/prometheus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Exporter is a Prometheus Exporter that embeds the OTel metric.Reader
// interface for easy instantiation with a MeterProvider.
type Exporter struct {
	metric.Reader
	Collector prometheus.Collector
}

// collector is used to implement prometheus.Collector.
type collector struct {
	metric.Reader
}

// config is added here to allow for options expansion in the future.
type config struct{}

// Option may be used in the future to apply options to a Prometheus Exporter config.
type Option interface {
	apply(config) config
}

// New returns a Prometheus Exporter.
func New(_ ...Option) Exporter {
	// this assumes that the default temporality selector will always return cumulative.
	// we only support cumulative temporality, so building our own reader enforces this.
	reader := metric.NewManualReader()
	e := Exporter{
		Reader: reader,
		Collector: &collector{
			Reader: reader,
		},
	}
	return e
}

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	metrics, err := c.Reader.Collect(context.TODO())
	if err != nil {
		otel.Handle(err)
	}
	for _, metricData := range getMetricData(metrics) {
		ch <- metricData.description
	}
}

// Collect implements prometheus.Collector.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	metrics, err := c.Reader.Collect(context.TODO())
	if err != nil {
		otel.Handle(err)
	}

	// TODO(#3166): convert otel resource to target_info
	// see https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/data-model.md#resource-attributes-1
	for _, metricData := range getMetricData(metrics) {
		if metricData.valueType == prometheus.UntypedValue {
			m, err := prometheus.NewConstHistogram(metricData.description, metricData.histogramCount, metricData.histogramSum, metricData.histogramBuckets, metricData.attributeValues...)
			if err != nil {
				otel.Handle(err)
				continue
			}
			ch <- m
		} else {
			m, err := prometheus.NewConstMetric(metricData.description, metricData.valueType, metricData.value, metricData.attributeValues...)
			if err != nil {
				otel.Handle(err)
				continue
			}
			ch <- m
		}
	}
}

// metricData holds the metadata as well as values for individual data points.
type metricData struct {
	// name should include the unit as a suffix (before _total on counters)
	// see https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/data-model.md#metric-metadata-1
	name             string
	description      *prometheus.Desc
	attributeValues  []string
	valueType        prometheus.ValueType
	value            float64
	histogramCount   uint64
	histogramSum     float64
	histogramBuckets map[float64]uint64
}

func getMetricData(metrics metricdata.ResourceMetrics) []*metricData {
	allMetrics := make([]*metricData, 0)
	for _, scopeMetrics := range metrics.ScopeMetrics {
		for _, m := range scopeMetrics.Metrics {
			switch v := m.Data.(type) {
			case metricdata.Histogram:
				allMetrics = append(allMetrics, getHistogramMetricData(v, m)...)
			case metricdata.Sum[int64]:
				allMetrics = append(allMetrics, getSumMetricData(v, m)...)
			case metricdata.Sum[float64]:
				allMetrics = append(allMetrics, getSumMetricData(v, m)...)
			case metricdata.Gauge[int64]:
				allMetrics = append(allMetrics, getGaugeMetricData(v, m)...)
			case metricdata.Gauge[float64]:
				allMetrics = append(allMetrics, getGaugeMetricData(v, m)...)
			}
		}
	}

	return allMetrics
}

func getHistogramMetricData(histogram metricdata.Histogram, m metricdata.Metrics) []*metricData {
	// TODO(https://github.com/open-telemetry/opentelemetry-go/issues/3163): support exemplars
	dataPoints := make([]*metricData, 0, len(histogram.DataPoints))
	for _, dp := range histogram.DataPoints {
		keys, values := getAttrs(dp.Attributes)
		desc := prometheus.NewDesc(sanitizeName(m.Name), m.Description, keys, nil)
		buckets := make(map[float64]uint64, len(dp.Bounds))
		for i, bound := range dp.Bounds {
			buckets[bound] = dp.BucketCounts[i]
		}
		md := &metricData{
			name:             m.Name,
			description:      desc,
			attributeValues:  values,
			valueType:        prometheus.UntypedValue,
			histogramCount:   dp.Count,
			histogramSum:     dp.Sum,
			histogramBuckets: buckets,
		}
		dataPoints = append(dataPoints, md)
	}
	return dataPoints
}

func getSumMetricData[N int64 | float64](sum metricdata.Sum[N], m metricdata.Metrics) []*metricData {
	dataPoints := make([]*metricData, 0, len(sum.DataPoints))
	for _, dp := range sum.DataPoints {
		keys, values := getAttrs(dp.Attributes)
		desc := prometheus.NewDesc(sanitizeName(m.Name), m.Description, keys, nil)
		md := &metricData{
			name:            m.Name,
			description:     desc,
			attributeValues: values,
			valueType:       prometheus.CounterValue,
			value:           float64(dp.Value),
		}
		dataPoints = append(dataPoints, md)
	}
	return dataPoints
}

func getGaugeMetricData[N int64 | float64](gauge metricdata.Gauge[N], m metricdata.Metrics) []*metricData {
	dataPoints := make([]*metricData, 0, len(gauge.DataPoints))
	for _, dp := range gauge.DataPoints {
		keys, values := getAttrs(dp.Attributes)
		desc := prometheus.NewDesc(sanitizeName(m.Name), m.Description, keys, nil)
		md := &metricData{
			name:            m.Name,
			description:     desc,
			attributeValues: values,
			valueType:       prometheus.GaugeValue,
			value:           float64(dp.Value),
		}
		dataPoints = append(dataPoints, md)
	}
	return dataPoints
}

// getAttrs parses the attribute.Set to two lists of matching Prometheus-style
// keys and values. It sanitizes invalid characters and handles duplicate keys
// (due to sanitization) by sorting and concatenating the values following the spec.
func getAttrs(attrs attribute.Set) ([]string, []string) {
	keysMap := make(map[string][]string)
	itr := attrs.Iter()
	for itr.Next() {
		kv := itr.Attribute()
		key := strings.Map(sanitizeRune, string(kv.Key))
		if _, ok := keysMap[key]; !ok {
			keysMap[key] = []string{kv.Value.Emit()}
		} else {
			// if the sanitized key is a duplicate, append to the list of keys
			keysMap[key] = append(keysMap[key], kv.Value.Emit())
		}
	}

	keys := make([]string, 0, attrs.Len())
	values := make([]string, 0, attrs.Len())
	for key, vals := range keysMap {
		keys = append(keys, key)
		sort.Slice(vals, func(i, j int) bool {
			return i < j
		})
		values = append(values, strings.Join(vals, ";"))
	}
	return keys, values
}

func sanitizeRune(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ':' || r == '_' {
		return r
	}
	return '_'
}

func sanitizeName(n string) string {
	// This algorithm is based on strings.Map from Go 1.19.
	const replacement = '_'

	valid := func(i int, r rune) bool {
		// Taken from
		// https://github.com/prometheus/common/blob/dfbc25bd00225c70aca0d94c3c4bb7744f28ace0/model/metric.go#L92-L102
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == ':' || (r >= '0' && r <= '9' && i > 0) {
			return true
		}
		return false
	}

	// This output buffer b is initialized on demand, the first time a
	// character needs to be replaced.
	var b strings.Builder
	for i, c := range n {
		if valid(i, c) {
			continue
		}

		if i == 0 && c >= '0' && c <= '9' {
			// Prefix leading number with replacement character.
			b.Grow(len(n) + 1)
			b.WriteByte(byte(replacement))
			break
		}
		b.Grow(len(n))
		b.WriteString(n[:i])
		b.WriteByte(byte(replacement))
		width := utf8.RuneLen(c)
		n = n[i+width:]
		break
	}

	// Fast path for unchanged input.
	if b.Cap() == 0 { // b.Grow was not called above.
		return n
	}

	for _, c := range n {
		// Due to inlining, it is more performant to invoke WriteByte rather then
		// WriteRune.
		if valid(1, c) { // We are guaranteed to not be at the start.
			b.WriteByte(byte(c))
		} else {
			b.WriteByte(byte(replacement))
		}
	}

	return b.String()
}
