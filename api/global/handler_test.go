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
	"log"
	"testing"
	"time"

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

	origHandler *handler

	errLogger *errLogger
}

func (s *HandlerTestSuite) SetupSuite() {
	s.errLogger = new(errLogger)
	s.origHandler = defaultHandler
	defaultHandler = &handler{
		l: log.New(s.errLogger, "", 0),
	}
}

func (s *HandlerTestSuite) TearDownSuite() {
	defaultHandler = s.origHandler
}

func (s *HandlerTestSuite) SetupTest() {
	s.errLogger.Reset()
}

func (s *HandlerTestSuite) TestGlobalHandler() {
	errs := []string{"one", "two"}
	Handler().Handle(errors.New(errs[0]))
	Handle(errors.New(errs[1]))
	s.Assert().Equal(errs, s.errLogger.Got())
}

func (s *HandlerTestSuite) TestNoDropsOnDelegate() {
	var sent int
	err := errors.New("")
	stop := make(chan struct{})
	beat := make(chan struct{})
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				done <- struct{}{}
				return
			default:
				sent++
				Handle(err)
			}

			select {
			case beat <- struct{}{}:
			default:
			}
		}
	}()

	// Wait for the spice to flow
	select {
	case <-time.Tick(2 * time.Millisecond):
		s.T().Fatal("no errors were sent in 2ms")
	case <-beat:
	}

	// Change to another Handler. We are testing this is loss-less.
	newErrLogger := new(errLogger)
	secondary := &handler{
		l: log.New(newErrLogger, "", 0),
	}
	SetHandler(secondary)

	select {
	case <-time.Tick(2 * time.Millisecond):
		s.T().Fatal("no errors were sent within 2ms after SetHandler")
	case <-beat:
	}

	// Now beat is clear, wait for a fresh send.
	select {
	case <-time.Tick(2 * time.Millisecond):
		s.T().Fatal("no fresh errors were sent within 2ms after SetHandler")
	case <-beat:
	}

	// Stop sending errors.
	stop <- struct{}{}
	// Ensure we do not lose any straglers.
	<-done

	got := append(s.errLogger.Got(), newErrLogger.Got()...)
	s.Assert().Greater(len(got), 1, "at least 2 errors should have been sent")
	s.Assert().Len(got, sent)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
