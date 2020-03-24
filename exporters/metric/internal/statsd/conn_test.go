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

package statsd_test

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/unit"
	"go.opentelemetry.io/otel/exporters/metric/internal/statsd"
	"go.opentelemetry.io/otel/exporters/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

// withTagsAdapter tests a dogstatsd-style statsd exporter.
type withTagsAdapter struct {
	*statsd.LabelEncoder
}

func (*withTagsAdapter) AppendName(rec export.Record, buf *bytes.Buffer) {
	_, _ = buf.WriteString(rec.Descriptor().Name())
}

func (ta *withTagsAdapter) AppendTags(rec export.Record, buf *bytes.Buffer) {
	encoded, _ := ta.LabelEncoder.ForceEncode(rec.Labels())
	_, _ = buf.WriteString(encoded)
}

func newWithTagsAdapter() *withTagsAdapter {
	return &withTagsAdapter{
		statsd.NewLabelEncoder(),
	}
}

// noTagsAdapter simulates a plain-statsd exporter that appends tag
// values to the metric name.
type noTagsAdapter struct {
}

func (*noTagsAdapter) AppendName(rec export.Record, buf *bytes.Buffer) {
	_, _ = buf.WriteString(rec.Descriptor().Name())

	iter := rec.Labels().Iter()
	for iter.Next() {
		tag := iter.Label()
		_, _ = buf.WriteString(".")
		_, _ = buf.WriteString(tag.Value.Emit())
	}
}

func (*noTagsAdapter) AppendTags(rec export.Record, buf *bytes.Buffer) {
}

func newNoTagsAdapter() *noTagsAdapter {
	return &noTagsAdapter{}
}

type testWriter struct {
	vec []string
}

func (w *testWriter) Write(b []byte) (int, error) {
	w.vec = append(w.vec, string(b))
	return len(b), nil
}

func TestBasicFormat(t *testing.T) {
	type adapterOutput struct {
		adapter  statsd.Adapter
		expected string
	}

	for _, ao := range []adapterOutput{{
		adapter: newWithTagsAdapter(),
		expected: `counter:%s|c|#A:B,C:D
observer:%s|g|#A:B,C:D
measure:%s|h|#A:B,C:D
timer:%s|ms|#A:B,C:D
`}, {
		adapter: newNoTagsAdapter(),
		expected: `counter.B.D:%s|c
observer.B.D:%s|g
measure.B.D:%s|h
timer.B.D:%s|ms
`},
	} {
		adapter := ao.adapter
		expected := ao.expected
		t.Run(fmt.Sprintf("%T", adapter), func(t *testing.T) {
			for _, nkind := range []core.NumberKind{
				core.Float64NumberKind,
				core.Int64NumberKind,
			} {
				t.Run(nkind.String(), func(t *testing.T) {
					ctx := context.Background()
					writer := &testWriter{}
					config := statsd.Config{
						Writer:        writer,
						MaxPacketSize: 1024,
					}
					exp, err := statsd.NewExporter(config, adapter)
					if err != nil {
						t.Fatal("New error: ", err)
					}

					checkpointSet := test.NewCheckpointSet(sdk.NewDefaultLabelEncoder())
					cdesc := metric.NewDescriptor(
						"counter", metric.CounterKind, nkind)
					gdesc := metric.NewDescriptor(
						"observer", metric.ObserverKind, nkind)
					mdesc := metric.NewDescriptor(
						"measure", metric.MeasureKind, nkind)
					tdesc := metric.NewDescriptor(
						"timer", metric.MeasureKind, nkind, metric.WithUnit(unit.Milliseconds))

					labels := []core.KeyValue{
						key.New("A").String("B"),
						key.New("C").String("D"),
					}
					const value = 123.456

					checkpointSet.AddCounter(&cdesc, value, labels...)
					checkpointSet.AddLastValue(&gdesc, value, labels...)
					checkpointSet.AddMeasure(&mdesc, value, labels...)
					checkpointSet.AddMeasure(&tdesc, value, labels...)

					err = exp.Export(ctx, checkpointSet)
					require.Nil(t, err)

					var vfmt string
					if nkind == core.Int64NumberKind {
						fv := value
						vfmt = strconv.FormatInt(int64(fv), 10)
					} else {
						vfmt = strconv.FormatFloat(value, 'g', -1, 64)
					}

					require.Equal(t, 1, len(writer.vec))
					require.Equal(t, fmt.Sprintf(expected, vfmt, vfmt, vfmt, vfmt), writer.vec[0])
				})
			}
		})
	}
}

