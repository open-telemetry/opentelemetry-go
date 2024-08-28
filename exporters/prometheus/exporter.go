// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	targetInfoMetricName  = "target_info"
	targetInfoDescription = "Target metadata"

	scopeInfoMetricName  = "otel_scope_info"
	scopeInfoDescription = "Instrumentation Scope metadata"

	traceIDExemplarKey = "trace_id"
	spanIDExemplarKey  = "span_id"
)

var (
	scopeInfoKeys = [2]string{"otel_scope_name", "otel_scope_version"}

	errScopeInvalid = errors.New("invalid scope")
)

// Exporter is a Prometheus Exporter that embeds the OTel metric.Reader
// interface for easy instantiation with a MeterProvider.
type Exporter struct {
	metric.Reader
}

// MarshalLog returns logging data about the Exporter.
func (e *Exporter) MarshalLog() interface{} {
	const t = "Prometheus exporter"

	if r, ok := e.Reader.(*metric.ManualReader); ok {
		under := r.MarshalLog()
		if data, ok := under.(struct {
			Type       string
			Registered bool
			Shutdown   bool
		}); ok {
			data.Type = t
			return data
		}
	}

	return struct{ Type string }{Type: t}
}

var _ metric.Reader = &Exporter{}

// keyVals is used to store resource attribute key value pairs.
type keyVals struct {
	keys []string
	vals []string
}

// collector is used to implement prometheus.Collector.
type collector struct {
	reader metric.Reader

	withoutUnits             bool
	withoutCounterSuffixes   bool
	disableScopeInfo         bool
	namespace                string
	resourceAttributesFilter attribute.Filter

	mu                sync.Mutex // mu protects all members below from the concurrent access.
	disableTargetInfo bool
	targetInfo        prometheus.Metric
	scopeInfos        map[instrumentation.Scope]prometheus.Metric
	scopeInfosInvalid map[instrumentation.Scope]struct{}
	metricFamilies    map[string]*dto.MetricFamily
	resourceKeyVals   keyVals
}

// prometheus counters MUST have a _total suffix by default:
// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/compatibility/prometheus_and_openmetrics.md
const counterSuffix = "_total"

// New returns a Prometheus Exporter.
func New(opts ...Option) (*Exporter, error) {
	cfg := newConfig(opts...)

	// this assumes that the default temporality selector will always return cumulative.
	// we only support cumulative temporality, so building our own reader enforces this.
	// TODO (#3244): Enable some way to configure the reader, but not change temporality.
	reader := metric.NewManualReader(cfg.readerOpts...)

	collector := &collector{
		reader:                   reader,
		disableTargetInfo:        cfg.disableTargetInfo,
		withoutUnits:             cfg.withoutUnits,
		withoutCounterSuffixes:   cfg.withoutCounterSuffixes,
		disableScopeInfo:         cfg.disableScopeInfo,
		scopeInfos:               make(map[instrumentation.Scope]prometheus.Metric),
		scopeInfosInvalid:        make(map[instrumentation.Scope]struct{}),
		metricFamilies:           make(map[string]*dto.MetricFamily),
		namespace:                cfg.namespace,
		resourceAttributesFilter: cfg.resourceAttributesFilter,
	}

	if err := cfg.registerer.Register(collector); err != nil {
		return nil, fmt.Errorf("cannot register the collector: %w", err)
	}

	e := &Exporter{
		Reader: reader,
	}

	return e, nil
}

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	// The Opentelemetry SDK doesn't have information on which will exist when the collector
	// is registered. By returning nothing we are an "unchecked" collector in Prometheus,
	// and assume responsibility for consistency of the metrics produced.
	//
	// See https://pkg.go.dev/github.com/prometheus/client_golang@v1.13.0/prometheus#hdr-Custom_Collectors_and_constant_Metrics
}

