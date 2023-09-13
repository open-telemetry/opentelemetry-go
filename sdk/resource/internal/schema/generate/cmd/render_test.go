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

package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	sUtil "go.opentelemetry.io/otel/schema/v1.1"
)

var schema = `
file_format: 1.1.0
schema_url: https://opentelemetry.io/schemas/1.21.0
versions:
  1.21.0:
    spans:
      changes:
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3336
        - rename_attributes:
            attribute_map:
              messaging.kafka.client_id: messaging.client_id
              messaging.rocketmq.client_id: messaging.client_id
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3402
        - rename_attributes:
            attribute_map:
              # net.peer.(name|port) attributes were usually populated on client side
              # so they should be usually translated to server.(address|port)
              # net.host.* attributes were only populated on server side
              net.host.name: server.address
              net.host.port: server.port
              # was only populated on client side
              net.sock.peer.name: server.socket.domain
              # net.sock.peer.(addr|port) mapping is not possible
              # since they applied to both client and server side
              # were only populated on server side
              net.sock.host.addr: server.socket.address
              net.sock.host.port: server.socket.port
              http.client_ip: client.address
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3426
        - rename_attributes:
            attribute_map:
              net.protocol.name: network.protocol.name
              net.protocol.version: network.protocol.version
              net.host.connection.type: network.connection.type
              net.host.connection.subtype: network.connection.subtype
              net.host.carrier.name: network.carrier.name
              net.host.carrier.mcc: network.carrier.mcc
              net.host.carrier.mnc: network.carrier.mnc
              net.host.carrier.icc: network.carrier.icc
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3355
        - rename_attributes:
            attribute_map:
              http.method: http.request.method
              http.status_code: http.response.status_code
              http.scheme: url.scheme
              http.url: url.full
              http.request_content_length: http.request.body.size
              http.response_content_length: http.response.body.size
    metrics:
      changes:
        # https://github.com/open-telemetry/semantic-conventions/pull/53
        - rename_metrics:
            process.runtime.jvm.cpu.utilization: process.runtime.jvm.cpu.recent_utilization
  1.20.0:
    spans:
      changes:
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3272
        - rename_attributes:
            attribute_map:
              net.app.protocol.name: net.protocol.name
              net.app.protocol.version: net.protocol.version
  1.19.0:
    spans:
      changes:
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3209
        - rename_attributes:
            attribute_map:
              faas.execution: faas.invocation_id
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3188
        - rename_attributes:
            attribute_map:
              faas.id: cloud.resource_id
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3190
        - rename_attributes:
            attribute_map:
              http.user_agent: user_agent.original
    resources:
      changes:
        # https://github.com/open-telemetry/opentelemetry-specification/pull/3190
        - rename_attributes:
            attribute_map:
              browser.user_agent: user_agent.original
  1.18.0:
  1.17.0:
    spans:
      changes:
        # https://github.com/open-telemetry/opentelemetry-specification/pull/2957
        - rename_attributes:
            attribute_map:
              messaging.consumer_id: messaging.consumer.id
              messaging.protocol: net.app.protocol.name
              messaging.protocol_version: net.app.protocol.version
              messaging.destination: messaging.destination.name
              messaging.temp_destination: messaging.destination.temporary
              messaging.destination_kind: messaging.destination.kind
              messaging.message_id: messaging.message.id
              messaging.conversation_id: messaging.message.conversation_id
              messaging.message_payload_size_bytes: messaging.message.payload_size_bytes
              messaging.message_payload_compressed_size_bytes: messaging.message.payload_compressed_size_bytes
              messaging.rabbitmq.routing_key: messaging.rabbitmq.destination.routing_key
              messaging.kafka.message_key: messaging.kafka.message.key
              messaging.kafka.partition: messaging.kafka.destination.partition
              messaging.kafka.tombstone: messaging.kafka.message.tombstone
              messaging.rocketmq.message_type: messaging.rocketmq.message.type
              messaging.rocketmq.message_tag: messaging.rocketmq.message.tag
              messaging.rocketmq.message_keys: messaging.rocketmq.message.keys
              messaging.kafka.consumer_group: messaging.kafka.consumer.group
  1.16.0:
  1.15.0:
    spans:
      changes:
        # https://github.com/open-telemetry/opentelemetry-specification/pull/2743
        - rename_attributes:
            attribute_map:
              http.retry_count: http.resend_count
  1.14.0:
  1.13.0:
    spans:
      changes:
        # https://github.com/open-telemetry/opentelemetry-specification/pull/2614
        - rename_attributes:
            attribute_map:
              net.peer.ip: net.sock.peer.addr
              net.host.ip: net.sock.host.addr
  1.12.0:
  1.11.0:
  1.10.0:
  1.9.0:
  1.8.0:
    spans:
      changes:
        - rename_attributes:
            attribute_map:
              db.cassandra.keyspace: db.name
              db.hbase.namespace: db.name
  1.7.0:
  1.6.1:
  1.5.0:
  1.4.0:
`

