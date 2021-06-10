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
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockResolver struct {
	mock.Mock
}

func (m *mockResolver) ResolveUDPAddr(network string, hostPort string) (*net.UDPAddr, error) {
	args := m.Called(network, hostPort)

	a0 := args.Get(0)
	if a0 == nil {
		return (*net.UDPAddr)(nil), args.Error(1)
	}
	return a0.(*net.UDPAddr), args.Error(1)
}

type mockDialer struct {
	mock.Mock
}

func (m *mockDialer) DialUDP(network string, laddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	args := m.Called(network, laddr, raddr)

	a0 := args.Get(0)
	if a0 == nil {
		return (*net.UDPConn)(nil), args.Error(1)
	}

	return a0.(*net.UDPConn), args.Error(1)
}

func newUDPListener() (net.PacketConn, error) {
	return net.ListenPacket("udp", "127.0.0.1:0")
}

func newUDPConn() (net.PacketConn, *net.UDPConn, error) {
	mockServer, err := newUDPListener()
	if err != nil {
		return nil, nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", mockServer.LocalAddr().String())
	if err != nil {
		mockServer.Close()
		return nil, nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		mockServer.Close()
		return nil, nil, err
	}

	return mockServer, conn, nil
}

func assertConnWritable(t *testing.T, conn udpConn, serverConn net.PacketConn) {
	expectedString := "yo this is a test"
	_, err := conn.Write([]byte(expectedString))
	require.NoError(t, err)

	var buf = make([]byte, len(expectedString))
	err = serverConn.SetReadDeadline(time.Now().Add(time.Second))
	require.NoError(t, err)

	_, _, err = serverConn.ReadFrom(buf)
	require.NoError(t, err)
	require.Equal(t, []byte(expectedString), buf)
}

func waitForCallWithTimeout(call *mock.Call) bool {
	called := make(chan struct{})
	call.Run(func(args mock.Arguments) {
		if !isChannelClosed(called) {
			close(called)
		}
	})

	var wasCalled bool
	// wait at most 100 milliseconds for the second call of ResolveUDPAddr that is supposed to fail
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	select {
	case <-called:
		wasCalled = true
	case <-ctx.Done():
		fmt.Println("timed out")
	}
	cancel()

	return wasCalled
}

func isChannelClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

func waitForConnCondition(conn *reconnectingUDPConn, condition func(conn *reconnectingUDPConn) bool) bool {
	var conditionVal bool
	for i := 0; i < 10; i++ {
		conn.connMtx.RLock()
		conditionVal = condition(conn)
		conn.connMtx.RUnlock()
		if conditionVal || i >= 9 {
			break
		}

		time.Sleep(time.Millisecond * 10)
	}

	return conditionVal
}

func newMockUDPAddr(t *testing.T, port int) *net.UDPAddr {
	var buf = make([]byte, 4)
	// random is not seeded to ensure tests are deterministic (also doesnt matter if ip is valid)
	_, err := rand.Read(buf)
	require.NoError(t, err)

	return &net.UDPAddr{
		IP:   net.IPv4(buf[0], buf[1], buf[2], buf[3]),
		Port: port,
	}
}

func TestNewResolvedUDPConn(t *testing.T) {
	hostPort := "blahblah:34322"

	mockServer, clientConn, err := newUDPConn()
	require.NoError(t, err)
	defer mockServer.Close()

	mockUDPAddr := newMockUDPAddr(t, 34322)

	resolver := mockResolver{}
	resolver.
		On("ResolveUDPAddr", "udp", hostPort).
		Return(mockUDPAddr, nil).
		Once()

	dialer := mockDialer{}
	dialer.
		On("DialUDP", "udp", (*net.UDPAddr)(nil), mockUDPAddr).
		Return(clientConn, nil).
		Once()

	conn, err := newReconnectingUDPConn(hostPort, udpPacketMaxLength, time.Hour, resolver.ResolveUDPAddr, dialer.DialUDP, nil)
	assert.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.Close()
	assert.NoError(t, err)

	// assert the actual connection was closed
	assert.Error(t, clientConn.Close())

	resolver.AssertExpectations(t)
	dialer.AssertExpectations(t)
}

func TestResolvedUDPConnWrites(t *testing.T) {
	hostPort := "blahblah:34322"

	mockServer, clientConn, err := newUDPConn()
	require.NoError(t, err)
	defer mockServer.Close()

	mockUDPAddr := newMockUDPAddr(t, 34322)

	resolver := mockResolver{}
	resolver.
		On("ResolveUDPAddr", "udp", hostPort).
		Return(mockUDPAddr, nil).
		Once()

	dialer := mockDialer{}
	dialer.
		On("DialUDP", "udp", (*net.UDPAddr)(nil), mockUDPAddr).
		Return(clientConn, nil).
		Once()

	conn, err := newReconnectingUDPConn(hostPort, udpPacketMaxLength, time.Hour, resolver.ResolveUDPAddr, dialer.DialUDP, nil)
	assert.NoError(t, err)
	require.NotNil(t, conn)

	assertConnWritable(t, conn, mockServer)

	err = conn.Close()
	assert.NoError(t, err)

	// assert the actual connection was closed
	assert.Error(t, clientConn.Close())

	resolver.AssertExpectations(t)
	dialer.AssertExpectations(t)
}

func TestResolvedUDPConnEventuallyDials(t *testing.T) {
	hostPort := "blahblah:34322"

	mockServer, clientConn, err := newUDPConn()
	require.NoError(t, err)
	defer mockServer.Close()

	mockUDPAddr := newMockUDPAddr(t, 34322)

	resolver := mockResolver{}
	resolver.
		On("ResolveUDPAddr", "udp", hostPort).
		Return(nil, fmt.Errorf("failed to resolve")).Once().
		On("ResolveUDPAddr", "udp", hostPort).
		Return(mockUDPAddr, nil)

	dialer := mockDialer{}
	dialCall := dialer.
		On("DialUDP", "udp", (*net.UDPAddr)(nil), mockUDPAddr).
		Return(clientConn, nil).Once()

	conn, err := newReconnectingUDPConn(hostPort, udpPacketMaxLength, time.Millisecond*10, resolver.ResolveUDPAddr, dialer.DialUDP, nil)
	assert.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.SetWriteBuffer(udpPacketMaxLength)
	assert.NoError(t, err)

	wasCalled := waitForCallWithTimeout(dialCall)
	assert.True(t, wasCalled)

	connEstablished := waitForConnCondition(conn, func(conn *reconnectingUDPConn) bool {
		return conn.conn != nil
	})

	assert.True(t, connEstablished)

	assertConnWritable(t, conn, mockServer)
	assertSockBufferSize(t, udpPacketMaxLength, clientConn)

	err = conn.Close()
	assert.NoError(t, err)

	// assert the actual connection was closed
	assert.Error(t, clientConn.Close())

	resolver.AssertExpectations(t)
	dialer.AssertExpectations(t)
}

func TestResolvedUDPConnNoSwapIfFail(t *testing.T) {
	hostPort := "blahblah:34322"

	mockServer, clientConn, err := newUDPConn()
	require.NoError(t, err)
	defer mockServer.Close()

	mockUDPAddr := newMockUDPAddr(t, 34322)

	resolver := mockResolver{}
	resolver.
		On("ResolveUDPAddr", "udp", hostPort).
		Return(mockUDPAddr, nil).Once()

	failCall := resolver.On("ResolveUDPAddr", "udp", hostPort).
		Return(nil, fmt.Errorf("resolve failed"))

	dialer := mockDialer{}
	dialer.
		On("DialUDP", "udp", (*net.UDPAddr)(nil), mockUDPAddr).
		Return(clientConn, nil).Once()

	conn, err := newReconnectingUDPConn(hostPort, udpPacketMaxLength, time.Millisecond*10, resolver.ResolveUDPAddr, dialer.DialUDP, nil)
	assert.NoError(t, err)
	require.NotNil(t, conn)

	wasCalled := waitForCallWithTimeout(failCall)

	assert.True(t, wasCalled)

	assertConnWritable(t, conn, mockServer)

	err = conn.Close()
	assert.NoError(t, err)

	// assert the actual connection was closed
	assert.Error(t, clientConn.Close())

	resolver.AssertExpectations(t)
	dialer.AssertExpectations(t)
}

func TestResolvedUDPConnWriteRetry(t *testing.T) {
	hostPort := "blahblah:34322"

	mockServer, clientConn, err := newUDPConn()
	require.NoError(t, err)
	defer mockServer.Close()

	mockUDPAddr := newMockUDPAddr(t, 34322)

	resolver := mockResolver{}
	resolver.
		On("ResolveUDPAddr", "udp", hostPort).
		Return(nil, fmt.Errorf("failed to resolve")).Once().
		On("ResolveUDPAddr", "udp", hostPort).
		Return(mockUDPAddr, nil).Once()

	dialer := mockDialer{}
	dialer.
		On("DialUDP", "udp", (*net.UDPAddr)(nil), mockUDPAddr).
		Return(clientConn, nil).Once()

	conn, err := newReconnectingUDPConn(hostPort, udpPacketMaxLength, time.Millisecond*10, resolver.ResolveUDPAddr, dialer.DialUDP, nil)
	assert.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.SetWriteBuffer(udpPacketMaxLength)
	assert.NoError(t, err)

	assertConnWritable(t, conn, mockServer)
	assertSockBufferSize(t, udpPacketMaxLength, clientConn)

	err = conn.Close()
	assert.NoError(t, err)

	// assert the actual connection was closed
	assert.Error(t, clientConn.Close())

	resolver.AssertExpectations(t)
	dialer.AssertExpectations(t)
}

func TestResolvedUDPConnWriteRetryFails(t *testing.T) {
	hostPort := "blahblah:34322"

	resolver := mockResolver{}
	resolver.
		On("ResolveUDPAddr", "udp", hostPort).
		Return(nil, fmt.Errorf("failed to resolve")).Twice()

	dialer := mockDialer{}

	conn, err := newReconnectingUDPConn(hostPort, udpPacketMaxLength, time.Millisecond*10, resolver.ResolveUDPAddr, dialer.DialUDP, nil)
	assert.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.SetWriteBuffer(udpPacketMaxLength)
	assert.NoError(t, err)

	_, err = conn.Write([]byte("yo this is a test"))

	assert.Error(t, err)

	err = conn.Close()
	assert.NoError(t, err)

	resolver.AssertExpectations(t)
	dialer.AssertExpectations(t)
}

func TestResolvedUDPConnChanges(t *testing.T) {
	hostPort := "blahblah:34322"

	mockServer, clientConn, err := newUDPConn()
	require.NoError(t, err)
	defer mockServer.Close()

	mockUDPAddr1 := newMockUDPAddr(t, 34322)

	mockServer2, clientConn2, err := newUDPConn()
	require.NoError(t, err)
	defer mockServer2.Close()

	mockUDPAddr2 := newMockUDPAddr(t, 34322)

	// ensure address doesn't duplicate mockUDPAddr1
	for i := 0; i < 10 && mockUDPAddr2.IP.Equal(mockUDPAddr1.IP); i++ {
		mockUDPAddr2 = newMockUDPAddr(t, 34322)
	}

	// this is really unlikely to ever fail the test, but its here as a safeguard
	require.False(t, mockUDPAddr2.IP.Equal(mockUDPAddr1.IP))

	resolver := mockResolver{}
	resolver.
		On("ResolveUDPAddr", "udp", hostPort).
		Return(mockUDPAddr1, nil).Once().
		On("ResolveUDPAddr", "udp", hostPort).
		Return(mockUDPAddr2, nil)

	dialer := mockDialer{}
	dialer.
		On("DialUDP", "udp", (*net.UDPAddr)(nil), mockUDPAddr1).
		Return(clientConn, nil).Once()

	secondDial := dialer.
		On("DialUDP", "udp", (*net.UDPAddr)(nil), mockUDPAddr2).
		Return(clientConn2, nil).Once()

	conn, err := newReconnectingUDPConn(hostPort, udpPacketMaxLength, time.Millisecond*10, resolver.ResolveUDPAddr, dialer.DialUDP, nil)
	assert.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.SetWriteBuffer(udpPacketMaxLength)
	assert.NoError(t, err)

	wasCalled := waitForCallWithTimeout(secondDial)
	assert.True(t, wasCalled)

	connSwapped := waitForConnCondition(conn, func(conn *reconnectingUDPConn) bool {
		return conn.conn == clientConn2
	})

	assert.True(t, connSwapped)

	assertConnWritable(t, conn, mockServer2)
	assertSockBufferSize(t, udpPacketMaxLength, clientConn2)

	err = conn.Close()
	assert.NoError(t, err)

	// assert the prev connection was closed
	assert.Error(t, clientConn.Close())

	// assert the actual connection was closed
	assert.Error(t, clientConn2.Close())

	resolver.AssertExpectations(t)
	dialer.AssertExpectations(t)
}

func TestResolvedUDPConnLoopWithoutChanges(t *testing.T) {
	hostPort := "blahblah:34322"

	mockServer, clientConn, err := newUDPConn()
	require.NoError(t, err)
	defer mockServer.Close()

	mockUDPAddr := newMockUDPAddr(t, 34322)

	resolver := mockResolver{}
	resolver.
		On("ResolveUDPAddr", "udp", hostPort).
		Return(mockUDPAddr, nil)

	dialer := mockDialer{}
	dialer.
		On("DialUDP", "udp", (*net.UDPAddr)(nil), mockUDPAddr).
		Return(clientConn, nil).
		Once()

	resolveTimeout := 500 * time.Millisecond
	conn, err := newReconnectingUDPConn(hostPort, udpPacketMaxLength, resolveTimeout, resolver.ResolveUDPAddr, dialer.DialUDP, nil)
	assert.NoError(t, err)
	require.NotNil(t, conn)
	assert.Equal(t, mockUDPAddr, conn.destAddr)

	// Waiting for one round of loop
	time.Sleep(3 * resolveTimeout)
	assert.Equal(t, mockUDPAddr, conn.destAddr)

	err = conn.Close()
	assert.NoError(t, err)

	// assert the actual connection was closed
	assert.Error(t, clientConn.Close())

	resolver.AssertExpectations(t)
	dialer.AssertExpectations(t)
}
