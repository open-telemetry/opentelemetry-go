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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
)

// Semantic conventions for service resource attribute keys.
const (
	// Name of the service.
	ServiceNameKey = semconv.ServiceNameKey

	// A namespace for `service.name`. This needs to have meaning that helps
	// to distinguish a group of services. For example, the team name that
	// owns a group of services. `service.name` is expected to be unique
	// within the same namespace.
	//
	// [Modified in 1.4.0]. TODO: revert this change, it is done for illustration only,
	// in reality the convention for Service Namespace has not changed in 1.4.0.
	ServiceNamespaceKey = attribute.Key("service.namespace.name")

	// A unique identifier of the service instance. In conjunction with the
	// `service.name` and `service.namespace` this must be unique.
	ServiceInstanceIDKey = semconv.ServiceInstanceIDKey

	// The version of the service API.
	ServiceVersionKey = semconv.ServiceVersionKey
)

// Semantic conventions for telemetry SDK resource attribute keys.
const (
	// The name of the telemetry SDK.
	//
	// The default OpenTelemetry SDK provided by the OpenTelemetry project
	// MUST set telemetry.sdk.name to the value `opentelemetry`.
	//
	// If another SDK is used, this attribute MUST be set to the import path
	// of that SDK's package.
	//
	// The value `opentelemetry` is reserved and MUST NOT be used by
	// non-OpenTelemetry SDKs.
	TelemetrySDKNameKey = semconv.TelemetrySDKNameKey

	// The language of the telemetry SDK.
	TelemetrySDKLanguageKey = semconv.TelemetrySDKLanguageKey

	// The version string of the telemetry SDK.
	TelemetrySDKVersionKey = semconv.TelemetrySDKVersionKey
)

// Semantic conventions for telemetry SDK resource attributes.
var (
	TelemetrySDKLanguageGo = semconv.TelemetrySDKLanguageGo
)

// Semantic conventions for container resource attribute keys.
const (
	// A uniquely identifying name for the Container.
	ContainerNameKey = semconv.ContainerNameKey

	// Container ID, usually a UUID, as for example used to
	// identify Docker containers. The UUID might be abbreviated.
	ContainerIDKey = semconv.ContainerIDKey

	// Name of the image the container was built on.
	ContainerImageNameKey = semconv.ContainerImageNameKey

	// Container image tag.
	ContainerImageTagKey = semconv.ContainerImageTagKey
)

// Semantic conventions for Function-as-a-Service resource attribute keys.
const (
	// A uniquely identifying name for the FaaS.
	FaaSNameKey = semconv.FaaSNameKey

	// The unique name of the function being executed.
	FaaSIDKey = semconv.FaaSIDKey

	// The version of the function being executed.
	FaaSVersionKey = semconv.FaaSVersionKey

	// The execution environment identifier.
	FaaSInstanceKey = semconv.FaaSInstanceKey
)

// Semantic conventions for operating system process resource attribute keys.
const (
	// Process identifier (PID).
	ProcessPIDKey = semconv.ProcessPIDKey
	// The name of the process executable. On Linux based systems, can be
	// set to the `Name` in `proc/[pid]/status`. On Windows, can be set to
	// the base name of `GetProcessImageFileNameW`.
	ProcessExecutableNameKey = semconv.ProcessExecutableNameKey
	// The full path to the process executable. On Linux based systems, can
	// be set to the target of `proc/[pid]/exe`. On Windows, can be set to
	// the result of `GetProcessImageFileNameW`.
	ProcessExecutablePathKey = semconv.ProcessExecutablePathKey
	// The command used to launch the process (i.e. the command name). On
	// Linux based systems, can be set to the zeroth string in
	// `proc/[pid]/cmdline`. On Windows, can be set to the first parameter
	// extracted from `GetCommandLineW`.
	ProcessCommandKey = semconv.ProcessCommandKey
	// The full command used to launch the process. The value can be either
	// a list of strings representing the ordered list of arguments, or a
	// single string representing the full command. On Linux based systems,
	// can be set to the list of null-delimited strings extracted from
	// `proc/[pid]/cmdline`. On Windows, can be set to the result of
	// `GetCommandLineW`.
	ProcessCommandLineKey = semconv.ProcessCommandLineKey
	// All the command arguments (including the command/executable itself)
	// as received by the process. On Linux-based systems (and some other
	// Unixoid systems supporting procfs), can be set according to the list
	// of null-delimited strings extracted from `proc/[pid]/cmdline`. For
	// libc-based executables, this would be the full argv vector passed to
	// `main`.
	ProcessCommandArgsKey = semconv.ProcessCommandArgsKey
	// The username of the user that owns the process.
	ProcessOwnerKey = semconv.ProcessOwnerKey
	// The name of the runtime of this process. For compiled native
	// binaries, this SHOULD be the name of the compiler.
	ProcessRuntimeNameKey = semconv.ProcessRuntimeNameKey
	// The version of the runtime of this process, as returned by the
	// runtime without modification.
	ProcessRuntimeVersionKey = semconv.ProcessRuntimeVersionKey
	// An additional description about the runtime of the process, for
	// example a specific vendor customization of the runtime environment.
	ProcessRuntimeDescriptionKey = semconv.ProcessRuntimeDescriptionKey
)

