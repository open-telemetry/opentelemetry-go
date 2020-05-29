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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel/api/oterror"
)

type mock []error

func (m *mock) Handle(err error) {
	(*m) = append(*m, err)
}

type HandlerTestSuite struct {
	suite.Suite

	origHandler oterror.Handler

	errs []error
}

func (s *HandlerTestSuite) Handle(err error) {
	s.errs = append(s.errs, err)
}

func (s *HandlerTestSuite) SetupSuite() {
	s.origHandler = globalHandler
	globalHandler = s
}

func (s *HandlerTestSuite) TearDownSuite() {
	globalHandler = s.origHandler
}

func (s *HandlerTestSuite) SetupTest() {
	s.errs = []error{}
}

func (s *HandlerTestSuite) TestGlocalHandler() {
	err1 := errors.New("one")
	err2 := errors.New("two")
	Handler().Handle(err1)
	Handle(err2)
	s.Assert().Equal([]error{err1, err2}, s.errs)
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
	secondary := new(mock)
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

	s.Assert().Greater(len(s.errs), 1, "at least 2 errors should have been sent")
	s.Assert().Len(s.errs, sent)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
