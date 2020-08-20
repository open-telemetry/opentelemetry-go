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
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAgentClientUDPWithParamsBadHostport(t *testing.T) {
	hostPort := "blahblah"

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		HostPort: hostPort,
	})

	assert.Error(t, err)
	assert.Nil(t, agentClient)
}

func TestNewAgentClientUDPWithParams(t *testing.T) {
	mockServer, err := newUDPListener()
	require.NoError(t, err)
	defer mockServer.Close()

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		HostPort:            mockServer.LocalAddr().String(),
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

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		HostPort:            mockServer.LocalAddr().String(),
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

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		HostPort:            mockServer.LocalAddr().String(),
		Logger:              nil,
		AttemptReconnecting: false,
	})
	assert.NoError(t, err)
	assert.NotNil(t, agentClient)
	assert.Equal(t, udpPacketMaxLength, agentClient.maxPacketSize)

	assert.IsType(t, &net.UDPConn{}, agentClient.connUDP)

	assert.NoError(t, agentClient.Close())
}
