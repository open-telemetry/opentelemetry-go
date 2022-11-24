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

package internal // import "go.opentelemetry.io/otel/semconv/internal/v2"

import "go.opentelemetry.io/otel/attribute"

// NetConv are the network semantic convention attributes defined for a version
// of the OpenTelemetry specification.
type NetConv struct {
	NetHostNameKey     attribute.Key
	NetHostPortKey     attribute.Key
	NetPeerNameKey     attribute.Key
	NetPeerPortKey     attribute.Key
	NetSockPeerAddrKey attribute.Key
	NetSockPeerPortKey attribute.Key
	NetTransportOther  attribute.KeyValue
	NetTransportTCP    attribute.KeyValue
	NetTransportUDP    attribute.KeyValue
	NetTransportInProc attribute.KeyValue
}

func (c *NetConv) Transport(network string) attribute.KeyValue {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return c.NetTransportTCP
	case "udp", "udp4", "udp6":
		return c.NetTransportUDP
	case "unix", "unixgram", "unixpacket":
		return c.NetTransportInProc
	default:
		// "ip:*", "ip4:*", and "ip6:*" all are considered other.
		return c.NetTransportOther
	}
}

func (c *NetConv) HostName(name string) attribute.KeyValue {
	return c.NetHostNameKey.String(name)
}

func (c *NetConv) HostPort(port int) attribute.KeyValue {
	return c.NetHostPortKey.Int(port)
}

func (c *NetConv) PeerName(name string) attribute.KeyValue {
	return c.NetPeerNameKey.String(name)
}

func (c *NetConv) PeerPort(port int) attribute.KeyValue {
	return c.NetPeerPortKey.Int(port)
}

func (c *NetConv) SockPeerAddr(addr string) attribute.KeyValue {
	return c.NetSockPeerAddrKey.String(addr)
}

func (c *NetConv) SockPeerPort(port int) attribute.KeyValue {
	return c.NetSockPeerPortKey.Int(port)
}
