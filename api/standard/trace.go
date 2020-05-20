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

package standard

import "go.opentelemetry.io/otel/api/kv"

// Standard attribute keys used for network related operations.
const (
	// Transport protocol used.
	NetTransportKey = kv.Key("net.transport")

	// Remote address of the peer.
	NetPeerIPKey = kv.Key("net.peer.ip")

	// Remote port number.
	NetPeerPortKey = kv.Key("net.peer.port")

	// Remote hostname or similar.
	NetPeerNameKey = kv.Key("net.peer.name")

	// Local host IP. Useful in case of a multi-IP host.
	NetHostIPKey = kv.Key("net.host.ip")

	// Local host port.
	NetHostPortKey = kv.Key("net.host.port")

	// Local hostname or similar.
	NetHostNameKey = kv.Key("net.host.name")
)

var (
	NetTransportTCP    = NetTransportKey.String("IP.TCP")
	NetTransportUDP    = NetTransportKey.String("IP.UDP")
	NetTransportIP     = NetTransportKey.String("IP")
	NetTransportUnix   = NetTransportKey.String("Unix")
	NetTransportPipe   = NetTransportKey.String("pipe")
	NetTransportInProc = NetTransportKey.String("inproc")
	NetTransportOther  = NetTransportKey.String("other")
)

// Standard attribute keys used to identify an authorized enduser.
const (
	// Username or the client identifier extracted from the access token or
	// authorization header in the inbound request from outside the system.
	EnduserIDKey = kv.Key("enduser.id")

	// Actual or assumed role the client is making the request with.
	EnduserRoleKey = kv.Key("enduser.role")

	// Scopes or granted authorities the client currently possesses.
	EnduserScopeKey = kv.Key("enduser.scope")
)

// Standard attribute keys for HTTP.
const (
	// HTTP request method.
	HTTPMethodKey = kv.Key("http.method")

	// Full HTTP request URL in the form:
	// scheme://host[:port]/path?query[#fragment].
	HTTPUrlKey = kv.Key("http.url")

	// The full request target as passed in a HTTP request line or
	// equivalent, e.g. "/path/12314/?q=ddds#123".
	HTTPTargetKey = kv.Key("http.target")

	// The value of the HTTP host header.
	HTTPHostKey = kv.Key("http.host")

	// The URI scheme identifying the used protocol.
	HTTPSchemeKey = kv.Key("http.scheme")

	// HTTP response status code.
	HTTPStatusCodeKey = kv.Key("http.status_code")

	// HTTP reason phrase.
	HTTPStatusTextKey = kv.Key("http.status_text")

	// Kind of HTTP protocol used.
	HTTPFlavorKey = kv.Key("http.flavor")

	// Value of the HTTP User-Agent header sent by the client.
	HTTPUserAgentKey = kv.Key("http.user_agent")

	// The primary server name of the matched virtual host.
	HTTPServerNameKey = kv.Key("http.server_name")

	// The matched route served (path template). For example,
	// "/users/:userID?".
	HTTPRouteKey = kv.Key("http.route")

	// The IP address of the original client behind all proxies, if known
	// (e.g. from X-Forwarded-For).
	HTTPClientIPKey = kv.Key("http.client_ip")
)

var (
	HTTPSchemeHTTP  = HTTPSchemeKey.String("http")
	HTTPSchemeHTTPS = HTTPSchemeKey.String("https")

	HTTPFlavor1_0  = HTTPFlavorKey.String("1.0")
	HTTPFlavor1_1  = HTTPFlavorKey.String("1.1")
	HTTPFlavor2    = HTTPFlavorKey.String("2")
	HTTPFlavorSPDY = HTTPFlavorKey.String("SPDY")
	HTTPFlavorQUIC = HTTPFlavorKey.String("QUIC")
)

// Standard attribute keys for database clients.
const (
	// Database type. For any SQL database, "sql". For others, the
	// lower-case database category, e.g. "cassandra", "hbase", or "redis".
	DBTypeKey = kv.Key("db.type")

	// Database instance name.
	DBInstanceKey = kv.Key("db.instance")

	// A database statement for the given database type.
	DBStatementKey = kv.Key("db.statement")

	// Username for accessing database.
	DBUserKey = kv.Key("db.user")

	// Database URL.
	DBUrlKey = kv.Key("db.url")
)

