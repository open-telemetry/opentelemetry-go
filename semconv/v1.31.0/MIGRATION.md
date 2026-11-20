# Semantic Convention Changes

The `go.opentelemetry.io/otel/semconv/v1.31.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.30.0` with the following exceptions.

## Dropped deprecations

The following declarations have been deprecated in the [OpenTelemetry Semantic Conventions].
Refer to the respective documentation in that repository for deprecation instructions for each type.

- `CodeNamespace`
- `CodeNamespaceKey`
- `DBClientConnectionsCreateTimeDescription`
- `DBClientConnectionsCreateTimeName`
- `DBClientConnectionsCreateTimeUnit`
- `DBClientConnectionsIdleMaxDescription`
- `DBClientConnectionsIdleMaxName`
- `DBClientConnectionsIdleMaxUnit`
- `DBClientConnectionsIdleMinDescription`
- `DBClientConnectionsIdleMinName`
- `DBClientConnectionsIdleMinUnit`
- `DBClientConnectionsMaxDescription`
- `DBClientConnectionsMaxName`
- `DBClientConnectionsMaxUnit`
- `DBClientConnectionsPendingRequestsDescription`
- `DBClientConnectionsPendingRequestsName`
- `DBClientConnectionsPendingRequestsUnit`
- `DBClientConnectionsTimeoutsDescription`
- `DBClientConnectionsTimeoutsName`
- `DBClientConnectionsTimeoutsUnit`
- `DBClientConnectionsUsageDescription`
- `DBClientConnectionsUsageName`
- `DBClientConnectionsUsageUnit`
- `DBClientConnectionsUseTimeDescription`
- `DBClientConnectionsUseTimeName`
- `DBClientConnectionsUseTimeUnit`
- `DBClientConnectionsWaitTimeDescription`
- `DBClientConnectionsWaitTimeName`
- `DBClientConnectionsWaitTimeUnit`
- `DBClientCosmosDBActiveInstanceCountDescription`
- `DBClientCosmosDBActiveInstanceCountName`
- `DBClientCosmosDBActiveInstanceCountUnit`
- `DBClientCosmosDBOperationRequestChargeDescription`
- `DBClientCosmosDBOperationRequestChargeName`
- `DBClientCosmosDBOperationRequestChargeUnit`
- `GenAIOpenaiRequestResponseFormatJSONObject`
- `GenAIOpenaiRequestResponseFormatJSONSchema`
- `GenAIOpenaiRequestResponseFormatKey`
- `GenAIOpenaiRequestResponseFormatText`
- `MessagingClientPublishedMessagesDescription`
- `MessagingClientPublishedMessagesName`
- `MessagingClientPublishedMessagesUnit`
- `MessagingProcessMessagesDescription`
- `MessagingProcessMessagesName`
- `MessagingProcessMessagesUnit`
- `MessagingPublishDurationDescription`
- `MessagingPublishDurationName`
- `MessagingPublishDurationUnit`
- `MessagingPublishMessagesDescription`
- `MessagingPublishMessagesName`
- `MessagingPublishMessagesUnit`
- `MessagingReceiveDurationDescription`
- `MessagingReceiveDurationName`
- `MessagingReceiveDurationUnit`
- `MessagingReceiveMessagesDescription`
- `MessagingReceiveMessagesName`
- `MessagingReceiveMessagesUnit`
- `SystemCPUFrequencyDescription`
- `SystemCPUFrequencyName`
- `SystemCPUFrequencyUnit`
- `SystemCPUTimeDescription`
- `SystemCPUTimeName`
- `SystemCPUTimeUnit`
- `SystemCPUUtilizationDescription`
- `SystemCPUUtilizationName`
- `SystemCPUUtilizationUnit`

[OpenTelemetry Semantic Conventions]: https://github.com/open-telemetry/semantic-conventions
