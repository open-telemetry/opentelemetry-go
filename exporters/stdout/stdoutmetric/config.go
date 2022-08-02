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

//go:build go1.18
// +build go1.18

package stdoutmetric // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"encoding/json"
	"os"
)

// config contains options for the exporter.
type config struct {
	encoder *encoderHolder
}

// newConfig creates a validated config configured with options.
func newConfig(options ...Option) (config, error) {
	cfg := config{}
	for _, opt := range options {
		cfg = opt.apply(cfg)
	}

	if cfg.encoder == nil {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		cfg.encoder = &encoderHolder{encoder: enc}
	}

	return cfg, nil
}

// Option sets exporter option values.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (o optionFunc) apply(c config) config {
	return o(c)
}

// WithEncoder sets the exporter to use ecoder to encode all the metric
// data-types to an output.
func WithEncoder(encoder Encoder) Option {
	return optionFunc(func(c config) config {
		if encoder != nil {
			c.encoder = &encoderHolder{encoder: encoder}
		}
		return c
	})
}