// Standard attribute keys for RPC.
const (
	// The RPC service name.
	RPCServiceKey = kv.Key("rpc.service")

	// Name of message transmitted or received.
	RPCNameKey = kv.Key("name")

	// Type of message transmitted or received.
	RPCMessageTypeKey = kv.Key("message.type")

	// Identifier of message transmitted or received.
	RPCMessageIDKey = kv.Key("message.id")

	// The compressed size of the message transmitted or received in bytes.
	RPCMessageCompressedSizeKey = kv.Key("message.compressed_size")

	// The uncompressed size of the message transmitted or received in
	// bytes.
	RPCMessageUncompressedSizeKey = kv.Key("message.uncompressed_size")
)

var (
	RPCNameMessage = RPCNameKey.String("message")

	RPCMessageTypeSent     = RPCMessageTypeKey.String("SENT")
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")
)

// Standard attribute keys for messaging systems.
const (
	// A unique identifier describing the messaging system. For example,
	// kafka, rabbitmq or activemq.
	MessagingSystemKey = kv.Key("messaging.system")

	// The message destination name, e.g. MyQueue or MyTopic.
	MessagingDestinationKey = kv.Key("messaging.destination")

	// The kind of message destination.
	MessagingDestinationKindKey = kv.Key("messaging.destination_kind")

	// Describes if the destination is temporary or not.
	MessagingTempDestinationKey = kv.Key("messaging.temp_destination")

	// The name of the transport protocol.
	MessagingProtocolKey = kv.Key("messaging.protocol")

	// The version of the transport protocol.
	MessagingProtocolVersionKey = kv.Key("messaging.protocol_version")

	// Messaging service URL.
	MessagingURLKey = kv.Key("messaging.url")

	// Identifier used by the messaging system for a message.
	MessagingMessageIDKey = kv.Key("messaging.message_id")

	// Identifier used by the messaging system for a conversation.
	MessagingConversationIDKey = kv.Key("messaging.conversation_id")

	// The (uncompressed) size of the message payload in bytes.
	MessagingMessagePayloadSizeBytesKey = kv.Key("messaging.message_payload_size_bytes")

	// The compressed size of the message payload in bytes.
	MessagingMessagePayloadCompressedSizeBytesKey = kv.Key("messaging.message_payload_compressed_size_bytes")

	// Identifies which part and kind of message consumption is being
	// preformed.
	MessagingOperationKey = kv.Key("messaging.operation")

	// RabbitMQ specific attribute describing the destination routing key.
	MessagingRabbitMQRoutingKeyKey = kv.Key("messaging.rabbitmq.routing_key")
)

var (
	MessagingDestinationKindKeyQueue = MessagingDestinationKindKey.String("queue")
	MessagingDestinationKindKeyTopic = MessagingDestinationKindKey.String("topic")

	MessagingTempDestination = MessagingTempDestinationKey.Bool(true)

	MessagingOperationReceive = MessagingOperationKey.String("receive")
	MessagingOperationProcess = MessagingOperationKey.String("process")
)

// Standard attribute keys for FaaS systems.
const (

	// Type of the trigger on which the function is executed.
	FaaSTriggerKey = kv.Key("faas.trigger")

	// String containing the execution identifier of the function.
	FaaSExecutionKey = kv.Key("faas.execution")

	// The name of the source on which the operation was performed.
	// For example, in Cloud Storage or S3 corresponds to the bucket name,
	// and in Cosmos DB to the database name.
	FaaSDocumentCollectionKey = kv.Key("faas.document.collection")

	// The type of the operation that was performed on the data.
	FaaSDocumentOperationKey = kv.Key("faas.document.operation")

	// A string containing the time when the data was accessed.
	FaaSDocumentTimeKey = kv.Key("faas.document.time")

	// The document name/table subjected to the operation.
	FaaSDocumentNameKey = kv.Key("faas.document.name")

	// The function invocation time.
	FaaSTimeKey = kv.Key("faas.time")

	// The schedule period as Cron Expression.
	FaaSCronKey = kv.Key("faas.cron")
)

var (
	FaasTriggerDatasource = FaaSTriggerKey.String("datasource")
	FaasTriggerHTTP       = FaaSTriggerKey.String("http")
	FaasTriggerPubSub     = FaaSTriggerKey.String("pubsub")
	FaasTriggerTimer      = FaaSTriggerKey.String("timer")
	FaasTriggerOther      = FaaSTriggerKey.String("other")

	FaaSDocumentOperationInsert = FaaSDocumentOperationKey.String("insert")
	FaaSDocumentOperationEdit   = FaaSDocumentOperationKey.String("edit")
	FaaSDocumentOperationDelete = FaaSDocumentOperationKey.String("delete")
)
