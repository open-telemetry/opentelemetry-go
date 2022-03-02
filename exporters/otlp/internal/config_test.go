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

package internal

import "testing"

func TestCleanPath(t *testing.T) {
	type args struct {
		urlPath     string
		defaultPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "clean empty path",
			args: args{
				urlPath:     "",
				defaultPath: "DefaultPath",
			},
			want: "DefaultPath",
		},
		{
			name: "clean metrics path",
			args: args{
				urlPath:     "/prefix/v1/metrics",
				defaultPath: "DefaultMetricsPath",
			},
			want: "/prefix/v1/metrics",
		},
		{
			name: "clean traces path",
			args: args{
				urlPath:     "https://env_endpoint",
				defaultPath: "DefaultTracesPath",
			},
			want: "/https:/env_endpoint",
		},
		{
			name: "spaces trimmed",
			args: args{
				urlPath: " /dir",
			},
			want: "/dir",
		},
		{
			name: "clean path empty",
			args: args{
				urlPath:     "dir/..",
				defaultPath: "DefaultTracesPath",
			},
			want: "DefaultTracesPath",
		},
		{
			name: "make absolute",
			args: args{
				urlPath: "dir/a",
			},
			want: "/dir/a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanPath(tt.args.urlPath, tt.args.defaultPath); got != tt.want {
				t.Errorf("CleanPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
