// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package semconv_test

import (
	"context"
	"reflect"
	"slices"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/v1.45.0/dbconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/k8sconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/rpcconv"
)

func TestNewMetricContracts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, unit, kind string
		sync, observable func(metric.Meter) (metricMetadata, error)
	}{
		{
			"k8s.container.memory.limit.utilization",
			"1",
			"Float64Gauge",
			metricConstructor(k8sconv.NewContainerMemoryLimitUtilization),
			metricConstructor(k8sconv.NewContainerMemoryLimitUtilizationObservable),
		},
		{
			"k8s.container.memory.request.utilization",
			"1",
			"Float64Gauge",
			metricConstructor(k8sconv.NewContainerMemoryRequestUtilization),
			metricConstructor(k8sconv.NewContainerMemoryRequestUtilizationObservable),
		},
		{
			"k8s.node.pid.limit",
			"{thread}",
			"Int64UpDownCounter",
			metricConstructor(k8sconv.NewNodePIDLimit),
			metricConstructor(k8sconv.NewNodePIDLimitObservable),
		},
		{
			"k8s.node.pid.usage",
			"{thread}",
			"Int64UpDownCounter",
			metricConstructor(k8sconv.NewNodePIDUsage),
			metricConstructor(k8sconv.NewNodePIDUsageObservable),
		},
		{
			"k8s.pod.cpu.limit.current",
			"{cpu}",
			"Float64UpDownCounter",
			metricConstructor(k8sconv.NewPodCPULimitCurrent),
			metricConstructor(k8sconv.NewPodCPULimitCurrentObservable),
		},
		{
			"k8s.pod.cpu.limit.desired",
			"{cpu}",
			"Float64UpDownCounter",
			metricConstructor(k8sconv.NewPodCPULimitDesired),
			metricConstructor(k8sconv.NewPodCPULimitDesiredObservable),
		},
		{
			"k8s.pod.cpu.limit.utilization",
			"1",
			"Float64Gauge",
			metricConstructor(k8sconv.NewPodCPULimitUtilization),
			metricConstructor(k8sconv.NewPodCPULimitUtilizationObservable),
		},
		{
			"k8s.pod.cpu.request.current",
			"{cpu}",
			"Float64UpDownCounter",
			metricConstructor(k8sconv.NewPodCPURequestCurrent),
			metricConstructor(k8sconv.NewPodCPURequestCurrentObservable),
		},
		{
			"k8s.pod.cpu.request.desired",
			"{cpu}",
			"Float64UpDownCounter",
			metricConstructor(k8sconv.NewPodCPURequestDesired),
			metricConstructor(k8sconv.NewPodCPURequestDesiredObservable),
		},
		{
			"k8s.pod.cpu.request.utilization",
			"1",
			"Float64Gauge",
			metricConstructor(k8sconv.NewPodCPURequestUtilization),
			metricConstructor(k8sconv.NewPodCPURequestUtilizationObservable),
		},
		{
			"k8s.pod.memory.limit.current",
			"By",
			"Int64UpDownCounter",
			metricConstructor(k8sconv.NewPodMemoryLimitCurrent),
			metricConstructor(k8sconv.NewPodMemoryLimitCurrentObservable),
		},
		{
			"k8s.pod.memory.limit.desired",
			"By",
			"Int64UpDownCounter",
			metricConstructor(k8sconv.NewPodMemoryLimitDesired),
			metricConstructor(k8sconv.NewPodMemoryLimitDesiredObservable),
		},
		{
			"k8s.pod.memory.limit.utilization",
			"1",
			"Float64Gauge",
			metricConstructor(k8sconv.NewPodMemoryLimitUtilization),
			metricConstructor(k8sconv.NewPodMemoryLimitUtilizationObservable),
		},
		{
			"k8s.pod.memory.request.current",
			"By",
			"Int64UpDownCounter",
			metricConstructor(k8sconv.NewPodMemoryRequestCurrent),
			metricConstructor(k8sconv.NewPodMemoryRequestCurrentObservable),
		},
		{
			"k8s.pod.memory.request.desired",
			"By",
			"Int64UpDownCounter",
			metricConstructor(k8sconv.NewPodMemoryRequestDesired),
			metricConstructor(k8sconv.NewPodMemoryRequestDesiredObservable),
		},
		{
			"k8s.pod.memory.request.utilization",
			"1",
			"Float64Gauge",
			metricConstructor(k8sconv.NewPodMemoryRequestUtilization),
			metricConstructor(k8sconv.NewPodMemoryRequestUtilizationObservable),
		},
		{
			"k8s.statefulset.pod.available",
			"{pod}",
			"Int64UpDownCounter",
			metricConstructor(k8sconv.NewStatefulSetPodAvailable),
			metricConstructor(k8sconv.NewStatefulSetPodAvailableObservable),
		},
		{
			"db.client.connection.limit",
			"{connection}",
			"Int64UpDownCounter",
			metricConstructor(dbconv.NewClientConnectionLimit),
			metricConstructor(dbconv.NewClientConnectionLimitObservable),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			for _, variant := range []struct {
				kind string
				new  func(metric.Meter) (metricMetadata, error)
			}{
				{test.kind, test.sync},
				{strings.Replace(test.kind, "64", "64Observable", 1), test.observable},
			} {
				t.Run(variant.kind, func(t *testing.T) {
					m := &constructorMeter{Meter: noop.NewMeterProvider().Meter("test")}
					i, err := variant.new(m)
					if err != nil {
						t.Fatal(err)
					}
					if got := reflect.TypeOf(m.config).Name(); got != variant.kind+"Config" {
						t.Errorf("instrument config = %s, want %sConfig", got, variant.kind)
					}
					if m.name != test.name || i.Name() != test.name {
						t.Errorf("meter name = %q, instrument name = %q, want %q", m.name, i.Name(), test.name)
					}
					if m.config.Unit() != test.unit || i.Unit() != test.unit {
						t.Errorf("meter unit = %q, instrument unit = %q, want %q", m.config.Unit(), i.Unit(), test.unit)
					}
				})
			}
		})
	}
}

