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

// Package internal contains common functionality for all OTLP exporters.
package internal // import "go.opentelemetry.io/otel/exporters/otlp/internal"

import (
	"fmt"
	"path"
	"strings"
)

// CleanPath returns cleaned URL path. Replace with default path if path is nil
func CleanPath(URLPath string, defaultPath string) string {
	tmp := strings.TrimSpace(URLPath)
	if tmp == "" {
		return defaultPath
	}
	tmp = path.Clean(tmp)
	if !path.IsAbs(tmp) {
		tmp = fmt.Sprintf("/%s", tmp)
	}
	return tmp
}
