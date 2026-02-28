// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package errorhandler

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fnErrHandler func(error)

func (f fnErrHandler) Handle(err error) { f(err) }

var noopEH = fnErrHandler(func(error) {})

type nonComparableErrorHandler struct {
	ErrorHandler

	nonComparable func() //nolint:unused  // This is not called.
}

func resetForTest(t testing.TB) {
	t.Cleanup(func() {
		globalErrorHandler = defaultErrorHandler()
		delegateErrorHandlerOnce = sync.Once{}
	})
}

func TestErrDelegator(t *testing.T) {
	buf := new(bytes.Buffer)
	log.Default().SetOutput(buf)
	t.Cleanup(func() { log.Default().SetOutput(os.Stderr) })

	e := &ErrDelegator{}

	err := errors.New("testing")
	e.Handle(err)

	got := buf.String()
	if !strings.Contains(got, err.Error()) {
		t.Error("default handler did not log")
	}
	buf.Reset()

	var gotErr error
	e.setDelegate(fnErrHandler(func(e error) { gotErr = e }))
	e.Handle(err)

	if buf.String() != "" {
		t.Error("delegate not set")
	} else if !errors.Is(gotErr, err) {
		t.Error("error not passed to delegate")
	}
}

func TestSetErrorHandler(t *testing.T) {
	t.Run("Set With default is a noop", func(t *testing.T) {
		resetForTest(t)
		SetErrorHandler(GetErrorHandler())

		eh, ok := GetErrorHandler().(*ErrDelegator)
		if !ok {
			t.Fatal("Global ErrorHandler should be the default ErrorHandler")
		}

		if eh.delegate.Load() != nil {
			t.Fatal("ErrorHandler should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		resetForTest(t)

		SetErrorHandler(noopEH)

		_, ok := GetErrorHandler().(*ErrDelegator)
		if ok {
			t.Fatal("Global ErrorHandler was not changed")
		}
	})

	t.Run("Set() should delegate existing ErrorHandlers", func(t *testing.T) {
		resetForTest(t)

		eh := GetErrorHandler()
		SetErrorHandler(noopEH)

		errDel, ok := eh.(*ErrDelegator)
		if !ok {
			t.Fatal("Wrong ErrorHandler returned")
		}

		if errDel.delegate.Load() == nil {
			t.Fatal("The ErrDelegator should have a delegate")
		}
	})

	t.Run("non-comparable types should not panic", func(t *testing.T) {
		resetForTest(t)

		eh := nonComparableErrorHandler{}
		assert.NotPanics(t, func() { SetErrorHandler(eh) }, "delegate")
		assert.NotPanics(t, func() { SetErrorHandler(eh) }, "replacement")
	})
}
