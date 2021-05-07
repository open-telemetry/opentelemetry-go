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

package v1_4_0

import (
	"go.opentelemetry.io/otel/semconv"
)

// Semantic conventions for attribute keys used for network related
// operations.
const (
	// Transport protocol used.
	NetTransportKey = semconv.NetTransportKey

	// Remote address of the peer.
	NetPeerIPKey = semconv.NetPeerIPKey

	// Remote port number.
	NetPeerPortKey = semconv.NetPeerPortKey

	// Remote hostname or similar.
	NetPeerNameKey = semconv.NetPeerNameKey

	// Local host IP. Useful in case of a multi-IP host.
	NetHostIPKey = semconv.NetHostIPKey

	// Local host port.
	NetHostPortKey = semconv.NetHostPortKey

	// Local hostname or similar.
	NetHostNameKey = semconv.NetHostNameKey
)

// Semantic conventions for common transport protocol attributes.
var (
	NetTransportTCP    = semconv.NetTransportTCP
	NetTransportUDP    = semconv.NetTransportUDP
	NetTransportIP     = semconv.NetTransportIP
	NetTransportUnix   = semconv.NetTransportUnix
	NetTransportPipe   = semconv.NetTransportPipe
	NetTransportInProc = semconv.NetTransportInProc
	NetTransportOther  = semconv.NetTransportOther
)

// General attribute keys for spans.
const (
	// Service name of the remote service. Should equal the actual
	// `service.name` resource attribute of the remote service, if any.
	PeerServiceKey = semconv.PeerServiceKey
)

// Semantic conventions for attribute keys used to identify an authorized
// user.
const (
	// Username or the client identifier extracted from the access token or
	// authorization header in the inbound request from outside the system.
	EnduserIDKey = semconv.EnduserIDKey

	// Actual or assumed role the client is making the request with.
	EnduserRoleKey = semconv.EnduserRoleKey

	// Scopes or granted authorities the client currently possesses.
	EnduserScopeKey = semconv.EnduserScopeKey
)

// Semantic conventions for attribute keys for HTTP.
const (
	// HTTP request method.
	HTTPMethodKey = semconv.HTTPMethodKey

	// Full HTTP request URL in the form:
	// scheme://host[:port]/path?query[#fragment].
	HTTPURLKey = semconv.HTTPURLKey

	// The full request target as passed in a HTTP request line or
	// equivalent, e.g. "/path/12314/?q=ddds#123".
	HTTPTargetKey = semconv.HTTPTargetKey

	// The value of the HTTP host header.
	HTTPHostKey = semconv.HTTPHostKey

	// The URI scheme identifying the used protocol.
	HTTPSchemeKey = semconv.HTTPSchemeKey

	// HTTP response status code.
	HTTPStatusCodeKey = semconv.HTTPStatusCodeKey

	// Kind of HTTP protocol used.
	HTTPFlavorKey = semconv.HTTPFlavorKey

	// Value of the HTTP User-Agent header sent by the client.
	HTTPUserAgentKey = semconv.HTTPUserAgentKey

	// The primary server name of the matched virtual host.
	HTTPServerNameKey = semconv.HTTPServerNameKey

	// The matched route served (path template). For example,
	// "/users/:userID?".
	HTTPRouteKey = semconv.HTTPRouteKey

	// The IP address of the original client behind all proxies, if known
	// (e.g. from X-Forwarded-For).
	HTTPClientIPKey = semconv.HTTPClientIPKey

	// The size of the request payload body in bytes.
	HTTPRequestContentLengthKey = semconv.HTTPRequestContentLengthKey

	// The size of the uncompressed request payload body after transport decoding.
	// Not set if transport encoding not used.
	HTTPRequestContentLengthUncompressedKey = semconv.HTTPRequestContentLengthUncompressedKey

	// The size of the response payload body in bytes.
	HTTPResponseContentLengthKey = semconv.HTTPResponseContentLengthKey

	// The size of the uncompressed response payload body after transport decoding.
	// Not set if transport encoding not used.
	HTTPResponseContentLengthUncompressedKey = semconv.HTTPResponseContentLengthUncompressedKey
)

// Semantic conventions for common HTTP attributes.
var (
	// Semantic conventions for HTTP(S) URI schemes.
	HTTPSchemeHTTP  = semconv.HTTPSchemeHTTP
	HTTPSchemeHTTPS = semconv.HTTPSchemeHTTPS

	// Semantic conventions for HTTP protocols.
	HTTPFlavor1_0  = semconv.HTTPFlavor1_0
	HTTPFlavor1_1  = semconv.HTTPFlavor1_1
	HTTPFlavor2    = semconv.HTTPFlavor2
	HTTPFlavorSPDY = semconv.HTTPFlavorSPDY
	HTTPFlavorQUIC = semconv.HTTPFlavorQUIC
)

