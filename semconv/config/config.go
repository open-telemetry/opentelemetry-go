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

package internal // import "go.opentelemetry.io/otel/semconv/internal"

import (
	"context"
	"net/http"
)

var defaultRequestSanitizer = func(r *http.Request) *http.Request {
	sr := r.Clone(context.Background())

	// remove any username/password info that may be in the URL
	sr.URL.User = nil

	return sr
}

// Config holds the configurable parameters for semconv methods.
type Config struct {
	RequestSanitizer func(*http.Request) *http.Request
}

// New creates a new config struct with default values set.
func New(opts ...Option) *Config {
	c := &Config{
		RequestSanitizer: defaultRequestSanitizer,
	}
	for _, opt := range opts {
		opt.apply(c)
	}

	return c
}

// Option interface used for setting optional config properties.
type Option interface {
	apply(*Config)
}

type optionFunc func(*Config)

func (o optionFunc) apply(c *Config) {
	o(c)
}

// WithRequestSanitizer specifies a custom URL sanitizer used when setting
// attributes with data coming from the HTTP request.
func WithRequestSanitizer(fn func(*http.Request) *http.Request) Option {
	return optionFunc(func(cfg *Config) {
		cfg.RequestSanitizer = fn
	})
}
