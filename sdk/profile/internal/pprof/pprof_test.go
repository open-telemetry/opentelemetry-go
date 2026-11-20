// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pprof_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"debug/elf"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		// Go's heap profile has no DefaultSampleType; filtered 2-type order is dense [1, 0].
		assert.Equal(t, []attribute.KeyValue{
			semconv.PprofScopeSampleTypeOrder(1, 0),
		}, scopeAttrs)

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

		contend := func() {
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
		}

		for _, name := range []string{"mutex", "block"} {
			var profs []profiledata.Profile
			deadline := time.Now().Add(2 * time.Second)
			for time.Now().Before(deadline) {
				contend()
				var buf bytes.Buffer
				require.NoError(t, runtimepprof.Lookup(name).WriteTo(&buf, 0))
				parsed, err := pprof.Parse(buf.Bytes())
				require.NoError(t, err)
				profs, _ = pprof.Convert(parsed, pprof.ConvertOptions{FrameType: semconv.ProfileFrameTypeGo})
				if len(profs) > 0 {
					break
				}
			}
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

	t.Run("nil_and_mismatched_sample_types", func(t *testing.T) {
		assert.Nil(t, pprof.Delta(prev, nil))
		assert.Same(t, cur, pprof.Delta(nil, cur))

		diffLen := &pprof.Profile{
			StringTable: []string{"", "alloc_objects", "count"},
			SampleType:  []pprof.ValueType{{Type: 1, Unit: 2}},
		}
		assert.Same(t, cur, pprof.Delta(diffLen, cur))

		diffName := &pprof.Profile{
			StringTable: []string{"", "alloc_objects", "count", "other", "bytes"},
			SampleType:  []pprof.ValueType{{Type: 1, Unit: 2}, {Type: 3, Unit: 4}},
		}
		assert.Same(t, cur, pprof.Delta(diffName, cur))
	})

	t.Run("aggregation_clamping_multilabel_and_collision_safety", func(t *testing.T) {
		st := []string{
			"", "alloc_objects", "count", "alloc_space", "bytes",
			"main.work", "main.go",
			"k1", "k2", "v1", "v2", "u1", "u2",
			"a=s:1", "2", "a", "1=s:2",
		}
		makeProf := func(tNanos int64, samples []pprof.Sample) *pprof.Profile {
			return &pprof.Profile{
				StringTable: st,
				SampleType:  []pprof.ValueType{{Type: 1, Unit: 2}, {Type: 3, Unit: 4}},
				Function:    []pprof.Function{{ID: 1, Name: 5, Filename: 6}},
				Location:    []pprof.Location{{ID: 1, Address: 0x100, Line: []pprof.Line{{FunctionID: 1, Line: 10}}}},
				Sample:      samples,
				TimeNanos:   tNanos,
			}
		}

		p1 := makeProf(100, []pprof.Sample{
			// Duplicate unlabeled samples aggregate before subtraction.
			{LocationID: []uint64{1}, Value: []int64{3, 30}},
			{LocationID: []uint64{1}, Value: []int64{2, 20}},
			// Partial negative clamping: first value drops below prev, second grows.
			{LocationID: []uint64{999}, Value: []int64{10, 100}},
			// Multi-label sample with unsorted labels exercising full comparator.
			{
				LocationID: []uint64{1},
				Value:      []int64{1, 10},
				Label: []pprof.Label{
					{Key: 8, Num: 2, NumUnit: 12},
					{Key: 8, Num: 2, NumUnit: 11},
					{Key: 8, Num: 1, NumUnit: 11},
					{Key: 8, Str: 10},
					{Key: 8, Str: 9},
					{Key: 7, Str: 9},
				},
			},
			// Delimiter collision pair: must stay distinct.
			{LocationID: []uint64{1}, Value: []int64{5, 50}, Label: []pprof.Label{{Key: 13, Str: 14}}},
		})
		p2 := makeProf(250, []pprof.Sample{
			// Duplicate in cur aggregates to [8, 90] -> delta [3, 40].
			{LocationID: []uint64{1}, Value: []int64{4, 40}},
			{LocationID: []uint64{1}, Value: []int64{4, 50}},
			// Partial negative: [8, 120] - [10, 100] -> [0, 20].
			{LocationID: []uint64{999}, Value: []int64{8, 120}},
			// Same multi-label sample in different label order -> delta [2, 20].
			{
				LocationID: []uint64{1},
				Value:      []int64{3, 30},
				Label: []pprof.Label{
					{Key: 7, Str: 9},
					{Key: 8, Str: 9},
					{Key: 8, Str: 10},
					{Key: 8, Num: 1, NumUnit: 11},
					{Key: 8, Num: 2, NumUnit: 11},
					{Key: 8, Num: 2, NumUnit: 12},
				},
			},
			// First collision-pair key unchanged (delta 0, dropped), second is brand-new ([7, 70]).
			{LocationID: []uint64{1}, Value: []int64{5, 50}, Label: []pprof.Label{{Key: 13, Str: 14}}},
			{LocationID: []uint64{1}, Value: []int64{7, 70}, Label: []pprof.Label{{Key: 15, Str: 16}}},
		})

		d1 := pprof.Delta(p1, p2)
		require.Len(t, d1.Sample, 4)
		assert.Equal(t, []int64{3, 40}, d1.Sample[0].Value)
		assert.Equal(t, []int64{0, 20}, d1.Sample[1].Value)
		assert.Equal(t, []int64{2, 20}, d1.Sample[2].Value)
		assert.Equal(t, []int64{7, 70}, d1.Sample[3].Value)

		// Verify deterministic sample ordering across repeated calls.
		for range 10 {
			d2 := pprof.Delta(p1, p2)
			require.Len(t, d2.Sample, 4)
			for i := range d1.Sample {
				assert.Equal(t, d1.Sample[i].Value, d2.Sample[i].Value)
			}
		}
	})
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
		assert.Equal(t, pprof.BuildIDs{}, cache.Lookup("relative/path"))
		assert.Equal(t, pprof.BuildIDs{}, (*pprof.BuildIDCache)(nil).Lookup(exe))
		assert.Equal(t, pprof.BuildIDs{}, cache.Lookup(t.TempDir()))
		assert.Equal(t, pprof.BuildIDs{}, cache.Lookup(filepath.Join(t.TempDir(), "does-not-exist")))

		// Cache eviction past maxBuildIDCache (256) succeeds without unbounded growth.
		dir := t.TempDir()
		for i := range 260 {
			_ = cache.Lookup(filepath.Join(dir, fmt.Sprintf("missing-%d", i)))
		}
		assert.Equal(t, ids1, cache.Lookup(exe))
	})

	t.Run("boundary_sizes_and_read_errors", func(t *testing.T) {
		_, err := pprof.ComputeHTLHash(bytes.NewReader(nil), -1)
		assert.ErrorIs(t, err, io.ErrUnexpectedEOF)

		for _, sz := range []int{0, 4096, 5000, 8192} {
			data := make([]byte, sz)
			for i := range data {
				data[i] = byte(i*31 + 7)
			}
			got, err := pprof.ComputeHTLHash(bytes.NewReader(data), int64(sz))
			require.NoError(t, err)
			assert.Len(t, got, 32)
		}

		// Full read accompanied by io.EOF must succeed.
		eofReader := readerAtFunc(func(p []byte, off int64) (int, error) {
			for i := range p {
				p[i] = byte(off) + byte(i)
			}
			return len(p), io.EOF
		})
		got, err := pprof.ComputeHTLHash(eofReader, 16)
		require.NoError(t, err)
		assert.Len(t, got, 32)

		// Short head read with io.EOF must fail with io.ErrUnexpectedEOF.
		_, err = pprof.ComputeHTLHash(bytes.NewReader([]byte("short")), 100)
		assert.ErrorIs(t, err, io.ErrUnexpectedEOF)

		// Short tail read with io.EOF must fail with io.ErrUnexpectedEOF.
		_, err = pprof.ComputeHTLHash(bytes.NewReader(make([]byte, 5000)), 10000)
		assert.ErrorIs(t, err, io.ErrUnexpectedEOF)

		// Non-EOF read errors propagate on both head and tail.
		wantErr := errors.New("disk failure")
		_, err = pprof.ComputeHTLHash(readerAtFunc(func([]byte, int64) (int, error) {
			return 0, wantErr
		}), 100)
		assert.ErrorIs(t, err, wantErr)

		_, err = pprof.ComputeHTLHash(readerAtFunc(func(p []byte, off int64) (int, error) {
			if off > 0 {
				return 0, wantErr
			}
			return len(p), nil
		}), 10000)
		assert.ErrorIs(t, err, wantErr)
	})
}

