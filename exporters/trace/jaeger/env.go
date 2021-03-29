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

package jaeger // import "go.opentelemetry.io/otel/exporters/trace/jaeger"

import (
	"os"
	"strconv"
)

// Environment variable names
const (
	// Whether the exporter is disabled or not. (default false).
	envDisabled = "OTEL_EXPORTER_JAEGER_DISABLED"
	// The HTTP endpoint for sending spans directly to a collector,
	// i.e. http://jaeger-collector:14268/api/traces.
	envEndpoint = "OTEL_EXPORTER_JAEGER_ENDPOINT"
	// Username to send as part of "Basic" authentication to the collector endpoint.
	envUser = "OTEL_EXPORTER_JAEGER_USER"
	// Password to send as part of "Basic" authentication to the collector endpoint.
	envPassword = "OTEL_EXPORTER_JAEGER_PASSWORD"
)

// CollectorEndpointFromEnv return environment variable value of JAEGER_ENDPOINT
func CollectorEndpointFromEnv() string {
	return os.Getenv(envEndpoint)
}

// WithCollectorEndpointOptionFromEnv uses environment variables to set the username and password
// if basic auth is required.
func WithCollectorEndpointOptionFromEnv() CollectorEndpointOption {
	return func(o *CollectorEndpointOptions) {
		if e := os.Getenv(envUser); e != "" {
			o.username = e
		}
		if e := os.Getenv(envPassword); e != "" {
			o.password = os.Getenv(envPassword)
		}
	}
}

// WithDisabledFromEnv uses environment variables and overrides disabled field.
func WithDisabledFromEnv() Option {
	return func(o *options) {
		if e := os.Getenv(envDisabled); e != "" {
			if v, err := strconv.ParseBool(e); err == nil {
				o.Disabled = v
			}
		}
	}
}
