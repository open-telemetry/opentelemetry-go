// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pprof_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"os"
	"runtime"
	runtimepprof "runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/profile/internal/pprof"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	semconv "go.opentelemetry.io/otel/semconv/v1.45.0"
	"go.opentelemetry.io/otel/trace"
)

func TestParseAndConvertRuntimeProfiles(t *testing.T) {
	cache := pprof.NewBuildIDCache()

	t.Run("heap", func(t *testing.T) {
		// Retain allocations across a GC cycle so the heap profile always
		// contains live inuse_space and inuse_objects samples on fresh runs.
		sink := make([][]byte, 64)
		for i := range sink {
			sink[i] = make([]byte, 32*1024)
		}
		runtime.GC()

		var buf bytes.Buffer
		require.NoError(t, runtimepprof.Lookup("heap").WriteTo(&buf, 0))
		runtime.KeepAlive(sink)

		parsed, err := pprof.Parse(buf.Bytes())
		require.NoError(t, err)
		require.NotEmpty(t, parsed.SampleType)
		require.NotEmpty(t, parsed.StringTable)
		assert.Empty(t, parsed.StringTable[0])

		profs, scopeAttrs := pprof.Convert(parsed, pprof.ConvertOptions{
			SampleTypes:       []string{"inuse_space", "inuse_objects"},
			DefaultSampleType: "inuse_space",
			FrameType:         semconv.ProfileFrameTypeGo,
			BuildIDs:          cache,
			Instantaneous:     true,
		})
		require.NotEmpty(t, profs)
		assert.Equal(t, profiledata.ValueType{Type: "inuse_space", Unit: "bytes"}, profs[0].SampleType)
		assert.Zero(t, profs[0].Duration)
		assert.Equal(t, profiledata.ValueType{Type: "space", Unit: "bytes"}, profs[0].PeriodType)
		assert.Positive(t, profs[0].Period)
		assert.NotEmpty(t, scopeAttrs)

		// Verify Go executable mapping carries htlhash (and Go/GNU build ID on Linux).
		var foundHTLHash bool
		for _, s := range profs[0].Samples {
			for _, loc := range s.Stack {
				assert.Contains(t, loc.Attributes, profiledata.Attribute{
					Key:   semconv.ProfileFrameTypeKey,
					Value: semconv.ProfileFrameTypeGo.Value,
				})
				if loc.Mapping != nil {
					for _, a := range loc.Mapping.Attributes {
						if a.Key == semconv.ProcessExecutableBuildIDHtlhashKey && a.Value.AsString() != "" {
							foundHTLHash = true
						}
					}
				}
			}
		}
		assert.True(t, foundHTLHash, "expected process.executable.build_id.htlhash on executable mapping")
	})

	t.Run("goroutine", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, runtimepprof.Lookup("goroutine").WriteTo(&buf, 0))

		parsed, err := pprof.Parse(buf.Bytes())
		require.NoError(t, err)

		profs, _ := pprof.Convert(parsed, pprof.ConvertOptions{
			FrameType:     semconv.ProfileFrameTypeGo,
			BuildIDs:      cache,
			Instantaneous: true,
		})
		require.Len(t, profs, 1)
		assert.Equal(t, profiledata.ValueType{Type: "goroutine", Unit: "count"}, profs[0].SampleType)
		assert.NotEmpty(t, profs[0].Samples)
		assert.Zero(t, profs[0].Duration)
	})

	t.Run("cpu_with_trace_labels", func(t *testing.T) {
		const (
			wantTraceHex = "0102030405060708090a0b0c0d0e0f10"
			wantSpanHex  = "0102030405060708"
		)

		var buf bytes.Buffer
		require.NoError(t, runtimepprof.StartCPUProfile(&buf))
		runtimepprof.Do(
			t.Context(),
			runtimepprof.Labels("otel.trace_id", wantTraceHex, "otel.span_id", wantSpanHex, "env", "test"),
			func(context.Context) {
				deadline := time.Now().Add(40 * time.Millisecond)
				acc := 0
				for time.Now().Before(deadline) {
					acc++
				}
				_ = acc
			},
		)
		runtimepprof.StopCPUProfile()

		parsed, err := pprof.Parse(buf.Bytes())
		require.NoError(t, err)

		profs, _ := pprof.Convert(parsed, pprof.ConvertOptions{
			FrameType: semconv.ProfileFrameTypeGo,
			BuildIDs:  cache,
		})
		if len(profs) == 0 {
			// On very slow or virtualized CPUs a 40ms spin may collect 0 samples.
			t.Skip("no CPU samples collected during test spin")
		}
		assert.Equal(t, profiledata.ValueType{Type: "cpu", Unit: "nanoseconds"}, profs[0].SampleType)
		assert.Positive(t, profs[0].Duration)

		wantTraceID, err := trace.TraceIDFromHex(wantTraceHex)
		require.NoError(t, err)
		wantSpanID, err := trace.SpanIDFromHex(wantSpanHex)
		require.NoError(t, err)

		var foundLink, foundEnv bool
		for _, s := range profs[0].Samples {
			if s.Link.TraceID == wantTraceID && s.Link.SpanID == wantSpanID {
				foundLink = true
				for _, a := range s.Attributes {
					assert.NotEqual(t, attribute.Key("otel.trace_id"), a.Key)
					assert.NotEqual(t, attribute.Key("otel.span_id"), a.Key)
					if a.Key == "env" && a.Value.AsString() == "test" {
						foundEnv = true
					}
				}
			}
		}
		assert.True(t, foundLink, "expected otel.trace_id + otel.span_id labels to convert into Sample.Link")
		assert.True(t, foundEnv, "expected non-link labels to remain as Sample.Attributes")
	})

	t.Run("mutex_and_block", func(t *testing.T) {
		prevMutex := runtime.SetMutexProfileFraction(1)
		t.Cleanup(func() { runtime.SetMutexProfileFraction(prevMutex) })
		runtime.SetBlockProfileRate(1)
		t.Cleanup(func() { runtime.SetBlockProfileRate(0) })

		var mu sync.Mutex
		var wg sync.WaitGroup
		ready := make(chan struct{})
		mu.Lock()
		wg.Go(func() {
			close(ready)
			mu.Lock()
			defer mu.Unlock()
		})
		<-ready
		time.Sleep(5 * time.Millisecond)
		mu.Unlock()
		wg.Wait()

		for _, name := range []string{"mutex", "block"} {
			var buf bytes.Buffer
			require.NoError(t, runtimepprof.Lookup(name).WriteTo(&buf, 0))
			parsed, err := pprof.Parse(buf.Bytes())
			require.NoError(t, err)
			profs, _ := pprof.Convert(parsed, pprof.ConvertOptions{FrameType: semconv.ProfileFrameTypeGo})
			require.NotEmpty(t, profs, "expected samples for %s profile", name)
			assert.Equal(t, profiledata.ValueType{Type: "delay", Unit: "nanoseconds"}, profs[0].SampleType)
		}
	})
}

