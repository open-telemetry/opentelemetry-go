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

package internal

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	GlobalScope        = unsafe.Pointer(newAtomicValue())
	GlobalDelegateOnce = unsafe.Pointer(newSyncOnce())
)

type currentScopeKeyType struct{}

var currentScopeKey = &currentScopeKeyType{}

func SetScopeImpl(ctx context.Context, si interface{}) context.Context {
	return context.WithValue(ctx, currentScopeKey, si)
}

func ScopeImpl(ctx context.Context) interface{} {
	if ctx == nil {
		return nil
	}
	return ctx.Value(currentScopeKey)
}

func newAtomicValue() *atomic.Value {
	av := &atomic.Value{}
	av.Store(int(1))
	return av
}

func newSyncOnce() *sync.Once {
	return &sync.Once{}
}
