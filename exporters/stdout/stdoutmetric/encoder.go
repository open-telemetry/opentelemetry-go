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

//go:build go1.18
// +build go1.18

package stdoutmetric

import "errors"

var ErrUnrecognized = errors.New("unrecognized metric data")

// Encoder encodes and outputs OpenTelemetry metric data-types as human
// readable text.
type Encoder interface {
	// Encode handles the encoding and writing OpenTelemetry metric data-types
	// that the exporter will pass to it.
	//
	// Any data-type that is not recognized by the encoder and not output to
	// the user, will have an ErrUnrecognized returned from Encode.
	Encode(v any) error
}

// encoderHolder is the concrete type used to wrap an Encoder so it can be
// used as a atomic.Value type.
type encoderHolder struct {
	encoder Encoder
}

func (e encoderHolder) Encode(v any) error { return e.encoder.Encode(v) }

// shutdownEncoder is used when the exporter is shutdown. It always returns
// errShutdown when Encode is called.
type shutdownEncoder struct{}

var errShutdown = errors.New("exporter shutdown")

func (shutdownEncoder) Encode(any) error { return errShutdown }
