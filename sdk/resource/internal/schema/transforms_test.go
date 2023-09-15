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

package schema

import "testing"

func TestTransformsOldesVersion(t *testing.T) {
	got := transforms[0].Version.String()
	if got != "1.4.0" {
		t.Errorf("oldest transform not v1.4.0: %s", got)
	}
}

func TestTransformsSorted(t *testing.T) {
	for i := len(transforms) - 1; i > 0; i-- {
		vI, vIMinus1 := transforms[i].Version, transforms[i-1].Version
		if vI.LessThan(vIMinus1) {
			t.Errorf("transforms are not in sorted order: %s is not less than %s", vI, vIMinus1)
		}
	}
}
