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
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type errLogger struct {
	errs []string
}

func (l *errLogger) Handle(err error) {
	l.errs = append(l.errs, err.Error())
}

func causeErr(text string) {
	Handle(errors.New(text))
}

type HandlerTestSuite struct {
	suite.Suite

	origHandler ErrorHandler
	errLogger   *errLogger
}

func (s *HandlerTestSuite) SetupSuite() {
	s.errLogger = &errLogger{errs: []string{}}
	s.origHandler = globalErrorHandler.delegate

	globalErrorHandler.delegate = s.errLogger
}

func (s *HandlerTestSuite) TearDownSuite() {
	globalErrorHandler.delegate = s.origHandler
	delegateErrorHandlerOnce = sync.Once{}
}

func (s *HandlerTestSuite) SetupTest() {
	s.errLogger.errs = []string{}
	SetErrorHandler(s.errLogger)
}

func (s *HandlerTestSuite) TearDownTest() {}

type bufferedErrorHandler struct {
	buf *strings.Builder
}

func (h *bufferedErrorHandler) Handle(err error) {
	if h.buf != nil {
		h.buf = &strings.Builder{}
	}
	h.buf.WriteString(err.Error())
}

func (s *HandlerTestSuite) TestGlobalHandler() {
	errs := []string{"one", "two"}
	GetErrorHandler().Handle(errors.New(errs[0]))
	Handle(errors.New(errs[1]))
	s.Assert().Equal(errs, s.errLogger.errs)
}

func (s *HandlerTestSuite) TestDelegatedHandler() {
	eh := GetErrorHandler()

	newErrLogger := &errLogger{errs: []string{}}
	SetErrorHandler(newErrLogger)

	errs := []string{"TestDelegatedHandler"}
	eh.Handle(errors.New(errs[0]))
	s.Assert().Equal(errs, newErrLogger.errs)
}

func (s *HandlerTestSuite) TestNoDropsOnDelegate() {
	causeErr("")
	s.Require().Len(s.errLogger.errs, 1)

	// Change to another Handler. We are testing this is loss-less.
	newErrLogger := &errLogger{errs: []string{}}
	SetErrorHandler(newErrLogger)

	causeErr("")
	s.Assert().Len(s.errLogger.errs, 1, "original Handler used after delegation")
	s.Assert().Len(newErrLogger.errs, 1, "new Handler not used after delegation")
}

func (s *HandlerTestSuite) TestAllowMultipleSets() {
	secondary := &errLogger{errs: []string{}}
	SetErrorHandler(secondary)
	s.Require().Same(GetErrorHandler().(*errorHandlerDelegate).delegate, secondary, "new Handler not set")

	tertiary := &errLogger{errs: []string{}}
	SetErrorHandler(tertiary)
	s.Assert().Same(GetErrorHandler().(*errorHandlerDelegate).delegate, tertiary, "user Handler not overridden")
}

func (s *HandlerTestSuite) TestGetErrorHandlerAlwaysIsCurrent() {
	orig := GetErrorHandler()

	newErrLogger := &errLogger{errs: []string{}}
	SetErrorHandler(newErrLogger)

	orig.Handle(errors.New("error"))

	s.Assert().Len(newErrLogger.errs, 1, "original Handler did not update")
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func BenchmarkErrorHandler(b *testing.B) {
	primary := DiscardErrorHandler{}
	secondary := DiscardErrorHandler{}
	tertiary := DiscardErrorHandler{}

	err := errors.New("BenchmarkErrorHandler")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetErrorHandler(primary)
		GetErrorHandler().Handle(err)
		Handle(err)

		SetErrorHandler(secondary)
		GetErrorHandler().Handle(err)
		Handle(err)

		SetErrorHandler(tertiary)
		GetErrorHandler().Handle(err)
		Handle(err)

	}

	b.StopTimer()
	reset()
}

var eh ErrorHandler

func BenchmarkGetDefaultErrorHandler(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		eh = GetErrorHandler()
	}
}

func BenchmarkGetDelegatedErrorHandler(b *testing.B) {
	SetErrorHandler(DiscardErrorHandler{})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh = GetErrorHandler()
	}

	b.StopTimer()
	reset()
}

func BenchmarkDefaultErrorHandlerHandle(b *testing.B) {
	SetErrorHandler(DiscardErrorHandler{})

	eh := GetErrorHandler()
	err := errors.New("BenchmarkDefaultErrorHandlerHandle")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh.Handle(err)
	}

	b.StopTimer()
	reset()
}

func BenchmarkDelegatedErrorHandlerHandle(b *testing.B) {
	eh := GetErrorHandler()
	SetErrorHandler(DiscardErrorHandler{})
	err := errors.New("BenchmarkDelegatedErrorHandlerHandle")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh.Handle(err)
	}

	b.StopTimer()
	reset()
}

func BenchmarkSetErrorHandlerDelegation(b *testing.B) {
	alt := DiscardErrorHandler{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetErrorHandler(alt)

		b.StopTimer()
		reset()
		b.StartTimer()
	}
}

func BenchmarkSetErrorHandlerNoDelegation(b *testing.B) {
	eh := []ErrorHandler{
		DiscardErrorHandler{},
		DiscardErrorHandler{},
	}
	mod := len(eh)
	// Do not measure delegation.
	SetErrorHandler(eh[1])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetErrorHandler(eh[i%mod])
	}

	b.StopTimer()
	reset()
}

func reset() {
	globalErrorHandler = &errorHandlerDelegate{
		delegate: &defaultErrorHandler{},
	}
}