type metricMetadata interface {
	Name() string
	Unit() string
}

func metricConstructor[O any, I metricMetadata](
	newInstrument func(metric.Meter, ...O) (I, error),
) func(metric.Meter) (metricMetadata, error) {
	return func(m metric.Meter) (metricMetadata, error) {
		return newInstrument(m)
	}
}

func TestRPCCallDurationBucketContract(t *testing.T) {
	t.Parallel()
	// RPC and gRPC deliberately recommend different bucket boundaries.
	general := []float64{
		0.0001,
		0.00025,
		0.0005,
		0.001,
		0.0025,
		0.005,
		0.01,
		0.025,
		0.05,
		0.1,
		0.25,
		0.5,
		1,
		2.5,
		5,
		10,
		30,
		60,
	}
	grpc := []float64{
		0,
		0.00001,
		0.00005,
		0.0001,
		0.0003,
		0.0006,
		0.0008,
		0.001,
		0.002,
		0.003,
		0.004,
		0.005,
		0.006,
		0.008,
		0.01,
		0.013,
		0.016,
		0.02,
		0.025,
		0.03,
		0.04,
		0.05,
		0.065,
		0.08,
		0.1,
		0.13,
		0.16,
		0.2,
		0.25,
		0.3,
		0.4,
		0.5,
		0.65,
		0.8,
		1,
		2,
		5,
		10,
		20,
		50,
		100,
	}
	tests := []struct {
		name string
		new  func(metric.Meter, ...metric.Float64HistogramOption) error
	}{
		{"client", func(m metric.Meter, options ...metric.Float64HistogramOption) error {
			_, err := rpcconv.NewClientCallDuration(m, options...)
			return err
		}},
		{"server", func(m metric.Meter, options ...metric.Float64HistogramOption) error {
			_, err := rpcconv.NewServerCallDuration(m, options...)
			return err
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			for _, protocol := range []struct {
				name    string
				want    []float64
				options []metric.Float64HistogramOption
			}{
				{"general", general, nil},
				{"grpc", grpc, []metric.Float64HistogramOption{metric.WithExplicitBucketBoundaries(grpc...)}},
			} {
				t.Run(protocol.name, func(t *testing.T) {
					m := &constructorMeter{Meter: noop.NewMeterProvider().Meter("test")}
					if err := test.new(m, protocol.options...); err != nil {
						t.Fatal(err)
					}
					config, ok := m.config.(metric.Float64HistogramConfig)
					if !ok {
						t.Fatalf("instrument config = %T, want metric.Float64HistogramConfig", m.config)
					}
					if got := config.ExplicitBucketBoundaries(); !slices.Equal(got, protocol.want) {
						t.Errorf("bucket boundaries = %v, want %v", got, protocol.want)
					}
				})
			}
		})
	}
}

