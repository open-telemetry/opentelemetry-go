// Copyright The OpenTelemetry Authors
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

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

const (
	addr = "127.0.0.1"
	port = 1834
)

var nc = &NetConv{
	NetHostNameKey:     attribute.Key("net.host.name"),
	NetHostPortKey:     attribute.Key("net.host.port"),
	NetPeerNameKey:     attribute.Key("net.peer.name"),
	NetPeerPortKey:     attribute.Key("net.peer.port"),
	NetSockPeerAddrKey: attribute.Key("net.sock.peer.addr"),
	NetSockPeerPortKey: attribute.Key("net.sock.peer.port"),
	NetTransportOther:  attribute.String("net.transport", "other"),
	NetTransportTCP:    attribute.String("net.transport", "ip_tcp"),
	NetTransportUDP:    attribute.String("net.transport", "ip_udp"),
	NetTransportInProc: attribute.String("net.transport", "inproc"),
}

func TestNetTransport(t *testing.T) {
	transports := map[string]attribute.KeyValue{
		"tcp":        attribute.String("net.transport", "ip_tcp"),
		"tcp4":       attribute.String("net.transport", "ip_tcp"),
		"tcp6":       attribute.String("net.transport", "ip_tcp"),
		"udp":        attribute.String("net.transport", "ip_udp"),
		"udp4":       attribute.String("net.transport", "ip_udp"),
		"udp6":       attribute.String("net.transport", "ip_udp"),
		"unix":       attribute.String("net.transport", "inproc"),
		"unixgram":   attribute.String("net.transport", "inproc"),
		"unixpacket": attribute.String("net.transport", "inproc"),
		"ip:1":       attribute.String("net.transport", "other"),
		"ip:icmp":    attribute.String("net.transport", "other"),
		"ip4:proto":  attribute.String("net.transport", "other"),
		"ip6:proto":  attribute.String("net.transport", "other"),
	}

	for net, want := range transports {
		assert.Equal(t, want, nc.Transport(net))
	}
}

func TestNetHostName(t *testing.T) {
	expected := attribute.Key("net.host.name").String(addr)
	assert.Equal(t, expected, nc.HostName(addr))
}

func TestNetHostPort(t *testing.T) {
	expected := attribute.Key("net.host.port").Int(port)
	assert.Equal(t, expected, nc.HostPort(port))
}

func TestNetPeerName(t *testing.T) {
	expected := attribute.Key("net.peer.name").String(addr)
	assert.Equal(t, expected, nc.PeerName(addr))
}

func TestNetPeerPort(t *testing.T) {
	expected := attribute.Key("net.peer.port").Int(port)
	assert.Equal(t, expected, nc.PeerPort(port))
}

func TestNetSockPeerName(t *testing.T) {
	expected := attribute.Key("net.sock.peer.addr").String(addr)
	assert.Equal(t, expected, nc.SockPeerAddr(addr))
}

func TestNetSockPeerPort(t *testing.T) {
	expected := attribute.Key("net.sock.peer.port").Int(port)
	assert.Equal(t, expected, nc.SockPeerPort(port))
}
