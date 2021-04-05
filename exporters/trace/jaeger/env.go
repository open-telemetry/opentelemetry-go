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
)

// Environment variable names
const (
	// Hostname for the Jaeger agent, part of address where exporter sends spans
	// i.e.	"localhost"
	envAgentHost = "OTEL_EXPORTER_JAEGER_AGENT_HOST"
	// Port for the Jaeger agent, part of address where exporter sends spans
	// i.e. 6832
	envAgentPort = "OTEL_EXPORTER_JAEGER_AGENT_PORT"
	// The HTTP endpoint for sending spans directly to a collector,
	// i.e. http://jaeger-collector:14268/api/traces.
	envEndpoint = "OTEL_EXPORTER_JAEGER_ENDPOINT"
	// Username to send as part of "Basic" authentication to the collector endpoint.
	envUser = "OTEL_EXPORTER_JAEGER_USER"
	// Password to send as part of "Basic" authentication to the collector endpoint.
	envPassword = "OTEL_EXPORTER_JAEGER_PASSWORD"
)

// agentEndpointFromEnv returns env vars values of OTEL_EXPORTER_JAEGER_AGENT_HOST and OTEL_EXPORTER_JAEGER_AGENT_PORT
func agentEndpointFromEnv() (string, string) {
	h := os.Getenv(envAgentHost)
	p := os.Getenv(envAgentPort)
	return h, p
}

// CollectorEndpointFromEnv return environment variable value of OTEL_EXPORTER_JAEGER_ENDPOINT
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
			o.password = e
		}
	}
}
