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

package basic // import "go.opentelemetry.io/otel/sdk/metric/controller/basic"

import (
	"time"

	"go.opentelemetry.io/otel/sdk/resource"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

// Config contains configuration for a basic Controller.
type Config struct {
	// Resource is the OpenTelemetry resource associated with all Meters
	// created by the Controller.
	Resource *resource.Resource

	// CollectPeriod is the interval between calls to Collect a
	// checkpoint.
	//
	// When pulling metrics and not exporting, this is the minimum
	// time between calls to Collect.  In a pull-only
	// configuration, collection is performed on demand; set
	// CollectPeriod to 0 always recompute the export record set.
	//
	// When exporting metrics, this must be > 0.
	//
	// Default value is 10s.
	CollectPeriod time.Duration

	// CollectTimeout is the timeout of the Context passed to
	// Collect() and subsequently to Observer instrument callbacks.
	// 
	// Defaults to the configured collection period (see CollectPeriod).
	CollectTimeout time.Duration

	// Exporter is used for pushing metric data.
	Exporter export.Exporter

	// ExportTimeout is the timeout of the Context when an Exporter is configured.
	//
	// Defaults to the configured collection period (see CollectPeriod).
	ExportTimeout time.Duration
}

// Option is the interface that applies the value to a configuration option.
type Option interface {
	// Apply sets the Option value of a Config.
	Apply(*Config)
}

// WithResource sets the Resource configuration option of a Config.
func WithResource(r *resource.Resource) Option {
	return resourceOption{r}
}

type resourceOption struct{ *resource.Resource }

func (o resourceOption) Apply(config *Config) {
	config.Resource = o.Resource
}

// WithCollectPeriod sets the CollectPeriod configuration option of a Config.
func WithCollectPeriod(period time.Duration) Option {
	return collectPeriodOption(period)
}

type collectPeriodOption time.Duration

func (o collectPeriodOption) Apply(config *Config) {
	config.CollectPeriod = time.Duration(o)
}

// WithCollectTimeout sets the CollectTimeout configuration option of a Config.
func WithCollectTimeout(timeout time.Duration) Option {
	return collectTimeoutOption(timeout)
}

type collectTimeoutOption time.Duration

func (o collectTimeoutOption) Apply(config *Config) {
	config.CollectTimeout = time.Duration(o)
}

// WithExporter sets the Exporter configuration option of a Config.
func WithExporter(exporter export.Exporter) Option {
	return exporterOption{exporter}
}

type exporterOption struct { export.Exporter }

func (o exporterOption) Apply(config *Config) {
	config.Exporter = o.Exporter
}

// WithExportTimeout sets the ExportTimeout configuration option of a Config.
func WithExportTimeout(timeout time.Duration) Option {
	return exportTimeoutOption(timeout)
}

type exportTimeoutOption time.Duration

func (o exportTimeoutOption) Apply(config *Config) {
	config.ExportTimeout = time.Duration(o)
}

