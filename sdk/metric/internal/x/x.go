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

// Package x contains support for OTel metric SDK experimental features.
package x // import "go.opentelemetry.io/otel/sdk/metric/internal/x"

import (
	"os"
	"strings"
)

const EnvKeyRoot = "OTEL_GO_X_"

var (
	CardinalityLimit = Feature{
		EnvKeySuffix: "CARDINALITY_LIMIT",
		// TODO: support accepting number values here to set the cardinality
		// limit.
		EnablementVals: []string{"true"},
	}
)

type Feature struct {
	// EnvKeySuffix is the environment variable key suffix the xFeature is
	// stored at. It is assumed EnvKeyRoot is the base of the environment
	// variable key.
	EnvKeySuffix string
	// EnablementVals are the case-insensitive comparison values that indicate
	// the Feature is enabled.
	EnablementVals []string
}

// Enabled returns if the Feature is enabled.
func Enabled(f Feature) bool {
	key := EnvKeyRoot + f.EnvKeySuffix
	vRaw, present := os.LookupEnv(key)
	if !present {
		return false
	}

	v := strings.ToLower(vRaw)
	for _, allowed := range f.EnablementVals {
		if v == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}
