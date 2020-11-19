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
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcConnection struct {
	// Ensure pointer is 64-bit aligned for atomic operations on both 32 and 64 bit machines.
	lastConnectErrPtr unsafe.Pointer

	// mu protects the connection as it is accessed by the
	// exporter goroutines and background connection goroutine
	mu sync.Mutex
	cc *grpc.ClientConn

	// these fields are read-only after constructor is finished
	c                    config
	metadata             metadata.MD
	newConnectionHandler func(cc *grpc.ClientConn) error

	// these channels are created once
	disconnectedCh             chan bool
	backgroundConnectionDoneCh chan struct{}
	stopCh                     chan struct{}

	// this is for tests, so they can replace the closing
	// routine without a worry of modifying some global variable
	// or changing it back to original after the test is done
	closeBackgroundConnectionDoneCh func(ch chan struct{})
}

func newGRPCConnection(c config, handler func(cc *grpc.ClientConn) error) *grpcConnection {
	conn := new(grpcConnection)
	conn.newConnectionHandler = handler
	if c.collectorAddr == "" {
		c.collectorAddr = fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorPort)
	}
	conn.c = c
	if len(conn.c.headers) > 0 {
		conn.metadata = metadata.New(conn.c.headers)
	}
	conn.closeBackgroundConnectionDoneCh = func(ch chan struct{}) {
		close(ch)
	}
	return conn
}

func (oc *grpcConnection) startConnection(ctx context.Context) {
	oc.stopCh = make(chan struct{})
	oc.disconnectedCh = make(chan bool)
	oc.backgroundConnectionDoneCh = make(chan struct{})

	if err := oc.connect(ctx); err == nil {
		oc.setStateConnected()
	} else {
		oc.setStateDisconnected(err)
	}
	go oc.indefiniteBackgroundConnection()
}

func (oc *grpcConnection) lastConnectError() error {
	errPtr := (*error)(atomic.LoadPointer(&oc.lastConnectErrPtr))
	if errPtr == nil {
		return nil
	}
	return *errPtr
}

func (oc *grpcConnection) saveLastConnectError(err error) {
	var errPtr *error
	if err != nil {
		errPtr = &err
	}
	atomic.StorePointer(&oc.lastConnectErrPtr, unsafe.Pointer(errPtr))
}

func (oc *grpcConnection) setStateDisconnected(err error) {
	oc.saveLastConnectError(err)
	select {
	case oc.disconnectedCh <- true:
	default:
	}
	_ = oc.newConnectionHandler(nil)
}

func (oc *grpcConnection) setStateConnected() {
	oc.saveLastConnectError(nil)
}

func (oc *grpcConnection) connected() bool {
	return oc.lastConnectError() == nil
}

const defaultConnReattemptPeriod = 10 * time.Second

func (oc *grpcConnection) indefiniteBackgroundConnection() {
	defer func() {
		oc.closeBackgroundConnectionDoneCh(oc.backgroundConnectionDoneCh)
	}()

	connReattemptPeriod := oc.c.reconnectionPeriod
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
		case <-oc.stopCh:
			return

		case <-oc.disconnectedCh:
			// Quickly check if we haven't stopped at the
			// same time.
			select {
			case <-oc.stopCh:
				return

			default:
			}

			// Normal scenario that we'll wait for
		}

		if err := oc.connect(context.Background()); err == nil {
			oc.setStateConnected()
		} else {
			oc.setStateDisconnected(err)
		}

		// Apply some jitter to avoid lockstep retrials of other
		// collector-exporters. Lockstep retrials could result in an
		// innocent DDOS, by clogging the machine's resources and network.
		jitter := time.Duration(rng.Int63n(maxJitterNanos))
		select {
		case <-oc.stopCh:
			return
		case <-time.After(connReattemptPeriod + jitter):
		}
	}
}

func (oc *grpcConnection) connect(ctx context.Context) error {
	cc, err := oc.dialToCollector(ctx)
	if err != nil {
		return err
	}
	oc.setConnection(cc)
	return oc.newConnectionHandler(cc)
}

// setConnection sets cc as the client connection and returns true if
// the connection state changed.
func (oc *grpcConnection) setConnection(cc *grpc.ClientConn) bool {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	// If previous clientConn is same as the current then just return.
	// This doesn't happen right now as this func is only called with new ClientConn.
	// It is more about future-proofing.
	if oc.cc == cc {
		return false
	}

	// If the previous clientConn was non-nil, close it
	if oc.cc != nil {
		_ = oc.cc.Close()
	}
	oc.cc = cc
	return true
}

func (oc *grpcConnection) dialToCollector(ctx context.Context) (*grpc.ClientConn, error) {
	addr := oc.c.collectorAddr

	dialOpts := []grpc.DialOption{}
	if oc.c.grpcServiceConfig != "" {
		dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(oc.c.grpcServiceConfig))
	}
	if oc.c.clientCredentials != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(oc.c.clientCredentials))
	} else if oc.c.canDialInsecure {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	if oc.c.compressor != "" {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor(oc.c.compressor)))
	}
	if len(oc.c.grpcDialOptions) != 0 {
		dialOpts = append(dialOpts, oc.c.grpcDialOptions...)
	}

	ctx, cancel := oc.contextWithStop(ctx)
	defer cancel()
	ctx = oc.contextWithMetadata(ctx)
	return grpc.DialContext(ctx, addr, dialOpts...)
}

func (oc *grpcConnection) contextWithMetadata(ctx context.Context) context.Context {
	if oc.metadata.Len() > 0 {
		return metadata.NewOutgoingContext(ctx, oc.metadata)
	}
	return ctx
}

func (oc *grpcConnection) shutdown(ctx context.Context) error {
	close(oc.stopCh)
	// Ensure that the backgroundConnector returns
	select {
	case <-oc.backgroundConnectionDoneCh:
	case <-ctx.Done():
		return ctx.Err()
	}

	close(oc.disconnectedCh)

	oc.mu.Lock()
	cc := oc.cc
	oc.cc = nil
	oc.mu.Unlock()

	if cc != nil {
		return cc.Close()
	}

	return nil
}

func (oc *grpcConnection) contextWithStop(ctx context.Context) (context.Context, context.CancelFunc) {
	// Unify the parent context Done signal with the connection's
	// stop channel.
	ctx, cancel := context.WithCancel(ctx)
	go func(ctx context.Context, cancel context.CancelFunc) {
		select {
		case <-ctx.Done():
			// Nothing to do, either cancelled or deadline
			// happened.
		case <-oc.stopCh:
			cancel()
		}
	}(ctx, cancel)
	return ctx, cancel
}
