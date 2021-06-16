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

package main

import (
	"sync/atomic"
)

type A struct {
	_ int64 // b is accessed via a sync/atomic function
	_ int32
	d int64 // d should NOT be accessed via a sync/atomic function as it's not 64bits aligned.
}

func main() {
	a := A{}
	// Atomically increment a.b. This is valid, b is the first field of the struct, so it's
	// guaranteed to be 64-bits aligned.
	atomic.AddInt64(&a.d, 1)
}
