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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"os"
	"strings"
)

const xEnvKeyRoot = "OTEL_GO_X_"

var (
	xExemplar = xFeature{
		envKeySuffix:   "EXEMPLAR",
		enablementVals: []string{"true"},
	}
)

type xFeature struct {
	// envKey is the environment variable key suffix the xFeature is stored at.
	// It is assumed xEnvKeyRoot is the base of the environment variable key.
	envKeySuffix string
	// enablementVals are the case-insensitive comparison values that indicate
	// the xFeature is enabled.
	enablementVals []string
}

// xEnabled returns if the xFeature is enabled.
func xEnabled(xf xFeature) bool {
	key := xEnvKeyRoot + xf.envKeySuffix
	vRaw, present := os.LookupEnv(key)
	if !present {
		return false
	}

	v := strings.ToLower(vRaw)
	for _, allowed := range xf.enablementVals {
		if v == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}