type readerAtFunc func(p []byte, off int64) (int, error)

func (f readerAtFunc) ReadAt(p []byte, off int64) (int, error) { return f(p, off) }

func buildSyntheticELFWithNotes(noteSections ...[]byte) []byte {
	const (
		ehdrSize = 64
		shdrSize = 64
	)
	var out []byte
	// Reserve ELF64 header, filled once section header table offset is known.
	out = append(out, make([]byte, ehdrSize)...)

	type secMeta struct {
		offset uint64
		size   uint64
	}
	metas := make([]secMeta, len(noteSections))
	for i, data := range noteSections {
		metas[i] = secMeta{offset: uint64(len(out)), size: uint64(len(data))}
		out = append(out, data...)
	}

	shoff := uint64(len(out))
	// Section 0: SHT_NULL.
	out = append(out, make([]byte, shdrSize)...)
	// Sections 1..N: SHT_NOTE.
	for _, m := range metas {
		var sh [shdrSize]byte
		binary.LittleEndian.PutUint32(sh[4:8], uint32(elf.SHT_NOTE))
		binary.LittleEndian.PutUint64(sh[24:32], m.offset)
		binary.LittleEndian.PutUint64(sh[32:40], m.size)
		out = append(out, sh[:]...)
	}

	copy(out[0:4], elf.ELFMAG)
	out[elf.EI_CLASS] = byte(elf.ELFCLASS64)
	out[elf.EI_DATA] = byte(elf.ELFDATA2LSB)
	out[elf.EI_VERSION] = byte(elf.EV_CURRENT)
	binary.LittleEndian.PutUint16(out[16:18], uint16(elf.ET_EXEC))
	binary.LittleEndian.PutUint16(out[18:20], uint16(elf.EM_X86_64))
	binary.LittleEndian.PutUint32(out[20:24], uint32(elf.EV_CURRENT))
	binary.LittleEndian.PutUint64(out[40:48], shoff)
	binary.LittleEndian.PutUint16(out[52:54], ehdrSize)
	binary.LittleEndian.PutUint16(out[58:60], shdrSize)
	binary.LittleEndian.PutUint16(out[60:62], uint16(1+len(noteSections))) //nolint:gosec // Small test slice.
	return out
}

