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

package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const scopeName = "go.opentelemetry.io/otel/bridge/opencensus"

// newTraceConfig returns a config configured with options.
func newTraceConfig(options []TraceOption) traceConfig {
	conf := traceConfig{tp: otel.GetTracerProvider()}
	for _, o := range options {
		conf = o.apply(conf)
	}
	return conf
}

type traceConfig struct {
	tp trace.TracerProvider
}

// TraceOption applies a configuration option value to an OpenCensus bridge
// Tracer.
type TraceOption interface {
	apply(traceConfig) traceConfig
}

// traceOptionFunc applies a set of options to a config.
type traceOptionFunc func(traceConfig) traceConfig

// apply returns a config with option(s) applied.
func (o traceOptionFunc) apply(conf traceConfig) traceConfig {
	return o(conf)
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
func WithTracerProvider(tp trace.TracerProvider) TraceOption {
	return traceOptionFunc(func(conf traceConfig) traceConfig {
		conf.tp = tp
		return conf
	})
}

type metricConfig struct{}

// MetricOption applies a configuration option value to an OpenCensus bridge
// MetricProducer.
type MetricOption interface {
	apply(metricConfig) metricConfig
}
