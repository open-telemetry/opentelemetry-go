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

package internaltest // import "go.opentelemetry.io/otel/internal/internaltest"

import (
	"fmt"
	"testing"
)

type ErrorHandler struct {
	errors []error
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) Handle(err error) {
	e.errors = append(e.errors, err)
}

func (e *ErrorHandler) Errors() []error {
	cp := make([]error, len(e.errors))
	copy(cp, e.errors)
	return cp
}

func (e *ErrorHandler) Len() int {
	return len(e.errors)
}

func (e *ErrorHandler) Reset() {
	if e.Len() > 0 {
		e.errors = e.errors[:0]
	}
}

func (e *ErrorHandler) RequireNoErrors(t *testing.T, msgAndArgs ...interface{}) {
	t.Helper()
	if e.hasErrors(t, msgAndArgs) {
		t.FailNow()
	}
}

func (e *ErrorHandler) AssertNoErrors(t *testing.T, msgAndArgs ...interface{}) bool {
	t.Helper()
	if e.hasErrors(t, msgAndArgs) {
		t.Fail()
		e.Reset()
		return false
	}
	return true
}

func (e *ErrorHandler) hasErrors(t *testing.T, msgAndArgs ...interface{}) bool {
	t.Helper()
	if n := e.Len(); n > 0 {
		t.Logf("Received unexpected errors (%d):\n%+v", n, e.errors)
		if msg := messageFromMsgAndArgs(msgAndArgs); len(msg) > 0 {
			t.Logf("Message: %s", msg)
		}
		return true
	}
	return false
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	// From https://github.com/stretchr/testify/blob/0ab3ce1249292a7221058b9e370472bca8f04813/assert/assertions.go#L179
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}
