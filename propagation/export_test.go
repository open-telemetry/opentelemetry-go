// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package propagation

import "sync"

// ResetHandleExtractErrOnce resets handleExtractErrOnce for tests.
func ResetHandleExtractErrOnce() {
	handleExtractErrOnce = sync.Once{}
}