func appendELFNote(b []byte, name string, noteType uint32, desc []byte) []byte {
	nameBytes := append([]byte(name), 0)
	b = binary.LittleEndian.AppendUint32(b, uint32(len(nameBytes))) //nolint:gosec // Small test note.
	b = binary.LittleEndian.AppendUint32(b, uint32(len(desc)))      //nolint:gosec // Small test note.
	b = binary.LittleEndian.AppendUint32(b, noteType)
	b = append(b, nameBytes...)
	for len(b)%4 != 0 {
		b = append(b, 0)
	}
	b = append(b, desc...)
	for len(b)%4 != 0 {
		b = append(b, 0)
	}
	return b
}

func TestExtractELFBuildIDsSynthetic(t *testing.T) {
	// Non-ELF data returns empty IDs.
	goID, gnuID := pprof.ExtractELFBuildIDs(bytes.NewReader([]byte("not an elf")))
	assert.Empty(t, goID)
	assert.Empty(t, gnuID)

	// Valid GNU (type 3) + Go (type 4) notes with non-4-aligned lengths.
	var notes []byte
	notes = appendELFNote(notes, "GNU", 3, []byte{0xde, 0xad, 0xbe, 0xef, 0x01})
	notes = appendELFNote(notes, "Go", 4, []byte("go-build-id/abc"))
	// Oversized note section (>64 KiB) must be skipped safely.
	oversized := make([]byte, 65*1024)
	elfBytes := buildSyntheticELFWithNotes(oversized, notes)

	goID, gnuID = pprof.ExtractELFBuildIDs(bytes.NewReader(elfBytes))
	assert.Equal(t, "go-build-id/abc", goID)
	assert.Equal(t, "deadbeef01", gnuID)

	// Overflowing nameSz/descSz (0xfffffffd) and truncated note must not panic.
	var overflowNote [12]byte
	binary.LittleEndian.PutUint32(overflowNote[0:4], 0xfffffffd)
	binary.LittleEndian.PutUint32(overflowNote[4:8], 0)
	binary.LittleEndian.PutUint32(overflowNote[8:12], 3)

	var overflowDesc [12]byte
	binary.LittleEndian.PutUint32(overflowDesc[0:4], 0)
	binary.LittleEndian.PutUint32(overflowDesc[4:8], 0xfffffffd)
	binary.LittleEndian.PutUint32(overflowDesc[8:12], 4)

	malformedELF := buildSyntheticELFWithNotes(overflowNote[:], overflowDesc[:])
	goID, gnuID = pprof.ExtractELFBuildIDs(bytes.NewReader(malformedELF))
	assert.Empty(t, goID)
	assert.Empty(t, gnuID)
}

func TestBuildIDCacheConcurrentSafe(t *testing.T) {
	exe, err := os.Executable()
	require.NoError(t, err)

	cache := pprof.NewBuildIDCache()
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			ids := cache.Lookup(exe)
			assert.Len(t, ids.HTLHash, 32)
			assert.Equal(t, pprof.BuildIDs{}, cache.Lookup("[vdso]"))
		})
	}
	wg.Wait()
}

