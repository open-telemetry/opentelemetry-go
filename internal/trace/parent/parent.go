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

package parent

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

func GetSpanReferenceAndLinks(ctx context.Context, ignoreContext bool) (otel.SpanReference, bool, []otel.Link) {
	lsref := otel.SpanFromContext(ctx).SpanReference()
	rsref := otel.RemoteSpanReferenceFromContext(ctx)

	if ignoreContext {
		links := addLinkIfValid(nil, lsref, "current")
		links = addLinkIfValid(links, rsref, "remote")

		return otel.SpanReference{}, false, links
	}
	if lsref.IsValid() {
		return lsref, false, nil
	}
	if rsref.IsValid() {
		return rsref, true, nil
	}
	return otel.SpanReference{}, false, nil
}

func addLinkIfValid(links []otel.Link, sr otel.SpanReference, kind string) []otel.Link {
	if !sr.IsValid() {
		return links
	}
	return append(links, otel.Link{
		SpanReference: sr,
		Attributes: []label.KeyValue{
			label.String("ignored-on-demand", kind),
		},
	})
}
