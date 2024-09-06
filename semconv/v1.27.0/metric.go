// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.27.0"

const (

	// AspnetcoreRoutingMatchAttempts is the metric conforming to the "aspnetcore.routing.match_attempts" semantic conventions. It represents the number of requests that were attempted to be matched to an endpoint
	// Instrument: counter
	// Unit: {match_attempt}
	// Stability: stable
	AspnetcoreRoutingMatchAttemptsName        = "aspnetcore.routing.match_attempts"
	AspnetcoreRoutingMatchAttemptsUnit        = "{match_attempt}"
	AspnetcoreRoutingMatchAttemptsDescription = "Number of requests that were attempted to be matched to an endpoint."

	// AspnetcoreDiagnosticsExceptions is the metric conforming to the "aspnetcore.diagnostics.exceptions" semantic conventions. It represents the number of exceptions caught by exception handling middleware
	// Instrument: counter
	// Unit: {exception}
	// Stability: stable
	AspnetcoreDiagnosticsExceptionsName        = "aspnetcore.diagnostics.exceptions"
	AspnetcoreDiagnosticsExceptionsUnit        = "{exception}"
	AspnetcoreDiagnosticsExceptionsDescription = "Number of exceptions caught by exception handling middleware."

	// AspnetcoreRateLimitingActiveRequestLeases is the metric conforming to the "aspnetcore.rate_limiting.active_request_leases" semantic conventions. It represents the number of requests that are currently active on the server that hold a rate limiting lease
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: stable
	AspnetcoreRateLimitingActiveRequestLeasesName        = "aspnetcore.rate_limiting.active_request_leases"
	AspnetcoreRateLimitingActiveRequestLeasesUnit        = "{request}"
	AspnetcoreRateLimitingActiveRequestLeasesDescription = "Number of requests that are currently active on the server that hold a rate limiting lease."

	// AspnetcoreRateLimitingRequestLeaseDuration is the metric conforming to the "aspnetcore.rate_limiting.request_lease.duration" semantic conventions. It represents the duration of rate limiting lease held by requests on the server
	// Instrument: histogram
	// Unit: s
	// Stability: stable
	AspnetcoreRateLimitingRequestLeaseDurationName        = "aspnetcore.rate_limiting.request_lease.duration"
	AspnetcoreRateLimitingRequestLeaseDurationUnit        = "s"
	AspnetcoreRateLimitingRequestLeaseDurationDescription = "The duration of rate limiting lease held by requests on the server."

	// AspnetcoreRateLimitingRequestTimeInQueue is the metric conforming to the "aspnetcore.rate_limiting.request.time_in_queue" semantic conventions. It represents the time the request spent in a queue waiting to acquire a rate limiting lease
	// Instrument: histogram
	// Unit: s
	// Stability: stable
	AspnetcoreRateLimitingRequestTimeInQueueName        = "aspnetcore.rate_limiting.request.time_in_queue"
	AspnetcoreRateLimitingRequestTimeInQueueUnit        = "s"
	AspnetcoreRateLimitingRequestTimeInQueueDescription = "The time the request spent in a queue waiting to acquire a rate limiting lease."

	// AspnetcoreRateLimitingQueuedRequests is the metric conforming to the "aspnetcore.rate_limiting.queued_requests" semantic conventions. It represents the number of requests that are currently queued, waiting to acquire a rate limiting lease
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: stable
	AspnetcoreRateLimitingQueuedRequestsName        = "aspnetcore.rate_limiting.queued_requests"
	AspnetcoreRateLimitingQueuedRequestsUnit        = "{request}"
	AspnetcoreRateLimitingQueuedRequestsDescription = "Number of requests that are currently queued, waiting to acquire a rate limiting lease."

	// AspnetcoreRateLimitingRequests is the metric conforming to the "aspnetcore.rate_limiting.requests" semantic conventions. It represents the number of requests that tried to acquire a rate limiting lease
	// Instrument: counter
	// Unit: {request}
	// Stability: stable
	AspnetcoreRateLimitingRequestsName        = "aspnetcore.rate_limiting.requests"
	AspnetcoreRateLimitingRequestsUnit        = "{request}"
	AspnetcoreRateLimitingRequestsDescription = "Number of requests that tried to acquire a rate limiting lease."

	// ContainerCpuTime is the metric conforming to the "container.cpu.time" semantic conventions. It represents the total CPU time consumed
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	ContainerCpuTimeName        = "container.cpu.time"
	ContainerCpuTimeUnit        = "s"
	ContainerCpuTimeDescription = "Total CPU time consumed"

	// ContainerCpuUsage is the metric conforming to the "container.cpu.usage" semantic conventions. It represents the container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs
	// Instrument: gauge
	// Unit: {cpu}
	// Stability: experimental
	ContainerCpuUsageName        = "container.cpu.usage"
	ContainerCpuUsageUnit        = "{cpu}"
	ContainerCpuUsageDescription = "Container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"

	// ContainerMemoryUsage is the metric conforming to the "container.memory.usage" semantic conventions. It represents the memory usage of the container
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	ContainerMemoryUsageName        = "container.memory.usage"
	ContainerMemoryUsageUnit        = "By"
	ContainerMemoryUsageDescription = "Memory usage of the container."

	// ContainerDiskIo is the metric conforming to the "container.disk.io" semantic conventions. It represents the disk bytes for the container
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	ContainerDiskIoName        = "container.disk.io"
	ContainerDiskIoUnit        = "By"
	ContainerDiskIoDescription = "Disk bytes for the container."

	// ContainerNetworkIo is the metric conforming to the "container.network.io" semantic conventions. It represents the network bytes for the container
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	ContainerNetworkIoName        = "container.network.io"
	ContainerNetworkIoUnit        = "By"
	ContainerNetworkIoDescription = "Network bytes for the container."

	// DbClientOperationDuration is the metric conforming to the "db.client.operation.duration" semantic conventions. It represents the duration of database client operations
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	DbClientOperationDurationName        = "db.client.operation.duration"
	DbClientOperationDurationUnit        = "s"
	DbClientOperationDurationDescription = "Duration of database client operations."

	// DbClientConnectionCount is the metric conforming to the "db.client.connection.count" semantic conventions. It represents the number of connections that are currently in state described by the `state` attribute
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	DbClientConnectionCountName        = "db.client.connection.count"
	DbClientConnectionCountUnit        = "{connection}"
	DbClientConnectionCountDescription = "The number of connections that are currently in state described by the `state` attribute"

	// DbClientConnectionIdleMax is the metric conforming to the "db.client.connection.idle.max" semantic conventions. It represents the maximum number of idle open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	DbClientConnectionIdleMaxName        = "db.client.connection.idle.max"
	DbClientConnectionIdleMaxUnit        = "{connection}"
	DbClientConnectionIdleMaxDescription = "The maximum number of idle open connections allowed"

	// DbClientConnectionIdleMin is the metric conforming to the "db.client.connection.idle.min" semantic conventions. It represents the minimum number of idle open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	DbClientConnectionIdleMinName        = "db.client.connection.idle.min"
	DbClientConnectionIdleMinUnit        = "{connection}"
	DbClientConnectionIdleMinDescription = "The minimum number of idle open connections allowed"

	// DbClientConnectionMax is the metric conforming to the "db.client.connection.max" semantic conventions. It represents the maximum number of open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	DbClientConnectionMaxName        = "db.client.connection.max"
	DbClientConnectionMaxUnit        = "{connection}"
	DbClientConnectionMaxDescription = "The maximum number of open connections allowed"

	// DbClientConnectionPendingRequests is the metric conforming to the "db.client.connection.pending_requests" semantic conventions. It represents the number of current pending requests for an open connection
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: experimental
	DbClientConnectionPendingRequestsName        = "db.client.connection.pending_requests"
	DbClientConnectionPendingRequestsUnit        = "{request}"
	DbClientConnectionPendingRequestsDescription = "The number of current pending requests for an open connection"

	// DbClientConnectionTimeouts is the metric conforming to the "db.client.connection.timeouts" semantic conventions. It represents the number of connection timeouts that have occurred trying to obtain a connection from the pool
	// Instrument: counter
	// Unit: {timeout}
	// Stability: experimental
	DbClientConnectionTimeoutsName        = "db.client.connection.timeouts"
	DbClientConnectionTimeoutsUnit        = "{timeout}"
	DbClientConnectionTimeoutsDescription = "The number of connection timeouts that have occurred trying to obtain a connection from the pool"

	// DbClientConnectionCreateTime is the metric conforming to the "db.client.connection.create_time" semantic conventions. It represents the time it took to create a new connection
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	DbClientConnectionCreateTimeName        = "db.client.connection.create_time"
	DbClientConnectionCreateTimeUnit        = "s"
	DbClientConnectionCreateTimeDescription = "The time it took to create a new connection"

	// DbClientConnectionWaitTime is the metric conforming to the "db.client.connection.wait_time" semantic conventions. It represents the time it took to obtain an open connection from the pool
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	DbClientConnectionWaitTimeName        = "db.client.connection.wait_time"
	DbClientConnectionWaitTimeUnit        = "s"
	DbClientConnectionWaitTimeDescription = "The time it took to obtain an open connection from the pool"

	// DbClientConnectionUseTime is the metric conforming to the "db.client.connection.use_time" semantic conventions. It represents the time between borrowing a connection and returning it to the pool
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	DbClientConnectionUseTimeName        = "db.client.connection.use_time"
	DbClientConnectionUseTimeUnit        = "s"
	DbClientConnectionUseTimeDescription = "The time between borrowing a connection and returning it to the pool"

	// DbClientConnectionsUsage is the metric conforming to the "db.client.connections.usage" semantic conventions. It represents the deprecated, use `db.client.connection.count` instead
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.count`.
	DbClientConnectionsUsageName        = "db.client.connections.usage"
	DbClientConnectionsUsageUnit        = "{connection}"
	DbClientConnectionsUsageDescription = "Deprecated, use `db.client.connection.count` instead."

	// DbClientConnectionsIdleMax is the metric conforming to the "db.client.connections.idle.max" semantic conventions. It represents the deprecated, use `db.client.connection.idle.max` instead
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.idle.max`.
	DbClientConnectionsIdleMaxName        = "db.client.connections.idle.max"
	DbClientConnectionsIdleMaxUnit        = "{connection}"
	DbClientConnectionsIdleMaxDescription = "Deprecated, use `db.client.connection.idle.max` instead."

	// DbClientConnectionsIdleMin is the metric conforming to the "db.client.connections.idle.min" semantic conventions. It represents the deprecated, use `db.client.connection.idle.min` instead
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.idle.min`.
	DbClientConnectionsIdleMinName        = "db.client.connections.idle.min"
	DbClientConnectionsIdleMinUnit        = "{connection}"
	DbClientConnectionsIdleMinDescription = "Deprecated, use `db.client.connection.idle.min` instead."

	// DbClientConnectionsMax is the metric conforming to the "db.client.connections.max" semantic conventions. It represents the deprecated, use `db.client.connection.max` instead
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.max`.
	DbClientConnectionsMaxName        = "db.client.connections.max"
	DbClientConnectionsMaxUnit        = "{connection}"
	DbClientConnectionsMaxDescription = "Deprecated, use `db.client.connection.max` instead."

	// DbClientConnectionsPendingRequests is the metric conforming to the "db.client.connections.pending_requests" semantic conventions. It represents the deprecated, use `db.client.connection.pending_requests` instead
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.pending_requests`.
	DbClientConnectionsPendingRequestsName        = "db.client.connections.pending_requests"
	DbClientConnectionsPendingRequestsUnit        = "{request}"
	DbClientConnectionsPendingRequestsDescription = "Deprecated, use `db.client.connection.pending_requests` instead."

	// DbClientConnectionsTimeouts is the metric conforming to the "db.client.connections.timeouts" semantic conventions. It represents the deprecated, use `db.client.connection.timeouts` instead
	// Instrument: counter
	// Unit: {timeout}
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.timeouts`.
	DbClientConnectionsTimeoutsName        = "db.client.connections.timeouts"
	DbClientConnectionsTimeoutsUnit        = "{timeout}"
	DbClientConnectionsTimeoutsDescription = "Deprecated, use `db.client.connection.timeouts` instead."

	// DbClientConnectionsCreateTime is the metric conforming to the "db.client.connections.create_time" semantic conventions. It represents the deprecated, use `db.client.connection.create_time` instead. Note: the unit also changed from `ms` to `s`
	// Instrument: histogram
	// Unit: ms
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.create_time`. Note: the unit also changed from `ms` to `s`.
	DbClientConnectionsCreateTimeName        = "db.client.connections.create_time"
	DbClientConnectionsCreateTimeUnit        = "ms"
	DbClientConnectionsCreateTimeDescription = "Deprecated, use `db.client.connection.create_time` instead. Note: the unit also changed from `ms` to `s`."

	// DbClientConnectionsWaitTime is the metric conforming to the "db.client.connections.wait_time" semantic conventions. It represents the deprecated, use `db.client.connection.wait_time` instead. Note: the unit also changed from `ms` to `s`
	// Instrument: histogram
	// Unit: ms
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.wait_time`. Note: the unit also changed from `ms` to `s`.
	DbClientConnectionsWaitTimeName        = "db.client.connections.wait_time"
	DbClientConnectionsWaitTimeUnit        = "ms"
	DbClientConnectionsWaitTimeDescription = "Deprecated, use `db.client.connection.wait_time` instead. Note: the unit also changed from `ms` to `s`."

	// DbClientConnectionsUseTime is the metric conforming to the "db.client.connections.use_time" semantic conventions. It represents the deprecated, use `db.client.connection.use_time` instead. Note: the unit also changed from `ms` to `s`
	// Instrument: histogram
	// Unit: ms
	// Stability: experimental
	// Deprecated: Replaced by `db.client.connection.use_time`. Note: the unit also changed from `ms` to `s`.
	DbClientConnectionsUseTimeName        = "db.client.connections.use_time"
	DbClientConnectionsUseTimeUnit        = "ms"
	DbClientConnectionsUseTimeDescription = "Deprecated, use `db.client.connection.use_time` instead. Note: the unit also changed from `ms` to `s`."

	// DnsLookupDuration is the metric conforming to the "dns.lookup.duration" semantic conventions. It represents the measures the time taken to perform a DNS lookup
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	DnsLookupDurationName        = "dns.lookup.duration"
	DnsLookupDurationUnit        = "s"
	DnsLookupDurationDescription = "Measures the time taken to perform a DNS lookup."

	// DotnetProcessCpuCount is the metric conforming to the "dotnet.process.cpu.count" semantic conventions. It represents the number of processors available to the process
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: experimental
	DotnetProcessCpuCountName        = "dotnet.process.cpu.count"
	DotnetProcessCpuCountUnit        = "{cpu}"
	DotnetProcessCpuCountDescription = "The number of processors available to the process."

	// DotnetProcessCpuTime is the metric conforming to the "dotnet.process.cpu.time" semantic conventions. It represents the cPU time used by the process
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	DotnetProcessCpuTimeName        = "dotnet.process.cpu.time"
	DotnetProcessCpuTimeUnit        = "s"
	DotnetProcessCpuTimeDescription = "CPU time used by the process."

	// DotnetProcessMemoryWorkingSet is the metric conforming to the "dotnet.process.memory.working_set" semantic conventions. It represents the number of bytes of physical memory mapped to the process context
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	DotnetProcessMemoryWorkingSetName        = "dotnet.process.memory.working_set"
	DotnetProcessMemoryWorkingSetUnit        = "By"
	DotnetProcessMemoryWorkingSetDescription = "The number of bytes of physical memory mapped to the process context."

	// DotnetGcCollections is the metric conforming to the "dotnet.gc.collections" semantic conventions. It represents the number of garbage collections that have occurred since the process has started
	// Instrument: counter
	// Unit: {collection}
	// Stability: experimental
	DotnetGcCollectionsName        = "dotnet.gc.collections"
	DotnetGcCollectionsUnit        = "{collection}"
	DotnetGcCollectionsDescription = "The number of garbage collections that have occurred since the process has started."

	// DotnetGcHeapTotalAllocated is the metric conforming to the "dotnet.gc.heap.total_allocated" semantic conventions. It represents the *approximate* number of bytes allocated on the managed GC heap since the process has started. The returned value does not include any native allocations.
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	DotnetGcHeapTotalAllocatedName        = "dotnet.gc.heap.total_allocated"
	DotnetGcHeapTotalAllocatedUnit        = "By"
	DotnetGcHeapTotalAllocatedDescription = "The *approximate* number of bytes allocated on the managed GC heap since the process has started. The returned value does not include any native allocations."

	// DotnetGcLastCollectionMemoryCommittedSize is the metric conforming to the "dotnet.gc.last_collection.memory.committed_size" semantic conventions. It represents the amount of committed virtual memory in use by the .NET GC, as observed during the latest garbage collection.
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	DotnetGcLastCollectionMemoryCommittedSizeName        = "dotnet.gc.last_collection.memory.committed_size"
	DotnetGcLastCollectionMemoryCommittedSizeUnit        = "By"
	DotnetGcLastCollectionMemoryCommittedSizeDescription = "The amount of committed virtual memory in use by the .NET GC, as observed during the latest garbage collection."

	// DotnetGcLastCollectionHeapSize is the metric conforming to the "dotnet.gc.last_collection.heap.size" semantic conventions. It represents the managed GC heap size (including fragmentation), as observed during the latest garbage collection.
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	DotnetGcLastCollectionHeapSizeName        = "dotnet.gc.last_collection.heap.size"
	DotnetGcLastCollectionHeapSizeUnit        = "By"
	DotnetGcLastCollectionHeapSizeDescription = "The managed GC heap size (including fragmentation), as observed during the latest garbage collection."

	// DotnetGcLastCollectionHeapFragmentationSize is the metric conforming to the "dotnet.gc.last_collection.heap.fragmentation.size" semantic conventions. It represents the heap fragmentation, as observed during the latest garbage collection.
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	DotnetGcLastCollectionHeapFragmentationSizeName        = "dotnet.gc.last_collection.heap.fragmentation.size"
	DotnetGcLastCollectionHeapFragmentationSizeUnit        = "By"
	DotnetGcLastCollectionHeapFragmentationSizeDescription = "The heap fragmentation, as observed during the latest garbage collection."

	// DotnetGcPauseTime is the metric conforming to the "dotnet.gc.pause.time" semantic conventions. It represents the total amount of time paused in GC since the process has started
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	DotnetGcPauseTimeName        = "dotnet.gc.pause.time"
	DotnetGcPauseTimeUnit        = "s"
	DotnetGcPauseTimeDescription = "The total amount of time paused in GC since the process has started."

	// DotnetJitCompiledIlSize is the metric conforming to the "dotnet.jit.compiled_il.size" semantic conventions. It represents the count of bytes of intermediate language that have been compiled since the process has started
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	DotnetJitCompiledIlSizeName        = "dotnet.jit.compiled_il.size"
	DotnetJitCompiledIlSizeUnit        = "By"
	DotnetJitCompiledIlSizeDescription = "Count of bytes of intermediate language that have been compiled since the process has started."

	// DotnetJitCompiledMethods is the metric conforming to the "dotnet.jit.compiled_methods" semantic conventions. It represents the number of times the JIT compiler (re)compiled methods since the process has started.
	// Instrument: counter
	// Unit: {method}
	// Stability: experimental
	DotnetJitCompiledMethodsName        = "dotnet.jit.compiled_methods"
	DotnetJitCompiledMethodsUnit        = "{method}"
	DotnetJitCompiledMethodsDescription = "The number of times the JIT compiler (re)compiled methods since the process has started."

	// DotnetJitCompilationTime is the metric conforming to the "dotnet.jit.compilation.time" semantic conventions. It represents the amount of time the JIT compiler has spent compiling methods since the process has started.
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	DotnetJitCompilationTimeName        = "dotnet.jit.compilation.time"
	DotnetJitCompilationTimeUnit        = "s"
	DotnetJitCompilationTimeDescription = "The amount of time the JIT compiler has spent compiling methods since the process has started."

	// DotnetMonitorLockContentions is the metric conforming to the "dotnet.monitor.lock_contentions" semantic conventions. It represents the number of times there was contention when trying to acquire a monitor lock since the process has started.
	// Instrument: counter
	// Unit: {contention}
	// Stability: experimental
	DotnetMonitorLockContentionsName        = "dotnet.monitor.lock_contentions"
	DotnetMonitorLockContentionsUnit        = "{contention}"
	DotnetMonitorLockContentionsDescription = "The number of times there was contention when trying to acquire a monitor lock since the process has started."

	// DotnetThreadPoolThreadCount is the metric conforming to the "dotnet.thread_pool.thread.count" semantic conventions. It represents the number of thread pool threads that currently exist
	// Instrument: updowncounter
	// Unit: {thread}
	// Stability: experimental
	DotnetThreadPoolThreadCountName        = "dotnet.thread_pool.thread.count"
	DotnetThreadPoolThreadCountUnit        = "{thread}"
	DotnetThreadPoolThreadCountDescription = "The number of thread pool threads that currently exist."

	// DotnetThreadPoolWorkItemCount is the metric conforming to the "dotnet.thread_pool.work_item.count" semantic conventions. It represents the number of work items that the thread pool has completed since the process has started.
	// Instrument: counter
	// Unit: {work_item}
	// Stability: experimental
	DotnetThreadPoolWorkItemCountName        = "dotnet.thread_pool.work_item.count"
	DotnetThreadPoolWorkItemCountUnit        = "{work_item}"
	DotnetThreadPoolWorkItemCountDescription = "The number of work items that the thread pool has completed since the process has started."

	// DotnetThreadPoolQueueLength is the metric conforming to the "dotnet.thread_pool.queue.length" semantic conventions. It represents the number of work items that are currently queued to be processed by the thread pool.
	// Instrument: updowncounter
	// Unit: {work_item}
	// Stability: experimental
	DotnetThreadPoolQueueLengthName        = "dotnet.thread_pool.queue.length"
	DotnetThreadPoolQueueLengthUnit        = "{work_item}"
	DotnetThreadPoolQueueLengthDescription = "The number of work items that are currently queued to be processed by the thread pool."

	// DotnetTimerCount is the metric conforming to the "dotnet.timer.count" semantic conventions. It represents the number of timer instances that are currently active
	// Instrument: updowncounter
	// Unit: {timer}
	// Stability: experimental
	DotnetTimerCountName        = "dotnet.timer.count"
	DotnetTimerCountUnit        = "{timer}"
	DotnetTimerCountDescription = "The number of timer instances that are currently active."

	// DotnetAssemblyCount is the metric conforming to the "dotnet.assembly.count" semantic conventions. It represents the number of .NET assemblies that are currently loaded
	// Instrument: updowncounter
	// Unit: {assembly}
	// Stability: experimental
	DotnetAssemblyCountName        = "dotnet.assembly.count"
	DotnetAssemblyCountUnit        = "{assembly}"
	DotnetAssemblyCountDescription = "The number of .NET assemblies that are currently loaded."

	// DotnetExceptions is the metric conforming to the "dotnet.exceptions" semantic conventions. It represents the number of exceptions that have been thrown in managed code
	// Instrument: counter
	// Unit: {exception}
	// Stability: experimental
	DotnetExceptionsName        = "dotnet.exceptions"
	DotnetExceptionsUnit        = "{exception}"
	DotnetExceptionsDescription = "The number of exceptions that have been thrown in managed code."

	// FaasInvokeDuration is the metric conforming to the "faas.invoke_duration" semantic conventions. It represents the measures the duration of the function's logic execution
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	FaasInvokeDurationName        = "faas.invoke_duration"
	FaasInvokeDurationUnit        = "s"
	FaasInvokeDurationDescription = "Measures the duration of the function's logic execution"

	// FaasInitDuration is the metric conforming to the "faas.init_duration" semantic conventions. It represents the measures the duration of the function's initialization, such as a cold start
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	FaasInitDurationName        = "faas.init_duration"
	FaasInitDurationUnit        = "s"
	FaasInitDurationDescription = "Measures the duration of the function's initialization, such as a cold start"

	// FaasColdstarts is the metric conforming to the "faas.coldstarts" semantic conventions. It represents the number of invocation cold starts
	// Instrument: counter
	// Unit: {coldstart}
	// Stability: experimental
	FaasColdstartsName        = "faas.coldstarts"
	FaasColdstartsUnit        = "{coldstart}"
	FaasColdstartsDescription = "Number of invocation cold starts"

	// FaasErrors is the metric conforming to the "faas.errors" semantic conventions. It represents the number of invocation errors
	// Instrument: counter
	// Unit: {error}
	// Stability: experimental
	FaasErrorsName        = "faas.errors"
	FaasErrorsUnit        = "{error}"
	FaasErrorsDescription = "Number of invocation errors"

	// FaasInvocations is the metric conforming to the "faas.invocations" semantic conventions. It represents the number of successful invocations
	// Instrument: counter
	// Unit: {invocation}
	// Stability: experimental
	FaasInvocationsName        = "faas.invocations"
	FaasInvocationsUnit        = "{invocation}"
	FaasInvocationsDescription = "Number of successful invocations"

	// FaasTimeouts is the metric conforming to the "faas.timeouts" semantic conventions. It represents the number of invocation timeouts
	// Instrument: counter
	// Unit: {timeout}
	// Stability: experimental
	FaasTimeoutsName        = "faas.timeouts"
	FaasTimeoutsUnit        = "{timeout}"
	FaasTimeoutsDescription = "Number of invocation timeouts"

	// FaasMemUsage is the metric conforming to the "faas.mem_usage" semantic conventions. It represents the distribution of max memory usage per invocation
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	FaasMemUsageName        = "faas.mem_usage"
	FaasMemUsageUnit        = "By"
	FaasMemUsageDescription = "Distribution of max memory usage per invocation"

	// FaasCpuUsage is the metric conforming to the "faas.cpu_usage" semantic conventions. It represents the distribution of CPU usage per invocation
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	FaasCpuUsageName        = "faas.cpu_usage"
	FaasCpuUsageUnit        = "s"
	FaasCpuUsageDescription = "Distribution of CPU usage per invocation"

	// FaasNetIo is the metric conforming to the "faas.net_io" semantic conventions. It represents the distribution of net I/O usage per invocation
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	FaasNetIoName        = "faas.net_io"
	FaasNetIoUnit        = "By"
	FaasNetIoDescription = "Distribution of net I/O usage per invocation"

	// GenAiClientTokenUsage is the metric conforming to the "gen_ai.client.token.usage" semantic conventions. It represents the measures number of input and output tokens used
	// Instrument: histogram
	// Unit: {token}
	// Stability: experimental
	GenAiClientTokenUsageName        = "gen_ai.client.token.usage"
	GenAiClientTokenUsageUnit        = "{token}"
	GenAiClientTokenUsageDescription = "Measures number of input and output tokens used"

	// GenAiClientOperationDuration is the metric conforming to the "gen_ai.client.operation.duration" semantic conventions. It represents the genAI operation duration
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	GenAiClientOperationDurationName        = "gen_ai.client.operation.duration"
	GenAiClientOperationDurationUnit        = "s"
	GenAiClientOperationDurationDescription = "GenAI operation duration"

	// GenAiServerRequestDuration is the metric conforming to the "gen_ai.server.request.duration" semantic conventions. It represents the generative AI server request duration such as time-to-last byte or last output token
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	GenAiServerRequestDurationName        = "gen_ai.server.request.duration"
	GenAiServerRequestDurationUnit        = "s"
	GenAiServerRequestDurationDescription = "Generative AI server request duration such as time-to-last byte or last output token"

	// GenAiServerTimePerOutputToken is the metric conforming to the "gen_ai.server.time_per_output_token" semantic conventions. It represents the time per output token generated after the first token for successful responses
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	GenAiServerTimePerOutputTokenName        = "gen_ai.server.time_per_output_token"
	GenAiServerTimePerOutputTokenUnit        = "s"
	GenAiServerTimePerOutputTokenDescription = "Time per output token generated after the first token for successful responses"

	// GenAiServerTimeToFirstToken is the metric conforming to the "gen_ai.server.time_to_first_token" semantic conventions. It represents the time to generate first token for successful responses
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	GenAiServerTimeToFirstTokenName        = "gen_ai.server.time_to_first_token"
	GenAiServerTimeToFirstTokenUnit        = "s"
	GenAiServerTimeToFirstTokenDescription = "Time to generate first token for successful responses"

	// GoMemoryUsed is the metric conforming to the "go.memory.used" semantic conventions. It represents the memory used by the Go runtime
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	GoMemoryUsedName        = "go.memory.used"
	GoMemoryUsedUnit        = "By"
	GoMemoryUsedDescription = "Memory used by the Go runtime."

	// GoMemoryLimit is the metric conforming to the "go.memory.limit" semantic conventions. It represents the go runtime memory limit configured by the user, if a limit exists
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	GoMemoryLimitName        = "go.memory.limit"
	GoMemoryLimitUnit        = "By"
	GoMemoryLimitDescription = "Go runtime memory limit configured by the user, if a limit exists."

	// GoMemoryAllocated is the metric conforming to the "go.memory.allocated" semantic conventions. It represents the memory allocated to the heap by the application
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	GoMemoryAllocatedName        = "go.memory.allocated"
	GoMemoryAllocatedUnit        = "By"
	GoMemoryAllocatedDescription = "Memory allocated to the heap by the application."

	// GoMemoryAllocations is the metric conforming to the "go.memory.allocations" semantic conventions. It represents the count of allocations to the heap by the application
	// Instrument: counter
	// Unit: {allocation}
	// Stability: experimental
	GoMemoryAllocationsName        = "go.memory.allocations"
	GoMemoryAllocationsUnit        = "{allocation}"
	GoMemoryAllocationsDescription = "Count of allocations to the heap by the application."

	// GoMemoryGcGoal is the metric conforming to the "go.memory.gc.goal" semantic conventions. It represents the heap size target for the end of the GC cycle
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	GoMemoryGcGoalName        = "go.memory.gc.goal"
	GoMemoryGcGoalUnit        = "By"
	GoMemoryGcGoalDescription = "Heap size target for the end of the GC cycle."

	// GoGoroutineCount is the metric conforming to the "go.goroutine.count" semantic conventions. It represents the count of live goroutines
	// Instrument: updowncounter
	// Unit: {goroutine}
	// Stability: experimental
	GoGoroutineCountName        = "go.goroutine.count"
	GoGoroutineCountUnit        = "{goroutine}"
	GoGoroutineCountDescription = "Count of live goroutines."

	// GoProcessorLimit is the metric conforming to the "go.processor.limit" semantic conventions. It represents the number of OS threads that can execute user-level Go code simultaneously
	// Instrument: updowncounter
	// Unit: {thread}
	// Stability: experimental
	GoProcessorLimitName        = "go.processor.limit"
	GoProcessorLimitUnit        = "{thread}"
	GoProcessorLimitDescription = "The number of OS threads that can execute user-level Go code simultaneously."

	// GoScheduleDuration is the metric conforming to the "go.schedule.duration" semantic conventions. It represents the time goroutines have spent in the scheduler in a runnable state before actually running
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	GoScheduleDurationName        = "go.schedule.duration"
	GoScheduleDurationUnit        = "s"
	GoScheduleDurationDescription = "The time goroutines have spent in the scheduler in a runnable state before actually running."

	// GoConfigGogc is the metric conforming to the "go.config.gogc" semantic conventions. It represents the heap size target percentage configured by the user, otherwise 100
	// Instrument: updowncounter
	// Unit: %
	// Stability: experimental
	GoConfigGogcName        = "go.config.gogc"
	GoConfigGogcUnit        = "%"
	GoConfigGogcDescription = "Heap size target percentage configured by the user, otherwise 100."

	// HwEnergy is the metric conforming to the "hw.energy" semantic conventions. It represents the energy consumed by the component
	// Instrument: counter
	// Unit: J
	// Stability: experimental
	HwEnergyName        = "hw.energy"
	HwEnergyUnit        = "J"
	HwEnergyDescription = "Energy consumed by the component"

	// HwErrors is the metric conforming to the "hw.errors" semantic conventions. It represents the number of errors encountered by the component
	// Instrument: counter
	// Unit: {error}
	// Stability: experimental
	HwErrorsName        = "hw.errors"
	HwErrorsUnit        = "{error}"
	HwErrorsDescription = "Number of errors encountered by the component"

	// HwPower is the metric conforming to the "hw.power" semantic conventions. It represents the instantaneous power consumed by the component
	// Instrument: gauge
	// Unit: W
	// Stability: experimental
	HwPowerName        = "hw.power"
	HwPowerUnit        = "W"
	HwPowerDescription = "Instantaneous power consumed by the component"

	// HwStatus is the metric conforming to the "hw.status" semantic conventions. It represents the operational status: `1` (true) or `0` (false) for each of the possible states
	// Instrument: updowncounter
	// Unit: 1
	// Stability: experimental
	HwStatusName        = "hw.status"
	HwStatusUnit        = "1"
	HwStatusDescription = "Operational status: `1` (true) or `0` (false) for each of the possible states"

	// HttpServerRequestDuration is the metric conforming to the "http.server.request.duration" semantic conventions. It represents the duration of HTTP server requests
	// Instrument: histogram
	// Unit: s
	// Stability: stable
	HttpServerRequestDurationName        = "http.server.request.duration"
	HttpServerRequestDurationUnit        = "s"
	HttpServerRequestDurationDescription = "Duration of HTTP server requests."

	// HttpServerActiveRequests is the metric conforming to the "http.server.active_requests" semantic conventions. It represents the number of active HTTP server requests
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: experimental
	HttpServerActiveRequestsName        = "http.server.active_requests"
	HttpServerActiveRequestsUnit        = "{request}"
	HttpServerActiveRequestsDescription = "Number of active HTTP server requests."

	// HttpServerRequestBodySize is the metric conforming to the "http.server.request.body.size" semantic conventions. It represents the size of HTTP server request bodies
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	HttpServerRequestBodySizeName        = "http.server.request.body.size"
	HttpServerRequestBodySizeUnit        = "By"
	HttpServerRequestBodySizeDescription = "Size of HTTP server request bodies."

	// HttpServerResponseBodySize is the metric conforming to the "http.server.response.body.size" semantic conventions. It represents the size of HTTP server response bodies
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	HttpServerResponseBodySizeName        = "http.server.response.body.size"
	HttpServerResponseBodySizeUnit        = "By"
	HttpServerResponseBodySizeDescription = "Size of HTTP server response bodies."

	// HttpClientRequestDuration is the metric conforming to the "http.client.request.duration" semantic conventions. It represents the duration of HTTP client requests
	// Instrument: histogram
	// Unit: s
	// Stability: stable
	HttpClientRequestDurationName        = "http.client.request.duration"
	HttpClientRequestDurationUnit        = "s"
	HttpClientRequestDurationDescription = "Duration of HTTP client requests."

	// HttpClientRequestBodySize is the metric conforming to the "http.client.request.body.size" semantic conventions. It represents the size of HTTP client request bodies
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	HttpClientRequestBodySizeName        = "http.client.request.body.size"
	HttpClientRequestBodySizeUnit        = "By"
	HttpClientRequestBodySizeDescription = "Size of HTTP client request bodies."

	// HttpClientResponseBodySize is the metric conforming to the "http.client.response.body.size" semantic conventions. It represents the size of HTTP client response bodies
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	HttpClientResponseBodySizeName        = "http.client.response.body.size"
	HttpClientResponseBodySizeUnit        = "By"
	HttpClientResponseBodySizeDescription = "Size of HTTP client response bodies."

	// HttpClientOpenConnections is the metric conforming to the "http.client.open_connections" semantic conventions. It represents the number of outbound HTTP connections that are currently active or idle on the client
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	HttpClientOpenConnectionsName        = "http.client.open_connections"
	HttpClientOpenConnectionsUnit        = "{connection}"
	HttpClientOpenConnectionsDescription = "Number of outbound HTTP connections that are currently active or idle on the client."

	// HttpClientConnectionDuration is the metric conforming to the "http.client.connection.duration" semantic conventions. It represents the duration of the successfully established outbound HTTP connections
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	HttpClientConnectionDurationName        = "http.client.connection.duration"
	HttpClientConnectionDurationUnit        = "s"
	HttpClientConnectionDurationDescription = "The duration of the successfully established outbound HTTP connections."

	// HttpClientActiveRequests is the metric conforming to the "http.client.active_requests" semantic conventions. It represents the number of active HTTP requests
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: experimental
	HttpClientActiveRequestsName        = "http.client.active_requests"
	HttpClientActiveRequestsUnit        = "{request}"
	HttpClientActiveRequestsDescription = "Number of active HTTP requests."

	// JvmMemoryInit is the metric conforming to the "jvm.memory.init" semantic conventions. It represents the measure of initial memory requested
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	JvmMemoryInitName        = "jvm.memory.init"
	JvmMemoryInitUnit        = "By"
	JvmMemoryInitDescription = "Measure of initial memory requested."

	// JvmSystemCpuUtilization is the metric conforming to the "jvm.system.cpu.utilization" semantic conventions. It represents the recent CPU utilization for the whole system as reported by the JVM
	// Instrument: gauge
	// Unit: 1
	// Stability: experimental
	JvmSystemCpuUtilizationName        = "jvm.system.cpu.utilization"
	JvmSystemCpuUtilizationUnit        = "1"
	JvmSystemCpuUtilizationDescription = "Recent CPU utilization for the whole system as reported by the JVM."

	// JvmSystemCpuLoad1m is the metric conforming to the "jvm.system.cpu.load_1m" semantic conventions. It represents the average CPU load of the whole system for the last minute as reported by the JVM
	// Instrument: gauge
	// Unit: {run_queue_item}
	// Stability: experimental
	JvmSystemCpuLoad1mName        = "jvm.system.cpu.load_1m"
	JvmSystemCpuLoad1mUnit        = "{run_queue_item}"
	JvmSystemCpuLoad1mDescription = "Average CPU load of the whole system for the last minute as reported by the JVM."

	// JvmBufferMemoryUsed is the metric conforming to the "jvm.buffer.memory.used" semantic conventions. It represents the measure of memory used by buffers
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	JvmBufferMemoryUsedName        = "jvm.buffer.memory.used"
	JvmBufferMemoryUsedUnit        = "By"
	JvmBufferMemoryUsedDescription = "Measure of memory used by buffers."

	// JvmBufferMemoryLimit is the metric conforming to the "jvm.buffer.memory.limit" semantic conventions. It represents the measure of total memory capacity of buffers
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	JvmBufferMemoryLimitName        = "jvm.buffer.memory.limit"
	JvmBufferMemoryLimitUnit        = "By"
	JvmBufferMemoryLimitDescription = "Measure of total memory capacity of buffers."

	// JvmBufferCount is the metric conforming to the "jvm.buffer.count" semantic conventions. It represents the number of buffers in the pool
	// Instrument: updowncounter
	// Unit: {buffer}
	// Stability: experimental
	JvmBufferCountName        = "jvm.buffer.count"
	JvmBufferCountUnit        = "{buffer}"
	JvmBufferCountDescription = "Number of buffers in the pool."

	// JvmMemoryUsed is the metric conforming to the "jvm.memory.used" semantic conventions. It represents the measure of memory used
	// Instrument: updowncounter
	// Unit: By
	// Stability: stable
	JvmMemoryUsedName        = "jvm.memory.used"
	JvmMemoryUsedUnit        = "By"
	JvmMemoryUsedDescription = "Measure of memory used."

	// JvmMemoryCommitted is the metric conforming to the "jvm.memory.committed" semantic conventions. It represents the measure of memory committed
	// Instrument: updowncounter
	// Unit: By
	// Stability: stable
	JvmMemoryCommittedName        = "jvm.memory.committed"
	JvmMemoryCommittedUnit        = "By"
	JvmMemoryCommittedDescription = "Measure of memory committed."

	// JvmMemoryLimit is the metric conforming to the "jvm.memory.limit" semantic conventions. It represents the measure of max obtainable memory
	// Instrument: updowncounter
	// Unit: By
	// Stability: stable
	JvmMemoryLimitName        = "jvm.memory.limit"
	JvmMemoryLimitUnit        = "By"
	JvmMemoryLimitDescription = "Measure of max obtainable memory."

	// JvmMemoryUsedAfterLastGc is the metric conforming to the "jvm.memory.used_after_last_gc" semantic conventions. It represents the measure of memory used, as measured after the most recent garbage collection event on this pool
	// Instrument: updowncounter
	// Unit: By
	// Stability: stable
	JvmMemoryUsedAfterLastGcName        = "jvm.memory.used_after_last_gc"
	JvmMemoryUsedAfterLastGcUnit        = "By"
	JvmMemoryUsedAfterLastGcDescription = "Measure of memory used, as measured after the most recent garbage collection event on this pool."

	// JvmGcDuration is the metric conforming to the "jvm.gc.duration" semantic conventions. It represents the duration of JVM garbage collection actions
	// Instrument: histogram
	// Unit: s
	// Stability: stable
	JvmGcDurationName        = "jvm.gc.duration"
	JvmGcDurationUnit        = "s"
	JvmGcDurationDescription = "Duration of JVM garbage collection actions."

	// JvmThreadCount is the metric conforming to the "jvm.thread.count" semantic conventions. It represents the number of executing platform threads
	// Instrument: updowncounter
	// Unit: {thread}
	// Stability: stable
	JvmThreadCountName        = "jvm.thread.count"
	JvmThreadCountUnit        = "{thread}"
	JvmThreadCountDescription = "Number of executing platform threads."

	// JvmClassLoaded is the metric conforming to the "jvm.class.loaded" semantic conventions. It represents the number of classes loaded since JVM start
	// Instrument: counter
	// Unit: {class}
	// Stability: stable
	JvmClassLoadedName        = "jvm.class.loaded"
	JvmClassLoadedUnit        = "{class}"
	JvmClassLoadedDescription = "Number of classes loaded since JVM start."

	// JvmClassUnloaded is the metric conforming to the "jvm.class.unloaded" semantic conventions. It represents the number of classes unloaded since JVM start
	// Instrument: counter
	// Unit: {class}
	// Stability: stable
	JvmClassUnloadedName        = "jvm.class.unloaded"
	JvmClassUnloadedUnit        = "{class}"
	JvmClassUnloadedDescription = "Number of classes unloaded since JVM start."

	// JvmClassCount is the metric conforming to the "jvm.class.count" semantic conventions. It represents the number of classes currently loaded
	// Instrument: updowncounter
	// Unit: {class}
	// Stability: stable
	JvmClassCountName        = "jvm.class.count"
	JvmClassCountUnit        = "{class}"
	JvmClassCountDescription = "Number of classes currently loaded."

	// JvmCpuCount is the metric conforming to the "jvm.cpu.count" semantic conventions. It represents the number of processors available to the Java virtual machine
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: stable
	JvmCpuCountName        = "jvm.cpu.count"
	JvmCpuCountUnit        = "{cpu}"
	JvmCpuCountDescription = "Number of processors available to the Java virtual machine."

	// JvmCpuTime is the metric conforming to the "jvm.cpu.time" semantic conventions. It represents the cPU time used by the process as reported by the JVM
	// Instrument: counter
	// Unit: s
	// Stability: stable
	JvmCpuTimeName        = "jvm.cpu.time"
	JvmCpuTimeUnit        = "s"
	JvmCpuTimeDescription = "CPU time used by the process as reported by the JVM."

	// JvmCpuRecentUtilization is the metric conforming to the "jvm.cpu.recent_utilization" semantic conventions. It represents the recent CPU utilization for the process as reported by the JVM
	// Instrument: gauge
	// Unit: 1
	// Stability: stable
	JvmCpuRecentUtilizationName        = "jvm.cpu.recent_utilization"
	JvmCpuRecentUtilizationUnit        = "1"
	JvmCpuRecentUtilizationDescription = "Recent CPU utilization for the process as reported by the JVM."

	// JvmBufferMemoryUsage is the metric conforming to the "jvm.buffer.memory.usage" semantic conventions. It represents the deprecated, use `jvm.buffer.memory.used` instead
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	// Deprecated: Replaced by `jvm.buffer.memory.used`.
	JvmBufferMemoryUsageName        = "jvm.buffer.memory.usage"
	JvmBufferMemoryUsageUnit        = "By"
	JvmBufferMemoryUsageDescription = "Deprecated, use `jvm.buffer.memory.used` instead."

	// K8sPodCpuTime is the metric conforming to the "k8s.pod.cpu.time" semantic conventions. It represents the total CPU time consumed
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	K8sPodCpuTimeName        = "k8s.pod.cpu.time"
	K8sPodCpuTimeUnit        = "s"
	K8sPodCpuTimeDescription = "Total CPU time consumed"

	// K8sPodCpuUsage is the metric conforming to the "k8s.pod.cpu.usage" semantic conventions. It represents the pod's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs
	// Instrument: gauge
	// Unit: {cpu}
	// Stability: experimental
	K8sPodCpuUsageName        = "k8s.pod.cpu.usage"
	K8sPodCpuUsageUnit        = "{cpu}"
	K8sPodCpuUsageDescription = "Pod's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"

	// K8sNodeCpuTime is the metric conforming to the "k8s.node.cpu.time" semantic conventions. It represents the total CPU time consumed
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	K8sNodeCpuTimeName        = "k8s.node.cpu.time"
	K8sNodeCpuTimeUnit        = "s"
	K8sNodeCpuTimeDescription = "Total CPU time consumed"

	// K8sNodeCpuUsage is the metric conforming to the "k8s.node.cpu.usage" semantic conventions. It represents the node's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs
	// Instrument: gauge
	// Unit: {cpu}
	// Stability: experimental
	K8sNodeCpuUsageName        = "k8s.node.cpu.usage"
	K8sNodeCpuUsageUnit        = "{cpu}"
	K8sNodeCpuUsageDescription = "Node's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"

	// KestrelActiveConnections is the metric conforming to the "kestrel.active_connections" semantic conventions. It represents the number of connections that are currently active on the server
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: stable
	KestrelActiveConnectionsName        = "kestrel.active_connections"
	KestrelActiveConnectionsUnit        = "{connection}"
	KestrelActiveConnectionsDescription = "Number of connections that are currently active on the server."

	// KestrelConnectionDuration is the metric conforming to the "kestrel.connection.duration" semantic conventions. It represents the duration of connections on the server
	// Instrument: histogram
	// Unit: s
	// Stability: stable
	KestrelConnectionDurationName        = "kestrel.connection.duration"
	KestrelConnectionDurationUnit        = "s"
	KestrelConnectionDurationDescription = "The duration of connections on the server."

	// KestrelRejectedConnections is the metric conforming to the "kestrel.rejected_connections" semantic conventions. It represents the number of connections rejected by the server
	// Instrument: counter
	// Unit: {connection}
	// Stability: stable
	KestrelRejectedConnectionsName        = "kestrel.rejected_connections"
	KestrelRejectedConnectionsUnit        = "{connection}"
	KestrelRejectedConnectionsDescription = "Number of connections rejected by the server."

	// KestrelQueuedConnections is the metric conforming to the "kestrel.queued_connections" semantic conventions. It represents the number of connections that are currently queued and are waiting to start
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: stable
	KestrelQueuedConnectionsName        = "kestrel.queued_connections"
	KestrelQueuedConnectionsUnit        = "{connection}"
	KestrelQueuedConnectionsDescription = "Number of connections that are currently queued and are waiting to start."

	// KestrelQueuedRequests is the metric conforming to the "kestrel.queued_requests" semantic conventions. It represents the number of HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are currently queued and are waiting to start
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: stable
	KestrelQueuedRequestsName        = "kestrel.queued_requests"
	KestrelQueuedRequestsUnit        = "{request}"
	KestrelQueuedRequestsDescription = "Number of HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are currently queued and are waiting to start."

	// KestrelUpgradedConnections is the metric conforming to the "kestrel.upgraded_connections" semantic conventions. It represents the number of connections that are currently upgraded (WebSockets).
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: stable
	KestrelUpgradedConnectionsName        = "kestrel.upgraded_connections"
	KestrelUpgradedConnectionsUnit        = "{connection}"
	KestrelUpgradedConnectionsDescription = "Number of connections that are currently upgraded (WebSockets). ."

	// KestrelTlsHandshakeDuration is the metric conforming to the "kestrel.tls_handshake.duration" semantic conventions. It represents the duration of TLS handshakes on the server
	// Instrument: histogram
	// Unit: s
	// Stability: stable
	KestrelTlsHandshakeDurationName        = "kestrel.tls_handshake.duration"
	KestrelTlsHandshakeDurationUnit        = "s"
	KestrelTlsHandshakeDurationDescription = "The duration of TLS handshakes on the server."

	// KestrelActiveTlsHandshakes is the metric conforming to the "kestrel.active_tls_handshakes" semantic conventions. It represents the number of TLS handshakes that are currently in progress on the server
	// Instrument: updowncounter
	// Unit: {handshake}
	// Stability: stable
	KestrelActiveTlsHandshakesName        = "kestrel.active_tls_handshakes"
	KestrelActiveTlsHandshakesUnit        = "{handshake}"
	KestrelActiveTlsHandshakesDescription = "Number of TLS handshakes that are currently in progress on the server."

	// MessagingClientOperationDuration is the metric conforming to the "messaging.client.operation.duration" semantic conventions. It represents the duration of messaging operation initiated by a producer or consumer client
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	MessagingClientOperationDurationName        = "messaging.client.operation.duration"
	MessagingClientOperationDurationUnit        = "s"
	MessagingClientOperationDurationDescription = "Duration of messaging operation initiated by a producer or consumer client."

	// MessagingProcessDuration is the metric conforming to the "messaging.process.duration" semantic conventions. It represents the duration of processing operation
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	MessagingProcessDurationName        = "messaging.process.duration"
	MessagingProcessDurationUnit        = "s"
	MessagingProcessDurationDescription = "Duration of processing operation."

	// MessagingClientPublishedMessages is the metric conforming to the "messaging.client.published.messages" semantic conventions. It represents the number of messages producer attempted to publish to the broker
	// Instrument: counter
	// Unit: {message}
	// Stability: experimental
	MessagingClientPublishedMessagesName        = "messaging.client.published.messages"
	MessagingClientPublishedMessagesUnit        = "{message}"
	MessagingClientPublishedMessagesDescription = "Number of messages producer attempted to publish to the broker."

	// MessagingClientConsumedMessages is the metric conforming to the "messaging.client.consumed.messages" semantic conventions. It represents the number of messages that were delivered to the application
	// Instrument: counter
	// Unit: {message}
	// Stability: experimental
	MessagingClientConsumedMessagesName        = "messaging.client.consumed.messages"
	MessagingClientConsumedMessagesUnit        = "{message}"
	MessagingClientConsumedMessagesDescription = "Number of messages that were delivered to the application."

	// MessagingPublishDuration is the metric conforming to the "messaging.publish.duration" semantic conventions. It represents the deprecated. Use `messaging.client.operation.duration` instead
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	// Deprecated: Replaced by `messaging.client.operation.duration`.
	MessagingPublishDurationName        = "messaging.publish.duration"
	MessagingPublishDurationUnit        = "s"
	MessagingPublishDurationDescription = "Deprecated. Use `messaging.client.operation.duration` instead."

	// MessagingReceiveDuration is the metric conforming to the "messaging.receive.duration" semantic conventions. It represents the deprecated. Use `messaging.client.operation.duration` instead
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	// Deprecated: Replaced by `messaging.client.operation.duration`.
	MessagingReceiveDurationName        = "messaging.receive.duration"
	MessagingReceiveDurationUnit        = "s"
	MessagingReceiveDurationDescription = "Deprecated. Use `messaging.client.operation.duration` instead."

	// MessagingProcessMessages is the metric conforming to the "messaging.process.messages" semantic conventions. It represents the deprecated. Use `messaging.client.consumed.messages` instead
	// Instrument: counter
	// Unit: {message}
	// Stability: experimental
	// Deprecated: Replaced by `messaging.client.consumed.messages`.
	MessagingProcessMessagesName        = "messaging.process.messages"
	MessagingProcessMessagesUnit        = "{message}"
	MessagingProcessMessagesDescription = "Deprecated. Use `messaging.client.consumed.messages` instead."

	// MessagingPublishMessages is the metric conforming to the "messaging.publish.messages" semantic conventions. It represents the deprecated. Use `messaging.client.produced.messages` instead
	// Instrument: counter
	// Unit: {message}
	// Stability: experimental
	// Deprecated: Replaced by `messaging.client.produced.messages`.
	MessagingPublishMessagesName        = "messaging.publish.messages"
	MessagingPublishMessagesUnit        = "{message}"
	MessagingPublishMessagesDescription = "Deprecated. Use `messaging.client.produced.messages` instead."

	// MessagingReceiveMessages is the metric conforming to the "messaging.receive.messages" semantic conventions. It represents the deprecated. Use `messaging.client.consumed.messages` instead
	// Instrument: counter
	// Unit: {message}
	// Stability: experimental
	// Deprecated: Replaced by `messaging.client.consumed.messages`.
	MessagingReceiveMessagesName        = "messaging.receive.messages"
	MessagingReceiveMessagesUnit        = "{message}"
	MessagingReceiveMessagesDescription = "Deprecated. Use `messaging.client.consumed.messages` instead."

	// NodejsEventloopDelayMin is the metric conforming to the "nodejs.eventloop.delay.min" semantic conventions. It represents the event loop minimum delay
	// Instrument: gauge
	// Unit: s
	// Stability: experimental
	NodejsEventloopDelayMinName        = "nodejs.eventloop.delay.min"
	NodejsEventloopDelayMinUnit        = "s"
	NodejsEventloopDelayMinDescription = "Event loop minimum delay."

	// NodejsEventloopDelayMax is the metric conforming to the "nodejs.eventloop.delay.max" semantic conventions. It represents the event loop maximum delay
	// Instrument: gauge
	// Unit: s
	// Stability: experimental
	NodejsEventloopDelayMaxName        = "nodejs.eventloop.delay.max"
	NodejsEventloopDelayMaxUnit        = "s"
	NodejsEventloopDelayMaxDescription = "Event loop maximum delay."

	// NodejsEventloopDelayMean is the metric conforming to the "nodejs.eventloop.delay.mean" semantic conventions. It represents the event loop mean delay
	// Instrument: gauge
	// Unit: s
	// Stability: experimental
	NodejsEventloopDelayMeanName        = "nodejs.eventloop.delay.mean"
	NodejsEventloopDelayMeanUnit        = "s"
	NodejsEventloopDelayMeanDescription = "Event loop mean delay."

	// NodejsEventloopDelayStddev is the metric conforming to the "nodejs.eventloop.delay.stddev" semantic conventions. It represents the event loop standard deviation delay
	// Instrument: gauge
	// Unit: s
	// Stability: experimental
	NodejsEventloopDelayStddevName        = "nodejs.eventloop.delay.stddev"
	NodejsEventloopDelayStddevUnit        = "s"
	NodejsEventloopDelayStddevDescription = "Event loop standard deviation delay."

	// NodejsEventloopDelayP50 is the metric conforming to the "nodejs.eventloop.delay.p50" semantic conventions. It represents the event loop 50 percentile delay
	// Instrument: gauge
	// Unit: s
	// Stability: experimental
	NodejsEventloopDelayP50Name        = "nodejs.eventloop.delay.p50"
	NodejsEventloopDelayP50Unit        = "s"
	NodejsEventloopDelayP50Description = "Event loop 50 percentile delay."

	// NodejsEventloopDelayP90 is the metric conforming to the "nodejs.eventloop.delay.p90" semantic conventions. It represents the event loop 90 percentile delay
	// Instrument: gauge
	// Unit: s
	// Stability: experimental
	NodejsEventloopDelayP90Name        = "nodejs.eventloop.delay.p90"
	NodejsEventloopDelayP90Unit        = "s"
	NodejsEventloopDelayP90Description = "Event loop 90 percentile delay."

	// NodejsEventloopDelayP99 is the metric conforming to the "nodejs.eventloop.delay.p99" semantic conventions. It represents the event loop 99 percentile delay
	// Instrument: gauge
	// Unit: s
	// Stability: experimental
	NodejsEventloopDelayP99Name        = "nodejs.eventloop.delay.p99"
	NodejsEventloopDelayP99Unit        = "s"
	NodejsEventloopDelayP99Description = "Event loop 99 percentile delay."

	// NodejsEventloopUtilization is the metric conforming to the "nodejs.eventloop.utilization" semantic conventions. It represents the event loop utilization
	// Instrument: gauge
	// Unit: 1
	// Stability: experimental
	NodejsEventloopUtilizationName        = "nodejs.eventloop.utilization"
	NodejsEventloopUtilizationUnit        = "1"
	NodejsEventloopUtilizationDescription = "Event loop utilization."

	// NodejsEventloopTime is the metric conforming to the "nodejs.eventloop.time" semantic conventions. It represents the cumulative duration of time the event loop has been in each state
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	NodejsEventloopTimeName        = "nodejs.eventloop.time"
	NodejsEventloopTimeUnit        = "s"
	NodejsEventloopTimeDescription = "Cumulative duration of time the event loop has been in each state."

	// ProcessCpuTime is the metric conforming to the "process.cpu.time" semantic conventions. It represents the total CPU seconds broken down by different states
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	ProcessCpuTimeName        = "process.cpu.time"
	ProcessCpuTimeUnit        = "s"
	ProcessCpuTimeDescription = "Total CPU seconds broken down by different states."

	// ProcessCpuUtilization is the metric conforming to the "process.cpu.utilization" semantic conventions. It represents the difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process
	// Instrument: gauge
	// Unit: 1
	// Stability: experimental
	ProcessCpuUtilizationName        = "process.cpu.utilization"
	ProcessCpuUtilizationUnit        = "1"
	ProcessCpuUtilizationDescription = "Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process."

	// ProcessMemoryUsage is the metric conforming to the "process.memory.usage" semantic conventions. It represents the amount of physical memory in use
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	ProcessMemoryUsageName        = "process.memory.usage"
	ProcessMemoryUsageUnit        = "By"
	ProcessMemoryUsageDescription = "The amount of physical memory in use."

	// ProcessMemoryVirtual is the metric conforming to the "process.memory.virtual" semantic conventions. It represents the amount of committed virtual memory
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	ProcessMemoryVirtualName        = "process.memory.virtual"
	ProcessMemoryVirtualUnit        = "By"
	ProcessMemoryVirtualDescription = "The amount of committed virtual memory."

	// ProcessDiskIo is the metric conforming to the "process.disk.io" semantic conventions. It represents the disk bytes transferred
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	ProcessDiskIoName        = "process.disk.io"
	ProcessDiskIoUnit        = "By"
	ProcessDiskIoDescription = "Disk bytes transferred."

	// ProcessNetworkIo is the metric conforming to the "process.network.io" semantic conventions. It represents the network bytes transferred
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	ProcessNetworkIoName        = "process.network.io"
	ProcessNetworkIoUnit        = "By"
	ProcessNetworkIoDescription = "Network bytes transferred."

	// ProcessThreadCount is the metric conforming to the "process.thread.count" semantic conventions. It represents the process threads count
	// Instrument: updowncounter
	// Unit: {thread}
	// Stability: experimental
	ProcessThreadCountName        = "process.thread.count"
	ProcessThreadCountUnit        = "{thread}"
	ProcessThreadCountDescription = "Process threads count."

	// ProcessOpenFileDescriptorCount is the metric conforming to the "process.open_file_descriptor.count" semantic conventions. It represents the number of file descriptors in use by the process
	// Instrument: updowncounter
	// Unit: {count}
	// Stability: experimental
	ProcessOpenFileDescriptorCountName        = "process.open_file_descriptor.count"
	ProcessOpenFileDescriptorCountUnit        = "{count}"
	ProcessOpenFileDescriptorCountDescription = "Number of file descriptors in use by the process."

	// ProcessContextSwitches is the metric conforming to the "process.context_switches" semantic conventions. It represents the number of times the process has been context switched
	// Instrument: counter
	// Unit: {count}
	// Stability: experimental
	ProcessContextSwitchesName        = "process.context_switches"
	ProcessContextSwitchesUnit        = "{count}"
	ProcessContextSwitchesDescription = "Number of times the process has been context switched."

	// ProcessPagingFaults is the metric conforming to the "process.paging.faults" semantic conventions. It represents the number of page faults the process has made
	// Instrument: counter
	// Unit: {fault}
	// Stability: experimental
	ProcessPagingFaultsName        = "process.paging.faults"
	ProcessPagingFaultsUnit        = "{fault}"
	ProcessPagingFaultsDescription = "Number of page faults the process has made."

	// ProcessUptime is the metric conforming to the "process.uptime" semantic conventions. It represents the time the process has been running
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	ProcessUptimeName        = "process.uptime"
	ProcessUptimeUnit        = "s"
	ProcessUptimeDescription = "The time the process has been running."

	// RpcServerDuration is the metric conforming to the "rpc.server.duration" semantic conventions. It represents the measures the duration of inbound RPC
	// Instrument: histogram
	// Unit: ms
	// Stability: experimental
	RpcServerDurationName        = "rpc.server.duration"
	RpcServerDurationUnit        = "ms"
	RpcServerDurationDescription = "Measures the duration of inbound RPC."

	// RpcServerRequestSize is the metric conforming to the "rpc.server.request.size" semantic conventions. It represents the measures the size of RPC request messages (uncompressed)
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	RpcServerRequestSizeName        = "rpc.server.request.size"
	RpcServerRequestSizeUnit        = "By"
	RpcServerRequestSizeDescription = "Measures the size of RPC request messages (uncompressed)."

	// RpcServerResponseSize is the metric conforming to the "rpc.server.response.size" semantic conventions. It represents the measures the size of RPC response messages (uncompressed)
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	RpcServerResponseSizeName        = "rpc.server.response.size"
	RpcServerResponseSizeUnit        = "By"
	RpcServerResponseSizeDescription = "Measures the size of RPC response messages (uncompressed)."

	// RpcServerRequestsPerRpc is the metric conforming to the "rpc.server.requests_per_rpc" semantic conventions. It represents the measures the number of messages received per RPC
	// Instrument: histogram
	// Unit: {count}
	// Stability: experimental
	RpcServerRequestsPerRpcName        = "rpc.server.requests_per_rpc"
	RpcServerRequestsPerRpcUnit        = "{count}"
	RpcServerRequestsPerRpcDescription = "Measures the number of messages received per RPC."

	// RpcServerResponsesPerRpc is the metric conforming to the "rpc.server.responses_per_rpc" semantic conventions. It represents the measures the number of messages sent per RPC
	// Instrument: histogram
	// Unit: {count}
	// Stability: experimental
	RpcServerResponsesPerRpcName        = "rpc.server.responses_per_rpc"
	RpcServerResponsesPerRpcUnit        = "{count}"
	RpcServerResponsesPerRpcDescription = "Measures the number of messages sent per RPC."

	// RpcClientDuration is the metric conforming to the "rpc.client.duration" semantic conventions. It represents the measures the duration of outbound RPC
	// Instrument: histogram
	// Unit: ms
	// Stability: experimental
	RpcClientDurationName        = "rpc.client.duration"
	RpcClientDurationUnit        = "ms"
	RpcClientDurationDescription = "Measures the duration of outbound RPC."

	// RpcClientRequestSize is the metric conforming to the "rpc.client.request.size" semantic conventions. It represents the measures the size of RPC request messages (uncompressed)
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	RpcClientRequestSizeName        = "rpc.client.request.size"
	RpcClientRequestSizeUnit        = "By"
	RpcClientRequestSizeDescription = "Measures the size of RPC request messages (uncompressed)."

	// RpcClientResponseSize is the metric conforming to the "rpc.client.response.size" semantic conventions. It represents the measures the size of RPC response messages (uncompressed)
	// Instrument: histogram
	// Unit: By
	// Stability: experimental
	RpcClientResponseSizeName        = "rpc.client.response.size"
	RpcClientResponseSizeUnit        = "By"
	RpcClientResponseSizeDescription = "Measures the size of RPC response messages (uncompressed)."

	// RpcClientRequestsPerRpc is the metric conforming to the "rpc.client.requests_per_rpc" semantic conventions. It represents the measures the number of messages received per RPC
	// Instrument: histogram
	// Unit: {count}
	// Stability: experimental
	RpcClientRequestsPerRpcName        = "rpc.client.requests_per_rpc"
	RpcClientRequestsPerRpcUnit        = "{count}"
	RpcClientRequestsPerRpcDescription = "Measures the number of messages received per RPC."

	// RpcClientResponsesPerRpc is the metric conforming to the "rpc.client.responses_per_rpc" semantic conventions. It represents the measures the number of messages sent per RPC
	// Instrument: histogram
	// Unit: {count}
	// Stability: experimental
	RpcClientResponsesPerRpcName        = "rpc.client.responses_per_rpc"
	RpcClientResponsesPerRpcUnit        = "{count}"
	RpcClientResponsesPerRpcDescription = "Measures the number of messages sent per RPC."

	// SignalrServerConnectionDuration is the metric conforming to the "signalr.server.connection.duration" semantic conventions. It represents the duration of connections on the server
	// Instrument: histogram
	// Unit: s
	// Stability: stable
	SignalrServerConnectionDurationName        = "signalr.server.connection.duration"
	SignalrServerConnectionDurationUnit        = "s"
	SignalrServerConnectionDurationDescription = "The duration of connections on the server."

	// SignalrServerActiveConnections is the metric conforming to the "signalr.server.active_connections" semantic conventions. It represents the number of connections that are currently active on the server
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: stable
	SignalrServerActiveConnectionsName        = "signalr.server.active_connections"
	SignalrServerActiveConnectionsUnit        = "{connection}"
	SignalrServerActiveConnectionsDescription = "Number of connections that are currently active on the server."

	// SystemCpuTime is the metric conforming to the "system.cpu.time" semantic conventions. It represents the seconds each logical CPU spent on each mode
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	SystemCpuTimeName        = "system.cpu.time"
	SystemCpuTimeUnit        = "s"
	SystemCpuTimeDescription = "Seconds each logical CPU spent on each mode"

	// SystemCpuUtilization is the metric conforming to the "system.cpu.utilization" semantic conventions. It represents the difference in system.cpu.time since the last measurement, divided by the elapsed time and number of logical CPUs
	// Instrument: gauge
	// Unit: 1
	// Stability: experimental
	SystemCpuUtilizationName        = "system.cpu.utilization"
	SystemCpuUtilizationUnit        = "1"
	SystemCpuUtilizationDescription = "Difference in system.cpu.time since the last measurement, divided by the elapsed time and number of logical CPUs"

	// SystemCpuFrequency is the metric conforming to the "system.cpu.frequency" semantic conventions. It represents the reports the current frequency of the CPU in Hz
	// Instrument: gauge
	// Unit: {Hz}
	// Stability: experimental
	SystemCpuFrequencyName        = "system.cpu.frequency"
	SystemCpuFrequencyUnit        = "{Hz}"
	SystemCpuFrequencyDescription = "Reports the current frequency of the CPU in Hz"

	// SystemCpuPhysicalCount is the metric conforming to the "system.cpu.physical.count" semantic conventions. It represents the reports the number of actual physical processor cores on the hardware
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: experimental
	SystemCpuPhysicalCountName        = "system.cpu.physical.count"
	SystemCpuPhysicalCountUnit        = "{cpu}"
	SystemCpuPhysicalCountDescription = "Reports the number of actual physical processor cores on the hardware"

	// SystemCpuLogicalCount is the metric conforming to the "system.cpu.logical.count" semantic conventions. It represents the reports the number of logical (virtual) processor cores created by the operating system to manage multitasking
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: experimental
	SystemCpuLogicalCountName        = "system.cpu.logical.count"
	SystemCpuLogicalCountUnit        = "{cpu}"
	SystemCpuLogicalCountDescription = "Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking"

	// SystemMemoryUsage is the metric conforming to the "system.memory.usage" semantic conventions. It represents the reports memory in use by state
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	SystemMemoryUsageName        = "system.memory.usage"
	SystemMemoryUsageUnit        = "By"
	SystemMemoryUsageDescription = "Reports memory in use by state."

	// SystemMemoryLimit is the metric conforming to the "system.memory.limit" semantic conventions. It represents the total memory available in the system
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	SystemMemoryLimitName        = "system.memory.limit"
	SystemMemoryLimitUnit        = "By"
	SystemMemoryLimitDescription = "Total memory available in the system."

	// SystemMemoryShared is the metric conforming to the "system.memory.shared" semantic conventions. It represents the shared memory used (mostly by tmpfs)
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	SystemMemorySharedName        = "system.memory.shared"
	SystemMemorySharedUnit        = "By"
	SystemMemorySharedDescription = "Shared memory used (mostly by tmpfs)."

	// SystemMemoryUtilization is the metric conforming to the "system.memory.utilization" semantic conventions
	// Instrument: gauge
	// Unit: 1
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemMemoryUtilizationName = "system.memory.utilization"
	SystemMemoryUtilizationUnit = "1"

	// SystemPagingUsage is the metric conforming to the "system.paging.usage" semantic conventions. It represents the unix swap or windows pagefile usage
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	SystemPagingUsageName        = "system.paging.usage"
	SystemPagingUsageUnit        = "By"
	SystemPagingUsageDescription = "Unix swap or windows pagefile usage"

	// SystemPagingUtilization is the metric conforming to the "system.paging.utilization" semantic conventions
	// Instrument: gauge
	// Unit: 1
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemPagingUtilizationName = "system.paging.utilization"
	SystemPagingUtilizationUnit = "1"

	// SystemPagingFaults is the metric conforming to the "system.paging.faults" semantic conventions
	// Instrument: counter
	// Unit: {fault}
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemPagingFaultsName = "system.paging.faults"
	SystemPagingFaultsUnit = "{fault}"

	// SystemPagingOperations is the metric conforming to the "system.paging.operations" semantic conventions
	// Instrument: counter
	// Unit: {operation}
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemPagingOperationsName = "system.paging.operations"
	SystemPagingOperationsUnit = "{operation}"

	// SystemDiskIo is the metric conforming to the "system.disk.io" semantic conventions
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemDiskIoName = "system.disk.io"
	SystemDiskIoUnit = "By"

	// SystemDiskOperations is the metric conforming to the "system.disk.operations" semantic conventions
	// Instrument: counter
	// Unit: {operation}
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemDiskOperationsName = "system.disk.operations"
	SystemDiskOperationsUnit = "{operation}"

	// SystemDiskIoTime is the metric conforming to the "system.disk.io_time" semantic conventions. It represents the time disk spent activated
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	SystemDiskIoTimeName        = "system.disk.io_time"
	SystemDiskIoTimeUnit        = "s"
	SystemDiskIoTimeDescription = "Time disk spent activated"

	// SystemDiskOperationTime is the metric conforming to the "system.disk.operation_time" semantic conventions. It represents the sum of the time each operation took to complete
	// Instrument: counter
	// Unit: s
	// Stability: experimental
	SystemDiskOperationTimeName        = "system.disk.operation_time"
	SystemDiskOperationTimeUnit        = "s"
	SystemDiskOperationTimeDescription = "Sum of the time each operation took to complete"

	// SystemDiskMerged is the metric conforming to the "system.disk.merged" semantic conventions
	// Instrument: counter
	// Unit: {operation}
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemDiskMergedName = "system.disk.merged"
	SystemDiskMergedUnit = "{operation}"

	// SystemDiskLimit is the metric conforming to the "system.disk.limit" semantic conventions. It represents the total storage capacity of the disk
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	SystemDiskLimitName        = "system.disk.limit"
	SystemDiskLimitUnit        = "By"
	SystemDiskLimitDescription = "The total storage capacity of the disk"

	// SystemFilesystemUsage is the metric conforming to the "system.filesystem.usage" semantic conventions
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemFilesystemUsageName = "system.filesystem.usage"
	SystemFilesystemUsageUnit = "By"

	// SystemFilesystemUtilization is the metric conforming to the "system.filesystem.utilization" semantic conventions
	// Instrument: gauge
	// Unit: 1
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemFilesystemUtilizationName = "system.filesystem.utilization"
	SystemFilesystemUtilizationUnit = "1"

	// SystemNetworkDropped is the metric conforming to the "system.network.dropped" semantic conventions. It represents the count of packets that are dropped or discarded even though there was no error
	// Instrument: counter
	// Unit: {packet}
	// Stability: experimental
	SystemNetworkDroppedName        = "system.network.dropped"
	SystemNetworkDroppedUnit        = "{packet}"
	SystemNetworkDroppedDescription = "Count of packets that are dropped or discarded even though there was no error"

	// SystemNetworkPackets is the metric conforming to the "system.network.packets" semantic conventions
	// Instrument: counter
	// Unit: {packet}
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemNetworkPacketsName = "system.network.packets"
	SystemNetworkPacketsUnit = "{packet}"

	// SystemNetworkErrors is the metric conforming to the "system.network.errors" semantic conventions. It represents the count of network errors detected
	// Instrument: counter
	// Unit: {error}
	// Stability: experimental
	SystemNetworkErrorsName        = "system.network.errors"
	SystemNetworkErrorsUnit        = "{error}"
	SystemNetworkErrorsDescription = "Count of network errors detected"

	// SystemNetworkIo is the metric conforming to the "system.network.io" semantic conventions
	// Instrument: counter
	// Unit: By
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemNetworkIoName = "system.network.io"
	SystemNetworkIoUnit = "By"

	// SystemNetworkConnections is the metric conforming to the "system.network.connections" semantic conventions
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: experimental
	// NOTE: The description (brief) for this metric is not defined in the semantic-conventions repository.
	SystemNetworkConnectionsName = "system.network.connections"
	SystemNetworkConnectionsUnit = "{connection}"

	// SystemProcessCount is the metric conforming to the "system.process.count" semantic conventions. It represents the total number of processes in each state
	// Instrument: updowncounter
	// Unit: {process}
	// Stability: experimental
	SystemProcessCountName        = "system.process.count"
	SystemProcessCountUnit        = "{process}"
	SystemProcessCountDescription = "Total number of processes in each state"

	// SystemProcessCreated is the metric conforming to the "system.process.created" semantic conventions. It represents the total number of processes created over uptime of the host
	// Instrument: counter
	// Unit: {process}
	// Stability: experimental
	SystemProcessCreatedName        = "system.process.created"
	SystemProcessCreatedUnit        = "{process}"
	SystemProcessCreatedDescription = "Total number of processes created over uptime of the host"

	// SystemLinuxMemoryAvailable is the metric conforming to the "system.linux.memory.available" semantic conventions. It represents an estimate of how much memory is available for starting new applications, without causing swapping
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	SystemLinuxMemoryAvailableName        = "system.linux.memory.available"
	SystemLinuxMemoryAvailableUnit        = "By"
	SystemLinuxMemoryAvailableDescription = "An estimate of how much memory is available for starting new applications, without causing swapping"

	// SystemLinuxMemorySlabUsage is the metric conforming to the "system.linux.memory.slab.usage" semantic conventions. It represents the reports the memory used by the Linux kernel for managing caches of frequently used objects
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	SystemLinuxMemorySlabUsageName        = "system.linux.memory.slab.usage"
	SystemLinuxMemorySlabUsageUnit        = "By"
	SystemLinuxMemorySlabUsageDescription = "Reports the memory used by the Linux kernel for managing caches of frequently used objects."

	// V8jsGcDuration is the metric conforming to the "v8js.gc.duration" semantic conventions. It represents the garbage collection duration
	// Instrument: histogram
	// Unit: s
	// Stability: experimental
	V8jsGcDurationName        = "v8js.gc.duration"
	V8jsGcDurationUnit        = "s"
	V8jsGcDurationDescription = "Garbage collection duration."

	// V8jsMemoryHeapLimit is the metric conforming to the "v8js.memory.heap.limit" semantic conventions. It represents the total heap memory size pre-allocated
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	V8jsMemoryHeapLimitName        = "v8js.memory.heap.limit"
	V8jsMemoryHeapLimitUnit        = "By"
	V8jsMemoryHeapLimitDescription = "Total heap memory size pre-allocated."

	// V8jsMemoryHeapUsed is the metric conforming to the "v8js.memory.heap.used" semantic conventions. It represents the heap Memory size allocated
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	V8jsMemoryHeapUsedName        = "v8js.memory.heap.used"
	V8jsMemoryHeapUsedUnit        = "By"
	V8jsMemoryHeapUsedDescription = "Heap Memory size allocated."

	// V8jsHeapSpaceAvailableSize is the metric conforming to the "v8js.heap.space.available_size" semantic conventions. It represents the heap space available size
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	V8jsHeapSpaceAvailableSizeName        = "v8js.heap.space.available_size"
	V8jsHeapSpaceAvailableSizeUnit        = "By"
	V8jsHeapSpaceAvailableSizeDescription = "Heap space available size."

	// V8jsHeapSpacePhysicalSize is the metric conforming to the "v8js.heap.space.physical_size" semantic conventions. It represents the committed size of a heap space
	// Instrument: updowncounter
	// Unit: By
	// Stability: experimental
	V8jsHeapSpacePhysicalSizeName        = "v8js.heap.space.physical_size"
	V8jsHeapSpacePhysicalSizeUnit        = "By"
	V8jsHeapSpacePhysicalSizeDescription = "Committed size of a heap space."
)
