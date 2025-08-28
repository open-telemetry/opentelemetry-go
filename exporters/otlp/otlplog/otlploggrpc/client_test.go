// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	lpb "go.opentelemetry.io/proto/otlp/logs/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var (
	// Sat Jan 01 2000 00:00:00 GMT+0000.
	ts  = time.Date(2000, time.January, 0o1, 0, 0, 0, 0, time.FixedZone("GMT", 0))
	obs = ts.Add(30 * time.Second)

	kvAlice = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "alice"},
	}}
	kvBob = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "bob"},
	}}
	kvSrvName = &cpb.KeyValue{Key: "service.name", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "test server"},
	}}
	kvSrvVer = &cpb.KeyValue{Key: "service.version", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "v0.1.0"},
	}}

	pbSevA = lpb.SeverityNumber_SEVERITY_NUMBER_INFO
	pbSevB = lpb.SeverityNumber_SEVERITY_NUMBER_ERROR

	pbBodyA = &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{
			StringValue: "a",
		},
	}
	pbBodyB = &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{
			StringValue: "b",
		},
	}

	spanIDA  = []byte{0, 0, 0, 0, 0, 0, 0, 1}
	spanIDB  = []byte{0, 0, 0, 0, 0, 0, 0, 2}
	traceIDA = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	traceIDB = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	flagsA   = byte(1)
	flagsB   = byte(0)

	logRecords = []*lpb.LogRecord{
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevA,
			SeverityText:         "A",
			Body:                 pbBodyA,
			Attributes:           []*cpb.KeyValue{kvAlice},
			Flags:                uint32(flagsA),
			TraceId:              traceIDA,
			SpanId:               spanIDA,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevA,
			SeverityText:         "A",
			Body:                 pbBodyA,
			Attributes:           []*cpb.KeyValue{kvBob},
			Flags:                uint32(flagsA),
			TraceId:              traceIDA,
			SpanId:               spanIDA,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevB,
			SeverityText:         "B",
			Body:                 pbBodyB,
			Attributes:           []*cpb.KeyValue{kvAlice},
			Flags:                uint32(flagsB),
			TraceId:              traceIDB,
			SpanId:               spanIDB,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevB,
			SeverityText:         "B",
			Body:                 pbBodyB,
			Attributes:           []*cpb.KeyValue{kvBob},
			Flags:                uint32(flagsB),
			TraceId:              traceIDB,
			SpanId:               spanIDB,
		},
	}

	scope = &cpb.InstrumentationScope{
		Name:    "test/code/path",
		Version: "v0.1.0",
	}
	scopeLogs = []*lpb.ScopeLogs{
		{
			Scope:      scope,
			LogRecords: logRecords,
			SchemaUrl:  semconv.SchemaURL,
		},
	}

	res = &rpb.Resource{
		Attributes: []*cpb.KeyValue{kvSrvName, kvSrvVer},
	}
	resourceLogs = []*lpb.ResourceLogs{{
		Resource:  res,
		ScopeLogs: scopeLogs,
		SchemaUrl: semconv.SchemaURL,
	}}
)

func TestThrottleDelay(t *testing.T) {
	c := codes.ResourceExhausted
	testcases := []struct {
		status       *status.Status
		wantOK       bool
		wantDuration time.Duration
	}{
		{
			status:       status.New(c, "NoRetryInfo"),
			wantOK:       false,
			wantDuration: 0,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "SingleRetryInfo").WithDetails(
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(15 * time.Millisecond),
					},
				)
				require.NoError(t, err)
				return s
			}(),
			wantOK:       true,
			wantDuration: 15 * time.Millisecond,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "ErrorInfo").WithDetails(
					&errdetails.ErrorInfo{Reason: "no throttle detail"},
				)
				require.NoError(t, err)
				return s
			}(),
			wantOK:       false,
			wantDuration: 0,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "ErrorAndRetryInfo").WithDetails(
					&errdetails.ErrorInfo{Reason: "with throttle detail"},
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(13 * time.Minute),
					},
				)
				require.NoError(t, err)
				return s
			}(),
			wantOK:       true,
			wantDuration: 13 * time.Minute,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "DoubleRetryInfo").WithDetails(
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(13 * time.Minute),
					},
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(15 * time.Minute),
					},
				)
				require.NoError(t, err)
				return s
			}(),
			wantOK:       true,
			wantDuration: 13 * time.Minute,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.status.Message(), func(t *testing.T) {
			ok, d := throttleDelay(tc.status)
			assert.Equal(t, tc.wantOK, ok)
			assert.Equal(t, tc.wantDuration, d)
		})
	}
}