var want = `// Code generated by "go.opentelemetry.io/otel/sdk/resource/internal/schema/gen"; DO NOT EDIT.
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

package schema

import (
	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	types10 "go.opentelemetry.io/otel/schema/v1.0/types"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
	"go.opentelemetry.io/otel/schema/v1.1/types"
)

var Schemas = map[string]*ast.Schema{
	"https://opentelemetry.io/schemas/1.21.0": {
		FileFormat: "1.1.0",
		SchemaURL:  "https://opentelemetry.io/schemas/1.21.0",
		Versions: map[types.TelemetryVersion]ast.VersionDef{
			"1.10.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.11.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.12.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.13.0": {
				All:       ast10.Attributes{},
				Resources: ast10.Attributes{},
				Spans: ast10.Spans{
					Changes: []ast10.SpansChange{
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"net.host.ip": "net.sock.host.addr",
									"net.peer.ip": "net.sock.peer.addr",
								},
							},
						},
					},
				},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.14.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.15.0": {
				All:       ast10.Attributes{},
				Resources: ast10.Attributes{},
				Spans: ast10.Spans{
					Changes: []ast10.SpansChange{
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"http.retry_count": "http.resend_count",
								},
							},
						},
					},
				},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.16.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.17.0": {
				All:       ast10.Attributes{},
				Resources: ast10.Attributes{},
				Spans: ast10.Spans{
					Changes: []ast10.SpansChange{
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"messaging.consumer_id":                           "messaging.consumer.id",
									"messaging.conversation_id":                       "messaging.message.conversation_id",
									"messaging.destination":                           "messaging.destination.name",
									"messaging.destination_kind":                      "messaging.destination.kind",
									"messaging.kafka.consumer_group":                  "messaging.kafka.consumer.group",
									"messaging.kafka.message_key":                     "messaging.kafka.message.key",
									"messaging.kafka.partition":                       "messaging.kafka.destination.partition",
									"messaging.kafka.tombstone":                       "messaging.kafka.message.tombstone",
									"messaging.message_id":                            "messaging.message.id",
									"messaging.message_payload_compressed_size_bytes": "messaging.message.payload_compressed_size_bytes",
									"messaging.message_payload_size_bytes":            "messaging.message.payload_size_bytes",
									"messaging.protocol":                              "net.app.protocol.name",
									"messaging.protocol_version":                      "net.app.protocol.version",
									"messaging.rabbitmq.routing_key":                  "messaging.rabbitmq.destination.routing_key",
									"messaging.rocketmq.message_keys":                 "messaging.rocketmq.message.keys",
									"messaging.rocketmq.message_tag":                  "messaging.rocketmq.message.tag",
									"messaging.rocketmq.message_type":                 "messaging.rocketmq.message.type",
									"messaging.temp_destination":                      "messaging.destination.temporary",
								},
							},
						},
					},
				},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.18.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.19.0": {
				All: ast10.Attributes{},
				Resources: ast10.Attributes{
					Changes: []ast10.AttributeChange{
						{
							RenameAttributes: &ast10.RenameAttributes{
								AttributeMap: map[string]string{
									"browser.user_agent": "user_agent.original",
								},
							},
						},
					},
				},
				Spans: ast10.Spans{
					Changes: []ast10.SpansChange{
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"faas.execution": "faas.invocation_id",
								},
							},
						},
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"faas.id": "cloud.resource_id",
								},
							},
						},
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"http.user_agent": "user_agent.original",
								},
							},
						},
					},
				},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.20.0": {
				All:       ast10.Attributes{},
				Resources: ast10.Attributes{},
				Spans: ast10.Spans{
					Changes: []ast10.SpansChange{
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"net.app.protocol.name":    "net.protocol.name",
									"net.app.protocol.version": "net.protocol.version",
								},
							},
						},
					},
				},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.21.0": {
				All:       ast10.Attributes{},
				Resources: ast10.Attributes{},
				Spans: ast10.Spans{
					Changes: []ast10.SpansChange{
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"messaging.kafka.client_id":    "messaging.client_id",
									"messaging.rocketmq.client_id": "messaging.client_id",
								},
							},
						},
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"http.client_ip":     "client.address",
									"net.host.name":      "server.address",
									"net.host.port":      "server.port",
									"net.sock.host.addr": "server.socket.address",
									"net.sock.host.port": "server.socket.port",
									"net.sock.peer.name": "server.socket.domain",
								},
							},
						},
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"net.host.carrier.icc":        "network.carrier.icc",
									"net.host.carrier.mcc":        "network.carrier.mcc",
									"net.host.carrier.mnc":        "network.carrier.mnc",
									"net.host.carrier.name":       "network.carrier.name",
									"net.host.connection.subtype": "network.connection.subtype",
									"net.host.connection.type":    "network.connection.type",
									"net.protocol.name":           "network.protocol.name",
									"net.protocol.version":        "network.protocol.version",
								},
							},
						},
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"http.method":                  "http.request.method",
									"http.request_content_length":  "http.request.body.size",
									"http.response_content_length": "http.response.body.size",
									"http.scheme":                  "url.scheme",
									"http.status_code":             "http.response.status_code",
									"http.url":                     "url.full",
								},
							},
						},
					},
				},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics: ast.Metrics{
					Changes: []ast.MetricsChange{
						{
							RenameMetrics: map[types10.MetricName]types10.MetricName{
								"process.runtime.jvm.cpu.utilization": "process.runtime.jvm.cpu.recent_utilization",
							},
							RenameAttributes: &ast10.AttributeMapForMetrics{},
							Split:            &ast.SplitMetric{},
						},
					},
				},
			},
			"1.4.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.5.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.6.1": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.7.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.8.0": {
				All:       ast10.Attributes{},
				Resources: ast10.Attributes{},
				Spans: ast10.Spans{
					Changes: []ast10.SpansChange{
						{
							RenameAttributes: &ast10.AttributeMapForSpans{
								ApplyToSpans: []types10.SpanName{},
								AttributeMap: map[string]string{
									"db.cassandra.keyspace": "db.name",
									"db.hbase.namespace":    "db.name",
								},
							},
						},
					},
				},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
			"1.9.0": {
				All:        ast10.Attributes{},
				Resources:  ast10.Attributes{},
				Spans:      ast10.Spans{},
				SpanEvents: ast10.SpanEvents{},
				Logs:       ast10.Logs{},
				Metrics:    ast.Metrics{},
			},
		},
	}}
`

func TestRender(t *testing.T) {
	f := bytes.NewReader([]byte(schema))
	s, err := sUtil.Parse(f)
	require.NoError(t, err)

	e, err := newEntry(s)
	require.NoError(t, err)

	var dest bytes.Buffer
	render(&dest, []entry{e})
}
