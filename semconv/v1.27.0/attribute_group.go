// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.27.0"

import "go.opentelemetry.io/otel/attribute"

// Namespace: android
const (

	// Uniquely identifies the framework API revision offered by a version (`os.version`) of the android operating system. More information can be found [here].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "33",
	// "32",
	//
	// [here]: https://developer.android.com/guide/topics/manifest/uses-sdk-element#ApiLevels
	AndroidOsApiLevelKey = attribute.Key("android.os.api_level")
	// Deprecated use the `device.app.lifecycle` event definition including `android.state` as a payload field instead.
	//
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `device.app.lifecycle`.
	//
	// Examples: undefined
	// Note: The Android lifecycle states are defined in [Activity lifecycle callbacks], and from which the `OS identifiers` are derived
	//
	// [Activity lifecycle callbacks]: https://developer.android.com/guide/components/activities/activity-lifecycle#lc
	AndroidStateKey = attribute.Key("android.state")
)

// AndroidOsApiLevel returns an attribute KeyValue conforming to the "android.os.api_level"semantic conventions
// It represents the Uniquely identifies the framework API revision offered by a version (`os.version`) of the android operating system. More information can be found [here]
//
// [here]: https://developer.android.com/guide/topics/manifest/uses-sdk-element#ApiLevels
func AndroidOsApiLevel(val string) attribute.KeyValue {
	return AndroidOsApiLevelKey.String(val)
}

// AndroidState returns an attribute KeyValue conforming to the "android.state"semantic conventions
// It represents the Deprecated use the `device.app.lifecycle` event definition including `android.state` as a payload field instead
func AndroidState(val string) attribute.KeyValue {
	return AndroidStateKey.String(val)
}

// Enum values for android.state
var (

	// Any time before Activity.onResume() or, if the app has no Activity, Context.startService() has been called in the app for the first time.
	//
	// Stability: experimental
	AndroidStateCreated = AndroidStateKey.String("created")
	// Any time after Activity.onPause() or, if the app has no Activity, Context.stopService() has been called when the app was in the foreground state.
	//
	// Stability: experimental
	AndroidStateBackground = AndroidStateKey.String("background")
	// Any time after Activity.onResume() or, if the app has no Activity, Context.startService() has been called when the app was in either the created or background states.
	//
	// Stability: experimental
	AndroidStateForeground = AndroidStateKey.String("foreground")
)

// Namespace: artifact
const (

	// The provenance filename of the built attestation which directly relates to the build artifact filename. This filename SHOULD accompany the artifact at publish time. See the [SLSA Relationship] specification for more information.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "golang-binary-amd64-v0.1.0.attestation",
	// "docker-image-amd64-v0.1.0.intoto.json1",
	// "release-1.tar.gz.attestation",
	// "file-name-package.tar.gz.intoto.json1",
	//
	// [SLSA Relationship]: https://slsa.dev/spec/v1.0/distributing-provenance#relationship-between-artifacts-and-attestations
	ArtifactAttestationFilenameKey = attribute.Key("artifact.attestation.filename")
	// The full [hash value (see glossary)], of the built attestation. Some envelopes in the software attestation space also refer to this as the [digest].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "1b31dfcd5b7f9267bf2ff47651df1cfb9147b9e4df1f335accf65b4cda498408",
	//
	// [hash value (see glossary)]: https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-5.pdf
	// [digest]: https://github.com/in-toto/attestation/blob/main/spec/README.md#in-toto-attestation-framework-spec
	ArtifactAttestationHashKey = attribute.Key("artifact.attestation.hash")
	// The id of the build [software attestation].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "123",
	//
	// [software attestation]: https://slsa.dev/attestation-model
	ArtifactAttestationIdKey = attribute.Key("artifact.attestation.id")
	// The human readable file name of the artifact, typically generated during build and release processes. Often includes the package name and version in the file name.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "golang-binary-amd64-v0.1.0",
	// "docker-image-amd64-v0.1.0",
	// "release-1.tar.gz",
	// "file-name-package.tar.gz",
	//
	// Note: This file name can also act as the [Package Name]
	// in cases where the package ecosystem maps accordingly.
	// Additionally, the artifact [can be published]
	// for others, but that is not a guarantee
	//
	// [Package Name]: https://slsa.dev/spec/v1.0/terminology#package-model
	// [can be published]: https://slsa.dev/spec/v1.0/terminology#software-supply-chain
	ArtifactFilenameKey = attribute.Key("artifact.filename")
	// The full [hash value (see glossary)], often found in checksum.txt on a release of the artifact and used to verify package integrity.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "9ff4c52759e2c4ac70b7d517bc7fcdc1cda631ca0045271ddd1b192544f8a3e9",
	//
	// Note: The specific algorithm used to create the cryptographic hash value is
	// not defined. In situations where an artifact has multiple
	// cryptographic hashes, it is up to the implementer to choose which
	// hash value to set here; this should be the most secure hash algorithm
	// that is suitable for the situation and consistent with the
	// corresponding attestation. The implementer can then provide the other
	// hash values through an additional set of attribute extensions as they
	// deem necessary
	//
	// [hash value (see glossary)]: https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-5.pdf
	ArtifactHashKey = attribute.Key("artifact.hash")
	// The [Package URL] of the [package artifact] provides a standard way to identify and locate the packaged artifact.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "pkg:github/package-url/purl-spec@1209109710924",
	// "pkg:npm/foo@12.12.3",
	//
	// [Package URL]: https://github.com/package-url/purl-spec
	// [package artifact]: https://slsa.dev/spec/v1.0/terminology#package-model
	ArtifactPurlKey = attribute.Key("artifact.purl")
	// The version of the artifact.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "v0.1.0",
	// "1.2.1",
	// "122691-build",
	ArtifactVersionKey = attribute.Key("artifact.version")
)

// ArtifactAttestationFilename returns an attribute KeyValue conforming to the "artifact.attestation.filename"semantic conventions
// It represents the The provenance filename of the built attestation which directly relates to the build artifact filename. This filename SHOULD accompany the artifact at publish time. See the [SLSA Relationship] specification for more information
//
// [SLSA Relationship]: https://slsa.dev/spec/v1.0/distributing-provenance#relationship-between-artifacts-and-attestations
func ArtifactAttestationFilename(val string) attribute.KeyValue {
	return ArtifactAttestationFilenameKey.String(val)
}

// ArtifactAttestationHash returns an attribute KeyValue conforming to the "artifact.attestation.hash"semantic conventions
// It represents the The full [hash value (see glossary)], of the built attestation. Some envelopes in the software attestation space also refer to this as the [digest]
//
// [hash value (see glossary)]: https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-5.pdf
// [digest]: https://github.com/in-toto/attestation/blob/main/spec/README.md#in-toto-attestation-framework-spec
func ArtifactAttestationHash(val string) attribute.KeyValue {
	return ArtifactAttestationHashKey.String(val)
}

// ArtifactAttestationId returns an attribute KeyValue conforming to the "artifact.attestation.id"semantic conventions
// It represents the The id of the build [software attestation]
//
// [software attestation]: https://slsa.dev/attestation-model
func ArtifactAttestationId(val string) attribute.KeyValue {
	return ArtifactAttestationIdKey.String(val)
}

// ArtifactFilename returns an attribute KeyValue conforming to the "artifact.filename"semantic conventions
// It represents the The human readable file name of the artifact, typically generated during build and release processes. Often includes the package name and version in the file name
func ArtifactFilename(val string) attribute.KeyValue {
	return ArtifactFilenameKey.String(val)
}

// ArtifactHash returns an attribute KeyValue conforming to the "artifact.hash"semantic conventions
// It represents the The full [hash value (see glossary)], often found in checksum.txt on a release of the artifact and used to verify package integrity
//
// [hash value (see glossary)]: https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-5.pdf
func ArtifactHash(val string) attribute.KeyValue {
	return ArtifactHashKey.String(val)
}

// ArtifactPurl returns an attribute KeyValue conforming to the "artifact.purl"semantic conventions
// It represents the The [Package URL] of the [package artifact] provides a standard way to identify and locate the packaged artifact
//
// [Package URL]: https://github.com/package-url/purl-spec
// [package artifact]: https://slsa.dev/spec/v1.0/terminology#package-model
func ArtifactPurl(val string) attribute.KeyValue {
	return ArtifactPurlKey.String(val)
}

// ArtifactVersion returns an attribute KeyValue conforming to the "artifact.version"semantic conventions
// It represents the The version of the artifact
func ArtifactVersion(val string) attribute.KeyValue {
	return ArtifactVersionKey.String(val)
}

// Namespace: aspnetcore
const (

	// ASP.NET Core exception middleware handling result
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "handled",
	// "unhandled",
	AspnetcoreDiagnosticsExceptionResultKey = attribute.Key("aspnetcore.diagnostics.exception.result")
	// Full type name of the [`IExceptionHandler`] implementation that handled the exception.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "Contoso.MyHandler",
	//
	// [`IExceptionHandler`]: https://learn.microsoft.com/dotnet/api/microsoft.aspnetcore.diagnostics.iexceptionhandler
	AspnetcoreDiagnosticsHandlerTypeKey = attribute.Key("aspnetcore.diagnostics.handler.type")
	// Rate limiting policy name.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "fixed",
	// "sliding",
	// "token",
	AspnetcoreRateLimitingPolicyKey = attribute.Key("aspnetcore.rate_limiting.policy")
	// Rate-limiting result, shows whether the lease was acquired or contains a rejection reason
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "acquired",
	// "request_canceled",
	AspnetcoreRateLimitingResultKey = attribute.Key("aspnetcore.rate_limiting.result")
	// Flag indicating if request was handled by the application pipeline.
	// Stability: Stable
	// Type: boolean
	//
	// Examples:
	// true,
	AspnetcoreRequestIsUnhandledKey = attribute.Key("aspnetcore.request.is_unhandled")
	// A value that indicates whether the matched route is a fallback route.
	// Stability: Stable
	// Type: boolean
	//
	// Examples:
	// true,
	AspnetcoreRoutingIsFallbackKey = attribute.Key("aspnetcore.routing.is_fallback")
	// Match result - success or failure
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "success",
	// "failure",
	AspnetcoreRoutingMatchStatusKey = attribute.Key("aspnetcore.routing.match_status")
)

// AspnetcoreDiagnosticsExceptionResult returns an attribute KeyValue conforming to the "aspnetcore.diagnostics.exception.result"semantic conventions
// It represents the ASP.NET Core exception middleware handling result
func AspnetcoreDiagnosticsExceptionResult(val string) attribute.KeyValue {
	return AspnetcoreDiagnosticsExceptionResultKey.String(val)
}

// AspnetcoreDiagnosticsHandlerType returns an attribute KeyValue conforming to the "aspnetcore.diagnostics.handler.type"semantic conventions
// It represents the Full type name of the [`IExceptionHandler`] implementation that handled the exception
//
// [`IExceptionHandler`]: https://learn.microsoft.com/dotnet/api/microsoft.aspnetcore.diagnostics.iexceptionhandler
func AspnetcoreDiagnosticsHandlerType(val string) attribute.KeyValue {
	return AspnetcoreDiagnosticsHandlerTypeKey.String(val)
}

// AspnetcoreRateLimitingPolicy returns an attribute KeyValue conforming to the "aspnetcore.rate_limiting.policy"semantic conventions
// It represents the Rate limiting policy name
func AspnetcoreRateLimitingPolicy(val string) attribute.KeyValue {
	return AspnetcoreRateLimitingPolicyKey.String(val)
}

// AspnetcoreRateLimitingResult returns an attribute KeyValue conforming to the "aspnetcore.rate_limiting.result"semantic conventions
// It represents the Rate-limiting result, shows whether the lease was acquired or contains a rejection reason
func AspnetcoreRateLimitingResult(val string) attribute.KeyValue {
	return AspnetcoreRateLimitingResultKey.String(val)
}

// AspnetcoreRequestIsUnhandled returns an attribute KeyValue conforming to the "aspnetcore.request.is_unhandled"semantic conventions
// It represents the Flag indicating if request was handled by the application pipeline
func AspnetcoreRequestIsUnhandled(val bool) attribute.KeyValue {
	return AspnetcoreRequestIsUnhandledKey.Bool(val)
}

// AspnetcoreRoutingIsFallback returns an attribute KeyValue conforming to the "aspnetcore.routing.is_fallback"semantic conventions
// It represents the A value that indicates whether the matched route is a fallback route
func AspnetcoreRoutingIsFallback(val bool) attribute.KeyValue {
	return AspnetcoreRoutingIsFallbackKey.Bool(val)
}

// AspnetcoreRoutingMatchStatus returns an attribute KeyValue conforming to the "aspnetcore.routing.match_status"semantic conventions
// It represents the Match result - success or failure
func AspnetcoreRoutingMatchStatus(val string) attribute.KeyValue {
	return AspnetcoreRoutingMatchStatusKey.String(val)
}

// Enum values for aspnetcore.diagnostics.exception.result
var (

	// Exception was handled by the exception handling middleware.
	// Stability: stable
	AspnetcoreDiagnosticsExceptionResultHandled = AspnetcoreDiagnosticsExceptionResultKey.String("handled")
	// Exception was not handled by the exception handling middleware.
	// Stability: stable
	AspnetcoreDiagnosticsExceptionResultUnhandled = AspnetcoreDiagnosticsExceptionResultKey.String("unhandled")
	// Exception handling was skipped because the response had started.
	// Stability: stable
	AspnetcoreDiagnosticsExceptionResultSkipped = AspnetcoreDiagnosticsExceptionResultKey.String("skipped")
	// Exception handling didn't run because the request was aborted.
	// Stability: stable
	AspnetcoreDiagnosticsExceptionResultAborted = AspnetcoreDiagnosticsExceptionResultKey.String("aborted")
)

// Enum values for aspnetcore.rate_limiting.result
var (

	// Lease was acquired
	// Stability: stable
	AspnetcoreRateLimitingResultAcquired = AspnetcoreRateLimitingResultKey.String("acquired")
	// Lease request was rejected by the endpoint limiter
	// Stability: stable
	AspnetcoreRateLimitingResultEndpointLimiter = AspnetcoreRateLimitingResultKey.String("endpoint_limiter")
	// Lease request was rejected by the global limiter
	// Stability: stable
	AspnetcoreRateLimitingResultGlobalLimiter = AspnetcoreRateLimitingResultKey.String("global_limiter")
	// Lease request was canceled
	// Stability: stable
	AspnetcoreRateLimitingResultRequestCanceled = AspnetcoreRateLimitingResultKey.String("request_canceled")
)

// Enum values for aspnetcore.routing.match_status
var (

	// Match succeeded
	// Stability: stable
	AspnetcoreRoutingMatchStatusSuccess = AspnetcoreRoutingMatchStatusKey.String("success")
	// Match failed
	// Stability: stable
	AspnetcoreRoutingMatchStatusFailure = AspnetcoreRoutingMatchStatusKey.String("failure")
)

// Namespace: aws
const (

	// The JSON-serialized value of each item in the `AttributeDefinitions` request field.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "{ "AttributeName": "string", "AttributeType": "string" }",
	// ],
	AwsDynamodbAttributeDefinitionsKey = attribute.Key("aws.dynamodb.attribute_definitions")
	// The value of the `AttributesToGet` request parameter.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "lives",
	// "id",
	// ],
	AwsDynamodbAttributesToGetKey = attribute.Key("aws.dynamodb.attributes_to_get")
	// The value of the `ConsistentRead` request parameter.
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	AwsDynamodbConsistentReadKey = attribute.Key("aws.dynamodb.consistent_read")
	// The JSON-serialized value of each item in the `ConsumedCapacity` response field.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "{ "CapacityUnits": number, "GlobalSecondaryIndexes": { "string" : { "CapacityUnits": number, "ReadCapacityUnits": number, "WriteCapacityUnits": number } }, "LocalSecondaryIndexes": { "string" : { "CapacityUnits": number, "ReadCapacityUnits": number, "WriteCapacityUnits": number } }, "ReadCapacityUnits": number, "Table": { "CapacityUnits": number, "ReadCapacityUnits": number, "WriteCapacityUnits": number }, "TableName": "string", "WriteCapacityUnits": number }",
	// ],
	AwsDynamodbConsumedCapacityKey = attribute.Key("aws.dynamodb.consumed_capacity")
	// The value of the `Count` response parameter.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 10,
	AwsDynamodbCountKey = attribute.Key("aws.dynamodb.count")
	// The value of the `ExclusiveStartTableName` request parameter.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Users",
	// "CatsTable",
	AwsDynamodbExclusiveStartTableKey = attribute.Key("aws.dynamodb.exclusive_start_table")
	// The JSON-serialized value of each item in the `GlobalSecondaryIndexUpdates` request field.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "{ "Create": { "IndexName": "string", "KeySchema": [ { "AttributeName": "string", "KeyType": "string" } ], "Projection": { "NonKeyAttributes": [ "string" ], "ProjectionType": "string" }, "ProvisionedThroughput": { "ReadCapacityUnits": number, "WriteCapacityUnits": number } }",
	// ],
	AwsDynamodbGlobalSecondaryIndexUpdatesKey = attribute.Key("aws.dynamodb.global_secondary_index_updates")
	// The JSON-serialized value of each item of the `GlobalSecondaryIndexes` request field
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "{ "IndexName": "string", "KeySchema": [ { "AttributeName": "string", "KeyType": "string" } ], "Projection": { "NonKeyAttributes": [ "string" ], "ProjectionType": "string" }, "ProvisionedThroughput": { "ReadCapacityUnits": number, "WriteCapacityUnits": number } }",
	// ],
	AwsDynamodbGlobalSecondaryIndexesKey = attribute.Key("aws.dynamodb.global_secondary_indexes")
	// The value of the `IndexName` request parameter.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "name_to_group",
	AwsDynamodbIndexNameKey = attribute.Key("aws.dynamodb.index_name")
	// The JSON-serialized value of the `ItemCollectionMetrics` response field.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "{ "string" : [ { "ItemCollectionKey": { "string" : { "B": blob, "BOOL": boolean, "BS": [ blob ], "L": [ "AttributeValue" ], "M": { "string" : "AttributeValue" }, "N": "string", "NS": [ "string" ], "NULL": boolean, "S": "string", "SS": [ "string" ] } }, "SizeEstimateRangeGB": [ number ] } ] }",
	AwsDynamodbItemCollectionMetricsKey = attribute.Key("aws.dynamodb.item_collection_metrics")
	// The value of the `Limit` request parameter.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 10,
	AwsDynamodbLimitKey = attribute.Key("aws.dynamodb.limit")
	// The JSON-serialized value of each item of the `LocalSecondaryIndexes` request field.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "{ "IndexArn": "string", "IndexName": "string", "IndexSizeBytes": number, "ItemCount": number, "KeySchema": [ { "AttributeName": "string", "KeyType": "string" } ], "Projection": { "NonKeyAttributes": [ "string" ], "ProjectionType": "string" } }",
	// ],
	AwsDynamodbLocalSecondaryIndexesKey = attribute.Key("aws.dynamodb.local_secondary_indexes")
	// The value of the `ProjectionExpression` request parameter.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Title",
	// "Title, Price, Color",
	// "Title, Description, RelatedItems, ProductReviews",
	AwsDynamodbProjectionKey = attribute.Key("aws.dynamodb.projection")
	// The value of the `ProvisionedThroughput.ReadCapacityUnits` request parameter.
	// Stability: Experimental
	// Type: double
	//
	// Examples:
	// 1.0,
	// 2.0,
	AwsDynamodbProvisionedReadCapacityKey = attribute.Key("aws.dynamodb.provisioned_read_capacity")
	// The value of the `ProvisionedThroughput.WriteCapacityUnits` request parameter.
	// Stability: Experimental
	// Type: double
	//
	// Examples:
	// 1.0,
	// 2.0,
	AwsDynamodbProvisionedWriteCapacityKey = attribute.Key("aws.dynamodb.provisioned_write_capacity")
	// The value of the `ScanIndexForward` request parameter.
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	AwsDynamodbScanForwardKey = attribute.Key("aws.dynamodb.scan_forward")
	// The value of the `ScannedCount` response parameter.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 50,
	AwsDynamodbScannedCountKey = attribute.Key("aws.dynamodb.scanned_count")
	// The value of the `Segment` request parameter.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 10,
	AwsDynamodbSegmentKey = attribute.Key("aws.dynamodb.segment")
	// The value of the `Select` request parameter.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "ALL_ATTRIBUTES",
	// "COUNT",
	AwsDynamodbSelectKey = attribute.Key("aws.dynamodb.select")
	// The number of items in the `TableNames` response parameter.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 20,
	AwsDynamodbTableCountKey = attribute.Key("aws.dynamodb.table_count")
	// The keys in the `RequestItems` object field.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "Users",
	// "Cats",
	// ],
	AwsDynamodbTableNamesKey = attribute.Key("aws.dynamodb.table_names")
	// The value of the `TotalSegments` request parameter.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 100,
	AwsDynamodbTotalSegmentsKey = attribute.Key("aws.dynamodb.total_segments")
	// The ARN of an [ECS cluster].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "arn:aws:ecs:us-west-2:123456789123:cluster/my-cluster",
	//
	// [ECS cluster]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/clusters.html
	AwsEcsClusterArnKey = attribute.Key("aws.ecs.cluster.arn")
	// The Amazon Resource Name (ARN) of an [ECS container instance].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "arn:aws:ecs:us-west-1:123456789123:container/32624152-9086-4f0e-acae-1a75b14fe4d9",
	//
	// [ECS container instance]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_instances.html
	AwsEcsContainerArnKey = attribute.Key("aws.ecs.container.arn")
	// The [launch type] for an ECS task.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	//
	// [launch type]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_types.html
	AwsEcsLaunchtypeKey = attribute.Key("aws.ecs.launchtype")
	// The ARN of a running [ECS task].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "arn:aws:ecs:us-west-1:123456789123:task/10838bed-421f-43ef-870a-f43feacbbb5b",
	// "arn:aws:ecs:us-west-1:123456789123:task/my-cluster/task-id/23ebb8ac-c18f-46c6-8bbe-d55d0e37cfbd",
	//
	// [ECS task]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-account-settings.html#ecs-resource-ids
	AwsEcsTaskArnKey = attribute.Key("aws.ecs.task.arn")
	// The family name of the [ECS task definition] used to create the ECS task.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry-family",
	//
	// [ECS task definition]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definitions.html
	AwsEcsTaskFamilyKey = attribute.Key("aws.ecs.task.family")
	// The ID of a running ECS task. The ID MUST be extracted from `task.arn`.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "10838bed-421f-43ef-870a-f43feacbbb5b",
	// "23ebb8ac-c18f-46c6-8bbe-d55d0e37cfbd",
	AwsEcsTaskIdKey = attribute.Key("aws.ecs.task.id")
	// The revision for the task definition used to create the ECS task.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "8",
	// "26",
	AwsEcsTaskRevisionKey = attribute.Key("aws.ecs.task.revision")
	// The ARN of an EKS cluster.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "arn:aws:ecs:us-west-2:123456789123:cluster/my-cluster",
	AwsEksClusterArnKey = attribute.Key("aws.eks.cluster.arn")
	// The full invoked ARN as provided on the `Context` passed to the function (`Lambda-Runtime-Invoked-Function-Arn` header on the `/runtime/invocation/next` applicable).
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "arn:aws:lambda:us-east-1:123456:function:myfunction:myalias",
	//
	// Note: This may be different from `cloud.resource_id` if an alias is involved
	AwsLambdaInvokedArnKey = attribute.Key("aws.lambda.invoked_arn")
	// The Amazon Resource Name(s) (ARN) of the AWS log group(s).
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "arn:aws:logs:us-west-1:123456789012:log-group:/aws/my/group:*",
	// ],
	//
	// Note: See the [log group ARN format documentation]
	//
	// [log group ARN format documentation]: https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/iam-access-control-overview-cwl.html#CWL_ARN_Format
	AwsLogGroupArnsKey = attribute.Key("aws.log.group.arns")
	// The name(s) of the AWS log group(s) an application is writing to.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "/aws/lambda/my-function",
	// "opentelemetry-service",
	// ],
	//
	// Note: Multiple log groups must be supported for cases like multi-container applications, where a single application has sidecar containers, and each write to their own log group
	AwsLogGroupNamesKey = attribute.Key("aws.log.group.names")
	// The ARN(s) of the AWS log stream(s).
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "arn:aws:logs:us-west-1:123456789012:log-group:/aws/my/group:log-stream:logs/main/10838bed-421f-43ef-870a-f43feacbbb5b",
	// ],
	//
	// Note: See the [log stream ARN format documentation]. One log group can contain several log streams, so these ARNs necessarily identify both a log group and a log stream
	//
	// [log stream ARN format documentation]: https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/iam-access-control-overview-cwl.html#CWL_ARN_Format
	AwsLogStreamArnsKey = attribute.Key("aws.log.stream.arns")
	// The name(s) of the AWS log stream(s) an application is writing to.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "logs/main/10838bed-421f-43ef-870a-f43feacbbb5b",
	// ],
	AwsLogStreamNamesKey = attribute.Key("aws.log.stream.names")
	// The AWS request ID as returned in the response headers `x-amz-request-id` or `x-amz-requestid`.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "79b9da39-b7ae-508a-a6bc-864b2829c622",
	// "C9ER4AJX75574TDJ",
	AwsRequestIdKey = attribute.Key("aws.request_id")
	// The S3 bucket name the request refers to. Corresponds to the `--bucket` parameter of the [S3 API] operations.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "some-bucket-name",
	//
	// Note: The `bucket` attribute is applicable to all S3 operations that reference a bucket, i.e. that require the bucket name as a mandatory parameter.
	// This applies to almost all S3 operations except `list-buckets`
	//
	// [S3 API]: https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html
	AwsS3BucketKey = attribute.Key("aws.s3.bucket")
	// The source object (in the form `bucket`/`key`) for the copy operation.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "someFile.yml",
	//
	// Note: The `copy_source` attribute applies to S3 copy operations and corresponds to the `--copy-source` parameter
	// of the [copy-object operation within the S3 API].
	// This applies in particular to the following operations:
	//
	//   - [copy-object]
	//   - [upload-part-copy]
	// [copy-object operation within the S3 API]: https://docs.aws.amazon.com/cli/latest/reference/s3api/copy-object.html
	// [copy-object]: https://docs.aws.amazon.com/cli/latest/reference/s3api/copy-object.html
	// [upload-part-copy]: https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part-copy.html
	AwsS3CopySourceKey = attribute.Key("aws.s3.copy_source")
	// The delete request container that specifies the objects to be deleted.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Objects=[{Key=string,VersionId=string},{Key=string,VersionId=string}],Quiet=boolean",
	//
	// Note: The `delete` attribute is only applicable to the [delete-object] operation.
	// The `delete` attribute corresponds to the `--delete` parameter of the
	// [delete-objects operation within the S3 API]
	//
	// [delete-object]: https://docs.aws.amazon.com/cli/latest/reference/s3api/delete-object.html
	// [delete-objects operation within the S3 API]: https://docs.aws.amazon.com/cli/latest/reference/s3api/delete-objects.html
	AwsS3DeleteKey = attribute.Key("aws.s3.delete")
	// The S3 object key the request refers to. Corresponds to the `--key` parameter of the [S3 API] operations.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "someFile.yml",
	//
	// Note: The `key` attribute is applicable to all object-related S3 operations, i.e. that require the object key as a mandatory parameter.
	// This applies in particular to the following operations:
	//
	//   - [copy-object]
	//   - [delete-object]
	//   - [get-object]
	//   - [head-object]
	//   - [put-object]
	//   - [restore-object]
	//   - [select-object-content]
	//   - [abort-multipart-upload]
	//   - [complete-multipart-upload]
	//   - [create-multipart-upload]
	//   - [list-parts]
	//   - [upload-part]
	//   - [upload-part-copy]
	// [S3 API]: https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html
	// [copy-object]: https://docs.aws.amazon.com/cli/latest/reference/s3api/copy-object.html
	// [delete-object]: https://docs.aws.amazon.com/cli/latest/reference/s3api/delete-object.html
	// [get-object]: https://docs.aws.amazon.com/cli/latest/reference/s3api/get-object.html
	// [head-object]: https://docs.aws.amazon.com/cli/latest/reference/s3api/head-object.html
	// [put-object]: https://docs.aws.amazon.com/cli/latest/reference/s3api/put-object.html
	// [restore-object]: https://docs.aws.amazon.com/cli/latest/reference/s3api/restore-object.html
	// [select-object-content]: https://docs.aws.amazon.com/cli/latest/reference/s3api/select-object-content.html
	// [abort-multipart-upload]: https://docs.aws.amazon.com/cli/latest/reference/s3api/abort-multipart-upload.html
	// [complete-multipart-upload]: https://docs.aws.amazon.com/cli/latest/reference/s3api/complete-multipart-upload.html
	// [create-multipart-upload]: https://docs.aws.amazon.com/cli/latest/reference/s3api/create-multipart-upload.html
	// [list-parts]: https://docs.aws.amazon.com/cli/latest/reference/s3api/list-parts.html
	// [upload-part]: https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part.html
	// [upload-part-copy]: https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part-copy.html
	AwsS3KeyKey = attribute.Key("aws.s3.key")
	// The part number of the part being uploaded in a multipart-upload operation. This is a positive integer between 1 and 10,000.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 3456,
	//
	// Note: The `part_number` attribute is only applicable to the [upload-part]
	// and [upload-part-copy] operations.
	// The `part_number` attribute corresponds to the `--part-number` parameter of the
	// [upload-part operation within the S3 API]
	//
	// [upload-part]: https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part.html
	// [upload-part-copy]: https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part-copy.html
	// [upload-part operation within the S3 API]: https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part.html
	AwsS3PartNumberKey = attribute.Key("aws.s3.part_number")
	// Upload ID that identifies the multipart upload.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "dfRtDYWFbkRONycy.Yxwh66Yjlx.cph0gtNBtJ",
	//
	// Note: The `upload_id` attribute applies to S3 multipart-upload operations and corresponds to the `--upload-id` parameter
	// of the [S3 API] multipart operations.
	// This applies in particular to the following operations:
	//
	//   - [abort-multipart-upload]
	//   - [complete-multipart-upload]
	//   - [list-parts]
	//   - [upload-part]
	//   - [upload-part-copy]
	// [S3 API]: https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html
	// [abort-multipart-upload]: https://docs.aws.amazon.com/cli/latest/reference/s3api/abort-multipart-upload.html
	// [complete-multipart-upload]: https://docs.aws.amazon.com/cli/latest/reference/s3api/complete-multipart-upload.html
	// [list-parts]: https://docs.aws.amazon.com/cli/latest/reference/s3api/list-parts.html
	// [upload-part]: https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part.html
	// [upload-part-copy]: https://docs.aws.amazon.com/cli/latest/reference/s3api/upload-part-copy.html
	AwsS3UploadIdKey = attribute.Key("aws.s3.upload_id")
)

// AwsDynamodbAttributeDefinitions returns an attribute KeyValue conforming to the "aws.dynamodb.attribute_definitions"semantic conventions
// It represents the The JSON-serialized value of each item in the `AttributeDefinitions` request field
func AwsDynamodbAttributeDefinitions(val []string) attribute.KeyValue {
	return AwsDynamodbAttributeDefinitionsKey.StringSlice(val)
}

// AwsDynamodbAttributesToGet returns an attribute KeyValue conforming to the "aws.dynamodb.attributes_to_get"semantic conventions
// It represents the The value of the `AttributesToGet` request parameter
func AwsDynamodbAttributesToGet(val []string) attribute.KeyValue {
	return AwsDynamodbAttributesToGetKey.StringSlice(val)
}

// AwsDynamodbConsistentRead returns an attribute KeyValue conforming to the "aws.dynamodb.consistent_read"semantic conventions
// It represents the The value of the `ConsistentRead` request parameter
func AwsDynamodbConsistentRead(val bool) attribute.KeyValue {
	return AwsDynamodbConsistentReadKey.Bool(val)
}

// AwsDynamodbConsumedCapacity returns an attribute KeyValue conforming to the "aws.dynamodb.consumed_capacity"semantic conventions
// It represents the The JSON-serialized value of each item in the `ConsumedCapacity` response field
func AwsDynamodbConsumedCapacity(val []string) attribute.KeyValue {
	return AwsDynamodbConsumedCapacityKey.StringSlice(val)
}

// AwsDynamodbCount returns an attribute KeyValue conforming to the "aws.dynamodb.count"semantic conventions
// It represents the The value of the `Count` response parameter
func AwsDynamodbCount(val int) attribute.KeyValue {
	return AwsDynamodbCountKey.Int(val)
}

// AwsDynamodbExclusiveStartTable returns an attribute KeyValue conforming to the "aws.dynamodb.exclusive_start_table"semantic conventions
// It represents the The value of the `ExclusiveStartTableName` request parameter
func AwsDynamodbExclusiveStartTable(val string) attribute.KeyValue {
	return AwsDynamodbExclusiveStartTableKey.String(val)
}

// AwsDynamodbGlobalSecondaryIndexUpdates returns an attribute KeyValue conforming to the "aws.dynamodb.global_secondary_index_updates"semantic conventions
// It represents the The JSON-serialized value of each item in the `GlobalSecondaryIndexUpdates` request field
func AwsDynamodbGlobalSecondaryIndexUpdates(val []string) attribute.KeyValue {
	return AwsDynamodbGlobalSecondaryIndexUpdatesKey.StringSlice(val)
}

// AwsDynamodbGlobalSecondaryIndexes returns an attribute KeyValue conforming to the "aws.dynamodb.global_secondary_indexes"semantic conventions
// It represents the The JSON-serialized value of each item of the `GlobalSecondaryIndexes` request field
func AwsDynamodbGlobalSecondaryIndexes(val []string) attribute.KeyValue {
	return AwsDynamodbGlobalSecondaryIndexesKey.StringSlice(val)
}

// AwsDynamodbIndexName returns an attribute KeyValue conforming to the "aws.dynamodb.index_name"semantic conventions
// It represents the The value of the `IndexName` request parameter
func AwsDynamodbIndexName(val string) attribute.KeyValue {
	return AwsDynamodbIndexNameKey.String(val)
}

// AwsDynamodbItemCollectionMetrics returns an attribute KeyValue conforming to the "aws.dynamodb.item_collection_metrics"semantic conventions
// It represents the The JSON-serialized value of the `ItemCollectionMetrics` response field
func AwsDynamodbItemCollectionMetrics(val string) attribute.KeyValue {
	return AwsDynamodbItemCollectionMetricsKey.String(val)
}

// AwsDynamodbLimit returns an attribute KeyValue conforming to the "aws.dynamodb.limit"semantic conventions
// It represents the The value of the `Limit` request parameter
func AwsDynamodbLimit(val int) attribute.KeyValue {
	return AwsDynamodbLimitKey.Int(val)
}

// AwsDynamodbLocalSecondaryIndexes returns an attribute KeyValue conforming to the "aws.dynamodb.local_secondary_indexes"semantic conventions
// It represents the The JSON-serialized value of each item of the `LocalSecondaryIndexes` request field
func AwsDynamodbLocalSecondaryIndexes(val []string) attribute.KeyValue {
	return AwsDynamodbLocalSecondaryIndexesKey.StringSlice(val)
}

// AwsDynamodbProjection returns an attribute KeyValue conforming to the "aws.dynamodb.projection"semantic conventions
// It represents the The value of the `ProjectionExpression` request parameter
func AwsDynamodbProjection(val string) attribute.KeyValue {
	return AwsDynamodbProjectionKey.String(val)
}

// AwsDynamodbProvisionedReadCapacity returns an attribute KeyValue conforming to the "aws.dynamodb.provisioned_read_capacity"semantic conventions
// It represents the The value of the `ProvisionedThroughput.ReadCapacityUnits` request parameter
func AwsDynamodbProvisionedReadCapacity(val float64) attribute.KeyValue {
	return AwsDynamodbProvisionedReadCapacityKey.Float64(val)
}

// AwsDynamodbProvisionedWriteCapacity returns an attribute KeyValue conforming to the "aws.dynamodb.provisioned_write_capacity"semantic conventions
// It represents the The value of the `ProvisionedThroughput.WriteCapacityUnits` request parameter
func AwsDynamodbProvisionedWriteCapacity(val float64) attribute.KeyValue {
	return AwsDynamodbProvisionedWriteCapacityKey.Float64(val)
}

// AwsDynamodbScanForward returns an attribute KeyValue conforming to the "aws.dynamodb.scan_forward"semantic conventions
// It represents the The value of the `ScanIndexForward` request parameter
func AwsDynamodbScanForward(val bool) attribute.KeyValue {
	return AwsDynamodbScanForwardKey.Bool(val)
}

// AwsDynamodbScannedCount returns an attribute KeyValue conforming to the "aws.dynamodb.scanned_count"semantic conventions
// It represents the The value of the `ScannedCount` response parameter
func AwsDynamodbScannedCount(val int) attribute.KeyValue {
	return AwsDynamodbScannedCountKey.Int(val)
}

// AwsDynamodbSegment returns an attribute KeyValue conforming to the "aws.dynamodb.segment"semantic conventions
// It represents the The value of the `Segment` request parameter
func AwsDynamodbSegment(val int) attribute.KeyValue {
	return AwsDynamodbSegmentKey.Int(val)
}

// AwsDynamodbSelect returns an attribute KeyValue conforming to the "aws.dynamodb.select"semantic conventions
// It represents the The value of the `Select` request parameter
func AwsDynamodbSelect(val string) attribute.KeyValue {
	return AwsDynamodbSelectKey.String(val)
}

// AwsDynamodbTableCount returns an attribute KeyValue conforming to the "aws.dynamodb.table_count"semantic conventions
// It represents the The number of items in the `TableNames` response parameter
func AwsDynamodbTableCount(val int) attribute.KeyValue {
	return AwsDynamodbTableCountKey.Int(val)
}

// AwsDynamodbTableNames returns an attribute KeyValue conforming to the "aws.dynamodb.table_names"semantic conventions
// It represents the The keys in the `RequestItems` object field
func AwsDynamodbTableNames(val []string) attribute.KeyValue {
	return AwsDynamodbTableNamesKey.StringSlice(val)
}

// AwsDynamodbTotalSegments returns an attribute KeyValue conforming to the "aws.dynamodb.total_segments"semantic conventions
// It represents the The value of the `TotalSegments` request parameter
func AwsDynamodbTotalSegments(val int) attribute.KeyValue {
	return AwsDynamodbTotalSegmentsKey.Int(val)
}

// AwsEcsClusterArn returns an attribute KeyValue conforming to the "aws.ecs.cluster.arn"semantic conventions
// It represents the The ARN of an [ECS cluster]
//
// [ECS cluster]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/clusters.html
func AwsEcsClusterArn(val string) attribute.KeyValue {
	return AwsEcsClusterArnKey.String(val)
}

// AwsEcsContainerArn returns an attribute KeyValue conforming to the "aws.ecs.container.arn"semantic conventions
// It represents the The Amazon Resource Name (ARN) of an [ECS container instance]
//
// [ECS container instance]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_instances.html
func AwsEcsContainerArn(val string) attribute.KeyValue {
	return AwsEcsContainerArnKey.String(val)
}

// AwsEcsLaunchtype returns an attribute KeyValue conforming to the "aws.ecs.launchtype"semantic conventions
// It represents the The [launch type] for an ECS task
//
// [launch type]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_types.html
func AwsEcsLaunchtype(val string) attribute.KeyValue {
	return AwsEcsLaunchtypeKey.String(val)
}

// AwsEcsTaskArn returns an attribute KeyValue conforming to the "aws.ecs.task.arn"semantic conventions
// It represents the The ARN of a running [ECS task]
//
// [ECS task]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-account-settings.html#ecs-resource-ids
func AwsEcsTaskArn(val string) attribute.KeyValue {
	return AwsEcsTaskArnKey.String(val)
}

// AwsEcsTaskFamily returns an attribute KeyValue conforming to the "aws.ecs.task.family"semantic conventions
// It represents the The family name of the [ECS task definition] used to create the ECS task
//
// [ECS task definition]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definitions.html
func AwsEcsTaskFamily(val string) attribute.KeyValue {
	return AwsEcsTaskFamilyKey.String(val)
}

// AwsEcsTaskId returns an attribute KeyValue conforming to the "aws.ecs.task.id"semantic conventions
// It represents the The ID of a running ECS task. The ID MUST be extracted from `task.arn`
func AwsEcsTaskId(val string) attribute.KeyValue {
	return AwsEcsTaskIdKey.String(val)
}

// AwsEcsTaskRevision returns an attribute KeyValue conforming to the "aws.ecs.task.revision"semantic conventions
// It represents the The revision for the task definition used to create the ECS task
func AwsEcsTaskRevision(val string) attribute.KeyValue {
	return AwsEcsTaskRevisionKey.String(val)
}

// AwsEksClusterArn returns an attribute KeyValue conforming to the "aws.eks.cluster.arn"semantic conventions
// It represents the The ARN of an EKS cluster
func AwsEksClusterArn(val string) attribute.KeyValue {
	return AwsEksClusterArnKey.String(val)
}

// AwsLambdaInvokedArn returns an attribute KeyValue conforming to the "aws.lambda.invoked_arn"semantic conventions
// It represents the The full invoked ARN as provided on the `Context` passed to the function (`Lambda-Runtime-Invoked-Function-Arn` header on the `/runtime/invocation/next` applicable)
func AwsLambdaInvokedArn(val string) attribute.KeyValue {
	return AwsLambdaInvokedArnKey.String(val)
}

// AwsLogGroupArns returns an attribute KeyValue conforming to the "aws.log.group.arns"semantic conventions
// It represents the The Amazon Resource Name(s) (ARN) of the AWS log group(s)
func AwsLogGroupArns(val []string) attribute.KeyValue {
	return AwsLogGroupArnsKey.StringSlice(val)
}

// AwsLogGroupNames returns an attribute KeyValue conforming to the "aws.log.group.names"semantic conventions
// It represents the The name(s) of the AWS log group(s) an application is writing to
func AwsLogGroupNames(val []string) attribute.KeyValue {
	return AwsLogGroupNamesKey.StringSlice(val)
}

// AwsLogStreamArns returns an attribute KeyValue conforming to the "aws.log.stream.arns"semantic conventions
// It represents the The ARN(s) of the AWS log stream(s)
func AwsLogStreamArns(val []string) attribute.KeyValue {
	return AwsLogStreamArnsKey.StringSlice(val)
}

// AwsLogStreamNames returns an attribute KeyValue conforming to the "aws.log.stream.names"semantic conventions
// It represents the The name(s) of the AWS log stream(s) an application is writing to
func AwsLogStreamNames(val []string) attribute.KeyValue {
	return AwsLogStreamNamesKey.StringSlice(val)
}

// AwsRequestId returns an attribute KeyValue conforming to the "aws.request_id"semantic conventions
// It represents the The AWS request ID as returned in the response headers `x-amz-request-id` or `x-amz-requestid`
func AwsRequestId(val string) attribute.KeyValue {
	return AwsRequestIdKey.String(val)
}

// AwsS3Bucket returns an attribute KeyValue conforming to the "aws.s3.bucket"semantic conventions
// It represents the The S3 bucket name the request refers to. Corresponds to the `--bucket` parameter of the [S3 API] operations
//
// [S3 API]: https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html
func AwsS3Bucket(val string) attribute.KeyValue {
	return AwsS3BucketKey.String(val)
}

// AwsS3CopySource returns an attribute KeyValue conforming to the "aws.s3.copy_source"semantic conventions
// It represents the The source object (in the form `bucket`/`key`) for the copy operation
func AwsS3CopySource(val string) attribute.KeyValue {
	return AwsS3CopySourceKey.String(val)
}

// AwsS3Delete returns an attribute KeyValue conforming to the "aws.s3.delete"semantic conventions
// It represents the The delete request container that specifies the objects to be deleted
func AwsS3Delete(val string) attribute.KeyValue {
	return AwsS3DeleteKey.String(val)
}

// AwsS3Key returns an attribute KeyValue conforming to the "aws.s3.key"semantic conventions
// It represents the The S3 object key the request refers to. Corresponds to the `--key` parameter of the [S3 API] operations
//
// [S3 API]: https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html
func AwsS3Key(val string) attribute.KeyValue {
	return AwsS3KeyKey.String(val)
}

// AwsS3PartNumber returns an attribute KeyValue conforming to the "aws.s3.part_number"semantic conventions
// It represents the The part number of the part being uploaded in a multipart-upload operation. This is a positive integer between 1 and 10,000
func AwsS3PartNumber(val int) attribute.KeyValue {
	return AwsS3PartNumberKey.Int(val)
}

// AwsS3UploadId returns an attribute KeyValue conforming to the "aws.s3.upload_id"semantic conventions
// It represents the Upload ID that identifies the multipart upload
func AwsS3UploadId(val string) attribute.KeyValue {
	return AwsS3UploadIdKey.String(val)
}

// Enum values for aws.ecs.launchtype
var (

	// ec2
	// Stability: experimental
	AwsEcsLaunchtypeEc2 = AwsEcsLaunchtypeKey.String("ec2")
	// fargate
	// Stability: experimental
	AwsEcsLaunchtypeFargate = AwsEcsLaunchtypeKey.String("fargate")
)

// Namespace: az
const (

	// [Azure Resource Provider Namespace] as recognized by the client.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Microsoft.Storage",
	// "Microsoft.KeyVault",
	// "Microsoft.ServiceBus",
	//
	// [Azure Resource Provider Namespace]: https://learn.microsoft.com/azure/azure-resource-manager/management/azure-services-resource-providers
	AzNamespaceKey = attribute.Key("az.namespace")
	// The unique identifier of the service request. It's generated by the Azure service and returned with the response.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "00000000-0000-0000-0000-000000000000",
	AzServiceRequestIdKey = attribute.Key("az.service_request_id")
)

// AzNamespace returns an attribute KeyValue conforming to the "az.namespace"semantic conventions
// It represents the [Azure Resource Provider Namespace] as recognized by the client
//
// [Azure Resource Provider Namespace]: https://learn.microsoft.com/azure/azure-resource-manager/management/azure-services-resource-providers
func AzNamespace(val string) attribute.KeyValue {
	return AzNamespaceKey.String(val)
}

// AzServiceRequestId returns an attribute KeyValue conforming to the "az.service_request_id"semantic conventions
// It represents the The unique identifier of the service request. It's generated by the Azure service and returned with the response
func AzServiceRequestId(val string) attribute.KeyValue {
	return AzServiceRequestIdKey.String(val)
}

// Namespace: browser
const (

	// Array of brand name and version separated by a space
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// " Not A;Brand 99",
	// "Chromium 99",
	// "Chrome 99",
	// ],
	//
	// Note: This value is intended to be taken from the [UA client hints API] (`navigator.userAgentData.brands`)
	//
	// [UA client hints API]: https://wicg.github.io/ua-client-hints/#interface
	BrowserBrandsKey = attribute.Key("browser.brands")
	// Preferred language of the user using the browser
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "en",
	// "en-US",
	// "fr",
	// "fr-FR",
	//
	// Note: This value is intended to be taken from the Navigator API `navigator.language`
	BrowserLanguageKey = attribute.Key("browser.language")
	// A boolean that is true if the browser is running on a mobile device
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	// Note: This value is intended to be taken from the [UA client hints API] (`navigator.userAgentData.mobile`). If unavailable, this attribute SHOULD be left unset
	//
	// [UA client hints API]: https://wicg.github.io/ua-client-hints/#interface
	BrowserMobileKey = attribute.Key("browser.mobile")
	// The platform on which the browser is running
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Windows",
	// "macOS",
	// "Android",
	//
	// Note: This value is intended to be taken from the [UA client hints API] (`navigator.userAgentData.platform`). If unavailable, the legacy `navigator.platform` API SHOULD NOT be used instead and this attribute SHOULD be left unset in order for the values to be consistent.
	// The list of possible values is defined in the [W3C User-Agent Client Hints specification]. Note that some (but not all) of these values can overlap with values in the [`os.type` and `os.name` attributes]. However, for consistency, the values in the `browser.platform` attribute should capture the exact value that the user agent provides
	//
	// [UA client hints API]: https://wicg.github.io/ua-client-hints/#interface
	// [W3C User-Agent Client Hints specification]: https://wicg.github.io/ua-client-hints/#sec-ch-ua-platform
	// [`os.type` and `os.name` attributes]: ./os.md
	BrowserPlatformKey = attribute.Key("browser.platform")
)

// BrowserBrands returns an attribute KeyValue conforming to the "browser.brands"semantic conventions
// It represents the Array of brand name and version separated by a space
func BrowserBrands(val []string) attribute.KeyValue {
	return BrowserBrandsKey.StringSlice(val)
}

// BrowserLanguage returns an attribute KeyValue conforming to the "browser.language"semantic conventions
// It represents the Preferred language of the user using the browser
func BrowserLanguage(val string) attribute.KeyValue {
	return BrowserLanguageKey.String(val)
}

// BrowserMobile returns an attribute KeyValue conforming to the "browser.mobile"semantic conventions
// It represents the A boolean that is true if the browser is running on a mobile device
func BrowserMobile(val bool) attribute.KeyValue {
	return BrowserMobileKey.Bool(val)
}

// BrowserPlatform returns an attribute KeyValue conforming to the "browser.platform"semantic conventions
// It represents the The platform on which the browser is running
func BrowserPlatform(val string) attribute.KeyValue {
	return BrowserPlatformKey.String(val)
}

// Namespace: cicd
const (

	// The human readable name of the pipeline within a CI/CD system.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Build and Test",
	// "Lint",
	// "Deploy Go Project",
	// "deploy_to_environment",
	CicdPipelineNameKey = attribute.Key("cicd.pipeline.name")
	// The unique identifier of a pipeline run within a CI/CD system.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "120912",
	CicdPipelineRunIdKey = attribute.Key("cicd.pipeline.run.id")
	// The human readable name of a task within a pipeline. Task here most closely aligns with a [computing process] in a pipeline. Other terms for tasks include commands, steps, and procedures.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Run GoLang Linter",
	// "Go Build",
	// "go-test",
	// "deploy_binary",
	//
	// [computing process]: https://en.wikipedia.org/wiki/Pipeline_(computing)
	CicdPipelineTaskNameKey = attribute.Key("cicd.pipeline.task.name")
	// The unique identifier of a task run within a pipeline.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "12097",
	CicdPipelineTaskRunIdKey = attribute.Key("cicd.pipeline.task.run.id")
	// The [URL] of the pipeline run providing the complete address in order to locate and identify the pipeline run.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "https://github.com/open-telemetry/semantic-conventions/actions/runs/9753949763/job/26920038674?pr=1075",
	//
	// [URL]: https://en.wikipedia.org/wiki/URL
	CicdPipelineTaskRunUrlFullKey = attribute.Key("cicd.pipeline.task.run.url.full")
	// The type of the task within a pipeline.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "build",
	// "test",
	// "deploy",
	CicdPipelineTaskTypeKey = attribute.Key("cicd.pipeline.task.type")
)

// CicdPipelineName returns an attribute KeyValue conforming to the "cicd.pipeline.name"semantic conventions
// It represents the The human readable name of the pipeline within a CI/CD system
func CicdPipelineName(val string) attribute.KeyValue {
	return CicdPipelineNameKey.String(val)
}

// CicdPipelineRunId returns an attribute KeyValue conforming to the "cicd.pipeline.run.id"semantic conventions
// It represents the The unique identifier of a pipeline run within a CI/CD system
func CicdPipelineRunId(val string) attribute.KeyValue {
	return CicdPipelineRunIdKey.String(val)
}

// CicdPipelineTaskName returns an attribute KeyValue conforming to the "cicd.pipeline.task.name"semantic conventions
// It represents the The human readable name of a task within a pipeline. Task here most closely aligns with a [computing process] in a pipeline. Other terms for tasks include commands, steps, and procedures
//
// [computing process]: https://en.wikipedia.org/wiki/Pipeline_(computing)
func CicdPipelineTaskName(val string) attribute.KeyValue {
	return CicdPipelineTaskNameKey.String(val)
}

// CicdPipelineTaskRunId returns an attribute KeyValue conforming to the "cicd.pipeline.task.run.id"semantic conventions
// It represents the The unique identifier of a task run within a pipeline
func CicdPipelineTaskRunId(val string) attribute.KeyValue {
	return CicdPipelineTaskRunIdKey.String(val)
}

// CicdPipelineTaskRunUrlFull returns an attribute KeyValue conforming to the "cicd.pipeline.task.run.url.full"semantic conventions
// It represents the The [URL] of the pipeline run providing the complete address in order to locate and identify the pipeline run
//
// [URL]: https://en.wikipedia.org/wiki/URL
func CicdPipelineTaskRunUrlFull(val string) attribute.KeyValue {
	return CicdPipelineTaskRunUrlFullKey.String(val)
}

// CicdPipelineTaskType returns an attribute KeyValue conforming to the "cicd.pipeline.task.type"semantic conventions
// It represents the The type of the task within a pipeline
func CicdPipelineTaskType(val string) attribute.KeyValue {
	return CicdPipelineTaskTypeKey.String(val)
}

// Enum values for cicd.pipeline.task.type
var (

	// build
	// Stability: experimental
	CicdPipelineTaskTypeBuild = CicdPipelineTaskTypeKey.String("build")
	// test
	// Stability: experimental
	CicdPipelineTaskTypeTest = CicdPipelineTaskTypeKey.String("test")
	// deploy
	// Stability: experimental
	CicdPipelineTaskTypeDeploy = CicdPipelineTaskTypeKey.String("deploy")
)

// Namespace: client
const (

	// Client address - domain name if available without reverse DNS lookup; otherwise, IP address or Unix domain socket name.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "client.example.com",
	// "10.1.2.80",
	// "/tmp/my.sock",
	//
	// Note: When observed from the server side, and when communicating through an intermediary, `client.address` SHOULD represent the client address behind any intermediaries,  for example proxies, if it's available
	ClientAddressKey = attribute.Key("client.address")
	// Client port number.
	// Stability: Stable
	// Type: int
	//
	// Examples:
	// 65123,
	//
	// Note: When observed from the server side, and when communicating through an intermediary, `client.port` SHOULD represent the client port behind any intermediaries,  for example proxies, if it's available
	ClientPortKey = attribute.Key("client.port")
)

// ClientAddress returns an attribute KeyValue conforming to the "client.address"semantic conventions
// It represents the Client address - domain name if available without reverse DNS lookup; otherwise, IP address or Unix domain socket name
func ClientAddress(val string) attribute.KeyValue {
	return ClientAddressKey.String(val)
}

// ClientPort returns an attribute KeyValue conforming to the "client.port"semantic conventions
// It represents the Client port number
func ClientPort(val int) attribute.KeyValue {
	return ClientPortKey.Int(val)
}

// Namespace: cloud
const (

	// The cloud account ID the resource is assigned to.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "111111111111",
	// "opentelemetry",
	CloudAccountIdKey = attribute.Key("cloud.account.id")
	// Cloud regions often have multiple, isolated locations known as zones to increase availability. Availability zone represents the zone where the resource is running.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "us-east-1c",
	//
	// Note: Availability zones are called "zones" on Alibaba Cloud and Google Cloud
	CloudAvailabilityZoneKey = attribute.Key("cloud.availability_zone")
	// The cloud platform in use.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: The prefix of the service SHOULD match the one specified in `cloud.provider`
	CloudPlatformKey = attribute.Key("cloud.platform")
	// Name of the cloud provider.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	CloudProviderKey = attribute.Key("cloud.provider")
	// The geographical region the resource is running.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "us-central1",
	// "us-east-1",
	//
	// Note: Refer to your provider's docs to see the available regions, for example [Alibaba Cloud regions], [AWS regions], [Azure regions], [Google Cloud regions], or [Tencent Cloud regions]
	//
	// [Alibaba Cloud regions]: https://www.alibabacloud.com/help/doc-detail/40654.htm
	// [AWS regions]: https://aws.amazon.com/about-aws/global-infrastructure/regions_az/
	// [Azure regions]: https://azure.microsoft.com/global-infrastructure/geographies/
	// [Google Cloud regions]: https://cloud.google.com/about/locations
	// [Tencent Cloud regions]: https://www.tencentcloud.com/document/product/213/6091
	CloudRegionKey = attribute.Key("cloud.region")
	// Cloud provider-specific native identifier of the monitored cloud resource (e.g. an [ARN] on AWS, a [fully qualified resource ID] on Azure, a [full resource name] on GCP)
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "arn:aws:lambda:REGION:ACCOUNT_ID:function:my-function",
	// "//run.googleapis.com/projects/PROJECT_ID/locations/LOCATION_ID/services/SERVICE_ID",
	// "/subscriptions/<SUBSCRIPTION_GUID>/resourceGroups/<RG>/providers/Microsoft.Web/sites/<FUNCAPP>/functions/<FUNC>",
	//
	// Note: On some cloud providers, it may not be possible to determine the full ID at startup,
	// so it may be necessary to set `cloud.resource_id` as a span attribute instead.
	//
	// The exact value to use for `cloud.resource_id` depends on the cloud provider.
	// The following well-known definitions MUST be used if you set this attribute and they apply:
	//
	//   - **AWS Lambda:** The function [ARN].
	//     Take care not to use the "invoked ARN" directly but replace any
	//     [alias suffix]
	//     with the resolved function version, as the same runtime instance may be invocable with
	//     multiple different aliases.
	//   - **GCP:** The [URI of the resource]
	//   - **Azure:** The [Fully Qualified Resource ID] of the invoked function,
	//     *not* the function app, having the form
	//     `/subscriptions/<SUBSCRIPTION_GUID>/resourceGroups/<RG>/providers/Microsoft.Web/sites/<FUNCAPP>/functions/<FUNC>`.
	//     This means that a span attribute MUST be used, as an Azure function app can host multiple functions that would usually share
	//     a TracerProvider
	// [ARN]: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
	// [fully qualified resource ID]: https://learn.microsoft.com/rest/api/resources/resources/get-by-id
	// [full resource name]: https://cloud.google.com/apis/design/resource_names#full_resource_name
	// [ARN]: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
	// [alias suffix]: https://docs.aws.amazon.com/lambda/latest/dg/configuration-aliases.html
	// [URI of the resource]: https://cloud.google.com/iam/docs/full-resource-names
	// [Fully Qualified Resource ID]: https://docs.microsoft.com/rest/api/resources/resources/get-by-id
	CloudResourceIdKey = attribute.Key("cloud.resource_id")
)

// CloudAccountId returns an attribute KeyValue conforming to the "cloud.account.id"semantic conventions
// It represents the The cloud account ID the resource is assigned to
func CloudAccountId(val string) attribute.KeyValue {
	return CloudAccountIdKey.String(val)
}

// CloudAvailabilityZone returns an attribute KeyValue conforming to the "cloud.availability_zone"semantic conventions
// It represents the Cloud regions often have multiple, isolated locations known as zones to increase availability. Availability zone represents the zone where the resource is running
func CloudAvailabilityZone(val string) attribute.KeyValue {
	return CloudAvailabilityZoneKey.String(val)
}

// CloudPlatform returns an attribute KeyValue conforming to the "cloud.platform"semantic conventions
// It represents the The cloud platform in use
func CloudPlatform(val string) attribute.KeyValue {
	return CloudPlatformKey.String(val)
}

// CloudProvider returns an attribute KeyValue conforming to the "cloud.provider"semantic conventions
// It represents the Name of the cloud provider
func CloudProvider(val string) attribute.KeyValue {
	return CloudProviderKey.String(val)
}

// CloudRegion returns an attribute KeyValue conforming to the "cloud.region"semantic conventions
// It represents the The geographical region the resource is running
func CloudRegion(val string) attribute.KeyValue {
	return CloudRegionKey.String(val)
}

// CloudResourceId returns an attribute KeyValue conforming to the "cloud.resource_id"semantic conventions
// It represents the Cloud provider-specific native identifier of the monitored cloud resource (e.g. an [ARN] on AWS, a [fully qualified resource ID] on Azure, a [full resource name] on GCP)
//
// [ARN]: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
// [fully qualified resource ID]: https://learn.microsoft.com/rest/api/resources/resources/get-by-id
// [full resource name]: https://cloud.google.com/apis/design/resource_names#full_resource_name
func CloudResourceId(val string) attribute.KeyValue {
	return CloudResourceIdKey.String(val)
}

// Enum values for cloud.platform
var (

	// Alibaba Cloud Elastic Compute Service
	// Stability: experimental
	CloudPlatformAlibabaCloudEcs = CloudPlatformKey.String("alibaba_cloud_ecs")
	// Alibaba Cloud Function Compute
	// Stability: experimental
	CloudPlatformAlibabaCloudFc = CloudPlatformKey.String("alibaba_cloud_fc")
	// Red Hat OpenShift on Alibaba Cloud
	// Stability: experimental
	CloudPlatformAlibabaCloudOpenshift = CloudPlatformKey.String("alibaba_cloud_openshift")
	// AWS Elastic Compute Cloud
	// Stability: experimental
	CloudPlatformAwsEc2 = CloudPlatformKey.String("aws_ec2")
	// AWS Elastic Container Service
	// Stability: experimental
	CloudPlatformAwsEcs = CloudPlatformKey.String("aws_ecs")
	// AWS Elastic Kubernetes Service
	// Stability: experimental
	CloudPlatformAwsEks = CloudPlatformKey.String("aws_eks")
	// AWS Lambda
	// Stability: experimental
	CloudPlatformAwsLambda = CloudPlatformKey.String("aws_lambda")
	// AWS Elastic Beanstalk
	// Stability: experimental
	CloudPlatformAwsElasticBeanstalk = CloudPlatformKey.String("aws_elastic_beanstalk")
	// AWS App Runner
	// Stability: experimental
	CloudPlatformAwsAppRunner = CloudPlatformKey.String("aws_app_runner")
	// Red Hat OpenShift on AWS (ROSA)
	// Stability: experimental
	CloudPlatformAwsOpenshift = CloudPlatformKey.String("aws_openshift")
	// Azure Virtual Machines
	// Stability: experimental
	CloudPlatformAzureVm = CloudPlatformKey.String("azure_vm")
	// Azure Container Apps
	// Stability: experimental
	CloudPlatformAzureContainerApps = CloudPlatformKey.String("azure_container_apps")
	// Azure Container Instances
	// Stability: experimental
	CloudPlatformAzureContainerInstances = CloudPlatformKey.String("azure_container_instances")
	// Azure Kubernetes Service
	// Stability: experimental
	CloudPlatformAzureAks = CloudPlatformKey.String("azure_aks")
	// Azure Functions
	// Stability: experimental
	CloudPlatformAzureFunctions = CloudPlatformKey.String("azure_functions")
	// Azure App Service
	// Stability: experimental
	CloudPlatformAzureAppService = CloudPlatformKey.String("azure_app_service")
	// Azure Red Hat OpenShift
	// Stability: experimental
	CloudPlatformAzureOpenshift = CloudPlatformKey.String("azure_openshift")
	// Google Bare Metal Solution (BMS)
	// Stability: experimental
	CloudPlatformGcpBareMetalSolution = CloudPlatformKey.String("gcp_bare_metal_solution")
	// Google Cloud Compute Engine (GCE)
	// Stability: experimental
	CloudPlatformGcpComputeEngine = CloudPlatformKey.String("gcp_compute_engine")
	// Google Cloud Run
	// Stability: experimental
	CloudPlatformGcpCloudRun = CloudPlatformKey.String("gcp_cloud_run")
	// Google Cloud Kubernetes Engine (GKE)
	// Stability: experimental
	CloudPlatformGcpKubernetesEngine = CloudPlatformKey.String("gcp_kubernetes_engine")
	// Google Cloud Functions (GCF)
	// Stability: experimental
	CloudPlatformGcpCloudFunctions = CloudPlatformKey.String("gcp_cloud_functions")
	// Google Cloud App Engine (GAE)
	// Stability: experimental
	CloudPlatformGcpAppEngine = CloudPlatformKey.String("gcp_app_engine")
	// Red Hat OpenShift on Google Cloud
	// Stability: experimental
	CloudPlatformGcpOpenshift = CloudPlatformKey.String("gcp_openshift")
	// Red Hat OpenShift on IBM Cloud
	// Stability: experimental
	CloudPlatformIbmCloudOpenshift = CloudPlatformKey.String("ibm_cloud_openshift")
	// Tencent Cloud Cloud Virtual Machine (CVM)
	// Stability: experimental
	CloudPlatformTencentCloudCvm = CloudPlatformKey.String("tencent_cloud_cvm")
	// Tencent Cloud Elastic Kubernetes Service (EKS)
	// Stability: experimental
	CloudPlatformTencentCloudEks = CloudPlatformKey.String("tencent_cloud_eks")
	// Tencent Cloud Serverless Cloud Function (SCF)
	// Stability: experimental
	CloudPlatformTencentCloudScf = CloudPlatformKey.String("tencent_cloud_scf")
)

// Enum values for cloud.provider
var (

	// Alibaba Cloud
	// Stability: experimental
	CloudProviderAlibabaCloud = CloudProviderKey.String("alibaba_cloud")
	// Amazon Web Services
	// Stability: experimental
	CloudProviderAws = CloudProviderKey.String("aws")
	// Microsoft Azure
	// Stability: experimental
	CloudProviderAzure = CloudProviderKey.String("azure")
	// Google Cloud Platform
	// Stability: experimental
	CloudProviderGcp = CloudProviderKey.String("gcp")
	// Heroku Platform as a Service
	// Stability: experimental
	CloudProviderHeroku = CloudProviderKey.String("heroku")
	// IBM Cloud
	// Stability: experimental
	CloudProviderIbmCloud = CloudProviderKey.String("ibm_cloud")
	// Tencent Cloud
	// Stability: experimental
	CloudProviderTencentCloud = CloudProviderKey.String("tencent_cloud")
)

// Namespace: cloudevents
const (

	// The [event_id] uniquely identifies the event.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "123e4567-e89b-12d3-a456-426614174000",
	// "0001",
	//
	// [event_id]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#id
	CloudeventsEventIdKey = attribute.Key("cloudevents.event_id")
	// The [source] identifies the context in which an event happened.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "https://github.com/cloudevents",
	// "/cloudevents/spec/pull/123",
	// "my-service",
	//
	// [source]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#source-1
	CloudeventsEventSourceKey = attribute.Key("cloudevents.event_source")
	// The [version of the CloudEvents specification] which the event uses.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "1.0"
	//
	// [version of the CloudEvents specification]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#specversion
	CloudeventsEventSpecVersionKey = attribute.Key("cloudevents.event_spec_version")
	// The [subject] of the event in the context of the event producer (identified by source).
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "mynewfile.jpg"
	//
	// [subject]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#subject
	CloudeventsEventSubjectKey = attribute.Key("cloudevents.event_subject")
	// The [event_type] contains a value describing the type of event related to the originating occurrence.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "com.github.pull_request.opened",
	// "com.example.object.deleted.v2",
	//
	// [event_type]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#type
	CloudeventsEventTypeKey = attribute.Key("cloudevents.event_type")
)

// CloudeventsEventId returns an attribute KeyValue conforming to the "cloudevents.event_id"semantic conventions
// It represents the The [event_id] uniquely identifies the event
//
// [event_id]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#id
func CloudeventsEventId(val string) attribute.KeyValue {
	return CloudeventsEventIdKey.String(val)
}

// CloudeventsEventSource returns an attribute KeyValue conforming to the "cloudevents.event_source"semantic conventions
// It represents the The [source] identifies the context in which an event happened
//
// [source]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#source-1
func CloudeventsEventSource(val string) attribute.KeyValue {
	return CloudeventsEventSourceKey.String(val)
}

// CloudeventsEventSpecVersion returns an attribute KeyValue conforming to the "cloudevents.event_spec_version"semantic conventions
// It represents the The [version of the CloudEvents specification] which the event uses
//
// [version of the CloudEvents specification]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#specversion
func CloudeventsEventSpecVersion(val string) attribute.KeyValue {
	return CloudeventsEventSpecVersionKey.String(val)
}

// CloudeventsEventSubject returns an attribute KeyValue conforming to the "cloudevents.event_subject"semantic conventions
// It represents the The [subject] of the event in the context of the event producer (identified by source)
//
// [subject]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#subject
func CloudeventsEventSubject(val string) attribute.KeyValue {
	return CloudeventsEventSubjectKey.String(val)
}

// CloudeventsEventType returns an attribute KeyValue conforming to the "cloudevents.event_type"semantic conventions
// It represents the The [event_type] contains a value describing the type of event related to the originating occurrence
//
// [event_type]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#type
func CloudeventsEventType(val string) attribute.KeyValue {
	return CloudeventsEventTypeKey.String(val)
}

// Namespace: code
const (

	// The column number in `code.filepath` best representing the operation. It SHOULD point within the code unit named in `code.function`.
	//
	// Stability: Experimental
	// Type: int
	CodeColumnKey = attribute.Key("code.column")
	// The source code file name that identifies the code unit as uniquely as possible (preferably an absolute file path).
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "/usr/local/MyApplication/content_root/app/index.php"
	CodeFilepathKey = attribute.Key("code.filepath")
	// The method or function name, or equivalent (usually rightmost part of the code unit's name).
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "serveRequest"
	CodeFunctionKey = attribute.Key("code.function")
	// The line number in `code.filepath` best representing the operation. It SHOULD point within the code unit named in `code.function`.
	//
	// Stability: Experimental
	// Type: int
	CodeLinenoKey = attribute.Key("code.lineno")
	// The "namespace" within which `code.function` is defined. Usually the qualified class or module name, such that `code.namespace` + some separator + `code.function` form a unique identifier for the code unit.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "com.example.MyHttpService"
	CodeNamespaceKey = attribute.Key("code.namespace")
	// A stacktrace as a string in the natural representation for the language runtime. The representation is to be determined and documented by each language SIG.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "at com.example.GenerateTrace.methodB(GenerateTrace.java:13)\n at com.example.GenerateTrace.methodA(GenerateTrace.java:9)\n at com.example.GenerateTrace.main(GenerateTrace.java:5)\n"
	CodeStacktraceKey = attribute.Key("code.stacktrace")
)

// CodeColumn returns an attribute KeyValue conforming to the "code.column"semantic conventions
// It represents the The column number in `code.filepath` best representing the operation. It SHOULD point within the code unit named in `code.function`
func CodeColumn(val int) attribute.KeyValue {
	return CodeColumnKey.Int(val)
}

// CodeFilepath returns an attribute KeyValue conforming to the "code.filepath"semantic conventions
// It represents the The source code file name that identifies the code unit as uniquely as possible (preferably an absolute file path)
func CodeFilepath(val string) attribute.KeyValue {
	return CodeFilepathKey.String(val)
}

// CodeFunction returns an attribute KeyValue conforming to the "code.function"semantic conventions
// It represents the The method or function name, or equivalent (usually rightmost part of the code unit's name)
func CodeFunction(val string) attribute.KeyValue {
	return CodeFunctionKey.String(val)
}

// CodeLineno returns an attribute KeyValue conforming to the "code.lineno"semantic conventions
// It represents the The line number in `code.filepath` best representing the operation. It SHOULD point within the code unit named in `code.function`
func CodeLineno(val int) attribute.KeyValue {
	return CodeLinenoKey.Int(val)
}

// CodeNamespace returns an attribute KeyValue conforming to the "code.namespace"semantic conventions
// It represents the The "namespace" within which `code.function` is defined. Usually the qualified class or module name, such that `code.namespace` + some separator + `code.function` form a unique identifier for the code unit
func CodeNamespace(val string) attribute.KeyValue {
	return CodeNamespaceKey.String(val)
}

// CodeStacktrace returns an attribute KeyValue conforming to the "code.stacktrace"semantic conventions
// It represents the A stacktrace as a string in the natural representation for the language runtime. The representation is to be determined and documented by each language SIG
func CodeStacktrace(val string) attribute.KeyValue {
	return CodeStacktraceKey.String(val)
}

// Namespace: container
const (

	// The command used to run the container (i.e. the command name).
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "otelcontribcol",
	//
	// Note: If using embedded credentials or sensitive data, it is recommended to remove them to prevent potential leakage
	ContainerCommandKey = attribute.Key("container.command")
	// All the command arguments (including the command/executable itself) run by the container. [2]
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "otelcontribcol",
	// "--config",
	// "config.yaml",
	// ],
	ContainerCommandArgsKey = attribute.Key("container.command_args")
	// The full command run by the container as a single string representing the full command. [2]
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "otelcontribcol --config config.yaml",
	ContainerCommandLineKey = attribute.Key("container.command_line")
	// Deprecated, use `cpu.mode` instead.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `cpu.mode`
	//
	// Examples:
	// "user",
	// "kernel",
	ContainerCpuStateKey = attribute.Key("container.cpu.state")
	// Container ID. Usually a UUID, as for example used to [identify Docker containers]. The UUID might be abbreviated.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "a3bf90e006b2",
	//
	// [identify Docker containers]: https://docs.docker.com/engine/containers/run/#container-identification
	ContainerIdKey = attribute.Key("container.id")
	// Runtime specific image identifier. Usually a hash algorithm followed by a UUID.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "sha256:19c92d0a00d1b66d897bceaa7319bee0dd38a10a851c60bcec9474aa3f01e50f",
	//
	// Note: Docker defines a sha256 of the image id; `container.image.id` corresponds to the `Image` field from the Docker container inspect [API] endpoint.
	// K8s defines a link to the container registry repository with digest `"imageID": "registry.azurecr.io /namespace/service/dockerfile@sha256:bdeabd40c3a8a492eaf9e8e44d0ebbb84bac7ee25ac0cf8a7159d25f62555625"`.
	// The ID is assigned by the container runtime and can vary in different environments. Consider using `oci.manifest.digest` if it is important to identify the same image in different environments/runtimes
	//
	// [API]: https://docs.docker.com/engine/api/v1.43/#tag/Container/operation/ContainerInspect
	ContainerImageIdKey = attribute.Key("container.image.id")
	// Name of the image the container was built on.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "gcr.io/opentelemetry/operator",
	ContainerImageNameKey = attribute.Key("container.image.name")
	// Repo digests of the container image as provided by the container runtime.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "example@sha256:afcc7f1ac1b49db317a7196c902e61c6c3c4607d63599ee1a82d702d249a0ccb",
	// "internal.registry.example.com:5000/example@sha256:b69959407d21e8a062e0416bf13405bb2b71ed7a84dde4158ebafacfa06f5578",
	// ],
	//
	// Note: [Docker] and [CRI] report those under the `RepoDigests` field
	//
	// [Docker]: https://docs.docker.com/engine/api/v1.43/#tag/Image/operation/ImageInspect
	// [CRI]: https://github.com/kubernetes/cri-api/blob/c75ef5b473bbe2d0a4fc92f82235efd665ea8e9f/pkg/apis/runtime/v1/api.proto#L1237-L1238
	ContainerImageRepoDigestsKey = attribute.Key("container.image.repo_digests")
	// Container image tags. An example can be found in [Docker Image Inspect]. Should be only the `<tag>` section of the full name for example from `registry.example.com/my-org/my-image:<tag>`.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "v1.27.1",
	// "3.5.7-0",
	// ],
	//
	// [Docker Image Inspect]: https://docs.docker.com/engine/api/v1.43/#tag/Image/operation/ImageInspect
	ContainerImageTagsKey = attribute.Key("container.image.tags")
	// Container name used by container runtime.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry-autoconf",
	ContainerNameKey = attribute.Key("container.name")
	// The container runtime managing this container.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "docker",
	// "containerd",
	// "rkt",
	ContainerRuntimeKey = attribute.Key("container.runtime")
)

// ContainerCommand returns an attribute KeyValue conforming to the "container.command"semantic conventions
// It represents the The command used to run the container (i.e. the command name)
func ContainerCommand(val string) attribute.KeyValue {
	return ContainerCommandKey.String(val)
}

// ContainerCommandArgs returns an attribute KeyValue conforming to the "container.command_args"semantic conventions
// It represents the All the command arguments (including the command/executable itself) run by the container. [2]
func ContainerCommandArgs(val []string) attribute.KeyValue {
	return ContainerCommandArgsKey.StringSlice(val)
}

// ContainerCommandLine returns an attribute KeyValue conforming to the "container.command_line"semantic conventions
// It represents the The full command run by the container as a single string representing the full command. [2]
func ContainerCommandLine(val string) attribute.KeyValue {
	return ContainerCommandLineKey.String(val)
}

// ContainerCpuState returns an attribute KeyValue conforming to the "container.cpu.state"semantic conventions
// It represents the Deprecated, use `cpu.mode` instead
func ContainerCpuState(val string) attribute.KeyValue {
	return ContainerCpuStateKey.String(val)
}

// ContainerId returns an attribute KeyValue conforming to the "container.id"semantic conventions
// It represents the Container ID. Usually a UUID, as for example used to [identify Docker containers]. The UUID might be abbreviated
//
// [identify Docker containers]: https://docs.docker.com/engine/containers/run/#container-identification
func ContainerId(val string) attribute.KeyValue {
	return ContainerIdKey.String(val)
}

// ContainerImageId returns an attribute KeyValue conforming to the "container.image.id"semantic conventions
// It represents the Runtime specific image identifier. Usually a hash algorithm followed by a UUID
func ContainerImageId(val string) attribute.KeyValue {
	return ContainerImageIdKey.String(val)
}

// ContainerImageName returns an attribute KeyValue conforming to the "container.image.name"semantic conventions
// It represents the Name of the image the container was built on
func ContainerImageName(val string) attribute.KeyValue {
	return ContainerImageNameKey.String(val)
}

// ContainerImageRepoDigests returns an attribute KeyValue conforming to the "container.image.repo_digests"semantic conventions
// It represents the Repo digests of the container image as provided by the container runtime
func ContainerImageRepoDigests(val []string) attribute.KeyValue {
	return ContainerImageRepoDigestsKey.StringSlice(val)
}

// ContainerImageTags returns an attribute KeyValue conforming to the "container.image.tags"semantic conventions
// It represents the Container image tags. An example can be found in [Docker Image Inspect]. Should be only the `<tag>` section of the full name for example from `registry.example.com/my-org/my-image:<tag>`
//
// [Docker Image Inspect]: https://docs.docker.com/engine/api/v1.43/#tag/Image/operation/ImageInspect
func ContainerImageTags(val []string) attribute.KeyValue {
	return ContainerImageTagsKey.StringSlice(val)
}

// ContainerName returns an attribute KeyValue conforming to the "container.name"semantic conventions
// It represents the Container name used by container runtime
func ContainerName(val string) attribute.KeyValue {
	return ContainerNameKey.String(val)
}

// ContainerRuntime returns an attribute KeyValue conforming to the "container.runtime"semantic conventions
// It represents the The container runtime managing this container
func ContainerRuntime(val string) attribute.KeyValue {
	return ContainerRuntimeKey.String(val)
}

// Enum values for container.cpu.state
var (

	// When tasks of the cgroup are in user mode (Linux). When all container processes are in user mode (Windows).
	// Stability: experimental
	ContainerCpuStateUser = ContainerCpuStateKey.String("user")
	// When CPU is used by the system (host OS)
	// Stability: experimental
	ContainerCpuStateSystem = ContainerCpuStateKey.String("system")
	// When tasks of the cgroup are in kernel mode (Linux). When all container processes are in kernel mode (Windows).
	// Stability: experimental
	ContainerCpuStateKernel = ContainerCpuStateKey.String("kernel")
)

// Namespace: cpu
const (

	// The mode of the CPU
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "user",
	// "system",
	CpuModeKey = attribute.Key("cpu.mode")
)

// CpuMode returns an attribute KeyValue conforming to the "cpu.mode"semantic conventions
// It represents the The mode of the CPU
func CpuMode(val string) attribute.KeyValue {
	return CpuModeKey.String(val)
}

// Enum values for cpu.mode
var (

	// user
	// Stability: experimental
	CpuModeUser = CpuModeKey.String("user")
	// system
	// Stability: experimental
	CpuModeSystem = CpuModeKey.String("system")
	// nice
	// Stability: experimental
	CpuModeNice = CpuModeKey.String("nice")
	// idle
	// Stability: experimental
	CpuModeIdle = CpuModeKey.String("idle")
	// iowait
	// Stability: experimental
	CpuModeIowait = CpuModeKey.String("iowait")
	// interrupt
	// Stability: experimental
	CpuModeInterrupt = CpuModeKey.String("interrupt")
	// steal
	// Stability: experimental
	CpuModeSteal = CpuModeKey.String("steal")
	// kernel
	// Stability: experimental
	CpuModeKernel = CpuModeKey.String("kernel")
)

// Namespace: db
const (

	// The consistency level of the query. Based on consistency values from [CQL].
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	//
	// [CQL]: https://docs.datastax.com/en/cassandra-oss/3.0/cassandra/dml/dmlConfigConsistency.html
	DbCassandraConsistencyLevelKey = attribute.Key("db.cassandra.consistency_level")
	// The data center of the coordinating node for a query.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "us-west-2"
	DbCassandraCoordinatorDcKey = attribute.Key("db.cassandra.coordinator.dc")
	// The ID of the coordinating node for a query.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "be13faa2-8574-4d71-926d-27f16cf8a7af"
	DbCassandraCoordinatorIdKey = attribute.Key("db.cassandra.coordinator.id")
	// Whether or not the query is idempotent.
	//
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	DbCassandraIdempotenceKey = attribute.Key("db.cassandra.idempotence")
	// The fetch size used for paging, i.e. how many rows will be returned at once.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 5000,
	DbCassandraPageSizeKey = attribute.Key("db.cassandra.page_size")
	// The number of times a query was speculatively executed. Not set or `0` if the query was not executed speculatively.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 0,
	// 2,
	DbCassandraSpeculativeExecutionCountKey = attribute.Key("db.cassandra.speculative_execution_count")
	// Deprecated, use `db.collection.name` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.collection.name`.
	//
	// Examples: "mytable"
	DbCassandraTableKey = attribute.Key("db.cassandra.table")
	// The name of the connection pool; unique within the instrumented application. In case the connection pool implementation doesn't provide a name, instrumentation SHOULD use a combination of parameters that would make the name unique, for example, combining attributes `server.address`, `server.port`, and `db.namespace`, formatted as `server.address:server.port/db.namespace`. Instrumentations that generate connection pool name following different patterns SHOULD document it.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "myDataSource",
	DbClientConnectionPoolNameKey = attribute.Key("db.client.connection.pool.name")
	// The state of a connection in the pool
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "idle",
	DbClientConnectionStateKey = attribute.Key("db.client.connection.state")
	// Deprecated, use `db.client.connection.pool.name` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.client.connection.pool.name`.
	//
	// Examples:
	// "myDataSource",
	DbClientConnectionsPoolNameKey = attribute.Key("db.client.connections.pool.name")
	// Deprecated, use `db.client.connection.state` instead.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `db.client.connection.state`.
	//
	// Examples:
	// "idle",
	DbClientConnectionsStateKey = attribute.Key("db.client.connections.state")
	// The name of a collection (table, container) within the database.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "public.users",
	// "customers",
	//
	// Note: It is RECOMMENDED to capture the value as provided by the application without attempting to do any case normalization.
	// If the collection name is parsed from the query text, it SHOULD be the first collection name found in the query and it SHOULD match the value provided in the query text including any schema and database name prefix.
	// For batch operations, if the individual operations are known to have the same collection name then that collection name SHOULD be used, otherwise `db.collection.name` SHOULD NOT be captured
	DbCollectionNameKey = attribute.Key("db.collection.name")
	// Deprecated, use `server.address`, `server.port` attributes instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `server.address` and `server.port`.
	//
	// Examples: "Server=(localdb)\v11.0;Integrated Security=true;"
	DbConnectionStringKey = attribute.Key("db.connection_string")
	// Unique Cosmos client instance id.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "3ba4827d-4422-483f-b59f-85b74211c11d"
	DbCosmosdbClientIdKey = attribute.Key("db.cosmosdb.client_id")
	// Cosmos client connection mode.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	DbCosmosdbConnectionModeKey = attribute.Key("db.cosmosdb.connection_mode")
	// Deprecated, use `db.collection.name` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.collection.name`.
	//
	// Examples: "mytable"
	DbCosmosdbContainerKey = attribute.Key("db.cosmosdb.container")
	// Cosmos DB Operation Type.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	DbCosmosdbOperationTypeKey = attribute.Key("db.cosmosdb.operation_type")
	// RU consumed for that operation
	// Stability: Experimental
	// Type: double
	//
	// Examples:
	// 46.18,
	// 1.0,
	DbCosmosdbRequestChargeKey = attribute.Key("db.cosmosdb.request_charge")
	// Request payload size in bytes
	// Stability: Experimental
	// Type: int
	//
	// Examples: undefined
	DbCosmosdbRequestContentLengthKey = attribute.Key("db.cosmosdb.request_content_length")
	// Cosmos DB status code.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 200,
	// 201,
	DbCosmosdbStatusCodeKey = attribute.Key("db.cosmosdb.status_code")
	// Cosmos DB sub status code.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 1000,
	// 1002,
	DbCosmosdbSubStatusCodeKey = attribute.Key("db.cosmosdb.sub_status_code")
	// Deprecated, use `db.namespace` instead.
	//
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.namespace`.
	//
	// Examples:
	// "e9106fc68e3044f0b1475b04bf4ffd5f",
	DbElasticsearchClusterNameKey = attribute.Key("db.elasticsearch.cluster.name")
	// Represents the human-readable identifier of the node/instance to which a request was routed.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "instance-0000000001",
	DbElasticsearchNodeNameKey = attribute.Key("db.elasticsearch.node.name")
	// Deprecated, no general replacement at this time. For Elasticsearch, use `db.elasticsearch.node.name` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Deprecated, no general replacement at this time. For Elasticsearch, use `db.elasticsearch.node.name` instead.
	//
	// Examples: "mysql-e26b99z.example.com"
	DbInstanceIdKey = attribute.Key("db.instance.id")
	// Removed, no replacement at this time.
	// Stability: Experimental
	// Type: string
	// Deprecated: Removed as not used.
	//
	// Examples:
	// "org.postgresql.Driver",
	// "com.microsoft.sqlserver.jdbc.SQLServerDriver",
	DbJdbcDriverClassnameKey = attribute.Key("db.jdbc.driver_classname")
	// Deprecated, use `db.collection.name` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.collection.name`.
	//
	// Examples: "mytable"
	DbMongodbCollectionKey = attribute.Key("db.mongodb.collection")
	// Deprecated, SQL Server instance is now populated as a part of `db.namespace` attribute.
	// Stability: Experimental
	// Type: string
	// Deprecated: Deprecated, no replacement at this time.
	//
	// Examples: "MSSQLSERVER"
	DbMssqlInstanceNameKey = attribute.Key("db.mssql.instance_name")
	// Deprecated, use `db.namespace` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.namespace`.
	//
	// Examples:
	// "customers",
	// "main",
	DbNameKey = attribute.Key("db.name")
	// The name of the database, fully qualified within the server address and port.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "customers",
	// "test.users",
	//
	// Note: If a database system has multiple namespace components, they SHOULD be concatenated (potentially using database system specific conventions) from most general to most specific namespace component, and more specific namespaces SHOULD NOT be captured without the more general namespaces, to ensure that "startswith" queries for the more general namespaces will be valid.
	// Semantic conventions for individual database systems SHOULD document what `db.namespace` means in the context of that system.
	// It is RECOMMENDED to capture the value as provided by the application without attempting to do any case normalization
	DbNamespaceKey = attribute.Key("db.namespace")
	// Deprecated, use `db.operation.name` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.operation.name`.
	//
	// Examples:
	// "findAndModify",
	// "HMSET",
	// "SELECT",
	DbOperationKey = attribute.Key("db.operation")
	// The number of queries included in a batch operation.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 2,
	// 3,
	// 4,
	//
	// Note: Operations are only considered batches when they contain two or more operations, and so `db.operation.batch.size` SHOULD never be `1`
	DbOperationBatchSizeKey = attribute.Key("db.operation.batch.size")
	// The name of the operation or command being executed.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "findAndModify",
	// "HMSET",
	// "SELECT",
	//
	// Note: It is RECOMMENDED to capture the value as provided by the application without attempting to do any case normalization.
	// If the operation name is parsed from the query text, it SHOULD be the first operation name found in the query.
	// For batch operations, if the individual operations are known to have the same operation name then that operation name SHOULD be used prepended by `BATCH `, otherwise `db.operation.name` SHOULD be `BATCH` or some other database system specific term if more applicable
	DbOperationNameKey = attribute.Key("db.operation.name")
	// The database query being executed.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "SELECT * FROM wuser_table where username = ?",
	// "SET mykey "WuValue"",
	//
	// Note: For sanitization see [Sanitization of `db.query.text`].
	// For batch operations, if the individual operations are known to have the same query text then that query text SHOULD be used, otherwise all of the individual query texts SHOULD be concatenated with separator `; ` or some other database system specific separator if more applicable.
	// Even though parameterized query text can potentially have sensitive data, by using a parameterized query the user is giving a strong signal that any sensitive data will be passed as parameter values, and the benefit to observability of capturing the static part of the query text by default outweighs the risk
	//
	// [Sanitization of `db.query.text`]: ../../docs/database/database-spans.md#sanitization-of-dbquerytext
	DbQueryTextKey = attribute.Key("db.query.text")
	// Deprecated, use `db.namespace` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `db.namespace`.
	//
	// Examples:
	// 0,
	// 1,
	// 15,
	DbRedisDatabaseIndexKey = attribute.Key("db.redis.database_index")
	// Deprecated, use `db.collection.name` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.collection.name`.
	//
	// Examples: "mytable"
	DbSqlTableKey = attribute.Key("db.sql.table")
	// The database statement being executed.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.query.text`.
	//
	// Examples:
	// "SELECT * FROM wuser_table",
	// "SET mykey "WuValue"",
	DbStatementKey = attribute.Key("db.statement")
	// The database management system (DBMS) product as identified by the client instrumentation.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: The actual DBMS may differ from the one identified by the client. For example, when using PostgreSQL client libraries to connect to a CockroachDB, the `db.system` is set to `postgresql` based on the instrumentation's best knowledge
	DbSystemKey = attribute.Key("db.system")
	// Deprecated, no replacement at this time.
	// Stability: Experimental
	// Type: string
	// Deprecated: No replacement at this time.
	//
	// Examples:
	// "readonly_user",
	// "reporting_user",
	DbUserKey = attribute.Key("db.user")
)

// DbCassandraConsistencyLevel returns an attribute KeyValue conforming to the "db.cassandra.consistency_level"semantic conventions
// It represents the The consistency level of the query. Based on consistency values from [CQL]
//
// [CQL]: https://docs.datastax.com/en/cassandra-oss/3.0/cassandra/dml/dmlConfigConsistency.html
func DbCassandraConsistencyLevel(val string) attribute.KeyValue {
	return DbCassandraConsistencyLevelKey.String(val)
}

// DbCassandraCoordinatorDc returns an attribute KeyValue conforming to the "db.cassandra.coordinator.dc"semantic conventions
// It represents the The data center of the coordinating node for a query
func DbCassandraCoordinatorDc(val string) attribute.KeyValue {
	return DbCassandraCoordinatorDcKey.String(val)
}

// DbCassandraCoordinatorId returns an attribute KeyValue conforming to the "db.cassandra.coordinator.id"semantic conventions
// It represents the The ID of the coordinating node for a query
func DbCassandraCoordinatorId(val string) attribute.KeyValue {
	return DbCassandraCoordinatorIdKey.String(val)
}

// DbCassandraIdempotence returns an attribute KeyValue conforming to the "db.cassandra.idempotence"semantic conventions
// It represents the Whether or not the query is idempotent
func DbCassandraIdempotence(val bool) attribute.KeyValue {
	return DbCassandraIdempotenceKey.Bool(val)
}

// DbCassandraPageSize returns an attribute KeyValue conforming to the "db.cassandra.page_size"semantic conventions
// It represents the The fetch size used for paging, i.e. how many rows will be returned at once
func DbCassandraPageSize(val int) attribute.KeyValue {
	return DbCassandraPageSizeKey.Int(val)
}

// DbCassandraSpeculativeExecutionCount returns an attribute KeyValue conforming to the "db.cassandra.speculative_execution_count"semantic conventions
// It represents the The number of times a query was speculatively executed. Not set or `0` if the query was not executed speculatively
func DbCassandraSpeculativeExecutionCount(val int) attribute.KeyValue {
	return DbCassandraSpeculativeExecutionCountKey.Int(val)
}

// DbCassandraTable returns an attribute KeyValue conforming to the "db.cassandra.table"semantic conventions
// It represents the Deprecated, use `db.collection.name` instead
func DbCassandraTable(val string) attribute.KeyValue {
	return DbCassandraTableKey.String(val)
}

// DbClientConnectionPoolName returns an attribute KeyValue conforming to the "db.client.connection.pool.name"semantic conventions
// It represents the The name of the connection pool; unique within the instrumented application. In case the connection pool implementation doesn't provide a name, instrumentation SHOULD use a combination of parameters that would make the name unique, for example, combining attributes `server.address`, `server.port`, and `db.namespace`, formatted as `server.address:server.port/db.namespace`. Instrumentations that generate connection pool name following different patterns SHOULD document it
func DbClientConnectionPoolName(val string) attribute.KeyValue {
	return DbClientConnectionPoolNameKey.String(val)
}

// DbClientConnectionState returns an attribute KeyValue conforming to the "db.client.connection.state"semantic conventions
// It represents the The state of a connection in the pool
func DbClientConnectionState(val string) attribute.KeyValue {
	return DbClientConnectionStateKey.String(val)
}

// DbClientConnectionsPoolName returns an attribute KeyValue conforming to the "db.client.connections.pool.name"semantic conventions
// It represents the Deprecated, use `db.client.connection.pool.name` instead
func DbClientConnectionsPoolName(val string) attribute.KeyValue {
	return DbClientConnectionsPoolNameKey.String(val)
}

// DbClientConnectionsState returns an attribute KeyValue conforming to the "db.client.connections.state"semantic conventions
// It represents the Deprecated, use `db.client.connection.state` instead
func DbClientConnectionsState(val string) attribute.KeyValue {
	return DbClientConnectionsStateKey.String(val)
}

// DbCollectionName returns an attribute KeyValue conforming to the "db.collection.name"semantic conventions
// It represents the The name of a collection (table, container) within the database
func DbCollectionName(val string) attribute.KeyValue {
	return DbCollectionNameKey.String(val)
}

// DbConnectionString returns an attribute KeyValue conforming to the "db.connection_string"semantic conventions
// It represents the Deprecated, use `server.address`, `server.port` attributes instead
func DbConnectionString(val string) attribute.KeyValue {
	return DbConnectionStringKey.String(val)
}

// DbCosmosdbClientId returns an attribute KeyValue conforming to the "db.cosmosdb.client_id"semantic conventions
// It represents the Unique Cosmos client instance id
func DbCosmosdbClientId(val string) attribute.KeyValue {
	return DbCosmosdbClientIdKey.String(val)
}

// DbCosmosdbConnectionMode returns an attribute KeyValue conforming to the "db.cosmosdb.connection_mode"semantic conventions
// It represents the Cosmos client connection mode
func DbCosmosdbConnectionMode(val string) attribute.KeyValue {
	return DbCosmosdbConnectionModeKey.String(val)
}

// DbCosmosdbContainer returns an attribute KeyValue conforming to the "db.cosmosdb.container"semantic conventions
// It represents the Deprecated, use `db.collection.name` instead
func DbCosmosdbContainer(val string) attribute.KeyValue {
	return DbCosmosdbContainerKey.String(val)
}

// DbCosmosdbOperationType returns an attribute KeyValue conforming to the "db.cosmosdb.operation_type"semantic conventions
// It represents the Cosmos DB Operation Type
func DbCosmosdbOperationType(val string) attribute.KeyValue {
	return DbCosmosdbOperationTypeKey.String(val)
}

// DbCosmosdbRequestCharge returns an attribute KeyValue conforming to the "db.cosmosdb.request_charge"semantic conventions
// It represents the RU consumed for that operation
func DbCosmosdbRequestCharge(val float64) attribute.KeyValue {
	return DbCosmosdbRequestChargeKey.Float64(val)
}

// DbCosmosdbRequestContentLength returns an attribute KeyValue conforming to the "db.cosmosdb.request_content_length"semantic conventions
// It represents the Request payload size in bytes
func DbCosmosdbRequestContentLength(val int) attribute.KeyValue {
	return DbCosmosdbRequestContentLengthKey.Int(val)
}

// DbCosmosdbStatusCode returns an attribute KeyValue conforming to the "db.cosmosdb.status_code"semantic conventions
// It represents the Cosmos DB status code
func DbCosmosdbStatusCode(val int) attribute.KeyValue {
	return DbCosmosdbStatusCodeKey.Int(val)
}

// DbCosmosdbSubStatusCode returns an attribute KeyValue conforming to the "db.cosmosdb.sub_status_code"semantic conventions
// It represents the Cosmos DB sub status code
func DbCosmosdbSubStatusCode(val int) attribute.KeyValue {
	return DbCosmosdbSubStatusCodeKey.Int(val)
}

// DbElasticsearchClusterName returns an attribute KeyValue conforming to the "db.elasticsearch.cluster.name"semantic conventions
// It represents the Deprecated, use `db.namespace` instead
func DbElasticsearchClusterName(val string) attribute.KeyValue {
	return DbElasticsearchClusterNameKey.String(val)
}

// DbElasticsearchNodeName returns an attribute KeyValue conforming to the "db.elasticsearch.node.name"semantic conventions
// It represents the Represents the human-readable identifier of the node/instance to which a request was routed
func DbElasticsearchNodeName(val string) attribute.KeyValue {
	return DbElasticsearchNodeNameKey.String(val)
}

// DbInstanceId returns an attribute KeyValue conforming to the "db.instance.id"semantic conventions
// It represents the Deprecated, no general replacement at this time. For Elasticsearch, use `db.elasticsearch.node.name` instead
func DbInstanceId(val string) attribute.KeyValue {
	return DbInstanceIdKey.String(val)
}

// DbJdbcDriverClassname returns an attribute KeyValue conforming to the "db.jdbc.driver_classname"semantic conventions
// It represents the Removed, no replacement at this time
func DbJdbcDriverClassname(val string) attribute.KeyValue {
	return DbJdbcDriverClassnameKey.String(val)
}

// DbMongodbCollection returns an attribute KeyValue conforming to the "db.mongodb.collection"semantic conventions
// It represents the Deprecated, use `db.collection.name` instead
func DbMongodbCollection(val string) attribute.KeyValue {
	return DbMongodbCollectionKey.String(val)
}

// DbMssqlInstanceName returns an attribute KeyValue conforming to the "db.mssql.instance_name"semantic conventions
// It represents the Deprecated, SQL Server instance is now populated as a part of `db.namespace` attribute
func DbMssqlInstanceName(val string) attribute.KeyValue {
	return DbMssqlInstanceNameKey.String(val)
}

// DbName returns an attribute KeyValue conforming to the "db.name"semantic conventions
// It represents the Deprecated, use `db.namespace` instead
func DbName(val string) attribute.KeyValue {
	return DbNameKey.String(val)
}

// DbNamespace returns an attribute KeyValue conforming to the "db.namespace"semantic conventions
// It represents the The name of the database, fully qualified within the server address and port
func DbNamespace(val string) attribute.KeyValue {
	return DbNamespaceKey.String(val)
}

// DbOperation returns an attribute KeyValue conforming to the "db.operation"semantic conventions
// It represents the Deprecated, use `db.operation.name` instead
func DbOperation(val string) attribute.KeyValue {
	return DbOperationKey.String(val)
}

// DbOperationBatchSize returns an attribute KeyValue conforming to the "db.operation.batch.size"semantic conventions
// It represents the The number of queries included in a batch operation
func DbOperationBatchSize(val int) attribute.KeyValue {
	return DbOperationBatchSizeKey.Int(val)
}

// DbOperationName returns an attribute KeyValue conforming to the "db.operation.name"semantic conventions
// It represents the The name of the operation or command being executed
func DbOperationName(val string) attribute.KeyValue {
	return DbOperationNameKey.String(val)
}

// DbQueryText returns an attribute KeyValue conforming to the "db.query.text"semantic conventions
// It represents the The database query being executed
func DbQueryText(val string) attribute.KeyValue {
	return DbQueryTextKey.String(val)
}

// DbRedisDatabaseIndex returns an attribute KeyValue conforming to the "db.redis.database_index"semantic conventions
// It represents the Deprecated, use `db.namespace` instead
func DbRedisDatabaseIndex(val int) attribute.KeyValue {
	return DbRedisDatabaseIndexKey.Int(val)
}

// DbSqlTable returns an attribute KeyValue conforming to the "db.sql.table"semantic conventions
// It represents the Deprecated, use `db.collection.name` instead
func DbSqlTable(val string) attribute.KeyValue {
	return DbSqlTableKey.String(val)
}

// DbStatement returns an attribute KeyValue conforming to the "db.statement"semantic conventions
// It represents the The database statement being executed
func DbStatement(val string) attribute.KeyValue {
	return DbStatementKey.String(val)
}

// DbSystem returns an attribute KeyValue conforming to the "db.system"semantic conventions
// It represents the The database management system (DBMS) product as identified by the client instrumentation
func DbSystem(val string) attribute.KeyValue {
	return DbSystemKey.String(val)
}

// DbUser returns an attribute KeyValue conforming to the "db.user"semantic conventions
// It represents the Deprecated, no replacement at this time
func DbUser(val string) attribute.KeyValue {
	return DbUserKey.String(val)
}

// Enum values for db.cassandra.consistency_level
var (

	// all
	// Stability: experimental
	DbCassandraConsistencyLevelAll = DbCassandraConsistencyLevelKey.String("all")
	// each_quorum
	// Stability: experimental
	DbCassandraConsistencyLevelEachQuorum = DbCassandraConsistencyLevelKey.String("each_quorum")
	// quorum
	// Stability: experimental
	DbCassandraConsistencyLevelQuorum = DbCassandraConsistencyLevelKey.String("quorum")
	// local_quorum
	// Stability: experimental
	DbCassandraConsistencyLevelLocalQuorum = DbCassandraConsistencyLevelKey.String("local_quorum")
	// one
	// Stability: experimental
	DbCassandraConsistencyLevelOne = DbCassandraConsistencyLevelKey.String("one")
	// two
	// Stability: experimental
	DbCassandraConsistencyLevelTwo = DbCassandraConsistencyLevelKey.String("two")
	// three
	// Stability: experimental
	DbCassandraConsistencyLevelThree = DbCassandraConsistencyLevelKey.String("three")
	// local_one
	// Stability: experimental
	DbCassandraConsistencyLevelLocalOne = DbCassandraConsistencyLevelKey.String("local_one")
	// any
	// Stability: experimental
	DbCassandraConsistencyLevelAny = DbCassandraConsistencyLevelKey.String("any")
	// serial
	// Stability: experimental
	DbCassandraConsistencyLevelSerial = DbCassandraConsistencyLevelKey.String("serial")
	// local_serial
	// Stability: experimental
	DbCassandraConsistencyLevelLocalSerial = DbCassandraConsistencyLevelKey.String("local_serial")
)

// Enum values for db.client.connection.state
var (

	// idle
	// Stability: experimental
	DbClientConnectionStateIdle = DbClientConnectionStateKey.String("idle")
	// used
	// Stability: experimental
	DbClientConnectionStateUsed = DbClientConnectionStateKey.String("used")
)

// Enum values for db.client.connections.state
var (

	// idle
	// Stability: experimental
	DbClientConnectionsStateIdle = DbClientConnectionsStateKey.String("idle")
	// used
	// Stability: experimental
	DbClientConnectionsStateUsed = DbClientConnectionsStateKey.String("used")
)

// Enum values for db.cosmosdb.connection_mode
var (

	// Gateway (HTTP) connections mode
	// Stability: experimental
	DbCosmosdbConnectionModeGateway = DbCosmosdbConnectionModeKey.String("gateway")
	// Direct connection.
	// Stability: experimental
	DbCosmosdbConnectionModeDirect = DbCosmosdbConnectionModeKey.String("direct")
)

// Enum values for db.cosmosdb.operation_type
var (

	// batch
	// Stability: experimental
	DbCosmosdbOperationTypeBatch = DbCosmosdbOperationTypeKey.String("batch")
	// create
	// Stability: experimental
	DbCosmosdbOperationTypeCreate = DbCosmosdbOperationTypeKey.String("create")
	// delete
	// Stability: experimental
	DbCosmosdbOperationTypeDelete = DbCosmosdbOperationTypeKey.String("delete")
	// execute
	// Stability: experimental
	DbCosmosdbOperationTypeExecute = DbCosmosdbOperationTypeKey.String("execute")
	// execute_javascript
	// Stability: experimental
	DbCosmosdbOperationTypeExecuteJavascript = DbCosmosdbOperationTypeKey.String("execute_javascript")
	// invalid
	// Stability: experimental
	DbCosmosdbOperationTypeInvalid = DbCosmosdbOperationTypeKey.String("invalid")
	// head
	// Stability: experimental
	DbCosmosdbOperationTypeHead = DbCosmosdbOperationTypeKey.String("head")
	// head_feed
	// Stability: experimental
	DbCosmosdbOperationTypeHeadFeed = DbCosmosdbOperationTypeKey.String("head_feed")
	// patch
	// Stability: experimental
	DbCosmosdbOperationTypePatch = DbCosmosdbOperationTypeKey.String("patch")
	// query
	// Stability: experimental
	DbCosmosdbOperationTypeQuery = DbCosmosdbOperationTypeKey.String("query")
	// query_plan
	// Stability: experimental
	DbCosmosdbOperationTypeQueryPlan = DbCosmosdbOperationTypeKey.String("query_plan")
	// read
	// Stability: experimental
	DbCosmosdbOperationTypeRead = DbCosmosdbOperationTypeKey.String("read")
	// read_feed
	// Stability: experimental
	DbCosmosdbOperationTypeReadFeed = DbCosmosdbOperationTypeKey.String("read_feed")
	// replace
	// Stability: experimental
	DbCosmosdbOperationTypeReplace = DbCosmosdbOperationTypeKey.String("replace")
	// upsert
	// Stability: experimental
	DbCosmosdbOperationTypeUpsert = DbCosmosdbOperationTypeKey.String("upsert")
)

// Enum values for db.system
var (

	// Some other SQL database. Fallback only. See notes.
	// Stability: experimental
	DbSystemOtherSql = DbSystemKey.String("other_sql")
	// Adabas (Adaptable Database System)
	// Stability: experimental
	DbSystemAdabas = DbSystemKey.String("adabas")
	// Deprecated, use `intersystems_cache` instead.
	// Stability: experimental// Replaced by `intersystems_cache`
	DbSystemCache = DbSystemKey.String("cache")
	// InterSystems Caché
	// Stability: experimental
	DbSystemIntersystemsCache = DbSystemKey.String("intersystems_cache")
	// Apache Cassandra
	// Stability: experimental
	DbSystemCassandra = DbSystemKey.String("cassandra")
	// ClickHouse
	// Stability: experimental
	DbSystemClickhouse = DbSystemKey.String("clickhouse")
	// Deprecated, use `other_sql` instead.
	// Stability: experimental// Replaced by `other_sql`
	DbSystemCloudscape = DbSystemKey.String("cloudscape")
	// CockroachDB
	// Stability: experimental
	DbSystemCockroachdb = DbSystemKey.String("cockroachdb")
	// Deprecated, no replacement at this time.
	// Stability: experimental// Removed
	DbSystemColdfusion = DbSystemKey.String("coldfusion")
	// Microsoft Azure Cosmos DB
	// Stability: experimental
	DbSystemCosmosdb = DbSystemKey.String("cosmosdb")
	// Couchbase
	// Stability: experimental
	DbSystemCouchbase = DbSystemKey.String("couchbase")
	// CouchDB
	// Stability: experimental
	DbSystemCouchdb = DbSystemKey.String("couchdb")
	// IBM Db2
	// Stability: experimental
	DbSystemDb2 = DbSystemKey.String("db2")
	// Apache Derby
	// Stability: experimental
	DbSystemDerby = DbSystemKey.String("derby")
	// Amazon DynamoDB
	// Stability: experimental
	DbSystemDynamodb = DbSystemKey.String("dynamodb")
	// EnterpriseDB
	// Stability: experimental
	DbSystemEdb = DbSystemKey.String("edb")
	// Elasticsearch
	// Stability: experimental
	DbSystemElasticsearch = DbSystemKey.String("elasticsearch")
	// FileMaker
	// Stability: experimental
	DbSystemFilemaker = DbSystemKey.String("filemaker")
	// Firebird
	// Stability: experimental
	DbSystemFirebird = DbSystemKey.String("firebird")
	// Deprecated, use `other_sql` instead.
	// Stability: experimental// Replaced by `other_sql`
	DbSystemFirstsql = DbSystemKey.String("firstsql")
	// Apache Geode
	// Stability: experimental
	DbSystemGeode = DbSystemKey.String("geode")
	// H2
	// Stability: experimental
	DbSystemH2 = DbSystemKey.String("h2")
	// SAP HANA
	// Stability: experimental
	DbSystemHanadb = DbSystemKey.String("hanadb")
	// Apache HBase
	// Stability: experimental
	DbSystemHbase = DbSystemKey.String("hbase")
	// Apache Hive
	// Stability: experimental
	DbSystemHive = DbSystemKey.String("hive")
	// HyperSQL DataBase
	// Stability: experimental
	DbSystemHsqldb = DbSystemKey.String("hsqldb")
	// InfluxDB
	// Stability: experimental
	DbSystemInfluxdb = DbSystemKey.String("influxdb")
	// Informix
	// Stability: experimental
	DbSystemInformix = DbSystemKey.String("informix")
	// Ingres
	// Stability: experimental
	DbSystemIngres = DbSystemKey.String("ingres")
	// InstantDB
	// Stability: experimental
	DbSystemInstantdb = DbSystemKey.String("instantdb")
	// InterBase
	// Stability: experimental
	DbSystemInterbase = DbSystemKey.String("interbase")
	// MariaDB
	// Stability: experimental
	DbSystemMariadb = DbSystemKey.String("mariadb")
	// SAP MaxDB
	// Stability: experimental
	DbSystemMaxdb = DbSystemKey.String("maxdb")
	// Memcached
	// Stability: experimental
	DbSystemMemcached = DbSystemKey.String("memcached")
	// MongoDB
	// Stability: experimental
	DbSystemMongodb = DbSystemKey.String("mongodb")
	// Microsoft SQL Server
	// Stability: experimental
	DbSystemMssql = DbSystemKey.String("mssql")
	// Deprecated, Microsoft SQL Server Compact is discontinued.
	// Stability: experimental// Removed, use `other_sql` instead
	DbSystemMssqlcompact = DbSystemKey.String("mssqlcompact")
	// MySQL
	// Stability: experimental
	DbSystemMysql = DbSystemKey.String("mysql")
	// Neo4j
	// Stability: experimental
	DbSystemNeo4j = DbSystemKey.String("neo4j")
	// Netezza
	// Stability: experimental
	DbSystemNetezza = DbSystemKey.String("netezza")
	// OpenSearch
	// Stability: experimental
	DbSystemOpensearch = DbSystemKey.String("opensearch")
	// Oracle Database
	// Stability: experimental
	DbSystemOracle = DbSystemKey.String("oracle")
	// Pervasive PSQL
	// Stability: experimental
	DbSystemPervasive = DbSystemKey.String("pervasive")
	// PointBase
	// Stability: experimental
	DbSystemPointbase = DbSystemKey.String("pointbase")
	// PostgreSQL
	// Stability: experimental
	DbSystemPostgresql = DbSystemKey.String("postgresql")
	// Progress Database
	// Stability: experimental
	DbSystemProgress = DbSystemKey.String("progress")
	// Redis
	// Stability: experimental
	DbSystemRedis = DbSystemKey.String("redis")
	// Amazon Redshift
	// Stability: experimental
	DbSystemRedshift = DbSystemKey.String("redshift")
	// Cloud Spanner
	// Stability: experimental
	DbSystemSpanner = DbSystemKey.String("spanner")
	// SQLite
	// Stability: experimental
	DbSystemSqlite = DbSystemKey.String("sqlite")
	// Sybase
	// Stability: experimental
	DbSystemSybase = DbSystemKey.String("sybase")
	// Teradata
	// Stability: experimental
	DbSystemTeradata = DbSystemKey.String("teradata")
	// Trino
	// Stability: experimental
	DbSystemTrino = DbSystemKey.String("trino")
	// Vertica
	// Stability: experimental
	DbSystemVertica = DbSystemKey.String("vertica")
)

// Namespace: deployment
const (

	// 'Deprecated, use `deployment.environment.name` instead.'
	//
	// Stability: Experimental
	// Type: string
	// Deprecated: Deprecated, use `deployment.environment.name` instead.
	//
	// Examples:
	// "staging",
	// "production",
	DeploymentEnvironmentKey = attribute.Key("deployment.environment")
	// Name of the [deployment environment] (aka deployment tier).
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "staging",
	// "production",
	//
	// Note: `deployment.environment.name` does not affect the uniqueness constraints defined through
	// the `service.namespace`, `service.name` and `service.instance.id` resource attributes.
	// This implies that resources carrying the following attribute combinations MUST be
	// considered to be identifying the same service:
	//
	//   - `service.name=frontend`, `deployment.environment.name=production`
	//   - `service.name=frontend`, `deployment.environment.name=staging`
	// [deployment environment]: https://wikipedia.org/wiki/Deployment_environment
	DeploymentEnvironmentNameKey = attribute.Key("deployment.environment.name")
	// The id of the deployment.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "1208",
	DeploymentIdKey = attribute.Key("deployment.id")
	// The name of the deployment.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "deploy my app",
	// "deploy-frontend",
	DeploymentNameKey = attribute.Key("deployment.name")
	// The status of the deployment.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	DeploymentStatusKey = attribute.Key("deployment.status")
)

// DeploymentEnvironment returns an attribute KeyValue conforming to the "deployment.environment"semantic conventions
// It represents the 'Deprecated, use `deployment.environment.name` instead.'
func DeploymentEnvironment(val string) attribute.KeyValue {
	return DeploymentEnvironmentKey.String(val)
}

// DeploymentEnvironmentName returns an attribute KeyValue conforming to the "deployment.environment.name"semantic conventions
// It represents the Name of the [deployment environment] (aka deployment tier)
//
// [deployment environment]: https://wikipedia.org/wiki/Deployment_environment
func DeploymentEnvironmentName(val string) attribute.KeyValue {
	return DeploymentEnvironmentNameKey.String(val)
}

// DeploymentId returns an attribute KeyValue conforming to the "deployment.id"semantic conventions
// It represents the The id of the deployment
func DeploymentId(val string) attribute.KeyValue {
	return DeploymentIdKey.String(val)
}

// DeploymentName returns an attribute KeyValue conforming to the "deployment.name"semantic conventions
// It represents the The name of the deployment
func DeploymentName(val string) attribute.KeyValue {
	return DeploymentNameKey.String(val)
}

// DeploymentStatus returns an attribute KeyValue conforming to the "deployment.status"semantic conventions
// It represents the The status of the deployment
func DeploymentStatus(val string) attribute.KeyValue {
	return DeploymentStatusKey.String(val)
}

// Enum values for deployment.status
var (

	// failed
	// Stability: experimental
	DeploymentStatusFailed = DeploymentStatusKey.String("failed")
	// succeeded
	// Stability: experimental
	DeploymentStatusSucceeded = DeploymentStatusKey.String("succeeded")
)

// Namespace: destination
const (

	// Destination address - domain name if available without reverse DNS lookup; otherwise, IP address or Unix domain socket name.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "destination.example.com",
	// "10.1.2.80",
	// "/tmp/my.sock",
	//
	// Note: When observed from the source side, and when communicating through an intermediary, `destination.address` SHOULD represent the destination address behind any intermediaries, for example proxies, if it's available
	DestinationAddressKey = attribute.Key("destination.address")
	// Destination port number
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 3389,
	// 2888,
	DestinationPortKey = attribute.Key("destination.port")
)

// DestinationAddress returns an attribute KeyValue conforming to the "destination.address"semantic conventions
// It represents the Destination address - domain name if available without reverse DNS lookup; otherwise, IP address or Unix domain socket name
func DestinationAddress(val string) attribute.KeyValue {
	return DestinationAddressKey.String(val)
}

// DestinationPort returns an attribute KeyValue conforming to the "destination.port"semantic conventions
// It represents the Destination port number
func DestinationPort(val int) attribute.KeyValue {
	return DestinationPortKey.Int(val)
}

// Namespace: device
const (

	// A unique identifier representing the device
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2ab2916d-a51f-4ac8-80ee-45ac31a28092",
	//
	// Note: The device identifier MUST only be defined using the values outlined below. This value is not an advertising identifier and MUST NOT be used as such. On iOS (Swift or Objective-C), this value MUST be equal to the [vendor identifier]. On Android (Java or Kotlin), this value MUST be equal to the Firebase Installation ID or a globally unique UUID which is persisted across sessions in your application. More information can be found [here] on best practices and exact implementation details. Caution should be taken when storing personal data or anything which can identify a user. GDPR and data protection laws may apply, ensure you do your own due diligence
	//
	// [vendor identifier]: https://developer.apple.com/documentation/uikit/uidevice/1620059-identifierforvendor
	// [here]: https://developer.android.com/training/articles/user-data-ids
	DeviceIdKey = attribute.Key("device.id")
	// The name of the device manufacturer
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Apple",
	// "Samsung",
	//
	// Note: The Android OS provides this field via [Build]. iOS apps SHOULD hardcode the value `Apple`
	//
	// [Build]: https://developer.android.com/reference/android/os/Build#MANUFACTURER
	DeviceManufacturerKey = attribute.Key("device.manufacturer")
	// The model identifier for the device
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "iPhone3,4",
	// "SM-G920F",
	//
	// Note: It's recommended this value represents a machine-readable version of the model identifier rather than the market or consumer-friendly name of the device
	DeviceModelIdentifierKey = attribute.Key("device.model.identifier")
	// The marketing name for the device model
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "iPhone 6s Plus",
	// "Samsung Galaxy S6",
	//
	// Note: It's recommended this value represents a human-readable version of the device model rather than a machine-readable alternative
	DeviceModelNameKey = attribute.Key("device.model.name")
)

// DeviceId returns an attribute KeyValue conforming to the "device.id"semantic conventions
// It represents the A unique identifier representing the device
func DeviceId(val string) attribute.KeyValue {
	return DeviceIdKey.String(val)
}

// DeviceManufacturer returns an attribute KeyValue conforming to the "device.manufacturer"semantic conventions
// It represents the The name of the device manufacturer
func DeviceManufacturer(val string) attribute.KeyValue {
	return DeviceManufacturerKey.String(val)
}

// DeviceModelIdentifier returns an attribute KeyValue conforming to the "device.model.identifier"semantic conventions
// It represents the The model identifier for the device
func DeviceModelIdentifier(val string) attribute.KeyValue {
	return DeviceModelIdentifierKey.String(val)
}

// DeviceModelName returns an attribute KeyValue conforming to the "device.model.name"semantic conventions
// It represents the The marketing name for the device model
func DeviceModelName(val string) attribute.KeyValue {
	return DeviceModelNameKey.String(val)
}

// Namespace: disk
const (

	// The disk IO operation direction.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "read",
	DiskIoDirectionKey = attribute.Key("disk.io.direction")
)

// DiskIoDirection returns an attribute KeyValue conforming to the "disk.io.direction"semantic conventions
// It represents the The disk IO operation direction
func DiskIoDirection(val string) attribute.KeyValue {
	return DiskIoDirectionKey.String(val)
}

// Enum values for disk.io.direction
var (

	// read
	// Stability: experimental
	DiskIoDirectionRead = DiskIoDirectionKey.String("read")
	// write
	// Stability: experimental
	DiskIoDirectionWrite = DiskIoDirectionKey.String("write")
)

// Namespace: dns
const (

	// The name being queried.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "www.example.com",
	// "opentelemetry.io",
	//
	// Note: If the name field contains non-printable characters (below 32 or above 126), those characters should be represented as escaped base 10 integers (\DDD). Back slashes and quotes should be escaped. Tabs, carriage returns, and line feeds should be converted to \t, \r, and \n respectively
	DnsQuestionNameKey = attribute.Key("dns.question.name")
)

// DnsQuestionName returns an attribute KeyValue conforming to the "dns.question.name"semantic conventions
// It represents the The name being queried
func DnsQuestionName(val string) attribute.KeyValue {
	return DnsQuestionNameKey.String(val)
}

// Namespace: dotnet
const (

	// Name of the garbage collector managed heap generation.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "gen0",
	// "gen1",
	// "gen2",
	DotnetGcHeapGenerationKey = attribute.Key("dotnet.gc.heap.generation")
)

// DotnetGcHeapGeneration returns an attribute KeyValue conforming to the "dotnet.gc.heap.generation"semantic conventions
// It represents the Name of the garbage collector managed heap generation
func DotnetGcHeapGeneration(val string) attribute.KeyValue {
	return DotnetGcHeapGenerationKey.String(val)
}

// Enum values for dotnet.gc.heap.generation
var (

	// Generation 0
	// Stability: experimental
	DotnetGcHeapGenerationGen0 = DotnetGcHeapGenerationKey.String("gen0")
	// Generation 1
	// Stability: experimental
	DotnetGcHeapGenerationGen1 = DotnetGcHeapGenerationKey.String("gen1")
	// Generation 2
	// Stability: experimental
	DotnetGcHeapGenerationGen2 = DotnetGcHeapGenerationKey.String("gen2")
	// Large Object Heap
	// Stability: experimental
	DotnetGcHeapGenerationLoh = DotnetGcHeapGenerationKey.String("loh")
	// Pinned Object Heap
	// Stability: experimental
	DotnetGcHeapGenerationPoh = DotnetGcHeapGenerationKey.String("poh")
)

// Namespace: enduser
const (

	// Deprecated, use `user.id` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `user.id` attribute.
	//
	// Examples: "username"
	EnduserIdKey = attribute.Key("enduser.id")
	// Deprecated, use `user.roles` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `user.roles` attribute.
	//
	// Examples: "admin"
	EnduserRoleKey = attribute.Key("enduser.role")
	// Deprecated, no replacement at this time.
	// Stability: Experimental
	// Type: string
	// Deprecated: Removed.
	//
	// Examples: "read:message, write:files"
	EnduserScopeKey = attribute.Key("enduser.scope")
)

// EnduserId returns an attribute KeyValue conforming to the "enduser.id"semantic conventions
// It represents the Deprecated, use `user.id` instead
func EnduserId(val string) attribute.KeyValue {
	return EnduserIdKey.String(val)
}

// EnduserRole returns an attribute KeyValue conforming to the "enduser.role"semantic conventions
// It represents the Deprecated, use `user.roles` instead
func EnduserRole(val string) attribute.KeyValue {
	return EnduserRoleKey.String(val)
}

// EnduserScope returns an attribute KeyValue conforming to the "enduser.scope"semantic conventions
// It represents the Deprecated, no replacement at this time
func EnduserScope(val string) attribute.KeyValue {
	return EnduserScopeKey.String(val)
}

// Namespace: error
const (

	// Describes a class of error the operation ended with.
	//
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "timeout",
	// "java.net.UnknownHostException",
	// "server_certificate_invalid",
	// "500",
	//
	// Note: The `error.type` SHOULD be predictable, and SHOULD have low cardinality.
	//
	// When `error.type` is set to a type (e.g., an exception type), its
	// canonical class name identifying the type within the artifact SHOULD be used.
	//
	// Instrumentations SHOULD document the list of errors they report.
	//
	// The cardinality of `error.type` within one instrumentation library SHOULD be low.
	// Telemetry consumers that aggregate data from multiple instrumentation libraries and applications
	// should be prepared for `error.type` to have high cardinality at query time when no
	// additional filters are applied.
	//
	// If the operation has completed successfully, instrumentations SHOULD NOT set `error.type`.
	//
	// If a specific domain defines its own set of error identifiers (such as HTTP or gRPC status codes),
	// it's RECOMMENDED to:
	//
	//   - Use a domain-specific attribute
	//   - Set `error.type` to capture all errors, regardless of whether they are defined within the domain-specific set or not
	ErrorTypeKey = attribute.Key("error.type")
)

// ErrorType returns an attribute KeyValue conforming to the "error.type"semantic conventions
// It represents the Describes a class of error the operation ended with
func ErrorType(val string) attribute.KeyValue {
	return ErrorTypeKey.String(val)
}

// Enum values for error.type
var (

	// A fallback error value to be used when the instrumentation doesn't define a custom value.
	//
	// Stability: stable
	ErrorTypeOther = ErrorTypeKey.String("_OTHER")
)

// Namespace: event
const (

	// Identifies the class / type of event.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "browser.mouse.click",
	// "device.app.lifecycle",
	//
	// Note: Event names are subject to the same rules as [attribute names]. Notably, event names are namespaced to avoid collisions and provide a clean separation of semantics for events in separate domains like browser, mobile, and kubernetes
	//
	// [attribute names]: /docs/general/attribute-naming.md
	EventNameKey = attribute.Key("event.name")
)

// EventName returns an attribute KeyValue conforming to the "event.name"semantic conventions
// It represents the Identifies the class / type of event
func EventName(val string) attribute.KeyValue {
	return EventNameKey.String(val)
}

// Namespace: exception
const (

	// SHOULD be set to true if the exception event is recorded at a point where it is known that the exception is escaping the scope of the span.
	//
	// Stability: Stable
	// Type: boolean
	//
	// Examples: undefined
	// Note: An exception is considered to have escaped (or left) the scope of a span,
	// if that span is ended while the exception is still logically "in flight".
	// This may be actually "in flight" in some languages (e.g. if the exception
	// is passed to a Context manager's `__exit__` method in Python) but will
	// usually be caught at the point of recording the exception in most languages.
	//
	// It is usually not possible to determine at the point where an exception is thrown
	// whether it will escape the scope of a span.
	// However, it is trivial to know that an exception
	// will escape, if one checks for an active exception just before ending the span,
	// as done in the [example for recording span exceptions].
	//
	// It follows that an exception may still escape the scope of the span
	// even if the `exception.escaped` attribute was not set or set to false,
	// since the event might have been recorded at a time where it was not
	// clear whether the exception will escape
	//
	// [example for recording span exceptions]: https://opentelemetry.io/docs/specs/semconv/exceptions/exceptions-spans/#recording-an-exception
	ExceptionEscapedKey = attribute.Key("exception.escaped")
	// The exception message.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "Division by zero",
	// "Can't convert 'int' object to str implicitly",
	ExceptionMessageKey = attribute.Key("exception.message")
	// A stacktrace as a string in the natural representation for the language runtime. The representation is to be determined and documented by each language SIG.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples: "Exception in thread "main" java.lang.RuntimeException: Test exception\n at com.example.GenerateTrace.methodB(GenerateTrace.java:13)\n at com.example.GenerateTrace.methodA(GenerateTrace.java:9)\n at com.example.GenerateTrace.main(GenerateTrace.java:5)\n"
	ExceptionStacktraceKey = attribute.Key("exception.stacktrace")
	// The type of the exception (its fully-qualified class name, if applicable). The dynamic type of the exception should be preferred over the static type in languages that support it.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "java.net.ConnectException",
	// "OSError",
	ExceptionTypeKey = attribute.Key("exception.type")
)

// ExceptionEscaped returns an attribute KeyValue conforming to the "exception.escaped"semantic conventions
// It represents the SHOULD be set to true if the exception event is recorded at a point where it is known that the exception is escaping the scope of the span
func ExceptionEscaped(val bool) attribute.KeyValue {
	return ExceptionEscapedKey.Bool(val)
}

// ExceptionMessage returns an attribute KeyValue conforming to the "exception.message"semantic conventions
// It represents the The exception message
func ExceptionMessage(val string) attribute.KeyValue {
	return ExceptionMessageKey.String(val)
}

// ExceptionStacktrace returns an attribute KeyValue conforming to the "exception.stacktrace"semantic conventions
// It represents the A stacktrace as a string in the natural representation for the language runtime. The representation is to be determined and documented by each language SIG
func ExceptionStacktrace(val string) attribute.KeyValue {
	return ExceptionStacktraceKey.String(val)
}

// ExceptionType returns an attribute KeyValue conforming to the "exception.type"semantic conventions
// It represents the The type of the exception (its fully-qualified class name, if applicable). The dynamic type of the exception should be preferred over the static type in languages that support it
func ExceptionType(val string) attribute.KeyValue {
	return ExceptionTypeKey.String(val)
}

// Namespace: faas
const (

	// A boolean that is true if the serverless function is executed for the first time (aka cold-start).
	//
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	FaasColdstartKey = attribute.Key("faas.coldstart")
	// A string containing the schedule period as [Cron Expression].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "0/5 * * * ? *"
	//
	// [Cron Expression]: https://docs.oracle.com/cd/E12058_01/doc/doc.1014/e12030/cron_expressions.htm
	FaasCronKey = attribute.Key("faas.cron")
	// The name of the source on which the triggering operation was performed. For example, in Cloud Storage or S3 corresponds to the bucket name, and in Cosmos DB to the database name.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "myBucketName",
	// "myDbName",
	FaasDocumentCollectionKey = attribute.Key("faas.document.collection")
	// The document name/table subjected to the operation. For example, in Cloud Storage or S3 is the name of the file, and in Cosmos DB the table name.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "myFile.txt",
	// "myTableName",
	FaasDocumentNameKey = attribute.Key("faas.document.name")
	// Describes the type of the operation that was performed on the data.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	FaasDocumentOperationKey = attribute.Key("faas.document.operation")
	// A string containing the time when the data was accessed in the [ISO 8601] format expressed in [UTC].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "2020-01-23T13:47:06Z"
	//
	// [ISO 8601]: https://www.iso.org/iso-8601-date-and-time-format.html
	// [UTC]: https://www.w3.org/TR/NOTE-datetime
	FaasDocumentTimeKey = attribute.Key("faas.document.time")
	// The execution environment ID as a string, that will be potentially reused for other invocations to the same function/function version.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2021/06/28/[$LATEST]2f399eb14537447da05ab2a2e39309de",
	//
	// Note: * **AWS Lambda:** Use the (full) log stream name
	FaasInstanceKey = attribute.Key("faas.instance")
	// The invocation ID of the current function invocation.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "af9d5aa4-a685-4c5f-a22b-444f80b3cc28"
	FaasInvocationIdKey = attribute.Key("faas.invocation_id")
	// The name of the invoked function.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "my-function"
	// Note: SHOULD be equal to the `faas.name` resource attribute of the invoked function
	FaasInvokedNameKey = attribute.Key("faas.invoked_name")
	// The cloud provider of the invoked function.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: SHOULD be equal to the `cloud.provider` resource attribute of the invoked function
	FaasInvokedProviderKey = attribute.Key("faas.invoked_provider")
	// The cloud region of the invoked function.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "eu-central-1"
	// Note: SHOULD be equal to the `cloud.region` resource attribute of the invoked function
	FaasInvokedRegionKey = attribute.Key("faas.invoked_region")
	// The amount of memory available to the serverless function converted to Bytes.
	//
	// Stability: Experimental
	// Type: int
	//
	// Note: It's recommended to set this attribute since e.g. too little memory can easily stop a Java AWS Lambda function from working correctly. On AWS Lambda, the environment variable `AWS_LAMBDA_FUNCTION_MEMORY_SIZE` provides this information (which must be multiplied by 1,048,576)
	FaasMaxMemoryKey = attribute.Key("faas.max_memory")
	// The name of the single function that this runtime instance executes.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "my-function",
	// "myazurefunctionapp/some-function-name",
	//
	// Note: This is the name of the function as configured/deployed on the FaaS
	// platform and is usually different from the name of the callback
	// function (which may be stored in the
	// [`code.namespace`/`code.function`]
	// span attributes).
	//
	// For some cloud providers, the above definition is ambiguous. The following
	// definition of function name MUST be used for this attribute
	// (and consequently the span name) for the listed cloud providers/products:
	//
	//   - **Azure:**  The full name `<FUNCAPP>/<FUNC>`, i.e., function app name
	//     followed by a forward slash followed by the function name (this form
	//     can also be seen in the resource JSON for the function).
	//     This means that a span attribute MUST be used, as an Azure function
	//     app can host multiple functions that would usually share
	//     a TracerProvider (see also the `cloud.resource_id` attribute)
	// [`code.namespace`/`code.function`]: /docs/general/attributes.md#source-code-attributes
	FaasNameKey = attribute.Key("faas.name")
	// A string containing the function invocation time in the [ISO 8601] format expressed in [UTC].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "2020-01-23T13:47:06Z"
	//
	// [ISO 8601]: https://www.iso.org/iso-8601-date-and-time-format.html
	// [UTC]: https://www.w3.org/TR/NOTE-datetime
	FaasTimeKey = attribute.Key("faas.time")
	// Type of the trigger which caused this function invocation.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	FaasTriggerKey = attribute.Key("faas.trigger")
	// The immutable version of the function being executed.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "26",
	// "pinkfroid-00002",
	//
	// Note: Depending on the cloud provider and platform, use:
	//
	//   - **AWS Lambda:** The [function version]
	//     (an integer represented as a decimal string).
	//   - **Google Cloud Run (Services):** The [revision]
	//     (i.e., the function name plus the revision suffix).
	//   - **Google Cloud Functions:** The value of the
	//     [`K_REVISION` environment variable].
	//   - **Azure Functions:** Not applicable. Do not set this attribute
	// [function version]: https://docs.aws.amazon.com/lambda/latest/dg/configuration-versions.html
	// [revision]: https://cloud.google.com/run/docs/managing/revisions
	// [`K_REVISION` environment variable]: https://cloud.google.com/functions/docs/env-var#runtime_environment_variables_set_automatically
	FaasVersionKey = attribute.Key("faas.version")
)

// FaasColdstart returns an attribute KeyValue conforming to the "faas.coldstart"semantic conventions
// It represents the A boolean that is true if the serverless function is executed for the first time (aka cold-start)
func FaasColdstart(val bool) attribute.KeyValue {
	return FaasColdstartKey.Bool(val)
}

// FaasCron returns an attribute KeyValue conforming to the "faas.cron"semantic conventions
// It represents the A string containing the schedule period as [Cron Expression]
//
// [Cron Expression]: https://docs.oracle.com/cd/E12058_01/doc/doc.1014/e12030/cron_expressions.htm
func FaasCron(val string) attribute.KeyValue {
	return FaasCronKey.String(val)
}

// FaasDocumentCollection returns an attribute KeyValue conforming to the "faas.document.collection"semantic conventions
// It represents the The name of the source on which the triggering operation was performed. For example, in Cloud Storage or S3 corresponds to the bucket name, and in Cosmos DB to the database name
func FaasDocumentCollection(val string) attribute.KeyValue {
	return FaasDocumentCollectionKey.String(val)
}

// FaasDocumentName returns an attribute KeyValue conforming to the "faas.document.name"semantic conventions
// It represents the The document name/table subjected to the operation. For example, in Cloud Storage or S3 is the name of the file, and in Cosmos DB the table name
func FaasDocumentName(val string) attribute.KeyValue {
	return FaasDocumentNameKey.String(val)
}

// FaasDocumentOperation returns an attribute KeyValue conforming to the "faas.document.operation"semantic conventions
// It represents the Describes the type of the operation that was performed on the data
func FaasDocumentOperation(val string) attribute.KeyValue {
	return FaasDocumentOperationKey.String(val)
}

// FaasDocumentTime returns an attribute KeyValue conforming to the "faas.document.time"semantic conventions
// It represents the A string containing the time when the data was accessed in the [ISO 8601] format expressed in [UTC]
//
// [ISO 8601]: https://www.iso.org/iso-8601-date-and-time-format.html
// [UTC]: https://www.w3.org/TR/NOTE-datetime
func FaasDocumentTime(val string) attribute.KeyValue {
	return FaasDocumentTimeKey.String(val)
}

// FaasInstance returns an attribute KeyValue conforming to the "faas.instance"semantic conventions
// It represents the The execution environment ID as a string, that will be potentially reused for other invocations to the same function/function version
func FaasInstance(val string) attribute.KeyValue {
	return FaasInstanceKey.String(val)
}

// FaasInvocationId returns an attribute KeyValue conforming to the "faas.invocation_id"semantic conventions
// It represents the The invocation ID of the current function invocation
func FaasInvocationId(val string) attribute.KeyValue {
	return FaasInvocationIdKey.String(val)
}

// FaasInvokedName returns an attribute KeyValue conforming to the "faas.invoked_name"semantic conventions
// It represents the The name of the invoked function
func FaasInvokedName(val string) attribute.KeyValue {
	return FaasInvokedNameKey.String(val)
}

// FaasInvokedProvider returns an attribute KeyValue conforming to the "faas.invoked_provider"semantic conventions
// It represents the The cloud provider of the invoked function
func FaasInvokedProvider(val string) attribute.KeyValue {
	return FaasInvokedProviderKey.String(val)
}

// FaasInvokedRegion returns an attribute KeyValue conforming to the "faas.invoked_region"semantic conventions
// It represents the The cloud region of the invoked function
func FaasInvokedRegion(val string) attribute.KeyValue {
	return FaasInvokedRegionKey.String(val)
}

// FaasMaxMemory returns an attribute KeyValue conforming to the "faas.max_memory"semantic conventions
// It represents the The amount of memory available to the serverless function converted to Bytes
func FaasMaxMemory(val int) attribute.KeyValue {
	return FaasMaxMemoryKey.Int(val)
}

// FaasName returns an attribute KeyValue conforming to the "faas.name"semantic conventions
// It represents the The name of the single function that this runtime instance executes
func FaasName(val string) attribute.KeyValue {
	return FaasNameKey.String(val)
}

// FaasTime returns an attribute KeyValue conforming to the "faas.time"semantic conventions
// It represents the A string containing the function invocation time in the [ISO 8601] format expressed in [UTC]
//
// [ISO 8601]: https://www.iso.org/iso-8601-date-and-time-format.html
// [UTC]: https://www.w3.org/TR/NOTE-datetime
func FaasTime(val string) attribute.KeyValue {
	return FaasTimeKey.String(val)
}

// FaasTrigger returns an attribute KeyValue conforming to the "faas.trigger"semantic conventions
// It represents the Type of the trigger which caused this function invocation
func FaasTrigger(val string) attribute.KeyValue {
	return FaasTriggerKey.String(val)
}

// FaasVersion returns an attribute KeyValue conforming to the "faas.version"semantic conventions
// It represents the The immutable version of the function being executed
func FaasVersion(val string) attribute.KeyValue {
	return FaasVersionKey.String(val)
}

// Enum values for faas.document.operation
var (

	// When a new object is created.
	// Stability: experimental
	FaasDocumentOperationInsert = FaasDocumentOperationKey.String("insert")
	// When an object is modified.
	// Stability: experimental
	FaasDocumentOperationEdit = FaasDocumentOperationKey.String("edit")
	// When an object is deleted.
	// Stability: experimental
	FaasDocumentOperationDelete = FaasDocumentOperationKey.String("delete")
)

// Enum values for faas.invoked_provider
var (

	// Alibaba Cloud
	// Stability: experimental
	FaasInvokedProviderAlibabaCloud = FaasInvokedProviderKey.String("alibaba_cloud")
	// Amazon Web Services
	// Stability: experimental
	FaasInvokedProviderAws = FaasInvokedProviderKey.String("aws")
	// Microsoft Azure
	// Stability: experimental
	FaasInvokedProviderAzure = FaasInvokedProviderKey.String("azure")
	// Google Cloud Platform
	// Stability: experimental
	FaasInvokedProviderGcp = FaasInvokedProviderKey.String("gcp")
	// Tencent Cloud
	// Stability: experimental
	FaasInvokedProviderTencentCloud = FaasInvokedProviderKey.String("tencent_cloud")
)

// Enum values for faas.trigger
var (

	// A response to some data source operation such as a database or filesystem read/write
	// Stability: experimental
	FaasTriggerDatasource = FaasTriggerKey.String("datasource")
	// To provide an answer to an inbound HTTP request
	// Stability: experimental
	FaasTriggerHttp = FaasTriggerKey.String("http")
	// A function is set to be executed when messages are sent to a messaging system
	// Stability: experimental
	FaasTriggerPubsub = FaasTriggerKey.String("pubsub")
	// A function is scheduled to be executed regularly
	// Stability: experimental
	FaasTriggerTimer = FaasTriggerKey.String("timer")
	// If none of the others apply
	// Stability: experimental
	FaasTriggerOther = FaasTriggerKey.String("other")
)

// Namespace: feature_flag
const (

	// The unique identifier of the feature flag.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "logo-color",
	FeatureFlagKeyKey = attribute.Key("feature_flag.key")
	// The name of the service provider that performs the flag evaluation.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Flag Manager",
	FeatureFlagProviderNameKey = attribute.Key("feature_flag.provider_name")
	// SHOULD be a semantic identifier for a value. If one is unavailable, a stringified version of the value can be used.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "red",
	// "true",
	// "on",
	//
	// Note: A semantic identifier, commonly referred to as a variant, provides a means
	// for referring to a value without including the value itself. This can
	// provide additional context for understanding the meaning behind a value.
	// For example, the variant `red` maybe be used for the value `#c05543`.
	//
	// A stringified version of the value can be used in situations where a
	// semantic identifier is unavailable. String representation of the value
	// should be determined by the implementer
	FeatureFlagVariantKey = attribute.Key("feature_flag.variant")
)

// FeatureFlagKey returns an attribute KeyValue conforming to the "feature_flag.key"semantic conventions
// It represents the The unique identifier of the feature flag
func FeatureFlagKey(val string) attribute.KeyValue {
	return FeatureFlagKeyKey.String(val)
}

// FeatureFlagProviderName returns an attribute KeyValue conforming to the "feature_flag.provider_name"semantic conventions
// It represents the The name of the service provider that performs the flag evaluation
func FeatureFlagProviderName(val string) attribute.KeyValue {
	return FeatureFlagProviderNameKey.String(val)
}

// FeatureFlagVariant returns an attribute KeyValue conforming to the "feature_flag.variant"semantic conventions
// It represents the SHOULD be a semantic identifier for a value. If one is unavailable, a stringified version of the value can be used
func FeatureFlagVariant(val string) attribute.KeyValue {
	return FeatureFlagVariantKey.String(val)
}

// Namespace: file
const (

	// Directory where the file is located. It should include the drive letter, when appropriate.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/home/user",
	// "C:\Program Files\MyApp",
	FileDirectoryKey = attribute.Key("file.directory")
	// File extension, excluding the leading dot.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "png",
	// "gz",
	//
	// Note: When the file name has multiple extensions (example.tar.gz), only the last one should be captured ("gz", not "tar.gz")
	FileExtensionKey = attribute.Key("file.extension")
	// Name of the file including the extension, without the directory.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "example.png",
	FileNameKey = attribute.Key("file.name")
	// Full path to the file, including the file name. It should include the drive letter, when appropriate.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/home/alice/example.png",
	// "C:\Program Files\MyApp\myapp.exe",
	FilePathKey = attribute.Key("file.path")
	// File size in bytes.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples: undefined
	FileSizeKey = attribute.Key("file.size")
)

// FileDirectory returns an attribute KeyValue conforming to the "file.directory"semantic conventions
// It represents the Directory where the file is located. It should include the drive letter, when appropriate
func FileDirectory(val string) attribute.KeyValue {
	return FileDirectoryKey.String(val)
}

// FileExtension returns an attribute KeyValue conforming to the "file.extension"semantic conventions
// It represents the File extension, excluding the leading dot
func FileExtension(val string) attribute.KeyValue {
	return FileExtensionKey.String(val)
}

// FileName returns an attribute KeyValue conforming to the "file.name"semantic conventions
// It represents the Name of the file including the extension, without the directory
func FileName(val string) attribute.KeyValue {
	return FileNameKey.String(val)
}

// FilePath returns an attribute KeyValue conforming to the "file.path"semantic conventions
// It represents the Full path to the file, including the file name. It should include the drive letter, when appropriate
func FilePath(val string) attribute.KeyValue {
	return FilePathKey.String(val)
}

// FileSize returns an attribute KeyValue conforming to the "file.size"semantic conventions
// It represents the File size in bytes
func FileSize(val int) attribute.KeyValue {
	return FileSizeKey.Int(val)
}

// Namespace: gcp
const (

	// Identifies the Google Cloud service for which the official client library is intended.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "appengine",
	// "run",
	// "firestore",
	// "alloydb",
	// "spanner",
	//
	// Note: Intended to be a stable identifier for Google Cloud client libraries that is uniform across implementation languages. The value should be derived from the canonical service domain for the service; for example, 'foo.googleapis.com' should result in a value of 'foo'
	GcpClientServiceKey = attribute.Key("gcp.client.service")
	// The name of the Cloud Run [execution] being run for the Job, as set by the [`CLOUD_RUN_EXECUTION`] environment variable.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "job-name-xxxx",
	// "sample-job-mdw84",
	//
	// [execution]: https://cloud.google.com/run/docs/managing/job-executions
	// [`CLOUD_RUN_EXECUTION`]: https://cloud.google.com/run/docs/container-contract#jobs-env-vars
	GcpCloudRunJobExecutionKey = attribute.Key("gcp.cloud_run.job.execution")
	// The index for a task within an execution as provided by the [`CLOUD_RUN_TASK_INDEX`] environment variable.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 0,
	// 1,
	//
	// [`CLOUD_RUN_TASK_INDEX`]: https://cloud.google.com/run/docs/container-contract#jobs-env-vars
	GcpCloudRunJobTaskIndexKey = attribute.Key("gcp.cloud_run.job.task_index")
	// The hostname of a GCE instance. This is the full value of the default or [custom hostname].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "my-host1234.example.com",
	// "sample-vm.us-west1-b.c.my-project.internal",
	//
	// [custom hostname]: https://cloud.google.com/compute/docs/instances/custom-hostname-vm
	GcpGceInstanceHostnameKey = attribute.Key("gcp.gce.instance.hostname")
	// The instance name of a GCE instance. This is the value provided by `host.name`, the visible name of the instance in the Cloud Console UI, and the prefix for the default hostname of the instance as defined by the [default internal DNS name].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "instance-1",
	// "my-vm-name",
	//
	// [default internal DNS name]: https://cloud.google.com/compute/docs/internal-dns#instance-fully-qualified-domain-names
	GcpGceInstanceNameKey = attribute.Key("gcp.gce.instance.name")
)

// GcpClientService returns an attribute KeyValue conforming to the "gcp.client.service"semantic conventions
// It represents the Identifies the Google Cloud service for which the official client library is intended
func GcpClientService(val string) attribute.KeyValue {
	return GcpClientServiceKey.String(val)
}

// GcpCloudRunJobExecution returns an attribute KeyValue conforming to the "gcp.cloud_run.job.execution"semantic conventions
// It represents the The name of the Cloud Run [execution] being run for the Job, as set by the [`CLOUD_RUN_EXECUTION`] environment variable
//
// [execution]: https://cloud.google.com/run/docs/managing/job-executions
// [`CLOUD_RUN_EXECUTION`]: https://cloud.google.com/run/docs/container-contract#jobs-env-vars
func GcpCloudRunJobExecution(val string) attribute.KeyValue {
	return GcpCloudRunJobExecutionKey.String(val)
}

// GcpCloudRunJobTaskIndex returns an attribute KeyValue conforming to the "gcp.cloud_run.job.task_index"semantic conventions
// It represents the The index for a task within an execution as provided by the [`CLOUD_RUN_TASK_INDEX`] environment variable
//
// [`CLOUD_RUN_TASK_INDEX`]: https://cloud.google.com/run/docs/container-contract#jobs-env-vars
func GcpCloudRunJobTaskIndex(val int) attribute.KeyValue {
	return GcpCloudRunJobTaskIndexKey.Int(val)
}

// GcpGceInstanceHostname returns an attribute KeyValue conforming to the "gcp.gce.instance.hostname"semantic conventions
// It represents the The hostname of a GCE instance. This is the full value of the default or [custom hostname]
//
// [custom hostname]: https://cloud.google.com/compute/docs/instances/custom-hostname-vm
func GcpGceInstanceHostname(val string) attribute.KeyValue {
	return GcpGceInstanceHostnameKey.String(val)
}

// GcpGceInstanceName returns an attribute KeyValue conforming to the "gcp.gce.instance.name"semantic conventions
// It represents the The instance name of a GCE instance. This is the value provided by `host.name`, the visible name of the instance in the Cloud Console UI, and the prefix for the default hostname of the instance as defined by the [default internal DNS name]
//
// [default internal DNS name]: https://cloud.google.com/compute/docs/internal-dns#instance-fully-qualified-domain-names
func GcpGceInstanceName(val string) attribute.KeyValue {
	return GcpGceInstanceNameKey.String(val)
}

// Namespace: gen_ai
const (

	// The full response received from the GenAI model.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "[{'role': 'assistant', 'content': 'The capital of France is Paris.'}]",
	//
	// Note: It's RECOMMENDED to format completions as JSON string matching [OpenAI messages format]
	//
	// [OpenAI messages format]: https://platform.openai.com/docs/guides/text-generation
	GenAiCompletionKey = attribute.Key("gen_ai.completion")
	// The name of the operation being performed.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: If one of the predefined values applies, but specific system uses a different name it's RECOMMENDED to document it in the semantic conventions for specific GenAI system and use system-specific name in the instrumentation. If a different name is not documented, instrumentation libraries SHOULD use applicable predefined value
	GenAiOperationNameKey = attribute.Key("gen_ai.operation.name")
	// The full prompt sent to the GenAI model.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "[{'role': 'user', 'content': 'What is the capital of France?'}]",
	//
	// Note: It's RECOMMENDED to format prompts as JSON string matching [OpenAI messages format]
	//
	// [OpenAI messages format]: https://platform.openai.com/docs/guides/text-generation
	GenAiPromptKey = attribute.Key("gen_ai.prompt")
	// The frequency penalty setting for the GenAI request.
	// Stability: Experimental
	// Type: double
	//
	// Examples:
	// 0.1,
	GenAiRequestFrequencyPenaltyKey = attribute.Key("gen_ai.request.frequency_penalty")
	// The maximum number of tokens the model generates for a request.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 100,
	GenAiRequestMaxTokensKey = attribute.Key("gen_ai.request.max_tokens")
	// The name of the GenAI model a request is being made to.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "gpt-4"
	GenAiRequestModelKey = attribute.Key("gen_ai.request.model")
	// The presence penalty setting for the GenAI request.
	// Stability: Experimental
	// Type: double
	//
	// Examples:
	// 0.1,
	GenAiRequestPresencePenaltyKey = attribute.Key("gen_ai.request.presence_penalty")
	// List of sequences that the model will use to stop generating further tokens.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "forest",
	// "lived",
	// ],
	GenAiRequestStopSequencesKey = attribute.Key("gen_ai.request.stop_sequences")
	// The temperature setting for the GenAI request.
	// Stability: Experimental
	// Type: double
	//
	// Examples:
	// 0.0,
	GenAiRequestTemperatureKey = attribute.Key("gen_ai.request.temperature")
	// The top_k sampling setting for the GenAI request.
	// Stability: Experimental
	// Type: double
	//
	// Examples:
	// 1.0,
	GenAiRequestTopKKey = attribute.Key("gen_ai.request.top_k")
	// The top_p sampling setting for the GenAI request.
	// Stability: Experimental
	// Type: double
	//
	// Examples:
	// 1.0,
	GenAiRequestTopPKey = attribute.Key("gen_ai.request.top_p")
	// Array of reasons the model stopped generating tokens, corresponding to each generation received.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "stop",
	// ],
	// [
	// "stop",
	// "length",
	// ],
	GenAiResponseFinishReasonsKey = attribute.Key("gen_ai.response.finish_reasons")
	// The unique identifier for the completion.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "chatcmpl-123",
	GenAiResponseIdKey = attribute.Key("gen_ai.response.id")
	// The name of the model that generated the response.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "gpt-4-0613",
	GenAiResponseModelKey = attribute.Key("gen_ai.response.model")
	// The Generative AI product as identified by the client or server instrumentation.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: "openai"
	// Note: The `gen_ai.system` describes a family of GenAI models with specific model identified
	// by `gen_ai.request.model` and `gen_ai.response.model` attributes.
	//
	// The actual GenAI product may differ from the one identified by the client.
	// For example, when using OpenAI client libraries to communicate with Mistral, the `gen_ai.system`
	// is set to `openai` based on the instrumentation's best knowledge.
	//
	// For custom model, a custom friendly name SHOULD be used.
	// If none of these options apply, the `gen_ai.system` SHOULD be set to `_OTHER`
	GenAiSystemKey = attribute.Key("gen_ai.system")
	// The type of token being counted.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "input",
	// "output",
	GenAiTokenTypeKey = attribute.Key("gen_ai.token.type")
	// Deprecated, use `gen_ai.usage.output_tokens` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `gen_ai.usage.output_tokens` attribute.
	//
	// Examples:
	// 42,
	GenAiUsageCompletionTokensKey = attribute.Key("gen_ai.usage.completion_tokens")
	// The number of tokens used in the GenAI input (prompt).
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 100,
	GenAiUsageInputTokensKey = attribute.Key("gen_ai.usage.input_tokens")
	// The number of tokens used in the GenAI response (completion).
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 180,
	GenAiUsageOutputTokensKey = attribute.Key("gen_ai.usage.output_tokens")
	// Deprecated, use `gen_ai.usage.input_tokens` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `gen_ai.usage.input_tokens` attribute.
	//
	// Examples:
	// 42,
	GenAiUsagePromptTokensKey = attribute.Key("gen_ai.usage.prompt_tokens")
)

// GenAiCompletion returns an attribute KeyValue conforming to the "gen_ai.completion"semantic conventions
// It represents the The full response received from the GenAI model
func GenAiCompletion(val string) attribute.KeyValue {
	return GenAiCompletionKey.String(val)
}

// GenAiOperationName returns an attribute KeyValue conforming to the "gen_ai.operation.name"semantic conventions
// It represents the The name of the operation being performed
func GenAiOperationName(val string) attribute.KeyValue {
	return GenAiOperationNameKey.String(val)
}

// GenAiPrompt returns an attribute KeyValue conforming to the "gen_ai.prompt"semantic conventions
// It represents the The full prompt sent to the GenAI model
func GenAiPrompt(val string) attribute.KeyValue {
	return GenAiPromptKey.String(val)
}

// GenAiRequestFrequencyPenalty returns an attribute KeyValue conforming to the "gen_ai.request.frequency_penalty"semantic conventions
// It represents the The frequency penalty setting for the GenAI request
func GenAiRequestFrequencyPenalty(val float64) attribute.KeyValue {
	return GenAiRequestFrequencyPenaltyKey.Float64(val)
}

// GenAiRequestMaxTokens returns an attribute KeyValue conforming to the "gen_ai.request.max_tokens"semantic conventions
// It represents the The maximum number of tokens the model generates for a request
func GenAiRequestMaxTokens(val int) attribute.KeyValue {
	return GenAiRequestMaxTokensKey.Int(val)
}

// GenAiRequestModel returns an attribute KeyValue conforming to the "gen_ai.request.model"semantic conventions
// It represents the The name of the GenAI model a request is being made to
func GenAiRequestModel(val string) attribute.KeyValue {
	return GenAiRequestModelKey.String(val)
}

// GenAiRequestPresencePenalty returns an attribute KeyValue conforming to the "gen_ai.request.presence_penalty"semantic conventions
// It represents the The presence penalty setting for the GenAI request
func GenAiRequestPresencePenalty(val float64) attribute.KeyValue {
	return GenAiRequestPresencePenaltyKey.Float64(val)
}

// GenAiRequestStopSequences returns an attribute KeyValue conforming to the "gen_ai.request.stop_sequences"semantic conventions
// It represents the List of sequences that the model will use to stop generating further tokens
func GenAiRequestStopSequences(val []string) attribute.KeyValue {
	return GenAiRequestStopSequencesKey.StringSlice(val)
}

// GenAiRequestTemperature returns an attribute KeyValue conforming to the "gen_ai.request.temperature"semantic conventions
// It represents the The temperature setting for the GenAI request
func GenAiRequestTemperature(val float64) attribute.KeyValue {
	return GenAiRequestTemperatureKey.Float64(val)
}

// GenAiRequestTopK returns an attribute KeyValue conforming to the "gen_ai.request.top_k"semantic conventions
// It represents the The top_k sampling setting for the GenAI request
func GenAiRequestTopK(val float64) attribute.KeyValue {
	return GenAiRequestTopKKey.Float64(val)
}

// GenAiRequestTopP returns an attribute KeyValue conforming to the "gen_ai.request.top_p"semantic conventions
// It represents the The top_p sampling setting for the GenAI request
func GenAiRequestTopP(val float64) attribute.KeyValue {
	return GenAiRequestTopPKey.Float64(val)
}

// GenAiResponseFinishReasons returns an attribute KeyValue conforming to the "gen_ai.response.finish_reasons"semantic conventions
// It represents the Array of reasons the model stopped generating tokens, corresponding to each generation received
func GenAiResponseFinishReasons(val []string) attribute.KeyValue {
	return GenAiResponseFinishReasonsKey.StringSlice(val)
}

// GenAiResponseId returns an attribute KeyValue conforming to the "gen_ai.response.id"semantic conventions
// It represents the The unique identifier for the completion
func GenAiResponseId(val string) attribute.KeyValue {
	return GenAiResponseIdKey.String(val)
}

// GenAiResponseModel returns an attribute KeyValue conforming to the "gen_ai.response.model"semantic conventions
// It represents the The name of the model that generated the response
func GenAiResponseModel(val string) attribute.KeyValue {
	return GenAiResponseModelKey.String(val)
}

// GenAiSystem returns an attribute KeyValue conforming to the "gen_ai.system"semantic conventions
// It represents the The Generative AI product as identified by the client or server instrumentation
func GenAiSystem(val string) attribute.KeyValue {
	return GenAiSystemKey.String(val)
}

// GenAiTokenType returns an attribute KeyValue conforming to the "gen_ai.token.type"semantic conventions
// It represents the The type of token being counted
func GenAiTokenType(val string) attribute.KeyValue {
	return GenAiTokenTypeKey.String(val)
}

// GenAiUsageCompletionTokens returns an attribute KeyValue conforming to the "gen_ai.usage.completion_tokens"semantic conventions
// It represents the Deprecated, use `gen_ai.usage.output_tokens` instead
func GenAiUsageCompletionTokens(val int) attribute.KeyValue {
	return GenAiUsageCompletionTokensKey.Int(val)
}

// GenAiUsageInputTokens returns an attribute KeyValue conforming to the "gen_ai.usage.input_tokens"semantic conventions
// It represents the The number of tokens used in the GenAI input (prompt)
func GenAiUsageInputTokens(val int) attribute.KeyValue {
	return GenAiUsageInputTokensKey.Int(val)
}

// GenAiUsageOutputTokens returns an attribute KeyValue conforming to the "gen_ai.usage.output_tokens"semantic conventions
// It represents the The number of tokens used in the GenAI response (completion)
func GenAiUsageOutputTokens(val int) attribute.KeyValue {
	return GenAiUsageOutputTokensKey.Int(val)
}

// GenAiUsagePromptTokens returns an attribute KeyValue conforming to the "gen_ai.usage.prompt_tokens"semantic conventions
// It represents the Deprecated, use `gen_ai.usage.input_tokens` instead
func GenAiUsagePromptTokens(val int) attribute.KeyValue {
	return GenAiUsagePromptTokensKey.Int(val)
}

// Enum values for gen_ai.operation.name
var (

	// Chat completion operation such as [OpenAI Chat API]
	// Stability: experimental
	//
	// [OpenAI Chat API]: https://platform.openai.com/docs/api-reference/chat
	GenAiOperationNameChat = GenAiOperationNameKey.String("chat")
	// Text completions operation such as [OpenAI Completions API (Legacy)]
	// Stability: experimental
	//
	// [OpenAI Completions API (Legacy)]: https://platform.openai.com/docs/api-reference/completions
	GenAiOperationNameTextCompletion = GenAiOperationNameKey.String("text_completion")
)

// Enum values for gen_ai.system
var (

	// OpenAI
	// Stability: experimental
	GenAiSystemOpenai = GenAiSystemKey.String("openai")
	// Vertex AI
	// Stability: experimental
	GenAiSystemVertexAi = GenAiSystemKey.String("vertex_ai")
	// Anthropic
	// Stability: experimental
	GenAiSystemAnthropic = GenAiSystemKey.String("anthropic")
	// Cohere
	// Stability: experimental
	GenAiSystemCohere = GenAiSystemKey.String("cohere")
)

// Enum values for gen_ai.token.type
var (

	// Input tokens (prompt, input, etc.)
	// Stability: experimental
	GenAiTokenTypeInput = GenAiTokenTypeKey.String("input")
	// Output tokens (completion, response, etc.)
	// Stability: experimental
	GenAiTokenTypeCompletion = GenAiTokenTypeKey.String("output")
)

// Namespace: go
const (

	// The type of memory.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "other",
	// "stack",
	GoMemoryTypeKey = attribute.Key("go.memory.type")
)

// GoMemoryType returns an attribute KeyValue conforming to the "go.memory.type"semantic conventions
// It represents the The type of memory
func GoMemoryType(val string) attribute.KeyValue {
	return GoMemoryTypeKey.String(val)
}

// Enum values for go.memory.type
var (

	// Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.
	// Stability: experimental
	GoMemoryTypeStack = GoMemoryTypeKey.String("stack")
	// Memory used by the Go runtime, excluding other categories of memory usage described in this enumeration.
	// Stability: experimental
	GoMemoryTypeOther = GoMemoryTypeKey.String("other")
)

// Namespace: graphql
const (

	// The GraphQL document being executed.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "query findBookById { bookById(id: ?) { name } }"
	// Note: The value may be sanitized to exclude sensitive information
	GraphqlDocumentKey = attribute.Key("graphql.document")
	// The name of the operation being executed.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "findBookById"
	GraphqlOperationNameKey = attribute.Key("graphql.operation.name")
	// The type of the operation being executed.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "query",
	// "mutation",
	// "subscription",
	GraphqlOperationTypeKey = attribute.Key("graphql.operation.type")
)

// GraphqlDocument returns an attribute KeyValue conforming to the "graphql.document"semantic conventions
// It represents the The GraphQL document being executed
func GraphqlDocument(val string) attribute.KeyValue {
	return GraphqlDocumentKey.String(val)
}

// GraphqlOperationName returns an attribute KeyValue conforming to the "graphql.operation.name"semantic conventions
// It represents the The name of the operation being executed
func GraphqlOperationName(val string) attribute.KeyValue {
	return GraphqlOperationNameKey.String(val)
}

// GraphqlOperationType returns an attribute KeyValue conforming to the "graphql.operation.type"semantic conventions
// It represents the The type of the operation being executed
func GraphqlOperationType(val string) attribute.KeyValue {
	return GraphqlOperationTypeKey.String(val)
}

// Enum values for graphql.operation.type
var (

	// GraphQL query
	// Stability: experimental
	GraphqlOperationTypeQuery = GraphqlOperationTypeKey.String("query")
	// GraphQL mutation
	// Stability: experimental
	GraphqlOperationTypeMutation = GraphqlOperationTypeKey.String("mutation")
	// GraphQL subscription
	// Stability: experimental
	GraphqlOperationTypeSubscription = GraphqlOperationTypeKey.String("subscription")
)

// Namespace: heroku
const (

	// Unique identifier for the application
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2daa2797-e42b-4624-9322-ec3f968df4da",
	HerokuAppIdKey = attribute.Key("heroku.app.id")
	// Commit hash for the current release
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "e6134959463efd8966b20e75b913cafe3f5ec",
	HerokuReleaseCommitKey = attribute.Key("heroku.release.commit")
	// Time and date the release was created
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2022-10-23T18:00:42Z",
	HerokuReleaseCreationTimestampKey = attribute.Key("heroku.release.creation_timestamp")
)

// HerokuAppId returns an attribute KeyValue conforming to the "heroku.app.id"semantic conventions
// It represents the Unique identifier for the application
func HerokuAppId(val string) attribute.KeyValue {
	return HerokuAppIdKey.String(val)
}

// HerokuReleaseCommit returns an attribute KeyValue conforming to the "heroku.release.commit"semantic conventions
// It represents the Commit hash for the current release
func HerokuReleaseCommit(val string) attribute.KeyValue {
	return HerokuReleaseCommitKey.String(val)
}

// HerokuReleaseCreationTimestamp returns an attribute KeyValue conforming to the "heroku.release.creation_timestamp"semantic conventions
// It represents the Time and date the release was created
func HerokuReleaseCreationTimestamp(val string) attribute.KeyValue {
	return HerokuReleaseCreationTimestampKey.String(val)
}

// Namespace: host
const (

	// The CPU architecture the host system is running on.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	HostArchKey = attribute.Key("host.arch")
	// The amount of level 2 memory cache available to the processor (in Bytes).
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 12288000,
	HostCpuCacheL2SizeKey = attribute.Key("host.cpu.cache.l2.size")
	// Family or generation of the CPU.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "6",
	// "PA-RISC 1.1e",
	HostCpuFamilyKey = attribute.Key("host.cpu.family")
	// Model identifier. It provides more granular information about the CPU, distinguishing it from other CPUs within the same family.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "6",
	// "9000/778/B180L",
	HostCpuModelIdKey = attribute.Key("host.cpu.model.id")
	// Model designation of the processor.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz",
	HostCpuModelNameKey = attribute.Key("host.cpu.model.name")
	// Stepping or core revisions.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "1",
	// "r1p1",
	HostCpuSteppingKey = attribute.Key("host.cpu.stepping")
	// Processor manufacturer identifier. A maximum 12-character string.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "GenuineIntel",
	//
	// Note: [CPUID] command returns the vendor ID string in EBX, EDX and ECX registers. Writing these to memory in this order results in a 12-character string
	//
	// [CPUID]: https://wiki.osdev.org/CPUID
	HostCpuVendorIdKey = attribute.Key("host.cpu.vendor.id")
	// Unique host ID. For Cloud, this must be the instance_id assigned by the cloud provider. For non-containerized systems, this should be the `machine-id`. See the table below for the sources to use to determine the `machine-id` based on operating system.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "fdbf79e8af94cb7f9e8df36789187052",
	HostIdKey = attribute.Key("host.id")
	// VM image ID or host OS image ID. For Cloud, this value is from the provider.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "ami-07b06b442921831e5",
	HostImageIdKey = attribute.Key("host.image.id")
	// Name of the VM image or OS install the host was instantiated from.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "infra-ami-eks-worker-node-7d4ec78312",
	// "CentOS-8-x86_64-1905",
	HostImageNameKey = attribute.Key("host.image.name")
	// The version string of the VM image or host OS as defined in [Version Attributes].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "0.1",
	//
	// [Version Attributes]: /docs/resource/README.md#version-attributes
	HostImageVersionKey = attribute.Key("host.image.version")
	// Available IP addresses of the host, excluding loopback interfaces.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "192.168.1.140",
	// "fe80::abc2:4a28:737a:609e",
	// ],
	//
	// Note: IPv4 Addresses MUST be specified in dotted-quad notation. IPv6 addresses MUST be specified in the [RFC 5952] format
	//
	// [RFC 5952]: https://www.rfc-editor.org/rfc/rfc5952.html
	HostIpKey = attribute.Key("host.ip")
	// Available MAC addresses of the host, excluding loopback interfaces.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "AC-DE-48-23-45-67",
	// "AC-DE-48-23-45-67-01-9F",
	// ],
	//
	// Note: MAC Addresses MUST be represented in [IEEE RA hexadecimal form]: as hyphen-separated octets in uppercase hexadecimal form from most to least significant
	//
	// [IEEE RA hexadecimal form]: https://standards.ieee.org/wp-content/uploads/import/documents/tutorials/eui.pdf
	HostMacKey = attribute.Key("host.mac")
	// Name of the host. On Unix systems, it may contain what the hostname command returns, or the fully qualified hostname, or another name specified by the user.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry-test",
	HostNameKey = attribute.Key("host.name")
	// Type of host. For Cloud, this must be the machine type.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "n1-standard-1",
	HostTypeKey = attribute.Key("host.type")
)

// HostArch returns an attribute KeyValue conforming to the "host.arch"semantic conventions
// It represents the The CPU architecture the host system is running on
func HostArch(val string) attribute.KeyValue {
	return HostArchKey.String(val)
}

// HostCpuCacheL2Size returns an attribute KeyValue conforming to the "host.cpu.cache.l2.size"semantic conventions
// It represents the The amount of level 2 memory cache available to the processor (in Bytes)
func HostCpuCacheL2Size(val int) attribute.KeyValue {
	return HostCpuCacheL2SizeKey.Int(val)
}

// HostCpuFamily returns an attribute KeyValue conforming to the "host.cpu.family"semantic conventions
// It represents the Family or generation of the CPU
func HostCpuFamily(val string) attribute.KeyValue {
	return HostCpuFamilyKey.String(val)
}

// HostCpuModelId returns an attribute KeyValue conforming to the "host.cpu.model.id"semantic conventions
// It represents the Model identifier. It provides more granular information about the CPU, distinguishing it from other CPUs within the same family
func HostCpuModelId(val string) attribute.KeyValue {
	return HostCpuModelIdKey.String(val)
}

// HostCpuModelName returns an attribute KeyValue conforming to the "host.cpu.model.name"semantic conventions
// It represents the Model designation of the processor
func HostCpuModelName(val string) attribute.KeyValue {
	return HostCpuModelNameKey.String(val)
}

// HostCpuStepping returns an attribute KeyValue conforming to the "host.cpu.stepping"semantic conventions
// It represents the Stepping or core revisions
func HostCpuStepping(val string) attribute.KeyValue {
	return HostCpuSteppingKey.String(val)
}

// HostCpuVendorId returns an attribute KeyValue conforming to the "host.cpu.vendor.id"semantic conventions
// It represents the Processor manufacturer identifier. A maximum 12-character string
func HostCpuVendorId(val string) attribute.KeyValue {
	return HostCpuVendorIdKey.String(val)
}

// HostId returns an attribute KeyValue conforming to the "host.id"semantic conventions
// It represents the Unique host ID. For Cloud, this must be the instance_id assigned by the cloud provider. For non-containerized systems, this should be the `machine-id`. See the table below for the sources to use to determine the `machine-id` based on operating system
func HostId(val string) attribute.KeyValue {
	return HostIdKey.String(val)
}

// HostImageId returns an attribute KeyValue conforming to the "host.image.id"semantic conventions
// It represents the VM image ID or host OS image ID. For Cloud, this value is from the provider
func HostImageId(val string) attribute.KeyValue {
	return HostImageIdKey.String(val)
}

// HostImageName returns an attribute KeyValue conforming to the "host.image.name"semantic conventions
// It represents the Name of the VM image or OS install the host was instantiated from
func HostImageName(val string) attribute.KeyValue {
	return HostImageNameKey.String(val)
}

// HostImageVersion returns an attribute KeyValue conforming to the "host.image.version"semantic conventions
// It represents the The version string of the VM image or host OS as defined in [Version Attributes]
//
// [Version Attributes]: /docs/resource/README.md#version-attributes
func HostImageVersion(val string) attribute.KeyValue {
	return HostImageVersionKey.String(val)
}

// HostIp returns an attribute KeyValue conforming to the "host.ip"semantic conventions
// It represents the Available IP addresses of the host, excluding loopback interfaces
func HostIp(val []string) attribute.KeyValue {
	return HostIpKey.StringSlice(val)
}

// HostMac returns an attribute KeyValue conforming to the "host.mac"semantic conventions
// It represents the Available MAC addresses of the host, excluding loopback interfaces
func HostMac(val []string) attribute.KeyValue {
	return HostMacKey.StringSlice(val)
}

// HostName returns an attribute KeyValue conforming to the "host.name"semantic conventions
// It represents the Name of the host. On Unix systems, it may contain what the hostname command returns, or the fully qualified hostname, or another name specified by the user
func HostName(val string) attribute.KeyValue {
	return HostNameKey.String(val)
}

// HostType returns an attribute KeyValue conforming to the "host.type"semantic conventions
// It represents the Type of host. For Cloud, this must be the machine type
func HostType(val string) attribute.KeyValue {
	return HostTypeKey.String(val)
}

// Enum values for host.arch
var (

	// AMD64
	// Stability: experimental
	HostArchAmd64 = HostArchKey.String("amd64")
	// ARM32
	// Stability: experimental
	HostArchArm32 = HostArchKey.String("arm32")
	// ARM64
	// Stability: experimental
	HostArchArm64 = HostArchKey.String("arm64")
	// Itanium
	// Stability: experimental
	HostArchIa64 = HostArchKey.String("ia64")
	// 32-bit PowerPC
	// Stability: experimental
	HostArchPpc32 = HostArchKey.String("ppc32")
	// 64-bit PowerPC
	// Stability: experimental
	HostArchPpc64 = HostArchKey.String("ppc64")
	// IBM z/Architecture
	// Stability: experimental
	HostArchS390x = HostArchKey.String("s390x")
	// 32-bit x86
	// Stability: experimental
	HostArchX86 = HostArchKey.String("x86")
)

// Namespace: http
const (

	// Deprecated, use `client.address` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `client.address`.
	//
	// Examples: "83.164.160.102"
	HttpClientIpKey = attribute.Key("http.client_ip")
	// State of the HTTP connection in the HTTP connection pool.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "active",
	// "idle",
	HttpConnectionStateKey = attribute.Key("http.connection.state")
	// Deprecated, use `network.protocol.name` instead.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `network.protocol.name`.
	//
	// Examples: undefined
	HttpFlavorKey = attribute.Key("http.flavor")
	// Deprecated, use one of `server.address`, `client.address` or `http.request.header.host` instead, depending on the usage.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by one of `server.address`, `client.address` or `http.request.header.host`, depending on the usage.
	//
	// Examples:
	// "www.example.org",
	HttpHostKey = attribute.Key("http.host")
	// Deprecated, use `http.request.method` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `http.request.method`.
	//
	// Examples:
	// "GET",
	// "POST",
	// "HEAD",
	HttpMethodKey = attribute.Key("http.method")
	// The size of the request payload body in bytes. This is the number of bytes transferred excluding headers and is often, but not always, present as the [Content-Length] header. For requests using transport encoding, this should be the compressed size.
	//
	// Stability: Experimental
	// Type: int
	//
	// [Content-Length]: https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length
	HttpRequestBodySizeKey = attribute.Key("http.request.body.size")
	// HTTP request method.
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "GET",
	// "POST",
	// "HEAD",
	//
	// Note: HTTP request method value SHOULD be "known" to the instrumentation.
	// By default, this convention defines "known" methods as the ones listed in [RFC9110]
	// and the PATCH method defined in [RFC5789].
	//
	// If the HTTP request method is not known to instrumentation, it MUST set the `http.request.method` attribute to `_OTHER`.
	//
	// If the HTTP instrumentation could end up converting valid HTTP request methods to `_OTHER`, then it MUST provide a way to override
	// the list of known HTTP methods. If this override is done via environment variable, then the environment variable MUST be named
	// OTEL_INSTRUMENTATION_HTTP_KNOWN_METHODS and support a comma-separated list of case-sensitive known HTTP methods
	// (this list MUST be a full override of the default known method, it is not a list of known methods in addition to the defaults).
	//
	// HTTP method names are case-sensitive and `http.request.method` attribute value MUST match a known HTTP method name exactly.
	// Instrumentations for specific web frameworks that consider HTTP methods to be case insensitive, SHOULD populate a canonical equivalent.
	// Tracing instrumentations that do so, MUST also set `http.request.method_original` to the original value
	//
	// [RFC9110]: https://www.rfc-editor.org/rfc/rfc9110.html#name-methods
	// [RFC5789]: https://www.rfc-editor.org/rfc/rfc5789.html
	HttpRequestMethodKey = attribute.Key("http.request.method")
	// Original HTTP method sent by the client in the request line.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "GeT",
	// "ACL",
	// "foo",
	HttpRequestMethodOriginalKey = attribute.Key("http.request.method_original")
	// The ordinal number of request resending attempt (for any reason, including redirects).
	//
	// Stability: Stable
	// Type: int
	//
	// Note: The resend count SHOULD be updated each time an HTTP request gets resent by the client, regardless of what was the cause of the resending (e.g. redirection, authorization failure, 503 Server Unavailable, network issues, or any other)
	HttpRequestResendCountKey = attribute.Key("http.request.resend_count")
	// The total size of the request in bytes. This should be the total number of bytes sent over the wire, including the request line (HTTP/1.1), framing (HTTP/2 and HTTP/3), headers, and request body if any.
	//
	// Stability: Experimental
	// Type: int
	HttpRequestSizeKey = attribute.Key("http.request.size")
	// Deprecated, use `http.request.header.content-length` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `http.request.header.content-length`
	HttpRequestContentLengthKey = attribute.Key("http.request_content_length")
	// Deprecated, use `http.request.body.size` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `http.request.body.size`
	HttpRequestContentLengthUncompressedKey = attribute.Key("http.request_content_length_uncompressed")
	// The size of the response payload body in bytes. This is the number of bytes transferred excluding headers and is often, but not always, present as the [Content-Length] header. For requests using transport encoding, this should be the compressed size.
	//
	// Stability: Experimental
	// Type: int
	//
	// [Content-Length]: https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length
	HttpResponseBodySizeKey = attribute.Key("http.response.body.size")
	// The total size of the response in bytes. This should be the total number of bytes sent over the wire, including the status line (HTTP/1.1), framing (HTTP/2 and HTTP/3), headers, and response body and trailers if any.
	//
	// Stability: Experimental
	// Type: int
	HttpResponseSizeKey = attribute.Key("http.response.size")
	// [HTTP response status code].
	// Stability: Stable
	// Type: int
	//
	// Examples:
	// 200,
	//
	// [HTTP response status code]: https://tools.ietf.org/html/rfc7231#section-6
	HttpResponseStatusCodeKey = attribute.Key("http.response.status_code")
	// Deprecated, use `http.response.header.content-length` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `http.response.header.content-length`
	HttpResponseContentLengthKey = attribute.Key("http.response_content_length")
	// Deprecated, use `http.response.body.size` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replace by `http.response.body.size`
	HttpResponseContentLengthUncompressedKey = attribute.Key("http.response_content_length_uncompressed")
	// The matched route, that is, the path template in the format used by the respective server framework.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "/users/:userID?",
	// "{controller}/{action}/{id?}",
	//
	// Note: MUST NOT be populated when this is not supported by the HTTP server framework as the route attribute should have low-cardinality and the URI path can NOT substitute it.
	// SHOULD include the [application root] if there is one
	//
	// [application root]: /docs/http/http-spans.md#http-server-definitions
	HttpRouteKey = attribute.Key("http.route")
	// Deprecated, use `url.scheme` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `url.scheme` instead.
	//
	// Examples:
	// "http",
	// "https",
	HttpSchemeKey = attribute.Key("http.scheme")
	// Deprecated, use `server.address` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `server.address`.
	//
	// Examples:
	// "example.com",
	HttpServerNameKey = attribute.Key("http.server_name")
	// Deprecated, use `http.response.status_code` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `http.response.status_code`.
	//
	// Examples:
	// 200,
	HttpStatusCodeKey = attribute.Key("http.status_code")
	// Deprecated, use `url.path` and `url.query` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Split to `url.path` and `url.query.
	//
	// Examples:
	// "/search?q=OpenTelemetry#SemConv",
	HttpTargetKey = attribute.Key("http.target")
	// Deprecated, use `url.full` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `url.full`.
	//
	// Examples:
	// "https://www.foo.bar/search?q=OpenTelemetry#SemConv",
	HttpUrlKey = attribute.Key("http.url")
	// Deprecated, use `user_agent.original` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `user_agent.original`.
	//
	// Examples:
	// "CERN-LineMode/2.15 libwww/2.17b3",
	// "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Mobile/15E148 Safari/604.1",
	HttpUserAgentKey = attribute.Key("http.user_agent")
)

// HttpClientIp returns an attribute KeyValue conforming to the "http.client_ip"semantic conventions
// It represents the Deprecated, use `client.address` instead
func HttpClientIp(val string) attribute.KeyValue {
	return HttpClientIpKey.String(val)
}

// HttpConnectionState returns an attribute KeyValue conforming to the "http.connection.state"semantic conventions
// It represents the State of the HTTP connection in the HTTP connection pool
func HttpConnectionState(val string) attribute.KeyValue {
	return HttpConnectionStateKey.String(val)
}

// HttpFlavor returns an attribute KeyValue conforming to the "http.flavor"semantic conventions
// It represents the Deprecated, use `network.protocol.name` instead
func HttpFlavor(val string) attribute.KeyValue {
	return HttpFlavorKey.String(val)
}

// HttpHost returns an attribute KeyValue conforming to the "http.host"semantic conventions
// It represents the Deprecated, use one of `server.address`, `client.address` or `http.request.header.host` instead, depending on the usage
func HttpHost(val string) attribute.KeyValue {
	return HttpHostKey.String(val)
}

// HttpMethod returns an attribute KeyValue conforming to the "http.method"semantic conventions
// It represents the Deprecated, use `http.request.method` instead
func HttpMethod(val string) attribute.KeyValue {
	return HttpMethodKey.String(val)
}

// HttpRequestBodySize returns an attribute KeyValue conforming to the "http.request.body.size"semantic conventions
// It represents the The size of the request payload body in bytes. This is the number of bytes transferred excluding headers and is often, but not always, present as the [Content-Length] header. For requests using transport encoding, this should be the compressed size
//
// [Content-Length]: https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length
func HttpRequestBodySize(val int) attribute.KeyValue {
	return HttpRequestBodySizeKey.Int(val)
}

// HttpRequestMethod returns an attribute KeyValue conforming to the "http.request.method"semantic conventions
// It represents the HTTP request method
func HttpRequestMethod(val string) attribute.KeyValue {
	return HttpRequestMethodKey.String(val)
}

// HttpRequestMethodOriginal returns an attribute KeyValue conforming to the "http.request.method_original"semantic conventions
// It represents the Original HTTP method sent by the client in the request line
func HttpRequestMethodOriginal(val string) attribute.KeyValue {
	return HttpRequestMethodOriginalKey.String(val)
}

// HttpRequestResendCount returns an attribute KeyValue conforming to the "http.request.resend_count"semantic conventions
// It represents the The ordinal number of request resending attempt (for any reason, including redirects)
func HttpRequestResendCount(val int) attribute.KeyValue {
	return HttpRequestResendCountKey.Int(val)
}

// HttpRequestSize returns an attribute KeyValue conforming to the "http.request.size"semantic conventions
// It represents the The total size of the request in bytes. This should be the total number of bytes sent over the wire, including the request line (HTTP/1.1), framing (HTTP/2 and HTTP/3), headers, and request body if any
func HttpRequestSize(val int) attribute.KeyValue {
	return HttpRequestSizeKey.Int(val)
}

// HttpRequestContentLength returns an attribute KeyValue conforming to the "http.request_content_length"semantic conventions
// It represents the Deprecated, use `http.request.header.content-length` instead
func HttpRequestContentLength(val int) attribute.KeyValue {
	return HttpRequestContentLengthKey.Int(val)
}

// HttpRequestContentLengthUncompressed returns an attribute KeyValue conforming to the "http.request_content_length_uncompressed"semantic conventions
// It represents the Deprecated, use `http.request.body.size` instead
func HttpRequestContentLengthUncompressed(val int) attribute.KeyValue {
	return HttpRequestContentLengthUncompressedKey.Int(val)
}

// HttpResponseBodySize returns an attribute KeyValue conforming to the "http.response.body.size"semantic conventions
// It represents the The size of the response payload body in bytes. This is the number of bytes transferred excluding headers and is often, but not always, present as the [Content-Length] header. For requests using transport encoding, this should be the compressed size
//
// [Content-Length]: https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length
func HttpResponseBodySize(val int) attribute.KeyValue {
	return HttpResponseBodySizeKey.Int(val)
}

// HttpResponseSize returns an attribute KeyValue conforming to the "http.response.size"semantic conventions
// It represents the The total size of the response in bytes. This should be the total number of bytes sent over the wire, including the status line (HTTP/1.1), framing (HTTP/2 and HTTP/3), headers, and response body and trailers if any
func HttpResponseSize(val int) attribute.KeyValue {
	return HttpResponseSizeKey.Int(val)
}

// HttpResponseStatusCode returns an attribute KeyValue conforming to the "http.response.status_code"semantic conventions
// It represents the [HTTP response status code]
//
// [HTTP response status code]: https://tools.ietf.org/html/rfc7231#section-6
func HttpResponseStatusCode(val int) attribute.KeyValue {
	return HttpResponseStatusCodeKey.Int(val)
}

// HttpResponseContentLength returns an attribute KeyValue conforming to the "http.response_content_length"semantic conventions
// It represents the Deprecated, use `http.response.header.content-length` instead
func HttpResponseContentLength(val int) attribute.KeyValue {
	return HttpResponseContentLengthKey.Int(val)
}

// HttpResponseContentLengthUncompressed returns an attribute KeyValue conforming to the "http.response_content_length_uncompressed"semantic conventions
// It represents the Deprecated, use `http.response.body.size` instead
func HttpResponseContentLengthUncompressed(val int) attribute.KeyValue {
	return HttpResponseContentLengthUncompressedKey.Int(val)
}

// HttpRoute returns an attribute KeyValue conforming to the "http.route"semantic conventions
// It represents the The matched route, that is, the path template in the format used by the respective server framework
func HttpRoute(val string) attribute.KeyValue {
	return HttpRouteKey.String(val)
}

// HttpScheme returns an attribute KeyValue conforming to the "http.scheme"semantic conventions
// It represents the Deprecated, use `url.scheme` instead
func HttpScheme(val string) attribute.KeyValue {
	return HttpSchemeKey.String(val)
}

// HttpServerName returns an attribute KeyValue conforming to the "http.server_name"semantic conventions
// It represents the Deprecated, use `server.address` instead
func HttpServerName(val string) attribute.KeyValue {
	return HttpServerNameKey.String(val)
}

// HttpStatusCode returns an attribute KeyValue conforming to the "http.status_code"semantic conventions
// It represents the Deprecated, use `http.response.status_code` instead
func HttpStatusCode(val int) attribute.KeyValue {
	return HttpStatusCodeKey.Int(val)
}

// HttpTarget returns an attribute KeyValue conforming to the "http.target"semantic conventions
// It represents the Deprecated, use `url.path` and `url.query` instead
func HttpTarget(val string) attribute.KeyValue {
	return HttpTargetKey.String(val)
}

// HttpUrl returns an attribute KeyValue conforming to the "http.url"semantic conventions
// It represents the Deprecated, use `url.full` instead
func HttpUrl(val string) attribute.KeyValue {
	return HttpUrlKey.String(val)
}

// HttpUserAgent returns an attribute KeyValue conforming to the "http.user_agent"semantic conventions
// It represents the Deprecated, use `user_agent.original` instead
func HttpUserAgent(val string) attribute.KeyValue {
	return HttpUserAgentKey.String(val)
}

// Enum values for http.connection.state
var (

	// active state.
	// Stability: experimental
	HttpConnectionStateActive = HttpConnectionStateKey.String("active")
	// idle state.
	// Stability: experimental
	HttpConnectionStateIdle = HttpConnectionStateKey.String("idle")
)

// Enum values for http.flavor
var (

	// HTTP/1.0
	// Stability: experimental
	HttpFlavorHttp10 = HttpFlavorKey.String("1.0")
	// HTTP/1.1
	// Stability: experimental
	HttpFlavorHttp11 = HttpFlavorKey.String("1.1")
	// HTTP/2
	// Stability: experimental
	HttpFlavorHttp20 = HttpFlavorKey.String("2.0")
	// HTTP/3
	// Stability: experimental
	HttpFlavorHttp30 = HttpFlavorKey.String("3.0")
	// SPDY protocol.
	// Stability: experimental
	HttpFlavorSpdy = HttpFlavorKey.String("SPDY")
	// QUIC protocol.
	// Stability: experimental
	HttpFlavorQuic = HttpFlavorKey.String("QUIC")
)

// Enum values for http.request.method
var (

	// CONNECT method.
	// Stability: stable
	HttpRequestMethodConnect = HttpRequestMethodKey.String("CONNECT")
	// DELETE method.
	// Stability: stable
	HttpRequestMethodDelete = HttpRequestMethodKey.String("DELETE")
	// GET method.
	// Stability: stable
	HttpRequestMethodGet = HttpRequestMethodKey.String("GET")
	// HEAD method.
	// Stability: stable
	HttpRequestMethodHead = HttpRequestMethodKey.String("HEAD")
	// OPTIONS method.
	// Stability: stable
	HttpRequestMethodOptions = HttpRequestMethodKey.String("OPTIONS")
	// PATCH method.
	// Stability: stable
	HttpRequestMethodPatch = HttpRequestMethodKey.String("PATCH")
	// POST method.
	// Stability: stable
	HttpRequestMethodPost = HttpRequestMethodKey.String("POST")
	// PUT method.
	// Stability: stable
	HttpRequestMethodPut = HttpRequestMethodKey.String("PUT")
	// TRACE method.
	// Stability: stable
	HttpRequestMethodTrace = HttpRequestMethodKey.String("TRACE")
	// Any HTTP method that the instrumentation has no prior knowledge of.
	// Stability: stable
	HttpRequestMethodOther = HttpRequestMethodKey.String("_OTHER")
)

// Namespace: hw
const (

	// An identifier for the hardware component, unique within the monitored host
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "win32battery_battery_testsysa33_1",
	HwIdKey = attribute.Key("hw.id")
	// An easily-recognizable name for the hardware component
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "eth0",
	HwNameKey = attribute.Key("hw.name")
	// Unique identifier of the parent component (typically the `hw.id` attribute of the enclosure, or disk controller)
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "dellStorage_perc_0",
	HwParentKey = attribute.Key("hw.parent")
	// The current state of the component
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	HwStateKey = attribute.Key("hw.state")
	// Type of the component
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: Describes the category of the hardware component for which `hw.state` is being reported. For example, `hw.type=temperature` along with `hw.state=degraded` would indicate that the temperature of the hardware component has been reported as `degraded`
	HwTypeKey = attribute.Key("hw.type")
)

// HwId returns an attribute KeyValue conforming to the "hw.id"semantic conventions
// It represents the An identifier for the hardware component, unique within the monitored host
func HwId(val string) attribute.KeyValue {
	return HwIdKey.String(val)
}

// HwName returns an attribute KeyValue conforming to the "hw.name"semantic conventions
// It represents the An easily-recognizable name for the hardware component
func HwName(val string) attribute.KeyValue {
	return HwNameKey.String(val)
}

// HwParent returns an attribute KeyValue conforming to the "hw.parent"semantic conventions
// It represents the Unique identifier of the parent component (typically the `hw.id` attribute of the enclosure, or disk controller)
func HwParent(val string) attribute.KeyValue {
	return HwParentKey.String(val)
}

// HwState returns an attribute KeyValue conforming to the "hw.state"semantic conventions
// It represents the The current state of the component
func HwState(val string) attribute.KeyValue {
	return HwStateKey.String(val)
}

// HwType returns an attribute KeyValue conforming to the "hw.type"semantic conventions
// It represents the Type of the component
func HwType(val string) attribute.KeyValue {
	return HwTypeKey.String(val)
}

// Enum values for hw.state
var (

	// Ok
	// Stability: experimental
	HwStateOk = HwStateKey.String("ok")
	// Degraded
	// Stability: experimental
	HwStateDegraded = HwStateKey.String("degraded")
	// Failed
	// Stability: experimental
	HwStateFailed = HwStateKey.String("failed")
)

// Enum values for hw.type
var (

	// Battery
	// Stability: experimental
	HwTypeBattery = HwTypeKey.String("battery")
	// CPU
	// Stability: experimental
	HwTypeCpu = HwTypeKey.String("cpu")
	// Disk controller
	// Stability: experimental
	HwTypeDiskController = HwTypeKey.String("disk_controller")
	// Enclosure
	// Stability: experimental
	HwTypeEnclosure = HwTypeKey.String("enclosure")
	// Fan
	// Stability: experimental
	HwTypeFan = HwTypeKey.String("fan")
	// GPU
	// Stability: experimental
	HwTypeGpu = HwTypeKey.String("gpu")
	// Logical disk
	// Stability: experimental
	HwTypeLogicalDisk = HwTypeKey.String("logical_disk")
	// Memory
	// Stability: experimental
	HwTypeMemory = HwTypeKey.String("memory")
	// Network
	// Stability: experimental
	HwTypeNetwork = HwTypeKey.String("network")
	// Physical disk
	// Stability: experimental
	HwTypePhysicalDisk = HwTypeKey.String("physical_disk")
	// Power supply
	// Stability: experimental
	HwTypePowerSupply = HwTypeKey.String("power_supply")
	// Tape drive
	// Stability: experimental
	HwTypeTapeDrive = HwTypeKey.String("tape_drive")
	// Temperature
	// Stability: experimental
	HwTypeTemperature = HwTypeKey.String("temperature")
	// Voltage
	// Stability: experimental
	HwTypeVoltage = HwTypeKey.String("voltage")
)

// Namespace: ios
const (

	// Deprecated use the `device.app.lifecycle` event definition including `ios.state` as a payload field instead.
	//
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Moved to a payload field of `device.app.lifecycle`.
	//
	// Examples: undefined
	// Note: The iOS lifecycle states are defined in the [UIApplicationDelegate documentation], and from which the `OS terminology` column values are derived
	//
	// [UIApplicationDelegate documentation]: https://developer.apple.com/documentation/uikit/uiapplicationdelegate#1656902
	IosStateKey = attribute.Key("ios.state")
)

// IosState returns an attribute KeyValue conforming to the "ios.state"semantic conventions
// It represents the Deprecated use the `device.app.lifecycle` event definition including `ios.state` as a payload field instead
func IosState(val string) attribute.KeyValue {
	return IosStateKey.String(val)
}

// Enum values for ios.state
var (

	// The app has become `active`. Associated with UIKit notification `applicationDidBecomeActive`.
	//
	// Stability: experimental
	IosStateActive = IosStateKey.String("active")
	// The app is now `inactive`. Associated with UIKit notification `applicationWillResignActive`.
	//
	// Stability: experimental
	IosStateInactive = IosStateKey.String("inactive")
	// The app is now in the background. This value is associated with UIKit notification `applicationDidEnterBackground`.
	//
	// Stability: experimental
	IosStateBackground = IosStateKey.String("background")
	// The app is now in the foreground. This value is associated with UIKit notification `applicationWillEnterForeground`.
	//
	// Stability: experimental
	IosStateForeground = IosStateKey.String("foreground")
	// The app is about to terminate. Associated with UIKit notification `applicationWillTerminate`.
	//
	// Stability: experimental
	IosStateTerminate = IosStateKey.String("terminate")
)

// Namespace: jvm
const (

	// Name of the buffer pool.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "mapped",
	// "direct",
	//
	// Note: Pool names are generally obtained via [BufferPoolMXBean#getName()]
	//
	// [BufferPoolMXBean#getName()]: https://docs.oracle.com/en/java/javase/11/docs/api/java.management/java/lang/management/BufferPoolMXBean.html#getName()
	JvmBufferPoolNameKey = attribute.Key("jvm.buffer.pool.name")
	// Name of the garbage collector action.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "end of minor GC",
	// "end of major GC",
	//
	// Note: Garbage collector action is generally obtained via [GarbageCollectionNotificationInfo#getGcAction()]
	//
	// [GarbageCollectionNotificationInfo#getGcAction()]: https://docs.oracle.com/en/java/javase/11/docs/api/jdk.management/com/sun/management/GarbageCollectionNotificationInfo.html#getGcAction()
	JvmGcActionKey = attribute.Key("jvm.gc.action")
	// Name of the garbage collector.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "G1 Young Generation",
	// "G1 Old Generation",
	//
	// Note: Garbage collector name is generally obtained via [GarbageCollectionNotificationInfo#getGcName()]
	//
	// [GarbageCollectionNotificationInfo#getGcName()]: https://docs.oracle.com/en/java/javase/11/docs/api/jdk.management/com/sun/management/GarbageCollectionNotificationInfo.html#getGcName()
	JvmGcNameKey = attribute.Key("jvm.gc.name")
	// Name of the memory pool.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "G1 Old Gen",
	// "G1 Eden space",
	// "G1 Survivor Space",
	//
	// Note: Pool names are generally obtained via [MemoryPoolMXBean#getName()]
	//
	// [MemoryPoolMXBean#getName()]: https://docs.oracle.com/en/java/javase/11/docs/api/java.management/java/lang/management/MemoryPoolMXBean.html#getName()
	JvmMemoryPoolNameKey = attribute.Key("jvm.memory.pool.name")
	// The type of memory.
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "heap",
	// "non_heap",
	JvmMemoryTypeKey = attribute.Key("jvm.memory.type")
	// Whether the thread is daemon or not.
	// Stability: Stable
	// Type: boolean
	//
	// Examples: undefined
	JvmThreadDaemonKey = attribute.Key("jvm.thread.daemon")
	// State of the thread.
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "runnable",
	// "blocked",
	JvmThreadStateKey = attribute.Key("jvm.thread.state")
)

// JvmBufferPoolName returns an attribute KeyValue conforming to the "jvm.buffer.pool.name"semantic conventions
// It represents the Name of the buffer pool
func JvmBufferPoolName(val string) attribute.KeyValue {
	return JvmBufferPoolNameKey.String(val)
}

// JvmGcAction returns an attribute KeyValue conforming to the "jvm.gc.action"semantic conventions
// It represents the Name of the garbage collector action
func JvmGcAction(val string) attribute.KeyValue {
	return JvmGcActionKey.String(val)
}

// JvmGcName returns an attribute KeyValue conforming to the "jvm.gc.name"semantic conventions
// It represents the Name of the garbage collector
func JvmGcName(val string) attribute.KeyValue {
	return JvmGcNameKey.String(val)
}

// JvmMemoryPoolName returns an attribute KeyValue conforming to the "jvm.memory.pool.name"semantic conventions
// It represents the Name of the memory pool
func JvmMemoryPoolName(val string) attribute.KeyValue {
	return JvmMemoryPoolNameKey.String(val)
}

// JvmMemoryType returns an attribute KeyValue conforming to the "jvm.memory.type"semantic conventions
// It represents the The type of memory
func JvmMemoryType(val string) attribute.KeyValue {
	return JvmMemoryTypeKey.String(val)
}

// JvmThreadDaemon returns an attribute KeyValue conforming to the "jvm.thread.daemon"semantic conventions
// It represents the Whether the thread is daemon or not
func JvmThreadDaemon(val bool) attribute.KeyValue {
	return JvmThreadDaemonKey.Bool(val)
}

// JvmThreadState returns an attribute KeyValue conforming to the "jvm.thread.state"semantic conventions
// It represents the State of the thread
func JvmThreadState(val string) attribute.KeyValue {
	return JvmThreadStateKey.String(val)
}

// Enum values for jvm.memory.type
var (

	// Heap memory.
	// Stability: stable
	JvmMemoryTypeHeap = JvmMemoryTypeKey.String("heap")
	// Non-heap memory
	// Stability: stable
	JvmMemoryTypeNonHeap = JvmMemoryTypeKey.String("non_heap")
)

// Enum values for jvm.thread.state
var (

	// A thread that has not yet started is in this state.
	// Stability: stable
	JvmThreadStateNew = JvmThreadStateKey.String("new")
	// A thread executing in the Java virtual machine is in this state.
	// Stability: stable
	JvmThreadStateRunnable = JvmThreadStateKey.String("runnable")
	// A thread that is blocked waiting for a monitor lock is in this state.
	// Stability: stable
	JvmThreadStateBlocked = JvmThreadStateKey.String("blocked")
	// A thread that is waiting indefinitely for another thread to perform a particular action is in this state.
	// Stability: stable
	JvmThreadStateWaiting = JvmThreadStateKey.String("waiting")
	// A thread that is waiting for another thread to perform an action for up to a specified waiting time is in this state.
	// Stability: stable
	JvmThreadStateTimedWaiting = JvmThreadStateKey.String("timed_waiting")
	// A thread that has exited is in this state.
	// Stability: stable
	JvmThreadStateTerminated = JvmThreadStateKey.String("terminated")
)

// Namespace: k8s
const (

	// The name of the cluster.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry-cluster",
	K8sClusterNameKey = attribute.Key("k8s.cluster.name")
	// A pseudo-ID for the cluster, set to the UID of the `kube-system` namespace.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "218fc5a9-a5f1-4b54-aa05-46717d0ab26d",
	//
	// Note: K8s doesn't have support for obtaining a cluster ID. If this is ever
	// added, we will recommend collecting the `k8s.cluster.uid` through the
	// official APIs. In the meantime, we are able to use the `uid` of the
	// `kube-system` namespace as a proxy for cluster ID. Read on for the
	// rationale.
	//
	// Every object created in a K8s cluster is assigned a distinct UID. The
	// `kube-system` namespace is used by Kubernetes itself and will exist
	// for the lifetime of the cluster. Using the `uid` of the `kube-system`
	// namespace is a reasonable proxy for the K8s ClusterID as it will only
	// change if the cluster is rebuilt. Furthermore, Kubernetes UIDs are
	// UUIDs as standardized by
	// [ISO/IEC 9834-8 and ITU-T X.667].
	// Which states:
	//
	// > If generated according to one of the mechanisms defined in Rec.
	// > ITU-T X.667 | ISO/IEC 9834-8, a UUID is either guaranteed to be
	// > different from all other UUIDs generated before 3603 A.D., or is
	// > extremely likely to be different (depending on the mechanism chosen).
	//
	// Therefore, UIDs between clusters should be extremely unlikely to
	// conflict
	//
	// [ISO/IEC 9834-8 and ITU-T X.667]: https://www.itu.int/ITU-T/studygroups/com17/oid.html
	K8sClusterUidKey = attribute.Key("k8s.cluster.uid")
	// The name of the Container from Pod specification, must be unique within a Pod. Container runtime usually uses different globally unique name (`container.name`).
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "redis",
	K8sContainerNameKey = attribute.Key("k8s.container.name")
	// Number of times the container was restarted. This attribute can be used to identify a particular container (running or stopped) within a container spec.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples: undefined
	K8sContainerRestartCountKey = attribute.Key("k8s.container.restart_count")
	// Last terminated reason of the Container.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Evicted",
	// "Error",
	K8sContainerStatusLastTerminatedReasonKey = attribute.Key("k8s.container.status.last_terminated_reason")
	// The name of the CronJob.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry",
	K8sCronjobNameKey = attribute.Key("k8s.cronjob.name")
	// The UID of the CronJob.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "275ecb36-5aa8-4c2a-9c47-d8bb681b9aff",
	K8sCronjobUidKey = attribute.Key("k8s.cronjob.uid")
	// The name of the DaemonSet.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry",
	K8sDaemonsetNameKey = attribute.Key("k8s.daemonset.name")
	// The UID of the DaemonSet.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "275ecb36-5aa8-4c2a-9c47-d8bb681b9aff",
	K8sDaemonsetUidKey = attribute.Key("k8s.daemonset.uid")
	// The name of the Deployment.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry",
	K8sDeploymentNameKey = attribute.Key("k8s.deployment.name")
	// The UID of the Deployment.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "275ecb36-5aa8-4c2a-9c47-d8bb681b9aff",
	K8sDeploymentUidKey = attribute.Key("k8s.deployment.uid")
	// The name of the Job.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry",
	K8sJobNameKey = attribute.Key("k8s.job.name")
	// The UID of the Job.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "275ecb36-5aa8-4c2a-9c47-d8bb681b9aff",
	K8sJobUidKey = attribute.Key("k8s.job.uid")
	// The name of the namespace that the pod is running in.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "default",
	K8sNamespaceNameKey = attribute.Key("k8s.namespace.name")
	// The name of the Node.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "node-1",
	K8sNodeNameKey = attribute.Key("k8s.node.name")
	// The UID of the Node.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "1eb3a0c6-0477-4080-a9cb-0cb7db65c6a2",
	K8sNodeUidKey = attribute.Key("k8s.node.uid")
	// The name of the Pod.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry-pod-autoconf",
	K8sPodNameKey = attribute.Key("k8s.pod.name")
	// The UID of the Pod.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "275ecb36-5aa8-4c2a-9c47-d8bb681b9aff",
	K8sPodUidKey = attribute.Key("k8s.pod.uid")
	// The name of the ReplicaSet.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry",
	K8sReplicasetNameKey = attribute.Key("k8s.replicaset.name")
	// The UID of the ReplicaSet.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "275ecb36-5aa8-4c2a-9c47-d8bb681b9aff",
	K8sReplicasetUidKey = attribute.Key("k8s.replicaset.uid")
	// The name of the StatefulSet.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "opentelemetry",
	K8sStatefulsetNameKey = attribute.Key("k8s.statefulset.name")
	// The UID of the StatefulSet.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "275ecb36-5aa8-4c2a-9c47-d8bb681b9aff",
	K8sStatefulsetUidKey = attribute.Key("k8s.statefulset.uid")
	// The name of the K8s volume.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "volume0",
	K8sVolumeNameKey = attribute.Key("k8s.volume.name")
	// The type of the K8s volume.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "emptyDir",
	// "persistentVolumeClaim",
	K8sVolumeTypeKey = attribute.Key("k8s.volume.type")
)

// K8sClusterName returns an attribute KeyValue conforming to the "k8s.cluster.name"semantic conventions
// It represents the The name of the cluster
func K8sClusterName(val string) attribute.KeyValue {
	return K8sClusterNameKey.String(val)
}

// K8sClusterUid returns an attribute KeyValue conforming to the "k8s.cluster.uid"semantic conventions
// It represents the A pseudo-ID for the cluster, set to the UID of the `kube-system` namespace
func K8sClusterUid(val string) attribute.KeyValue {
	return K8sClusterUidKey.String(val)
}

// K8sContainerName returns an attribute KeyValue conforming to the "k8s.container.name"semantic conventions
// It represents the The name of the Container from Pod specification, must be unique within a Pod. Container runtime usually uses different globally unique name (`container.name`)
func K8sContainerName(val string) attribute.KeyValue {
	return K8sContainerNameKey.String(val)
}

// K8sContainerRestartCount returns an attribute KeyValue conforming to the "k8s.container.restart_count"semantic conventions
// It represents the Number of times the container was restarted. This attribute can be used to identify a particular container (running or stopped) within a container spec
func K8sContainerRestartCount(val int) attribute.KeyValue {
	return K8sContainerRestartCountKey.Int(val)
}

// K8sContainerStatusLastTerminatedReason returns an attribute KeyValue conforming to the "k8s.container.status.last_terminated_reason"semantic conventions
// It represents the Last terminated reason of the Container
func K8sContainerStatusLastTerminatedReason(val string) attribute.KeyValue {
	return K8sContainerStatusLastTerminatedReasonKey.String(val)
}

// K8sCronjobName returns an attribute KeyValue conforming to the "k8s.cronjob.name"semantic conventions
// It represents the The name of the CronJob
func K8sCronjobName(val string) attribute.KeyValue {
	return K8sCronjobNameKey.String(val)
}

// K8sCronjobUid returns an attribute KeyValue conforming to the "k8s.cronjob.uid"semantic conventions
// It represents the The UID of the CronJob
func K8sCronjobUid(val string) attribute.KeyValue {
	return K8sCronjobUidKey.String(val)
}

// K8sDaemonsetName returns an attribute KeyValue conforming to the "k8s.daemonset.name"semantic conventions
// It represents the The name of the DaemonSet
func K8sDaemonsetName(val string) attribute.KeyValue {
	return K8sDaemonsetNameKey.String(val)
}

// K8sDaemonsetUid returns an attribute KeyValue conforming to the "k8s.daemonset.uid"semantic conventions
// It represents the The UID of the DaemonSet
func K8sDaemonsetUid(val string) attribute.KeyValue {
	return K8sDaemonsetUidKey.String(val)
}

// K8sDeploymentName returns an attribute KeyValue conforming to the "k8s.deployment.name"semantic conventions
// It represents the The name of the Deployment
func K8sDeploymentName(val string) attribute.KeyValue {
	return K8sDeploymentNameKey.String(val)
}

// K8sDeploymentUid returns an attribute KeyValue conforming to the "k8s.deployment.uid"semantic conventions
// It represents the The UID of the Deployment
func K8sDeploymentUid(val string) attribute.KeyValue {
	return K8sDeploymentUidKey.String(val)
}

// K8sJobName returns an attribute KeyValue conforming to the "k8s.job.name"semantic conventions
// It represents the The name of the Job
func K8sJobName(val string) attribute.KeyValue {
	return K8sJobNameKey.String(val)
}

// K8sJobUid returns an attribute KeyValue conforming to the "k8s.job.uid"semantic conventions
// It represents the The UID of the Job
func K8sJobUid(val string) attribute.KeyValue {
	return K8sJobUidKey.String(val)
}

// K8sNamespaceName returns an attribute KeyValue conforming to the "k8s.namespace.name"semantic conventions
// It represents the The name of the namespace that the pod is running in
func K8sNamespaceName(val string) attribute.KeyValue {
	return K8sNamespaceNameKey.String(val)
}

// K8sNodeName returns an attribute KeyValue conforming to the "k8s.node.name"semantic conventions
// It represents the The name of the Node
func K8sNodeName(val string) attribute.KeyValue {
	return K8sNodeNameKey.String(val)
}

// K8sNodeUid returns an attribute KeyValue conforming to the "k8s.node.uid"semantic conventions
// It represents the The UID of the Node
func K8sNodeUid(val string) attribute.KeyValue {
	return K8sNodeUidKey.String(val)
}

// K8sPodName returns an attribute KeyValue conforming to the "k8s.pod.name"semantic conventions
// It represents the The name of the Pod
func K8sPodName(val string) attribute.KeyValue {
	return K8sPodNameKey.String(val)
}

// K8sPodUid returns an attribute KeyValue conforming to the "k8s.pod.uid"semantic conventions
// It represents the The UID of the Pod
func K8sPodUid(val string) attribute.KeyValue {
	return K8sPodUidKey.String(val)
}

// K8sReplicasetName returns an attribute KeyValue conforming to the "k8s.replicaset.name"semantic conventions
// It represents the The name of the ReplicaSet
func K8sReplicasetName(val string) attribute.KeyValue {
	return K8sReplicasetNameKey.String(val)
}

// K8sReplicasetUid returns an attribute KeyValue conforming to the "k8s.replicaset.uid"semantic conventions
// It represents the The UID of the ReplicaSet
func K8sReplicasetUid(val string) attribute.KeyValue {
	return K8sReplicasetUidKey.String(val)
}

// K8sStatefulsetName returns an attribute KeyValue conforming to the "k8s.statefulset.name"semantic conventions
// It represents the The name of the StatefulSet
func K8sStatefulsetName(val string) attribute.KeyValue {
	return K8sStatefulsetNameKey.String(val)
}

// K8sStatefulsetUid returns an attribute KeyValue conforming to the "k8s.statefulset.uid"semantic conventions
// It represents the The UID of the StatefulSet
func K8sStatefulsetUid(val string) attribute.KeyValue {
	return K8sStatefulsetUidKey.String(val)
}

// K8sVolumeName returns an attribute KeyValue conforming to the "k8s.volume.name"semantic conventions
// It represents the The name of the K8s volume
func K8sVolumeName(val string) attribute.KeyValue {
	return K8sVolumeNameKey.String(val)
}

// K8sVolumeType returns an attribute KeyValue conforming to the "k8s.volume.type"semantic conventions
// It represents the The type of the K8s volume
func K8sVolumeType(val string) attribute.KeyValue {
	return K8sVolumeTypeKey.String(val)
}

// Enum values for k8s.volume.type
var (

	// A [persistentVolumeClaim] volume
	// Stability: experimental
	//
	// [persistentVolumeClaim]: https://v1-29.docs.kubernetes.io/docs/concepts/storage/volumes/#persistentvolumeclaim
	K8sVolumeTypePersistentVolumeClaim = K8sVolumeTypeKey.String("persistentVolumeClaim")
	// A [configMap] volume
	// Stability: experimental
	//
	// [configMap]: https://v1-29.docs.kubernetes.io/docs/concepts/storage/volumes/#configmap
	K8sVolumeTypeConfigMap = K8sVolumeTypeKey.String("configMap")
	// A [downwardAPI] volume
	// Stability: experimental
	//
	// [downwardAPI]: https://v1-29.docs.kubernetes.io/docs/concepts/storage/volumes/#downwardapi
	K8sVolumeTypeDownwardApi = K8sVolumeTypeKey.String("downwardAPI")
	// An [emptyDir] volume
	// Stability: experimental
	//
	// [emptyDir]: https://v1-29.docs.kubernetes.io/docs/concepts/storage/volumes/#emptydir
	K8sVolumeTypeEmptyDir = K8sVolumeTypeKey.String("emptyDir")
	// A [secret] volume
	// Stability: experimental
	//
	// [secret]: https://v1-29.docs.kubernetes.io/docs/concepts/storage/volumes/#secret
	K8sVolumeTypeSecret = K8sVolumeTypeKey.String("secret")
	// A [local] volume
	// Stability: experimental
	//
	// [local]: https://v1-29.docs.kubernetes.io/docs/concepts/storage/volumes/#local
	K8sVolumeTypeLocal = K8sVolumeTypeKey.String("local")
)

// Namespace: linux
const (

	// The Linux Slab memory state
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "reclaimable",
	// "unreclaimable",
	LinuxMemorySlabStateKey = attribute.Key("linux.memory.slab.state")
)

// LinuxMemorySlabState returns an attribute KeyValue conforming to the "linux.memory.slab.state"semantic conventions
// It represents the The Linux Slab memory state
func LinuxMemorySlabState(val string) attribute.KeyValue {
	return LinuxMemorySlabStateKey.String(val)
}

// Enum values for linux.memory.slab.state
var (

	// reclaimable
	// Stability: experimental
	LinuxMemorySlabStateReclaimable = LinuxMemorySlabStateKey.String("reclaimable")
	// unreclaimable
	// Stability: experimental
	LinuxMemorySlabStateUnreclaimable = LinuxMemorySlabStateKey.String("unreclaimable")
)

// Namespace: log
const (

	// The basename of the file.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "audit.log",
	LogFileNameKey = attribute.Key("log.file.name")
	// The basename of the file, with symlinks resolved.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "uuid.log",
	LogFileNameResolvedKey = attribute.Key("log.file.name_resolved")
	// The full path to the file.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/var/log/mysql/audit.log",
	LogFilePathKey = attribute.Key("log.file.path")
	// The full path to the file, with symlinks resolved.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/var/lib/docker/uuid.log",
	LogFilePathResolvedKey = attribute.Key("log.file.path_resolved")
	// The stream associated with the log. See below for a list of well-known values.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	LogIostreamKey = attribute.Key("log.iostream")
	// The complete original Log Record.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "77 <86>1 2015-08-06T21:58:59.694Z 192.168.2.133 inactive - - - Something happened",
	// "[INFO] 8/3/24 12:34:56 Something happened",
	//
	// Note: This value MAY be added when processing a Log Record which was originally transmitted as a string or equivalent data type AND the Body field of the Log Record does not contain the same value. (e.g. a syslog or a log record read from a file.)
	LogRecordOriginalKey = attribute.Key("log.record.original")
	// A unique identifier for the Log Record.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "01ARZ3NDEKTSV4RRFFQ69G5FAV",
	//
	// Note: If an id is provided, other log records with the same id will be considered duplicates and can be removed safely. This means, that two distinguishable log records MUST have different values.
	// The id MAY be an [Universally Unique Lexicographically Sortable Identifier (ULID)], but other identifiers (e.g. UUID) may be used as needed
	//
	// [Universally Unique Lexicographically Sortable Identifier (ULID)]: https://github.com/ulid/spec
	LogRecordUidKey = attribute.Key("log.record.uid")
)

// LogFileName returns an attribute KeyValue conforming to the "log.file.name"semantic conventions
// It represents the The basename of the file
func LogFileName(val string) attribute.KeyValue {
	return LogFileNameKey.String(val)
}

// LogFileNameResolved returns an attribute KeyValue conforming to the "log.file.name_resolved"semantic conventions
// It represents the The basename of the file, with symlinks resolved
func LogFileNameResolved(val string) attribute.KeyValue {
	return LogFileNameResolvedKey.String(val)
}

// LogFilePath returns an attribute KeyValue conforming to the "log.file.path"semantic conventions
// It represents the The full path to the file
func LogFilePath(val string) attribute.KeyValue {
	return LogFilePathKey.String(val)
}

// LogFilePathResolved returns an attribute KeyValue conforming to the "log.file.path_resolved"semantic conventions
// It represents the The full path to the file, with symlinks resolved
func LogFilePathResolved(val string) attribute.KeyValue {
	return LogFilePathResolvedKey.String(val)
}

// LogIostream returns an attribute KeyValue conforming to the "log.iostream"semantic conventions
// It represents the The stream associated with the log. See below for a list of well-known values
func LogIostream(val string) attribute.KeyValue {
	return LogIostreamKey.String(val)
}

// LogRecordOriginal returns an attribute KeyValue conforming to the "log.record.original"semantic conventions
// It represents the The complete original Log Record
func LogRecordOriginal(val string) attribute.KeyValue {
	return LogRecordOriginalKey.String(val)
}

// LogRecordUid returns an attribute KeyValue conforming to the "log.record.uid"semantic conventions
// It represents the A unique identifier for the Log Record
func LogRecordUid(val string) attribute.KeyValue {
	return LogRecordUidKey.String(val)
}

// Enum values for log.iostream
var (

	// Logs from stdout stream
	// Stability: experimental
	LogIostreamStdout = LogIostreamKey.String("stdout")
	// Events from stderr stream
	// Stability: experimental
	LogIostreamStderr = LogIostreamKey.String("stderr")
)

// Namespace: message
const (

	// Deprecated, use `rpc.message.compressed_size` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `rpc.message.compressed_size`.
	//
	// Examples: undefined
	MessageCompressedSizeKey = attribute.Key("message.compressed_size")
	// Deprecated, use `rpc.message.id` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `rpc.message.id`.
	//
	// Examples: undefined
	MessageIdKey = attribute.Key("message.id")
	// Deprecated, use `rpc.message.type` instead.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `rpc.message.type`.
	//
	// Examples: undefined
	MessageTypeKey = attribute.Key("message.type")
	// Deprecated, use `rpc.message.uncompressed_size` instead.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `rpc.message.uncompressed_size`.
	//
	// Examples: undefined
	MessageUncompressedSizeKey = attribute.Key("message.uncompressed_size")
)

// MessageCompressedSize returns an attribute KeyValue conforming to the "message.compressed_size"semantic conventions
// It represents the Deprecated, use `rpc.message.compressed_size` instead
func MessageCompressedSize(val int) attribute.KeyValue {
	return MessageCompressedSizeKey.Int(val)
}

// MessageId returns an attribute KeyValue conforming to the "message.id"semantic conventions
// It represents the Deprecated, use `rpc.message.id` instead
func MessageId(val int) attribute.KeyValue {
	return MessageIdKey.Int(val)
}

// MessageType returns an attribute KeyValue conforming to the "message.type"semantic conventions
// It represents the Deprecated, use `rpc.message.type` instead
func MessageType(val string) attribute.KeyValue {
	return MessageTypeKey.String(val)
}

// MessageUncompressedSize returns an attribute KeyValue conforming to the "message.uncompressed_size"semantic conventions
// It represents the Deprecated, use `rpc.message.uncompressed_size` instead
func MessageUncompressedSize(val int) attribute.KeyValue {
	return MessageUncompressedSizeKey.Int(val)
}

// Enum values for message.type
var (

	// sent
	// Stability: experimental
	MessageTypeSent = MessageTypeKey.String("SENT")
	// received
	// Stability: experimental
	MessageTypeReceived = MessageTypeKey.String("RECEIVED")
)

// Namespace: messaging
const (

	// The number of messages sent, received, or processed in the scope of the batching operation.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 0,
	// 1,
	// 2,
	//
	// Note: Instrumentations SHOULD NOT set `messaging.batch.message_count` on spans that operate with a single message. When a messaging client library supports both batch and single-message API for the same operation, instrumentations SHOULD use `messaging.batch.message_count` for batching APIs and SHOULD NOT use it for single-message APIs
	MessagingBatchMessageCountKey = attribute.Key("messaging.batch.message_count")
	// A unique identifier for the client that consumes or produces a message.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "client-5",
	// "myhost@8742@s8083jm",
	MessagingClientIdKey = attribute.Key("messaging.client.id")
	// The name of the consumer group with which a consumer is associated.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "my-group",
	// "indexer",
	//
	// Note: Semantic conventions for individual messaging systems SHOULD document whether `messaging.consumer.group.name` is applicable and what it means in the context of that system
	MessagingConsumerGroupNameKey = attribute.Key("messaging.consumer.group.name")
	// A boolean that is true if the message destination is anonymous (could be unnamed or have auto-generated name).
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	MessagingDestinationAnonymousKey = attribute.Key("messaging.destination.anonymous")
	// The message destination name
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "MyQueue",
	// "MyTopic",
	//
	// Note: Destination name SHOULD uniquely identify a specific queue, topic or other entity within the broker. If
	// the broker doesn't have such notion, the destination name SHOULD uniquely identify the broker
	MessagingDestinationNameKey = attribute.Key("messaging.destination.name")
	// The identifier of the partition messages are sent to or received from, unique within the `messaging.destination.name`.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "1"
	MessagingDestinationPartitionIdKey = attribute.Key("messaging.destination.partition.id")
	// The name of the destination subscription from which a message is consumed.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "subscription-a",
	//
	// Note: Semantic conventions for individual messaging systems SHOULD document whether `messaging.destination.subscription.name` is applicable and what it means in the context of that system
	MessagingDestinationSubscriptionNameKey = attribute.Key("messaging.destination.subscription.name")
	// Low cardinality representation of the messaging destination name
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/customers/{customerId}",
	//
	// Note: Destination names could be constructed from templates. An example would be a destination name involving a user name or product id. Although the destination name in this case is of high cardinality, the underlying template is of low cardinality and can be effectively used for grouping and aggregation
	MessagingDestinationTemplateKey = attribute.Key("messaging.destination.template")
	// A boolean that is true if the message destination is temporary and might not exist anymore after messages are processed.
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	MessagingDestinationTemporaryKey = attribute.Key("messaging.destination.temporary")
	// Deprecated, no replacement at this time.
	// Stability: Experimental
	// Type: boolean
	// Deprecated: No replacement at this time.
	//
	// Examples: undefined
	MessagingDestinationPublishAnonymousKey = attribute.Key("messaging.destination_publish.anonymous")
	// Deprecated, no replacement at this time.
	// Stability: Experimental
	// Type: string
	// Deprecated: No replacement at this time.
	//
	// Examples:
	// "MyQueue",
	// "MyTopic",
	MessagingDestinationPublishNameKey = attribute.Key("messaging.destination_publish.name")
	// Deprecated, use `messaging.consumer.group.name` instead.
	//
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `messaging.consumer.group.name`.
	//
	// Examples: "$Default"
	MessagingEventhubsConsumerGroupKey = attribute.Key("messaging.eventhubs.consumer.group")
	// The UTC epoch seconds at which the message has been accepted and stored in the entity.
	//
	// Stability: Experimental
	// Type: int
	MessagingEventhubsMessageEnqueuedTimeKey = attribute.Key("messaging.eventhubs.message.enqueued_time")
	// The ack deadline in seconds set for the modify ack deadline request.
	//
	// Stability: Experimental
	// Type: int
	MessagingGcpPubsubMessageAckDeadlineKey = attribute.Key("messaging.gcp_pubsub.message.ack_deadline")
	// The ack id for a given message.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "ack_id"
	MessagingGcpPubsubMessageAckIdKey = attribute.Key("messaging.gcp_pubsub.message.ack_id")
	// The delivery attempt for a given message.
	//
	// Stability: Experimental
	// Type: int
	MessagingGcpPubsubMessageDeliveryAttemptKey = attribute.Key("messaging.gcp_pubsub.message.delivery_attempt")
	// The ordering key for a given message. If the attribute is not present, the message does not have an ordering key.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "ordering_key"
	MessagingGcpPubsubMessageOrderingKeyKey = attribute.Key("messaging.gcp_pubsub.message.ordering_key")
	// Deprecated, use `messaging.consumer.group.name` instead.
	//
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `messaging.consumer.group.name`.
	//
	// Examples: "my-group"
	MessagingKafkaConsumerGroupKey = attribute.Key("messaging.kafka.consumer.group")
	// Deprecated, use `messaging.destination.partition.id` instead.
	//
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `messaging.destination.partition.id`
	MessagingKafkaDestinationPartitionKey = attribute.Key("messaging.kafka.destination.partition")
	// Message keys in Kafka are used for grouping alike messages to ensure they're processed on the same partition. They differ from `messaging.message.id` in that they're not unique. If the key is `null`, the attribute MUST NOT be set.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "myKey"
	// Note: If the key type is not string, it's string representation has to be supplied for the attribute. If the key has no unambiguous, canonical string form, don't include its value
	MessagingKafkaMessageKeyKey = attribute.Key("messaging.kafka.message.key")
	// Deprecated, use `messaging.kafka.offset` instead.
	//
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `messaging.kafka.offset`
	MessagingKafkaMessageOffsetKey = attribute.Key("messaging.kafka.message.offset")
	// A boolean that is true if the message is a tombstone.
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	MessagingKafkaMessageTombstoneKey = attribute.Key("messaging.kafka.message.tombstone")
	// The offset of a record in the corresponding Kafka partition.
	//
	// Stability: Experimental
	// Type: int
	MessagingKafkaOffsetKey = attribute.Key("messaging.kafka.offset")
	// The size of the message body in bytes.
	//
	// Stability: Experimental
	// Type: int
	//
	// Note: This can refer to both the compressed or uncompressed body size. If both sizes are known, the uncompressed
	// body size should be used
	MessagingMessageBodySizeKey = attribute.Key("messaging.message.body.size")
	// The conversation ID identifying the conversation to which the message belongs, represented as a string. Sometimes called "Correlation ID".
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "MyConversationId"
	MessagingMessageConversationIdKey = attribute.Key("messaging.message.conversation_id")
	// The size of the message body and metadata in bytes.
	//
	// Stability: Experimental
	// Type: int
	//
	// Note: This can refer to both the compressed or uncompressed size. If both sizes are known, the uncompressed
	// size should be used
	MessagingMessageEnvelopeSizeKey = attribute.Key("messaging.message.envelope.size")
	// A value used by the messaging system as an identifier for the message, represented as a string.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "452a7c7c7c7048c2f887f61572b18fc2"
	MessagingMessageIdKey = attribute.Key("messaging.message.id")
	// Deprecated, use `messaging.operation.type` instead.
	//
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `messaging.operation.type`.
	//
	// Examples:
	// "publish",
	// "create",
	// "process",
	MessagingOperationKey = attribute.Key("messaging.operation")
	// The system-specific name of the messaging operation.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "ack",
	// "nack",
	// "send",
	MessagingOperationNameKey = attribute.Key("messaging.operation.name")
	// A string identifying the type of the messaging operation.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: If a custom value is used, it MUST be of low cardinality
	MessagingOperationTypeKey = attribute.Key("messaging.operation.type")
	// RabbitMQ message routing key.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "myKey"
	MessagingRabbitmqDestinationRoutingKeyKey = attribute.Key("messaging.rabbitmq.destination.routing_key")
	// RabbitMQ message delivery tag
	//
	// Stability: Experimental
	// Type: int
	MessagingRabbitmqMessageDeliveryTagKey = attribute.Key("messaging.rabbitmq.message.delivery_tag")
	// Deprecated, use `messaging.consumer.group.name` instead.
	//
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `messaging.consumer.group.name` on the consumer spans. No replacement for producer spans.
	//
	// Examples: "myConsumerGroup"
	MessagingRocketmqClientGroupKey = attribute.Key("messaging.rocketmq.client_group")
	// Model of message consumption. This only applies to consumer spans.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	MessagingRocketmqConsumptionModelKey = attribute.Key("messaging.rocketmq.consumption_model")
	// The delay time level for delay message, which determines the message delay time.
	//
	// Stability: Experimental
	// Type: int
	MessagingRocketmqMessageDelayTimeLevelKey = attribute.Key("messaging.rocketmq.message.delay_time_level")
	// The timestamp in milliseconds that the delay message is expected to be delivered to consumer.
	//
	// Stability: Experimental
	// Type: int
	MessagingRocketmqMessageDeliveryTimestampKey = attribute.Key("messaging.rocketmq.message.delivery_timestamp")
	// It is essential for FIFO message. Messages that belong to the same message group are always processed one by one within the same consumer group.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "myMessageGroup"
	MessagingRocketmqMessageGroupKey = attribute.Key("messaging.rocketmq.message.group")
	// Key(s) of message, another way to mark message besides message id.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "keyA",
	// "keyB",
	// ],
	MessagingRocketmqMessageKeysKey = attribute.Key("messaging.rocketmq.message.keys")
	// The secondary classifier of message besides topic.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "tagA"
	MessagingRocketmqMessageTagKey = attribute.Key("messaging.rocketmq.message.tag")
	// Type of message.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	MessagingRocketmqMessageTypeKey = attribute.Key("messaging.rocketmq.message.type")
	// Namespace of RocketMQ resources, resources in different namespaces are individual.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "myNamespace"
	MessagingRocketmqNamespaceKey = attribute.Key("messaging.rocketmq.namespace")
	// Deprecated, use `messaging.servicebus.destination.subscription_name` instead.
	//
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `messaging.servicebus.destination.subscription_name`.
	//
	// Examples: "subscription-a"
	MessagingServicebusDestinationSubscriptionNameKey = attribute.Key("messaging.servicebus.destination.subscription_name")
	// Describes the [settlement type].
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	//
	// [settlement type]: https://learn.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#peeklock
	MessagingServicebusDispositionStatusKey = attribute.Key("messaging.servicebus.disposition_status")
	// Number of deliveries that have been attempted for this message.
	//
	// Stability: Experimental
	// Type: int
	MessagingServicebusMessageDeliveryCountKey = attribute.Key("messaging.servicebus.message.delivery_count")
	// The UTC epoch seconds at which the message has been accepted and stored in the entity.
	//
	// Stability: Experimental
	// Type: int
	MessagingServicebusMessageEnqueuedTimeKey = attribute.Key("messaging.servicebus.message.enqueued_time")
	// The messaging system as identified by the client instrumentation.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: The actual messaging system may differ from the one known by the client. For example, when using Kafka client libraries to communicate with Azure Event Hubs, the `messaging.system` is set to `kafka` based on the instrumentation's best knowledge
	MessagingSystemKey = attribute.Key("messaging.system")
)

// MessagingBatchMessageCount returns an attribute KeyValue conforming to the "messaging.batch.message_count"semantic conventions
// It represents the The number of messages sent, received, or processed in the scope of the batching operation
func MessagingBatchMessageCount(val int) attribute.KeyValue {
	return MessagingBatchMessageCountKey.Int(val)
}

// MessagingClientId returns an attribute KeyValue conforming to the "messaging.client.id"semantic conventions
// It represents the A unique identifier for the client that consumes or produces a message
func MessagingClientId(val string) attribute.KeyValue {
	return MessagingClientIdKey.String(val)
}

// MessagingConsumerGroupName returns an attribute KeyValue conforming to the "messaging.consumer.group.name"semantic conventions
// It represents the The name of the consumer group with which a consumer is associated
func MessagingConsumerGroupName(val string) attribute.KeyValue {
	return MessagingConsumerGroupNameKey.String(val)
}

// MessagingDestinationAnonymous returns an attribute KeyValue conforming to the "messaging.destination.anonymous"semantic conventions
// It represents the A boolean that is true if the message destination is anonymous (could be unnamed or have auto-generated name)
func MessagingDestinationAnonymous(val bool) attribute.KeyValue {
	return MessagingDestinationAnonymousKey.Bool(val)
}

// MessagingDestinationName returns an attribute KeyValue conforming to the "messaging.destination.name"semantic conventions
// It represents the The message destination name
func MessagingDestinationName(val string) attribute.KeyValue {
	return MessagingDestinationNameKey.String(val)
}

// MessagingDestinationPartitionId returns an attribute KeyValue conforming to the "messaging.destination.partition.id"semantic conventions
// It represents the The identifier of the partition messages are sent to or received from, unique within the `messaging.destination.name`
func MessagingDestinationPartitionId(val string) attribute.KeyValue {
	return MessagingDestinationPartitionIdKey.String(val)
}

// MessagingDestinationSubscriptionName returns an attribute KeyValue conforming to the "messaging.destination.subscription.name"semantic conventions
// It represents the The name of the destination subscription from which a message is consumed
func MessagingDestinationSubscriptionName(val string) attribute.KeyValue {
	return MessagingDestinationSubscriptionNameKey.String(val)
}

// MessagingDestinationTemplate returns an attribute KeyValue conforming to the "messaging.destination.template"semantic conventions
// It represents the Low cardinality representation of the messaging destination name
func MessagingDestinationTemplate(val string) attribute.KeyValue {
	return MessagingDestinationTemplateKey.String(val)
}

// MessagingDestinationTemporary returns an attribute KeyValue conforming to the "messaging.destination.temporary"semantic conventions
// It represents the A boolean that is true if the message destination is temporary and might not exist anymore after messages are processed
func MessagingDestinationTemporary(val bool) attribute.KeyValue {
	return MessagingDestinationTemporaryKey.Bool(val)
}

// MessagingDestinationPublishAnonymous returns an attribute KeyValue conforming to the "messaging.destination_publish.anonymous"semantic conventions
// It represents the Deprecated, no replacement at this time
func MessagingDestinationPublishAnonymous(val bool) attribute.KeyValue {
	return MessagingDestinationPublishAnonymousKey.Bool(val)
}

// MessagingDestinationPublishName returns an attribute KeyValue conforming to the "messaging.destination_publish.name"semantic conventions
// It represents the Deprecated, no replacement at this time
func MessagingDestinationPublishName(val string) attribute.KeyValue {
	return MessagingDestinationPublishNameKey.String(val)
}

// MessagingEventhubsConsumerGroup returns an attribute KeyValue conforming to the "messaging.eventhubs.consumer.group"semantic conventions
// It represents the Deprecated, use `messaging.consumer.group.name` instead
func MessagingEventhubsConsumerGroup(val string) attribute.KeyValue {
	return MessagingEventhubsConsumerGroupKey.String(val)
}

// MessagingEventhubsMessageEnqueuedTime returns an attribute KeyValue conforming to the "messaging.eventhubs.message.enqueued_time"semantic conventions
// It represents the The UTC epoch seconds at which the message has been accepted and stored in the entity
func MessagingEventhubsMessageEnqueuedTime(val int) attribute.KeyValue {
	return MessagingEventhubsMessageEnqueuedTimeKey.Int(val)
}

// MessagingGcpPubsubMessageAckDeadline returns an attribute KeyValue conforming to the "messaging.gcp_pubsub.message.ack_deadline"semantic conventions
// It represents the The ack deadline in seconds set for the modify ack deadline request
func MessagingGcpPubsubMessageAckDeadline(val int) attribute.KeyValue {
	return MessagingGcpPubsubMessageAckDeadlineKey.Int(val)
}

// MessagingGcpPubsubMessageAckId returns an attribute KeyValue conforming to the "messaging.gcp_pubsub.message.ack_id"semantic conventions
// It represents the The ack id for a given message
func MessagingGcpPubsubMessageAckId(val string) attribute.KeyValue {
	return MessagingGcpPubsubMessageAckIdKey.String(val)
}

// MessagingGcpPubsubMessageDeliveryAttempt returns an attribute KeyValue conforming to the "messaging.gcp_pubsub.message.delivery_attempt"semantic conventions
// It represents the The delivery attempt for a given message
func MessagingGcpPubsubMessageDeliveryAttempt(val int) attribute.KeyValue {
	return MessagingGcpPubsubMessageDeliveryAttemptKey.Int(val)
}

// MessagingGcpPubsubMessageOrderingKey returns an attribute KeyValue conforming to the "messaging.gcp_pubsub.message.ordering_key"semantic conventions
// It represents the The ordering key for a given message. If the attribute is not present, the message does not have an ordering key
func MessagingGcpPubsubMessageOrderingKey(val string) attribute.KeyValue {
	return MessagingGcpPubsubMessageOrderingKeyKey.String(val)
}

// MessagingKafkaConsumerGroup returns an attribute KeyValue conforming to the "messaging.kafka.consumer.group"semantic conventions
// It represents the Deprecated, use `messaging.consumer.group.name` instead
func MessagingKafkaConsumerGroup(val string) attribute.KeyValue {
	return MessagingKafkaConsumerGroupKey.String(val)
}

// MessagingKafkaDestinationPartition returns an attribute KeyValue conforming to the "messaging.kafka.destination.partition"semantic conventions
// It represents the Deprecated, use `messaging.destination.partition.id` instead
func MessagingKafkaDestinationPartition(val int) attribute.KeyValue {
	return MessagingKafkaDestinationPartitionKey.Int(val)
}

// MessagingKafkaMessageKey returns an attribute KeyValue conforming to the "messaging.kafka.message.key"semantic conventions
// It represents the Message keys in Kafka are used for grouping alike messages to ensure they're processed on the same partition. They differ from `messaging.message.id` in that they're not unique. If the key is `null`, the attribute MUST NOT be set
func MessagingKafkaMessageKey(val string) attribute.KeyValue {
	return MessagingKafkaMessageKeyKey.String(val)
}

// MessagingKafkaMessageOffset returns an attribute KeyValue conforming to the "messaging.kafka.message.offset"semantic conventions
// It represents the Deprecated, use `messaging.kafka.offset` instead
func MessagingKafkaMessageOffset(val int) attribute.KeyValue {
	return MessagingKafkaMessageOffsetKey.Int(val)
}

// MessagingKafkaMessageTombstone returns an attribute KeyValue conforming to the "messaging.kafka.message.tombstone"semantic conventions
// It represents the A boolean that is true if the message is a tombstone
func MessagingKafkaMessageTombstone(val bool) attribute.KeyValue {
	return MessagingKafkaMessageTombstoneKey.Bool(val)
}

// MessagingKafkaOffset returns an attribute KeyValue conforming to the "messaging.kafka.offset"semantic conventions
// It represents the The offset of a record in the corresponding Kafka partition
func MessagingKafkaOffset(val int) attribute.KeyValue {
	return MessagingKafkaOffsetKey.Int(val)
}

// MessagingMessageBodySize returns an attribute KeyValue conforming to the "messaging.message.body.size"semantic conventions
// It represents the The size of the message body in bytes
func MessagingMessageBodySize(val int) attribute.KeyValue {
	return MessagingMessageBodySizeKey.Int(val)
}

// MessagingMessageConversationId returns an attribute KeyValue conforming to the "messaging.message.conversation_id"semantic conventions
// It represents the The conversation ID identifying the conversation to which the message belongs, represented as a string. Sometimes called "Correlation ID"
func MessagingMessageConversationId(val string) attribute.KeyValue {
	return MessagingMessageConversationIdKey.String(val)
}

// MessagingMessageEnvelopeSize returns an attribute KeyValue conforming to the "messaging.message.envelope.size"semantic conventions
// It represents the The size of the message body and metadata in bytes
func MessagingMessageEnvelopeSize(val int) attribute.KeyValue {
	return MessagingMessageEnvelopeSizeKey.Int(val)
}

// MessagingMessageId returns an attribute KeyValue conforming to the "messaging.message.id"semantic conventions
// It represents the A value used by the messaging system as an identifier for the message, represented as a string
func MessagingMessageId(val string) attribute.KeyValue {
	return MessagingMessageIdKey.String(val)
}

// MessagingOperation returns an attribute KeyValue conforming to the "messaging.operation"semantic conventions
// It represents the Deprecated, use `messaging.operation.type` instead
func MessagingOperation(val string) attribute.KeyValue {
	return MessagingOperationKey.String(val)
}

// MessagingOperationName returns an attribute KeyValue conforming to the "messaging.operation.name"semantic conventions
// It represents the The system-specific name of the messaging operation
func MessagingOperationName(val string) attribute.KeyValue {
	return MessagingOperationNameKey.String(val)
}

// MessagingOperationType returns an attribute KeyValue conforming to the "messaging.operation.type"semantic conventions
// It represents the A string identifying the type of the messaging operation
func MessagingOperationType(val string) attribute.KeyValue {
	return MessagingOperationTypeKey.String(val)
}

// MessagingRabbitmqDestinationRoutingKey returns an attribute KeyValue conforming to the "messaging.rabbitmq.destination.routing_key"semantic conventions
// It represents the RabbitMQ message routing key
func MessagingRabbitmqDestinationRoutingKey(val string) attribute.KeyValue {
	return MessagingRabbitmqDestinationRoutingKeyKey.String(val)
}

// MessagingRabbitmqMessageDeliveryTag returns an attribute KeyValue conforming to the "messaging.rabbitmq.message.delivery_tag"semantic conventions
// It represents the RabbitMQ message delivery tag
func MessagingRabbitmqMessageDeliveryTag(val int) attribute.KeyValue {
	return MessagingRabbitmqMessageDeliveryTagKey.Int(val)
}

// MessagingRocketmqClientGroup returns an attribute KeyValue conforming to the "messaging.rocketmq.client_group"semantic conventions
// It represents the Deprecated, use `messaging.consumer.group.name` instead
func MessagingRocketmqClientGroup(val string) attribute.KeyValue {
	return MessagingRocketmqClientGroupKey.String(val)
}

// MessagingRocketmqConsumptionModel returns an attribute KeyValue conforming to the "messaging.rocketmq.consumption_model"semantic conventions
// It represents the Model of message consumption. This only applies to consumer spans
func MessagingRocketmqConsumptionModel(val string) attribute.KeyValue {
	return MessagingRocketmqConsumptionModelKey.String(val)
}

// MessagingRocketmqMessageDelayTimeLevel returns an attribute KeyValue conforming to the "messaging.rocketmq.message.delay_time_level"semantic conventions
// It represents the The delay time level for delay message, which determines the message delay time
func MessagingRocketmqMessageDelayTimeLevel(val int) attribute.KeyValue {
	return MessagingRocketmqMessageDelayTimeLevelKey.Int(val)
}

// MessagingRocketmqMessageDeliveryTimestamp returns an attribute KeyValue conforming to the "messaging.rocketmq.message.delivery_timestamp"semantic conventions
// It represents the The timestamp in milliseconds that the delay message is expected to be delivered to consumer
func MessagingRocketmqMessageDeliveryTimestamp(val int) attribute.KeyValue {
	return MessagingRocketmqMessageDeliveryTimestampKey.Int(val)
}

// MessagingRocketmqMessageGroup returns an attribute KeyValue conforming to the "messaging.rocketmq.message.group"semantic conventions
// It represents the It is essential for FIFO message. Messages that belong to the same message group are always processed one by one within the same consumer group
func MessagingRocketmqMessageGroup(val string) attribute.KeyValue {
	return MessagingRocketmqMessageGroupKey.String(val)
}

// MessagingRocketmqMessageKeys returns an attribute KeyValue conforming to the "messaging.rocketmq.message.keys"semantic conventions
// It represents the Key(s) of message, another way to mark message besides message id
func MessagingRocketmqMessageKeys(val []string) attribute.KeyValue {
	return MessagingRocketmqMessageKeysKey.StringSlice(val)
}

// MessagingRocketmqMessageTag returns an attribute KeyValue conforming to the "messaging.rocketmq.message.tag"semantic conventions
// It represents the The secondary classifier of message besides topic
func MessagingRocketmqMessageTag(val string) attribute.KeyValue {
	return MessagingRocketmqMessageTagKey.String(val)
}

// MessagingRocketmqMessageType returns an attribute KeyValue conforming to the "messaging.rocketmq.message.type"semantic conventions
// It represents the Type of message
func MessagingRocketmqMessageType(val string) attribute.KeyValue {
	return MessagingRocketmqMessageTypeKey.String(val)
}

// MessagingRocketmqNamespace returns an attribute KeyValue conforming to the "messaging.rocketmq.namespace"semantic conventions
// It represents the Namespace of RocketMQ resources, resources in different namespaces are individual
func MessagingRocketmqNamespace(val string) attribute.KeyValue {
	return MessagingRocketmqNamespaceKey.String(val)
}

// MessagingServicebusDestinationSubscriptionName returns an attribute KeyValue conforming to the "messaging.servicebus.destination.subscription_name"semantic conventions
// It represents the Deprecated, use `messaging.servicebus.destination.subscription_name` instead
func MessagingServicebusDestinationSubscriptionName(val string) attribute.KeyValue {
	return MessagingServicebusDestinationSubscriptionNameKey.String(val)
}

// MessagingServicebusDispositionStatus returns an attribute KeyValue conforming to the "messaging.servicebus.disposition_status"semantic conventions
// It represents the Describes the [settlement type]
//
// [settlement type]: https://learn.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#peeklock
func MessagingServicebusDispositionStatus(val string) attribute.KeyValue {
	return MessagingServicebusDispositionStatusKey.String(val)
}

// MessagingServicebusMessageDeliveryCount returns an attribute KeyValue conforming to the "messaging.servicebus.message.delivery_count"semantic conventions
// It represents the Number of deliveries that have been attempted for this message
func MessagingServicebusMessageDeliveryCount(val int) attribute.KeyValue {
	return MessagingServicebusMessageDeliveryCountKey.Int(val)
}

// MessagingServicebusMessageEnqueuedTime returns an attribute KeyValue conforming to the "messaging.servicebus.message.enqueued_time"semantic conventions
// It represents the The UTC epoch seconds at which the message has been accepted and stored in the entity
func MessagingServicebusMessageEnqueuedTime(val int) attribute.KeyValue {
	return MessagingServicebusMessageEnqueuedTimeKey.Int(val)
}

// MessagingSystem returns an attribute KeyValue conforming to the "messaging.system"semantic conventions
// It represents the The messaging system as identified by the client instrumentation
func MessagingSystem(val string) attribute.KeyValue {
	return MessagingSystemKey.String(val)
}

// Enum values for messaging.operation.type
var (

	// One or more messages are provided for publishing to an intermediary. If a single message is published, the context of the "Publish" span can be used as the creation context and no "Create" span needs to be created.
	//
	// Stability: experimental
	MessagingOperationTypePublish = MessagingOperationTypeKey.String("publish")
	// A message is created. "Create" spans always refer to a single message and are used to provide a unique creation context for messages in batch publishing scenarios.
	//
	// Stability: experimental
	MessagingOperationTypeCreate = MessagingOperationTypeKey.String("create")
	// One or more messages are requested by a consumer. This operation refers to pull-based scenarios, where consumers explicitly call methods of messaging SDKs to receive messages.
	//
	// Stability: experimental
	MessagingOperationTypeReceive = MessagingOperationTypeKey.String("receive")
	// One or more messages are processed by a consumer.
	//
	// Stability: experimental
	MessagingOperationTypeProcess = MessagingOperationTypeKey.String("process")
	// One or more messages are settled.
	//
	// Stability: experimental
	MessagingOperationTypeSettle = MessagingOperationTypeKey.String("settle")
	// Deprecated. Use `process` instead.
	// Stability: experimental// Replaced by `process`
	MessagingOperationTypeDeliver = MessagingOperationTypeKey.String("deliver")
)

// Enum values for messaging.rocketmq.consumption_model
var (

	// Clustering consumption model
	// Stability: experimental
	MessagingRocketmqConsumptionModelClustering = MessagingRocketmqConsumptionModelKey.String("clustering")
	// Broadcasting consumption model
	// Stability: experimental
	MessagingRocketmqConsumptionModelBroadcasting = MessagingRocketmqConsumptionModelKey.String("broadcasting")
)

// Enum values for messaging.rocketmq.message.type
var (

	// Normal message
	// Stability: experimental
	MessagingRocketmqMessageTypeNormal = MessagingRocketmqMessageTypeKey.String("normal")
	// FIFO message
	// Stability: experimental
	MessagingRocketmqMessageTypeFifo = MessagingRocketmqMessageTypeKey.String("fifo")
	// Delay message
	// Stability: experimental
	MessagingRocketmqMessageTypeDelay = MessagingRocketmqMessageTypeKey.String("delay")
	// Transaction message
	// Stability: experimental
	MessagingRocketmqMessageTypeTransaction = MessagingRocketmqMessageTypeKey.String("transaction")
)

// Enum values for messaging.servicebus.disposition_status
var (

	// Message is completed
	// Stability: experimental
	MessagingServicebusDispositionStatusComplete = MessagingServicebusDispositionStatusKey.String("complete")
	// Message is abandoned
	// Stability: experimental
	MessagingServicebusDispositionStatusAbandon = MessagingServicebusDispositionStatusKey.String("abandon")
	// Message is sent to dead letter queue
	// Stability: experimental
	MessagingServicebusDispositionStatusDeadLetter = MessagingServicebusDispositionStatusKey.String("dead_letter")
	// Message is deferred
	// Stability: experimental
	MessagingServicebusDispositionStatusDefer = MessagingServicebusDispositionStatusKey.String("defer")
)

// Enum values for messaging.system
var (

	// Apache ActiveMQ
	// Stability: experimental
	MessagingSystemActivemq = MessagingSystemKey.String("activemq")
	// Amazon Simple Queue Service (SQS)
	// Stability: experimental
	MessagingSystemAwsSqs = MessagingSystemKey.String("aws_sqs")
	// Azure Event Grid
	// Stability: experimental
	MessagingSystemEventgrid = MessagingSystemKey.String("eventgrid")
	// Azure Event Hubs
	// Stability: experimental
	MessagingSystemEventhubs = MessagingSystemKey.String("eventhubs")
	// Azure Service Bus
	// Stability: experimental
	MessagingSystemServicebus = MessagingSystemKey.String("servicebus")
	// Google Cloud Pub/Sub
	// Stability: experimental
	MessagingSystemGcpPubsub = MessagingSystemKey.String("gcp_pubsub")
	// Java Message Service
	// Stability: experimental
	MessagingSystemJms = MessagingSystemKey.String("jms")
	// Apache Kafka
	// Stability: experimental
	MessagingSystemKafka = MessagingSystemKey.String("kafka")
	// RabbitMQ
	// Stability: experimental
	MessagingSystemRabbitmq = MessagingSystemKey.String("rabbitmq")
	// Apache RocketMQ
	// Stability: experimental
	MessagingSystemRocketmq = MessagingSystemKey.String("rocketmq")
	// Apache Pulsar
	// Stability: experimental
	MessagingSystemPulsar = MessagingSystemKey.String("pulsar")
)

// Namespace: net
const (

	// Deprecated, use `network.local.address`.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `network.local.address`.
	//
	// Examples: "192.168.0.1"
	NetHostIpKey = attribute.Key("net.host.ip")
	// Deprecated, use `server.address`.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `server.address`.
	//
	// Examples:
	// "example.com",
	NetHostNameKey = attribute.Key("net.host.name")
	// Deprecated, use `server.port`.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `server.port`.
	//
	// Examples:
	// 8080,
	NetHostPortKey = attribute.Key("net.host.port")
	// Deprecated, use `network.peer.address`.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `network.peer.address`.
	//
	// Examples: "127.0.0.1"
	NetPeerIpKey = attribute.Key("net.peer.ip")
	// Deprecated, use `server.address` on client spans and `client.address` on server spans.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `server.address` on client spans and `client.address` on server spans.
	//
	// Examples:
	// "example.com",
	NetPeerNameKey = attribute.Key("net.peer.name")
	// Deprecated, use `server.port` on client spans and `client.port` on server spans.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `server.port` on client spans and `client.port` on server spans.
	//
	// Examples:
	// 8080,
	NetPeerPortKey = attribute.Key("net.peer.port")
	// Deprecated, use `network.protocol.name`.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `network.protocol.name`.
	//
	// Examples:
	// "amqp",
	// "http",
	// "mqtt",
	NetProtocolNameKey = attribute.Key("net.protocol.name")
	// Deprecated, use `network.protocol.version`.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `network.protocol.version`.
	//
	// Examples: "3.1.1"
	NetProtocolVersionKey = attribute.Key("net.protocol.version")
	// Deprecated, use `network.transport` and `network.type`.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Split to `network.transport` and `network.type`.
	//
	// Examples: undefined
	NetSockFamilyKey = attribute.Key("net.sock.family")
	// Deprecated, use `network.local.address`.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `network.local.address`.
	//
	// Examples:
	// "/var/my.sock",
	NetSockHostAddrKey = attribute.Key("net.sock.host.addr")
	// Deprecated, use `network.local.port`.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `network.local.port`.
	//
	// Examples:
	// 8080,
	NetSockHostPortKey = attribute.Key("net.sock.host.port")
	// Deprecated, use `network.peer.address`.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `network.peer.address`.
	//
	// Examples:
	// "192.168.0.1",
	NetSockPeerAddrKey = attribute.Key("net.sock.peer.addr")
	// Deprecated, no replacement at this time.
	// Stability: Experimental
	// Type: string
	// Deprecated: Removed.
	//
	// Examples:
	// "/var/my.sock",
	NetSockPeerNameKey = attribute.Key("net.sock.peer.name")
	// Deprecated, use `network.peer.port`.
	// Stability: Experimental
	// Type: int
	// Deprecated: Replaced by `network.peer.port`.
	//
	// Examples:
	// 65531,
	NetSockPeerPortKey = attribute.Key("net.sock.peer.port")
	// Deprecated, use `network.transport`.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `network.transport`.
	//
	// Examples: undefined
	NetTransportKey = attribute.Key("net.transport")
)

// NetHostIp returns an attribute KeyValue conforming to the "net.host.ip"semantic conventions
// It represents the Deprecated, use `network.local.address`
func NetHostIp(val string) attribute.KeyValue {
	return NetHostIpKey.String(val)
}

// NetHostName returns an attribute KeyValue conforming to the "net.host.name"semantic conventions
// It represents the Deprecated, use `server.address`
func NetHostName(val string) attribute.KeyValue {
	return NetHostNameKey.String(val)
}

// NetHostPort returns an attribute KeyValue conforming to the "net.host.port"semantic conventions
// It represents the Deprecated, use `server.port`
func NetHostPort(val int) attribute.KeyValue {
	return NetHostPortKey.Int(val)
}

// NetPeerIp returns an attribute KeyValue conforming to the "net.peer.ip"semantic conventions
// It represents the Deprecated, use `network.peer.address`
func NetPeerIp(val string) attribute.KeyValue {
	return NetPeerIpKey.String(val)
}

// NetPeerName returns an attribute KeyValue conforming to the "net.peer.name"semantic conventions
// It represents the Deprecated, use `server.address` on client spans and `client.address` on server spans
func NetPeerName(val string) attribute.KeyValue {
	return NetPeerNameKey.String(val)
}

// NetPeerPort returns an attribute KeyValue conforming to the "net.peer.port"semantic conventions
// It represents the Deprecated, use `server.port` on client spans and `client.port` on server spans
func NetPeerPort(val int) attribute.KeyValue {
	return NetPeerPortKey.Int(val)
}

// NetProtocolName returns an attribute KeyValue conforming to the "net.protocol.name"semantic conventions
// It represents the Deprecated, use `network.protocol.name`
func NetProtocolName(val string) attribute.KeyValue {
	return NetProtocolNameKey.String(val)
}

// NetProtocolVersion returns an attribute KeyValue conforming to the "net.protocol.version"semantic conventions
// It represents the Deprecated, use `network.protocol.version`
func NetProtocolVersion(val string) attribute.KeyValue {
	return NetProtocolVersionKey.String(val)
}

// NetSockFamily returns an attribute KeyValue conforming to the "net.sock.family"semantic conventions
// It represents the Deprecated, use `network.transport` and `network.type`
func NetSockFamily(val string) attribute.KeyValue {
	return NetSockFamilyKey.String(val)
}

// NetSockHostAddr returns an attribute KeyValue conforming to the "net.sock.host.addr"semantic conventions
// It represents the Deprecated, use `network.local.address`
func NetSockHostAddr(val string) attribute.KeyValue {
	return NetSockHostAddrKey.String(val)
}

// NetSockHostPort returns an attribute KeyValue conforming to the "net.sock.host.port"semantic conventions
// It represents the Deprecated, use `network.local.port`
func NetSockHostPort(val int) attribute.KeyValue {
	return NetSockHostPortKey.Int(val)
}

// NetSockPeerAddr returns an attribute KeyValue conforming to the "net.sock.peer.addr"semantic conventions
// It represents the Deprecated, use `network.peer.address`
func NetSockPeerAddr(val string) attribute.KeyValue {
	return NetSockPeerAddrKey.String(val)
}

// NetSockPeerName returns an attribute KeyValue conforming to the "net.sock.peer.name"semantic conventions
// It represents the Deprecated, no replacement at this time
func NetSockPeerName(val string) attribute.KeyValue {
	return NetSockPeerNameKey.String(val)
}

// NetSockPeerPort returns an attribute KeyValue conforming to the "net.sock.peer.port"semantic conventions
// It represents the Deprecated, use `network.peer.port`
func NetSockPeerPort(val int) attribute.KeyValue {
	return NetSockPeerPortKey.Int(val)
}

// NetTransport returns an attribute KeyValue conforming to the "net.transport"semantic conventions
// It represents the Deprecated, use `network.transport`
func NetTransport(val string) attribute.KeyValue {
	return NetTransportKey.String(val)
}

// Enum values for net.sock.family
var (

	// IPv4 address
	// Stability: experimental
	NetSockFamilyInet = NetSockFamilyKey.String("inet")
	// IPv6 address
	// Stability: experimental
	NetSockFamilyInet6 = NetSockFamilyKey.String("inet6")
	// Unix domain socket path
	// Stability: experimental
	NetSockFamilyUnix = NetSockFamilyKey.String("unix")
)

// Enum values for net.transport
var (

	// ip_tcp
	// Stability: experimental
	NetTransportIpTcp = NetTransportKey.String("ip_tcp")
	// ip_udp
	// Stability: experimental
	NetTransportIpUdp = NetTransportKey.String("ip_udp")
	// Named or anonymous pipe.
	// Stability: experimental
	NetTransportPipe = NetTransportKey.String("pipe")
	// In-process communication.
	// Stability: experimental
	NetTransportInproc = NetTransportKey.String("inproc")
	// Something else (non IP-based).
	// Stability: experimental
	NetTransportOther = NetTransportKey.String("other")
)

// Namespace: network
const (

	// The ISO 3166-1 alpha-2 2-character country code associated with the mobile carrier network.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "DE"
	NetworkCarrierIccKey = attribute.Key("network.carrier.icc")
	// The mobile carrier country code.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "310"
	NetworkCarrierMccKey = attribute.Key("network.carrier.mcc")
	// The mobile carrier network code.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "001"
	NetworkCarrierMncKey = attribute.Key("network.carrier.mnc")
	// The name of the mobile carrier.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "sprint"
	NetworkCarrierNameKey = attribute.Key("network.carrier.name")
	// This describes more details regarding the connection.type. It may be the type of cell technology connection, but it could be used for describing details about a wifi connection.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: "LTE"
	NetworkConnectionSubtypeKey = attribute.Key("network.connection.subtype")
	// The internet connection type.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: "wifi"
	NetworkConnectionTypeKey = attribute.Key("network.connection.type")
	// The network IO operation direction.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "transmit",
	NetworkIoDirectionKey = attribute.Key("network.io.direction")
	// Local address of the network connection - IP address or Unix domain socket name.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "10.1.2.80",
	// "/tmp/my.sock",
	NetworkLocalAddressKey = attribute.Key("network.local.address")
	// Local port number of the network connection.
	// Stability: Stable
	// Type: int
	//
	// Examples:
	// 65123,
	NetworkLocalPortKey = attribute.Key("network.local.port")
	// Peer address of the network connection - IP address or Unix domain socket name.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "10.1.2.80",
	// "/tmp/my.sock",
	NetworkPeerAddressKey = attribute.Key("network.peer.address")
	// Peer port number of the network connection.
	// Stability: Stable
	// Type: int
	//
	// Examples:
	// 65123,
	NetworkPeerPortKey = attribute.Key("network.peer.port")
	// [OSI application layer] or non-OSI equivalent.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "amqp",
	// "http",
	// "mqtt",
	//
	// Note: The value SHOULD be normalized to lowercase
	//
	// [OSI application layer]: https://osi-model.com/application-layer/
	NetworkProtocolNameKey = attribute.Key("network.protocol.name")
	// The actual version of the protocol used for network communication.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "1.1",
	// "2",
	//
	// Note: If protocol version is subject to negotiation (for example using [ALPN]), this attribute SHOULD be set to the negotiated version. If the actual protocol version is not known, this attribute SHOULD NOT be set
	//
	// [ALPN]: https://www.rfc-editor.org/rfc/rfc7301.html
	NetworkProtocolVersionKey = attribute.Key("network.protocol.version")
	// [OSI transport layer] or [inter-process communication method].
	//
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "tcp",
	// "udp",
	//
	// Note: The value SHOULD be normalized to lowercase.
	//
	// Consider always setting the transport when setting a port number, since
	// a port number is ambiguous without knowing the transport. For example
	// different processes could be listening on TCP port 12345 and UDP port 12345
	//
	// [OSI transport layer]: https://osi-model.com/transport-layer/
	// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
	NetworkTransportKey = attribute.Key("network.transport")
	// [OSI network layer] or non-OSI equivalent.
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "ipv4",
	// "ipv6",
	//
	// Note: The value SHOULD be normalized to lowercase
	//
	// [OSI network layer]: https://osi-model.com/network-layer/
	NetworkTypeKey = attribute.Key("network.type")
)

// NetworkCarrierIcc returns an attribute KeyValue conforming to the "network.carrier.icc"semantic conventions
// It represents the The ISO 3166-1 alpha-2 2-character country code associated with the mobile carrier network
func NetworkCarrierIcc(val string) attribute.KeyValue {
	return NetworkCarrierIccKey.String(val)
}

// NetworkCarrierMcc returns an attribute KeyValue conforming to the "network.carrier.mcc"semantic conventions
// It represents the The mobile carrier country code
func NetworkCarrierMcc(val string) attribute.KeyValue {
	return NetworkCarrierMccKey.String(val)
}

// NetworkCarrierMnc returns an attribute KeyValue conforming to the "network.carrier.mnc"semantic conventions
// It represents the The mobile carrier network code
func NetworkCarrierMnc(val string) attribute.KeyValue {
	return NetworkCarrierMncKey.String(val)
}

// NetworkCarrierName returns an attribute KeyValue conforming to the "network.carrier.name"semantic conventions
// It represents the The name of the mobile carrier
func NetworkCarrierName(val string) attribute.KeyValue {
	return NetworkCarrierNameKey.String(val)
}

// NetworkConnectionSubtype returns an attribute KeyValue conforming to the "network.connection.subtype"semantic conventions
// It represents the This describes more details regarding the connection.type. It may be the type of cell technology connection, but it could be used for describing details about a wifi connection
func NetworkConnectionSubtype(val string) attribute.KeyValue {
	return NetworkConnectionSubtypeKey.String(val)
}

// NetworkConnectionType returns an attribute KeyValue conforming to the "network.connection.type"semantic conventions
// It represents the The internet connection type
func NetworkConnectionType(val string) attribute.KeyValue {
	return NetworkConnectionTypeKey.String(val)
}

// NetworkIoDirection returns an attribute KeyValue conforming to the "network.io.direction"semantic conventions
// It represents the The network IO operation direction
func NetworkIoDirection(val string) attribute.KeyValue {
	return NetworkIoDirectionKey.String(val)
}

// NetworkLocalAddress returns an attribute KeyValue conforming to the "network.local.address"semantic conventions
// It represents the Local address of the network connection - IP address or Unix domain socket name
func NetworkLocalAddress(val string) attribute.KeyValue {
	return NetworkLocalAddressKey.String(val)
}

// NetworkLocalPort returns an attribute KeyValue conforming to the "network.local.port"semantic conventions
// It represents the Local port number of the network connection
func NetworkLocalPort(val int) attribute.KeyValue {
	return NetworkLocalPortKey.Int(val)
}

// NetworkPeerAddress returns an attribute KeyValue conforming to the "network.peer.address"semantic conventions
// It represents the Peer address of the network connection - IP address or Unix domain socket name
func NetworkPeerAddress(val string) attribute.KeyValue {
	return NetworkPeerAddressKey.String(val)
}

// NetworkPeerPort returns an attribute KeyValue conforming to the "network.peer.port"semantic conventions
// It represents the Peer port number of the network connection
func NetworkPeerPort(val int) attribute.KeyValue {
	return NetworkPeerPortKey.Int(val)
}

// NetworkProtocolName returns an attribute KeyValue conforming to the "network.protocol.name"semantic conventions
// It represents the [OSI application layer] or non-OSI equivalent
//
// [OSI application layer]: https://osi-model.com/application-layer/
func NetworkProtocolName(val string) attribute.KeyValue {
	return NetworkProtocolNameKey.String(val)
}

// NetworkProtocolVersion returns an attribute KeyValue conforming to the "network.protocol.version"semantic conventions
// It represents the The actual version of the protocol used for network communication
func NetworkProtocolVersion(val string) attribute.KeyValue {
	return NetworkProtocolVersionKey.String(val)
}

// NetworkTransport returns an attribute KeyValue conforming to the "network.transport"semantic conventions
// It represents the [OSI transport layer] or [inter-process communication method]
//
// [OSI transport layer]: https://osi-model.com/transport-layer/
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func NetworkTransport(val string) attribute.KeyValue {
	return NetworkTransportKey.String(val)
}

// NetworkType returns an attribute KeyValue conforming to the "network.type"semantic conventions
// It represents the [OSI network layer] or non-OSI equivalent
//
// [OSI network layer]: https://osi-model.com/network-layer/
func NetworkType(val string) attribute.KeyValue {
	return NetworkTypeKey.String(val)
}

// Enum values for network.connection.subtype
var (

	// GPRS
	// Stability: experimental
	NetworkConnectionSubtypeGprs = NetworkConnectionSubtypeKey.String("gprs")
	// EDGE
	// Stability: experimental
	NetworkConnectionSubtypeEdge = NetworkConnectionSubtypeKey.String("edge")
	// UMTS
	// Stability: experimental
	NetworkConnectionSubtypeUmts = NetworkConnectionSubtypeKey.String("umts")
	// CDMA
	// Stability: experimental
	NetworkConnectionSubtypeCdma = NetworkConnectionSubtypeKey.String("cdma")
	// EVDO Rel. 0
	// Stability: experimental
	NetworkConnectionSubtypeEvdo0 = NetworkConnectionSubtypeKey.String("evdo_0")
	// EVDO Rev. A
	// Stability: experimental
	NetworkConnectionSubtypeEvdoA = NetworkConnectionSubtypeKey.String("evdo_a")
	// CDMA2000 1XRTT
	// Stability: experimental
	NetworkConnectionSubtypeCdma20001xrtt = NetworkConnectionSubtypeKey.String("cdma2000_1xrtt")
	// HSDPA
	// Stability: experimental
	NetworkConnectionSubtypeHsdpa = NetworkConnectionSubtypeKey.String("hsdpa")
	// HSUPA
	// Stability: experimental
	NetworkConnectionSubtypeHsupa = NetworkConnectionSubtypeKey.String("hsupa")
	// HSPA
	// Stability: experimental
	NetworkConnectionSubtypeHspa = NetworkConnectionSubtypeKey.String("hspa")
	// IDEN
	// Stability: experimental
	NetworkConnectionSubtypeIden = NetworkConnectionSubtypeKey.String("iden")
	// EVDO Rev. B
	// Stability: experimental
	NetworkConnectionSubtypeEvdoB = NetworkConnectionSubtypeKey.String("evdo_b")
	// LTE
	// Stability: experimental
	NetworkConnectionSubtypeLte = NetworkConnectionSubtypeKey.String("lte")
	// EHRPD
	// Stability: experimental
	NetworkConnectionSubtypeEhrpd = NetworkConnectionSubtypeKey.String("ehrpd")
	// HSPAP
	// Stability: experimental
	NetworkConnectionSubtypeHspap = NetworkConnectionSubtypeKey.String("hspap")
	// GSM
	// Stability: experimental
	NetworkConnectionSubtypeGsm = NetworkConnectionSubtypeKey.String("gsm")
	// TD-SCDMA
	// Stability: experimental
	NetworkConnectionSubtypeTdScdma = NetworkConnectionSubtypeKey.String("td_scdma")
	// IWLAN
	// Stability: experimental
	NetworkConnectionSubtypeIwlan = NetworkConnectionSubtypeKey.String("iwlan")
	// 5G NR (New Radio)
	// Stability: experimental
	NetworkConnectionSubtypeNr = NetworkConnectionSubtypeKey.String("nr")
	// 5G NRNSA (New Radio Non-Standalone)
	// Stability: experimental
	NetworkConnectionSubtypeNrnsa = NetworkConnectionSubtypeKey.String("nrnsa")
	// LTE CA
	// Stability: experimental
	NetworkConnectionSubtypeLteCa = NetworkConnectionSubtypeKey.String("lte_ca")
)

// Enum values for network.connection.type
var (

	// wifi
	// Stability: experimental
	NetworkConnectionTypeWifi = NetworkConnectionTypeKey.String("wifi")
	// wired
	// Stability: experimental
	NetworkConnectionTypeWired = NetworkConnectionTypeKey.String("wired")
	// cell
	// Stability: experimental
	NetworkConnectionTypeCell = NetworkConnectionTypeKey.String("cell")
	// unavailable
	// Stability: experimental
	NetworkConnectionTypeUnavailable = NetworkConnectionTypeKey.String("unavailable")
	// unknown
	// Stability: experimental
	NetworkConnectionTypeUnknown = NetworkConnectionTypeKey.String("unknown")
)

// Enum values for network.io.direction
var (

	// transmit
	// Stability: experimental
	NetworkIoDirectionTransmit = NetworkIoDirectionKey.String("transmit")
	// receive
	// Stability: experimental
	NetworkIoDirectionReceive = NetworkIoDirectionKey.String("receive")
)

// Enum values for network.transport
var (

	// TCP
	// Stability: stable
	NetworkTransportTcp = NetworkTransportKey.String("tcp")
	// UDP
	// Stability: stable
	NetworkTransportUdp = NetworkTransportKey.String("udp")
	// Named or anonymous pipe.
	// Stability: stable
	NetworkTransportPipe = NetworkTransportKey.String("pipe")
	// Unix domain socket
	// Stability: stable
	NetworkTransportUnix = NetworkTransportKey.String("unix")
	// QUIC
	// Stability: experimental
	NetworkTransportQuic = NetworkTransportKey.String("quic")
)

// Enum values for network.type
var (

	// IPv4
	// Stability: stable
	NetworkTypeIpv4 = NetworkTypeKey.String("ipv4")
	// IPv6
	// Stability: stable
	NetworkTypeIpv6 = NetworkTypeKey.String("ipv6")
)

// Namespace: nodejs
const (

	// The state of event loop time.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	NodejsEventloopStateKey = attribute.Key("nodejs.eventloop.state")
)

// NodejsEventloopState returns an attribute KeyValue conforming to the "nodejs.eventloop.state"semantic conventions
// It represents the The state of event loop time
func NodejsEventloopState(val string) attribute.KeyValue {
	return NodejsEventloopStateKey.String(val)
}

// Enum values for nodejs.eventloop.state
var (

	// Active time.
	// Stability: experimental
	NodejsEventloopStateActive = NodejsEventloopStateKey.String("active")
	// Idle time.
	// Stability: experimental
	NodejsEventloopStateIdle = NodejsEventloopStateKey.String("idle")
)

// Namespace: oci
const (

	// The digest of the OCI image manifest. For container images specifically is the digest by which the container image is known.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "sha256:e4ca62c0d62f3e886e684806dfe9d4e0cda60d54986898173c1083856cfda0f4",
	//
	// Note: Follows [OCI Image Manifest Specification], and specifically the [Digest property].
	// An example can be found in [Example Image Manifest]
	//
	// [OCI Image Manifest Specification]: https://github.com/opencontainers/image-spec/blob/main/manifest.md
	// [Digest property]: https://github.com/opencontainers/image-spec/blob/main/descriptor.md#digests
	// [Example Image Manifest]: https://docs.docker.com/registry/spec/manifest-v2-2/#example-image-manifest
	OciManifestDigestKey = attribute.Key("oci.manifest.digest")
)

// OciManifestDigest returns an attribute KeyValue conforming to the "oci.manifest.digest"semantic conventions
// It represents the The digest of the OCI image manifest. For container images specifically is the digest by which the container image is known
func OciManifestDigest(val string) attribute.KeyValue {
	return OciManifestDigestKey.String(val)
}

// Namespace: opentracing
const (

	// Parent-child Reference type
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: The causal relationship between a child Span and a parent Span
	OpentracingRefTypeKey = attribute.Key("opentracing.ref_type")
)

// OpentracingRefType returns an attribute KeyValue conforming to the "opentracing.ref_type"semantic conventions
// It represents the Parent-child Reference type
func OpentracingRefType(val string) attribute.KeyValue {
	return OpentracingRefTypeKey.String(val)
}

// Enum values for opentracing.ref_type
var (

	// The parent Span depends on the child Span in some capacity
	// Stability: experimental
	OpentracingRefTypeChildOf = OpentracingRefTypeKey.String("child_of")
	// The parent Span doesn't depend in any way on the result of the child Span
	// Stability: experimental
	OpentracingRefTypeFollowsFrom = OpentracingRefTypeKey.String("follows_from")
)

// Namespace: os
const (

	// Unique identifier for a particular build or compilation of the operating system.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "TQ3C.230805.001.B2",
	// "20E247",
	// "22621",
	OsBuildIdKey = attribute.Key("os.build_id")
	// Human readable (not intended to be parsed) OS version information, like e.g. reported by `ver` or `lsb_release -a` commands.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Microsoft Windows [Version 10.0.18363.778]",
	// "Ubuntu 18.04.1 LTS",
	OsDescriptionKey = attribute.Key("os.description")
	// Human readable operating system name.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "iOS",
	// "Android",
	// "Ubuntu",
	OsNameKey = attribute.Key("os.name")
	// The operating system type.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	OsTypeKey = attribute.Key("os.type")
	// The version string of the operating system as defined in [Version Attributes].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "14.2.1",
	// "18.04.1",
	//
	// [Version Attributes]: /docs/resource/README.md#version-attributes
	OsVersionKey = attribute.Key("os.version")
)

// OsBuildId returns an attribute KeyValue conforming to the "os.build_id"semantic conventions
// It represents the Unique identifier for a particular build or compilation of the operating system
func OsBuildId(val string) attribute.KeyValue {
	return OsBuildIdKey.String(val)
}

// OsDescription returns an attribute KeyValue conforming to the "os.description"semantic conventions
// It represents the Human readable (not intended to be parsed) OS version information, like e.g. reported by `ver` or `lsb_release -a` commands
func OsDescription(val string) attribute.KeyValue {
	return OsDescriptionKey.String(val)
}

// OsName returns an attribute KeyValue conforming to the "os.name"semantic conventions
// It represents the Human readable operating system name
func OsName(val string) attribute.KeyValue {
	return OsNameKey.String(val)
}

// OsType returns an attribute KeyValue conforming to the "os.type"semantic conventions
// It represents the The operating system type
func OsType(val string) attribute.KeyValue {
	return OsTypeKey.String(val)
}

// OsVersion returns an attribute KeyValue conforming to the "os.version"semantic conventions
// It represents the The version string of the operating system as defined in [Version Attributes]
//
// [Version Attributes]: /docs/resource/README.md#version-attributes
func OsVersion(val string) attribute.KeyValue {
	return OsVersionKey.String(val)
}

// Enum values for os.type
var (

	// Microsoft Windows
	// Stability: experimental
	OsTypeWindows = OsTypeKey.String("windows")
	// Linux
	// Stability: experimental
	OsTypeLinux = OsTypeKey.String("linux")
	// Apple Darwin
	// Stability: experimental
	OsTypeDarwin = OsTypeKey.String("darwin")
	// FreeBSD
	// Stability: experimental
	OsTypeFreebsd = OsTypeKey.String("freebsd")
	// NetBSD
	// Stability: experimental
	OsTypeNetbsd = OsTypeKey.String("netbsd")
	// OpenBSD
	// Stability: experimental
	OsTypeOpenbsd = OsTypeKey.String("openbsd")
	// DragonFly BSD
	// Stability: experimental
	OsTypeDragonflybsd = OsTypeKey.String("dragonflybsd")
	// HP-UX (Hewlett Packard Unix)
	// Stability: experimental
	OsTypeHpux = OsTypeKey.String("hpux")
	// AIX (Advanced Interactive eXecutive)
	// Stability: experimental
	OsTypeAix = OsTypeKey.String("aix")
	// SunOS, Oracle Solaris
	// Stability: experimental
	OsTypeSolaris = OsTypeKey.String("solaris")
	// IBM z/OS
	// Stability: experimental
	OsTypeZOs = OsTypeKey.String("z_os")
)

// Namespace: otel
const (

	// Deprecated. Use the `otel.scope.name` attribute
	// Stability: Experimental
	// Type: string
	// Deprecated: Use the `otel.scope.name` attribute.
	//
	// Examples:
	// "io.opentelemetry.contrib.mongodb",
	OtelLibraryNameKey = attribute.Key("otel.library.name")
	// Deprecated. Use the `otel.scope.version` attribute.
	// Stability: Experimental
	// Type: string
	// Deprecated: Use the `otel.scope.version` attribute.
	//
	// Examples:
	// "1.0.0",
	OtelLibraryVersionKey = attribute.Key("otel.library.version")
	// The name of the instrumentation scope - (`InstrumentationScope.Name` in OTLP).
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "io.opentelemetry.contrib.mongodb",
	OtelScopeNameKey = attribute.Key("otel.scope.name")
	// The version of the instrumentation scope - (`InstrumentationScope.Version` in OTLP).
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "1.0.0",
	OtelScopeVersionKey = attribute.Key("otel.scope.version")
	// Name of the code, either "OK" or "ERROR". MUST NOT be set if the status code is UNSET.
	// Stability: Stable
	// Type: Enum
	//
	// Examples: undefined
	OtelStatusCodeKey = attribute.Key("otel.status_code")
	// Description of the Status if it has a value, otherwise not set.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "resource not found",
	OtelStatusDescriptionKey = attribute.Key("otel.status_description")
)

// OtelLibraryName returns an attribute KeyValue conforming to the "otel.library.name"semantic conventions
// It represents the Deprecated. Use the `otel.scope.name` attribute
func OtelLibraryName(val string) attribute.KeyValue {
	return OtelLibraryNameKey.String(val)
}

// OtelLibraryVersion returns an attribute KeyValue conforming to the "otel.library.version"semantic conventions
// It represents the Deprecated. Use the `otel.scope.version` attribute
func OtelLibraryVersion(val string) attribute.KeyValue {
	return OtelLibraryVersionKey.String(val)
}

// OtelScopeName returns an attribute KeyValue conforming to the "otel.scope.name"semantic conventions
// It represents the The name of the instrumentation scope - (`InstrumentationScope.Name` in OTLP)
func OtelScopeName(val string) attribute.KeyValue {
	return OtelScopeNameKey.String(val)
}

// OtelScopeVersion returns an attribute KeyValue conforming to the "otel.scope.version"semantic conventions
// It represents the The version of the instrumentation scope - (`InstrumentationScope.Version` in OTLP)
func OtelScopeVersion(val string) attribute.KeyValue {
	return OtelScopeVersionKey.String(val)
}

// OtelStatusCode returns an attribute KeyValue conforming to the "otel.status_code"semantic conventions
// It represents the Name of the code, either "OK" or "ERROR". MUST NOT be set if the status code is UNSET
func OtelStatusCode(val string) attribute.KeyValue {
	return OtelStatusCodeKey.String(val)
}

// OtelStatusDescription returns an attribute KeyValue conforming to the "otel.status_description"semantic conventions
// It represents the Description of the Status if it has a value, otherwise not set
func OtelStatusDescription(val string) attribute.KeyValue {
	return OtelStatusDescriptionKey.String(val)
}

// Enum values for otel.status_code
var (

	// The operation has been validated by an Application developer or Operator to have completed successfully.
	// Stability: stable
	OtelStatusCodeOk = OtelStatusCodeKey.String("OK")
	// The operation contains an error.
	// Stability: stable
	OtelStatusCodeError = OtelStatusCodeKey.String("ERROR")
)

// Namespace: other
const (

	// Deprecated, use `db.client.connection.state` instead.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `db.client.connection.state`.
	//
	// Examples:
	// "idle",
	StateKey = attribute.Key("state")
)

// State returns an attribute KeyValue conforming to the "state"semantic conventions
// It represents the Deprecated, use `db.client.connection.state` instead
func State(val string) attribute.KeyValue {
	return StateKey.String(val)
}

// Enum values for state
var (

	// idle
	// Stability: experimental
	StateIdle = StateKey.String("idle")
	// used
	// Stability: experimental
	StateUsed = StateKey.String("used")
)

// Namespace: peer
const (

	// The [`service.name`] of the remote service. SHOULD be equal to the actual `service.name` resource attribute of the remote service if any.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "AuthTokenCache"
	//
	// [`service.name`]: /docs/resource/README.md#service
	PeerServiceKey = attribute.Key("peer.service")
)

// PeerService returns an attribute KeyValue conforming to the "peer.service"semantic conventions
// It represents the The [`service.name`] of the remote service. SHOULD be equal to the actual `service.name` resource attribute of the remote service if any
//
// [`service.name`]: /docs/resource/README.md#service
func PeerService(val string) attribute.KeyValue {
	return PeerServiceKey.String(val)
}

// Namespace: pool
const (

	// Deprecated, use `db.client.connection.pool.name` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `db.client.connection.pool.name`.
	//
	// Examples:
	// "myDataSource",
	PoolNameKey = attribute.Key("pool.name")
)

// PoolName returns an attribute KeyValue conforming to the "pool.name"semantic conventions
// It represents the Deprecated, use `db.client.connection.pool.name` instead
func PoolName(val string) attribute.KeyValue {
	return PoolNameKey.String(val)
}

// Namespace: process
const (

	// Length of the process.command_args array
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 4,
	//
	// Note: This field can be useful for querying or performing bucket analysis on how many arguments were provided to start a process. More arguments may be an indication of suspicious activity
	ProcessArgsCountKey = attribute.Key("process.args_count")
	// The command used to launch the process (i.e. the command name). On Linux based systems, can be set to the zeroth string in `proc/[pid]/cmdline`. On Windows, can be set to the first parameter extracted from `GetCommandLineW`.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "cmd/otelcol",
	ProcessCommandKey = attribute.Key("process.command")
	// All the command arguments (including the command/executable itself) as received by the process. On Linux-based systems (and some other Unixoid systems supporting procfs), can be set according to the list of null-delimited strings extracted from `proc/[pid]/cmdline`. For libc-based executables, this would be the full argv vector passed to `main`.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "cmd/otecol",
	// "--config=config.yaml",
	// ],
	ProcessCommandArgsKey = attribute.Key("process.command_args")
	// The full command used to launch the process as a single string representing the full command. On Windows, can be set to the result of `GetCommandLineW`. Do not set this if you have to assemble it just for monitoring; use `process.command_args` instead.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "C:\cmd\otecol --config="my directory\config.yaml"",
	ProcessCommandLineKey = attribute.Key("process.command_line")
	// Specifies whether the context switches for this data point were voluntary or involuntary.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	ProcessContextSwitchTypeKey = attribute.Key("process.context_switch_type")
	// Deprecated, use `cpu.mode` instead.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `cpu.mode`
	//
	// Examples: undefined
	ProcessCpuStateKey = attribute.Key("process.cpu.state")
	// The date and time the process was created, in ISO 8601 format.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2023-11-21T09:25:34.853Z",
	ProcessCreationTimeKey = attribute.Key("process.creation.time")
	// The name of the process executable. On Linux based systems, can be set to the `Name` in `proc/[pid]/status`. On Windows, can be set to the base name of `GetProcessImageFileNameW`.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "otelcol",
	ProcessExecutableNameKey = attribute.Key("process.executable.name")
	// The full path to the process executable. On Linux based systems, can be set to the target of `proc/[pid]/exe`. On Windows, can be set to the result of `GetProcessImageFileNameW`.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/usr/bin/cmd/otelcol",
	ProcessExecutablePathKey = attribute.Key("process.executable.path")
	// The exit code of the process.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 127,
	ProcessExitCodeKey = attribute.Key("process.exit.code")
	// The date and time the process exited, in ISO 8601 format.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2023-11-21T09:26:12.315Z",
	ProcessExitTimeKey = attribute.Key("process.exit.time")
	// The PID of the process's group leader. This is also the process group ID (PGID) of the process.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 23,
	ProcessGroupLeaderPidKey = attribute.Key("process.group_leader.pid")
	// Whether the process is connected to an interactive shell.
	//
	// Stability: Experimental
	// Type: boolean
	//
	// Examples: undefined
	ProcessInteractiveKey = attribute.Key("process.interactive")
	// The username of the user that owns the process.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "root",
	ProcessOwnerKey = attribute.Key("process.owner")
	// The type of page fault for this data point. Type `major` is for major/hard page faults, and `minor` is for minor/soft page faults.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	ProcessPagingFaultTypeKey = attribute.Key("process.paging.fault_type")
	// Parent Process identifier (PPID).
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 111,
	ProcessParentPidKey = attribute.Key("process.parent_pid")
	// Process identifier (PID).
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 1234,
	ProcessPidKey = attribute.Key("process.pid")
	// The real user ID (RUID) of the process.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 1000,
	ProcessRealUserIdKey = attribute.Key("process.real_user.id")
	// The username of the real user of the process.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "operator",
	ProcessRealUserNameKey = attribute.Key("process.real_user.name")
	// An additional description about the runtime of the process, for example a specific vendor customization of the runtime environment.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "Eclipse OpenJ9 Eclipse OpenJ9 VM openj9-0.21.0"
	ProcessRuntimeDescriptionKey = attribute.Key("process.runtime.description")
	// The name of the runtime of this process.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "OpenJDK Runtime Environment",
	ProcessRuntimeNameKey = attribute.Key("process.runtime.name")
	// The version of the runtime of this process, as returned by the runtime without modification.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "14.0.2"
	ProcessRuntimeVersionKey = attribute.Key("process.runtime.version")
	// The saved user ID (SUID) of the process.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 1002,
	ProcessSavedUserIdKey = attribute.Key("process.saved_user.id")
	// The username of the saved user.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "operator",
	ProcessSavedUserNameKey = attribute.Key("process.saved_user.name")
	// The PID of the process's session leader. This is also the session ID (SID) of the process.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 14,
	ProcessSessionLeaderPidKey = attribute.Key("process.session_leader.pid")
	// Process title (proctitle)
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "cat /etc/hostname",
	// "xfce4-session",
	// "bash",
	//
	// Note: In many Unix-like systems, process title (proctitle), is the string that represents the name or command line of a running process, displayed by system monitoring tools like ps, top, and htop
	ProcessTitleKey = attribute.Key("process.title")
	// The effective user ID (EUID) of the process.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 1001,
	ProcessUserIdKey = attribute.Key("process.user.id")
	// The username of the effective user of the process.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "root",
	ProcessUserNameKey = attribute.Key("process.user.name")
	// Virtual process identifier.
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 12,
	//
	// Note: The process ID within a PID namespace. This is not necessarily unique across all processes on the host but it is unique within the process namespace that the process exists within
	ProcessVpidKey = attribute.Key("process.vpid")
	// The working directory of the process.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/root",
	ProcessWorkingDirectoryKey = attribute.Key("process.working_directory")
)

// ProcessArgsCount returns an attribute KeyValue conforming to the "process.args_count"semantic conventions
// It represents the Length of the process.command_args array
func ProcessArgsCount(val int) attribute.KeyValue {
	return ProcessArgsCountKey.Int(val)
}

// ProcessCommand returns an attribute KeyValue conforming to the "process.command"semantic conventions
// It represents the The command used to launch the process (i.e. the command name). On Linux based systems, can be set to the zeroth string in `proc/[pid]/cmdline`. On Windows, can be set to the first parameter extracted from `GetCommandLineW`
func ProcessCommand(val string) attribute.KeyValue {
	return ProcessCommandKey.String(val)
}

// ProcessCommandArgs returns an attribute KeyValue conforming to the "process.command_args"semantic conventions
// It represents the All the command arguments (including the command/executable itself) as received by the process. On Linux-based systems (and some other Unixoid systems supporting procfs), can be set according to the list of null-delimited strings extracted from `proc/[pid]/cmdline`. For libc-based executables, this would be the full argv vector passed to `main`
func ProcessCommandArgs(val []string) attribute.KeyValue {
	return ProcessCommandArgsKey.StringSlice(val)
}

// ProcessCommandLine returns an attribute KeyValue conforming to the "process.command_line"semantic conventions
// It represents the The full command used to launch the process as a single string representing the full command. On Windows, can be set to the result of `GetCommandLineW`. Do not set this if you have to assemble it just for monitoring; use `process.command_args` instead
func ProcessCommandLine(val string) attribute.KeyValue {
	return ProcessCommandLineKey.String(val)
}

// ProcessContextSwitchType returns an attribute KeyValue conforming to the "process.context_switch_type"semantic conventions
// It represents the Specifies whether the context switches for this data point were voluntary or involuntary
func ProcessContextSwitchType(val string) attribute.KeyValue {
	return ProcessContextSwitchTypeKey.String(val)
}

// ProcessCpuState returns an attribute KeyValue conforming to the "process.cpu.state"semantic conventions
// It represents the Deprecated, use `cpu.mode` instead
func ProcessCpuState(val string) attribute.KeyValue {
	return ProcessCpuStateKey.String(val)
}

// ProcessCreationTime returns an attribute KeyValue conforming to the "process.creation.time"semantic conventions
// It represents the The date and time the process was created, in ISO 8601 format
func ProcessCreationTime(val string) attribute.KeyValue {
	return ProcessCreationTimeKey.String(val)
}

// ProcessExecutableName returns an attribute KeyValue conforming to the "process.executable.name"semantic conventions
// It represents the The name of the process executable. On Linux based systems, can be set to the `Name` in `proc/[pid]/status`. On Windows, can be set to the base name of `GetProcessImageFileNameW`
func ProcessExecutableName(val string) attribute.KeyValue {
	return ProcessExecutableNameKey.String(val)
}

// ProcessExecutablePath returns an attribute KeyValue conforming to the "process.executable.path"semantic conventions
// It represents the The full path to the process executable. On Linux based systems, can be set to the target of `proc/[pid]/exe`. On Windows, can be set to the result of `GetProcessImageFileNameW`
func ProcessExecutablePath(val string) attribute.KeyValue {
	return ProcessExecutablePathKey.String(val)
}

// ProcessExitCode returns an attribute KeyValue conforming to the "process.exit.code"semantic conventions
// It represents the The exit code of the process
func ProcessExitCode(val int) attribute.KeyValue {
	return ProcessExitCodeKey.Int(val)
}

// ProcessExitTime returns an attribute KeyValue conforming to the "process.exit.time"semantic conventions
// It represents the The date and time the process exited, in ISO 8601 format
func ProcessExitTime(val string) attribute.KeyValue {
	return ProcessExitTimeKey.String(val)
}

// ProcessGroupLeaderPid returns an attribute KeyValue conforming to the "process.group_leader.pid"semantic conventions
// It represents the The PID of the process's group leader. This is also the process group ID (PGID) of the process
func ProcessGroupLeaderPid(val int) attribute.KeyValue {
	return ProcessGroupLeaderPidKey.Int(val)
}

// ProcessInteractive returns an attribute KeyValue conforming to the "process.interactive"semantic conventions
// It represents the Whether the process is connected to an interactive shell
func ProcessInteractive(val bool) attribute.KeyValue {
	return ProcessInteractiveKey.Bool(val)
}

// ProcessOwner returns an attribute KeyValue conforming to the "process.owner"semantic conventions
// It represents the The username of the user that owns the process
func ProcessOwner(val string) attribute.KeyValue {
	return ProcessOwnerKey.String(val)
}

// ProcessPagingFaultType returns an attribute KeyValue conforming to the "process.paging.fault_type"semantic conventions
// It represents the The type of page fault for this data point. Type `major` is for major/hard page faults, and `minor` is for minor/soft page faults
func ProcessPagingFaultType(val string) attribute.KeyValue {
	return ProcessPagingFaultTypeKey.String(val)
}

// ProcessParentPid returns an attribute KeyValue conforming to the "process.parent_pid"semantic conventions
// It represents the Parent Process identifier (PPID)
func ProcessParentPid(val int) attribute.KeyValue {
	return ProcessParentPidKey.Int(val)
}

// ProcessPid returns an attribute KeyValue conforming to the "process.pid"semantic conventions
// It represents the Process identifier (PID)
func ProcessPid(val int) attribute.KeyValue {
	return ProcessPidKey.Int(val)
}

// ProcessRealUserId returns an attribute KeyValue conforming to the "process.real_user.id"semantic conventions
// It represents the The real user ID (RUID) of the process
func ProcessRealUserId(val int) attribute.KeyValue {
	return ProcessRealUserIdKey.Int(val)
}

// ProcessRealUserName returns an attribute KeyValue conforming to the "process.real_user.name"semantic conventions
// It represents the The username of the real user of the process
func ProcessRealUserName(val string) attribute.KeyValue {
	return ProcessRealUserNameKey.String(val)
}

// ProcessRuntimeDescription returns an attribute KeyValue conforming to the "process.runtime.description"semantic conventions
// It represents the An additional description about the runtime of the process, for example a specific vendor customization of the runtime environment
func ProcessRuntimeDescription(val string) attribute.KeyValue {
	return ProcessRuntimeDescriptionKey.String(val)
}

// ProcessRuntimeName returns an attribute KeyValue conforming to the "process.runtime.name"semantic conventions
// It represents the The name of the runtime of this process
func ProcessRuntimeName(val string) attribute.KeyValue {
	return ProcessRuntimeNameKey.String(val)
}

// ProcessRuntimeVersion returns an attribute KeyValue conforming to the "process.runtime.version"semantic conventions
// It represents the The version of the runtime of this process, as returned by the runtime without modification
func ProcessRuntimeVersion(val string) attribute.KeyValue {
	return ProcessRuntimeVersionKey.String(val)
}

// ProcessSavedUserId returns an attribute KeyValue conforming to the "process.saved_user.id"semantic conventions
// It represents the The saved user ID (SUID) of the process
func ProcessSavedUserId(val int) attribute.KeyValue {
	return ProcessSavedUserIdKey.Int(val)
}

// ProcessSavedUserName returns an attribute KeyValue conforming to the "process.saved_user.name"semantic conventions
// It represents the The username of the saved user
func ProcessSavedUserName(val string) attribute.KeyValue {
	return ProcessSavedUserNameKey.String(val)
}

// ProcessSessionLeaderPid returns an attribute KeyValue conforming to the "process.session_leader.pid"semantic conventions
// It represents the The PID of the process's session leader. This is also the session ID (SID) of the process
func ProcessSessionLeaderPid(val int) attribute.KeyValue {
	return ProcessSessionLeaderPidKey.Int(val)
}

// ProcessTitle returns an attribute KeyValue conforming to the "process.title"semantic conventions
// It represents the Process title (proctitle)
func ProcessTitle(val string) attribute.KeyValue {
	return ProcessTitleKey.String(val)
}

// ProcessUserId returns an attribute KeyValue conforming to the "process.user.id"semantic conventions
// It represents the The effective user ID (EUID) of the process
func ProcessUserId(val int) attribute.KeyValue {
	return ProcessUserIdKey.Int(val)
}

// ProcessUserName returns an attribute KeyValue conforming to the "process.user.name"semantic conventions
// It represents the The username of the effective user of the process
func ProcessUserName(val string) attribute.KeyValue {
	return ProcessUserNameKey.String(val)
}

// ProcessVpid returns an attribute KeyValue conforming to the "process.vpid"semantic conventions
// It represents the Virtual process identifier
func ProcessVpid(val int) attribute.KeyValue {
	return ProcessVpidKey.Int(val)
}

// ProcessWorkingDirectory returns an attribute KeyValue conforming to the "process.working_directory"semantic conventions
// It represents the The working directory of the process
func ProcessWorkingDirectory(val string) attribute.KeyValue {
	return ProcessWorkingDirectoryKey.String(val)
}

// Enum values for process.context_switch_type
var (

	// voluntary
	// Stability: experimental
	ProcessContextSwitchTypeVoluntary = ProcessContextSwitchTypeKey.String("voluntary")
	// involuntary
	// Stability: experimental
	ProcessContextSwitchTypeInvoluntary = ProcessContextSwitchTypeKey.String("involuntary")
)

// Enum values for process.cpu.state
var (

	// system
	// Stability: experimental
	ProcessCpuStateSystem = ProcessCpuStateKey.String("system")
	// user
	// Stability: experimental
	ProcessCpuStateUser = ProcessCpuStateKey.String("user")
	// wait
	// Stability: experimental
	ProcessCpuStateWait = ProcessCpuStateKey.String("wait")
)

// Enum values for process.paging.fault_type
var (

	// major
	// Stability: experimental
	ProcessPagingFaultTypeMajor = ProcessPagingFaultTypeKey.String("major")
	// minor
	// Stability: experimental
	ProcessPagingFaultTypeMinor = ProcessPagingFaultTypeKey.String("minor")
)

// Namespace: profile
const (

	// Describes the interpreter or compiler of a single frame.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "cpython",
	ProfileFrameTypeKey = attribute.Key("profile.frame.type")
)

// ProfileFrameType returns an attribute KeyValue conforming to the "profile.frame.type"semantic conventions
// It represents the Describes the interpreter or compiler of a single frame
func ProfileFrameType(val string) attribute.KeyValue {
	return ProfileFrameTypeKey.String(val)
}

// Enum values for profile.frame.type
var (

	// [.NET]
	//
	// Stability: experimental
	//
	// [.NET]: https://wikipedia.org/wiki/.NET
	ProfileFrameTypeDotnet = ProfileFrameTypeKey.String("dotnet")
	// [JVM]
	//
	// Stability: experimental
	//
	// [JVM]: https://wikipedia.org/wiki/Java_virtual_machine
	ProfileFrameTypeJvm = ProfileFrameTypeKey.String("jvm")
	// [Kernel]
	//
	// Stability: experimental
	//
	// [Kernel]: https://wikipedia.org/wiki/Kernel_(operating_system)
	ProfileFrameTypeKernel = ProfileFrameTypeKey.String("kernel")
	// [C], [C++], [Go], [Rust]
	//
	// Stability: experimental
	//
	// [C]: https://wikipedia.org/wiki/C_(programming_language)
	// [C++]: https://wikipedia.org/wiki/C%2B%2B
	// [Go]: https://wikipedia.org/wiki/Go_(programming_language)
	// [Rust]: https://wikipedia.org/wiki/Rust_(programming_language)
	ProfileFrameTypeNative = ProfileFrameTypeKey.String("native")
	// [Perl]
	//
	// Stability: experimental
	//
	// [Perl]: https://wikipedia.org/wiki/Perl
	ProfileFrameTypePerl = ProfileFrameTypeKey.String("perl")
	// [PHP]
	//
	// Stability: experimental
	//
	// [PHP]: https://wikipedia.org/wiki/PHP
	ProfileFrameTypePhp = ProfileFrameTypeKey.String("php")
	// [Python]
	//
	// Stability: experimental
	//
	// [Python]: https://wikipedia.org/wiki/Python_(programming_language)
	ProfileFrameTypeCpython = ProfileFrameTypeKey.String("cpython")
	// [Ruby]
	//
	// Stability: experimental
	//
	// [Ruby]: https://wikipedia.org/wiki/Ruby_(programming_language)
	ProfileFrameTypeRuby = ProfileFrameTypeKey.String("ruby")
	// [V8JS]
	//
	// Stability: experimental
	//
	// [V8JS]: https://wikipedia.org/wiki/V8_(JavaScript_engine)
	ProfileFrameTypeV8js = ProfileFrameTypeKey.String("v8js")
)

// Namespace: rpc
const (

	// The [error codes] of the Connect request. Error codes are always string values.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	//
	// [error codes]: https://connect.build/docs/protocol/#error-codes
	RpcConnectRpcErrorCodeKey = attribute.Key("rpc.connect_rpc.error_code")
	// The [numeric status code] of the gRPC request.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	//
	// [numeric status code]: https://github.com/grpc/grpc/blob/v1.33.2/doc/statuscodes.md
	RpcGrpcStatusCodeKey = attribute.Key("rpc.grpc.status_code")
	// `error.code` property of response if it is an error response.
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// -32700,
	// 100,
	RpcJsonrpcErrorCodeKey = attribute.Key("rpc.jsonrpc.error_code")
	// `error.message` property of response if it is an error response.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Parse error",
	// "User already exists",
	RpcJsonrpcErrorMessageKey = attribute.Key("rpc.jsonrpc.error_message")
	// `id` property of request or response. Since protocol allows id to be int, string, `null` or missing (for notifications), value is expected to be cast to string for simplicity. Use empty string in case of `null` value. Omit entirely if this is a notification.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "10",
	// "request-7",
	// "",
	RpcJsonrpcRequestIdKey = attribute.Key("rpc.jsonrpc.request_id")
	// Protocol version as in `jsonrpc` property of request/response. Since JSON-RPC 1.0 doesn't specify this, the value can be omitted.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2.0",
	// "1.0",
	RpcJsonrpcVersionKey = attribute.Key("rpc.jsonrpc.version")
	// Compressed size of the message in bytes.
	// Stability: Experimental
	// Type: int
	//
	// Examples: undefined
	RpcMessageCompressedSizeKey = attribute.Key("rpc.message.compressed_size")
	// MUST be calculated as two different counters starting from `1` one for sent messages and one for received message.
	// Stability: Experimental
	// Type: int
	//
	// Examples: undefined
	// Note: This way we guarantee that the values will be consistent between different implementations
	RpcMessageIdKey = attribute.Key("rpc.message.id")
	// Whether this is a received or sent message.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	RpcMessageTypeKey = attribute.Key("rpc.message.type")
	// Uncompressed size of the message in bytes.
	// Stability: Experimental
	// Type: int
	//
	// Examples: undefined
	RpcMessageUncompressedSizeKey = attribute.Key("rpc.message.uncompressed_size")
	// The name of the (logical) method being called, must be equal to the $method part in the span name.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "exampleMethod"
	// Note: This is the logical name of the method from the RPC interface perspective, which can be different from the name of any implementing method/function. The `code.function` attribute may be used to store the latter (e.g., method actually executing the call on the server side, RPC client stub method on the client side)
	RpcMethodKey = attribute.Key("rpc.method")
	// The full (logical) name of the service being called, including its package name, if applicable.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "myservice.EchoService"
	// Note: This is the logical name of the service from the RPC interface perspective, which can be different from the name of any implementing class. The `code.namespace` attribute may be used to store the latter (despite the attribute name, it may include a class name; e.g., class with method actually executing the call on the server side, RPC client stub class on the client side)
	RpcServiceKey = attribute.Key("rpc.service")
	// A string identifying the remoting system. See below for a list of well-known identifiers.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	RpcSystemKey = attribute.Key("rpc.system")
)

// RpcConnectRpcErrorCode returns an attribute KeyValue conforming to the "rpc.connect_rpc.error_code"semantic conventions
// It represents the The [error codes] of the Connect request. Error codes are always string values
//
// [error codes]: https://connect.build/docs/protocol/#error-codes
func RpcConnectRpcErrorCode(val string) attribute.KeyValue {
	return RpcConnectRpcErrorCodeKey.String(val)
}

// RpcGrpcStatusCode returns an attribute KeyValue conforming to the "rpc.grpc.status_code"semantic conventions
// It represents the The [numeric status code] of the gRPC request
//
// [numeric status code]: https://github.com/grpc/grpc/blob/v1.33.2/doc/statuscodes.md
func RpcGrpcStatusCode(val int) attribute.KeyValue {
	return RpcGrpcStatusCodeKey.Int(val)
}

// RpcJsonrpcErrorCode returns an attribute KeyValue conforming to the "rpc.jsonrpc.error_code"semantic conventions
// It represents the `error.code` property of response if it is an error response
func RpcJsonrpcErrorCode(val int) attribute.KeyValue {
	return RpcJsonrpcErrorCodeKey.Int(val)
}

// RpcJsonrpcErrorMessage returns an attribute KeyValue conforming to the "rpc.jsonrpc.error_message"semantic conventions
// It represents the `error.message` property of response if it is an error response
func RpcJsonrpcErrorMessage(val string) attribute.KeyValue {
	return RpcJsonrpcErrorMessageKey.String(val)
}

// RpcJsonrpcRequestId returns an attribute KeyValue conforming to the "rpc.jsonrpc.request_id"semantic conventions
// It represents the `id` property of request or response. Since protocol allows id to be int, string, `null` or missing (for notifications), value is expected to be cast to string for simplicity. Use empty string in case of `null` value. Omit entirely if this is a notification
func RpcJsonrpcRequestId(val string) attribute.KeyValue {
	return RpcJsonrpcRequestIdKey.String(val)
}

// RpcJsonrpcVersion returns an attribute KeyValue conforming to the "rpc.jsonrpc.version"semantic conventions
// It represents the Protocol version as in `jsonrpc` property of request/response. Since JSON-RPC 1.0 doesn't specify this, the value can be omitted
func RpcJsonrpcVersion(val string) attribute.KeyValue {
	return RpcJsonrpcVersionKey.String(val)
}

// RpcMessageCompressedSize returns an attribute KeyValue conforming to the "rpc.message.compressed_size"semantic conventions
// It represents the Compressed size of the message in bytes
func RpcMessageCompressedSize(val int) attribute.KeyValue {
	return RpcMessageCompressedSizeKey.Int(val)
}

// RpcMessageId returns an attribute KeyValue conforming to the "rpc.message.id"semantic conventions
// It represents the MUST be calculated as two different counters starting from `1` one for sent messages and one for received message
func RpcMessageId(val int) attribute.KeyValue {
	return RpcMessageIdKey.Int(val)
}

// RpcMessageType returns an attribute KeyValue conforming to the "rpc.message.type"semantic conventions
// It represents the Whether this is a received or sent message
func RpcMessageType(val string) attribute.KeyValue {
	return RpcMessageTypeKey.String(val)
}

// RpcMessageUncompressedSize returns an attribute KeyValue conforming to the "rpc.message.uncompressed_size"semantic conventions
// It represents the Uncompressed size of the message in bytes
func RpcMessageUncompressedSize(val int) attribute.KeyValue {
	return RpcMessageUncompressedSizeKey.Int(val)
}

// RpcMethod returns an attribute KeyValue conforming to the "rpc.method"semantic conventions
// It represents the The name of the (logical) method being called, must be equal to the $method part in the span name
func RpcMethod(val string) attribute.KeyValue {
	return RpcMethodKey.String(val)
}

// RpcService returns an attribute KeyValue conforming to the "rpc.service"semantic conventions
// It represents the The full (logical) name of the service being called, including its package name, if applicable
func RpcService(val string) attribute.KeyValue {
	return RpcServiceKey.String(val)
}

// RpcSystem returns an attribute KeyValue conforming to the "rpc.system"semantic conventions
// It represents the A string identifying the remoting system. See below for a list of well-known identifiers
func RpcSystem(val string) attribute.KeyValue {
	return RpcSystemKey.String(val)
}

// Enum values for rpc.connect_rpc.error_code
var (

	// cancelled
	// Stability: experimental
	RpcConnectRpcErrorCodeCancelled = RpcConnectRpcErrorCodeKey.String("cancelled")
	// unknown
	// Stability: experimental
	RpcConnectRpcErrorCodeUnknown = RpcConnectRpcErrorCodeKey.String("unknown")
	// invalid_argument
	// Stability: experimental
	RpcConnectRpcErrorCodeInvalidArgument = RpcConnectRpcErrorCodeKey.String("invalid_argument")
	// deadline_exceeded
	// Stability: experimental
	RpcConnectRpcErrorCodeDeadlineExceeded = RpcConnectRpcErrorCodeKey.String("deadline_exceeded")
	// not_found
	// Stability: experimental
	RpcConnectRpcErrorCodeNotFound = RpcConnectRpcErrorCodeKey.String("not_found")
	// already_exists
	// Stability: experimental
	RpcConnectRpcErrorCodeAlreadyExists = RpcConnectRpcErrorCodeKey.String("already_exists")
	// permission_denied
	// Stability: experimental
	RpcConnectRpcErrorCodePermissionDenied = RpcConnectRpcErrorCodeKey.String("permission_denied")
	// resource_exhausted
	// Stability: experimental
	RpcConnectRpcErrorCodeResourceExhausted = RpcConnectRpcErrorCodeKey.String("resource_exhausted")
	// failed_precondition
	// Stability: experimental
	RpcConnectRpcErrorCodeFailedPrecondition = RpcConnectRpcErrorCodeKey.String("failed_precondition")
	// aborted
	// Stability: experimental
	RpcConnectRpcErrorCodeAborted = RpcConnectRpcErrorCodeKey.String("aborted")
	// out_of_range
	// Stability: experimental
	RpcConnectRpcErrorCodeOutOfRange = RpcConnectRpcErrorCodeKey.String("out_of_range")
	// unimplemented
	// Stability: experimental
	RpcConnectRpcErrorCodeUnimplemented = RpcConnectRpcErrorCodeKey.String("unimplemented")
	// internal
	// Stability: experimental
	RpcConnectRpcErrorCodeInternal = RpcConnectRpcErrorCodeKey.String("internal")
	// unavailable
	// Stability: experimental
	RpcConnectRpcErrorCodeUnavailable = RpcConnectRpcErrorCodeKey.String("unavailable")
	// data_loss
	// Stability: experimental
	RpcConnectRpcErrorCodeDataLoss = RpcConnectRpcErrorCodeKey.String("data_loss")
	// unauthenticated
	// Stability: experimental
	RpcConnectRpcErrorCodeUnauthenticated = RpcConnectRpcErrorCodeKey.String("unauthenticated")
)

// Enum values for rpc.grpc.status_code
var (

	// OK
	// Stability: experimental
	RpcGrpcStatusCodeOk = RpcGrpcStatusCodeKey.Int(0)
	// CANCELLED
	// Stability: experimental
	RpcGrpcStatusCodeCancelled = RpcGrpcStatusCodeKey.Int(1)
	// UNKNOWN
	// Stability: experimental
	RpcGrpcStatusCodeUnknown = RpcGrpcStatusCodeKey.Int(2)
	// INVALID_ARGUMENT
	// Stability: experimental
	RpcGrpcStatusCodeInvalidArgument = RpcGrpcStatusCodeKey.Int(3)
	// DEADLINE_EXCEEDED
	// Stability: experimental
	RpcGrpcStatusCodeDeadlineExceeded = RpcGrpcStatusCodeKey.Int(4)
	// NOT_FOUND
	// Stability: experimental
	RpcGrpcStatusCodeNotFound = RpcGrpcStatusCodeKey.Int(5)
	// ALREADY_EXISTS
	// Stability: experimental
	RpcGrpcStatusCodeAlreadyExists = RpcGrpcStatusCodeKey.Int(6)
	// PERMISSION_DENIED
	// Stability: experimental
	RpcGrpcStatusCodePermissionDenied = RpcGrpcStatusCodeKey.Int(7)
	// RESOURCE_EXHAUSTED
	// Stability: experimental
	RpcGrpcStatusCodeResourceExhausted = RpcGrpcStatusCodeKey.Int(8)
	// FAILED_PRECONDITION
	// Stability: experimental
	RpcGrpcStatusCodeFailedPrecondition = RpcGrpcStatusCodeKey.Int(9)
	// ABORTED
	// Stability: experimental
	RpcGrpcStatusCodeAborted = RpcGrpcStatusCodeKey.Int(10)
	// OUT_OF_RANGE
	// Stability: experimental
	RpcGrpcStatusCodeOutOfRange = RpcGrpcStatusCodeKey.Int(11)
	// UNIMPLEMENTED
	// Stability: experimental
	RpcGrpcStatusCodeUnimplemented = RpcGrpcStatusCodeKey.Int(12)
	// INTERNAL
	// Stability: experimental
	RpcGrpcStatusCodeInternal = RpcGrpcStatusCodeKey.Int(13)
	// UNAVAILABLE
	// Stability: experimental
	RpcGrpcStatusCodeUnavailable = RpcGrpcStatusCodeKey.Int(14)
	// DATA_LOSS
	// Stability: experimental
	RpcGrpcStatusCodeDataLoss = RpcGrpcStatusCodeKey.Int(15)
	// UNAUTHENTICATED
	// Stability: experimental
	RpcGrpcStatusCodeUnauthenticated = RpcGrpcStatusCodeKey.Int(16)
)

// Enum values for rpc.message.type
var (

	// sent
	// Stability: experimental
	RpcMessageTypeSent = RpcMessageTypeKey.String("SENT")
	// received
	// Stability: experimental
	RpcMessageTypeReceived = RpcMessageTypeKey.String("RECEIVED")
)

// Enum values for rpc.system
var (

	// gRPC
	// Stability: experimental
	RpcSystemGrpc = RpcSystemKey.String("grpc")
	// Java RMI
	// Stability: experimental
	RpcSystemJavaRmi = RpcSystemKey.String("java_rmi")
	// .NET WCF
	// Stability: experimental
	RpcSystemDotnetWcf = RpcSystemKey.String("dotnet_wcf")
	// Apache Dubbo
	// Stability: experimental
	RpcSystemApacheDubbo = RpcSystemKey.String("apache_dubbo")
	// Connect RPC
	// Stability: experimental
	RpcSystemConnectRpc = RpcSystemKey.String("connect_rpc")
)

// Namespace: server
const (

	// Server domain name if available without reverse DNS lookup; otherwise, IP address or Unix domain socket name.
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "example.com",
	// "10.1.2.80",
	// "/tmp/my.sock",
	//
	// Note: When observed from the client side, and when communicating through an intermediary, `server.address` SHOULD represent the server address behind any intermediaries, for example proxies, if it's available
	ServerAddressKey = attribute.Key("server.address")
	// Server port number.
	// Stability: Stable
	// Type: int
	//
	// Examples:
	// 80,
	// 8080,
	// 443,
	//
	// Note: When observed from the client side, and when communicating through an intermediary, `server.port` SHOULD represent the server port behind any intermediaries, for example proxies, if it's available
	ServerPortKey = attribute.Key("server.port")
)

// ServerAddress returns an attribute KeyValue conforming to the "server.address"semantic conventions
// It represents the Server domain name if available without reverse DNS lookup; otherwise, IP address or Unix domain socket name
func ServerAddress(val string) attribute.KeyValue {
	return ServerAddressKey.String(val)
}

// ServerPort returns an attribute KeyValue conforming to the "server.port"semantic conventions
// It represents the Server port number
func ServerPort(val int) attribute.KeyValue {
	return ServerPortKey.Int(val)
}

// Namespace: service
const (

	// The string ID of the service instance.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "627cc493-f310-47de-96bd-71410b7dec09",
	//
	// Note: MUST be unique for each instance of the same `service.namespace,service.name` pair (in other words
	// `service.namespace,service.name,service.instance.id` triplet MUST be globally unique). The ID helps to
	// distinguish instances of the same service that exist at the same time (e.g. instances of a horizontally scaled
	// service).
	//
	// Implementations, such as SDKs, are recommended to generate a random Version 1 or Version 4 [RFC
	// 4122] UUID, but are free to use an inherent unique ID as the source of
	// this value if stability is desirable. In that case, the ID SHOULD be used as source of a UUID Version 5 and
	// SHOULD use the following UUID as the namespace: `4d63009a-8d0f-11ee-aad7-4c796ed8e320`.
	//
	// UUIDs are typically recommended, as only an opaque value for the purposes of identifying a service instance is
	// needed. Similar to what can be seen in the man page for the
	// [`/etc/machine-id`] file, the underlying
	// data, such as pod name and namespace should be treated as confidential, being the user's choice to expose it
	// or not via another resource attribute.
	//
	// For applications running behind an application server (like unicorn), we do not recommend using one identifier
	// for all processes participating in the application. Instead, it's recommended each division (e.g. a worker
	// thread in unicorn) to have its own instance.id.
	//
	// It's not recommended for a Collector to set `service.instance.id` if it can't unambiguously determine the
	// service instance that is generating that telemetry. For instance, creating an UUID based on `pod.name` will
	// likely be wrong, as the Collector might not know from which container within that pod the telemetry originated.
	// However, Collectors can set the `service.instance.id` if they can unambiguously determine the service instance
	// for that telemetry. This is typically the case for scraping receivers, as they know the target address and
	// port
	//
	// [RFC
	// 4122]: https://www.ietf.org/rfc/rfc4122.txt
	// [`/etc/machine-id`]: https://www.freedesktop.org/software/systemd/man/machine-id.html
	ServiceInstanceIdKey = attribute.Key("service.instance.id")
	// Logical name of the service.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "shoppingcart",
	//
	// Note: MUST be the same for all instances of horizontally scaled services. If the value was not specified, SDKs MUST fallback to `unknown_service:` concatenated with [`process.executable.name`], e.g. `unknown_service:bash`. If `process.executable.name` is not available, the value MUST be set to `unknown_service`
	//
	// [`process.executable.name`]: process.md
	ServiceNameKey = attribute.Key("service.name")
	// A namespace for `service.name`.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Shop",
	//
	// Note: A string value having a meaning that helps to distinguish a group of services, for example the team name that owns a group of services. `service.name` is expected to be unique within the same namespace. If `service.namespace` is not specified in the Resource then `service.name` is expected to be unique for all services that have no explicit namespace defined (so the empty/unspecified namespace is simply one more valid namespace). Zero-length namespace string is assumed equal to unspecified namespace
	ServiceNamespaceKey = attribute.Key("service.namespace")
	// The version string of the service API or implementation. The format is not defined by these conventions.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "2.0.0",
	// "a01dbef8a",
	ServiceVersionKey = attribute.Key("service.version")
)

// ServiceInstanceId returns an attribute KeyValue conforming to the "service.instance.id"semantic conventions
// It represents the The string ID of the service instance
func ServiceInstanceId(val string) attribute.KeyValue {
	return ServiceInstanceIdKey.String(val)
}

// ServiceName returns an attribute KeyValue conforming to the "service.name"semantic conventions
// It represents the Logical name of the service
func ServiceName(val string) attribute.KeyValue {
	return ServiceNameKey.String(val)
}

// ServiceNamespace returns an attribute KeyValue conforming to the "service.namespace"semantic conventions
// It represents the A namespace for `service.name`
func ServiceNamespace(val string) attribute.KeyValue {
	return ServiceNamespaceKey.String(val)
}

// ServiceVersion returns an attribute KeyValue conforming to the "service.version"semantic conventions
// It represents the The version string of the service API or implementation. The format is not defined by these conventions
func ServiceVersion(val string) attribute.KeyValue {
	return ServiceVersionKey.String(val)
}

// Namespace: session
const (

	// A unique id to identify a session.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "00112233-4455-6677-8899-aabbccddeeff"
	SessionIdKey = attribute.Key("session.id")
	// The previous `session.id` for this user, when known.
	// Stability: Experimental
	// Type: string
	//
	// Examples: "00112233-4455-6677-8899-aabbccddeeff"
	SessionPreviousIdKey = attribute.Key("session.previous_id")
)

// SessionId returns an attribute KeyValue conforming to the "session.id"semantic conventions
// It represents the A unique id to identify a session
func SessionId(val string) attribute.KeyValue {
	return SessionIdKey.String(val)
}

// SessionPreviousId returns an attribute KeyValue conforming to the "session.previous_id"semantic conventions
// It represents the The previous `session.id` for this user, when known
func SessionPreviousId(val string) attribute.KeyValue {
	return SessionPreviousIdKey.String(val)
}

// Namespace: signalr
const (

	// SignalR HTTP connection closure status.
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "app_shutdown",
	// "timeout",
	SignalrConnectionStatusKey = attribute.Key("signalr.connection.status")
	// [SignalR transport type]
	// Stability: Stable
	// Type: Enum
	//
	// Examples:
	// "web_sockets",
	// "long_polling",
	//
	// [SignalR transport type]: https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md
	SignalrTransportKey = attribute.Key("signalr.transport")
)

// SignalrConnectionStatus returns an attribute KeyValue conforming to the "signalr.connection.status"semantic conventions
// It represents the SignalR HTTP connection closure status
func SignalrConnectionStatus(val string) attribute.KeyValue {
	return SignalrConnectionStatusKey.String(val)
}

// SignalrTransport returns an attribute KeyValue conforming to the "signalr.transport"semantic conventions
// It represents the [SignalR transport type]
//
// [SignalR transport type]: https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md
func SignalrTransport(val string) attribute.KeyValue {
	return SignalrTransportKey.String(val)
}

// Enum values for signalr.connection.status
var (

	// The connection was closed normally.
	// Stability: stable
	SignalrConnectionStatusNormalClosure = SignalrConnectionStatusKey.String("normal_closure")
	// The connection was closed due to a timeout.
	// Stability: stable
	SignalrConnectionStatusTimeout = SignalrConnectionStatusKey.String("timeout")
	// The connection was closed because the app is shutting down.
	// Stability: stable
	SignalrConnectionStatusAppShutdown = SignalrConnectionStatusKey.String("app_shutdown")
)

// Enum values for signalr.transport
var (

	// ServerSentEvents protocol
	// Stability: stable
	SignalrTransportServerSentEvents = SignalrTransportKey.String("server_sent_events")
	// LongPolling protocol
	// Stability: stable
	SignalrTransportLongPolling = SignalrTransportKey.String("long_polling")
	// WebSockets protocol
	// Stability: stable
	SignalrTransportWebSockets = SignalrTransportKey.String("web_sockets")
)

// Namespace: source
const (

	// Source address - domain name if available without reverse DNS lookup; otherwise, IP address or Unix domain socket name.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "source.example.com",
	// "10.1.2.80",
	// "/tmp/my.sock",
	//
	// Note: When observed from the destination side, and when communicating through an intermediary, `source.address` SHOULD represent the source address behind any intermediaries, for example proxies, if it's available
	SourceAddressKey = attribute.Key("source.address")
	// Source port number
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 3389,
	// 2888,
	SourcePortKey = attribute.Key("source.port")
)

// SourceAddress returns an attribute KeyValue conforming to the "source.address"semantic conventions
// It represents the Source address - domain name if available without reverse DNS lookup; otherwise, IP address or Unix domain socket name
func SourceAddress(val string) attribute.KeyValue {
	return SourceAddressKey.String(val)
}

// SourcePort returns an attribute KeyValue conforming to the "source.port"semantic conventions
// It represents the Source port number
func SourcePort(val int) attribute.KeyValue {
	return SourcePortKey.Int(val)
}

// Namespace: system
const (

	// The logical CPU number [0..n-1]
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 1,
	SystemCpuLogicalNumberKey = attribute.Key("system.cpu.logical_number")
	// Deprecated, use `cpu.mode` instead.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `cpu.mode`
	//
	// Examples:
	// "idle",
	// "interrupt",
	SystemCpuStateKey = attribute.Key("system.cpu.state")
	// The device identifier
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "(identifier)",
	SystemDeviceKey = attribute.Key("system.device")
	// The filesystem mode
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "rw, ro",
	SystemFilesystemModeKey = attribute.Key("system.filesystem.mode")
	// The filesystem mount path
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/mnt/data",
	SystemFilesystemMountpointKey = attribute.Key("system.filesystem.mountpoint")
	// The filesystem state
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "used",
	SystemFilesystemStateKey = attribute.Key("system.filesystem.state")
	// The filesystem type
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "ext4",
	SystemFilesystemTypeKey = attribute.Key("system.filesystem.type")
	// The memory state
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "free",
	// "cached",
	SystemMemoryStateKey = attribute.Key("system.memory.state")
	// A stateless protocol MUST NOT set this attribute
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "close_wait",
	SystemNetworkStateKey = attribute.Key("system.network.state")
	// The paging access direction
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "in",
	SystemPagingDirectionKey = attribute.Key("system.paging.direction")
	// The memory paging state
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "free",
	SystemPagingStateKey = attribute.Key("system.paging.state")
	// The memory paging type
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "minor",
	SystemPagingTypeKey = attribute.Key("system.paging.type")
	// The process state, e.g., [Linux Process State Codes]
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "running",
	//
	// [Linux Process State Codes]: https://man7.org/linux/man-pages/man1/ps.1.html#PROCESS_STATE_CODES
	SystemProcessStatusKey = attribute.Key("system.process.status")
	// Deprecated, use `system.process.status` instead.
	// Stability: Experimental
	// Type: Enum
	// Deprecated: Replaced by `system.process.status`.
	//
	// Examples:
	// "running",
	SystemProcessesStatusKey = attribute.Key("system.processes.status")
)

// SystemCpuLogicalNumber returns an attribute KeyValue conforming to the "system.cpu.logical_number"semantic conventions
// It represents the The logical CPU number [0..n-1]
func SystemCpuLogicalNumber(val int) attribute.KeyValue {
	return SystemCpuLogicalNumberKey.Int(val)
}

// SystemCpuState returns an attribute KeyValue conforming to the "system.cpu.state"semantic conventions
// It represents the Deprecated, use `cpu.mode` instead
func SystemCpuState(val string) attribute.KeyValue {
	return SystemCpuStateKey.String(val)
}

// SystemDevice returns an attribute KeyValue conforming to the "system.device"semantic conventions
// It represents the The device identifier
func SystemDevice(val string) attribute.KeyValue {
	return SystemDeviceKey.String(val)
}

// SystemFilesystemMode returns an attribute KeyValue conforming to the "system.filesystem.mode"semantic conventions
// It represents the The filesystem mode
func SystemFilesystemMode(val string) attribute.KeyValue {
	return SystemFilesystemModeKey.String(val)
}

// SystemFilesystemMountpoint returns an attribute KeyValue conforming to the "system.filesystem.mountpoint"semantic conventions
// It represents the The filesystem mount path
func SystemFilesystemMountpoint(val string) attribute.KeyValue {
	return SystemFilesystemMountpointKey.String(val)
}

// SystemFilesystemState returns an attribute KeyValue conforming to the "system.filesystem.state"semantic conventions
// It represents the The filesystem state
func SystemFilesystemState(val string) attribute.KeyValue {
	return SystemFilesystemStateKey.String(val)
}

// SystemFilesystemType returns an attribute KeyValue conforming to the "system.filesystem.type"semantic conventions
// It represents the The filesystem type
func SystemFilesystemType(val string) attribute.KeyValue {
	return SystemFilesystemTypeKey.String(val)
}

// SystemMemoryState returns an attribute KeyValue conforming to the "system.memory.state"semantic conventions
// It represents the The memory state
func SystemMemoryState(val string) attribute.KeyValue {
	return SystemMemoryStateKey.String(val)
}

// SystemNetworkState returns an attribute KeyValue conforming to the "system.network.state"semantic conventions
// It represents the A stateless protocol MUST NOT set this attribute
func SystemNetworkState(val string) attribute.KeyValue {
	return SystemNetworkStateKey.String(val)
}

// SystemPagingDirection returns an attribute KeyValue conforming to the "system.paging.direction"semantic conventions
// It represents the The paging access direction
func SystemPagingDirection(val string) attribute.KeyValue {
	return SystemPagingDirectionKey.String(val)
}

// SystemPagingState returns an attribute KeyValue conforming to the "system.paging.state"semantic conventions
// It represents the The memory paging state
func SystemPagingState(val string) attribute.KeyValue {
	return SystemPagingStateKey.String(val)
}

// SystemPagingType returns an attribute KeyValue conforming to the "system.paging.type"semantic conventions
// It represents the The memory paging type
func SystemPagingType(val string) attribute.KeyValue {
	return SystemPagingTypeKey.String(val)
}

// SystemProcessStatus returns an attribute KeyValue conforming to the "system.process.status"semantic conventions
// It represents the The process state, e.g., [Linux Process State Codes]
//
// [Linux Process State Codes]: https://man7.org/linux/man-pages/man1/ps.1.html#PROCESS_STATE_CODES
func SystemProcessStatus(val string) attribute.KeyValue {
	return SystemProcessStatusKey.String(val)
}

// SystemProcessesStatus returns an attribute KeyValue conforming to the "system.processes.status"semantic conventions
// It represents the Deprecated, use `system.process.status` instead
func SystemProcessesStatus(val string) attribute.KeyValue {
	return SystemProcessesStatusKey.String(val)
}

// Enum values for system.cpu.state
var (

	// user
	// Stability: experimental
	SystemCpuStateUser = SystemCpuStateKey.String("user")
	// system
	// Stability: experimental
	SystemCpuStateSystem = SystemCpuStateKey.String("system")
	// nice
	// Stability: experimental
	SystemCpuStateNice = SystemCpuStateKey.String("nice")
	// idle
	// Stability: experimental
	SystemCpuStateIdle = SystemCpuStateKey.String("idle")
	// iowait
	// Stability: experimental
	SystemCpuStateIowait = SystemCpuStateKey.String("iowait")
	// interrupt
	// Stability: experimental
	SystemCpuStateInterrupt = SystemCpuStateKey.String("interrupt")
	// steal
	// Stability: experimental
	SystemCpuStateSteal = SystemCpuStateKey.String("steal")
)

// Enum values for system.filesystem.state
var (

	// used
	// Stability: experimental
	SystemFilesystemStateUsed = SystemFilesystemStateKey.String("used")
	// free
	// Stability: experimental
	SystemFilesystemStateFree = SystemFilesystemStateKey.String("free")
	// reserved
	// Stability: experimental
	SystemFilesystemStateReserved = SystemFilesystemStateKey.String("reserved")
)

// Enum values for system.filesystem.type
var (

	// fat32
	// Stability: experimental
	SystemFilesystemTypeFat32 = SystemFilesystemTypeKey.String("fat32")
	// exfat
	// Stability: experimental
	SystemFilesystemTypeExfat = SystemFilesystemTypeKey.String("exfat")
	// ntfs
	// Stability: experimental
	SystemFilesystemTypeNtfs = SystemFilesystemTypeKey.String("ntfs")
	// refs
	// Stability: experimental
	SystemFilesystemTypeRefs = SystemFilesystemTypeKey.String("refs")
	// hfsplus
	// Stability: experimental
	SystemFilesystemTypeHfsplus = SystemFilesystemTypeKey.String("hfsplus")
	// ext4
	// Stability: experimental
	SystemFilesystemTypeExt4 = SystemFilesystemTypeKey.String("ext4")
)

// Enum values for system.memory.state
var (

	// used
	// Stability: experimental
	SystemMemoryStateUsed = SystemMemoryStateKey.String("used")
	// free
	// Stability: experimental
	SystemMemoryStateFree = SystemMemoryStateKey.String("free")
	// shared
	// Stability: experimental// Removed, report shared memory usage with `metric.system.memory.shared` metric
	SystemMemoryStateShared = SystemMemoryStateKey.String("shared")
	// buffers
	// Stability: experimental
	SystemMemoryStateBuffers = SystemMemoryStateKey.String("buffers")
	// cached
	// Stability: experimental
	SystemMemoryStateCached = SystemMemoryStateKey.String("cached")
)

// Enum values for system.network.state
var (

	// close
	// Stability: experimental
	SystemNetworkStateClose = SystemNetworkStateKey.String("close")
	// close_wait
	// Stability: experimental
	SystemNetworkStateCloseWait = SystemNetworkStateKey.String("close_wait")
	// closing
	// Stability: experimental
	SystemNetworkStateClosing = SystemNetworkStateKey.String("closing")
	// delete
	// Stability: experimental
	SystemNetworkStateDelete = SystemNetworkStateKey.String("delete")
	// established
	// Stability: experimental
	SystemNetworkStateEstablished = SystemNetworkStateKey.String("established")
	// fin_wait_1
	// Stability: experimental
	SystemNetworkStateFinWait1 = SystemNetworkStateKey.String("fin_wait_1")
	// fin_wait_2
	// Stability: experimental
	SystemNetworkStateFinWait2 = SystemNetworkStateKey.String("fin_wait_2")
	// last_ack
	// Stability: experimental
	SystemNetworkStateLastAck = SystemNetworkStateKey.String("last_ack")
	// listen
	// Stability: experimental
	SystemNetworkStateListen = SystemNetworkStateKey.String("listen")
	// syn_recv
	// Stability: experimental
	SystemNetworkStateSynRecv = SystemNetworkStateKey.String("syn_recv")
	// syn_sent
	// Stability: experimental
	SystemNetworkStateSynSent = SystemNetworkStateKey.String("syn_sent")
	// time_wait
	// Stability: experimental
	SystemNetworkStateTimeWait = SystemNetworkStateKey.String("time_wait")
)

// Enum values for system.paging.direction
var (

	// in
	// Stability: experimental
	SystemPagingDirectionIn = SystemPagingDirectionKey.String("in")
	// out
	// Stability: experimental
	SystemPagingDirectionOut = SystemPagingDirectionKey.String("out")
)

// Enum values for system.paging.state
var (

	// used
	// Stability: experimental
	SystemPagingStateUsed = SystemPagingStateKey.String("used")
	// free
	// Stability: experimental
	SystemPagingStateFree = SystemPagingStateKey.String("free")
)

// Enum values for system.paging.type
var (

	// major
	// Stability: experimental
	SystemPagingTypeMajor = SystemPagingTypeKey.String("major")
	// minor
	// Stability: experimental
	SystemPagingTypeMinor = SystemPagingTypeKey.String("minor")
)

// Enum values for system.process.status
var (

	// running
	// Stability: experimental
	SystemProcessStatusRunning = SystemProcessStatusKey.String("running")
	// sleeping
	// Stability: experimental
	SystemProcessStatusSleeping = SystemProcessStatusKey.String("sleeping")
	// stopped
	// Stability: experimental
	SystemProcessStatusStopped = SystemProcessStatusKey.String("stopped")
	// defunct
	// Stability: experimental
	SystemProcessStatusDefunct = SystemProcessStatusKey.String("defunct")
)

// Enum values for system.processes.status
var (

	// running
	// Stability: experimental
	SystemProcessesStatusRunning = SystemProcessesStatusKey.String("running")
	// sleeping
	// Stability: experimental
	SystemProcessesStatusSleeping = SystemProcessesStatusKey.String("sleeping")
	// stopped
	// Stability: experimental
	SystemProcessesStatusStopped = SystemProcessesStatusKey.String("stopped")
	// defunct
	// Stability: experimental
	SystemProcessesStatusDefunct = SystemProcessesStatusKey.String("defunct")
)

// Namespace: telemetry
const (

	// The name of the auto instrumentation agent or distribution, if used.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "parts-unlimited-java",
	//
	// Note: Official auto instrumentation agents and distributions SHOULD set the `telemetry.distro.name` attribute to
	// a string starting with `opentelemetry-`, e.g. `opentelemetry-java-instrumentation`
	TelemetryDistroNameKey = attribute.Key("telemetry.distro.name")
	// The version string of the auto instrumentation agent or distribution, if used.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "1.2.3",
	TelemetryDistroVersionKey = attribute.Key("telemetry.distro.version")
	// The language of the telemetry SDK.
	//
	// Stability: Stable
	// Type: Enum
	//
	// Examples: undefined
	TelemetrySdkLanguageKey = attribute.Key("telemetry.sdk.language")
	// The name of the telemetry SDK as defined above.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "opentelemetry",
	//
	// Note: The OpenTelemetry SDK MUST set the `telemetry.sdk.name` attribute to `opentelemetry`.
	// If another SDK, like a fork or a vendor-provided implementation, is used, this SDK MUST set the
	// `telemetry.sdk.name` attribute to the fully-qualified class or module name of this SDK's main entry point
	// or another suitable identifier depending on the language.
	// The identifier `opentelemetry` is reserved and MUST NOT be used in this case.
	// All custom identifiers SHOULD be stable across different versions of an implementation
	TelemetrySdkNameKey = attribute.Key("telemetry.sdk.name")
	// The version string of the telemetry SDK.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "1.2.3",
	TelemetrySdkVersionKey = attribute.Key("telemetry.sdk.version")
)

// TelemetryDistroName returns an attribute KeyValue conforming to the "telemetry.distro.name"semantic conventions
// It represents the The name of the auto instrumentation agent or distribution, if used
func TelemetryDistroName(val string) attribute.KeyValue {
	return TelemetryDistroNameKey.String(val)
}

// TelemetryDistroVersion returns an attribute KeyValue conforming to the "telemetry.distro.version"semantic conventions
// It represents the The version string of the auto instrumentation agent or distribution, if used
func TelemetryDistroVersion(val string) attribute.KeyValue {
	return TelemetryDistroVersionKey.String(val)
}

// TelemetrySdkLanguage returns an attribute KeyValue conforming to the "telemetry.sdk.language"semantic conventions
// It represents the The language of the telemetry SDK
func TelemetrySdkLanguage(val string) attribute.KeyValue {
	return TelemetrySdkLanguageKey.String(val)
}

// TelemetrySdkName returns an attribute KeyValue conforming to the "telemetry.sdk.name"semantic conventions
// It represents the The name of the telemetry SDK as defined above
func TelemetrySdkName(val string) attribute.KeyValue {
	return TelemetrySdkNameKey.String(val)
}

// TelemetrySdkVersion returns an attribute KeyValue conforming to the "telemetry.sdk.version"semantic conventions
// It represents the The version string of the telemetry SDK
func TelemetrySdkVersion(val string) attribute.KeyValue {
	return TelemetrySdkVersionKey.String(val)
}

// Enum values for telemetry.sdk.language
var (

	// cpp
	// Stability: stable
	TelemetrySdkLanguageCpp = TelemetrySdkLanguageKey.String("cpp")
	// dotnet
	// Stability: stable
	TelemetrySdkLanguageDotnet = TelemetrySdkLanguageKey.String("dotnet")
	// erlang
	// Stability: stable
	TelemetrySdkLanguageErlang = TelemetrySdkLanguageKey.String("erlang")
	// go
	// Stability: stable
	TelemetrySdkLanguageGo = TelemetrySdkLanguageKey.String("go")
	// java
	// Stability: stable
	TelemetrySdkLanguageJava = TelemetrySdkLanguageKey.String("java")
	// nodejs
	// Stability: stable
	TelemetrySdkLanguageNodejs = TelemetrySdkLanguageKey.String("nodejs")
	// php
	// Stability: stable
	TelemetrySdkLanguagePhp = TelemetrySdkLanguageKey.String("php")
	// python
	// Stability: stable
	TelemetrySdkLanguagePython = TelemetrySdkLanguageKey.String("python")
	// ruby
	// Stability: stable
	TelemetrySdkLanguageRuby = TelemetrySdkLanguageKey.String("ruby")
	// rust
	// Stability: stable
	TelemetrySdkLanguageRust = TelemetrySdkLanguageKey.String("rust")
	// swift
	// Stability: stable
	TelemetrySdkLanguageSwift = TelemetrySdkLanguageKey.String("swift")
	// webjs
	// Stability: stable
	TelemetrySdkLanguageWebjs = TelemetrySdkLanguageKey.String("webjs")
)

// Namespace: test
const (

	// The fully qualified human readable name of the [test case].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "org.example.TestCase1.test1",
	// "example/tests/TestCase1.test1",
	// "ExampleTestCase1_test1",
	//
	// [test case]: https://en.wikipedia.org/wiki/Test_case
	TestCaseNameKey = attribute.Key("test.case.name")
	// The status of the actual test case result from test execution.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "pass",
	// "fail",
	TestCaseResultStatusKey = attribute.Key("test.case.result.status")
	// The human readable name of a [test suite].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "TestSuite1",
	//
	// [test suite]: https://en.wikipedia.org/wiki/Test_suite
	TestSuiteNameKey = attribute.Key("test.suite.name")
	// The status of the test suite run.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "success",
	// "failure",
	// "skipped",
	// "aborted",
	// "timed_out",
	// "in_progress",
	TestSuiteRunStatusKey = attribute.Key("test.suite.run.status")
)

// TestCaseName returns an attribute KeyValue conforming to the "test.case.name"semantic conventions
// It represents the The fully qualified human readable name of the [test case]
//
// [test case]: https://en.wikipedia.org/wiki/Test_case
func TestCaseName(val string) attribute.KeyValue {
	return TestCaseNameKey.String(val)
}

// TestCaseResultStatus returns an attribute KeyValue conforming to the "test.case.result.status"semantic conventions
// It represents the The status of the actual test case result from test execution
func TestCaseResultStatus(val string) attribute.KeyValue {
	return TestCaseResultStatusKey.String(val)
}

// TestSuiteName returns an attribute KeyValue conforming to the "test.suite.name"semantic conventions
// It represents the The human readable name of a [test suite]
//
// [test suite]: https://en.wikipedia.org/wiki/Test_suite
func TestSuiteName(val string) attribute.KeyValue {
	return TestSuiteNameKey.String(val)
}

// TestSuiteRunStatus returns an attribute KeyValue conforming to the "test.suite.run.status"semantic conventions
// It represents the The status of the test suite run
func TestSuiteRunStatus(val string) attribute.KeyValue {
	return TestSuiteRunStatusKey.String(val)
}

// Enum values for test.case.result.status
var (

	// pass
	// Stability: experimental
	TestCaseResultStatusPass = TestCaseResultStatusKey.String("pass")
	// fail
	// Stability: experimental
	TestCaseResultStatusFail = TestCaseResultStatusKey.String("fail")
)

// Enum values for test.suite.run.status
var (

	// success
	// Stability: experimental
	TestSuiteRunStatusSuccess = TestSuiteRunStatusKey.String("success")
	// failure
	// Stability: experimental
	TestSuiteRunStatusFailure = TestSuiteRunStatusKey.String("failure")
	// skipped
	// Stability: experimental
	TestSuiteRunStatusSkipped = TestSuiteRunStatusKey.String("skipped")
	// aborted
	// Stability: experimental
	TestSuiteRunStatusAborted = TestSuiteRunStatusKey.String("aborted")
	// timed_out
	// Stability: experimental
	TestSuiteRunStatusTimedOut = TestSuiteRunStatusKey.String("timed_out")
	// in_progress
	// Stability: experimental
	TestSuiteRunStatusInProgress = TestSuiteRunStatusKey.String("in_progress")
)

// Namespace: thread
const (

	// Current "managed" thread ID (as opposed to OS thread ID).
	//
	// Stability: Experimental
	// Type: int
	ThreadIdKey = attribute.Key("thread.id")
	// Current thread name.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples: "main"
	ThreadNameKey = attribute.Key("thread.name")
)

// ThreadId returns an attribute KeyValue conforming to the "thread.id"semantic conventions
// It represents the Current "managed" thread ID (as opposed to OS thread ID)
func ThreadId(val int) attribute.KeyValue {
	return ThreadIdKey.Int(val)
}

// ThreadName returns an attribute KeyValue conforming to the "thread.name"semantic conventions
// It represents the Current thread name
func ThreadName(val string) attribute.KeyValue {
	return ThreadNameKey.String(val)
}

// Namespace: tls
const (

	// String indicating the [cipher] used during the current connection.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	// "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	//
	// Note: The values allowed for `tls.cipher` MUST be one of the `Descriptions` of the [registered TLS Cipher Suits]
	//
	// [cipher]: https://datatracker.ietf.org/doc/html/rfc5246#appendix-A.5
	// [registered TLS Cipher Suits]: https://www.iana.org/assignments/tls-parameters/tls-parameters.xhtml#table-tls-parameters-4
	TlsCipherKey = attribute.Key("tls.cipher")
	// PEM-encoded stand-alone certificate offered by the client. This is usually mutually-exclusive of `client.certificate_chain` since this value also exists in that list.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "MII...",
	TlsClientCertificateKey = attribute.Key("tls.client.certificate")
	// Array of PEM-encoded certificates that make up the certificate chain offered by the client. This is usually mutually-exclusive of `client.certificate` since that value should be the first certificate in the chain.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "MII...",
	// "MI...",
	// ],
	TlsClientCertificateChainKey = attribute.Key("tls.client.certificate_chain")
	// Certificate fingerprint using the MD5 digest of DER-encoded version of certificate offered by the client. For consistency with other hash values, this value should be formatted as an uppercase hash.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "0F76C7F2C55BFD7D8E8B8F4BFBF0C9EC",
	TlsClientHashMd5Key = attribute.Key("tls.client.hash.md5")
	// Certificate fingerprint using the SHA1 digest of DER-encoded version of certificate offered by the client. For consistency with other hash values, this value should be formatted as an uppercase hash.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "9E393D93138888D288266C2D915214D1D1CCEB2A",
	TlsClientHashSha1Key = attribute.Key("tls.client.hash.sha1")
	// Certificate fingerprint using the SHA256 digest of DER-encoded version of certificate offered by the client. For consistency with other hash values, this value should be formatted as an uppercase hash.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "0687F666A054EF17A08E2F2162EAB4CBC0D265E1D7875BE74BF3C712CA92DAF0",
	TlsClientHashSha256Key = attribute.Key("tls.client.hash.sha256")
	// Distinguished name of [subject] of the issuer of the x.509 certificate presented by the client.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "CN=Example Root CA, OU=Infrastructure Team, DC=example, DC=com",
	//
	// [subject]: https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6
	TlsClientIssuerKey = attribute.Key("tls.client.issuer")
	// A hash that identifies clients based on how they perform an SSL/TLS handshake.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "d4e5b18d6b55c71272893221c96ba240",
	TlsClientJa3Key = attribute.Key("tls.client.ja3")
	// Date/Time indicating when client certificate is no longer considered valid.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2021-01-01T00:00:00.000Z",
	TlsClientNotAfterKey = attribute.Key("tls.client.not_after")
	// Date/Time indicating when client certificate is first considered valid.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "1970-01-01T00:00:00.000Z",
	TlsClientNotBeforeKey = attribute.Key("tls.client.not_before")
	// Deprecated, use `server.address` instead.
	// Stability: Experimental
	// Type: string
	// Deprecated: Replaced by `server.address`.
	//
	// Examples:
	// "opentelemetry.io",
	TlsClientServerNameKey = attribute.Key("tls.client.server_name")
	// Distinguished name of subject of the x.509 certificate presented by the client.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "CN=myclient, OU=Documentation Team, DC=example, DC=com",
	TlsClientSubjectKey = attribute.Key("tls.client.subject")
	// Array of ciphers offered by the client during the client hello.
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	// "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	// ],
	TlsClientSupportedCiphersKey = attribute.Key("tls.client.supported_ciphers")
	// String indicating the curve used for the given cipher, when applicable
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "secp256r1",
	TlsCurveKey = attribute.Key("tls.curve")
	// Boolean flag indicating if the TLS negotiation was successful and transitioned to an encrypted tunnel.
	// Stability: Experimental
	// Type: boolean
	//
	// Examples:
	// true,
	TlsEstablishedKey = attribute.Key("tls.established")
	// String indicating the protocol being tunneled. Per the values in the [IANA registry], this string should be lower case.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "http/1.1",
	//
	// [IANA registry]: https://www.iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml#alpn-protocol-ids
	TlsNextProtocolKey = attribute.Key("tls.next_protocol")
	// Normalized lowercase protocol name parsed from original string of the negotiated [SSL/TLS protocol version]
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	//
	// [SSL/TLS protocol version]: https://www.openssl.org/docs/man1.1.1/man3/SSL_get_version.html#RETURN-VALUES
	TlsProtocolNameKey = attribute.Key("tls.protocol.name")
	// Numeric part of the version parsed from the original string of the negotiated [SSL/TLS protocol version]
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "1.2",
	// "3",
	//
	// [SSL/TLS protocol version]: https://www.openssl.org/docs/man1.1.1/man3/SSL_get_version.html#RETURN-VALUES
	TlsProtocolVersionKey = attribute.Key("tls.protocol.version")
	// Boolean flag indicating if this TLS connection was resumed from an existing TLS negotiation.
	// Stability: Experimental
	// Type: boolean
	//
	// Examples:
	// true,
	TlsResumedKey = attribute.Key("tls.resumed")
	// PEM-encoded stand-alone certificate offered by the server. This is usually mutually-exclusive of `server.certificate_chain` since this value also exists in that list.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "MII...",
	TlsServerCertificateKey = attribute.Key("tls.server.certificate")
	// Array of PEM-encoded certificates that make up the certificate chain offered by the server. This is usually mutually-exclusive of `server.certificate` since that value should be the first certificate in the chain.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "MII...",
	// "MI...",
	// ],
	TlsServerCertificateChainKey = attribute.Key("tls.server.certificate_chain")
	// Certificate fingerprint using the MD5 digest of DER-encoded version of certificate offered by the server. For consistency with other hash values, this value should be formatted as an uppercase hash.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "0F76C7F2C55BFD7D8E8B8F4BFBF0C9EC",
	TlsServerHashMd5Key = attribute.Key("tls.server.hash.md5")
	// Certificate fingerprint using the SHA1 digest of DER-encoded version of certificate offered by the server. For consistency with other hash values, this value should be formatted as an uppercase hash.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "9E393D93138888D288266C2D915214D1D1CCEB2A",
	TlsServerHashSha1Key = attribute.Key("tls.server.hash.sha1")
	// Certificate fingerprint using the SHA256 digest of DER-encoded version of certificate offered by the server. For consistency with other hash values, this value should be formatted as an uppercase hash.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "0687F666A054EF17A08E2F2162EAB4CBC0D265E1D7875BE74BF3C712CA92DAF0",
	TlsServerHashSha256Key = attribute.Key("tls.server.hash.sha256")
	// Distinguished name of [subject] of the issuer of the x.509 certificate presented by the client.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "CN=Example Root CA, OU=Infrastructure Team, DC=example, DC=com",
	//
	// [subject]: https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6
	TlsServerIssuerKey = attribute.Key("tls.server.issuer")
	// A hash that identifies servers based on how they perform an SSL/TLS handshake.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "d4e5b18d6b55c71272893221c96ba240",
	TlsServerJa3sKey = attribute.Key("tls.server.ja3s")
	// Date/Time indicating when server certificate is no longer considered valid.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "2021-01-01T00:00:00.000Z",
	TlsServerNotAfterKey = attribute.Key("tls.server.not_after")
	// Date/Time indicating when server certificate is first considered valid.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "1970-01-01T00:00:00.000Z",
	TlsServerNotBeforeKey = attribute.Key("tls.server.not_before")
	// Distinguished name of subject of the x.509 certificate presented by the server.
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "CN=myserver, OU=Documentation Team, DC=example, DC=com",
	TlsServerSubjectKey = attribute.Key("tls.server.subject")
)

// TlsCipher returns an attribute KeyValue conforming to the "tls.cipher"semantic conventions
// It represents the String indicating the [cipher] used during the current connection
//
// [cipher]: https://datatracker.ietf.org/doc/html/rfc5246#appendix-A.5
func TlsCipher(val string) attribute.KeyValue {
	return TlsCipherKey.String(val)
}

// TlsClientCertificate returns an attribute KeyValue conforming to the "tls.client.certificate"semantic conventions
// It represents the PEM-encoded stand-alone certificate offered by the client. This is usually mutually-exclusive of `client.certificate_chain` since this value also exists in that list
func TlsClientCertificate(val string) attribute.KeyValue {
	return TlsClientCertificateKey.String(val)
}

// TlsClientCertificateChain returns an attribute KeyValue conforming to the "tls.client.certificate_chain"semantic conventions
// It represents the Array of PEM-encoded certificates that make up the certificate chain offered by the client. This is usually mutually-exclusive of `client.certificate` since that value should be the first certificate in the chain
func TlsClientCertificateChain(val []string) attribute.KeyValue {
	return TlsClientCertificateChainKey.StringSlice(val)
}

// TlsClientHashMd5 returns an attribute KeyValue conforming to the "tls.client.hash.md5"semantic conventions
// It represents the Certificate fingerprint using the MD5 digest of DER-encoded version of certificate offered by the client. For consistency with other hash values, this value should be formatted as an uppercase hash
func TlsClientHashMd5(val string) attribute.KeyValue {
	return TlsClientHashMd5Key.String(val)
}

// TlsClientHashSha1 returns an attribute KeyValue conforming to the "tls.client.hash.sha1"semantic conventions
// It represents the Certificate fingerprint using the SHA1 digest of DER-encoded version of certificate offered by the client. For consistency with other hash values, this value should be formatted as an uppercase hash
func TlsClientHashSha1(val string) attribute.KeyValue {
	return TlsClientHashSha1Key.String(val)
}

// TlsClientHashSha256 returns an attribute KeyValue conforming to the "tls.client.hash.sha256"semantic conventions
// It represents the Certificate fingerprint using the SHA256 digest of DER-encoded version of certificate offered by the client. For consistency with other hash values, this value should be formatted as an uppercase hash
func TlsClientHashSha256(val string) attribute.KeyValue {
	return TlsClientHashSha256Key.String(val)
}

// TlsClientIssuer returns an attribute KeyValue conforming to the "tls.client.issuer"semantic conventions
// It represents the Distinguished name of [subject] of the issuer of the x.509 certificate presented by the client
//
// [subject]: https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6
func TlsClientIssuer(val string) attribute.KeyValue {
	return TlsClientIssuerKey.String(val)
}

// TlsClientJa3 returns an attribute KeyValue conforming to the "tls.client.ja3"semantic conventions
// It represents the A hash that identifies clients based on how they perform an SSL/TLS handshake
func TlsClientJa3(val string) attribute.KeyValue {
	return TlsClientJa3Key.String(val)
}

// TlsClientNotAfter returns an attribute KeyValue conforming to the "tls.client.not_after"semantic conventions
// It represents the Date/Time indicating when client certificate is no longer considered valid
func TlsClientNotAfter(val string) attribute.KeyValue {
	return TlsClientNotAfterKey.String(val)
}

// TlsClientNotBefore returns an attribute KeyValue conforming to the "tls.client.not_before"semantic conventions
// It represents the Date/Time indicating when client certificate is first considered valid
func TlsClientNotBefore(val string) attribute.KeyValue {
	return TlsClientNotBeforeKey.String(val)
}

// TlsClientServerName returns an attribute KeyValue conforming to the "tls.client.server_name"semantic conventions
// It represents the Deprecated, use `server.address` instead
func TlsClientServerName(val string) attribute.KeyValue {
	return TlsClientServerNameKey.String(val)
}

// TlsClientSubject returns an attribute KeyValue conforming to the "tls.client.subject"semantic conventions
// It represents the Distinguished name of subject of the x.509 certificate presented by the client
func TlsClientSubject(val string) attribute.KeyValue {
	return TlsClientSubjectKey.String(val)
}

// TlsClientSupportedCiphers returns an attribute KeyValue conforming to the "tls.client.supported_ciphers"semantic conventions
// It represents the Array of ciphers offered by the client during the client hello
func TlsClientSupportedCiphers(val []string) attribute.KeyValue {
	return TlsClientSupportedCiphersKey.StringSlice(val)
}

// TlsCurve returns an attribute KeyValue conforming to the "tls.curve"semantic conventions
// It represents the String indicating the curve used for the given cipher, when applicable
func TlsCurve(val string) attribute.KeyValue {
	return TlsCurveKey.String(val)
}

// TlsEstablished returns an attribute KeyValue conforming to the "tls.established"semantic conventions
// It represents the Boolean flag indicating if the TLS negotiation was successful and transitioned to an encrypted tunnel
func TlsEstablished(val bool) attribute.KeyValue {
	return TlsEstablishedKey.Bool(val)
}

// TlsNextProtocol returns an attribute KeyValue conforming to the "tls.next_protocol"semantic conventions
// It represents the String indicating the protocol being tunneled. Per the values in the [IANA registry], this string should be lower case
//
// [IANA registry]: https://www.iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml#alpn-protocol-ids
func TlsNextProtocol(val string) attribute.KeyValue {
	return TlsNextProtocolKey.String(val)
}

// TlsProtocolName returns an attribute KeyValue conforming to the "tls.protocol.name"semantic conventions
// It represents the Normalized lowercase protocol name parsed from original string of the negotiated [SSL/TLS protocol version]
//
// [SSL/TLS protocol version]: https://www.openssl.org/docs/man1.1.1/man3/SSL_get_version.html#RETURN-VALUES
func TlsProtocolName(val string) attribute.KeyValue {
	return TlsProtocolNameKey.String(val)
}

// TlsProtocolVersion returns an attribute KeyValue conforming to the "tls.protocol.version"semantic conventions
// It represents the Numeric part of the version parsed from the original string of the negotiated [SSL/TLS protocol version]
//
// [SSL/TLS protocol version]: https://www.openssl.org/docs/man1.1.1/man3/SSL_get_version.html#RETURN-VALUES
func TlsProtocolVersion(val string) attribute.KeyValue {
	return TlsProtocolVersionKey.String(val)
}

// TlsResumed returns an attribute KeyValue conforming to the "tls.resumed"semantic conventions
// It represents the Boolean flag indicating if this TLS connection was resumed from an existing TLS negotiation
func TlsResumed(val bool) attribute.KeyValue {
	return TlsResumedKey.Bool(val)
}

// TlsServerCertificate returns an attribute KeyValue conforming to the "tls.server.certificate"semantic conventions
// It represents the PEM-encoded stand-alone certificate offered by the server. This is usually mutually-exclusive of `server.certificate_chain` since this value also exists in that list
func TlsServerCertificate(val string) attribute.KeyValue {
	return TlsServerCertificateKey.String(val)
}

// TlsServerCertificateChain returns an attribute KeyValue conforming to the "tls.server.certificate_chain"semantic conventions
// It represents the Array of PEM-encoded certificates that make up the certificate chain offered by the server. This is usually mutually-exclusive of `server.certificate` since that value should be the first certificate in the chain
func TlsServerCertificateChain(val []string) attribute.KeyValue {
	return TlsServerCertificateChainKey.StringSlice(val)
}

// TlsServerHashMd5 returns an attribute KeyValue conforming to the "tls.server.hash.md5"semantic conventions
// It represents the Certificate fingerprint using the MD5 digest of DER-encoded version of certificate offered by the server. For consistency with other hash values, this value should be formatted as an uppercase hash
func TlsServerHashMd5(val string) attribute.KeyValue {
	return TlsServerHashMd5Key.String(val)
}

// TlsServerHashSha1 returns an attribute KeyValue conforming to the "tls.server.hash.sha1"semantic conventions
// It represents the Certificate fingerprint using the SHA1 digest of DER-encoded version of certificate offered by the server. For consistency with other hash values, this value should be formatted as an uppercase hash
func TlsServerHashSha1(val string) attribute.KeyValue {
	return TlsServerHashSha1Key.String(val)
}

// TlsServerHashSha256 returns an attribute KeyValue conforming to the "tls.server.hash.sha256"semantic conventions
// It represents the Certificate fingerprint using the SHA256 digest of DER-encoded version of certificate offered by the server. For consistency with other hash values, this value should be formatted as an uppercase hash
func TlsServerHashSha256(val string) attribute.KeyValue {
	return TlsServerHashSha256Key.String(val)
}

// TlsServerIssuer returns an attribute KeyValue conforming to the "tls.server.issuer"semantic conventions
// It represents the Distinguished name of [subject] of the issuer of the x.509 certificate presented by the client
//
// [subject]: https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6
func TlsServerIssuer(val string) attribute.KeyValue {
	return TlsServerIssuerKey.String(val)
}

// TlsServerJa3s returns an attribute KeyValue conforming to the "tls.server.ja3s"semantic conventions
// It represents the A hash that identifies servers based on how they perform an SSL/TLS handshake
func TlsServerJa3s(val string) attribute.KeyValue {
	return TlsServerJa3sKey.String(val)
}

// TlsServerNotAfter returns an attribute KeyValue conforming to the "tls.server.not_after"semantic conventions
// It represents the Date/Time indicating when server certificate is no longer considered valid
func TlsServerNotAfter(val string) attribute.KeyValue {
	return TlsServerNotAfterKey.String(val)
}

// TlsServerNotBefore returns an attribute KeyValue conforming to the "tls.server.not_before"semantic conventions
// It represents the Date/Time indicating when server certificate is first considered valid
func TlsServerNotBefore(val string) attribute.KeyValue {
	return TlsServerNotBeforeKey.String(val)
}

// TlsServerSubject returns an attribute KeyValue conforming to the "tls.server.subject"semantic conventions
// It represents the Distinguished name of subject of the x.509 certificate presented by the server
func TlsServerSubject(val string) attribute.KeyValue {
	return TlsServerSubjectKey.String(val)
}

// Enum values for tls.protocol.name
var (

	// ssl
	// Stability: experimental
	TlsProtocolNameSsl = TlsProtocolNameKey.String("ssl")
	// tls
	// Stability: experimental
	TlsProtocolNameTls = TlsProtocolNameKey.String("tls")
)

// Namespace: url
const (

	// Domain extracted from the `url.full`, such as "opentelemetry.io".
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "www.foo.bar",
	// "opentelemetry.io",
	// "3.12.167.2",
	// "[1080:0:0:0:8:800:200C:417A]",
	//
	// Note: In some cases a URL may refer to an IP and/or port directly, without a domain name. In this case, the IP address would go to the domain field. If the URL contains a [literal IPv6 address] enclosed by `[` and `]`, the `[` and `]` characters should also be captured in the domain field
	//
	// [literal IPv6 address]: https://www.rfc-editor.org/rfc/rfc2732#section-2
	UrlDomainKey = attribute.Key("url.domain")
	// The file extension extracted from the `url.full`, excluding the leading dot.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "png",
	// "gz",
	//
	// Note: The file extension is only set if it exists, as not every url has a file extension. When the file name has multiple extensions `example.tar.gz`, only the last one should be captured `gz`, not `tar.gz`
	UrlExtensionKey = attribute.Key("url.extension")
	// The [URI fragment] component
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "SemConv",
	//
	// [URI fragment]: https://www.rfc-editor.org/rfc/rfc3986#section-3.5
	UrlFragmentKey = attribute.Key("url.fragment")
	// Absolute URL describing a network resource according to [RFC3986]
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "https://www.foo.bar/search?q=OpenTelemetry#SemConv",
	// "//localhost",
	//
	// Note: For network calls, URL usually has `scheme://host[:port][path][?query][#fragment]` format, where the fragment is not transmitted over HTTP, but if it is known, it SHOULD be included nevertheless.
	// `url.full` MUST NOT contain credentials passed via URL in form of `https://username:password@www.example.com/`. In such case username and password SHOULD be redacted and attribute's value SHOULD be `https://REDACTED:REDACTED@www.example.com/`.
	// `url.full` SHOULD capture the absolute URL when it is available (or can be reconstructed). Sensitive content provided in `url.full` SHOULD be scrubbed when instrumentations can identify it
	//
	// [RFC3986]: https://www.rfc-editor.org/rfc/rfc3986
	UrlFullKey = attribute.Key("url.full")
	// Unmodified original URL as seen in the event source.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "https://www.foo.bar/search?q=OpenTelemetry#SemConv",
	// "search?q=OpenTelemetry",
	//
	// Note: In network monitoring, the observed URL may be a full URL, whereas in access logs, the URL is often just represented as a path. This field is meant to represent the URL as it was observed, complete or not.
	// `url.original` might contain credentials passed via URL in form of `https://username:password@www.example.com/`. In such case password and username SHOULD NOT be redacted and attribute's value SHOULD remain the same
	UrlOriginalKey = attribute.Key("url.original")
	// The [URI path] component
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "/search",
	//
	// Note: Sensitive content provided in `url.path` SHOULD be scrubbed when instrumentations can identify it
	//
	// [URI path]: https://www.rfc-editor.org/rfc/rfc3986#section-3.3
	UrlPathKey = attribute.Key("url.path")
	// Port extracted from the `url.full`
	//
	// Stability: Experimental
	// Type: int
	//
	// Examples:
	// 443,
	UrlPortKey = attribute.Key("url.port")
	// The [URI query] component
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "q=OpenTelemetry",
	//
	// Note: Sensitive content provided in `url.query` SHOULD be scrubbed when instrumentations can identify it
	//
	// [URI query]: https://www.rfc-editor.org/rfc/rfc3986#section-3.4
	UrlQueryKey = attribute.Key("url.query")
	// The highest registered url domain, stripped of the subdomain.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "example.com",
	// "foo.co.uk",
	//
	// Note: This value can be determined precisely with the [public suffix list]. For example, the registered domain for `foo.example.com` is `example.com`. Trying to approximate this by simply taking the last two labels will not work well for TLDs such as `co.uk`
	//
	// [public suffix list]: http://publicsuffix.org
	UrlRegisteredDomainKey = attribute.Key("url.registered_domain")
	// The [URI scheme] component identifying the used protocol.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "https",
	// "ftp",
	// "telnet",
	//
	// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
	UrlSchemeKey = attribute.Key("url.scheme")
	// The subdomain portion of a fully qualified domain name includes all of the names except the host name under the registered_domain. In a partially qualified domain, or if the qualification level of the full name cannot be determined, subdomain contains all of the names below the registered domain.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "east",
	// "sub2.sub1",
	//
	// Note: The subdomain portion of `www.east.mydomain.co.uk` is `east`. If the domain has multiple levels of subdomain, such as `sub2.sub1.example.com`, the subdomain field should contain `sub2.sub1`, with no trailing period
	UrlSubdomainKey = attribute.Key("url.subdomain")
	// The low-cardinality template of an [absolute path reference].
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "/users/{id}",
	// "/users/:id",
	// "/users?id={id}",
	//
	// [absolute path reference]: https://www.rfc-editor.org/rfc/rfc3986#section-4.2
	UrlTemplateKey = attribute.Key("url.template")
	// The effective top level domain (eTLD), also known as the domain suffix, is the last part of the domain name. For example, the top level domain for example.com is `com`.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "com",
	// "co.uk",
	//
	// Note: This value can be determined precisely with the [public suffix list]
	//
	// [public suffix list]: http://publicsuffix.org
	UrlTopLevelDomainKey = attribute.Key("url.top_level_domain")
)

// UrlDomain returns an attribute KeyValue conforming to the "url.domain"semantic conventions
// It represents the Domain extracted from the `url.full`, such as "opentelemetry.io"
func UrlDomain(val string) attribute.KeyValue {
	return UrlDomainKey.String(val)
}

// UrlExtension returns an attribute KeyValue conforming to the "url.extension"semantic conventions
// It represents the The file extension extracted from the `url.full`, excluding the leading dot
func UrlExtension(val string) attribute.KeyValue {
	return UrlExtensionKey.String(val)
}

// UrlFragment returns an attribute KeyValue conforming to the "url.fragment"semantic conventions
// It represents the The [URI fragment] component
//
// [URI fragment]: https://www.rfc-editor.org/rfc/rfc3986#section-3.5
func UrlFragment(val string) attribute.KeyValue {
	return UrlFragmentKey.String(val)
}

// UrlFull returns an attribute KeyValue conforming to the "url.full"semantic conventions
// It represents the Absolute URL describing a network resource according to [RFC3986]
//
// [RFC3986]: https://www.rfc-editor.org/rfc/rfc3986
func UrlFull(val string) attribute.KeyValue {
	return UrlFullKey.String(val)
}

// UrlOriginal returns an attribute KeyValue conforming to the "url.original"semantic conventions
// It represents the Unmodified original URL as seen in the event source
func UrlOriginal(val string) attribute.KeyValue {
	return UrlOriginalKey.String(val)
}

// UrlPath returns an attribute KeyValue conforming to the "url.path"semantic conventions
// It represents the The [URI path] component
//
// [URI path]: https://www.rfc-editor.org/rfc/rfc3986#section-3.3
func UrlPath(val string) attribute.KeyValue {
	return UrlPathKey.String(val)
}

// UrlPort returns an attribute KeyValue conforming to the "url.port"semantic conventions
// It represents the Port extracted from the `url.full`
func UrlPort(val int) attribute.KeyValue {
	return UrlPortKey.Int(val)
}

// UrlQuery returns an attribute KeyValue conforming to the "url.query"semantic conventions
// It represents the The [URI query] component
//
// [URI query]: https://www.rfc-editor.org/rfc/rfc3986#section-3.4
func UrlQuery(val string) attribute.KeyValue {
	return UrlQueryKey.String(val)
}

// UrlRegisteredDomain returns an attribute KeyValue conforming to the "url.registered_domain"semantic conventions
// It represents the The highest registered url domain, stripped of the subdomain
func UrlRegisteredDomain(val string) attribute.KeyValue {
	return UrlRegisteredDomainKey.String(val)
}

// UrlScheme returns an attribute KeyValue conforming to the "url.scheme"semantic conventions
// It represents the The [URI scheme] component identifying the used protocol
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func UrlScheme(val string) attribute.KeyValue {
	return UrlSchemeKey.String(val)
}

// UrlSubdomain returns an attribute KeyValue conforming to the "url.subdomain"semantic conventions
// It represents the The subdomain portion of a fully qualified domain name includes all of the names except the host name under the registered_domain. In a partially qualified domain, or if the qualification level of the full name cannot be determined, subdomain contains all of the names below the registered domain
func UrlSubdomain(val string) attribute.KeyValue {
	return UrlSubdomainKey.String(val)
}

// UrlTemplate returns an attribute KeyValue conforming to the "url.template"semantic conventions
// It represents the The low-cardinality template of an [absolute path reference]
//
// [absolute path reference]: https://www.rfc-editor.org/rfc/rfc3986#section-4.2
func UrlTemplate(val string) attribute.KeyValue {
	return UrlTemplateKey.String(val)
}

// UrlTopLevelDomain returns an attribute KeyValue conforming to the "url.top_level_domain"semantic conventions
// It represents the The effective top level domain (eTLD), also known as the domain suffix, is the last part of the domain name. For example, the top level domain for example.com is `com`
func UrlTopLevelDomain(val string) attribute.KeyValue {
	return UrlTopLevelDomainKey.String(val)
}

// Namespace: user
const (

	// User email address.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "a.einstein@example.com",
	UserEmailKey = attribute.Key("user.email")
	// User's full name
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Albert Einstein",
	UserFullNameKey = attribute.Key("user.full_name")
	// Unique user hash to correlate information for a user in anonymized form.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "364fc68eaf4c8acec74a4e52d7d1feaa",
	//
	// Note: Useful if `user.id` or `user.name` contain confidential information and cannot be used
	UserHashKey = attribute.Key("user.hash")
	// Unique identifier of the user.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "S-1-5-21-202424912787-2692429404-2351956786-1000",
	UserIdKey = attribute.Key("user.id")
	// Short name or login/username of the user.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "a.einstein",
	UserNameKey = attribute.Key("user.name")
	// Array of user roles at the time of the event.
	//
	// Stability: Experimental
	// Type: string[]
	//
	// Examples:
	// [
	// "admin",
	// "reporting_user",
	// ],
	UserRolesKey = attribute.Key("user.roles")
)

// UserEmail returns an attribute KeyValue conforming to the "user.email"semantic conventions
// It represents the User email address
func UserEmail(val string) attribute.KeyValue {
	return UserEmailKey.String(val)
}

// UserFullName returns an attribute KeyValue conforming to the "user.full_name"semantic conventions
// It represents the User's full name
func UserFullName(val string) attribute.KeyValue {
	return UserFullNameKey.String(val)
}

// UserHash returns an attribute KeyValue conforming to the "user.hash"semantic conventions
// It represents the Unique user hash to correlate information for a user in anonymized form
func UserHash(val string) attribute.KeyValue {
	return UserHashKey.String(val)
}

// UserId returns an attribute KeyValue conforming to the "user.id"semantic conventions
// It represents the Unique identifier of the user
func UserId(val string) attribute.KeyValue {
	return UserIdKey.String(val)
}

// UserName returns an attribute KeyValue conforming to the "user.name"semantic conventions
// It represents the Short name or login/username of the user
func UserName(val string) attribute.KeyValue {
	return UserNameKey.String(val)
}

// UserRoles returns an attribute KeyValue conforming to the "user.roles"semantic conventions
// It represents the Array of user roles at the time of the event
func UserRoles(val []string) attribute.KeyValue {
	return UserRolesKey.StringSlice(val)
}

// Namespace: user_agent
const (

	// Name of the user-agent extracted from original. Usually refers to the browser's name.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Safari",
	// "YourApp",
	//
	// Note: [Example] of extracting browser's name from original string. In the case of using a user-agent for non-browser products, such as microservices with multiple names/versions inside the `user_agent.original`, the most significant name SHOULD be selected. In such a scenario it should align with `user_agent.version`
	//
	// [Example]: https://www.whatsmyua.info
	UserAgentNameKey = attribute.Key("user_agent.name")
	// Value of the [HTTP User-Agent] header sent by the client.
	//
	// Stability: Stable
	// Type: string
	//
	// Examples:
	// "CERN-LineMode/2.15 libwww/2.17b3",
	// "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Mobile/15E148 Safari/604.1",
	// "YourApp/1.0.0 grpc-java-okhttp/1.27.2",
	//
	// [HTTP User-Agent]: https://www.rfc-editor.org/rfc/rfc9110.html#field.user-agent
	UserAgentOriginalKey = attribute.Key("user_agent.original")
	// Version of the user-agent extracted from original. Usually refers to the browser's version
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "14.1.2",
	// "1.0.0",
	//
	// Note: [Example] of extracting browser's version from original string. In the case of using a user-agent for non-browser products, such as microservices with multiple names/versions inside the `user_agent.original`, the most significant version SHOULD be selected. In such a scenario it should align with `user_agent.name`
	//
	// [Example]: https://www.whatsmyua.info
	UserAgentVersionKey = attribute.Key("user_agent.version")
)

// UserAgentName returns an attribute KeyValue conforming to the "user_agent.name"semantic conventions
// It represents the Name of the user-agent extracted from original. Usually refers to the browser's name
func UserAgentName(val string) attribute.KeyValue {
	return UserAgentNameKey.String(val)
}

// UserAgentOriginal returns an attribute KeyValue conforming to the "user_agent.original"semantic conventions
// It represents the Value of the [HTTP User-Agent] header sent by the client
//
// [HTTP User-Agent]: https://www.rfc-editor.org/rfc/rfc9110.html#field.user-agent
func UserAgentOriginal(val string) attribute.KeyValue {
	return UserAgentOriginalKey.String(val)
}

// UserAgentVersion returns an attribute KeyValue conforming to the "user_agent.version"semantic conventions
// It represents the Version of the user-agent extracted from original. Usually refers to the browser's version
func UserAgentVersion(val string) attribute.KeyValue {
	return UserAgentVersionKey.String(val)
}

// Namespace: v8js
const (

	// The type of garbage collection.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	V8jsGcTypeKey = attribute.Key("v8js.gc.type")
	// The name of the space type of heap memory.
	// Stability: Experimental
	// Type: Enum
	//
	// Examples: undefined
	// Note: Value can be retrieved from value `space_name` of [`v8.getHeapSpaceStatistics()`]
	//
	// [`v8.getHeapSpaceStatistics()`]: https://nodejs.org/api/v8.html#v8getheapspacestatistics
	V8jsHeapSpaceNameKey = attribute.Key("v8js.heap.space.name")
)

// V8jsGcType returns an attribute KeyValue conforming to the "v8js.gc.type"semantic conventions
// It represents the The type of garbage collection
func V8jsGcType(val string) attribute.KeyValue {
	return V8jsGcTypeKey.String(val)
}

// V8jsHeapSpaceName returns an attribute KeyValue conforming to the "v8js.heap.space.name"semantic conventions
// It represents the The name of the space type of heap memory
func V8jsHeapSpaceName(val string) attribute.KeyValue {
	return V8jsHeapSpaceNameKey.String(val)
}

// Enum values for v8js.gc.type
var (

	// Major (Mark Sweep Compact).
	// Stability: experimental
	V8jsGcTypeMajor = V8jsGcTypeKey.String("major")
	// Minor (Scavenge).
	// Stability: experimental
	V8jsGcTypeMinor = V8jsGcTypeKey.String("minor")
	// Incremental (Incremental Marking).
	// Stability: experimental
	V8jsGcTypeIncremental = V8jsGcTypeKey.String("incremental")
	// Weak Callbacks (Process Weak Callbacks).
	// Stability: experimental
	V8jsGcTypeWeakcb = V8jsGcTypeKey.String("weakcb")
)

// Enum values for v8js.heap.space.name
var (

	// New memory space.
	// Stability: experimental
	V8jsHeapSpaceNameNewSpace = V8jsHeapSpaceNameKey.String("new_space")
	// Old memory space.
	// Stability: experimental
	V8jsHeapSpaceNameOldSpace = V8jsHeapSpaceNameKey.String("old_space")
	// Code memory space.
	// Stability: experimental
	V8jsHeapSpaceNameCodeSpace = V8jsHeapSpaceNameKey.String("code_space")
	// Map memory space.
	// Stability: experimental
	V8jsHeapSpaceNameMapSpace = V8jsHeapSpaceNameKey.String("map_space")
	// Large object memory space.
	// Stability: experimental
	V8jsHeapSpaceNameLargeObjectSpace = V8jsHeapSpaceNameKey.String("large_object_space")
)

// Namespace: vcs
const (

	// The ID of the change (pull request/merge request) if applicable. This is usually a unique (within repository) identifier generated by the VCS system.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "123",
	VcsRepositoryChangeIdKey = attribute.Key("vcs.repository.change.id")
	// The human readable title of the change (pull request/merge request). This title is often a brief summary of the change and may get merged in to a ref as the commit summary.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "Fixes broken thing",
	// "feat: add my new feature",
	// "[chore] update dependency",
	VcsRepositoryChangeTitleKey = attribute.Key("vcs.repository.change.title")
	// The name of the [reference] such as **branch** or **tag** in the repository.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "my-feature-branch",
	// "tag-1-test",
	//
	// [reference]: https://git-scm.com/docs/gitglossary#def_ref
	VcsRepositoryRefNameKey = attribute.Key("vcs.repository.ref.name")
	// The revision, literally [revised version], The revision most often refers to a commit object in Git, or a revision number in SVN.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "9d59409acf479dfa0df1aa568182e43e43df8bbe28d60fcf2bc52e30068802cc",
	// "main",
	// "123",
	// "HEAD",
	//
	// Note: The revision can be a full [hash value (see glossary)],
	// of the recorded change to a ref within a repository pointing to a
	// commit [commit] object. It does
	// not necessarily have to be a hash; it can simply define a
	// [revision number]
	// which is an integer that is monotonically increasing. In cases where
	// it is identical to the `ref.name`, it SHOULD still be included. It is
	// up to the implementer to decide which value to set as the revision
	// based on the VCS system and situational context
	//
	// [revised version]: https://www.merriam-webster.com/dictionary/revision
	// [hash value (see glossary)]: https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-5.pdf
	// [commit]: https://git-scm.com/docs/git-commit
	// [revision number]: https://svnbook.red-bean.com/en/1.7/svn.tour.revs.specifiers.html
	VcsRepositoryRefRevisionKey = attribute.Key("vcs.repository.ref.revision")
	// The type of the [reference] in the repository.
	//
	// Stability: Experimental
	// Type: Enum
	//
	// Examples:
	// "branch",
	// "tag",
	//
	// [reference]: https://git-scm.com/docs/gitglossary#def_ref
	VcsRepositoryRefTypeKey = attribute.Key("vcs.repository.ref.type")
	// The [URL] of the repository providing the complete address in order to locate and identify the repository.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "https://github.com/opentelemetry/open-telemetry-collector-contrib",
	// "https://gitlab.com/my-org/my-project/my-projects-project/repo",
	//
	// [URL]: https://en.wikipedia.org/wiki/URL
	VcsRepositoryUrlFullKey = attribute.Key("vcs.repository.url.full")
)

// VcsRepositoryChangeId returns an attribute KeyValue conforming to the "vcs.repository.change.id"semantic conventions
// It represents the The ID of the change (pull request/merge request) if applicable. This is usually a unique (within repository) identifier generated by the VCS system
func VcsRepositoryChangeId(val string) attribute.KeyValue {
	return VcsRepositoryChangeIdKey.String(val)
}

// VcsRepositoryChangeTitle returns an attribute KeyValue conforming to the "vcs.repository.change.title"semantic conventions
// It represents the The human readable title of the change (pull request/merge request). This title is often a brief summary of the change and may get merged in to a ref as the commit summary
func VcsRepositoryChangeTitle(val string) attribute.KeyValue {
	return VcsRepositoryChangeTitleKey.String(val)
}

// VcsRepositoryRefName returns an attribute KeyValue conforming to the "vcs.repository.ref.name"semantic conventions
// It represents the The name of the [reference] such as **branch** or **tag** in the repository
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
func VcsRepositoryRefName(val string) attribute.KeyValue {
	return VcsRepositoryRefNameKey.String(val)
}

// VcsRepositoryRefRevision returns an attribute KeyValue conforming to the "vcs.repository.ref.revision"semantic conventions
// It represents the The revision, literally [revised version], The revision most often refers to a commit object in Git, or a revision number in SVN
//
// [revised version]: https://www.merriam-webster.com/dictionary/revision
func VcsRepositoryRefRevision(val string) attribute.KeyValue {
	return VcsRepositoryRefRevisionKey.String(val)
}

// VcsRepositoryRefType returns an attribute KeyValue conforming to the "vcs.repository.ref.type"semantic conventions
// It represents the The type of the [reference] in the repository
//
// [reference]: https://git-scm.com/docs/gitglossary#def_ref
func VcsRepositoryRefType(val string) attribute.KeyValue {
	return VcsRepositoryRefTypeKey.String(val)
}

// VcsRepositoryUrlFull returns an attribute KeyValue conforming to the "vcs.repository.url.full"semantic conventions
// It represents the The [URL] of the repository providing the complete address in order to locate and identify the repository
//
// [URL]: https://en.wikipedia.org/wiki/URL
func VcsRepositoryUrlFull(val string) attribute.KeyValue {
	return VcsRepositoryUrlFullKey.String(val)
}

// Enum values for vcs.repository.ref.type
var (

	// [branch]
	// Stability: experimental
	//
	// [branch]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddefbranchabranch
	VcsRepositoryRefTypeBranch = VcsRepositoryRefTypeKey.String("branch")
	// [tag]
	// Stability: experimental
	//
	// [tag]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddeftagatag
	VcsRepositoryRefTypeTag = VcsRepositoryRefTypeKey.String("tag")
)

// Namespace: webengine
const (

	// Additional description of the web engine (e.g. detailed version and edition information).
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "WildFly Full 21.0.0.Final (WildFly Core 13.0.1.Final) - 2.2.2.Final",
	WebengineDescriptionKey = attribute.Key("webengine.description")
	// The name of the web engine.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "WildFly",
	WebengineNameKey = attribute.Key("webengine.name")
	// The version of the web engine.
	//
	// Stability: Experimental
	// Type: string
	//
	// Examples:
	// "21.0.0",
	WebengineVersionKey = attribute.Key("webengine.version")
)

// WebengineDescription returns an attribute KeyValue conforming to the "webengine.description"semantic conventions
// It represents the Additional description of the web engine (e.g. detailed version and edition information)
func WebengineDescription(val string) attribute.KeyValue {
	return WebengineDescriptionKey.String(val)
}

// WebengineName returns an attribute KeyValue conforming to the "webengine.name"semantic conventions
// It represents the The name of the web engine
func WebengineName(val string) attribute.KeyValue {
	return WebengineNameKey.String(val)
}

// WebengineVersion returns an attribute KeyValue conforming to the "webengine.version"semantic conventions
// It represents the The version of the web engine
func WebengineVersion(val string) attribute.KeyValue {
	return WebengineVersionKey.String(val)
}
