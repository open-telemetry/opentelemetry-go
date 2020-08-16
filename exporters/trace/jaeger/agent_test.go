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
		HostPort:      mockServer.LocalAddr().String(),
		MaxPacketSize: 25000,
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
	mockServer, err := newUDPListenerOnPort(6831)
	require.NoError(t, err)
	defer mockServer.Close()

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		HostPort: "localhost:6831",
	})
	assert.NoError(t, err)
	assert.NotNil(t, agentClient)
	assert.Equal(t, udpPacketMaxLength, agentClient.maxPacketSize)

	if assert.IsType(t, &reconnectingUDPConn{}, agentClient.connUDP) {
		assert.Equal(t, (*log.Logger)(nil), agentClient.connUDP.(*reconnectingUDPConn).logger)
	}

	assert.NoError(t, agentClient.Close())
}

func TestNewAgentClientUDPDefaults(t *testing.T) {
	mockServer, err := newUDPListenerOnPort(6831)
	require.NoError(t, err)
	defer mockServer.Close()

	agentClient, err := newAgentClientUDP(agentClientUDPParams{
		HostPort: "localhost:6831",
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
		HostPort:                   mockServer.LocalAddr().String(),
		Logger:                     nil,
		DisableAttemptReconnecting: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, agentClient)
	assert.Equal(t, udpPacketMaxLength, agentClient.maxPacketSize)

	assert.IsType(t, &net.UDPConn{}, agentClient.connUDP)

	assert.NoError(t, agentClient.Close())
}
