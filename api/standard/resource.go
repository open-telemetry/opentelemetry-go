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

// Constants for Service resources.
const (
	// A uniquely identifying name for a Service.
	ServiceName       = kv.Key("service.name")
	ServiceNamespace  = kv.Key("service.namespace")
	ServiceInstanceID = kv.Key("service.instance.id")
	ServiceVersion    = kv.Key("service.version")
)

// Constants for Library resources.
const (
	// A uniquely identifying name for a Library.
	LibraryName     = kv.Key("library.name")
	LibraryLanguage = kv.Key("library.language")
	LibraryVersion  = kv.Key("library.version")
)

// Constants for Kubernetes resources.
const (
	// A uniquely identifying name for the Kubernetes cluster. Kubernetes
	// does not have cluster names as an internal concept so this may be
	// set to any meaningful value within the environment. For example,
	// GKE clusters have a name which can be used for this label.
	K8SClusterName    = kv.Key("k8s.cluster.name")
	K8SNamespaceName  = kv.Key("k8s.namespace.name")
	K8SPodName        = kv.Key("k8s.pod.name")
	K8SDeploymentName = kv.Key("k8s.deployment.name")
)

// Constants for Container resources.
const (
	// A uniquely identifying name for the Container.
	ContainerName      = kv.Key("container.name")
	ContainerImageName = kv.Key("container.image.name")
	ContainerImageTag  = kv.Key("container.image.tag")
)

// Constants for Cloud resources.
const (
	CloudProvider  = kv.Key("cloud.provider")
	CloudAccountID = kv.Key("cloud.account.id")
	CloudRegion    = kv.Key("cloud.region")
	CloudZone      = kv.Key("cloud.zone")
)

// Constants for Host resources.
const (
	// A uniquely identifying name for the host.
	HostName = kv.Key("host.name")

	// A hostname as returned by the 'hostname' command on host machine.
	HostHostName     = kv.Key("host.hostname")
	HostID           = kv.Key("host.id")
	HostType         = kv.Key("host.type")
	HostImageName    = kv.Key("host.image.name")
	HostImageID      = kv.Key("host.image.id")
	HostImageVersion = kv.Key("host.image.version")
)
