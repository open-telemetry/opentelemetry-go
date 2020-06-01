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

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/example/grpc/api"
	"go.opentelemetry.io/otel/example/grpc/config"
	"go.opentelemetry.io/otel/instrumentation/grpctrace"

	"google.golang.org/grpc"
)

const (
	port = ":7777"
)

// server is used to implement api.HelloServiceServer
type server struct {
	api.HelloServiceServer
}

// SayHello implements api.HelloServiceServer
func (s *server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	log.Printf("Received: %v\n", in.GetGreeting())
	time.Sleep(50 * time.Millisecond)

	return &api.HelloResponse{Reply: "Hello " + in.Greeting}, nil
}

func (s *server) SayHelloServerStream(in *api.HelloRequest, out api.HelloService_SayHelloServerStreamServer) error {
	log.Printf("Received: %v\n", in.GetGreeting())

	for i := 0; i < 5; i++ {
		err := out.Send(&api.HelloResponse{Reply: "Hello " + in.Greeting})
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(i*50) * time.Millisecond)
	}

	return nil
}

func (s *server) SayHelloClientStream(stream api.HelloService_SayHelloClientStreamServer) error {
	i := 0

	for {
		in, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Non EOF error: %v\n", err)
			return err
		}

		log.Printf("Received: %v\n", in.GetGreeting())
		i++
	}

	time.Sleep(50 * time.Millisecond)

	return stream.SendAndClose(&api.HelloResponse{Reply: fmt.Sprintf("Hello (%v times)", i)})
}

func (s *server) SayHelloBidiStream(stream api.HelloService_SayHelloBidiStreamServer) error {
	for {
		in, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Non EOF error: %v\n", err)
			return err
		}

		time.Sleep(50 * time.Millisecond)

		log.Printf("Received: %v\n", in.GetGreeting())
		err = stream.Send(&api.HelloResponse{Reply: "Hello " + in.Greeting})

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	config.Init()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctrace.UnaryServerInterceptor(global.Tracer(""))),
		grpc.StreamInterceptor(grpctrace.StreamServerInterceptor(global.Tracer(""))),
	)

	api.RegisterHelloServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