func TestK8sFractionalValues(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	counter := &float64AddCapture{}
	for _, add := range []func(context.Context, float64, ...attribute.KeyValue){
		(k8sconv.PodCPULimitCurrent{Float64UpDownCounter: counter}).Add,
		(k8sconv.PodCPULimitDesired{Float64UpDownCounter: counter}).Add,
		(k8sconv.PodCPURequestCurrent{Float64UpDownCounter: counter}).Add,
		(k8sconv.PodCPURequestDesired{Float64UpDownCounter: counter}).Add,
	} {
		counter.calls = 0
		counter.value = 0
		add(ctx, 0.125)
		if counter.calls != 1 || counter.value != 0.125 {
			t.Errorf("CPU measurement = %d calls, value %v; want 1 call, value 0.125", counter.calls, counter.value)
		}
	}
	gauge := &float64RecordCapture{}
	for _, record := range []func(context.Context, float64, ...attribute.KeyValue){
		(k8sconv.ContainerMemoryLimitUtilization{Float64Gauge: gauge}).Record,
		(k8sconv.ContainerMemoryRequestUtilization{Float64Gauge: gauge}).Record,
		(k8sconv.PodCPULimitUtilization{Float64Gauge: gauge}).Record,
		(k8sconv.PodCPURequestUtilization{Float64Gauge: gauge}).Record,
		(k8sconv.PodMemoryLimitUtilization{Float64Gauge: gauge}).Record,
		(k8sconv.PodMemoryRequestUtilization{Float64Gauge: gauge}).Record,
	} {
		gauge.calls = 0
		gauge.value = 0
		record(ctx, 0.375)
		if gauge.calls != 1 || gauge.value != 0.375 {
			t.Errorf("utilization measurement = %d calls, value %v; want 1 call, value 0.375", gauge.calls, gauge.value)
		}
	}
	// Requests reserve resources and are not hard caps.
	gauge.calls = 0
	(k8sconv.PodCPURequestUtilization{Float64Gauge: gauge}).Record(ctx, 1.25)
	if gauge.calls != 1 || gauge.value != 1.25 {
		t.Errorf("CPU request measurement = %d calls, value %v; want 1 call, value 1.25", gauge.calls, gauge.value)
	}
	gauge.calls = 0
	(k8sconv.PodMemoryRequestUtilization{Float64Gauge: gauge}).Record(ctx, 1.5)
	if gauge.calls != 1 || gauge.value != 1.5 {
		t.Errorf("memory request measurement = %d calls, value %v; want 1 call, value 1.5", gauge.calls, gauge.value)
	}
}

func TestRPCCallDurationRequiredAttributes(t *testing.T) {
	t.Parallel()
	capture := &histogramAttributeCapture{}
	for _, record := range []func(context.Context, float64, rpcconv.SystemNameAttr, ...attribute.KeyValue){
		(rpcconv.ClientCallDuration{Float64Histogram: capture}).Record,
		(rpcconv.ServerCallDuration{Float64Histogram: capture}).Record,
	} {
		for _, attrs := range [][]attribute.KeyValue{nil, {attribute.String("rpc.method", "example.Service/Call")}} {
			capture.calls = 0
			record(t.Context(), 0.125, rpcconv.SystemNameGRPC, attrs...)
			want := attribute.NewSet(append(slices.Clone(attrs), attribute.String("rpc.system.name", "grpc"))...)
			if capture.calls != 1 || !capture.attrs.Equals(&want) {
				t.Errorf(
					"measurement = %d calls, attributes %v; want 1 call, attributes %v",
					capture.calls,
					capture.attrs,
					want,
				)
			}
		}
	}
}

func TestDBConnectionLimitAttributes(t *testing.T) {
	t.Parallel()
	capture := &int64AddCapture{}
	i := dbconv.ClientConnectionLimit{Int64UpDownCounter: capture}
	for _, attrs := range [][]attribute.KeyValue{
		nil,
		{i.AttrNamespace("orders"), i.AttrSystemName(dbconv.SystemNamePostgreSQL), i.AttrServerAddress("localhost"), i.AttrServerPort(5432)},
	} {
		capture.calls = 0
		i.Add(t.Context(), 7, "pool", attrs...)
		want := attribute.NewSet(attribute.String("db.client.connection.pool.name", "pool"))
		if len(attrs) != 0 {
			want = attribute.NewSet(attribute.String("db.client.connection.pool.name", "pool"),
				attribute.String("db.namespace", "orders"), attribute.String("db.system.name", "postgresql"),
				attribute.String("server.address", "localhost"), attribute.Int("server.port", 5432))
		}
		if capture.calls != 1 || capture.value != 7 || !capture.attrs.Equals(&want) {
			t.Errorf(
				"measurement = %d calls, value %d, attributes %v; want 1 call, value 7, attributes %v",
				capture.calls,
				capture.value,
				capture.attrs,
				want,
			)
		}
	}
}

type float64AddCapture struct {
	noop.Float64UpDownCounter
	calls int
	value float64
}

func (*float64AddCapture) Enabled(context.Context) bool { return true }

func (c *float64AddCapture) Add(_ context.Context, value float64, _ ...metric.AddOption) {
	c.calls++
	c.value = value
}

type float64RecordCapture struct {
	noop.Float64Gauge
	calls int
	value float64
}

func (*float64RecordCapture) Enabled(context.Context) bool { return true }

func (c *float64RecordCapture) Record(_ context.Context, value float64, _ ...metric.RecordOption) {
	c.calls++
	c.value = value
}

type histogramAttributeCapture struct {
	noop.Float64Histogram
	calls int
	attrs attribute.Set
}

func (*histogramAttributeCapture) Enabled(context.Context) bool { return true }

func (c *histogramAttributeCapture) Record(_ context.Context, _ float64, opts ...metric.RecordOption) {
	c.calls++
	c.attrs = metric.NewRecordConfig(opts).Attributes()
}

type int64AddCapture struct {
	noop.Int64UpDownCounter
	calls int
	value int64
	attrs attribute.Set
}

func (*int64AddCapture) Enabled(context.Context) bool { return true }

func (c *int64AddCapture) Add(_ context.Context, value int64, opts ...metric.AddOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewAddConfig(opts).Attributes()
}
