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

package gce

import (
	"cloud.google.com/go/compute/metadata"
	"context"
	"go.opencensus.io/resource/resourcekeys" // TODO: replace by "go.opentelemetry.io/otel/sdk/resource/resourcekeys" after publish
	"strings"

	"log"
	"strings"
)

// Detect detects associated resources when running on GCE hosts.
func Detect(ctx context.Context) (*resource.Resource, error) {
	if !metadata.OnGCE() {
		return nil, nil
	}

	labels := []kv.KeyValue{}

	labels = append(labels, kv.String(resourcekeys.CloudKeyProvider, resourcekeys.CloudProviderGCP))

	projectID, err := metadata.ProjectID()
	logError(err)
	if projectID != "" {
		labels = append(labels, kv.String(resourcekeys.CloudKeyAccountID, projectID))
	}

	zone, err := metadata.Zone()
	logError(err)
	if zone != "" {
		labels = append(labels, kv.String(resourcekeys.CloudKeyZone, zone))
	}

	labels = append(labels, kv.String(resourcekeys.CloudKeyRegion, ""))

	instanceID, err := metadata.InstanceID()
	logError(err)
	if instanceID != "" {
		labels = append(labels, kv.String(resourcekeys.HostKeyID, instanceID))
	}

	name, err := metadata.InstanceName()
	logError(err)
	if instanceID != "" {
		labels = append(labels, kv.String(resourcekeys.HostKeyName, name))
	}

	hostname, err := metadata.Hostname()
	logError(err)
	if instanceID != "" {
		labels = append(labels, kv.String(resourcekeys.HostKeyHostName, hostname))
	}

	hostType, err := metadata.InstanceAttributeValue("instance/machine-type")
	logError(err)
	if instanceID != "" {
		labels = append(labels, kv.String(resourcekeys.HostKeyType, hostType))
	}

	return resource.New(labels...), nil
}

// logError logs error only if the error is present and it is not 'not defined'
func logError(err error) {
	if err != nil {
		if !strings.Contains(err.Error(), "not defined") {
			log.Printf("Error retrieving gcp metadata: %v", err)
		}
	}
}
