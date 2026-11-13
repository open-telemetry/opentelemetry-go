<!-- Generated. DO NOT MODIFY. -->
# Migration from v1.44.0 to v1.45.0

The `go.opentelemetry.io/otel/semconv/v1.45.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.44.0` with the following exceptions.

## Metric API changes

The following changes also affect metric subpackages, which are not included in the root-package symbol comparison below.

### Database connection limit

Rename `dbconv.ClientConnectionMax` and `ClientConnectionMaxObservable` to `ClientConnectionLimit` and `ClientConnectionLimitObservable`, including their `New...` constructors.
The emitted name changes from `db.client.connection.max` to `db.client.connection.limit`; update queries, dashboards, and alerts that select the old name.
The instruments remain `Int64UpDownCounter` and `Int64ObservableUpDownCounter`, with unit `{connection}` and the corresponding constructor option types.

All DB connection-pool metric helpers now provide `AttrSystemName`, `AttrNamespace`, `AttrServerAddress`, and `AttrServerPort`.
Include `db.system.name` and `db.namespace` when available.
Populate server attributes from the client configuration, following the [database connection-pool conventions](https://github.com/open-telemetry/semantic-conventions/blob/v1.45.0/docs/db/database-metrics.md#connection-pools).
The required pool-name argument remains; `ClientConnectionCount.Add` also keeps its required connection-state argument.
Adding these attributes creates distinct metric series, so review aggregations that group by the previous attribute set.

### Sampled absolute values

The new Kubernetes CPU limit/request helpers use `Float64UpDownCounter` and `Float64ObservableUpDownCounter` with unit `{cpu}`.
Keep fractional CPU values as `float64` and use the corresponding `Float64` constructor options.
The new memory limit/request, node PID, and StatefulSet available-pod helpers use `Int64UpDownCounter` and `Int64ObservableUpDownCounter` with their corresponding `Int64` options.
Utilization helpers use `Float64Gauge` and `Float64ObservableGauge` with unit `1`.

Synchronous up/down counters take deltas for each attribute set.
For sampled absolute values, use the observable helpers, such as `k8sconv.NewPodCPULimitCurrentObservable` or `dbconv.NewClientConnectionLimitObservable`, with `metric.WithFloat64Callback` or `metric.WithInt64Callback`, respectively, and observe the current value.
For example, observe successive connection limits of 100 and 120; adding both readings to the synchronous helper would report 220.
If using the synchronous helper, add 100 and then 20 instead.

## RPC conventions

Replace `semconv.RPCResponseStatusCode` and `RPCResponseStatusCodeKey` with `RPCStatusCode` and `RPCStatusCodeKey`.
The emitted key changes from `rpc.response.status_code` to `rpc.status_code`, and the value remains a string.
In `rpcconv.ClientCallDuration` and `ServerCallDuration`, rename `AttrResponseStatusCode` to `AttrStatusCode`.
In `otelconv.SDKExporterOperationDuration`, rename `AttrRPCResponseStatusCode` to `AttrRPCStatusCode`.

Replace `semconv.RPCRequestMetadata(key, values...)` with `RPCRequestHeader(key, values...)`.
Split `RPCResponseMetadata` into `RPCResponseHeader` for metadata before the response payload and `RPCResponseTrailer` for metadata after it.
The emitted prefixes become `rpc.request.header.`, `rpc.response.header.`, and `rpc.response.trailer.`.
Continue to use lowercase header/trailer names and string-slice values, and explicitly configure which fields are captured.
Update telemetry queries that select the old status or metadata keys.

The generic `rpcconv` duration helpers retain `Record(ctx, seconds, systemName, attrs...)`; `systemName` remains the only required attribute argument.
Protocol-specific requirements still apply to the additional attributes.
For example, [gRPC client metrics](https://github.com/open-telemetry/semantic-conventions/blob/v1.45.0/docs/rpc/grpc.md#metric-rpcclientcallduration-grpc) require method, status code, and server address as well as `rpc.system.name="grpc"`.
Supply these through `AttrMethod`, `AttrStatusCode`, and `AttrServerAddress`.
When using `RecordSet`, include all required attributes in the set yourself.

### Duration histogram buckets

`rpcconv.NewClientCallDuration` and `NewServerCallDuration` now advise the [general RPC boundaries](https://github.com/open-telemetry/semantic-conventions/blob/v1.45.0/docs/rpc/rpc-metrics.md):
`[0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 60]`.
Both remain `Float64Histogram` instruments measuring seconds, with `metric.Float64HistogramOption` constructor options.

[gRPC recommends 41 boundaries](https://github.com/open-telemetry/semantic-conventions/blob/v1.45.0/docs/rpc/grpc.md#metrics), including zero.
Select them explicitly with `metric.WithExplicitBucketBoundaries` when constructing the instrument, for example:

```go
import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv/v1.45.0/rpcconv"
)

func newGRPCClientDuration(meter metric.Meter) (rpcconv.ClientCallDuration, error) {
	return rpcconv.NewClientCallDuration(meter,
		metric.WithExplicitBucketBoundaries(
			0, 0.00001, 0.00005, 0.0001, 0.0003, 0.0006, 0.0008,
			0.001, 0.002, 0.003, 0.004, 0.005, 0.006, 0.008, 0.01,
			0.013, 0.016, 0.02, 0.025, 0.03, 0.04, 0.05, 0.065,
			0.08, 0.1, 0.13, 0.16, 0.2, 0.25, 0.3, 0.4, 0.5,
			0.65, 0.8, 1, 2, 5, 10, 20, 50, 100,
		),
	)
}
```

Pass the same option to `NewServerCallDuration` for gRPC server metrics.
Protocol refinements share the `rpc.client.call.duration` and `rpc.server.call.duration` metric names; recording `SystemNameGRPC` does not change bucket advice.
SDK views can override instrument advice.
Review histogram queries and dashboards that rely on specific bucket boundaries.

## Messaging namespaces

Replace `messaging.system` with `messaging.system.name` and root-package `MessagingSystem...` declarations with `MessagingSystemName...` declarations.
In `messagingconv`, replace `SystemAttr` with `SystemNameAttr` and the `System...` values with the corresponding `SystemName...` values.
This changes the required system argument of `ClientConsumedMessages.Add`, `ClientSentMessages.Add`, `ClientOperationDuration.Record`, and `ProcessDuration.Record`.
For observable helpers, replace `AttrSystem` with `AttrSystemName`; update attribute sets passed to `AddSet` and `RecordSet` too.

The following well-known values also change:

| Previous value | New value | New `messagingconv` value |
| --- | --- | --- |
| `aws_sqs` | `aws.sqs` | `SystemNameAWSSQS` |
| `eventgrid` | `azure.eventgrid` | `SystemNameAzureEventGrid` |
| `eventhubs` | `azure.eventhubs` | `SystemNameAzureEventHubs` |
| `servicebus` | `azure.servicebus` | `SystemNameAzureServiceBus` |
| `gcp_pubsub` | `gcp.pubsub` | `SystemNameGCPPubSub` |

Messaging-system-specific attributes move to their own namespaces; the suffixes are unchanged:

| Previous prefix | New prefix | Previous Go declaration prefix | New Go declaration prefix |
| --- | --- | --- | --- |
| `messaging.kafka.` | `kafka.` | `MessagingKafka` | `Kafka` |
| `messaging.rabbitmq.` | `rabbitmq.` | `MessagingRabbitMQ` | `RabbitMQ` |
| `messaging.rocketmq.` | `rocketmq.` | `MessagingRocketMQ` | `RocketMQ` |
| `messaging.gcp_pubsub.` | `gcp.pubsub.` | `MessagingGCPPubSub` | `GCPPubSub` |
| `messaging.servicebus.` | `azure.servicebus.` | `MessagingServiceBus` | `AzureServiceBus` |
| `messaging.eventhubs.` | `azure.eventhubs.` | `MessagingEventHubs` | `AzureEventHubs` |

Update queries, dashboards, and attribute-based filters for both the new keys and the changed system values.
Messaging metric names and instrument types remain unchanged.

## FaaS invocation attributes

`semconv.FaaSInvocationID` and `FaaSInvocationIDKey` keep their Go names but emit `faas.invocation.id` instead of `faas.invocation_id`.
Replace `FaaSTime` and `FaaSTimeKey` with `FaaSInvocationStartTime` and `FaaSInvocationStartTimeKey`, which emit `faas.invocation.start_time` instead of `faas.time`.
The start time remains an ISO 8601 string in UTC.
Update attribute filters and queries for both renamed keys.

## System device and memory attributes

Replace `semconv.SystemDevice` and `SystemDeviceKey` with `SystemDeviceName` and `SystemDeviceNameKey`; the emitted key changes from `system.device` to `system.device.name`.
Rename `systemconv` metric helpers' `AttrDevice` methods to `AttrDeviceName`, including observable variants.
For `containerconv.DiskIO` and `DiskIOObservable`, rename `AttrSystemDevice` to `AttrSystemDeviceName`.
Update telemetry queries that filter or group by the old key.

`SystemMemoryStateFree`, `SystemMemoryStateBuffers`, and `SystemMemoryStateCached`, and the corresponding `systemconv.MemoryState...` values, are removed.
Report `used` and `available` states for `system.memory.usage` and `system.memory.utilization`, using `SystemMemoryStateAvailable` or `systemconv.MemoryStateAvailable` for the new state.
Available memory estimates capacity for new applications without swapping; it is not a direct relabeling of free memory.
On Linux, use `MemAvailable` when available.
The usage values for the two states should sum to `system.memory.limit`; their utilization values should sum to 1.
Update collection logic and state-based queries together.

## Removed

The following declarations have been removed.
Refer to the [OpenTelemetry Semantic Conventions documentation] for deprecation instructions.

If the type is not listed in the documentation as deprecated, it has been removed in this version due to lack of applicability or use.
If you use any of these non-deprecated declarations in your Go application, please [open an issue] describing your use-case.

- `EnduserID`
- `EnduserIDKey`
- `EnduserPseudoID`
- `EnduserPseudoIDKey`
- `FaaSTime`
- `FaaSTimeKey`
- `MessagingEventHubsMessageEnqueuedTime`
- `MessagingEventHubsMessageEnqueuedTimeKey`
- `MessagingGCPPubSubMessageAckDeadline`
- `MessagingGCPPubSubMessageAckDeadlineKey`
- `MessagingGCPPubSubMessageAckID`
- `MessagingGCPPubSubMessageAckIDKey`
- `MessagingGCPPubSubMessageDeliveryAttempt`
- `MessagingGCPPubSubMessageDeliveryAttemptKey`
- `MessagingGCPPubSubMessageOrderingKey`
- `MessagingGCPPubSubMessageOrderingKeyKey`
- `MessagingKafkaClusterID`
- `MessagingKafkaClusterIDKey`
- `MessagingKafkaMessageKey`
- `MessagingKafkaMessageKeyKey`
- `MessagingKafkaMessageTombstone`
- `MessagingKafkaMessageTombstoneKey`
- `MessagingKafkaOffset`
- `MessagingKafkaOffsetKey`
- `MessagingRabbitMQDestinationRoutingKey`
- `MessagingRabbitMQDestinationRoutingKeyKey`
- `MessagingRabbitMQMessageDeliveryTag`
- `MessagingRabbitMQMessageDeliveryTagKey`
- `MessagingRocketMQConsumptionModelBroadcasting`
- `MessagingRocketMQConsumptionModelClustering`
- `MessagingRocketMQConsumptionModelKey`
- `MessagingRocketMQMessageDelayTimeLevel`
- `MessagingRocketMQMessageDelayTimeLevelKey`
- `MessagingRocketMQMessageDeliveryTimestamp`
- `MessagingRocketMQMessageDeliveryTimestampKey`
- `MessagingRocketMQMessageGroup`
- `MessagingRocketMQMessageGroupKey`
- `MessagingRocketMQMessageKeys`
- `MessagingRocketMQMessageKeysKey`
- `MessagingRocketMQMessageTag`
- `MessagingRocketMQMessageTagKey`
- `MessagingRocketMQMessageTypeDelay`
- `MessagingRocketMQMessageTypeFifo`
- `MessagingRocketMQMessageTypeKey`
- `MessagingRocketMQMessageTypeNormal`
- `MessagingRocketMQMessageTypeTransaction`
- `MessagingRocketMQNamespace`
- `MessagingRocketMQNamespaceKey`
- `MessagingServiceBusDispositionStatusAbandon`
- `MessagingServiceBusDispositionStatusComplete`
- `MessagingServiceBusDispositionStatusDeadLetter`
- `MessagingServiceBusDispositionStatusDefer`
- `MessagingServiceBusDispositionStatusKey`
- `MessagingServiceBusMessageDeliveryCount`
- `MessagingServiceBusMessageDeliveryCountKey`
- `MessagingServiceBusMessageEnqueuedTime`
- `MessagingServiceBusMessageEnqueuedTimeKey`
- `MessagingSystemAWSSNS`
- `MessagingSystemAWSSQS`
- `MessagingSystemActiveMQ`
- `MessagingSystemEventGrid`
- `MessagingSystemEventHubs`
- `MessagingSystemGCPPubSub`
- `MessagingSystemJMS`
- `MessagingSystemKafka`
- `MessagingSystemKey`
- `MessagingSystemPulsar`
- `MessagingSystemRabbitMQ`
- `MessagingSystemRocketMQ`
- `MessagingSystemServiceBus`
- `OpenTracingRefTypeChildOf`
- `OpenTracingRefTypeFollowsFrom`
- `OpenTracingRefTypeKey`
- `RPCRequestMetadata`
- `RPCResponseMetadata`
- `RPCResponseStatusCode`
- `RPCResponseStatusCodeKey`
- `ServicePeerName`
- `ServicePeerNameKey`
- `ServicePeerNamespace`
- `ServicePeerNamespaceKey`
- `SystemDevice`
- `SystemDeviceKey`
- `SystemMemoryStateBuffers`
- `SystemMemoryStateCached`
- `SystemMemoryStateFree`
- `ZOSSmfID`
- `ZOSSmfIDKey`
- `ZOSSysplexName`
- `ZOSSysplexNameKey`

[OpenTelemetry Semantic Conventions documentation]: https://github.com/open-telemetry/semantic-conventions
[open an issue]: https://github.com/open-telemetry/opentelemetry-go/issues/new?template=Blank+issue
