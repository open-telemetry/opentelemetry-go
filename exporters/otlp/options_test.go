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

package otlp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectionsSetsDefaultOptions(t *testing.T) {

	config := NewConnections()

	expectedAddress := fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorPort)
	assert.Equal(t, expectedAddress, config.metrics.collectorAddr,
		"expected different metrics collector address")
	assert.Equal(t, expectedAddress, config.traces.collectorAddr,
		"expected different traces collector address")

	assert.Equal(t, DefaultNumWorkers, config.metrics.numWorkers,
		"expected different metrics number of workers")
	assert.Equal(t, DefaultNumWorkers, config.traces.numWorkers,
		"expected different traces number of workers")

	assert.Equal(t, DefaultGRPCServiceConfig, config.metrics.grpcServiceConfig,
		"expected different metrics grpc service config")
	assert.Equal(t, DefaultGRPCServiceConfig, config.traces.grpcServiceConfig,
		"expected different traces grpc service config")

}

func TestNewTraceConnectionsSetsDefaultOptionsOnlyForTraces(t *testing.T) {

	config := NewTracesConnection()

	expectedAddress := fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorPort)
	assert.Equal(t, "", config.metrics.collectorAddr,
		"expected different metrics collector address")
	assert.Equal(t, expectedAddress, config.traces.collectorAddr,
		"expected different traces collector address")

	assert.Equal(t, uint(0), config.metrics.numWorkers,
		"expected different metrics number of workers")
	assert.Equal(t, DefaultNumWorkers, config.traces.numWorkers,
		"expected different traces number of workers")

	assert.Equal(t, "", config.metrics.grpcServiceConfig,
		"expected different metrics grpc service config")
	assert.Equal(t, DefaultGRPCServiceConfig, config.traces.grpcServiceConfig,
		"expected different traces grpc service config")

}

func TestNewMetricsConnectionsSetsDefaultOptionsOnlyForMetrics(t *testing.T) {

	config := NewMetricsConnection()

	expectedAddress := fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorPort)
	assert.Equal(t, expectedAddress, config.metrics.collectorAddr,
		"expected different metrics collector address")
	assert.Equal(t, "", config.traces.collectorAddr,
		"expected different traces collector address")

	assert.Equal(t, DefaultNumWorkers, config.metrics.numWorkers,
		"expected different metrics number of workers")
	assert.Equal(t, uint(0), config.traces.numWorkers,
		"expected different traces number of workers")

	assert.Equal(t, DefaultGRPCServiceConfig, config.metrics.grpcServiceConfig,
		"expected different metrics grpc service config")
	assert.Equal(t, "", config.traces.grpcServiceConfig,
		"expected different traces grpc service config")

}

func TestNewConnectionsSetsCommonOptions(t *testing.T) {

	expectedAddress := "foo"
	options := []ExporterOption{
		WithInsecure(),
		WithAddress(expectedAddress),
	}

	config := NewConnections(options...)

	assert.True(t, config.metrics.canDialInsecure,
		"expected metrics connection to dial insecure")
	assert.True(t, config.traces.canDialInsecure,
		"expected traces connection to dial insecure")
	assert.Equal(t, config.metrics.collectorAddr, expectedAddress,
		"expected different metrics collector address")
	assert.Equal(t, config.traces.collectorAddr, expectedAddress,
		"expected different traces collector address")
}

func TestSetMetricOptionsOverridesCommonOptions(t *testing.T) {

	metricsAddress := "metrics"
	tracesAddress := "traces"
	config := NewConnections(
		WithInsecure(),
		WithAddress(tracesAddress),
	).SetMetricOptions(WithAddress(metricsAddress))

	assert.True(t, config.metrics.canDialInsecure,
		"expected metrics connection to dial insecure")
	assert.True(t, config.traces.canDialInsecure,
		"expected traces connection to dial insecure")

	assert.Equal(t, config.metrics.collectorAddr, metricsAddress,
		"expected different metrics collector address")
	assert.Equal(t, config.traces.collectorAddr, tracesAddress,
		"expected different traces collector address")
}

func TestSetTraceOptionsOverridesCommonOptions(t *testing.T) {
	metricsAddress := "metrics"
	tracesAddress := "traces"
	config := NewConnections(
		WithInsecure(),
		WithAddress(metricsAddress),
	).SetTraceOptions(WithAddress(tracesAddress))

	assert.True(t, config.metrics.canDialInsecure,
		"expected metrics connection to dial insecure")
	assert.True(t, config.traces.canDialInsecure,
		"expected traces connection to dial insecure")

	assert.Equal(t, config.metrics.collectorAddr, metricsAddress,
		"expected different metrics collector address")
	assert.Equal(t, config.traces.collectorAddr, tracesAddress,
		"expected different traces collector address")
}

func TestSetCommonOptionsOverridesDefaultOptions(t *testing.T) {
	commonAddress := "common"
	config := NewConnections().
		SetCommonOptions(WithAddress(commonAddress), WithInsecure())

	assert.True(t, config.metrics.canDialInsecure,
		"expected metrics connection to dial insecure")
	assert.True(t, config.traces.canDialInsecure,
		"expected traces connection to dial insecure")

	assert.Equal(t, config.metrics.collectorAddr, commonAddress,
		"expected different metrics collector address")
	assert.Equal(t, config.traces.collectorAddr, commonAddress,
		"expected different traces collector address")
}
