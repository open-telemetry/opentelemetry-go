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

package global

import (
	"bytes"
	"errors"
	"io"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testErrCatcher []string

func (l *testErrCatcher) Write(p []byte) (int, error) {
	msg := bytes.TrimRight(p, "\n")
	(*l) = append(*l, string(msg))
	return len(msg), nil
}

func (l *testErrCatcher) Reset() {
	*l = testErrCatcher([]string{})
}

func (l *testErrCatcher) Got() []string {
	return []string(*l)
}

func causeErr(text string) {
	Handle(errors.New(text))
}

type HandlerTestSuite struct {
	suite.Suite

	origHandler ErrorHandler
	errCatcher  *testErrCatcher
}

func (s *HandlerTestSuite) SetupSuite() {
	s.errCatcher = new(testErrCatcher)
	s.origHandler = GlobalErrorHandler.getDelegate()

	GlobalErrorHandler.setDelegate(&ErrLogger{l: log.New(s.errCatcher, "", 0)})
}

func (s *HandlerTestSuite) TearDownSuite() {
	GlobalErrorHandler.setDelegate(s.origHandler)
}

func (s *HandlerTestSuite) SetupTest() {
	s.errCatcher.Reset()
}

func (s *HandlerTestSuite) TearDownTest() {
	GlobalErrorHandler.setDelegate(&ErrLogger{l: log.New(s.errCatcher, "", 0)})
}

func (s *HandlerTestSuite) TestGlobalHandler() {
	errs := []string{"one", "two"}
	GetErrorHandler().Handle(errors.New(errs[0]))
	Handle(errors.New(errs[1]))
	s.Assert().Equal(errs, s.errCatcher.Got())
}

func (s *HandlerTestSuite) TestDelegatedHandler() {
	eh := GetErrorHandler()

	newErrLogger := new(testErrCatcher)
	SetErrorHandler(&ErrLogger{l: log.New(newErrLogger, "", 0)})

	errs := []string{"TestDelegatedHandler"}
	eh.Handle(errors.New(errs[0]))
	s.Assert().Equal(errs, newErrLogger.Got())
}

func (s *HandlerTestSuite) TestNoDropsOnDelegate() {
	causeErr("")
	s.Require().Len(s.errCatcher.Got(), 1)

	// Change to another Handler. We are testing this is loss-less.
	newErrLogger := new(testErrCatcher)
	secondary := &ErrLogger{
		l: log.New(newErrLogger, "", 0),
	}
	SetErrorHandler(secondary)

	causeErr("")
	s.Assert().Len(s.errCatcher.Got(), 1, "original Handler used after delegation")
	s.Assert().Len(newErrLogger.Got(), 1, "new Handler not used after delegation")
}

func (s *HandlerTestSuite) TestAllowMultipleSets() {
	notUsed := new(testErrCatcher)

	secondary := &ErrLogger{l: log.New(notUsed, "", 0)}
	SetErrorHandler(secondary)
	s.Require().Same(GetErrorHandler(), GlobalErrorHandler, "set changed globalErrorHandler")
	s.Require().Same(GlobalErrorHandler.getDelegate(), secondary, "new Handler not set")

	tertiary := &ErrLogger{l: log.New(notUsed, "", 0)}
	SetErrorHandler(tertiary)
	s.Require().Same(GetErrorHandler(), GlobalErrorHandler, "set changed globalErrorHandler")
	s.Assert().Same(GlobalErrorHandler.getDelegate(), tertiary, "user Handler not overridden")
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func TestHandlerConcurrentSafe(t *testing.T) {
	// In order not to pollute the test output.
	SetErrorHandler(&ErrLogger{log.New(io.Discard, "", 0)})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		SetErrorHandler(&ErrLogger{log.New(io.Discard, "", 0)})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		Handle(errors.New("error"))
	}()

	wg.Wait()
	reset()
}

func BenchmarkErrorHandler(b *testing.B) {
	primary := &ErrLogger{l: log.New(io.Discard, "", 0)}
	secondary := &ErrLogger{l: log.New(io.Discard, "", 0)}
	tertiary := &ErrLogger{l: log.New(io.Discard, "", 0)}

	GlobalErrorHandler.setDelegate(primary)

	err := errors.New("benchmark error handler")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetErrorHandler().Handle(err)
		Handle(err)

		SetErrorHandler(secondary)
		GetErrorHandler().Handle(err)
		Handle(err)

		SetErrorHandler(tertiary)
		GetErrorHandler().Handle(err)
		Handle(err)

		GlobalErrorHandler.setDelegate(primary)
	}

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
	SetErrorHandler(&ErrLogger{l: log.New(io.Discard, "", 0)})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh = GetErrorHandler()
	}

	reset()
}

func BenchmarkDefaultErrorHandlerHandle(b *testing.B) {
	GlobalErrorHandler.setDelegate(
		&ErrLogger{l: log.New(io.Discard, "", 0)},
	)

	eh := GetErrorHandler()
	err := errors.New("benchmark default error handler handle")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh.Handle(err)
	}

	reset()
}

func BenchmarkDelegatedErrorHandlerHandle(b *testing.B) {
	eh := GetErrorHandler()
	SetErrorHandler(&ErrLogger{l: log.New(io.Discard, "", 0)})
	err := errors.New("benchmark delegated error handler handle")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh.Handle(err)
	}

	reset()
}

func BenchmarkSetErrorHandlerDelegation(b *testing.B) {
	alt := &ErrLogger{l: log.New(io.Discard, "", 0)}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetErrorHandler(alt)

		reset()
	}
}

func reset() {
	GlobalErrorHandler = defaultErrorHandler()
}
