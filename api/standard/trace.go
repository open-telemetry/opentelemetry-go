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

// General attribute keys for spans.
const (
	// Service name of the remote service. Should equal the actual
	// `service.name` resource attribute of the remote service, if any.
	PeerServiceKey = kv.Key("peer.service")
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

	// The size of the request payload body in bytes.
	HTTPRequestContentLengthKey = kv.Key("http.request_content_length")

	// The size of the uncompressed request payload body after transport decoding.
	// Not set if transport encoding not used.
	HTTPRequestContentLengthUncompressedKey = kv.Key("http.request_content_length_uncompressed")

	// The size of the response payload body in bytes.
	HTTPResponseContentLengthKey = kv.Key("http.response_content_length")

	// The size of the uncompressed response payload body after transport decoding.
	// Not set if transport encoding not used.
	HTTPResponseContentLengthUncompressedKey = kv.Key("http.response_content_length_uncompressed")
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

// Standard attribute keys for database connections.
const (
	// Identifier for the database system (DBMS) being used.
	DBSystemKey = kv.Key("db.system")

	// Database Connection String with embedded credentials removed.
	DBConnectionStringKey = kv.Key("db.connection_string")

	// Username for accessing database.
	DBUserKey = kv.Key("db.user")
)

var (
	DBSystemDB2       = DBSystemKey.String("db2")        // IBM DB2
	DBSystemDerby     = DBSystemKey.String("derby")      // Apache Derby
	DBSystemHive      = DBSystemKey.String("hive")       // Apache Hive
	DBSystemMariaDB   = DBSystemKey.String("mariadb")    // MariaDB
	DBSystemMSSql     = DBSystemKey.String("mssql")      // Microsoft SQL Server
	DBSystemMySQL     = DBSystemKey.String("mysql")      // MySQL
	DBSystemOracle    = DBSystemKey.String("oracle")     // Oracle Database
	DBSystemPostgres  = DBSystemKey.String("postgresql") // PostgreSQL
	DBSystemSqlite    = DBSystemKey.String("sqlite")     // SQLite
	DBSystemTeradata  = DBSystemKey.String("teradata")   // Teradata
	DBSystemOtherSQL  = DBSystemKey.String("other_sql")  // Some other Sql database. Fallback only
	DBSystemCassandra = DBSystemKey.String("cassandra")  // Cassandra
	DBSystemCosmosDB  = DBSystemKey.String("cosmosdb")   // Microsoft Azure CosmosDB
	DBSystemCouchbase = DBSystemKey.String("couchbase")  // Couchbase
	DBSystemCouchDB   = DBSystemKey.String("couchdb")    // CouchDB
	DBSystemDynamoDB  = DBSystemKey.String("dynamodb")   // Amazon DynamoDB
	DBSystemHBase     = DBSystemKey.String("hbase")      // HBase
	DBSystemMongodb   = DBSystemKey.String("mongodb")    // MongoDB
	DBSystemNeo4j     = DBSystemKey.String("neo4j")      // Neo4j
	DBSystemRedis     = DBSystemKey.String("redis")      // Redis
)

// Standard attribute keys for database calls.
const (
	// Database instance name.
	DBNameKey = kv.Key("db.name")

	// A database statement for the given database type.
	DBStatementKey = kv.Key("db.statement")

	// A database operation for the given database type.
	DBOperationKey = kv.Key("db.operation")
)

// Database technology-specific attributes
const (
	// Name of the Cassandra keyspace accessed. Use instead of `db.name`.
	DBCassandraKeyspaceKey = kv.Key("db.cassandra.keyspace")

	// HBase namespace accessed. Use instead of `db.name`.
	DBHBaseNamespaceKey = kv.Key("db.hbase.namespace")

	// Index of Redis database accessed. Use instead of `db.name`.
	DBRedisDBIndexKey = kv.Key("db.redis.database_index")

	// Collection being accessed within the database in `db.name`.
	DBMongoDBCollectionKey = kv.Key("db.mongodb.collection")
)

// Standard attribute keys for RPC.
const (
	// A string identifying the remoting system.
	RPCSystemKey = kv.Key("rpc.system")

	// The full name of the service being called.
	RPCServiceKey = kv.Key("rpc.service")

	// The name of the method being called.
	RPCMethodKey = kv.Key("rpc.method")

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
	RPCSystemGRPC = RPCSystemKey.String("grpc")

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

	// A boolean indicating that the serverless function is executed
	// for the first time (aka cold start).
	FaaSColdstartKey = kv.Key("faas.coldstart")

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
