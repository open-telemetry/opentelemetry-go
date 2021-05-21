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

package connection

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestGetThrottleDuration(t *testing.T) {
	tts := []struct {
		stsFn    func() (*status.Status, error)
		throttle time.Duration
	}{
		{
			stsFn: func() (*status.Status, error) {
				return status.New(
					codes.OK,
					"status with no retry info",
				), nil
			},
			throttle: 0,
		},
		{
			stsFn: func() (*status.Status, error) {
				st := status.New(codes.ResourceExhausted, "status with retry info")
				return st.WithDetails(
					&errdetails.RetryInfo{RetryDelay: durationpb.New(15 * time.Millisecond)},
				)
			},
			throttle: 15 * time.Millisecond,
		},
		{
			stsFn: func() (*status.Status, error) {
				st := status.New(codes.ResourceExhausted, "status with error info detail")
				return st.WithDetails(
					&errdetails.ErrorInfo{Reason: "no throttle detail"},
				)
			},
			throttle: 0,
		},
		{
			stsFn: func() (*status.Status, error) {
				st := status.New(codes.ResourceExhausted, "status with error info and retry info")
				return st.WithDetails(
					&errdetails.ErrorInfo{Reason: "no throttle detail"},
					&errdetails.RetryInfo{RetryDelay: durationpb.New(13 * time.Minute)},
				)
			},
			throttle: 13 * time.Minute,
		},
		{
			stsFn: func() (*status.Status, error) {
				st := status.New(codes.ResourceExhausted, "status with two retry info should take the first")
				return st.WithDetails(
					&errdetails.RetryInfo{RetryDelay: durationpb.New(13 * time.Minute)},
					&errdetails.RetryInfo{RetryDelay: durationpb.New(18 * time.Minute)},
				)
			},
			throttle: 13 * time.Minute,
		},
	}

	for _, tt := range tts {
		sts, _ := tt.stsFn()
		t.Run(sts.Message(), func(t *testing.T) {
			th := getThrottleDuration(sts)
			require.Equal(t, tt.throttle, th)
		})
	}
}
