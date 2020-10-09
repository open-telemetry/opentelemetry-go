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

package zipkin

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	zkmodel "github.com/openzipkin/zipkin-go/model"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

func TestModelConversion(t *testing.T) {
	inputBatch := []*export.SpanData{
		// typical span data
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			SpanKind:     otel.SpanKindServer,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.Uint64("attr1", 42),
				label.String("attr2", "bar"),
			},
			MessageEvents: []export.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []label.KeyValue{
						label.Uint64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// span data with no parent (same as typical, but has
		// invalid parent)
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{},
			SpanKind:     otel.SpanKindServer,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.Uint64("attr1", 42),
				label.String("attr2", "bar"),
			},
			MessageEvents: []export.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []label.KeyValue{
						label.Uint64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// span data of unspecified kind
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			SpanKind:     otel.SpanKindUnspecified,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.Uint64("attr1", 42),
				label.String("attr2", "bar"),
			},
			MessageEvents: []export.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []label.KeyValue{
						label.Uint64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// span data of internal kind
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			SpanKind:     otel.SpanKindInternal,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.Uint64("attr1", 42),
				label.String("attr2", "bar"),
			},
			MessageEvents: []export.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []label.KeyValue{
						label.Uint64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// span data of client kind
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			SpanKind:     otel.SpanKindClient,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.Uint64("attr1", 42),
				label.String("attr2", "bar"),
			},
			MessageEvents: []export.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []label.KeyValue{
						label.Uint64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// span data of producer kind
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			SpanKind:     otel.SpanKindProducer,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.Uint64("attr1", 42),
				label.String("attr2", "bar"),
			},
			MessageEvents: []export.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []label.KeyValue{
						label.Uint64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// span data of consumer kind
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			SpanKind:     otel.SpanKindConsumer,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.Uint64("attr1", 42),
				label.String("attr2", "bar"),
			},
			MessageEvents: []export.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []label.KeyValue{
						label.Uint64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// span data with no events
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			SpanKind:     otel.SpanKindServer,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.Uint64("attr1", 42),
				label.String("attr2", "bar"),
			},
			MessageEvents: nil,
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// span data with an "error" attribute set to "false"
		{
			SpanContext: otel.SpanContext{
				TraceID: otel.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  otel.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID: otel.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			SpanKind:     otel.SpanKindServer,
			Name:         "foo",
			StartTime:    time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:      time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []label.KeyValue{
				label.String("error", "false"),
			},
			MessageEvents: []export.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []label.KeyValue{
						label.Uint64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
	}

	expectedOutputBatch := []zkmodel.SpanModel{
		// model for typical span data
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":                   "42",
				"attr2":                   "bar",
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model for span data with no parent
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: nil,
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":                   "42",
				"attr2":                   "bar",
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model for span data of unspecified kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":                   "42",
				"attr2":                   "bar",
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model for span data of internal kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":                   "42",
				"attr2":                   "bar",
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model for span data of client kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "CLIENT",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":                   "42",
				"attr2":                   "bar",
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model for span data of producer kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "PRODUCER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":                   "42",
				"attr2":                   "bar",
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model for span data of consumer kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "CONSUMER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":                   "42",
				"attr2":                   "bar",
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model for span data with no events
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations:    nil,
			Tags: map[string]string{
				"attr1":                   "42",
				"attr2":                   "bar",
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model for span data with an "error" attribute set to "false"
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
	}
	gottenOutputBatch := toZipkinSpanModels(inputBatch, "model-test")
	require.Equal(t, expectedOutputBatch, gottenOutputBatch)
}

func zkmodelIDPtr(n uint64) *zkmodel.ID {
	id := zkmodel.ID(n)
	return &id
}

func Test_toZipkinTags(t *testing.T) {
	keyValue := "value"
	doubleValue := 123.456
	uintValue := int64(123)
	statusMessage := "this is a problem"
	instrLibName := "instrumentation-library"
	instrLibVersion := "semver:1.0.0"

	tests := []struct {
		name string
		data *export.SpanData
		want map[string]string
	}{
		{
			name: "attributes",
			data: &export.SpanData{
				Attributes: []label.KeyValue{
					label.String("key", keyValue),
					label.Float64("double", doubleValue),
					label.Uint64("uint", uint64(uintValue)),
					label.Bool("ok", true),
				},
			},
			want: map[string]string{
				"double":                  fmt.Sprint(doubleValue),
				"key":                     keyValue,
				"ok":                      "true",
				"uint":                    strconv.FormatInt(uintValue, 10),
				"otel.status_code":        codes.Unset.String(),
				"otel.status_description": "",
			},
		},
		{
			name: "no attributes",
			data: &export.SpanData{},
			want: map[string]string{
				"otel.status_code":        codes.Unset.String(),
				"otel.status_description": "",
			},
		},
		{
			name: "omit-noerror",
			data: &export.SpanData{
				Attributes: []label.KeyValue{
					label.Bool("error", false),
				},
			},
			want: map[string]string{
				"otel.status_code":        codes.Unset.String(),
				"otel.status_description": "",
			},
		},
		{
			name: "statusCode",
			data: &export.SpanData{
				Attributes: []label.KeyValue{
					label.String("key", keyValue),
					label.Bool("error", true),
				},
				StatusCode:    codes.Error,
				StatusMessage: statusMessage,
			},
			want: map[string]string{
				"error":                   "true",
				"key":                     keyValue,
				"otel.status_code":        codes.Error.String(),
				"otel.status_description": statusMessage,
			},
		},
		{
			name: "instrLib-empty",
			data: &export.SpanData{
				InstrumentationLibrary: instrumentation.Library{},
			},
			want: map[string]string{
				"otel.status_code":        codes.Unset.String(),
				"otel.status_description": "",
			},
		},
		{
			name: "instrLib-noversion",
			data: &export.SpanData{
				Attributes: []label.KeyValue{},
				InstrumentationLibrary: instrumentation.Library{
					Name: instrLibName,
				},
			},
			want: map[string]string{
				"otel.instrumentation_library.name": instrLibName,
				"otel.status_code":                  codes.Unset.String(),
				"otel.status_description":           "",
			},
		},
		{
			name: "instrLib-with-version",
			data: &export.SpanData{
				Attributes: []label.KeyValue{},
				InstrumentationLibrary: instrumentation.Library{
					Name:    instrLibName,
					Version: instrLibVersion,
				},
			},
			want: map[string]string{
				"otel.instrumentation_library.name":    instrLibName,
				"otel.instrumentation_library.version": instrLibVersion,
				"otel.status_code":                     codes.Unset.String(),
				"otel.status_description":              "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toZipkinTags(tt.data)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Diff%v", diff)
			}
		})
	}
}