// Collect implements prometheus.Collector.
//
// This method is safe to call concurrently.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	// TODO (#3047): Use a sync.Pool instead of allocating metrics every Collect.
	metrics := metricdata.ResourceMetrics{}
	err := c.reader.Collect(context.TODO(), &metrics)
	if err != nil {
		if errors.Is(err, metric.ErrReaderShutdown) {
			return
		}
		otel.Handle(err)
		if errors.Is(err, metric.ErrReaderNotRegistered) {
			return
		}
	}

	global.Debug("Prometheus exporter export", "Data", metrics)

	// Initialize (once) targetInfo and disableTargetInfo.
	func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.targetInfo == nil && !c.disableTargetInfo {
			targetInfo, err := createInfoMetric(targetInfoMetricName, targetInfoDescription, metrics.Resource)
			if err != nil {
				// If the target info metric is invalid, disable sending it.
				c.disableTargetInfo = true
				otel.Handle(err)
				return
			}

			c.targetInfo = targetInfo
		}
	}()

	if !c.disableTargetInfo {
		ch <- c.targetInfo
	}

	if c.resourceAttributesFilter != nil && len(c.resourceKeyVals.keys) == 0 {
		c.createResourceAttributes(metrics.Resource)
	}

	for _, scopeMetrics := range metrics.ScopeMetrics {
		var keys, values [2]string

		if !c.disableScopeInfo {
			scopeInfo, err := c.scopeInfo(scopeMetrics.Scope)
			if errors.Is(err, errScopeInvalid) {
				// Do not report the same error multiple times.
				continue
			}
			if err != nil {
				otel.Handle(err)
				continue
			}

			ch <- scopeInfo

			keys = scopeInfoKeys
			values = [2]string{scopeMetrics.Scope.Name, scopeMetrics.Scope.Version}
		}

		for _, m := range scopeMetrics.Metrics {
			typ := c.metricType(m)
			if typ == nil {
				continue
			}
			name := c.getName(m, typ)

			drop, help := c.validateMetrics(name, m.Description, typ)
			if drop {
				continue
			}

			if help != "" {
				m.Description = help
			}

			switch v := m.Data.(type) {
			case metricdata.Histogram[int64]:
				addHistogramMetric(ch, v, m, keys, values, name, c.resourceKeyVals)
			case metricdata.Histogram[float64]:
				addHistogramMetric(ch, v, m, keys, values, name, c.resourceKeyVals)
			case metricdata.Sum[int64]:
				addSumMetric(ch, v, m, keys, values, name, c.resourceKeyVals)
			case metricdata.Sum[float64]:
				addSumMetric(ch, v, m, keys, values, name, c.resourceKeyVals)
			case metricdata.Gauge[int64]:
				addGaugeMetric(ch, v, m, keys, values, name, c.resourceKeyVals)
			case metricdata.Gauge[float64]:
				addGaugeMetric(ch, v, m, keys, values, name, c.resourceKeyVals)
			}
		}
	}
}

func addHistogramMetric[N int64 | float64](ch chan<- prometheus.Metric, histogram metricdata.Histogram[N], m metricdata.Metrics, ks, vs [2]string, name string, resourceKV keyVals) {
	for _, dp := range histogram.DataPoints {
		keys, values := getAttrs(dp.Attributes, ks, vs, resourceKV)

		desc := prometheus.NewDesc(name, m.Description, keys, nil)
		buckets := make(map[float64]uint64, len(dp.Bounds))

		cumulativeCount := uint64(0)
		for i, bound := range dp.Bounds {
			cumulativeCount += dp.BucketCounts[i]
			buckets[bound] = cumulativeCount
		}
		m, err := prometheus.NewConstHistogram(desc, dp.Count, float64(dp.Sum), buckets, values...)
		if err != nil {
			otel.Handle(err)
			continue
		}
		m = addExemplars(m, dp.Exemplars)
		ch <- m
	}
}

func addSumMetric[N int64 | float64](ch chan<- prometheus.Metric, sum metricdata.Sum[N], m metricdata.Metrics, ks, vs [2]string, name string, resourceKV keyVals) {
	valueType := prometheus.CounterValue
	if !sum.IsMonotonic {
		valueType = prometheus.GaugeValue
	}

	for _, dp := range sum.DataPoints {
		keys, values := getAttrs(dp.Attributes, ks, vs, resourceKV)

		desc := prometheus.NewDesc(name, m.Description, keys, nil)
		m, err := prometheus.NewConstMetric(desc, valueType, float64(dp.Value), values...)
		if err != nil {
			otel.Handle(err)
			continue
		}
		m = addExemplars(m, dp.Exemplars)
		ch <- m
	}
}