func TestDeltaSubtraction(t *testing.T) {
	// Location IDs intentionally differ between prev (ID=10) and cur (ID=99)
	// to verify semantic frame matching.
	prev := &pprof.Profile{
		StringTable: []string{"", "alloc_objects", "count", "alloc_space", "bytes", "main.work", "main.go", "bytes"},
		SampleType: []pprof.ValueType{
			{Type: 1, Unit: 2},
			{Type: 3, Unit: 4},
		},
		Function: []pprof.Function{
			{ID: 1, Name: 5, Filename: 6, StartLine: 10},
		},
		Location: []pprof.Location{
			{ID: 10, Address: 0x401000, Line: []pprof.Line{{FunctionID: 1, Line: 15}}},
		},
		Sample: []pprof.Sample{
			{
				LocationID: []uint64{10},
				Value:      []int64{10, 1024},
				Label:      []pprof.Label{{Key: 7, Num: 128, NumUnit: 4}},
			},
			{
				LocationID: []uint64{10},
				Value:      []int64{5, 512},
				Label:      []pprof.Label{{Key: 7, Num: 64, NumUnit: 4}},
			},
		},
		TimeNanos: 1_000_000_000,
	}

	cur := &pprof.Profile{
		StringTable: []string{"", "alloc_objects", "count", "alloc_space", "bytes", "main.work", "main.go", "bytes"},
		SampleType: []pprof.ValueType{
			{Type: 1, Unit: 2},
			{Type: 3, Unit: 4},
		},
		Function: []pprof.Function{
			{ID: 42, Name: 5, Filename: 6, StartLine: 10},
		},
		Location: []pprof.Location{
			{ID: 99, Address: 0x401000, Line: []pprof.Line{{FunctionID: 42, Line: 15}}},
		},
		Sample: []pprof.Sample{
			{
				// Grew by +4 objects, +512 bytes.
				LocationID: []uint64{99},
				Value:      []int64{14, 1536},
				Label:      []pprof.Label{{Key: 7, Num: 128, NumUnit: 4}},
			},
			{
				// Unchanged -> must be dropped.
				LocationID: []uint64{99},
				Value:      []int64{5, 512},
				Label:      []pprof.Label{{Key: 7, Num: 64, NumUnit: 4}},
			},
		},
		TimeNanos: 2_500_000_000,
	}

	delta := pprof.Delta(prev, cur)
	require.NotNil(t, delta)
	assert.Equal(t, int64(1_000_000_000), delta.TimeNanos)
	assert.Equal(t, int64(1_500_000_000), delta.DurationNanos)
	require.Len(t, delta.Sample, 1)
	assert.Equal(t, []int64{4, 512}, delta.Sample[0].Value)
}

