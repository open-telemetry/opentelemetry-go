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

// Package env provides a detector that loads resource information from
// the OTEL_RESOURCE environment variable. A list of labels of the form
// `<key1>=<value1>,<key2>=<value2>,...` is accepted. Domain names and
// paths are accepted as label keys.
package detect

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Environment variable used by "env" to decode a resource.
const envVar = "OTEL_RESOURCE"

// FromEnv is a detector that implments the ResourceDetector and collects resources
// from environment
type FromEnv struct{}

// Detect is a function that collects resources from environment
func (d *FromEnv) Detect(context.Context) (*resource.Resource, error) {
	labels := strings.TrimSpace(os.Getenv(envVar))

	if labels == "" {
		return resource.Empty(), nil
	}
	return constructOTResources(labels)
}

// labelRegex matches any key=value pair including a trailing comma or the end of the
// string. Captures the trimmed key & value parts, and ignores any superfluous spaces.
var labelRegex = regexp.MustCompile(`\s*([[:ascii:]]{1,256}?)\s*=\s*([[:ascii:]]{0,256}?)\s*(?:,|$)`)

func constructOTResources(s string) (*resource.Resource, error) {
	matches := labelRegex.FindAllStringSubmatchIndex(s, -1)

	if len(matches) == 0 {
		return resource.Empty(), fmt.Errorf("invalid resource format: %q", s)
	}

	labels := []kv.KeyValue{}
	prevIndex := 0

	for _, match := range matches {
		// if there is any text between matches, raise an error
		if prevIndex != match[0] {
			return resource.Empty(), fmt.Errorf("invalid resource format, invalid text: %q", s[prevIndex:match[0]])
		}

		key := s[match[2]:match[3]]
		value := s[match[4]:match[5]]

		var err error
		if value, err = url.QueryUnescape(value); err != nil {
			return resource.Empty(), fmt.Errorf("invalid resource format in attribute: %q, err: %s", s[match[0]:match[1]], err)
		}

		labels = append(labels, kv.String(key, value))

		prevIndex = match[1]
	}

	// if there is any text after the last match, raise an error
	if matches[len(matches)-1][1] != len(s) {
		return resource.Empty(), fmt.Errorf("invalid resource format, invalid text: %q", s[matches[len(matches)-1][1]:])
	}

	return resource.New(labels...), nil
}
