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
	"fmt"
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
	// max time to wait for goroutine to Handle an error.
	pause := 10 * time.Millisecond

	var sent int
	err := errors.New("")
	stop := make(chan struct{})
	beat := make(chan struct{})
	done := make(chan struct{})

	// Wait for a error to be submitted from the following goroutine.
	wait := func(d time.Duration) error {
		timer := time.NewTimer(d)
		select {
		case <-timer.C:
			// We are about to fail, stop the spawned goroutine.
			stop <- struct{}{}
			return fmt.Errorf("no errors sent in %v", d)
		case <-beat:
			// Allow the timer to be reclaimed by GC.
			timer.Stop()
			return nil
		}
	}

	go func() {
		// Slow down to speed up: do not overload the processor.
		ticker := time.NewTicker(100 * time.Microsecond)
		for {
			select {
			case <-stop:
				ticker.Stop()
				done <- struct{}{}
				return
			case <-ticker.C:
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
	s.Require().NoError(wait(pause), "starting error stream")

	// Change to another Handler. We are testing this is loss-less.
	newErrLogger := new(errLogger)
	secondary := &loggingErrorHandler{
		l: log.New(newErrLogger, "", 0),
	}
	SetErrorHandler(secondary)
	s.Require().NoError(wait(pause), "switched to new Handler")

	// Testing done, stop sending errors.
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
