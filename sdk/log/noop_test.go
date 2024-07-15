// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"reflect"
	"testing"
)

func TestNoopProcessorNoPanics(t *testing.T) {
	assertAllExportedMethodNoPanic(
		reflect.ValueOf(NewNoopProcessor()),
		reflect.TypeOf((*Processor)(nil)).Elem(),
	)(t)
}
