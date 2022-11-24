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

// Package netconv provides OpenTelemetry semantic convetions for the net
// package from the standard library.
package netconv // import "go.opentelemetry.io/otel/semconv/v1.13.0/netconv"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv/internal/v2"
	semconv "go.opentelemetry.io/otel/semconv/v1.13.0"
)

var nc = &internal.NetConv{
	NetHostNameKey:     semconv.NetHostNameKey,
	NetHostPortKey:     semconv.NetHostPortKey,
	NetPeerNameKey:     semconv.NetPeerNameKey,
	NetPeerPortKey:     semconv.NetPeerPortKey,
	NetSockPeerAddrKey: semconv.NetSockPeerAddrKey,
	NetSockPeerPortKey: semconv.NetSockPeerPortKey,
	NetTransportOther:  semconv.NetTransportOther,
	NetTransportTCP:    semconv.NetTransportTCP,
	NetTransportUDP:    semconv.NetTransportUDP,
	NetTransportInProc: semconv.NetTransportInProc,
}

// Transport returns an attribute describing the transport protocol of the
// passed network. See the net.Dial for information about acceptable network
// values.
func Transport(network string) attribute.KeyValue {
	return nc.Transport(network)
}
