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

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// config contains configuration for Resource creation.
type config struct {
	// detectors that will be evaluated.
	detectors []Detector
}

// Option is the interface that applies a configuration option.
type Option interface {
	// Apply sets the Option value of a config.
	Apply(*config)

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
}

type option struct{}

func (option) private() {}

// WithAttributes adds attributes to the configured Resource.
func WithAttributes(attributes ...attribute.KeyValue) Option {
	return WithDetectors(detectAttributes{attributes})
}

type detectAttributes struct {
	attributes []attribute.KeyValue
}

func (d detectAttributes) Detect(context.Context) (*Resource, error) {
	return NewWithAttributes(d.attributes...), nil
}

// WithDetectors adds detectors to be evaluated for the configured resource.
// Any use of WithDetectors disabled the default behavior, to reenable this
// inlcude a WithDetectors(BuiltinDetectors...),
// Examples:
// `New(ctx)`: Use builtin `Detector`s.
// `New(ctx, WithDetectors())`: Use no `Detector`s.

// `New(ctx, WithDetectors(d1, d2))`: Use Detector `d1`, then overlay Detector `d2`.
// ```
// New(ctx,
//      WithDetectors(BuiltinDetectors...),
//      WithDetectors(d1),
// )
// ```
// Use The `BuiltinDetectors`, then overlay Detector `d1`.
func WithDetectors(detectors ...Detector) Option {
	if len(detectors) == 0 {
		return detectorsOption{detectors: []Detector{noOp{}}}
	}
	return detectorsOption{detectors: detectors}
}

type detectorsOption struct {
	option
	detectors []Detector
}

// Apply implements Option.
func (o detectorsOption) Apply(cfg *config) {
	cfg.detectors = append(cfg.detectors, o.detectors...)
}

var BuiltinDetectors = []Detector{
	TelemetrySDK{},
	Host{},
	FromEnv{},
}

// New returns a Resource combined from the provided attributes,
// user-provided detectors or builtin detectors.
func New(ctx context.Context, opts ...Option) (*Resource, error) {
	cfg := config{}
	for _, opt := range opts {
		opt.Apply(&cfg)
	}

	if cfg.detectors == nil {
		cfg.detectors = BuiltinDetectors
	}

	return Detect(ctx, cfg.detectors...)
}
