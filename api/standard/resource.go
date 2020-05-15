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

package standard // import "go.opentelemetry.io/otel/api/standard"

import "go.opentelemetry.io/otel/api/kv"

// Standard service resource attribute keys.
const (
	// Name of the service.
	ServiceNameKey = kv.Key("service.name")

	// A namespace for `service.name`. This needs to have meaning that helps
	// to distinguish a group of services. For example, the team name that
	// owns a group of services. `service.name` is expected to be unique
	// within the same namespace.
	ServiceNamespaceKey = kv.Key("service.namespace")

	// A unique identifier of the service instance. In conjunction with the
	// `service.name` and `service.namespace` this must be unique.
	ServiceInstanceIDKey = kv.Key("service.instance.id")

	// The version of the service API.
	ServiceVersionKey = kv.Key("service.version")
)

// Standard telemetry SDK resource attribute keys.
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
	TelemetrySDKNameKey = kv.Key("telemetry.sdk.name")

	// The language of the telemetry SDK.
	TelemetrySDKLanguageKey = kv.Key("telemetry.sdk.language")

	// The version string of the telemetry SDK.
	TelemetrySDKVersionKey = kv.Key("telemetry.sdk.version")
)

// Standard telemetry SDK resource attributes.
var (
	TelemetrySDKLanguageGo = TelemetrySDKLanguageKey.String("go")
)

// Standard container resource attribute keys.
const (
	// A uniquely identifying name for the Container.
	ContainerNameKey = kv.Key("container.name")

	// Name of the image the container was built on.
	ContainerImageNameKey = kv.Key("container.image.name")

	// Container image tag.
	ContainerImageTagKey = kv.Key("container.image.tag")
)

// Standard Function-as-a-Service resource attribute keys.
const (
	// A uniquely identifying name for the FaaS.
	FaaSName = kv.Key("faas.name")

	// The unique name of the function being executed.
	FaaSID = kv.Key("faas.id")

	// The version of the function being executed.
	FaaSVersion = kv.Key("faas.version")

	// The execution environment identifier.
	FaaSInstance = kv.Key("faas.instance")
)

// Standard Kubernetes resource attribute keys.
const (
	// A uniquely identifying name for the Kubernetes cluster. Kubernetes
	// does not have cluster names as an internal concept so this may be
	// set to any meaningful value within the environment. For example,
	// GKE clusters have a name which can be used for this label.
	K8SClusterNameKey = kv.Key("k8s.cluster.name")

	// The name of the namespace that the pod is running in.
	K8SNamespaceNameKey = kv.Key("k8s.namespace.name")

	// The name of the pod.
	K8SPodNameKey = kv.Key("k8s.pod.name")

	// The name of the deployment.
	K8SDeploymentNameKey = kv.Key("k8s.deployment.name")
)

// Standard host resource attribute keys.
const (
	// A uniquely identifying name for the host.
	HostNameKey = kv.Key("host.name")

	// A hostname as returned by the 'hostname' command on host machine.
	HostHostNameKey = kv.Key("host.hostname")

	// Unique host ID. For cloud environments this will be the instance ID.
	HostIDKey = kv.Key("host.id")

	// Type of host. For cloud environments this will be the machine type.
	HostTypeKey = kv.Key("host.type")

	// Name of the OS or VM image the host is running.
	HostImageNameKey = kv.Key("host.image.name")

	// Identifier of the image the host is running.
	HostImageIDKey = kv.Key("host.image.id")

	// Version of the image the host is running.
	HostImageVersionKey = kv.Key("host.image.version")
)

// Standard cloud environment resource attribute keys.
const (
	// Name of the cloud provider.
	CloudProviderKey = kv.Key("cloud.provider")

	// The account ID from the cloud provider used for authorization.
	CloudAccountIDKey = kv.Key("cloud.account.id")

	// Geographical region where this resource is.
	CloudRegionKey = kv.Key("cloud.region")

	// Zone of the region where this resource is.
	CloudZoneKey = kv.Key("cloud.zone")
)

var (
	CloudProviderAWS   = CloudProviderKey.String("aws")
	CloudProviderAzure = CloudProviderKey.String("azure")
	CloudProviderGCP   = CloudProviderKey.String("gcp")
)