func makeLabels(offset, nkeys int) []core.KeyValue {
	r := make([]core.KeyValue, nkeys)
	for i := range r {
		r[i] = key.New(fmt.Sprint("k", offset+i)).String(fmt.Sprint("v", offset+i))
	}
	return r
}

type splitTestCase struct {
	name  string
	setup func(add func(int))
	check func(expected, got []string, t *testing.T)
}

var splitTestCases = []splitTestCase{
	// These test use the number of keys to control where packets
	// are split.
	{"Simple",
		func(add func(int)) {
			add(1)
			add(1000)
			add(1)
		},
		func(expected, got []string, t *testing.T) {
			require.EqualValues(t, expected, got)
		},
	},
	{"LastBig",
		func(add func(int)) {
			add(1)
			add(1)
			add(1000)
		},
		func(expected, got []string, t *testing.T) {
			require.Equal(t, 2, len(got))
			require.EqualValues(t, []string{
				expected[0] + expected[1],
				expected[2],
			}, got)
		},
	},
	{"FirstBig",
		func(add func(int)) {
			add(1000)
			add(1)
			add(1)
			add(1000)
			add(1)
			add(1)
		},
		func(expected, got []string, t *testing.T) {
			require.Equal(t, 4, len(got))
			require.EqualValues(t, []string{
				expected[0],
				expected[1] + expected[2],
				expected[3],
				expected[4] + expected[5],
			}, got)
		},
	},
	{"OneBig",
		func(add func(int)) {
			add(1000)
		},
		func(expected, got []string, t *testing.T) {
			require.EqualValues(t, expected, got)
		},
	},
	{"LastSmall",
		func(add func(int)) {
			add(1000)
			add(1)
		},
		func(expected, got []string, t *testing.T) {
			require.EqualValues(t, expected, got)
		},
	},
	{"Overflow",
		func(add func(int)) {
			for i := 0; i < 1000; i++ {
				add(1)
			}
		},
		func(expected, got []string, t *testing.T) {
			require.Less(t, 1, len(got))
			require.Equal(t, strings.Join(expected, ""), strings.Join(got, ""))
		},
	},
	{"Empty",
		func(add func(int)) {
		},
		func(expected, got []string, t *testing.T) {
			require.Equal(t, 0, len(got))
		},
	},
	{"AllBig",
		func(add func(int)) {
			add(1000)
			add(1000)
			add(1000)
		},
		func(expected, got []string, t *testing.T) {
			require.EqualValues(t, expected, got)
		},
	},
}

func TestPacketSplit(t *testing.T) {
	for _, tcase := range splitTestCases {
		t.Run(tcase.name, func(t *testing.T) {
			ctx := context.Background()
			writer := &testWriter{}
			config := statsd.Config{
				Writer:        writer,
				MaxPacketSize: 1024,
			}
			adapter := newWithTagsAdapter()
			exp, err := statsd.NewExporter(config, adapter)
			if err != nil {
				t.Fatal("New error: ", err)
			}

			checkpointSet := test.NewCheckpointSet(adapter.LabelEncoder)
			desc := metric.NewDescriptor("counter", metric.CounterKind, core.Int64NumberKind)

			var expected []string

			offset := 0
			tcase.setup(func(nkeys int) {
				labels := makeLabels(offset, nkeys)
				offset += nkeys
				iter := export.LabelSlice(labels).Iter()
				encoded := adapter.LabelEncoder.Encode(iter)
				expect := fmt.Sprint("counter:100|c", encoded, "\n")
				expected = append(expected, expect)
				checkpointSet.AddCounter(&desc, 100, labels...)
			})

			err = exp.Export(ctx, checkpointSet)
			require.Nil(t, err)

			tcase.check(expected, writer.vec, t)
		})
	}
}