// Semantic conventions for attribute keys for database connections.
const (
	// Identifier for the database system (DBMS) being used.
	DBSystemKey = semconv.DBSystemKey

	// Database Connection String with embedded credentials removed.
	DBConnectionStringKey = semconv.DBConnectionStringKey

	// Username for accessing database.
	DBUserKey = semconv.DBUserKey
)

// Semantic conventions for common database system attributes.
var (
	DBSystemDB2       = semconv.DBSystemDB2
	DBSystemDerby     = semconv.DBSystemDerby
	DBSystemHive      = semconv.DBSystemHive
	DBSystemMariaDB   = semconv.DBSystemMariaDB
	DBSystemMSSql     = semconv.DBSystemMSSql
	DBSystemMySQL     = semconv.DBSystemMySQL
	DBSystemOracle    = semconv.DBSystemOracle
	DBSystemPostgres  = semconv.DBSystemPostgres
	DBSystemSqlite    = semconv.DBSystemSqlite
	DBSystemTeradata  = semconv.DBSystemTeradata
	DBSystemOtherSQL  = semconv.DBSystemOtherSQL
	DBSystemCassandra = semconv.DBSystemCassandra
	DBSystemCosmosDB  = semconv.DBSystemCosmosDB
	DBSystemCouchbase = semconv.DBSystemCouchbase
	DBSystemCouchDB   = semconv.DBSystemCouchDB
	DBSystemDynamoDB  = semconv.DBSystemDynamoDB
	DBSystemHBase     = semconv.DBSystemHBase
	DBSystemMongodb   = semconv.DBSystemMongodb
	DBSystemNeo4j     = semconv.DBSystemNeo4j
	DBSystemRedis     = semconv.DBSystemRedis
)

// Semantic conventions for attribute keys for database calls.
const (
	// Database instance name.
	DBNameKey = semconv.DBNameKey

	// A database statement for the given database type.
	DBStatementKey = semconv.DBStatementKey

	// A database operation for the given database type.
	DBOperationKey = semconv.DBOperationKey
)

// Database technology-specific attributes
const (
	// Name of the Cassandra keyspace accessed. Use instead of `db.name`.
	DBCassandraKeyspaceKey = semconv.DBCassandraKeyspaceKey

	// HBase namespace accessed. Use instead of `db.name`.
	DBHBaseNamespaceKey = semconv.DBHBaseNamespaceKey

	// Index of Redis database accessed. Use instead of `db.name`.
	DBRedisDBIndexKey = semconv.DBRedisDBIndexKey

	// Collection being accessed within the database in `db.name`.
	DBMongoDBCollectionKey = semconv.DBMongoDBCollectionKey
)

// Semantic conventions for attribute keys for RPC.
const (
	// A string identifying the remoting system.
	RPCSystemKey = semconv.RPCSystemKey

	// The full name of the service being called.
	RPCServiceKey = semconv.RPCServiceKey

	// The name of the method being called.
	RPCMethodKey = semconv.RPCMethodKey

	// Name of message transmitted or received.
	RPCNameKey = semconv.RPCNameKey

	// Type of message transmitted or received.
	RPCMessageTypeKey = semconv.RPCMessageTypeKey

	// Identifier of message transmitted or received.
	RPCMessageIDKey = semconv.RPCMessageIDKey

	// The compressed size of the message transmitted or received in bytes.
	RPCMessageCompressedSizeKey = semconv.RPCMessageCompressedSizeKey

	// The uncompressed size of the message transmitted or received in
	// bytes.
	RPCMessageUncompressedSizeKey = semconv.RPCMessageUncompressedSizeKey
)

// Semantic conventions for common RPC attributes.
var (
	// Semantic convention for gRPC as the remoting system.
	RPCSystemGRPC = semconv.RPCSystemGRPC

	// Semantic convention for a message named message.
	RPCNameMessage = semconv.RPCNameMessage

	// Semantic conventions for RPC message types.
	RPCMessageTypeSent     = semconv.RPCMessageTypeSent
	RPCMessageTypeReceived = semconv.RPCMessageTypeReceived
)

