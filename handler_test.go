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
	"log"
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

type HandlerTestSuite struct {
	suite.Suite

	origHandler *loggingErrorHandler
	errLogger   *errLogger
}

func (s *HandlerTestSuite) SetupSuite() {
	s.errLogger = new(errLogger)
	s.origHandler = globalErrorHandler
	globalErrorHandler = &loggingErrorHandler{
		l: log.New(s.errLogger, "", 0),
	}
}

func (s *HandlerTestSuite) TearDownSuite() {
	globalErrorHandler = s.origHandler
}

func (s *HandlerTestSuite) SetupTest() {
	s.errLogger.Reset()
}

func (s *HandlerTestSuite) TestGlobalHandler() {
	errs := []string{"one", "two"}
	GetErrorHandler().Handle(errors.New(errs[0]))
	Handle(errors.New(errs[1]))
	s.Assert().Equal(errs, s.errLogger.Got())
}

func (s *HandlerTestSuite) TestNoDropsOnDelegate() {
	causeErr := func() func() {
		err := errors.New("")
		return func() {
			Handle(err)
		}
	}()

	causeErr()
	s.Require().Len(s.errLogger.Got(), 1)

	// Change to another Handler. We are testing this is loss-less.
	newErrLogger := new(errLogger)
	secondary := &loggingErrorHandler{
		l: log.New(newErrLogger, "", 0),
	}
	SetErrorHandler(secondary)

	causeErr()
	s.Assert().Len(s.errLogger.Got(), 1, "original Handler used after delegation")
	s.Assert().Len(newErrLogger.Got(), 1, "new Handler not used after delegation")
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