func appendVarintField(b []byte, num, v uint64) []byte {
	b = binary.AppendUvarint(b, num<<3)
	return binary.AppendUvarint(b, v)
}

func appendBytesField(b []byte, num uint64, data []byte) []byte {
	b = binary.AppendUvarint(b, (num<<3)|2)
	b = binary.AppendUvarint(b, uint64(len(data)))
	return append(b, data...)
}

func appendFixed64Field(b []byte, num, v uint64) []byte {
	b = binary.AppendUvarint(b, (num<<3)|1)
	return binary.LittleEndian.AppendUint64(b, v)
}

func appendFixed32Field(b []byte, num uint64, v uint32) []byte {
	b = binary.AppendUvarint(b, (num<<3)|5)
	return binary.LittleEndian.AppendUint32(b, v)
}

func appendUnknownFields(b []byte) []byte {
	b = appendVarintField(b, 99, 123)
	b = appendFixed64Field(b, 100, 0x0102030405060708)
	b = appendBytesField(b, 101, []byte("skipme"))
	b = appendFixed32Field(b, 102, 0x01020304)
	return b
}

func TestParseErrors(t *testing.T) {
	t.Run("corrupt_gzip", func(t *testing.T) {
		_, err := pprof.Parse([]byte{0x1f, 0x8b, 0x00, 0x01})
		assert.ErrorContains(t, err, "gzip reader")
	})

	t.Run("truncated_gzip_body", func(t *testing.T) {
		var gz bytes.Buffer
		gw := gzip.NewWriter(&gz)
		_, err := gw.Write([]byte("hello world"))
		require.NoError(t, err)
		require.NoError(t, gw.Close())

		raw := gz.Bytes()
		_, err = pprof.Parse(raw[:len(raw)-4])
		assert.ErrorContains(t, err, "gzip read")
	})

	t.Run("oversized_gzip_rejected", func(t *testing.T) {
		var gz bytes.Buffer
		gw, err := gzip.NewWriterLevel(&gz, gzip.BestSpeed)
		require.NoError(t, err)
		chunk := make([]byte, 1<<20)
		for range 65 {
			_, err = gw.Write(chunk)
			require.NoError(t, err)
		}
		require.NoError(t, gw.Close())

		_, err = pprof.Parse(gz.Bytes())
		assert.ErrorContains(t, err, "exceeds maximum size")
	})

	t.Run("truncation_errors", func(t *testing.T) {
		cases := []struct {
			name string
			raw  []byte
		}{
			{name: "truncated_tag", raw: []byte{0x80}},
			{name: "truncated_varint", raw: []byte{9 << 3, 0x80}},
			{name: "truncated_bytes_len", raw: []byte{(6 << 3) | 2, 0x80}},
			{name: "truncated_bytes_body", raw: []byte{(6 << 3) | 2, 5, 'a', 'b'}},
			{name: "truncated_skip_fixed64", raw: append(binary.AppendUvarint(nil, (99<<3)|1), 1, 2, 3)},
			{name: "truncated_skip_fixed32", raw: append(binary.AppendUvarint(nil, (99<<3)|5), 1, 2)},
			{
				name: "truncated_packed_uint64",
				raw:  appendBytesField(nil, 2, appendBytesField(nil, 1, []byte{0x80})),
			},
			{
				name: "truncated_packed_int64",
				raw:  appendBytesField(nil, 2, appendBytesField(nil, 2, []byte{0x80})),
			},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := pprof.Parse(tc.raw)
				assert.ErrorContains(t, err, "truncated")
			})
		}
	})

	t.Run("invalid_field_and_wire_types", func(t *testing.T) {
		_, err := pprof.Parse([]byte{0x00})
		assert.ErrorContains(t, err, "invalid protobuf field number")

		badWireCases := []struct {
			name string
			raw  []byte
		}{
			{name: "profile_start_group", raw: []byte{(1 << 3) | 3}},
			{name: "bytes_field_with_varint", raw: []byte{1 << 3, 0x01}},
			{name: "varint_field_with_bytes", raw: appendBytesField(nil, 9, []byte{1})},
			{name: "unknown_invalid_wire_type", raw: binary.AppendUvarint(nil, (99<<3)|3)},
			{name: "repeated_uint64_bad_wire", raw: appendBytesField(nil, 2, appendFixed32Field(nil, 1, 1))},
			{name: "repeated_int64_bad_wire", raw: appendBytesField(nil, 2, appendFixed64Field(nil, 2, 1))},
			{name: "nested_sample_type", raw: appendBytesField(nil, 1, []byte{0x00})},
			{name: "nested_sample", raw: appendBytesField(nil, 2, []byte{0x00})},
			{name: "nested_label", raw: appendBytesField(nil, 2, appendBytesField(nil, 3, []byte{0x00}))},
			{name: "nested_mapping", raw: appendBytesField(nil, 3, []byte{0x00})},
			{name: "nested_location", raw: appendBytesField(nil, 4, []byte{0x00})},
			{name: "nested_line", raw: appendBytesField(nil, 4, appendBytesField(nil, 4, []byte{0x00}))},
			{name: "nested_function", raw: appendBytesField(nil, 5, []byte{0x00})},
			{name: "nested_period_type", raw: appendBytesField(nil, 11, []byte{0x00})},
		}
		for _, tc := range badWireCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := pprof.Parse(tc.raw)
				assert.Error(t, err)
			})
		}
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

