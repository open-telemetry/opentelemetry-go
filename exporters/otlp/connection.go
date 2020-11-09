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

package otlp // import "go.opentelemetry.io/otel/exporters/otlp"

import (
	"math/rand"
	"sync/atomic"
	"time"
	"unsafe"
)

func (e *Exporter) lastConnectError() error {
	errPtr := (*error)(atomic.LoadPointer(&e.lastConnectErrPtr))
	if errPtr == nil {
		return nil
	}
	return *errPtr
}

func (e *Exporter) saveLastConnectError(err error) {
	var errPtr *error
	if err != nil {
		errPtr = &err
	}
	atomic.StorePointer(&e.lastConnectErrPtr, unsafe.Pointer(errPtr))
}

func (e *Exporter) setStateDisconnected(err error) {
	e.saveLastConnectError(err)
	select {
	case e.disconnectedCh <- true:
	default:
	}
}

func (e *Exporter) setStateConnected() {
	e.saveLastConnectError(nil)
}

func (e *Exporter) connected() bool {
	return e.lastConnectError() == nil
}

const defaultConnReattemptPeriod = 10 * time.Second

func (e *Exporter) indefiniteBackgroundConnection() {
	defer func() {
		e.backgroundConnectionDoneCh <- true
	}()

	connReattemptPeriod := e.c.reconnectionPeriod
	if connReattemptPeriod <= 0 {
		connReattemptPeriod = defaultConnReattemptPeriod
	}

	// No strong seeding required, nano time can
	// already help with pseudo uniqueness.
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63n(1024)))

	// maxJitterNanos: 70% of the connectionReattemptPeriod
	maxJitterNanos := int64(0.7 * float64(connReattemptPeriod))

	for {
		// Otherwise these will be the normal scenarios to enable
		// reconnection if we trip out.
		// 1. If we've stopped, return entirely
		// 2. Otherwise block until we are disconnected, and
		//    then retry connecting
		select {
		case <-e.stopCh:
			return

		case <-e.disconnectedCh:
			// Normal scenario that we'll wait for
		}

		if err := e.connect(); err == nil {
			e.setStateConnected()
		} else {
			e.setStateDisconnected(err)
		}

		// Apply some jitter to avoid lockstep retrials of other
		// collector-exporters. Lockstep retrials could result in an
		// innocent DDOS, by clogging the machine's resources and network.
		jitter := time.Duration(rng.Int63n(maxJitterNanos))
		select {
		case <-e.stopCh:
			return
		case <-time.After(connReattemptPeriod + jitter):
		}
	}
}

func (e *Exporter) connect() error {
	cc, err := e.dialToCollector()
	if err != nil {
		return err
	}
	return e.enableConnections(cc)
}
