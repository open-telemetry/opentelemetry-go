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

package instrument // import "go.opentelemetry.io/otel/metric/instrument"

import (
	"go.opentelemetry.io/otel/metric/unit"
)

// Float64Config contains options for Asynchronous instruments that
// observe float64 values.
type Float64Config struct {
	description string
	unit        unit.Unit
}

// Float64Config contains options for Synchronous instruments that record
// float64 values.
func NewFloat64Config(opts ...Float64Option) Float64Config {
	var config Float64Config
	for _, o := range opts {
		config = o.applyFloat64(config)
	}
	return config
}

// Description returns the Config description.
func (c Float64Config) Description() string {
	return c.description
}

// Unit returns the Config unit.
func (c Float64Config) Unit() unit.Unit {
	return c.unit
}

// Float64Option applies options to synchronous float64 instruments.
type Float64Option interface {
	applyFloat64(Float64Config) Float64Config
}
