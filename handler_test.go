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

package otel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testErrHandler struct {
	err error
}

var _ ErrorHandler = &testErrHandler{}

func (eh *testErrHandler) Handle(err error) { eh.err = err }

func TestGlobalErrorHandler(t *testing.T) {
	e1 := &testErrHandler{}
	SetErrorHandler(e1)
	Handle(assert.AnError)
	assert.ErrorIs(t, e1.err, assert.AnError)
	e1.err = nil

	e2 := &testErrHandler{}
	SetErrorHandler(e2)
	GetErrorHandler().Handle(assert.AnError)
	assert.ErrorIs(t, e2.err, assert.AnError)
}