func TestParseFullSyntheticProfileAndUnknownFields(t *testing.T) {
	vtBytes := appendUnknownFields(appendVarintField(appendVarintField(nil, 1, 1), 2, 2))
	lblBytes := appendUnknownFields(
		appendVarintField(
			appendVarintField(
				appendVarintField(appendVarintField(nil, 1, 3), 2, 0),
				3, 42,
			),
			4, 2,
		),
	)
	// Mix unpacked and packed repeated LocationID/Value fields.
	sampleBytes := appendUnknownFields(
		appendBytesField(
			appendVarintField(
				appendVarintField(nil, 1, 10),
				2, 99,
			),
			3, lblBytes,
		),
	)
	mappingBytes := appendUnknownFields(
		appendVarintField(
			appendVarintField(
				appendVarintField(
					appendVarintField(
						appendVarintField(
							appendVarintField(
								appendVarintField(
									appendVarintField(
										appendVarintField(
											appendVarintField(nil, 1, 1),
											2, 0x1000,
										),
										3, 0x2000,
									),
									4, 0,
								),
								5, 4,
							),
							6, 5,
						),
						7, 1,
					),
					8, 1,
				),
				9, 1,
			),
			10, 1,
		),
	)
	lineBytes := appendUnknownFields(
		appendVarintField(
			appendVarintField(appendVarintField(nil, 1, 7), 2, 12),
			3, 34,
		),
	)
	locBytes := appendUnknownFields(
		appendVarintField(
			appendBytesField(
				appendVarintField(
					appendVarintField(appendVarintField(nil, 1, 10), 2, 1),
					3, 0x1010,
				),
				4, lineBytes,
			),
			5, 1,
		),
	)
	fnBytes := appendUnknownFields(
		appendVarintField(
			appendVarintField(
				appendVarintField(
					appendVarintField(appendVarintField(nil, 1, 7), 2, 6),
					3, 6,
				),
				4, 7,
			),
			5, 10,
		),
	)

	var raw []byte
	for _, s := range []string{"", "cpu", "nanoseconds", "size", "/bin/app", "abcd", "main.run", "main.go", "drop.*", "keep.*", "note", "https://example.com"} {
		raw = appendBytesField(raw, 6, []byte(s))
	}
	raw = appendBytesField(raw, 1, vtBytes)
	raw = appendBytesField(raw, 2, sampleBytes)
	raw = appendBytesField(raw, 3, mappingBytes)
	raw = appendBytesField(raw, 4, locBytes)
	raw = appendBytesField(raw, 5, fnBytes)
	raw = appendVarintField(raw, 7, 8)
	raw = appendVarintField(raw, 8, 9)
	raw = appendVarintField(raw, 9, 1000)
	raw = appendVarintField(raw, 10, 500)
	raw = appendBytesField(raw, 11, vtBytes)
	raw = appendVarintField(raw, 12, 10)
	raw = appendVarintField(raw, 13, 10)
	raw = appendBytesField(raw, 13, binary.AppendUvarint(nil, 10))
	raw = appendVarintField(raw, 14, 1)
	raw = appendVarintField(raw, 15, 11)
	raw = appendUnknownFields(raw)

	p, err := pprof.Parse(raw)
	require.NoError(t, err)
	assert.Equal(t, int64(8), p.DropFrames)
	assert.Equal(t, int64(9), p.KeepFrames)
	assert.Equal(t, int64(1), p.DefaultSampleType)
	assert.Equal(t, int64(11), p.DocURL)
	assert.Equal(t, []int64{10, 10}, p.Comment)
	require.Len(t, p.Mapping, 1)
	assert.True(t, p.Mapping[0].HasFunctions)
	assert.True(t, p.Mapping[0].HasFilenames)
	assert.True(t, p.Mapping[0].HasLineNumbers)
	assert.True(t, p.Mapping[0].HasInlineFrames)
	require.Len(t, p.Location, 1)
	assert.True(t, p.Location[0].IsFolded)
	require.Len(t, p.Location[0].Line, 1)
	assert.Equal(t, int64(34), p.Location[0].Line[0].Column)
	require.Len(t, p.Sample, 1)
	assert.Equal(t, []uint64{10}, p.Sample[0].LocationID)
	assert.Equal(t, []int64{99}, p.Sample[0].Value)
	require.Len(t, p.Sample[0].Label, 1)
	assert.Equal(t, int64(2), p.Sample[0].Label[0].NumUnit)
}

