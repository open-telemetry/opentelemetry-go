// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/noop_helper_test.go.tmpl

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package noop

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertAllExportedMethodNoPanic(rVal reflect.Value, rType reflect.Type) func(*testing.T) {
	return func(t *testing.T) {
		for n := 0; n < rType.NumMethod(); n++ {
			mType := rType.Method(n)
			if !mType.IsExported() {
				t.Logf("ignoring unexported %s", mType.Name)
				continue
			}
			m := rVal.MethodByName(mType.Name)
			if !m.IsValid() {
				t.Errorf("unknown method for %s: %s", rVal.Type().Name(), mType.Name)
			}

			numIn := mType.Type.NumIn()
			if mType.Type.IsVariadic() {
				numIn--
			}
			args := make([]reflect.Value, numIn)
			ctx := context.Background()
			for i := range args {
				aType := mType.Type.In(i)
				if aType.Name() == "Context" {
					// Do not panic on a nil context.
					args[i] = reflect.ValueOf(ctx)
				} else {
					args[i] = reflect.New(aType).Elem()
				}
			}

			assert.NotPanicsf(t, func() {
				_ = m.Call(args)
			}, "%s.%s", rVal.Type().Name(), mType.Name)
		}
	}
}