// Semantic conventions for attribute keys for messaging systems.
const (
	// A unique identifier describing the messaging system. For example,
	// kafka, rabbitmq or activemq.
	MessagingSystemKey = semconv.MessagingSystemKey

	// The message destination name, e.g. MyQueue or MyTopic.
	MessagingDestinationKey = semconv.MessagingDestinationKey

	// The kind of message destination.
	MessagingDestinationKindKey = semconv.MessagingDestinationKindKey

	// Describes if the destination is temporary or not.
	MessagingTempDestinationKey = semconv.MessagingTempDestinationKey

	// The name of the transport protocol.
	MessagingProtocolKey = semconv.MessagingProtocolKey

	// The version of the transport protocol.
	MessagingProtocolVersionKey = semconv.MessagingProtocolVersionKey

	// Messaging service URL.
	MessagingURLKey = semconv.MessagingURLKey

	// Identifier used by the messaging system for a message.
	MessagingMessageIDKey = semconv.MessagingMessageIDKey

	// Identifier used by the messaging system for a conversation.
	MessagingConversationIDKey = semconv.MessagingConversationIDKey

	// The (uncompressed) size of the message payload in bytes.
	MessagingMessagePayloadSizeBytesKey = semconv.MessagingMessagePayloadSizeBytesKey

	// The compressed size of the message payload in bytes.
	MessagingMessagePayloadCompressedSizeBytesKey = semconv.MessagingMessagePayloadCompressedSizeBytesKey

	// Identifies which part and kind of message consumption is being
	// preformed.
	MessagingOperationKey = semconv.MessagingOperationKey

	// RabbitMQ specific attribute describing the destination routing key.
	MessagingRabbitMQRoutingKeyKey = semconv.MessagingRabbitMQRoutingKeyKey
)

// Semantic conventions for common messaging system attributes.
var (
	// Semantic conventions for message destinations.
	MessagingDestinationKindKeyQueue = semconv.MessagingDestinationKindKeyQueue
	MessagingDestinationKindKeyTopic = semconv.MessagingDestinationKindKeyTopic

	// Semantic convention for message destinations that are temporary.
	MessagingTempDestination = semconv.MessagingTempDestination

	// Semantic convention for the operation parts of message consumption.
	// This does not include a "send" attribute as that is explicitly not
	// allowed in the OpenTelemetry specification.
	MessagingOperationReceive = semconv.MessagingOperationReceive
	MessagingOperationProcess = semconv.MessagingOperationProcess
)

// Semantic conventions for attribute keys for FaaS systems.
const (

	// Type of the trigger on which the function is executed.
	FaaSTriggerKey = semconv.FaaSTriggerKey

	// String containing the execution identifier of the function.
	FaaSExecutionKey = semconv.FaaSExecutionKey

	// A boolean indicating that the serverless function is executed
	// for the first time (aka cold start).
	FaaSColdstartKey = semconv.FaaSColdstartKey

	// The name of the source on which the operation was performed.
	// For example, in Cloud Storage or S3 corresponds to the bucket name,
	// and in Cosmos DB to the database name.
	FaaSDocumentCollectionKey = semconv.FaaSDocumentCollectionKey

	// The type of the operation that was performed on the data.
	FaaSDocumentOperationKey = semconv.FaaSDocumentOperationKey

	// A string containing the time when the data was accessed.
	FaaSDocumentTimeKey = semconv.FaaSDocumentTimeKey

	// The document name/table subjected to the operation.
	FaaSDocumentNameKey = semconv.FaaSDocumentNameKey

	// The function invocation time.
	FaaSTimeKey = semconv.FaaSTimeKey

	// The schedule period as Cron Expression.
	FaaSCronKey = semconv.FaaSCronKey
)

// Semantic conventions for common FaaS system attributes.
var (
	// Semantic conventions for the types of triggers.
	FaasTriggerDatasource = semconv.FaasTriggerDatasource
	FaasTriggerHTTP       = semconv.FaasTriggerHTTP
	FaasTriggerPubSub     = semconv.FaasTriggerPubSub
	FaasTriggerTimer      = semconv.FaasTriggerTimer
	FaasTriggerOther      = semconv.FaasTriggerOther

	// Semantic conventions for the types of operations performed.
	FaaSDocumentOperationInsert = semconv.FaaSDocumentOperationInsert
	FaaSDocumentOperationEdit   = semconv.FaaSDocumentOperationEdit
	FaaSDocumentOperationDelete = semconv.FaaSDocumentOperationDelete
)

// Semantic conventions for source code attributes.
const (
	// The method or function name, or equivalent (usually rightmost part of
	// the code unit's name).
	CodeFunctionKey = semconv.CodeFunctionKey

	// The "namespace" within which `code.function` is defined. Usually the
	// qualified class or module name, such that
	// `code.namespace` + some separator + `code.function` form a unique
	// identifier for the code unit.
	CodeNamespaceKey = semconv.CodeNamespaceKey

	// The source code file name that identifies the code unit as uniquely as
	// possible (preferably an absolute file path).
	CodeFilepathKey = semconv.CodeFilepathKey

	// The line number in `code.filepath` best representing the operation.
	// It SHOULD point within the code unit named in `code.function`.
	CodeLineNumberKey = semconv.CodeLineNumberKey
)
