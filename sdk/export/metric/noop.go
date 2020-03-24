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

package metric

// NoopLabelEncoder does no encoding at all.
type NoopLabelEncoder struct{}

var _ LabelEncoder = NoopLabelEncoder{}

// Encode is a part of an implementation of the LabelEncoder
// interface. It returns an empty string.
func (NoopLabelEncoder) Encode(LabelIterator) string {
	return ""
}

// ID is a part of an implementation of the LabelEncoder interface.
func (NoopLabelEncoder) ID() int64 {
	return noopLabelEncoderID
}
