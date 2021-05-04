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
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
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

	// 1500 spans, size 79559, does not fit within one UDP packet with the default size of 65000.
	n := 1500
	s := make([]*tracesdk.SpanSnapshot, n)
	for i := 0; i < n; i++ {
		s[i] = &tracesdk.SpanSnapshot{}
	}

	exp, err := NewRawExporter(
		WithAgentEndpoint(WithAgentHost("localhost"), WithAgentPort("6831")),
	)
	require.NoError(t, err)

	ctx := context.Background()
	assert.NoError(t, exp.ExportSpans(ctx, s))
	assert.NoError(t, exp.Shutdown(ctx))
}

// generateALargeSpan generates a span with a long name.
func generateALargeSpan() *tracesdk.SpanSnapshot {
	span := &tracesdk.SpanSnapshot{
		Name: "a-longer-name-that-makes-it-exceeds-limit",
	}
	return span
}

func TestSpanExceedsMaxPacketLimit(t *testing.T) {
	otel.SetErrorHandler(errorHandler{t})

	// 106 is the serialized size of a span with default values.
	maxSize := 106
	span := generateALargeSpan()

	largeSpans := []*tracesdk.SpanSnapshot{span, {}}
	normalSpans := []*tracesdk.SpanSnapshot{{}, {}}

	exp, err := NewRawExporter(
		WithAgentEndpoint(WithAgentHost("localhost"), WithAgentPort("6831"), WithMaxPacketSize(maxSize+1)),
	)
	require.NoError(t, err)

	ctx := context.Background()
	assert.Error(t, exp.ExportSpans(ctx, largeSpans))
	assert.NoError(t, exp.ExportSpans(ctx, normalSpans))
	assert.NoError(t, exp.Shutdown(ctx))
}

func TestEmitBatchWithMultipleErrors(t *testing.T) {
	otel.SetErrorHandler(errorHandler{t})

	span := generateALargeSpan()
	largeSpans := []*tracesdk.SpanSnapshot{span, span}
	// make max packet size smaller than span
	maxSize := len(span.Name)
	exp, err := NewRawExporter(
		WithAgentEndpoint(WithAgentHost("localhost"), WithAgentPort("6831"), WithMaxPacketSize(maxSize)),
	)
	require.NoError(t, err)

	ctx := context.Background()
	err = exp.ExportSpans(ctx, largeSpans)
	assert.Error(t, err)
	require.Contains(t, err.Error(), "multiple errors")
}
