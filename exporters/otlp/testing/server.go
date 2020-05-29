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

package otlp_testing

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/exporters/otlp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Let the system define the port
const defaultServerAddr = "127.0.0.1:0"

type ServerSuite struct {
	suite.Suite

	ServerOpts     []grpc.ServerOption
	serverAddr     string
	ServerListener net.Listener
	Server         *grpc.Server

	requestCount  uint64
	FailureModulo uint64
	FailureCodes  []codes.Code

	ExporterOpts []otlp.ExporterOption
	Exporter     *otlp.Exporter
}

func (s *ServerSuite) SetupSuite() {
	s.serverAddr = defaultServerAddr

	var err error
	s.ServerListener, err = net.Listen("tcp", s.serverAddr)
	s.serverAddr = s.ServerListener.Addr().String()
	require.NoError(s.T(), err, "failed to allocate a port for server")

	s.Server = grpc.NewServer(s.ServerOpts...)
}

func (s *ServerSuite) StartServer() {
	go func() {
		s.Server.Serve(s.ServerListener)
	}()
	s.T().Logf("started grpc.Server at: %v", s.ServerAddr())

	if s.Exporter == nil {
		s.Exporter = s.NewExporter()
	}
}

func (s *ServerSuite) ServerAddr() string {
	return s.serverAddr
}

func (s *ServerSuite) RequestError() error {
	if s.FailureModulo <= 0 {
		return nil
	}

	count := atomic.AddUint64(&s.requestCount, 1)
	if count%s.FailureModulo == 0 {
		return nil
	}

	var c codes.Code
	if n := len(s.FailureCodes); n > 0 {
		/* Example to understand the indexing:
		*  - s.FailureModulo = 3
		*  - len(s.Codes) 5
		*
		* count - 1 | count / s.FailureModulo | index (mod 5)
		* ===================================================
		*      0    |            0            |       0
		*      1    |            0            |       1
		*      2    |            1            |     n/a (2)
		*      3    |            1            |       2
		*      4    |            1            |       3
		*      5    |            2            |     n/a (3)
		*      6    |            2            |       4
		*      7    |            2            |       0
		*      8    |            3            |     n/a (0)
		*      9    |            3            |       1
		 */
		c = s.FailureCodes[(count-1-(count/s.FailureModulo))%uint64(n)]
	} else {
		c = codes.Unavailable
	}
	return status.Errorf(c, "artificial error: count %d, modulo %d", count, s.FailureModulo)
}

func (s *ServerSuite) NewExporter() *otlp.Exporter {
	opts := []otlp.ExporterOption{
		otlp.WithInsecure(),
		otlp.WithAddress(s.serverAddr),
		otlp.WithReconnectionPeriod(10 * time.Millisecond),
	}
	exp, err := otlp.NewExporter(append(opts, s.ExporterOpts...)...)
	require.NoError(s.T(), err, "failed to create exporter")
	return exp
}

func (s *ServerSuite) TearDownSuite() {
	if s.ServerListener != nil {
		s.Server.GracefulStop()
		s.T().Logf("stopped grpc.Server at: %v", s.ServerAddr())
	}
	s.Exporter.Stop()
}