func TestRetryable(t *testing.T) {
	retryableCodes := map[codes.Code]bool{
		codes.OK:                 false,
		codes.Canceled:           true,
		codes.Unknown:            false,
		codes.InvalidArgument:    false,
		codes.DeadlineExceeded:   true,
		codes.NotFound:           false,
		codes.AlreadyExists:      false,
		codes.PermissionDenied:   false,
		codes.ResourceExhausted:  false,
		codes.FailedPrecondition: false,
		codes.Aborted:            true,
		codes.OutOfRange:         true,
		codes.Unimplemented:      false,
		codes.Internal:           false,
		codes.Unavailable:        true,
		codes.DataLoss:           true,
		codes.Unauthenticated:    false,
	}

	for c, want := range retryableCodes {
		got, _ := retryable(status.Error(c, ""))
		assert.Equalf(t, want, got, "evaluate(%s)", c)
	}
}

func TestRetryableGRPCStatusResourceExhaustedWithRetryInfo(t *testing.T) {
	delay := 15 * time.Millisecond
	s, err := status.New(codes.ResourceExhausted, "WithRetryInfo").WithDetails(
		&errdetails.RetryInfo{
			RetryDelay: durationpb.New(delay),
		},
	)
	require.NoError(t, err)

	ok, d := retryableGRPCStatus(s)
	assert.True(t, ok)
	assert.Equal(t, delay, d)
}

