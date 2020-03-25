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

// Package resourcekeys contains well known type and label keys for resources.
package resourcekeys // import "go.opentelemetry.io/otel/sdk/resource/resourcekeys"

// Constants for Service resources.
const (
	// A uniquely identifying name for a Service.
	ServiceKeyName       = "service.name"
	ServiceKeyNamespace  = "service.namespace"
	ServiceKeyInstanceID = "service.instance.id"
	ServiceKeyVersion    = "service.version"
)

// Constants for Library resources.
const (
	// A uniquely identifying name for a Library.
	LibraryKeyName     = "library.name"
	LibraryKeyLanguage = "library.language"
	LibraryKeyVersion  = "library.version"
)

// Constants for Kubernetes resources.
const (
	// A uniquely identifying name for the Kubernetes cluster. Kubernetes
	// does not have cluster names as an internal concept so this may be
	// set to any meaningful value within the environment. For example,
	// GKE clusters have a name which can be used for this label.
	K8SKeyClusterName    = "k8s.cluster.name"
	K8SKeyNamespaceName  = "k8s.namespace.name"
	K8SKeyPodName        = "k8s.pod.name"
	K8SKeyDeploymentName = "k8s.deployment.name"
)

// Constants for Container resources.
const (
	// A uniquely identifying name for the Container.
	ContainerKeyName      = "container.name"
	ContainerKeyImageName = "container.image.name"
	ContainerKeyImageTag  = "container.image.tag"
)

// Constants for Cloud resources.
const (
	CloudKeyProvider  = "cloud.provider"
	CloudKeyAccountID = "cloud.account.id"
	CloudKeyRegion    = "cloud.region"
	CloudKeyZone      = "cloud.zone"

	// Cloud Providers
	CloudProviderAWS   = "aws"
	CloudProviderGCP   = "gcp"
	CloudProviderAZURE = "azure"
)

// Constants for Host resources.
const (
	// A uniquely identifying name for the host.
	HostKeyName = "host.name"

	// A hostname as returned by the 'hostname' command on host machine.
	HostKeyHostName     = "host.hostname"
	HostKeyID           = "host.id"
	HostKeyType         = "host.type"
	HostKeyImageName    = "host.image.name"
	HostKeyImageID      = "host.image.id"
	HostKeyImageVersion = "host.image.version"
)
