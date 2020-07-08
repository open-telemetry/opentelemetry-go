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
	"context"
	"log"
	"strings"

	"cloud.google.com/go/compute/metadata"
	// TODO: import "go.opentelemetry.io/otel/sdk/resource/resourcekeys" after publishing it
	// for now, the resourcekeys is in const.go

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Detect detects associated resources when running on GCE hosts.
func Detect(ctx context.Context) (*resource.Resource, error) {
	if !metadata.OnGCE() {
		return nil, nil
	}

	labels := []kv.KeyValue{}

	labels = append(labels, kv.String(CloudKeyProvider, CloudProviderGCP))

	projectID, err := metadata.ProjectID()
	logError(err)
	if projectID != "" {
		labels = append(labels, kv.String(CloudKeyAccountID, projectID))
	}

	zone, err := metadata.Zone()
	logError(err)
	if zone != "" {
		labels = append(labels, kv.String(CloudKeyZone, zone))
	}

	labels = append(labels, kv.String(CloudKeyRegion, ""))

	instanceID, err := metadata.InstanceID()
	logError(err)
	if instanceID != "" {
		labels = append(labels, kv.String(HostKeyID, instanceID))
	}

	name, err := metadata.InstanceName()
	logError(err)
	if instanceID != "" {
		labels = append(labels, kv.String(HostKeyName, name))
	}

	hostname, err := metadata.Hostname()
	logError(err)
	if instanceID != "" {
		labels = append(labels, kv.String(HostKeyHostName, hostname))
	}

	hostType, err := metadata.InstanceAttributeValue("instance/machine-type")
	logError(err)
	if instanceID != "" {
		labels = append(labels, kv.String(HostKeyType, hostType))
	}

	return resource.New(labels...), nil

}

//logError logs error only if the error is present and it is not 'not defined'
func logError(err error) {
	if err != nil {
		if !strings.Contains(err.Error(), "not defined") {
			log.Printf("Error retrieving gcp metadata: %v", err)
		}
	}
}
