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

package pull

import (
	"time"

	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Config contains configuration for a push Controller.
type Config struct {
	// ErrorHandler is the function called when the Controller encounters an error.
	//
	// This option can be overridden after instantiation of the Controller
	// with the `SetErrorHandler` method.
	ErrorHandler sdk.ErrorHandler

	// Resource is the OpenTelemetry resource associated with all Meters
	// created by the Controller.
	Resource *resource.Resource

	// Stateful causes the controller to maintain state across
	// collection events, so that records in the exported
	// checkpoint set are cumulative.
	Stateful bool

	// CachePeriod is the period which a recently-computed result
	// will be returned without gathering metric data again.
	//
	// If the period is zero, caching of the result is disabled.
	// The default value is 10 seconds.
	CachePeriod time.Duration
}

// Option is the interface that applies the value to a configuration option.
type Option interface {
	// Apply sets the Option value of a Config.
	Apply(*Config)
}

// WithErrorHandler sets the ErrorHandler configuration option of a Config.
func WithErrorHandler(fn sdk.ErrorHandler) Option {
	return errorHandlerOption(fn)
}

type errorHandlerOption sdk.ErrorHandler

func (o errorHandlerOption) Apply(config *Config) {
	config.ErrorHandler = sdk.ErrorHandler(o)
}

// WithResource sets the Resource configuration option of a Config.
func WithResource(r *resource.Resource) Option {
	return resourceOption{r}
}

type resourceOption struct{ *resource.Resource }

func (o resourceOption) Apply(config *Config) {
	config.Resource = o.Resource
}

// WithStateful sets the Stateful configuration option of a Config.
func WithStateful(stateful bool) Option {
	return statefulOption(stateful)
}

type statefulOption bool

func (o statefulOption) Apply(config *Config) {
	config.Stateful = bool(o)
}

// WithCachePeriod sets the CachePeriod configuration option of a Config.
func WithCachePeriod(cachePeriod time.Duration) Option {
	return cachePeriodOption(cachePeriod)
}

type cachePeriodOption time.Duration

func (o cachePeriodOption) Apply(config *Config) {
	config.CachePeriod = time.Duration(o)
}