// Semantic conventions for Kubernetes resource attribute keys.
const (
	// A uniquely identifying name for the Kubernetes cluster. Kubernetes
	// does not have cluster names as an internal concept so this may be
	// set to any meaningful value within the environment. For example,
	// GKE clusters have a name which can be used for this attribute.
	K8SClusterNameKey = semconv.K8SClusterNameKey

	// The name of the Node.
	K8SNodeNameKey = semconv.K8SNodeNameKey

	// The UID of the Node.
	K8SNodeUIDKey = semconv.K8SNodeUIDKey

	// The name of the namespace that the pod is running in.
	K8SNamespaceNameKey = semconv.K8SNamespaceNameKey

	// The uid of the Pod.
	K8SPodUIDKey = semconv.K8SPodUIDKey

	// The name of the pod.
	K8SPodNameKey = semconv.K8SPodNameKey

	// The name of the Container in a Pod template.
	K8SContainerNameKey = semconv.K8SContainerNameKey

	// The uid of the ReplicaSet.
	K8SReplicaSetUIDKey = semconv.K8SReplicaSetUIDKey

	// The name of the ReplicaSet.
	K8SReplicaSetNameKey = semconv.K8SReplicaSetNameKey

	// The uid of the Deployment.
	K8SDeploymentUIDKey = semconv.K8SDeploymentUIDKey

	// The name of the deployment.
	K8SDeploymentNameKey = semconv.K8SDeploymentNameKey

	// The uid of the StatefulSet.
	K8SStatefulSetUIDKey = semconv.K8SStatefulSetUIDKey

	// The name of the StatefulSet.
	K8SStatefulSetNameKey = semconv.K8SStatefulSetNameKey

	// The uid of the DaemonSet.
	K8SDaemonSetUIDKey = semconv.K8SDaemonSetUIDKey

	// The name of the DaemonSet.
	K8SDaemonSetNameKey = semconv.K8SDaemonSetNameKey

	// The uid of the Job.
	K8SJobUIDKey = semconv.K8SJobUIDKey

	// The name of the Job.
	K8SJobNameKey = semconv.K8SJobNameKey

	// The uid of the CronJob.
	K8SCronJobUIDKey = semconv.K8SCronJobUIDKey

	// The name of the CronJob.
	K8SCronJobNameKey = semconv.K8SCronJobNameKey
)

// Semantic conventions for OS resource attribute keys.
const (
	// The operating system type.
	OSTypeKey = semconv.OSTypeKey
	// Human readable (not intended to be parsed) OS version information.
	OSDescriptionKey = semconv.OSDescriptionKey
)

// Semantic conventions for host resource attribute keys.
const (
	// A uniquely identifying name for the host: 'hostname', FQDN, or user specified name
	HostNameKey = semconv.HostNameKey

	// Unique host ID. For cloud environments this will be the instance ID.
	HostIDKey = semconv.HostIDKey

	// Type of host. For cloud environments this will be the machine type.
	HostTypeKey = semconv.HostTypeKey

	// Name of the OS or VM image the host is running.
	HostImageNameKey = semconv.HostImageNameKey

	// Identifier of the image the host is running.
	HostImageIDKey = semconv.HostImageIDKey

	// Version of the image the host is running.
	HostImageVersionKey = semconv.HostImageVersionKey
)

// Semantic conventions for cloud environment resource attribute keys.
const (
	// Name of the cloud provider.
	CloudProviderKey = semconv.CloudProviderKey

	// The account ID from the cloud provider used for authorization.
	CloudAccountIDKey = semconv.CloudAccountIDKey

	// Geographical region where this resource is.
	CloudRegionKey = semconv.CloudRegionKey

	// Availability zone of the region where this resource is.
	CloudAvailabilityZoneKey = semconv.CloudAvailabilityZoneKey
)

// Semantic conventions for common cloud provider resource attributes.
var (
	CloudProviderAWS   = semconv.CloudProviderAWS
	CloudProviderAzure = semconv.CloudProviderAzure
	CloudProviderGCP   = semconv.CloudProviderGCP
)

// Semantic conventions for deployment attributes.
const (
	// Name of the deployment environment (aka deployment tier); e.g. (staging, production).
	DeploymentEnvironmentKey = semconv.DeploymentEnvironmentKey
)