func TestNewClient(t *testing.T) {
	newGRPCClientFnSwap := newGRPCClientFn
	t.Cleanup(func() {
		newGRPCClientFn = newGRPCClientFnSwap
	})

	// The gRPC connection created by newClient.
	conn, err := grpc.NewClient("test", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	newGRPCClientFn = func(string, ...grpc.DialOption) (*grpc.ClientConn, error) {
		return conn, nil
	}

	// The gRPC connection created by users.
	userConn, err := grpc.NewClient("test 2", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	testCases := []struct {
		name string
		cfg  config
		cli  *client
	}{
		{
			name: "empty config",
			cli: &client{
				ourConn: true,
				conn:    conn,
				lsc:     collogpb.NewLogsServiceClient(conn),
			},
		},
		{
			name: "with headers",
			cfg: config{
				headers: newSetting(map[string]string{
					"key": "value",
				}),
			},
			cli: &client{
				ourConn:  true,
				conn:     conn,
				lsc:      collogpb.NewLogsServiceClient(conn),
				metadata: map[string][]string{"key": {"value"}},
			},
		},
		{
			name: "with gRPC connection",
			cfg: config{
				gRPCConn: newSetting(userConn),
			},
			cli: &client{
				ourConn: false,
				conn:    userConn,
				lsc:     collogpb.NewLogsServiceClient(userConn),
			},
		},
		{
			// It is not possible to compare grpc dial options directly, so we just check that the client is created
			// and no panic occurs.
			name: "with dial options",
			cfg: config{
				serviceConfig:      newSetting("service config"),
				gRPCCredentials:    newSetting(credentials.NewTLS(nil)),
				compression:        newSetting(GzipCompression),
				reconnectionPeriod: newSetting(10 * time.Second),
			},
			cli: &client{
				ourConn: true,
				conn:    conn,
				lsc:     collogpb.NewLogsServiceClient(conn),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cli, err := newClient(tc.cfg)
			require.NoError(t, err)

			assert.Equal(t, tc.cli.metadata, cli.metadata)
			assert.Equal(t, tc.cli.exportTimeout, cli.exportTimeout)
			assert.Equal(t, tc.cli.ourConn, cli.ourConn)
			assert.Equal(t, tc.cli.conn, cli.conn)
			assert.Equal(t, tc.cli.lsc, cli.lsc)
		})
	}
}

type exportResult struct {
	Response *collogpb.ExportLogsServiceResponse
	Err      error
}

// storage stores uploaded OTLP log data in their proto form.
type storage struct {
	dataMu sync.Mutex
	data   []*lpb.ResourceLogs
}

// newStorage returns a configure storage ready to store received requests.
func newStorage() *storage {
	return &storage{}
}

// Add adds the request to the Storage.
func (s *storage) Add(request *collogpb.ExportLogsServiceRequest) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()
	s.data = append(s.data, request.ResourceLogs...)
}

// Dump returns all added ResourceLogs and clears the storage.
func (s *storage) Dump() []*lpb.ResourceLogs {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var data []*lpb.ResourceLogs
	data, s.data = s.data, []*lpb.ResourceLogs{}
	return data
}

// grpcCollector is an OTLP gRPC server that collects all requests it receives.
type grpcCollector struct {
	collogpb.UnimplementedLogsServiceServer

	headersMu sync.Mutex
	headers   metadata.MD
	storage   *storage

	resultCh <-chan exportResult
	listener net.Listener
	srv      *grpc.Server
}

var _ collogpb.LogsServiceServer = (*grpcCollector)(nil)

// newGRPCCollector returns a *grpcCollector that is listening at the provided
// endpoint.
//
// If endpoint is an empty string, the returned collector will be listening on
// the localhost interface at an OS chosen port.
//
// If errCh is not nil, the collector will respond to Export calls with errors
// sent on that channel. This means that if errCh is not nil Export calls will
// block until an error is received.
func newGRPCCollector(endpoint string, resultCh <-chan exportResult) (*grpcCollector, error) {
	if endpoint == "" {
		endpoint = "localhost:0"
	}

	c := &grpcCollector{
		storage:  newStorage(),
		resultCh: resultCh,
	}

	var err error
	c.listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	c.srv = grpc.NewServer()
	collogpb.RegisterLogsServiceServer(c.srv, c)
	go func() { _ = c.srv.Serve(c.listener) }()

	return c, nil
}

// Export handles the export req.
func (c *grpcCollector) Export(
	ctx context.Context,
	req *collogpb.ExportLogsServiceRequest,
) (*collogpb.ExportLogsServiceResponse, error) {
	c.storage.Add(req)

	if h, ok := metadata.FromIncomingContext(ctx); ok {
		c.headersMu.Lock()
		c.headers = metadata.Join(c.headers, h)
		c.headersMu.Unlock()
	}

	if c.resultCh != nil {
		r := <-c.resultCh
		if r.Response == nil {
			return &collogpb.ExportLogsServiceResponse{}, r.Err
		}
		return r.Response, r.Err
	}
	return &collogpb.ExportLogsServiceResponse{}, nil
}

// Collect returns the Storage holding all collected requests.
func (c *grpcCollector) Collect() *storage {
	return c.storage
}

func clientFactory(t *testing.T, rCh <-chan exportResult) (*client, *grpcCollector) {
	t.Helper()
	coll, err := newGRPCCollector("", rCh)
	require.NoError(t, err)

	addr := coll.listener.Addr().String()
	opts := []Option{WithEndpoint(addr), WithInsecure()}
	cfg := newConfig(opts)
	client, err := newClient(cfg)
	require.NoError(t, err)
	return client, coll
}

func testCtxErrs(factory func() func(context.Context) error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		t.Run("DeadlineExceeded", func(t *testing.T) {
			innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
			t.Cleanup(innerCancel)
			<-innerCtx.Done()

			f := factory()
			assert.ErrorIs(t, f(innerCtx), context.DeadlineExceeded)
		})

		t.Run("Canceled", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			f := factory()
			assert.ErrorIs(t, f(innerCtx), context.Canceled)
		})
	}
}

