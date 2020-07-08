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

package detect

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
)

// ResourceDetector attempts to detect resource information.
// If the detector cannot find resource information, the returned resource is nil but no
// error is returned.
// An error is only returned on unexpected failures.
type ResourceDetector func(ctx context.Context) (*resource.Resource, error)

// AutoDetect calls all input detectors sequentially and merges each result with the previous one.
// It returns on the first error that a sub-detector encounters.
func AutoDetect(ctx context.Context, detectors ...ResourceDetector) (*resource.Resource, error) {
	var autoDetectedRes *resource.Resource
	for _, detector := range detectors {
		res, err := detector(ctx)
		if err != nil {
			return nil, err
		}
		autoDetectedRes = resource.Merge(autoDetectedRes, res)
	}
	return autoDetectedRes, nil
}
