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

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource/internal/schema"
	"go.opentelemetry.io/otel/sdk/resource/internal/schema/resconv"
)

var errUnknownSchema = errors.New("unknown schema")

// upgradeResource returns a copy of orig with the schema URL set to schemaURL
// and all attributes transformed based on the associated schema. If the schema
// transformation fails an error is returned.
func upgradeResource(schemaURL string, r *Resource) (*Resource, error) {
	if r == nil || r.Len() == 0 {
		return NewWithAttributes(schemaURL), nil
	}

	if r.SchemaURL() == schemaURL {
		// Resources are immutable, just return the ptr to the same value.
		return r, nil
	}

	s, ok := schema.Schemas[schemaURL]
	if !ok {
		return nil, fmt.Errorf("%w: %s", errUnknownSchema, schemaURL)
	}
	attrs := r.Attributes()
	if err := resconv.Upgrade(s, attrs); err != nil {
		return nil, err
	}
	return NewWithAttributes(schemaURL, attrs...), nil
}