func addGaugeMetric[N int64 | float64](ch chan<- prometheus.Metric, gauge metricdata.Gauge[N], m metricdata.Metrics, ks, vs [2]string, name string, resourceKV keyVals) {
	for _, dp := range gauge.DataPoints {
		keys, values := getAttrs(dp.Attributes, ks, vs, resourceKV)

		desc := prometheus.NewDesc(name, m.Description, keys, nil)
		m, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, float64(dp.Value), values...)
		if err != nil {
			otel.Handle(err)
			continue
		}
		ch <- m
	}
}

// getAttrs parses the attribute.Set to two lists of matching Prometheus-style
// keys and values.
func getAttrs(attrs attribute.Set, ks, vs [2]string, resourceKV keyVals) ([]string, []string) {
	keys := make([]string, 0, attrs.Len())
	values := make([]string, 0, attrs.Len())
	itr := attrs.Iter()

	if model.NameValidationScheme == model.UTF8Validation {
		// Do not perform sanitization if prometheus supports UTF-8.
		for itr.Next() {
			kv := itr.Attribute()
			keys = append(keys, string(kv.Key))
			values = append(values, kv.Value.Emit())
		}
	} else {
		// It sanitizes invalid characters and handles duplicate keys
		// (due to sanitization) by sorting and concatenating the values following the spec.
		keysMap := make(map[string][]string)
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
		for key, vals := range keysMap {
			keys = append(keys, key)
			slices.Sort(vals)
			values = append(values, strings.Join(vals, ";"))
		}
	}

	if ks[0] != "" {
		keys = append(keys, ks[:]...)
		values = append(values, vs[:]...)
	}

	for idx := range resourceKV.keys {
		keys = append(keys, resourceKV.keys[idx])
		values = append(values, resourceKV.vals[idx])
	}

	return keys, values
}

func createInfoMetric(name, description string, res *resource.Resource) (prometheus.Metric, error) {
	keys, values := getAttrs(*res.Set(), [2]string{}, [2]string{}, keyVals{})
	desc := prometheus.NewDesc(name, description, keys, nil)
	return prometheus.NewConstMetric(desc, prometheus.GaugeValue, float64(1), values...)
}

func createScopeInfoMetric(scope instrumentation.Scope) (prometheus.Metric, error) {
	keys := scopeInfoKeys[:]
	desc := prometheus.NewDesc(scopeInfoMetricName, scopeInfoDescription, keys, nil)
	return prometheus.NewConstMetric(desc, prometheus.GaugeValue, float64(1), scope.Name, scope.Version)
}

func sanitizeRune(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ':' || r == '_' {
		return r
	}
	return '_'
}

var unitSuffixes = map[string]string{
	// Time
	"d":   "_days",
	"h":   "_hours",
	"min": "_minutes",
	"s":   "_seconds",
	"ms":  "_milliseconds",
	"us":  "_microseconds",
	"ns":  "_nanoseconds",

	// Bytes
	"By":   "_bytes",
	"KiBy": "_kibibytes",
	"MiBy": "_mebibytes",
	"GiBy": "_gibibytes",
	"TiBy": "_tibibytes",
	"KBy":  "_kilobytes",
	"MBy":  "_megabytes",
	"GBy":  "_gigabytes",
	"TBy":  "_terabytes",

	// SI
	"m": "_meters",
	"V": "_volts",
	"A": "_amperes",
	"J": "_joules",
	"W": "_watts",
	"g": "_grams",

	// Misc
	"Cel": "_celsius",
	"Hz":  "_hertz",
	"1":   "_ratio",
	"%":   "_percent",
}

// getName returns the sanitized name, prefixed with the namespace and suffixed with unit.
func (c *collector) getName(m metricdata.Metrics, typ *dto.MetricType) string {
	name := m.Name
	if model.NameValidationScheme != model.UTF8Validation {
		// Only sanitize if prometheus does not support UTF-8.
		name = sanitizeName(m.Name)
	}
	addCounterSuffix := !c.withoutCounterSuffixes && *typ == dto.MetricType_COUNTER
	if addCounterSuffix {
		// Remove the _total suffix here, as we will re-add the total suffix
		// later, and it needs to come after the unit suffix.
		name = strings.TrimSuffix(name, counterSuffix)
	}
	if c.namespace != "" {
		name = c.namespace + name
	}
	if suffix, ok := unitSuffixes[m.Unit]; ok && !c.withoutUnits && !strings.HasSuffix(name, suffix) {
		name += suffix
	}
	if addCounterSuffix {
		name += counterSuffix
	}
	return name
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
			_ = b.WriteByte(byte(replacement))
			break
		}
		b.Grow(len(n))
		_, _ = b.WriteString(n[:i])
		_ = b.WriteByte(byte(replacement))
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
			_ = b.WriteByte(byte(c))
		} else {
			_ = b.WriteByte(byte(replacement))
		}
	}

	return b.String()
}

