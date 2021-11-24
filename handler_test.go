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
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
)

type errLogger []string

func (l *errLogger) Write(p []byte) (int, error) {
	msg := bytes.TrimRight(p, "\n")
	(*l) = append(*l, string(msg))
	return len(msg), nil
}

func (l *errLogger) Reset() {
	*l = errLogger([]string{})
}

func (l *errLogger) Got() []string {
	return []string(*l)
}

type logger struct {
	l *log.Logger
}

func (l *logger) Handle(err error) {
	l.l.Print(err)
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
	s.errLogger = new(errLogger)
	s.origHandler = globalErrorHandler.Load().(holder).eh

	globalErrorHandler.Store(holder{eh: &delegator{l: log.New(s.errLogger, "", 0)}})
}

func (s *HandlerTestSuite) TearDownSuite() {
	globalErrorHandler.Store(holder{eh: s.origHandler})
	delegateErrorHandlerOnce = sync.Once{}
}

func (s *HandlerTestSuite) SetupTest() {
	s.errLogger.Reset()
}

func (s *HandlerTestSuite) TearDownTest() {
	globalErrorHandler.Store(holder{eh: &delegator{l: log.New(s.errLogger, "", 0)}})
	delegateErrorHandlerOnce = sync.Once{}
}

func (s *HandlerTestSuite) TestGlobalHandler() {
	errs := []string{"one", "two"}
	GetErrorHandler().Handle(errors.New(errs[0]))
	Handle(errors.New(errs[1]))
	s.Assert().Equal(errs, s.errLogger.Got())
}

func (s *HandlerTestSuite) TestDelegatedHandler() {
	eh := GetErrorHandler()

	newErrLogger := new(errLogger)
	SetErrorHandler(&logger{l: log.New(newErrLogger, "", 0)})

	errs := []string{"TestDelegatedHandler"}
	eh.Handle(errors.New(errs[0]))
	s.Assert().Equal(errs, newErrLogger.Got())
}

func (s *HandlerTestSuite) TestSettingDefaultIsANoOp() {
	SetErrorHandler(GetErrorHandler())
	d := globalErrorHandler.Load().(holder).eh.(*delegator)
	s.Assert().Nil(d.delegate.Load())
}

func (s *HandlerTestSuite) TestNoDropsOnDelegate() {
	causeErr("")
	s.Require().Len(s.errLogger.Got(), 1)

	// Change to another Handler. We are testing this is loss-less.
	newErrLogger := new(errLogger)
	secondary := &logger{
		l: log.New(newErrLogger, "", 0),
	}
	SetErrorHandler(secondary)

	causeErr("")
	s.Assert().Len(s.errLogger.Got(), 1, "original Handler used after delegation")
	s.Assert().Len(newErrLogger.Got(), 1, "new Handler not used after delegation")
}

func (s *HandlerTestSuite) TestAllowMultipleSets() {
	notUsed := new(errLogger)

	secondary := &logger{l: log.New(notUsed, "", 0)}
	SetErrorHandler(secondary)
	s.Require().Same(GetErrorHandler(), secondary, "new Handler not set")

	tertiary := &logger{l: log.New(notUsed, "", 0)}
	SetErrorHandler(tertiary)
	s.Assert().Same(GetErrorHandler(), tertiary, "user Handler not overridden")
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func BenchmarkErrorHandler(b *testing.B) {
	primary := &delegator{l: log.New(ioutil.Discard, "", 0)}
	secondary := &logger{l: log.New(ioutil.Discard, "", 0)}
	tertiary := &logger{l: log.New(ioutil.Discard, "", 0)}

	globalErrorHandler.Store(holder{eh: primary})

	err := errors.New("BenchmarkErrorHandler")

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

		b.StopTimer()
		primary.delegate = atomic.Value{}
		globalErrorHandler.Store(holder{eh: primary})
		delegateErrorHandlerOnce = sync.Once{}
		b.StartTimer()
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
	SetErrorHandler(&logger{l: log.New(ioutil.Discard, "", 0)})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh = GetErrorHandler()
	}

	b.StopTimer()
	reset()
}

func BenchmarkDefaultErrorHandlerHandle(b *testing.B) {
	globalErrorHandler.Store(holder{
		eh: &delegator{l: log.New(ioutil.Discard, "", 0)},
	})

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
	SetErrorHandler(&logger{l: log.New(ioutil.Discard, "", 0)})
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
	alt := &logger{l: log.New(ioutil.Discard, "", 0)}

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
		&logger{l: log.New(ioutil.Discard, "", 0)},
		&logger{l: log.New(ioutil.Discard, "", 0)},
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
	globalErrorHandler = defaultErrorHandler()
	delegateErrorHandlerOnce = sync.Once{}
}
