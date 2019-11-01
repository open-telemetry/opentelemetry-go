// Copyright 2019, OpenTelemetry Authors
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

package array

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"go.opentelemetry.io/api/core"
)

func TestRawFloatSort(t *testing.T) {
	floats := []core.Number{
		core.NewFloat64Number(0),
		core.NewFloat64Number(+1. * 0.),
		core.NewFloat64Number(1. / math.Inf(-1)),
		core.NewFloat64Number(1. / math.Inf(+1)),
		core.NewFloat64Number(+1),
		core.NewFloat64Number(-1),
	}

	sort.Slice(floats, func(i, j int) bool {
		return int64(floats[i]) < int64(floats[j])
	})

	for _, n := range floats {
		fmt.Println("F=", n.Emit(core.Float64NumberKind))
	}

	panic("Anyway")
}