func TestComputeHTLHash(t *testing.T) {
	t.Run("small_file_under_4096_bytes", func(t *testing.T) {
		data := []byte("hello otel profiles")
		got, err := pprof.ComputeHTLHash(bytes.NewReader(data), int64(len(data)))
		require.NoError(t, err)

		// Per spec: Concat(File[:4096], File[-4096:], BigEndianUInt64(Len(File)))
		// overlaps without padding when Len < 4096.
		h := sha256.New()
		h.Write(data)
		h.Write(data)
		var lb [8]byte
		binary.BigEndian.PutUint64(lb[:], uint64(len(data)))
		h.Write(lb[:])
		want := hex.EncodeToString(h.Sum(nil)[:16])

		assert.Equal(t, want, got)
	})

	t.Run("large_file_over_8192_bytes", func(t *testing.T) {
		data := make([]byte, 10000)
		for i := range data {
			data[i] = byte(i % 251)
		}
		got, err := pprof.ComputeHTLHash(bytes.NewReader(data), int64(len(data)))
		require.NoError(t, err)

		h := sha256.New()
		h.Write(data[:4096])
		h.Write(data[len(data)-4096:])
		var lb [8]byte
		binary.BigEndian.PutUint64(lb[:], uint64(len(data)))
		h.Write(lb[:])
		want := hex.EncodeToString(h.Sum(nil)[:16])

		assert.Equal(t, want, got)
	})

	t.Run("current_test_binary_elf_and_cache", func(t *testing.T) {
		exe, err := os.Executable()
		require.NoError(t, err)

		cache := pprof.NewBuildIDCache()
		ids1 := cache.Lookup(exe)
		ids2 := cache.Lookup(exe)
		assert.Equal(t, ids1, ids2)
		assert.Len(t, ids1.HTLHash, 32)
		if runtime.GOOS == "linux" {
			assert.NotEmpty(t, ids1.GoBuildID)
		}

		// Virtual mapping paths must return empty IDs without error.
		assert.Equal(t, pprof.BuildIDs{}, cache.Lookup("[vdso]"))
		assert.Equal(t, pprof.BuildIDs{}, cache.Lookup(""))
	})
}

func TestParseErrors(t *testing.T) {
	t.Run("corrupt_gzip", func(t *testing.T) {
		_, err := pprof.Parse([]byte{0x1f, 0x8b, 0x00, 0x01})
		assert.Error(t, err)
	})

	t.Run("truncated_varint", func(t *testing.T) {
		_, err := pprof.Parse([]byte{0x08, 0x80})
		assert.Error(t, err)
	})

	t.Run("invalid_wire_type", func(t *testing.T) {
		// Field 1 with wire type 3 (start group, invalid in proto3 pprof).
		_, err := pprof.Parse([]byte{(1 << 3) | 3})
		assert.Error(t, err)
	})

	t.Run("empty_payload", func(t *testing.T) {
		p, err := pprof.Parse(nil)
		require.NoError(t, err)
		assert.Empty(t, p.Sample)
	})

	t.Run("uncompressed_valid_profile", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, runtimepprof.Lookup("goroutine").WriteTo(&buf, 0))
		gr, err := gzip.NewReader(&buf)
		require.NoError(t, err)
		var raw bytes.Buffer
		_, err = raw.ReadFrom(gr)
		require.NoError(t, err)
		require.NoError(t, gr.Close())

		p, err := pprof.Parse(raw.Bytes())
		require.NoError(t, err)
		assert.NotEmpty(t, p.Sample)
	})
}

