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
package jaeger

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestNewAgentClientUDPWithParamsBadHostport(t *testing.T) {
	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		Host: "blahblah",
		Port: "",
	})
	assert.Error(t, err)
	assert.Nil(t, agentClient)
}

func TestNewAgentClientUDPWithParams(t *testing.T) {
	mockServer, err := newUDPListener()
	require.NoError(t, err)
	defer mockServer.Close()
	host, port, err := net.SplitHostPort(mockServer.LocalAddr().String())
	assert.NoError(t, err)

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		Host:                host,
		Port:                port,
		MaxPacketSize:       25000,
		AttemptReconnecting: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, agentClient)
	assert.Equal(t, 25000, agentClient.maxPacketSize)

	if assert.IsType(t, &reconnectingUDPConn{}, agentClient.connUDP) {
		assert.Equal(t, (*log.Logger)(nil), agentClient.connUDP.(*reconnectingUDPConn).logger)
	}

	assert.NoError(t, agentClient.Close())
}

func TestNewAgentClientUDPWithParamsDefaults(t *testing.T) {
	mockServer, err := newUDPListener()
	require.NoError(t, err)
	defer mockServer.Close()
	host, port, err := net.SplitHostPort(mockServer.LocalAddr().String())
	assert.NoError(t, err)

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		Host:                host,
		Port:                port,
		AttemptReconnecting: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, agentClient)
	assert.Equal(t, udpPacketMaxLength, agentClient.maxPacketSize)

	if assert.IsType(t, &reconnectingUDPConn{}, agentClient.connUDP) {
		assert.Equal(t, (*log.Logger)(nil), agentClient.connUDP.(*reconnectingUDPConn).logger)
	}

	assert.NoError(t, agentClient.Close())
}

func TestNewAgentClientUDPWithParamsReconnectingDisabled(t *testing.T) {
	mockServer, err := newUDPListener()
	require.NoError(t, err)
	defer mockServer.Close()
	host, port, err := net.SplitHostPort(mockServer.LocalAddr().String())
	assert.NoError(t, err)

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		Host:                host,
		Port:                port,
		Logger:              nil,
		AttemptReconnecting: false,
	})
	assert.NoError(t, err)
	assert.NotNil(t, agentClient)
	assert.Equal(t, udpPacketMaxLength, agentClient.maxPacketSize)

	assert.IsType(t, &net.UDPConn{}, agentClient.connUDP)

	assert.NoError(t, agentClient.Close())
}

type errorHandler struct{ t *testing.T }

func (eh errorHandler) Handle(err error) { assert.NoError(eh.t, err) }

func TestJaegerAgentUDPLimitBatching(t *testing.T) {
	otel.SetErrorHandler(errorHandler{t})

	mockServer, err := newUDPListener()
	require.NoError(t, err)
	defer mockServer.Close()
	host, port, err := net.SplitHostPort(mockServer.LocalAddr().String())
	assert.NoError(t, err)

	// 1500 spans, size 79559, does not fit within one UDP packet with the default size of 65000.
	n := 1500
	s := make(tracetest.SpanStubs, n).Snapshots()

	exp, err := New(
		WithAgentEndpoint(WithAgentHost(host), WithAgentPort(port)),
	)
	require.NoError(t, err)

	ctx := context.Background()
	assert.NoError(t, exp.ExportSpans(ctx, s))
	assert.NoError(t, exp.Shutdown(ctx))
}

// generateALargeSpan generates a span with a long name.
func generateALargeSpan() tracetest.SpanStub {
	return tracetest.SpanStub{
		Name: "a-longer-name-that-makes-it-exceeds-limit",
	}
}

func TestSpanExceedsMaxPacketLimit(t *testing.T) {
	otel.SetErrorHandler(errorHandler{t})

	mockServer, err := newUDPListener()
	require.NoError(t, err)
	defer mockServer.Close()
	host, port, err := net.SplitHostPort(mockServer.LocalAddr().String())
	assert.NoError(t, err)

	// 106 is the serialized size of a span with default values.
	maxSize := 106

	largeSpans := tracetest.SpanStubs{generateALargeSpan(), {}}.Snapshots()
	normalSpans := tracetest.SpanStubs{{}, {}}.Snapshots()

	exp, err := New(
		WithAgentEndpoint(WithAgentHost(host), WithAgentPort(port), WithMaxPacketSize(maxSize+1)),
	)
	require.NoError(t, err)

	ctx := context.Background()
	assert.Error(t, exp.ExportSpans(ctx, largeSpans))
	assert.NoError(t, exp.ExportSpans(ctx, normalSpans))
	assert.NoError(t, exp.Shutdown(ctx))
}

func TestEmitBatchWithMultipleErrors(t *testing.T) {
	otel.SetErrorHandler(errorHandler{t})

	mockServer, err := newUDPListener()
	require.NoError(t, err)
	defer mockServer.Close()
	host, port, err := net.SplitHostPort(mockServer.LocalAddr().String())
	assert.NoError(t, err)

	span := generateALargeSpan()
	largeSpans := tracetest.SpanStubs{span, span}.Snapshots()
	// make max packet size smaller than span
	maxSize := len(span.Name)
	exp, err := New(
		WithAgentEndpoint(WithAgentHost(host), WithAgentPort(port), WithMaxPacketSize(maxSize)),
	)
	require.NoError(t, err)

	ctx := context.Background()
	err = exp.ExportSpans(ctx, largeSpans)
	assert.Error(t, err)
	require.Contains(t, err.Error(), "multiple errors")
}
