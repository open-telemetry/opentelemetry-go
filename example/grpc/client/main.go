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
	"io"
	"log"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/example/grpc/api"
	"go.opentelemetry.io/otel/example/grpc/config"
	"go.opentelemetry.io/otel/instrumentation/grpctrace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	config.Init()

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7777", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpctrace.UnaryClientInterceptor(global.Tracer(""))),
		grpc.WithStreamInterceptor(grpctrace.StreamClientInterceptor(global.Tracer(""))),
	)

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer func() { _ = conn.Close() }()

	c := api.NewHelloServiceClient(conn)

	callSayHello(c)
	callSayHelloClientStream(c)
	callSayHelloServerStream(c)
	callSayHelloBidiStream(c)

	time.Sleep(10 * time.Millisecond)
}

func callSayHello(c api.HelloServiceClient) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	response, err := c.SayHello(ctx, &api.HelloRequest{Greeting: "World"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Reply)
}

func callSayHelloClientStream(c api.HelloServiceClient) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.SayHelloClientStream(ctx)
	if err != nil {
		log.Fatalf("Error when opening SayHelloClientStream: %s", err)
	}

	for i := 0; i < 5; i++ {
		err := stream.Send(&api.HelloRequest{Greeting: "World"})

		time.Sleep(time.Duration(i*50) * time.Millisecond)

		if err != nil {
			log.Fatalf("Error when sending to SayHelloClientStream: %s", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error when closing SayHelloClientStream: %s", err)
	}

	log.Printf("Response from server: %s", response.Reply)
}

func callSayHelloServerStream(c api.HelloServiceClient) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.SayHelloServerStream(ctx, &api.HelloRequest{Greeting: "World"})
	if err != nil {
		log.Fatalf("Error when opening SayHelloServerStream: %s", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error when receiving from SayHelloServerStream: %s", err)
		}

		log.Printf("Response from server: %s", response.Reply)
		time.Sleep(50 * time.Millisecond)
	}
}

func callSayHelloBidiStream(c api.HelloServiceClient) {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.SayHelloBidiStream(ctx)
	if err != nil {
		log.Fatalf("Error when opening SayHelloBidiStream: %s", err)
	}

	serverClosed := make(chan struct{})
	clientClosed := make(chan struct{})

	go func() {
		for i := 0; i < 5; i++ {
			err := stream.Send(&api.HelloRequest{Greeting: "World"})

			if err != nil {
				log.Fatalf("Error when sending to SayHelloBidiStream: %s", err)
			}

			time.Sleep(50 * time.Millisecond)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Error when closing SayHelloBidiStream: %s", err)
		}

		clientClosed <- struct{}{}
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error when receiving from SayHelloBidiStream: %s", err)
			}

			log.Printf("Response from server: %s", response.Reply)
			time.Sleep(50 * time.Millisecond)
		}

		serverClosed <- struct{}{}
	}()

	// Wait until client and server both closed the connection.
	<-clientClosed
	<-serverClosed
}
