// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.27.0"

const (
  // ASPNETCoreDiagnosticsExceptions is the metric conforming to the
  // "aspnetcore.diagnostics.exceptions" semantic conventions. It represents the
  // number of exceptions caught by exception handling middleware.
  // Instrument: counter
  // Unit: {exception}
  // Stability: stable
  ASPNETCoreDiagnosticsExceptionsName = "aspnetcore.diagnostics.exceptions"
  ASPNETCoreDiagnosticsExceptionsUnit = "{exception}"
  ASPNETCoreDiagnosticsExceptionsDescription = "Number of exceptions caught by exception handling middleware."
  // ASPNETCoreRateLimitingActiveRequestLeases is the metric conforming to the
  // "aspnetcore.rate_limiting.active_request_leases" semantic conventions. It
  // represents the number of requests that are currently active on the server
  // that hold a rate limiting lease.
  // Instrument: updowncounter
  // Unit: {request}
  // Stability: stable
  ASPNETCoreRateLimitingActiveRequestLeasesName = "aspnetcore.rate_limiting.active_request_leases"
  ASPNETCoreRateLimitingActiveRequestLeasesUnit = "{request}"
  ASPNETCoreRateLimitingActiveRequestLeasesDescription = "Number of requests that are currently active on the server that hold a rate limiting lease."
  // ASPNETCoreRateLimitingQueuedRequests is the metric conforming to the
  // "aspnetcore.rate_limiting.queued_requests" semantic conventions. It
  // represents the number of requests that are currently queued, waiting to
  // acquire a rate limiting lease.
  // Instrument: updowncounter
  // Unit: {request}
  // Stability: stable
  ASPNETCoreRateLimitingQueuedRequestsName = "aspnetcore.rate_limiting.queued_requests"
  ASPNETCoreRateLimitingQueuedRequestsUnit = "{request}"
  ASPNETCoreRateLimitingQueuedRequestsDescription = "Number of requests that are currently queued, waiting to acquire a rate limiting lease."
  // ASPNETCoreRateLimitingRequestTimeInQueue is the metric conforming to the
  // "aspnetcore.rate_limiting.request.time_in_queue" semantic conventions. It
  // represents the time the request spent in a queue waiting to acquire a rate
  // limiting lease.
  // Instrument: histogram
  // Unit: s
  // Stability: stable
  ASPNETCoreRateLimitingRequestTimeInQueueName = "aspnetcore.rate_limiting.request.time_in_queue"
  ASPNETCoreRateLimitingRequestTimeInQueueUnit = "s"
  ASPNETCoreRateLimitingRequestTimeInQueueDescription = "The time the request spent in a queue waiting to acquire a rate limiting lease."
  // ASPNETCoreRateLimitingRequestLeaseDuration is the metric conforming to the
  // "aspnetcore.rate_limiting.request_lease.duration" semantic conventions. It
  // represents the duration of rate limiting lease held by requests on the
  // server.
  // Instrument: histogram
  // Unit: s
  // Stability: stable
  ASPNETCoreRateLimitingRequestLeaseDurationName = "aspnetcore.rate_limiting.request_lease.duration"
  ASPNETCoreRateLimitingRequestLeaseDurationUnit = "s"
  ASPNETCoreRateLimitingRequestLeaseDurationDescription = "The duration of rate limiting lease held by requests on the server."
  // ASPNETCoreRateLimitingRequests is the metric conforming to the
  // "aspnetcore.rate_limiting.requests" semantic conventions. It represents the
  // number of requests that tried to acquire a rate limiting lease.
  // Instrument: counter
  // Unit: {request}
  // Stability: stable
  ASPNETCoreRateLimitingRequestsName = "aspnetcore.rate_limiting.requests"
  ASPNETCoreRateLimitingRequestsUnit = "{request}"
  ASPNETCoreRateLimitingRequestsDescription = "Number of requests that tried to acquire a rate limiting lease."
  // ASPNETCoreRoutingMatchAttempts is the metric conforming to the
  // "aspnetcore.routing.match_attempts" semantic conventions. It represents the
  // number of requests that were attempted to be matched to an endpoint.
  // Instrument: counter
  // Unit: {match_attempt}
  // Stability: stable
  ASPNETCoreRoutingMatchAttemptsName = "aspnetcore.routing.match_attempts"
  ASPNETCoreRoutingMatchAttemptsUnit = "{match_attempt}"
  ASPNETCoreRoutingMatchAttemptsDescription = "Number of requests that were attempted to be matched to an endpoint."
  // ContainerCPUTime is the metric conforming to the "container.cpu.time"
  // semantic conventions. It represents the total CPU time consumed.
  // Instrument: counter
  // Unit: s
  // Stability: development
  ContainerCPUTimeName = "container.cpu.time"
  ContainerCPUTimeUnit = "s"
  ContainerCPUTimeDescription = "Total CPU time consumed"
  // ContainerDiskIo is the metric conforming to the "container.disk.io" semantic
  // conventions. It represents the disk bytes for the container.
  // Instrument: counter
  // Unit: By
  // Stability: development
  ContainerDiskIoName = "container.disk.io"
  ContainerDiskIoUnit = "By"
  ContainerDiskIoDescription = "Disk bytes for the container."
  // ContainerMemoryUsage is the metric conforming to the
  // "container.memory.usage" semantic conventions. It represents the memory
  // usage of the container.
  // Instrument: counter
  // Unit: By
  // Stability: development
  ContainerMemoryUsageName = "container.memory.usage"
  ContainerMemoryUsageUnit = "By"
  ContainerMemoryUsageDescription = "Memory usage of the container."
  // ContainerNetworkIo is the metric conforming to the "container.network.io"
  // semantic conventions. It represents the network bytes for the container.
  // Instrument: counter
  // Unit: By
  // Stability: development
  ContainerNetworkIoName = "container.network.io"
  ContainerNetworkIoUnit = "By"
  ContainerNetworkIoDescription = "Network bytes for the container."
  // DBClientConnectionCount is the metric conforming to the
  // "db.client.connection.count" semantic conventions. It represents the number
  // of connections that are currently in state described by the `state`
  // attribute.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  DBClientConnectionCountName = "db.client.connection.count"
  DBClientConnectionCountUnit = "{connection}"
  DBClientConnectionCountDescription = "The number of connections that are currently in state described by the `state` attribute"
  // DBClientConnectionCreateTime is the metric conforming to the
  // "db.client.connection.create_time" semantic conventions. It represents the
  // time it took to create a new connection.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  DBClientConnectionCreateTimeName = "db.client.connection.create_time"
  DBClientConnectionCreateTimeUnit = "s"
  DBClientConnectionCreateTimeDescription = "The time it took to create a new connection"
  // DBClientConnectionIdleMax is the metric conforming to the
  // "db.client.connection.idle.max" semantic conventions. It represents the
  // maximum number of idle open connections allowed.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  DBClientConnectionIdleMaxName = "db.client.connection.idle.max"
  DBClientConnectionIdleMaxUnit = "{connection}"
  DBClientConnectionIdleMaxDescription = "The maximum number of idle open connections allowed"
  // DBClientConnectionIdleMin is the metric conforming to the
  // "db.client.connection.idle.min" semantic conventions. It represents the
  // minimum number of idle open connections allowed.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  DBClientConnectionIdleMinName = "db.client.connection.idle.min"
  DBClientConnectionIdleMinUnit = "{connection}"
  DBClientConnectionIdleMinDescription = "The minimum number of idle open connections allowed"
  // DBClientConnectionMax is the metric conforming to the
  // "db.client.connection.max" semantic conventions. It represents the maximum
  // number of open connections allowed.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  DBClientConnectionMaxName = "db.client.connection.max"
  DBClientConnectionMaxUnit = "{connection}"
  DBClientConnectionMaxDescription = "The maximum number of open connections allowed"
  // DBClientConnectionPendingRequests is the metric conforming to the
  // "db.client.connection.pending_requests" semantic conventions. It represents
  // the number of pending requests for an open connection, cumulative for the
  // entire pool.
  // Instrument: updowncounter
  // Unit: {request}
  // Stability: development
  DBClientConnectionPendingRequestsName = "db.client.connection.pending_requests"
  DBClientConnectionPendingRequestsUnit = "{request}"
  DBClientConnectionPendingRequestsDescription = "The number of pending requests for an open connection, cumulative for the entire pool"
  // DBClientConnectionTimeouts is the metric conforming to the
  // "db.client.connection.timeouts" semantic conventions. It represents the
  // number of connection timeouts that have occurred trying to obtain a
  // connection from the pool.
  // Instrument: counter
  // Unit: {timeout}
  // Stability: development
  DBClientConnectionTimeoutsName = "db.client.connection.timeouts"
  DBClientConnectionTimeoutsUnit = "{timeout}"
  DBClientConnectionTimeoutsDescription = "The number of connection timeouts that have occurred trying to obtain a connection from the pool"
  // DBClientConnectionUseTime is the metric conforming to the
  // "db.client.connection.use_time" semantic conventions. It represents the time
  // between borrowing a connection and returning it to the pool.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  DBClientConnectionUseTimeName = "db.client.connection.use_time"
  DBClientConnectionUseTimeUnit = "s"
  DBClientConnectionUseTimeDescription = "The time between borrowing a connection and returning it to the pool"
  // DBClientConnectionWaitTime is the metric conforming to the
  // "db.client.connection.wait_time" semantic conventions. It represents the
  // time it took to obtain an open connection from the pool.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  DBClientConnectionWaitTimeName = "db.client.connection.wait_time"
  DBClientConnectionWaitTimeUnit = "s"
  DBClientConnectionWaitTimeDescription = "The time it took to obtain an open connection from the pool"
  // DBClientConnectionsCreateTime is the metric conforming to the
  // "db.client.connections.create_time" semantic conventions. It represents the
  // deprecated, use `db.client.connection.create_time` instead. Note: the unit
  // also changed from `ms` to `s`.
  // Instrument: histogram
  // Unit: ms
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.create_time`. Note: the unit also changed from `ms` to `s`.
  DBClientConnectionsCreateTimeName = "db.client.connections.create_time"
  DBClientConnectionsCreateTimeUnit = "ms"
  DBClientConnectionsCreateTimeDescription = "Deprecated, use `db.client.connection.create_time` instead. Note: the unit also changed from `ms` to `s`."
  // DBClientConnectionsIdleMax is the metric conforming to the
  // "db.client.connections.idle.max" semantic conventions. It represents the
  // deprecated, use `db.client.connection.idle.max` instead.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.idle.max`.
  DBClientConnectionsIdleMaxName = "db.client.connections.idle.max"
  DBClientConnectionsIdleMaxUnit = "{connection}"
  DBClientConnectionsIdleMaxDescription = "Deprecated, use `db.client.connection.idle.max` instead."
  // DBClientConnectionsIdleMin is the metric conforming to the
  // "db.client.connections.idle.min" semantic conventions. It represents the
  // deprecated, use `db.client.connection.idle.min` instead.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.idle.min`.
  DBClientConnectionsIdleMinName = "db.client.connections.idle.min"
  DBClientConnectionsIdleMinUnit = "{connection}"
  DBClientConnectionsIdleMinDescription = "Deprecated, use `db.client.connection.idle.min` instead."
  // DBClientConnectionsMax is the metric conforming to the
  // "db.client.connections.max" semantic conventions. It represents the
  // deprecated, use `db.client.connection.max` instead.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.max`.
  DBClientConnectionsMaxName = "db.client.connections.max"
  DBClientConnectionsMaxUnit = "{connection}"
  DBClientConnectionsMaxDescription = "Deprecated, use `db.client.connection.max` instead."
  // DBClientConnectionsPendingRequests is the metric conforming to the
  // "db.client.connections.pending_requests" semantic conventions. It represents
  // the deprecated, use `db.client.connection.pending_requests` instead.
  // Instrument: updowncounter
  // Unit: {request}
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.pending_requests`.
  DBClientConnectionsPendingRequestsName = "db.client.connections.pending_requests"
  DBClientConnectionsPendingRequestsUnit = "{request}"
  DBClientConnectionsPendingRequestsDescription = "Deprecated, use `db.client.connection.pending_requests` instead."
  // DBClientConnectionsTimeouts is the metric conforming to the
  // "db.client.connections.timeouts" semantic conventions. It represents the
  // deprecated, use `db.client.connection.timeouts` instead.
  // Instrument: counter
  // Unit: {timeout}
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.timeouts`.
  DBClientConnectionsTimeoutsName = "db.client.connections.timeouts"
  DBClientConnectionsTimeoutsUnit = "{timeout}"
  DBClientConnectionsTimeoutsDescription = "Deprecated, use `db.client.connection.timeouts` instead."
  // DBClientConnectionsUsage is the metric conforming to the
  // "db.client.connections.usage" semantic conventions. It represents the
  // deprecated, use `db.client.connection.count` instead.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.count`.
  DBClientConnectionsUsageName = "db.client.connections.usage"
  DBClientConnectionsUsageUnit = "{connection}"
  DBClientConnectionsUsageDescription = "Deprecated, use `db.client.connection.count` instead."
  // DBClientConnectionsUseTime is the metric conforming to the
  // "db.client.connections.use_time" semantic conventions. It represents the
  // deprecated, use `db.client.connection.use_time` instead. Note: the unit also
  // changed from `ms` to `s`.
  // Instrument: histogram
  // Unit: ms
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.use_time`. Note: the unit also changed from `ms` to `s`.
  DBClientConnectionsUseTimeName = "db.client.connections.use_time"
  DBClientConnectionsUseTimeUnit = "ms"
  DBClientConnectionsUseTimeDescription = "Deprecated, use `db.client.connection.use_time` instead. Note: the unit also changed from `ms` to `s`."
  // DBClientConnectionsWaitTime is the metric conforming to the
  // "db.client.connections.wait_time" semantic conventions. It represents the
  // deprecated, use `db.client.connection.wait_time` instead. Note: the unit
  // also changed from `ms` to `s`.
  // Instrument: histogram
  // Unit: ms
  // Stability: development
  // Deprecated: Replaced by `db.client.connection.wait_time`. Note: the unit also changed from `ms` to `s`.
  DBClientConnectionsWaitTimeName = "db.client.connections.wait_time"
  DBClientConnectionsWaitTimeUnit = "ms"
  DBClientConnectionsWaitTimeDescription = "Deprecated, use `db.client.connection.wait_time` instead. Note: the unit also changed from `ms` to `s`."
  // DBClientOperationDuration is the metric conforming to the
  // "db.client.operation.duration" semantic conventions. It represents the
  // duration of database client operations.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  DBClientOperationDurationName = "db.client.operation.duration"
  DBClientOperationDurationUnit = "s"
  DBClientOperationDurationDescription = "Duration of database client operations."
  // DNSLookupDuration is the metric conforming to the "dns.lookup.duration"
  // semantic conventions. It represents the measures the time taken to perform a
  // DNS lookup.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  DNSLookupDurationName = "dns.lookup.duration"
  DNSLookupDurationUnit = "s"
  DNSLookupDurationDescription = "Measures the time taken to perform a DNS lookup."
  // FaaSColdstarts is the metric conforming to the "faas.coldstarts" semantic
  // conventions. It represents the number of invocation cold starts.
  // Instrument: counter
  // Unit: {coldstart}
  // Stability: development
  FaaSColdstartsName = "faas.coldstarts"
  FaaSColdstartsUnit = "{coldstart}"
  FaaSColdstartsDescription = "Number of invocation cold starts"
  // FaaSCPUUsage is the metric conforming to the "faas.cpu_usage" semantic
  // conventions. It represents the distribution of CPU usage per invocation.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  FaaSCPUUsageName = "faas.cpu_usage"
  FaaSCPUUsageUnit = "s"
  FaaSCPUUsageDescription = "Distribution of CPU usage per invocation"
  // FaaSErrors is the metric conforming to the "faas.errors" semantic
  // conventions. It represents the number of invocation errors.
  // Instrument: counter
  // Unit: {error}
  // Stability: development
  FaaSErrorsName = "faas.errors"
  FaaSErrorsUnit = "{error}"
  FaaSErrorsDescription = "Number of invocation errors"
  // FaaSInitDuration is the metric conforming to the "faas.init_duration"
  // semantic conventions. It represents the measures the duration of the
  // function's initialization, such as a cold start.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  FaaSInitDurationName = "faas.init_duration"
  FaaSInitDurationUnit = "s"
  FaaSInitDurationDescription = "Measures the duration of the function's initialization, such as a cold start"
  // FaaSInvocations is the metric conforming to the "faas.invocations" semantic
  // conventions. It represents the number of successful invocations.
  // Instrument: counter
  // Unit: {invocation}
  // Stability: development
  FaaSInvocationsName = "faas.invocations"
  FaaSInvocationsUnit = "{invocation}"
  FaaSInvocationsDescription = "Number of successful invocations"
  // FaaSInvokeDuration is the metric conforming to the "faas.invoke_duration"
  // semantic conventions. It represents the measures the duration of the
  // function's logic execution.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  FaaSInvokeDurationName = "faas.invoke_duration"
  FaaSInvokeDurationUnit = "s"
  FaaSInvokeDurationDescription = "Measures the duration of the function's logic execution"
  // FaaSMemUsage is the metric conforming to the "faas.mem_usage" semantic
  // conventions. It represents the distribution of max memory usage per
  // invocation.
  // Instrument: histogram
  // Unit: By
  // Stability: development
  FaaSMemUsageName = "faas.mem_usage"
  FaaSMemUsageUnit = "By"
  FaaSMemUsageDescription = "Distribution of max memory usage per invocation"
  // FaaSNetIo is the metric conforming to the "faas.net_io" semantic
  // conventions. It represents the distribution of net I/O usage per invocation.
  // Instrument: histogram
  // Unit: By
  // Stability: development
  FaaSNetIoName = "faas.net_io"
  FaaSNetIoUnit = "By"
  FaaSNetIoDescription = "Distribution of net I/O usage per invocation"
  // FaaSTimeouts is the metric conforming to the "faas.timeouts" semantic
  // conventions. It represents the number of invocation timeouts.
  // Instrument: counter
  // Unit: {timeout}
  // Stability: development
  FaaSTimeoutsName = "faas.timeouts"
  FaaSTimeoutsUnit = "{timeout}"
  FaaSTimeoutsDescription = "Number of invocation timeouts"
  // GenAIClientOperationDuration is the metric conforming to the
  // "gen_ai.client.operation.duration" semantic conventions. It represents the
  // genAI operation duration.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  GenAIClientOperationDurationName = "gen_ai.client.operation.duration"
  GenAIClientOperationDurationUnit = "s"
  GenAIClientOperationDurationDescription = "GenAI operation duration"
  // GenAIClientTokenUsage is the metric conforming to the
  // "gen_ai.client.token.usage" semantic conventions. It represents the measures
  // number of input and output tokens used.
  // Instrument: histogram
  // Unit: {token}
  // Stability: development
  GenAIClientTokenUsageName = "gen_ai.client.token.usage"
  GenAIClientTokenUsageUnit = "{token}"
  GenAIClientTokenUsageDescription = "Measures number of input and output tokens used"
  // GenAIServerRequestDuration is the metric conforming to the
  // "gen_ai.server.request.duration" semantic conventions. It represents the
  // generative AI server request duration such as time-to-last byte or last
  // output token.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  GenAIServerRequestDurationName = "gen_ai.server.request.duration"
  GenAIServerRequestDurationUnit = "s"
  GenAIServerRequestDurationDescription = "Generative AI server request duration such as time-to-last byte or last output token"
  // GenAIServerTimePerOutputToken is the metric conforming to the
  // "gen_ai.server.time_per_output_token" semantic conventions. It represents
  // the time per output token generated after the first token for successful
  // responses.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  GenAIServerTimePerOutputTokenName = "gen_ai.server.time_per_output_token"
  GenAIServerTimePerOutputTokenUnit = "s"
  GenAIServerTimePerOutputTokenDescription = "Time per output token generated after the first token for successful responses"
  // GenAIServerTimeToFirstToken is the metric conforming to the
  // "gen_ai.server.time_to_first_token" semantic conventions. It represents the
  // time to generate first token for successful responses.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  GenAIServerTimeToFirstTokenName = "gen_ai.server.time_to_first_token"
  GenAIServerTimeToFirstTokenUnit = "s"
  GenAIServerTimeToFirstTokenDescription = "Time to generate first token for successful responses"
  // GoConfigGogc is the metric conforming to the "go.config.gogc" semantic
  // conventions. It represents the heap size target percentage configured by the
  // user, otherwise 100.
  // Instrument: updowncounter
  // Unit: %
  // Stability: development
  GoConfigGogcName = "go.config.gogc"
  GoConfigGogcUnit = "%"
  GoConfigGogcDescription = "Heap size target percentage configured by the user, otherwise 100."
  // GoGoroutineCount is the metric conforming to the "go.goroutine.count"
  // semantic conventions. It represents the count of live goroutines.
  // Instrument: updowncounter
  // Unit: {goroutine}
  // Stability: development
  GoGoroutineCountName = "go.goroutine.count"
  GoGoroutineCountUnit = "{goroutine}"
  GoGoroutineCountDescription = "Count of live goroutines."
  // GoMemoryAllocated is the metric conforming to the "go.memory.allocated"
  // semantic conventions. It represents the memory allocated to the heap by the
  // application.
  // Instrument: counter
  // Unit: By
  // Stability: development
  GoMemoryAllocatedName = "go.memory.allocated"
  GoMemoryAllocatedUnit = "By"
  GoMemoryAllocatedDescription = "Memory allocated to the heap by the application."
  // GoMemoryAllocations is the metric conforming to the "go.memory.allocations"
  // semantic conventions. It represents the count of allocations to the heap by
  // the application.
  // Instrument: counter
  // Unit: {allocation}
  // Stability: development
  GoMemoryAllocationsName = "go.memory.allocations"
  GoMemoryAllocationsUnit = "{allocation}"
  GoMemoryAllocationsDescription = "Count of allocations to the heap by the application."
  // GoMemoryGCGoal is the metric conforming to the "go.memory.gc.goal" semantic
  // conventions. It represents the heap size target for the end of the GC cycle.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  GoMemoryGCGoalName = "go.memory.gc.goal"
  GoMemoryGCGoalUnit = "By"
  GoMemoryGCGoalDescription = "Heap size target for the end of the GC cycle."
  // GoMemoryLimit is the metric conforming to the "go.memory.limit" semantic
  // conventions. It represents the go runtime memory limit configured by the
  // user, if a limit exists.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  GoMemoryLimitName = "go.memory.limit"
  GoMemoryLimitUnit = "By"
  GoMemoryLimitDescription = "Go runtime memory limit configured by the user, if a limit exists."
  // GoMemoryUsed is the metric conforming to the "go.memory.used" semantic
  // conventions. It represents the memory used by the Go runtime.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  GoMemoryUsedName = "go.memory.used"
  GoMemoryUsedUnit = "By"
  GoMemoryUsedDescription = "Memory used by the Go runtime."
  // GoProcessorLimit is the metric conforming to the "go.processor.limit"
  // semantic conventions. It represents the number of OS threads that can
  // execute user-level Go code simultaneously.
  // Instrument: updowncounter
  // Unit: {thread}
  // Stability: development
  GoProcessorLimitName = "go.processor.limit"
  GoProcessorLimitUnit = "{thread}"
  GoProcessorLimitDescription = "The number of OS threads that can execute user-level Go code simultaneously."
  // GoScheduleDuration is the metric conforming to the "go.schedule.duration"
  // semantic conventions. It represents the time goroutines have spent in the
  // scheduler in a runnable state before actually running.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  GoScheduleDurationName = "go.schedule.duration"
  GoScheduleDurationUnit = "s"
  GoScheduleDurationDescription = "The time goroutines have spent in the scheduler in a runnable state before actually running."
  // HTTPClientActiveRequests is the metric conforming to the
  // "http.client.active_requests" semantic conventions. It represents the number
  // of active HTTP requests.
  // Instrument: updowncounter
  // Unit: {request}
  // Stability: development
  HTTPClientActiveRequestsName = "http.client.active_requests"
  HTTPClientActiveRequestsUnit = "{request}"
  HTTPClientActiveRequestsDescription = "Number of active HTTP requests."
  // HTTPClientConnectionDuration is the metric conforming to the
  // "http.client.connection.duration" semantic conventions. It represents the
  // duration of the successfully established outbound HTTP connections.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  HTTPClientConnectionDurationName = "http.client.connection.duration"
  HTTPClientConnectionDurationUnit = "s"
  HTTPClientConnectionDurationDescription = "The duration of the successfully established outbound HTTP connections."
  // HTTPClientOpenConnections is the metric conforming to the
  // "http.client.open_connections" semantic conventions. It represents the
  // number of outbound HTTP connections that are currently active or idle on the
  // client.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  HTTPClientOpenConnectionsName = "http.client.open_connections"
  HTTPClientOpenConnectionsUnit = "{connection}"
  HTTPClientOpenConnectionsDescription = "Number of outbound HTTP connections that are currently active or idle on the client."
  // HTTPClientRequestBodySize is the metric conforming to the
  // "http.client.request.body.size" semantic conventions. It represents the size
  // of HTTP client request bodies.
  // Instrument: histogram
  // Unit: By
  // Stability: development
  HTTPClientRequestBodySizeName = "http.client.request.body.size"
  HTTPClientRequestBodySizeUnit = "By"
  HTTPClientRequestBodySizeDescription = "Size of HTTP client request bodies."
  // HTTPClientRequestDuration is the metric conforming to the
  // "http.client.request.duration" semantic conventions. It represents the
  // duration of HTTP client requests.
  // Instrument: histogram
  // Unit: s
  // Stability: stable
  HTTPClientRequestDurationName = "http.client.request.duration"
  HTTPClientRequestDurationUnit = "s"
  HTTPClientRequestDurationDescription = "Duration of HTTP client requests."
  // HTTPClientResponseBodySize is the metric conforming to the
  // "http.client.response.body.size" semantic conventions. It represents the
  // size of HTTP client response bodies.
  // Instrument: histogram
  // Unit: By
  // Stability: development
  HTTPClientResponseBodySizeName = "http.client.response.body.size"
  HTTPClientResponseBodySizeUnit = "By"
  HTTPClientResponseBodySizeDescription = "Size of HTTP client response bodies."
  // HTTPServerActiveRequests is the metric conforming to the
  // "http.server.active_requests" semantic conventions. It represents the number
  // of active HTTP server requests.
  // Instrument: updowncounter
  // Unit: {request}
  // Stability: development
  HTTPServerActiveRequestsName = "http.server.active_requests"
  HTTPServerActiveRequestsUnit = "{request}"
  HTTPServerActiveRequestsDescription = "Number of active HTTP server requests."
  // HTTPServerRequestBodySize is the metric conforming to the
  // "http.server.request.body.size" semantic conventions. It represents the size
  // of HTTP server request bodies.
  // Instrument: histogram
  // Unit: By
  // Stability: development
  HTTPServerRequestBodySizeName = "http.server.request.body.size"
  HTTPServerRequestBodySizeUnit = "By"
  HTTPServerRequestBodySizeDescription = "Size of HTTP server request bodies."
  // HTTPServerRequestDuration is the metric conforming to the
  // "http.server.request.duration" semantic conventions. It represents the
  // duration of HTTP server requests.
  // Instrument: histogram
  // Unit: s
  // Stability: stable
  HTTPServerRequestDurationName = "http.server.request.duration"
  HTTPServerRequestDurationUnit = "s"
  HTTPServerRequestDurationDescription = "Duration of HTTP server requests."
  // HTTPServerResponseBodySize is the metric conforming to the
  // "http.server.response.body.size" semantic conventions. It represents the
  // size of HTTP server response bodies.
  // Instrument: histogram
  // Unit: By
  // Stability: development
  HTTPServerResponseBodySizeName = "http.server.response.body.size"
  HTTPServerResponseBodySizeUnit = "By"
  HTTPServerResponseBodySizeDescription = "Size of HTTP server response bodies."
  // JVMBufferCount is the metric conforming to the "jvm.buffer.count" semantic
  // conventions. It represents the number of buffers in the pool.
  // Instrument: updowncounter
  // Unit: {buffer}
  // Stability: development
  JVMBufferCountName = "jvm.buffer.count"
  JVMBufferCountUnit = "{buffer}"
  JVMBufferCountDescription = "Number of buffers in the pool."
  // JVMBufferMemoryLimit is the metric conforming to the
  // "jvm.buffer.memory.limit" semantic conventions. It represents the measure of
  // total memory capacity of buffers.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  JVMBufferMemoryLimitName = "jvm.buffer.memory.limit"
  JVMBufferMemoryLimitUnit = "By"
  JVMBufferMemoryLimitDescription = "Measure of total memory capacity of buffers."
  // JVMBufferMemoryUsage is the metric conforming to the
  // "jvm.buffer.memory.usage" semantic conventions. It represents the
  // deprecated, use `jvm.buffer.memory.used` instead.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  // Deprecated: Replaced by `jvm.buffer.memory.used`.
  JVMBufferMemoryUsageName = "jvm.buffer.memory.usage"
  JVMBufferMemoryUsageUnit = "By"
  JVMBufferMemoryUsageDescription = "Deprecated, use `jvm.buffer.memory.used` instead."
  // JVMBufferMemoryUsed is the metric conforming to the "jvm.buffer.memory.used"
  // semantic conventions. It represents the measure of memory used by buffers.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  JVMBufferMemoryUsedName = "jvm.buffer.memory.used"
  JVMBufferMemoryUsedUnit = "By"
  JVMBufferMemoryUsedDescription = "Measure of memory used by buffers."
  // JVMClassCount is the metric conforming to the "jvm.class.count" semantic
  // conventions. It represents the number of classes currently loaded.
  // Instrument: updowncounter
  // Unit: {class}
  // Stability: stable
  JVMClassCountName = "jvm.class.count"
  JVMClassCountUnit = "{class}"
  JVMClassCountDescription = "Number of classes currently loaded."
  // JVMClassLoaded is the metric conforming to the "jvm.class.loaded" semantic
  // conventions. It represents the number of classes loaded since JVM start.
  // Instrument: counter
  // Unit: {class}
  // Stability: stable
  JVMClassLoadedName = "jvm.class.loaded"
  JVMClassLoadedUnit = "{class}"
  JVMClassLoadedDescription = "Number of classes loaded since JVM start."
  // JVMClassUnloaded is the metric conforming to the "jvm.class.unloaded"
  // semantic conventions. It represents the number of classes unloaded since JVM
  // start.
  // Instrument: counter
  // Unit: {class}
  // Stability: stable
  JVMClassUnloadedName = "jvm.class.unloaded"
  JVMClassUnloadedUnit = "{class}"
  JVMClassUnloadedDescription = "Number of classes unloaded since JVM start."
  // JVMCPUCount is the metric conforming to the "jvm.cpu.count" semantic
  // conventions. It represents the number of processors available to the Java
  // virtual machine.
  // Instrument: updowncounter
  // Unit: {cpu}
  // Stability: stable
  JVMCPUCountName = "jvm.cpu.count"
  JVMCPUCountUnit = "{cpu}"
  JVMCPUCountDescription = "Number of processors available to the Java virtual machine."
  // JVMCPURecentUtilization is the metric conforming to the
  // "jvm.cpu.recent_utilization" semantic conventions. It represents the recent
  // CPU utilization for the process as reported by the JVM.
  // Instrument: gauge
  // Unit: 1
  // Stability: stable
  JVMCPURecentUtilizationName = "jvm.cpu.recent_utilization"
  JVMCPURecentUtilizationUnit = "1"
  JVMCPURecentUtilizationDescription = "Recent CPU utilization for the process as reported by the JVM."
  // JVMCPUTime is the metric conforming to the "jvm.cpu.time" semantic
  // conventions. It represents the cPU time used by the process as reported by
  // the JVM.
  // Instrument: counter
  // Unit: s
  // Stability: stable
  JVMCPUTimeName = "jvm.cpu.time"
  JVMCPUTimeUnit = "s"
  JVMCPUTimeDescription = "CPU time used by the process as reported by the JVM."
  // JVMGCDuration is the metric conforming to the "jvm.gc.duration" semantic
  // conventions. It represents the duration of JVM garbage collection actions.
  // Instrument: histogram
  // Unit: s
  // Stability: stable
  JVMGCDurationName = "jvm.gc.duration"
  JVMGCDurationUnit = "s"
  JVMGCDurationDescription = "Duration of JVM garbage collection actions."
  // JVMMemoryCommitted is the metric conforming to the "jvm.memory.committed"
  // semantic conventions. It represents the measure of memory committed.
  // Instrument: updowncounter
  // Unit: By
  // Stability: stable
  JVMMemoryCommittedName = "jvm.memory.committed"
  JVMMemoryCommittedUnit = "By"
  JVMMemoryCommittedDescription = "Measure of memory committed."
  // JVMMemoryInit is the metric conforming to the "jvm.memory.init" semantic
  // conventions. It represents the measure of initial memory requested.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  JVMMemoryInitName = "jvm.memory.init"
  JVMMemoryInitUnit = "By"
  JVMMemoryInitDescription = "Measure of initial memory requested."
  // JVMMemoryLimit is the metric conforming to the "jvm.memory.limit" semantic
  // conventions. It represents the measure of max obtainable memory.
  // Instrument: updowncounter
  // Unit: By
  // Stability: stable
  JVMMemoryLimitName = "jvm.memory.limit"
  JVMMemoryLimitUnit = "By"
  JVMMemoryLimitDescription = "Measure of max obtainable memory."
  // JVMMemoryUsed is the metric conforming to the "jvm.memory.used" semantic
  // conventions. It represents the measure of memory used.
  // Instrument: updowncounter
  // Unit: By
  // Stability: stable
  JVMMemoryUsedName = "jvm.memory.used"
  JVMMemoryUsedUnit = "By"
  JVMMemoryUsedDescription = "Measure of memory used."
  // JVMMemoryUsedAfterLastGC is the metric conforming to the
  // "jvm.memory.used_after_last_gc" semantic conventions. It represents the
  // measure of memory used, as measured after the most recent garbage collection
  // event on this pool.
  // Instrument: updowncounter
  // Unit: By
  // Stability: stable
  JVMMemoryUsedAfterLastGCName = "jvm.memory.used_after_last_gc"
  JVMMemoryUsedAfterLastGCUnit = "By"
  JVMMemoryUsedAfterLastGCDescription = "Measure of memory used, as measured after the most recent garbage collection event on this pool."
  // JVMSystemCPULoad1m is the metric conforming to the "jvm.system.cpu.load_1m"
  // semantic conventions. It represents the average CPU load of the whole system
  // for the last minute as reported by the JVM.
  // Instrument: gauge
  // Unit: {run_queue_item}
  // Stability: development
  JVMSystemCPULoad1mName = "jvm.system.cpu.load_1m"
  JVMSystemCPULoad1mUnit = "{run_queue_item}"
  JVMSystemCPULoad1mDescription = "Average CPU load of the whole system for the last minute as reported by the JVM."
  // JVMSystemCPUUtilization is the metric conforming to the
  // "jvm.system.cpu.utilization" semantic conventions. It represents the recent
  // CPU utilization for the whole system as reported by the JVM.
  // Instrument: gauge
  // Unit: 1
  // Stability: development
  JVMSystemCPUUtilizationName = "jvm.system.cpu.utilization"
  JVMSystemCPUUtilizationUnit = "1"
  JVMSystemCPUUtilizationDescription = "Recent CPU utilization for the whole system as reported by the JVM."
  // JVMThreadCount is the metric conforming to the "jvm.thread.count" semantic
  // conventions. It represents the number of executing platform threads.
  // Instrument: updowncounter
  // Unit: {thread}
  // Stability: stable
  JVMThreadCountName = "jvm.thread.count"
  JVMThreadCountUnit = "{thread}"
  JVMThreadCountDescription = "Number of executing platform threads."
  // KestrelActiveConnections is the metric conforming to the
  // "kestrel.active_connections" semantic conventions. It represents the number
  // of connections that are currently active on the server.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: stable
  KestrelActiveConnectionsName = "kestrel.active_connections"
  KestrelActiveConnectionsUnit = "{connection}"
  KestrelActiveConnectionsDescription = "Number of connections that are currently active on the server."
  // KestrelActiveTLSHandshakes is the metric conforming to the
  // "kestrel.active_tls_handshakes" semantic conventions. It represents the
  // number of TLS handshakes that are currently in progress on the server.
  // Instrument: updowncounter
  // Unit: {handshake}
  // Stability: stable
  KestrelActiveTLSHandshakesName = "kestrel.active_tls_handshakes"
  KestrelActiveTLSHandshakesUnit = "{handshake}"
  KestrelActiveTLSHandshakesDescription = "Number of TLS handshakes that are currently in progress on the server."
  // KestrelConnectionDuration is the metric conforming to the
  // "kestrel.connection.duration" semantic conventions. It represents the
  // duration of connections on the server.
  // Instrument: histogram
  // Unit: s
  // Stability: stable
  KestrelConnectionDurationName = "kestrel.connection.duration"
  KestrelConnectionDurationUnit = "s"
  KestrelConnectionDurationDescription = "The duration of connections on the server."
  // KestrelQueuedConnections is the metric conforming to the
  // "kestrel.queued_connections" semantic conventions. It represents the number
  // of connections that are currently queued and are waiting to start.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: stable
  KestrelQueuedConnectionsName = "kestrel.queued_connections"
  KestrelQueuedConnectionsUnit = "{connection}"
  KestrelQueuedConnectionsDescription = "Number of connections that are currently queued and are waiting to start."
  // KestrelQueuedRequests is the metric conforming to the
  // "kestrel.queued_requests" semantic conventions. It represents the number of
  // HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are
  // currently queued and are waiting to start.
  // Instrument: updowncounter
  // Unit: {request}
  // Stability: stable
  KestrelQueuedRequestsName = "kestrel.queued_requests"
  KestrelQueuedRequestsUnit = "{request}"
  KestrelQueuedRequestsDescription = "Number of HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are currently queued and are waiting to start."
  // KestrelRejectedConnections is the metric conforming to the
  // "kestrel.rejected_connections" semantic conventions. It represents the
  // number of connections rejected by the server.
  // Instrument: counter
  // Unit: {connection}
  // Stability: stable
  KestrelRejectedConnectionsName = "kestrel.rejected_connections"
  KestrelRejectedConnectionsUnit = "{connection}"
  KestrelRejectedConnectionsDescription = "Number of connections rejected by the server."
  // KestrelTLSHandshakeDuration is the metric conforming to the
  // "kestrel.tls_handshake.duration" semantic conventions. It represents the
  // duration of TLS handshakes on the server.
  // Instrument: histogram
  // Unit: s
  // Stability: stable
  KestrelTLSHandshakeDurationName = "kestrel.tls_handshake.duration"
  KestrelTLSHandshakeDurationUnit = "s"
  KestrelTLSHandshakeDurationDescription = "The duration of TLS handshakes on the server."
  // KestrelUpgradedConnections is the metric conforming to the
  // "kestrel.upgraded_connections" semantic conventions. It represents the
  // number of connections that are currently upgraded (WebSockets). .
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: stable
  KestrelUpgradedConnectionsName = "kestrel.upgraded_connections"
  KestrelUpgradedConnectionsUnit = "{connection}"
  KestrelUpgradedConnectionsDescription = "Number of connections that are currently upgraded (WebSockets). ."
  // MessagingClientConsumedMessages is the metric conforming to the
  // "messaging.client.consumed.messages" semantic conventions. It represents the
  // number of messages that were delivered to the application.
  // Instrument: counter
  // Unit: {message}
  // Stability: development
  MessagingClientConsumedMessagesName = "messaging.client.consumed.messages"
  MessagingClientConsumedMessagesUnit = "{message}"
  MessagingClientConsumedMessagesDescription = "Number of messages that were delivered to the application."
  // MessagingClientOperationDuration is the metric conforming to the
  // "messaging.client.operation.duration" semantic conventions. It represents
  // the duration of messaging operation initiated by a producer or consumer
  // client.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  MessagingClientOperationDurationName = "messaging.client.operation.duration"
  MessagingClientOperationDurationUnit = "s"
  MessagingClientOperationDurationDescription = "Duration of messaging operation initiated by a producer or consumer client."
  // MessagingClientPublishedMessages is the metric conforming to the
  // "messaging.client.published.messages" semantic conventions. It represents
  // the number of messages producer attempted to publish to the broker.
  // Instrument: counter
  // Unit: {message}
  // Stability: development
  MessagingClientPublishedMessagesName = "messaging.client.published.messages"
  MessagingClientPublishedMessagesUnit = "{message}"
  MessagingClientPublishedMessagesDescription = "Number of messages producer attempted to publish to the broker."
  // MessagingProcessDuration is the metric conforming to the
  // "messaging.process.duration" semantic conventions. It represents the
  // duration of processing operation.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  MessagingProcessDurationName = "messaging.process.duration"
  MessagingProcessDurationUnit = "s"
  MessagingProcessDurationDescription = "Duration of processing operation."
  // MessagingProcessMessages is the metric conforming to the
  // "messaging.process.messages" semantic conventions. It represents the
  // deprecated. Use `messaging.client.consumed.messages` instead.
  // Instrument: counter
  // Unit: {message}
  // Stability: development
  // Deprecated: Replaced by `messaging.client.consumed.messages`.
  MessagingProcessMessagesName = "messaging.process.messages"
  MessagingProcessMessagesUnit = "{message}"
  MessagingProcessMessagesDescription = "Deprecated. Use `messaging.client.consumed.messages` instead."
  // MessagingPublishDuration is the metric conforming to the
  // "messaging.publish.duration" semantic conventions. It represents the
  // deprecated. Use `messaging.client.operation.duration` instead.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  // Deprecated: Replaced by `messaging.client.operation.duration`.
  MessagingPublishDurationName = "messaging.publish.duration"
  MessagingPublishDurationUnit = "s"
  MessagingPublishDurationDescription = "Deprecated. Use `messaging.client.operation.duration` instead."
  // MessagingPublishMessages is the metric conforming to the
  // "messaging.publish.messages" semantic conventions. It represents the
  // deprecated. Use `messaging.client.produced.messages` instead.
  // Instrument: counter
  // Unit: {message}
  // Stability: development
  // Deprecated: Replaced by `messaging.client.produced.messages`.
  MessagingPublishMessagesName = "messaging.publish.messages"
  MessagingPublishMessagesUnit = "{message}"
  MessagingPublishMessagesDescription = "Deprecated. Use `messaging.client.produced.messages` instead."
  // MessagingReceiveDuration is the metric conforming to the
  // "messaging.receive.duration" semantic conventions. It represents the
  // deprecated. Use `messaging.client.operation.duration` instead.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  // Deprecated: Replaced by `messaging.client.operation.duration`.
  MessagingReceiveDurationName = "messaging.receive.duration"
  MessagingReceiveDurationUnit = "s"
  MessagingReceiveDurationDescription = "Deprecated. Use `messaging.client.operation.duration` instead."
  // MessagingReceiveMessages is the metric conforming to the
  // "messaging.receive.messages" semantic conventions. It represents the
  // deprecated. Use `messaging.client.consumed.messages` instead.
  // Instrument: counter
  // Unit: {message}
  // Stability: development
  // Deprecated: Replaced by `messaging.client.consumed.messages`.
  MessagingReceiveMessagesName = "messaging.receive.messages"
  MessagingReceiveMessagesUnit = "{message}"
  MessagingReceiveMessagesDescription = "Deprecated. Use `messaging.client.consumed.messages` instead."
  // NodejsEventloopDelayMax is the metric conforming to the
  // "nodejs.eventloop.delay.max" semantic conventions. It represents the event
  // loop maximum delay.
  // Instrument: gauge
  // Unit: s
  // Stability: development
  NodejsEventloopDelayMaxName = "nodejs.eventloop.delay.max"
  NodejsEventloopDelayMaxUnit = "s"
  NodejsEventloopDelayMaxDescription = "Event loop maximum delay."
  // NodejsEventloopDelayMean is the metric conforming to the
  // "nodejs.eventloop.delay.mean" semantic conventions. It represents the event
  // loop mean delay.
  // Instrument: gauge
  // Unit: s
  // Stability: development
  NodejsEventloopDelayMeanName = "nodejs.eventloop.delay.mean"
  NodejsEventloopDelayMeanUnit = "s"
  NodejsEventloopDelayMeanDescription = "Event loop mean delay."
  // NodejsEventloopDelayMin is the metric conforming to the
  // "nodejs.eventloop.delay.min" semantic conventions. It represents the event
  // loop minimum delay.
  // Instrument: gauge
  // Unit: s
  // Stability: development
  NodejsEventloopDelayMinName = "nodejs.eventloop.delay.min"
  NodejsEventloopDelayMinUnit = "s"
  NodejsEventloopDelayMinDescription = "Event loop minimum delay."
  // NodejsEventloopDelayP50 is the metric conforming to the
  // "nodejs.eventloop.delay.p50" semantic conventions. It represents the event
  // loop 50 percentile delay.
  // Instrument: gauge
  // Unit: s
  // Stability: development
  NodejsEventloopDelayP50Name = "nodejs.eventloop.delay.p50"
  NodejsEventloopDelayP50Unit = "s"
  NodejsEventloopDelayP50Description = "Event loop 50 percentile delay."
  // NodejsEventloopDelayP90 is the metric conforming to the
  // "nodejs.eventloop.delay.p90" semantic conventions. It represents the event
  // loop 90 percentile delay.
  // Instrument: gauge
  // Unit: s
  // Stability: development
  NodejsEventloopDelayP90Name = "nodejs.eventloop.delay.p90"
  NodejsEventloopDelayP90Unit = "s"
  NodejsEventloopDelayP90Description = "Event loop 90 percentile delay."
  // NodejsEventloopDelayP99 is the metric conforming to the
  // "nodejs.eventloop.delay.p99" semantic conventions. It represents the event
  // loop 99 percentile delay.
  // Instrument: gauge
  // Unit: s
  // Stability: development
  NodejsEventloopDelayP99Name = "nodejs.eventloop.delay.p99"
  NodejsEventloopDelayP99Unit = "s"
  NodejsEventloopDelayP99Description = "Event loop 99 percentile delay."
  // NodejsEventloopDelayStddev is the metric conforming to the
  // "nodejs.eventloop.delay.stddev" semantic conventions. It represents the
  // event loop standard deviation delay.
  // Instrument: gauge
  // Unit: s
  // Stability: development
  NodejsEventloopDelayStddevName = "nodejs.eventloop.delay.stddev"
  NodejsEventloopDelayStddevUnit = "s"
  NodejsEventloopDelayStddevDescription = "Event loop standard deviation delay."
  // NodejsEventloopUtilization is the metric conforming to the
  // "nodejs.eventloop.utilization" semantic conventions. It represents the event
  // loop utilization.
  // Instrument: gauge
  // Unit: 1
  // Stability: development
  NodejsEventloopUtilizationName = "nodejs.eventloop.utilization"
  NodejsEventloopUtilizationUnit = "1"
  NodejsEventloopUtilizationDescription = "Event loop utilization."
  // ProcessContextSwitches is the metric conforming to the
  // "process.context_switches" semantic conventions. It represents the number of
  // times the process has been context switched.
  // Instrument: counter
  // Unit: {count}
  // Stability: development
  ProcessContextSwitchesName = "process.context_switches"
  ProcessContextSwitchesUnit = "{count}"
  ProcessContextSwitchesDescription = "Number of times the process has been context switched."
  // ProcessCPUTime is the metric conforming to the "process.cpu.time" semantic
  // conventions. It represents the total CPU seconds broken down by different
  // states.
  // Instrument: counter
  // Unit: s
  // Stability: development
  ProcessCPUTimeName = "process.cpu.time"
  ProcessCPUTimeUnit = "s"
  ProcessCPUTimeDescription = "Total CPU seconds broken down by different states."
  // ProcessCPUUtilization is the metric conforming to the
  // "process.cpu.utilization" semantic conventions. It represents the difference
  // in process.cpu.time since the last measurement, divided by the elapsed time
  // and number of CPUs available to the process.
  // Instrument: gauge
  // Unit: 1
  // Stability: development
  ProcessCPUUtilizationName = "process.cpu.utilization"
  ProcessCPUUtilizationUnit = "1"
  ProcessCPUUtilizationDescription = "Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process."
  // ProcessDiskIo is the metric conforming to the "process.disk.io" semantic
  // conventions. It represents the disk bytes transferred.
  // Instrument: counter
  // Unit: By
  // Stability: development
  ProcessDiskIoName = "process.disk.io"
  ProcessDiskIoUnit = "By"
  ProcessDiskIoDescription = "Disk bytes transferred."
  // ProcessMemoryUsage is the metric conforming to the "process.memory.usage"
  // semantic conventions. It represents the amount of physical memory in use.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  ProcessMemoryUsageName = "process.memory.usage"
  ProcessMemoryUsageUnit = "By"
  ProcessMemoryUsageDescription = "The amount of physical memory in use."
  // ProcessMemoryVirtual is the metric conforming to the
  // "process.memory.virtual" semantic conventions. It represents the amount of
  // committed virtual memory.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  ProcessMemoryVirtualName = "process.memory.virtual"
  ProcessMemoryVirtualUnit = "By"
  ProcessMemoryVirtualDescription = "The amount of committed virtual memory."
  // ProcessNetworkIo is the metric conforming to the "process.network.io"
  // semantic conventions. It represents the network bytes transferred.
  // Instrument: counter
  // Unit: By
  // Stability: development
  ProcessNetworkIoName = "process.network.io"
  ProcessNetworkIoUnit = "By"
  ProcessNetworkIoDescription = "Network bytes transferred."
  // ProcessOpenFileDescriptorCount is the metric conforming to the
  // "process.open_file_descriptor.count" semantic conventions. It represents the
  // number of file descriptors in use by the process.
  // Instrument: updowncounter
  // Unit: {count}
  // Stability: development
  ProcessOpenFileDescriptorCountName = "process.open_file_descriptor.count"
  ProcessOpenFileDescriptorCountUnit = "{count}"
  ProcessOpenFileDescriptorCountDescription = "Number of file descriptors in use by the process."
  // ProcessPagingFaults is the metric conforming to the "process.paging.faults"
  // semantic conventions. It represents the number of page faults the process
  // has made.
  // Instrument: counter
  // Unit: {fault}
  // Stability: development
  ProcessPagingFaultsName = "process.paging.faults"
  ProcessPagingFaultsUnit = "{fault}"
  ProcessPagingFaultsDescription = "Number of page faults the process has made."
  // ProcessThreadCount is the metric conforming to the "process.thread.count"
  // semantic conventions. It represents the process threads count.
  // Instrument: updowncounter
  // Unit: {thread}
  // Stability: development
  ProcessThreadCountName = "process.thread.count"
  ProcessThreadCountUnit = "{thread}"
  ProcessThreadCountDescription = "Process threads count."
  // RPCClientDuration is the metric conforming to the "rpc.client.duration"
  // semantic conventions. It represents the measures the duration of outbound
  // RPC.
  // Instrument: histogram
  // Unit: ms
  // Stability: development
  RPCClientDurationName = "rpc.client.duration"
  RPCClientDurationUnit = "ms"
  RPCClientDurationDescription = "Measures the duration of outbound RPC."
  // RPCClientRequestSize is the metric conforming to the
  // "rpc.client.request.size" semantic conventions. It represents the measures
  // the size of RPC request messages (uncompressed).
  // Instrument: histogram
  // Unit: By
  // Stability: development
  RPCClientRequestSizeName = "rpc.client.request.size"
  RPCClientRequestSizeUnit = "By"
  RPCClientRequestSizeDescription = "Measures the size of RPC request messages (uncompressed)."
  // RPCClientRequestsPerRPC is the metric conforming to the
  // "rpc.client.requests_per_rpc" semantic conventions. It represents the
  // measures the number of messages received per RPC.
  // Instrument: histogram
  // Unit: {count}
  // Stability: development
  RPCClientRequestsPerRPCName = "rpc.client.requests_per_rpc"
  RPCClientRequestsPerRPCUnit = "{count}"
  RPCClientRequestsPerRPCDescription = "Measures the number of messages received per RPC."
  // RPCClientResponseSize is the metric conforming to the
  // "rpc.client.response.size" semantic conventions. It represents the measures
  // the size of RPC response messages (uncompressed).
  // Instrument: histogram
  // Unit: By
  // Stability: development
  RPCClientResponseSizeName = "rpc.client.response.size"
  RPCClientResponseSizeUnit = "By"
  RPCClientResponseSizeDescription = "Measures the size of RPC response messages (uncompressed)."
  // RPCClientResponsesPerRPC is the metric conforming to the
  // "rpc.client.responses_per_rpc" semantic conventions. It represents the
  // measures the number of messages sent per RPC.
  // Instrument: histogram
  // Unit: {count}
  // Stability: development
  RPCClientResponsesPerRPCName = "rpc.client.responses_per_rpc"
  RPCClientResponsesPerRPCUnit = "{count}"
  RPCClientResponsesPerRPCDescription = "Measures the number of messages sent per RPC."
  // RPCServerDuration is the metric conforming to the "rpc.server.duration"
  // semantic conventions. It represents the measures the duration of inbound
  // RPC.
  // Instrument: histogram
  // Unit: ms
  // Stability: development
  RPCServerDurationName = "rpc.server.duration"
  RPCServerDurationUnit = "ms"
  RPCServerDurationDescription = "Measures the duration of inbound RPC."
  // RPCServerRequestSize is the metric conforming to the
  // "rpc.server.request.size" semantic conventions. It represents the measures
  // the size of RPC request messages (uncompressed).
  // Instrument: histogram
  // Unit: By
  // Stability: development
  RPCServerRequestSizeName = "rpc.server.request.size"
  RPCServerRequestSizeUnit = "By"
  RPCServerRequestSizeDescription = "Measures the size of RPC request messages (uncompressed)."
  // RPCServerRequestsPerRPC is the metric conforming to the
  // "rpc.server.requests_per_rpc" semantic conventions. It represents the
  // measures the number of messages received per RPC.
  // Instrument: histogram
  // Unit: {count}
  // Stability: development
  RPCServerRequestsPerRPCName = "rpc.server.requests_per_rpc"
  RPCServerRequestsPerRPCUnit = "{count}"
  RPCServerRequestsPerRPCDescription = "Measures the number of messages received per RPC."
  // RPCServerResponseSize is the metric conforming to the
  // "rpc.server.response.size" semantic conventions. It represents the measures
  // the size of RPC response messages (uncompressed).
  // Instrument: histogram
  // Unit: By
  // Stability: development
  RPCServerResponseSizeName = "rpc.server.response.size"
  RPCServerResponseSizeUnit = "By"
  RPCServerResponseSizeDescription = "Measures the size of RPC response messages (uncompressed)."
  // RPCServerResponsesPerRPC is the metric conforming to the
  // "rpc.server.responses_per_rpc" semantic conventions. It represents the
  // measures the number of messages sent per RPC.
  // Instrument: histogram
  // Unit: {count}
  // Stability: development
  RPCServerResponsesPerRPCName = "rpc.server.responses_per_rpc"
  RPCServerResponsesPerRPCUnit = "{count}"
  RPCServerResponsesPerRPCDescription = "Measures the number of messages sent per RPC."
  // SignalrServerActiveConnections is the metric conforming to the
  // "signalr.server.active_connections" semantic conventions. It represents the
  // number of connections that are currently active on the server.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: stable
  SignalrServerActiveConnectionsName = "signalr.server.active_connections"
  SignalrServerActiveConnectionsUnit = "{connection}"
  SignalrServerActiveConnectionsDescription = "Number of connections that are currently active on the server."
  // SignalrServerConnectionDuration is the metric conforming to the
  // "signalr.server.connection.duration" semantic conventions. It represents the
  // duration of connections on the server.
  // Instrument: histogram
  // Unit: s
  // Stability: stable
  SignalrServerConnectionDurationName = "signalr.server.connection.duration"
  SignalrServerConnectionDurationUnit = "s"
  SignalrServerConnectionDurationDescription = "The duration of connections on the server."
  // SystemCPUFrequency is the metric conforming to the "system.cpu.frequency"
  // semantic conventions. It represents the reports the current frequency of the
  // CPU in Hz.
  // Instrument: gauge
  // Unit: {Hz}
  // Stability: development
  SystemCPUFrequencyName = "system.cpu.frequency"
  SystemCPUFrequencyUnit = "{Hz}"
  SystemCPUFrequencyDescription = "Reports the current frequency of the CPU in Hz"
  // SystemCPULogicalCount is the metric conforming to the
  // "system.cpu.logical.count" semantic conventions. It represents the reports
  // the number of logical (virtual) processor cores created by the operating
  // system to manage multitasking.
  // Instrument: updowncounter
  // Unit: {cpu}
  // Stability: development
  SystemCPULogicalCountName = "system.cpu.logical.count"
  SystemCPULogicalCountUnit = "{cpu}"
  SystemCPULogicalCountDescription = "Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking"
  // SystemCPUPhysicalCount is the metric conforming to the
  // "system.cpu.physical.count" semantic conventions. It represents the reports
  // the number of actual physical processor cores on the hardware.
  // Instrument: updowncounter
  // Unit: {cpu}
  // Stability: development
  SystemCPUPhysicalCountName = "system.cpu.physical.count"
  SystemCPUPhysicalCountUnit = "{cpu}"
  SystemCPUPhysicalCountDescription = "Reports the number of actual physical processor cores on the hardware"
  // SystemCPUTime is the metric conforming to the "system.cpu.time" semantic
  // conventions. It represents the seconds each logical CPU spent on each mode.
  // Instrument: counter
  // Unit: s
  // Stability: development
  SystemCPUTimeName = "system.cpu.time"
  SystemCPUTimeUnit = "s"
  SystemCPUTimeDescription = "Seconds each logical CPU spent on each mode"
  // SystemCPUUtilization is the metric conforming to the
  // "system.cpu.utilization" semantic conventions. It represents the difference
  // in system.cpu.time since the last measurement, divided by the elapsed time
  // and number of logical CPUs.
  // Instrument: gauge
  // Unit: 1
  // Stability: development
  SystemCPUUtilizationName = "system.cpu.utilization"
  SystemCPUUtilizationUnit = "1"
  SystemCPUUtilizationDescription = "Difference in system.cpu.time since the last measurement, divided by the elapsed time and number of logical CPUs"
  // SystemDiskIo is the metric conforming to the "system.disk.io" semantic
  // conventions.
  // Instrument: counter
  // Unit: By
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemDiskIoName = "system.disk.io"
  SystemDiskIoUnit = "By"
  // SystemDiskIoTime is the metric conforming to the "system.disk.io_time"
  // semantic conventions. It represents the time disk spent activated.
  // Instrument: counter
  // Unit: s
  // Stability: development
  SystemDiskIoTimeName = "system.disk.io_time"
  SystemDiskIoTimeUnit = "s"
  SystemDiskIoTimeDescription = "Time disk spent activated"
  // SystemDiskMerged is the metric conforming to the "system.disk.merged"
  // semantic conventions.
  // Instrument: counter
  // Unit: {operation}
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemDiskMergedName = "system.disk.merged"
  SystemDiskMergedUnit = "{operation}"
  // SystemDiskOperationTime is the metric conforming to the
  // "system.disk.operation_time" semantic conventions. It represents the sum of
  // the time each operation took to complete.
  // Instrument: counter
  // Unit: s
  // Stability: development
  SystemDiskOperationTimeName = "system.disk.operation_time"
  SystemDiskOperationTimeUnit = "s"
  SystemDiskOperationTimeDescription = "Sum of the time each operation took to complete"
  // SystemDiskOperations is the metric conforming to the
  // "system.disk.operations" semantic conventions.
  // Instrument: counter
  // Unit: {operation}
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemDiskOperationsName = "system.disk.operations"
  SystemDiskOperationsUnit = "{operation}"
  // SystemFilesystemUsage is the metric conforming to the
  // "system.filesystem.usage" semantic conventions.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemFilesystemUsageName = "system.filesystem.usage"
  SystemFilesystemUsageUnit = "By"
  // SystemFilesystemUtilization is the metric conforming to the
  // "system.filesystem.utilization" semantic conventions.
  // Instrument: gauge
  // Unit: 1
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemFilesystemUtilizationName = "system.filesystem.utilization"
  SystemFilesystemUtilizationUnit = "1"
  // SystemLinuxMemoryAvailable is the metric conforming to the
  // "system.linux.memory.available" semantic conventions. It represents an
  // estimate of how much memory is available for starting new applications,
  // without causing swapping.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  SystemLinuxMemoryAvailableName = "system.linux.memory.available"
  SystemLinuxMemoryAvailableUnit = "By"
  SystemLinuxMemoryAvailableDescription = "An estimate of how much memory is available for starting new applications, without causing swapping"
  // SystemLinuxMemorySlabUsage is the metric conforming to the
  // "system.linux.memory.slab.usage" semantic conventions. It represents the
  // reports the memory used by the Linux kernel for managing caches of
  // frequently used objects.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  SystemLinuxMemorySlabUsageName = "system.linux.memory.slab.usage"
  SystemLinuxMemorySlabUsageUnit = "By"
  SystemLinuxMemorySlabUsageDescription = "Reports the memory used by the Linux kernel for managing caches of frequently used objects."
  // SystemMemoryLimit is the metric conforming to the "system.memory.limit"
  // semantic conventions. It represents the total memory available in the
  // system.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  SystemMemoryLimitName = "system.memory.limit"
  SystemMemoryLimitUnit = "By"
  SystemMemoryLimitDescription = "Total memory available in the system."
  // SystemMemoryShared is the metric conforming to the "system.memory.shared"
  // semantic conventions. It represents the shared memory used (mostly by
  // tmpfs).
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  SystemMemorySharedName = "system.memory.shared"
  SystemMemorySharedUnit = "By"
  SystemMemorySharedDescription = "Shared memory used (mostly by tmpfs)."
  // SystemMemoryUsage is the metric conforming to the "system.memory.usage"
  // semantic conventions. It represents the reports memory in use by state.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  SystemMemoryUsageName = "system.memory.usage"
  SystemMemoryUsageUnit = "By"
  SystemMemoryUsageDescription = "Reports memory in use by state."
  // SystemMemoryUtilization is the metric conforming to the
  // "system.memory.utilization" semantic conventions.
  // Instrument: gauge
  // Unit: 1
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemMemoryUtilizationName = "system.memory.utilization"
  SystemMemoryUtilizationUnit = "1"
  // SystemNetworkConnections is the metric conforming to the
  // "system.network.connections" semantic conventions.
  // Instrument: updowncounter
  // Unit: {connection}
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemNetworkConnectionsName = "system.network.connections"
  SystemNetworkConnectionsUnit = "{connection}"
  // SystemNetworkDropped is the metric conforming to the
  // "system.network.dropped" semantic conventions. It represents the count of
  // packets that are dropped or discarded even though there was no error.
  // Instrument: counter
  // Unit: {packet}
  // Stability: development
  SystemNetworkDroppedName = "system.network.dropped"
  SystemNetworkDroppedUnit = "{packet}"
  SystemNetworkDroppedDescription = "Count of packets that are dropped or discarded even though there was no error"
  // SystemNetworkErrors is the metric conforming to the "system.network.errors"
  // semantic conventions. It represents the count of network errors detected.
  // Instrument: counter
  // Unit: {error}
  // Stability: development
  SystemNetworkErrorsName = "system.network.errors"
  SystemNetworkErrorsUnit = "{error}"
  SystemNetworkErrorsDescription = "Count of network errors detected"
  // SystemNetworkIo is the metric conforming to the "system.network.io" semantic
  // conventions.
  // Instrument: counter
  // Unit: By
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemNetworkIoName = "system.network.io"
  SystemNetworkIoUnit = "By"
  // SystemNetworkPackets is the metric conforming to the
  // "system.network.packets" semantic conventions.
  // Instrument: counter
  // Unit: {packet}
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemNetworkPacketsName = "system.network.packets"
  SystemNetworkPacketsUnit = "{packet}"
  // SystemPagingFaults is the metric conforming to the "system.paging.faults"
  // semantic conventions.
  // Instrument: counter
  // Unit: {fault}
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemPagingFaultsName = "system.paging.faults"
  SystemPagingFaultsUnit = "{fault}"
  // SystemPagingOperations is the metric conforming to the
  // "system.paging.operations" semantic conventions.
  // Instrument: counter
  // Unit: {operation}
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemPagingOperationsName = "system.paging.operations"
  SystemPagingOperationsUnit = "{operation}"
  // SystemPagingUsage is the metric conforming to the "system.paging.usage"
  // semantic conventions. It represents the unix swap or windows pagefile usage.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  SystemPagingUsageName = "system.paging.usage"
  SystemPagingUsageUnit = "By"
  SystemPagingUsageDescription = "Unix swap or windows pagefile usage"
  // SystemPagingUtilization is the metric conforming to the
  // "system.paging.utilization" semantic conventions.
  // Instrument: gauge
  // Unit: 1
  // Stability: development
  // NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
  SystemPagingUtilizationName = "system.paging.utilization"
  SystemPagingUtilizationUnit = "1"
  // SystemProcessCount is the metric conforming to the "system.process.count"
  // semantic conventions. It represents the total number of processes in each
  // state.
  // Instrument: updowncounter
  // Unit: {process}
  // Stability: development
  SystemProcessCountName = "system.process.count"
  SystemProcessCountUnit = "{process}"
  SystemProcessCountDescription = "Total number of processes in each state"
  // SystemProcessCreated is the metric conforming to the
  // "system.process.created" semantic conventions. It represents the total
  // number of processes created over uptime of the host.
  // Instrument: counter
  // Unit: {process}
  // Stability: development
  SystemProcessCreatedName = "system.process.created"
  SystemProcessCreatedUnit = "{process}"
  SystemProcessCreatedDescription = "Total number of processes created over uptime of the host"
  // V8JSGCDuration is the metric conforming to the "v8js.gc.duration" semantic
  // conventions. It represents the garbage collection duration.
  // Instrument: histogram
  // Unit: s
  // Stability: development
  V8JSGCDurationName = "v8js.gc.duration"
  V8JSGCDurationUnit = "s"
  V8JSGCDurationDescription = "Garbage collection duration."
  // V8JSHeapSpaceAvailableSize is the metric conforming to the
  // "v8js.heap.space.available_size" semantic conventions. It represents the
  // heap space available size.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  V8JSHeapSpaceAvailableSizeName = "v8js.heap.space.available_size"
  V8JSHeapSpaceAvailableSizeUnit = "By"
  V8JSHeapSpaceAvailableSizeDescription = "Heap space available size."
  // V8JSHeapSpacePhysicalSize is the metric conforming to the
  // "v8js.heap.space.physical_size" semantic conventions. It represents the
  // committed size of a heap space.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  V8JSHeapSpacePhysicalSizeName = "v8js.heap.space.physical_size"
  V8JSHeapSpacePhysicalSizeUnit = "By"
  V8JSHeapSpacePhysicalSizeDescription = "Committed size of a heap space."
  // V8JSMemoryHeapLimit is the metric conforming to the "v8js.memory.heap.limit"
  // semantic conventions. It represents the total heap memory size
  // pre-allocated.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  V8JSMemoryHeapLimitName = "v8js.memory.heap.limit"
  V8JSMemoryHeapLimitUnit = "By"
  V8JSMemoryHeapLimitDescription = "Total heap memory size pre-allocated."
  // V8JSMemoryHeapUsed is the metric conforming to the "v8js.memory.heap.used"
  // semantic conventions. It represents the heap Memory size allocated.
  // Instrument: updowncounter
  // Unit: By
  // Stability: development
  V8JSMemoryHeapUsedName = "v8js.memory.heap.used"
  V8JSMemoryHeapUsedUnit = "By"
  V8JSMemoryHeapUsedDescription = "Heap Memory size allocated."
)