func TestClient(t *testing.T) {
	t.Run("ClientHonorsContextErrors", func(t *testing.T) {
		t.Run("Shutdown", testCtxErrs(func() func(context.Context) error {
			c, _ := clientFactory(t, nil)
			return c.Shutdown
		}))

		t.Run("UploadLog", testCtxErrs(func() func(context.Context) error {
			c, _ := clientFactory(t, nil)
			return func(ctx context.Context) error {
				return c.UploadLogs(ctx, nil)
			}
		}))
	})

	t.Run("UploadLogs", func(t *testing.T) {
		ctx := context.Background()
		client, coll := clientFactory(t, nil)

		require.NoError(t, client.UploadLogs(ctx, resourceLogs))
		require.NoError(t, client.Shutdown(ctx))
		got := coll.Collect().Dump()
		require.Len(t, got, 1, "upload of one ResourceLogs")
		diff := cmp.Diff(got[0], resourceLogs[0], cmp.Comparer(proto.Equal))
		if diff != "" {
			t.Fatalf("unexpected ResourceLogs:\n%s", diff)
		}
	})

	t.Run("PartialSuccess", func(t *testing.T) {
		const n, msg = 2, "bad data"
		rCh := make(chan exportResult, 3)
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{
				PartialSuccess: &collogpb.ExportLogsPartialSuccess{
					RejectedLogRecords: n,
					ErrorMessage:       msg,
				},
			},
		}
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{
				PartialSuccess: &collogpb.ExportLogsPartialSuccess{
					// Should not be logged.
					RejectedLogRecords: 0,
					ErrorMessage:       "",
				},
			},
		}
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{},
		}

		ctx := context.Background()
		client, _ := clientFactory(t, rCh)

		defer func(orig otel.ErrorHandler) {
			otel.SetErrorHandler(orig)
		}(otel.GetErrorHandler())

		var errs []error
		eh := otel.ErrorHandlerFunc(func(e error) { errs = append(errs, e) })
		otel.SetErrorHandler(eh)

		require.NoError(t, client.UploadLogs(ctx, resourceLogs))
		require.NoError(t, client.UploadLogs(ctx, resourceLogs))
		require.NoError(t, client.UploadLogs(ctx, resourceLogs))

		require.Len(t, errs, 1)
		want := fmt.Sprintf("%s (%d log records rejected)", msg, n)
		assert.ErrorContains(t, errs[0], want)
	})
}

func TestConfig(t *testing.T) {
	factoryFunc := func(rCh <-chan exportResult, o ...Option) (log.Exporter, *grpcCollector) {
		coll, err := newGRPCCollector("", rCh)
		require.NoError(t, err)

		ctx := context.Background()
		opts := append([]Option{
			WithEndpoint(coll.listener.Addr().String()),
			WithInsecure(),
		}, o...)
		exp, err := New(ctx, opts...)
		require.NoError(t, err)
		return exp, coll
	}

	t.Run("WithHeaders", func(t *testing.T) {
		key := "my-custom-header"
		headers := map[string]string{key: "custom-value"}
		exp, coll := factoryFunc(nil, WithHeaders(headers))
		t.Cleanup(coll.srv.Stop)

		ctx := context.Background()
		additionalKey := "additional-custom-header"
		ctx = metadata.AppendToOutgoingContext(ctx, additionalKey, "additional-value")
		require.NoError(t, exp.Export(ctx, make([]log.Record, 1)))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := metadata.Join(coll.headers)
		require.Regexp(t, "OTel Go OTLP over gRPC logs exporter/[01]\\..*", got)
		require.Contains(t, got, key)
		require.Contains(t, got, additionalKey)
		assert.Equal(t, []string{headers[key]}, got[key])
	})
}
