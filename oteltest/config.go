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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

type config struct {
	// SpanRecorder keeps track of spans.
	SpanRecorder *SpanRecorder
}

func newConfig(opts []Option) config {
	conf := config{}
	for _, opt := range opts {
		opt.apply(&conf)
	}
	return conf
}

// Option applies an option to a config.
type Option interface {
	apply(*config)
}

type spanRecorderOption struct {
	SpanRecorder *SpanRecorder
}

func (o spanRecorderOption) apply(c *config) {
	c.SpanRecorder = o.SpanRecorder
}

// WithSpanRecorder sets the SpanRecorder to use with the TracerProvider for
// testing.
func WithSpanRecorder(sr *SpanRecorder) Option {
	return spanRecorderOption{SpanRecorder: sr}
}