func (c *collector) metricType(m metricdata.Metrics) *dto.MetricType {
	switch v := m.Data.(type) {
	case metricdata.Histogram[int64], metricdata.Histogram[float64]:
		return dto.MetricType_HISTOGRAM.Enum()
	case metricdata.Sum[float64]:
		if v.IsMonotonic {
			return dto.MetricType_COUNTER.Enum()
		}
		return dto.MetricType_GAUGE.Enum()
	case metricdata.Sum[int64]:
		if v.IsMonotonic {
			return dto.MetricType_COUNTER.Enum()
		}
		return dto.MetricType_GAUGE.Enum()
	case metricdata.Gauge[int64], metricdata.Gauge[float64]:
		return dto.MetricType_GAUGE.Enum()
	}
	return nil
}

func (c *collector) createResourceAttributes(res *resource.Resource) {
	c.mu.Lock()
	defer c.mu.Unlock()

	resourceAttrs, _ := res.Set().Filter(c.resourceAttributesFilter)
	resourceKeys, resourceValues := getAttrs(resourceAttrs, [2]string{}, [2]string{}, keyVals{})
	c.resourceKeyVals = keyVals{keys: resourceKeys, vals: resourceValues}
}

func (c *collector) scopeInfo(scope instrumentation.Scope) (prometheus.Metric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	scopeInfo, ok := c.scopeInfos[scope]
	if ok {
		return scopeInfo, nil
	}

	if _, ok := c.scopeInfosInvalid[scope]; ok {
		return nil, errScopeInvalid
	}

	scopeInfo, err := createScopeInfoMetric(scope)
	if err != nil {
		c.scopeInfosInvalid[scope] = struct{}{}
		return nil, fmt.Errorf("cannot create scope info metric: %w", err)
	}

	c.scopeInfos[scope] = scopeInfo

	return scopeInfo, nil
}

func (c *collector) validateMetrics(name, description string, metricType *dto.MetricType) (drop bool, help string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	emf, exist := c.metricFamilies[name]

	if !exist {
		c.metricFamilies[name] = &dto.MetricFamily{
			Name: proto.String(name),
			Help: proto.String(description),
			Type: metricType,
		}
		return false, ""
	}

	if emf.GetType() != *metricType {
		global.Error(
			errors.New("instrument type conflict"),
			"Using existing type definition.",
			"instrument", name,
			"existing", emf.GetType(),
			"dropped", *metricType,
		)
		return true, ""
	}
	if emf.GetHelp() != description {
		global.Info(
			"Instrument description conflict, using existing",
			"instrument", name,
			"existing", emf.GetHelp(),
			"dropped", description,
		)
		return false, emf.GetHelp()
	}

	return false, ""
}

func addExemplars[N int64 | float64](m prometheus.Metric, exemplars []metricdata.Exemplar[N]) prometheus.Metric {
	if len(exemplars) == 0 {
		return m
	}
	promExemplars := make([]prometheus.Exemplar, len(exemplars))
	for i, exemplar := range exemplars {
		labels := attributesToLabels(exemplar.FilteredAttributes)
		// Overwrite any existing trace ID or span ID attributes
		labels[traceIDExemplarKey] = hex.EncodeToString(exemplar.TraceID[:])
		labels[spanIDExemplarKey] = hex.EncodeToString(exemplar.SpanID[:])
		promExemplars[i] = prometheus.Exemplar{
			Value:     float64(exemplar.Value),
			Timestamp: exemplar.Time,
			Labels:    labels,
		}
	}
	metricWithExemplar, err := prometheus.NewMetricWithExemplars(m, promExemplars...)
	if err != nil {
		// If there are errors creating the metric with exemplars, just warn
		// and return the metric without exemplars.
		otel.Handle(err)
		return m
	}
	return metricWithExemplar
}

func attributesToLabels(attrs []attribute.KeyValue) prometheus.Labels {
	labels := make(map[string]string)
	for _, attr := range attrs {
		labels[string(attr.Key)] = attr.Value.Emit()
	}
	return labels
}
