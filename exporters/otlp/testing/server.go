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
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/exporters/otlp"
	"google.golang.org/grpc"
)

// Let the system define the port
const defaultServerAddr = "127.0.0.1:0"

type ServerSuite struct {
	suite.Suite

	ServerOpts     []grpc.ServerOption
	serverAddr     string
	ServerListener net.Listener
	Server         *grpc.Server

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