func TestParseValidationErrors(t *testing.T) {
	baseStrings := func() []byte {
		var b []byte
		b = appendBytesField(b, 6, []byte(""))
		b = appendBytesField(b, 6, []byte("cpu"))
		return b
	}

	cases := []struct {
		name    string
		raw     []byte
		wantErr string
	}{
		{
			name:    "empty_string_table_non_empty_profile",
			raw:     appendBytesField(nil, 1, appendVarintField(nil, 1, 0)),
			wantErr: "empty string table",
		},
		{
			name:    "non_empty_string_table_zero",
			raw:     appendBytesField(nil, 6, []byte("nonempty")),
			wantErr: "string_table[0] must be empty",
		},
		{
			name:    "sample_type_str_oob",
			raw:     appendBytesField(baseStrings(), 1, appendVarintField(nil, 1, 99)),
			wantErr: "out of range",
		},
		{
			name:    "header_str_negative",
			raw:     appendVarintField(baseStrings(), 7, ^uint64(0)),
			wantErr: "out of range",
		},
		{
			name:    "comment_str_oob",
			raw:     appendVarintField(baseStrings(), 13, 99),
			wantErr: "out of range",
		},
		{
			name:    "mapping_str_oob",
			raw:     appendBytesField(baseStrings(), 3, appendVarintField(nil, 5, 99)),
			wantErr: "out of range",
		},
		{
			name:    "function_str_oob",
			raw:     appendBytesField(baseStrings(), 5, appendVarintField(nil, 2, 99)),
			wantErr: "out of range",
		},
		{
			name:    "sample_value_len_mismatch",
			raw:     appendBytesField(appendBytesField(baseStrings(), 1, appendVarintField(nil, 1, 1)), 2, nil),
			wantErr: "has 0 values, want 1",
		},
		{
			name: "sample_label_str_oob",
			raw: appendBytesField(
				appendBytesField(baseStrings(), 1, appendVarintField(nil, 1, 1)),
				2,
				appendBytesField(
					appendVarintField(nil, 2, 1),
					3,
					appendVarintField(nil, 1, 99),
				),
			),
			wantErr: "out of range",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := pprof.Parse(tc.raw)
			assert.ErrorContains(t, err, tc.wantErr)
		})
	}
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

	profs, scopeAttrs := pprof.Convert(p, pprof.ConvertOptions{})
	require.Len(t, profs, 1)
	assert.Empty(t, scopeAttrs)
	require.Len(t, profs[0].Samples, 1)
	assert.Equal(t, []int64{10}, profs[0].Samples[0].Values)
	require.Len(t, profs[0].Samples[0].Stack, 1)
	assert.Len(t, profs[0].Samples[0].Stack[0].Lines, 1)
	assert.Equal(t, []profiledata.Attribute{
		{Key: semconv.PprofLocationIsFoldedKey, Value: attribute.BoolValue(true)},
	}, profs[0].Samples[0].Stack[0].Attributes)
	assert.Equal(t, []profiledata.Attribute{
		{Key: semconv.PprofProfileDropFramesKey, Value: attribute.StringValue("drop.*")},
		{Key: semconv.PprofProfileKeepFramesKey, Value: attribute.StringValue("keep.*")},
		{Key: semconv.PprofProfileDocURLKey, Value: attribute.StringValue("https://example.com")},
		{Key: semconv.PprofProfileCommentKey, Value: attribute.StringSliceValue([]string{"note1"})},
	}, profs[0].Attributes)
}

