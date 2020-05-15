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

// Package resourcekey contains standard resource attribute keys as defined
// by the OpenTelemetry specification
// (https://github.com/open-telemetry/opentelemetry-specification/tree/v0.4.0/specification/resource/semantic_conventions).
package resourcekey // import "go.opentelemetry.io/otel/sdk/resource/resourcekey"

// Constants for Service resources.
const (
	// A uniquely identifying name for a Service.
	ServiceName       = "service.name"
	ServiceNamespace  = "service.namespace"
	ServiceInstanceID = "service.instance.id"
	ServiceVersion    = "service.version"
)

// Constants for Library resources.
const (
	// A uniquely identifying name for a Library.
	LibraryName     = "library.name"
	LibraryLanguage = "library.language"
	LibraryVersion  = "library.version"
)

// Constants for Kubernetes resources.
const (
	// A uniquely identifying name for the Kubernetes cluster. Kubernetes
	// does not have cluster names as an internal concept so this may be
	// set to any meaningful value within the environment. For example,
	// GKE clusters have a name which can be used for this label.
	K8SClusterName    = "k8s.cluster.name"
	K8SNamespaceName  = "k8s.namespace.name"
	K8SPodName        = "k8s.pod.name"
	K8SDeploymentName = "k8s.deployment.name"
)

// Constants for Container resources.
const (
	// A uniquely identifying name for the Container.
	ContainerName      = "container.name"
	ContainerImageName = "container.image.name"
	ContainerImageTag  = "container.image.tag"
)

// Constants for Cloud resources.
const (
	CloudProvider  = "cloud.provider"
	CloudAccountID = "cloud.account.id"
	CloudRegion    = "cloud.region"
	CloudZone      = "cloud.zone"

	// Cloud Providers
	CloudProviderAWS   = "aws"
	CloudProviderGCP   = "gcp"
	CloudProviderAZURE = "azure"
)

// Constants for Host resources.
const (
	// A uniquely identifying name for the host.
	HostName = "host.name"

	// A hostname as returned by the 'hostname' command on host machine.
	HostHostName     = "host.hostname"
	HostID           = "host.id"
	HostType         = "host.type"
	HostImageName    = "host.image.name"
	HostImageID      = "host.image.id"
	HostImageVersion = "host.image.version"
)