func TestConvertPartialOrInvalidTraceLabelsRemainAttributes(t *testing.T) {
	p := &pprof.Profile{
		StringTable: []string{"", "cpu", "nanoseconds", "otel.trace_id", "0102030405060708090a0b0c0d0e0f10", "main"},
		SampleType:  []pprof.ValueType{{Type: 1, Unit: 2}},
		Function:    []pprof.Function{{ID: 1, Name: 5}},
		Location:    []pprof.Location{{ID: 1, Line: []pprof.Line{{FunctionID: 1, Line: 1}}}},
		Sample: []pprof.Sample{
			{
				LocationID: []uint64{1},
				Value:      []int64{100},
				// Only trace_id without span_id cannot form a valid OTLP Link.
				Label: []pprof.Label{{Key: 3, Str: 4}},
			},
		},
	}

	profs, _ := pprof.Convert(p, pprof.ConvertOptions{})
	require.Len(t, profs, 1)
	require.Len(t, profs[0].Samples, 1)
	assert.False(t, profs[0].Samples[0].Link.IsValid())
	assert.Equal(t, []profiledata.Attribute{
		{Key: "otel.trace_id", Value: attribute.StringValue("0102030405060708090a0b0c0d0e0f10")},
	}, profs[0].Samples[0].Attributes)
}

func TestConvertMergesDuplicateSamplesAndDeduplicatesLines(t *testing.T) {
	p := &pprof.Profile{
		StringTable: []string{
			"",
			"samples",
			"count",
			"foo",
			"foo.go",
			"drop.*",
			"keep.*",
			"https://example.com",
			"note1",
		},
		SampleType: []pprof.ValueType{{Type: 1, Unit: 2}},
		Function:   []pprof.Function{{ID: 1, Name: 3, Filename: 4, StartLine: 1}},
		Location: []pprof.Location{
			{
				ID:       1,
				IsFolded: true,
				Line: []pprof.Line{
					{FunctionID: 1, Line: 10},
					{FunctionID: 1, Line: 10}, // Duplicate line should be deduplicated.
				},
			},
		},
		Sample: []pprof.Sample{
			{LocationID: []uint64{1}, Value: []int64{3}},
			{LocationID: []uint64{1}, Value: []int64{7}},
		},
		DropFrames: 5,
		KeepFrames: 6,
		DocURL:     7,
		Comment:    []int64{8},
	}

	profs, _ := pprof.Convert(p, pprof.ConvertOptions{})
	require.Len(t, profs, 1)
	require.Len(t, profs[0].Samples, 1)
	assert.Equal(t, []int64{10}, profs[0].Samples[0].Values)
	require.Len(t, profs[0].Samples[0].Stack, 1)
	assert.Len(t, profs[0].Samples[0].Stack[0].Lines, 1)
	assert.Len(t, profs[0].Attributes, 4)
}

var (
	benchHeapOnce sync.Once
	benchHeapRaw  []byte
)

// Recommended CLI invocation:
// go test ./sdk/profile/internal/pprof -run '^$' -bench '^BenchmarkParseAndConvertHeap$' -benchtime 1s -count 6 -cpu 2 -timeout 999m -benchmem | tee pprof_bench.txt.
func BenchmarkParseAndConvertHeap(b *testing.B) {
	benchHeapOnce.Do(func() {
		sink := make([][]byte, 64)
		for i := range sink {
			sink[i] = make([]byte, 32*1024)
		}
		runtime.GC()

		var buf bytes.Buffer
		if err := runtimepprof.Lookup("heap").WriteTo(&buf, 0); err != nil {
			b.Fatal(err)
		}
		runtime.KeepAlive(sink)
		benchHeapRaw = buf.Bytes()
	})
	cache := pprof.NewBuildIDCache()
	opts := pprof.ConvertOptions{
		FrameType:     semconv.ProfileFrameTypeGo,
		BuildIDs:      cache,
		Instantaneous: true,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		parsed, err := pprof.Parse(benchHeapRaw)
		if err != nil {
			b.Fatal(err)
		}
		profs, _ := pprof.Convert(parsed, opts)
		if len(profs) == 0 {
			b.Fatal("expected converted profiles")
		}
	}
}