func TestConvertEdgeCasesAndInvariants(t *testing.T) {
	t.Run("nil_empty_unmatched_and_all_zero", func(t *testing.T) {
		profs, attrs := pprof.Convert(nil, pprof.ConvertOptions{})
		assert.Nil(t, profs)
		assert.Nil(t, attrs)

		p := &pprof.Profile{
			StringTable: []string{"", "cpu", "nanoseconds"},
			SampleType:  []pprof.ValueType{{Type: 1, Unit: 2}},
			Sample:      []pprof.Sample{{Value: []int64{0}}},
		}
		profs, attrs = pprof.Convert(p, pprof.ConvertOptions{SampleTypes: []string{"missing"}})
		assert.Nil(t, profs)
		assert.Nil(t, attrs)

		profs, attrs = pprof.Convert(p, pprof.ConvertOptions{})
		assert.Nil(t, profs)
		assert.Nil(t, attrs)
	})

	t.Run("sample_type_order_and_default_sample_type", func(t *testing.T) {
		p := &pprof.Profile{
			StringTable: []string{"", "a", "count", "b", "bytes", "c", "ns", "d"},
			SampleType: []pprof.ValueType{
				{Type: 1, Unit: 2}, // 0: a
				{Type: 3, Unit: 4}, // 1: b
				{Type: 5, Unit: 6}, // 2: c
				{Type: 7, Unit: 2}, // 3: d
			},
			DefaultSampleType: 7, // "d"
			Sample: []pprof.Sample{
				{Value: []int64{1, 0, 2, 3}},
			},
		}

		// Full profile with index 1 ("b") zero-valued -> emitted order is d, a, c -> dense [2, 0, 1].
		profs, scopeAttrs := pprof.Convert(p, pprof.ConvertOptions{})
		require.Len(t, profs, 3)
		assert.Equal(t, "d", profs[0].SampleType.Type)
		assert.Equal(t, "a", profs[1].SampleType.Type)
		assert.Equal(t, "c", profs[2].SampleType.Type)
		assert.Equal(t, []attribute.KeyValue{
			semconv.PprofScopeDefaultSampleType("d"),
			semconv.PprofScopeSampleTypeOrder(2, 0, 1),
		}, scopeAttrs)

		// Filtered subset excluding p.DefaultSampleType ("d") must not emit default_sample_type="d".
		filtered, filteredAttrs := pprof.Convert(p, pprof.ConvertOptions{
			SampleTypes:       []string{"a", "c"},
			DefaultSampleType: "c",
		})
		require.Len(t, filtered, 2)
		assert.Equal(t, "c", filtered[0].SampleType.Type)
		assert.Equal(t, "a", filtered[1].SampleType.Type)
		assert.Equal(t, []attribute.KeyValue{
			semconv.PprofScopeSampleTypeOrder(1, 0),
		}, filteredAttrs)
	})

	t.Run("mappings_zero_ids_and_large_inline_locations", func(t *testing.T) {
		lines := make([]pprof.Line, 12)
		for i := range lines {
			lines[i] = pprof.Line{FunctionID: 1, Line: int64(i%10 + 1)}
		}
		p := &pprof.Profile{
			StringTable: []string{
				"", "samples", "count",
				"/bin/valid", "deadbeef",
				`C:\app\svc.exe`, `C:\app\svc.exe2026-11-20 09:00:00 +0000 UTC`,
				"fn", "fn.go",
			},
			SampleType: []pprof.ValueType{{Type: 1, Unit: 2}},
			Mapping: []pprof.Mapping{
				{ID: 0, Filename: 3, BuildID: 4}, // Zero ID must be ignored.
				{
					ID:              1,
					MemoryStart:     0x2000,
					MemoryLimit:     0x1000, // Invalid range -> clamped to 0, 0.
					Filename:        3,
					BuildID:         4, // Valid hex GNU build ID fallback.
					HasFunctions:    true,
					HasFilenames:    true,
					HasLineNumbers:  true,
					HasInlineFrames: true,
				},
				{
					ID:          2,
					MemoryStart: 0x1000,
					MemoryLimit: 0x2000,
					Filename:    5,
					BuildID:     6, // Non-hex Windows PE BuildID must not become build_id.gnu.
				},
			},
			Function: []pprof.Function{
				{ID: 0, Name: 7}, // Zero ID must be ignored.
				{ID: 1, Name: 7, Filename: 8},
			},
			Location: []pprof.Location{
				{ID: 0, Address: 0x1}, // Zero ID must be ignored.
				{ID: 1, MappingID: 1, Address: 0x10, Line: lines, IsFolded: true},
				{ID: 2, MappingID: 2, Address: 0x20},
				{ID: 3, MappingID: 0, Address: 0x30, Line: []pprof.Line{{FunctionID: 0, Line: 5}}},
			},
			Sample: []pprof.Sample{
				{LocationID: []uint64{1, 2, 3}, Value: []int64{5}},
			},
		}

		profs, _ := pprof.Convert(p, pprof.ConvertOptions{FrameType: semconv.ProfileFrameTypeGo})
		require.Len(t, profs, 1)
		stack := profs[0].Samples[0].Stack
		require.Len(t, stack, 3)

		// Location 1: clamped mapping range, GNU hex fallback + mapping bool attrs, 10 deduplicated lines.
		require.NotNil(t, stack[0].Mapping)
		assert.Zero(t, stack[0].Mapping.MemoryStart)
		assert.Zero(t, stack[0].Mapping.MemoryLimit)
		assert.Contains(t, stack[0].Mapping.Attributes, profiledata.Attribute{
			Key:   semconv.ProcessExecutableBuildIDGNUKey,
			Value: attribute.StringValue("deadbeef"),
		})
		assert.Contains(t, stack[0].Mapping.Attributes, profiledata.Attribute{
			Key:   semconv.PprofMappingHasInlineFramesKey,
			Value: attribute.BoolValue(true),
		})
		assert.Len(t, stack[0].Lines, 10)
		assert.Len(t, stack[0].Attributes, 2)

		// Location 2: Windows PE BuildID must NOT be tagged as build_id.gnu.
		require.NotNil(t, stack[1].Mapping)
		assert.Empty(t, stack[1].Mapping.Attributes)

		// Location 3: unset MappingID/FunctionID must not attach zero-ID entries.
		assert.Nil(t, stack[2].Mapping)
		require.Len(t, stack[2].Lines, 1)
		assert.Empty(t, stack[2].Lines[0].Function.Name)
	})

	t.Run("duplicate_trace_labels_sorting_and_key_collision_safety", func(t *testing.T) {
		const (
			traceHex1 = "0102030405060708090a0b0c0d0e0f10"
			traceHex2 = "1112131415161718191a1b1c1d1e1f20"
			spanHex   = "0102030405060708"
		)
		p := &pprof.Profile{
			StringTable: []string{
				"", "cpu", "ns",
				"otel.trace_id", traceHex1, traceHex2, "not-hex",
				"otel.span_id", spanHex,
				"zeta", "alpha", "bytes",
				"a=STRING:1", "2", "a", "1=STRING:2",
			},
			SampleType: []pprof.ValueType{{Type: 1, Unit: 2}},
			Location:   []pprof.Location{{ID: 1, Address: 0x10}},
			Sample: []pprof.Sample{
				// Duplicate valid otel.trace_id: last wins for Link and earlier value must NOT leak into Attributes.
				{
					LocationID: []uint64{1},
					Value:      []int64{10},
					Label: []pprof.Label{
						{Key: 3, Str: 4},
						{Key: 3, Str: 5},
						{Key: 7, Str: 8},
						{Key: 9, Num: 64, NumUnit: 11},
						{Key: 10, Str: 9},
						{Key: 0, Str: 9}, // Empty key skipped.
					},
				},
				// Valid trace_id overridden by invalid trace_id: last-wins means NO Link and invalid value kept.
				{
					LocationID: []uint64{1},
					Value:      []int64{20},
					Label: []pprof.Label{
						{Key: 3, Str: 4},
						{Key: 3, Str: 6},
						{Key: 7, Str: 8},
					},
				},
				// Delimiter-collision attribute pair: must NOT merge into one sample.
				{
					LocationID: []uint64{1},
					Value:      []int64{30},
					Label:      []pprof.Label{{Key: 12, Str: 13}},
				},
				{
					LocationID: []uint64{1},
					Value:      []int64{40},
					Label:      []pprof.Label{{Key: 14, Str: 15}},
				},
			},
		}

		profs, _ := pprof.Convert(p, pprof.ConvertOptions{})
		require.Len(t, profs, 1)
		require.Len(t, profs[0].Samples, 4)

		wantTID2, err := trace.TraceIDFromHex(traceHex2)
		require.NoError(t, err)
		wantSID, err := trace.SpanIDFromHex(spanHex)
		require.NoError(t, err)

		// Sample 0: Link uses traceHex2, no otel.trace_id/span_id leaked, attributes sorted by key.
		s0 := profs[0].Samples[0]
		assert.Equal(t, profiledata.Link{TraceID: wantTID2, SpanID: wantSID}, s0.Link)
		assert.Equal(t, []profiledata.Attribute{
			{Key: "alpha", Value: attribute.StringValue("zeta")},
			{Key: "zeta", Value: attribute.Int64Value(64), Unit: "bytes"},
		}, s0.Attributes)

		// Sample 1: Last-wins invalid trace_id prevents Link creation.
		s1 := profs[0].Samples[1]
		assert.False(t, s1.Link.IsValid())
		assert.Equal(t, []profiledata.Attribute{
			{Key: "otel.span_id", Value: attribute.StringValue(spanHex)},
			{Key: "otel.trace_id", Value: attribute.StringValue("not-hex")},
		}, s1.Attributes)

		// Samples 2 & 3 stayed distinct despite delimiter characters.
		assert.Equal(t, []int64{30}, profs[0].Samples[2].Values)
		assert.Equal(t, []int64{40}, profs[0].Samples[3].Values)
	})
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
