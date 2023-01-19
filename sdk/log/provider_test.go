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

package log

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

type storingHandler struct {
	errs []error
}

func (s *storingHandler) Handle(err error) {
	s.errs = append(s.errs, err)
}

func (s *storingHandler) Reset() {
	s.errs = nil
}

var (
	handler = &storingHandler{}
)

type basicLogRecordProcessor struct {
	flushed             bool
	closed              bool
	injectShutdownError error
}

func (t *basicLogRecordProcessor) Shutdown(context.Context) error {
	t.closed = true
	return t.injectShutdownError
}

func (t *basicLogRecordProcessor) OnEmit(context.Context, ReadWriteLogRecord) {}
func (t *basicLogRecordProcessor) ForceFlush(context.Context) error {
	t.flushed = true
	return nil
}

func TestForceFlushAndShutdownTraceProviderWithoutProcessor(t *testing.T) {
	stp := NewLoggerProvider()
	assert.NoError(t, stp.ForceFlush(context.Background()))
	assert.NoError(t, stp.Shutdown(context.Background()))
}

func TestShutdownTraceProvider(t *testing.T) {
	stp := NewLoggerProvider()
	sp := &basicLogRecordProcessor{}
	stp.RegisterLogRecordProcessor(sp)

	assert.NoError(t, stp.ForceFlush(context.Background()))
	assert.True(t, sp.flushed, "error ForceFlush basicLogRecordProcessor")
	assert.NoError(t, stp.Shutdown(context.Background()))
	assert.True(t, sp.closed, "error Shutdown basicLogRecordProcessor")
}

func TestFailedProcessorShutdown(t *testing.T) {
	stp := NewLoggerProvider()
	spErr := errors.New("basic span processor shutdown failure")
	sp := &basicLogRecordProcessor{
		injectShutdownError: spErr,
	}
	stp.RegisterLogRecordProcessor(sp)

	err := stp.Shutdown(context.Background())
	assert.Error(t, err)
	assert.Equal(t, err, spErr)
}

func TestFailedProcessorsShutdown(t *testing.T) {
	stp := NewLoggerProvider()
	spErr1 := errors.New("basic span processor shutdown failure1")
	spErr2 := errors.New("basic span processor shutdown failure2")
	sp1 := &basicLogRecordProcessor{
		injectShutdownError: spErr1,
	}
	sp2 := &basicLogRecordProcessor{
		injectShutdownError: spErr2,
	}
	stp.RegisterLogRecordProcessor(sp1)
	stp.RegisterLogRecordProcessor(sp2)

	err := stp.Shutdown(context.Background())
	assert.Error(t, err)
	assert.EqualError(t, err, "basic span processor shutdown failure1; basic span processor shutdown failure2")
	assert.True(t, sp1.closed)
	assert.True(t, sp2.closed)
}

func TestSchemaURL(t *testing.T) {
	stp := NewLoggerProvider()
	schemaURL := "https://opentelemetry.io/schemas/1.2.0"
	tracerIface := stp.Logger("tracername", log.WithSchemaURL(schemaURL))

	// Verify that the SchemaURL of the constructed Tracer is correctly populated.
	tracerStruct := tracerIface.(*logger)
	assert.EqualValues(t, schemaURL, tracerStruct.instrumentationScope.SchemaURL)
}
