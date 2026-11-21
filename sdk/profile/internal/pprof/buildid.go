// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pprof

import (
	"bytes"
	"crypto/sha256"
	"debug/elf"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	htlHashChunkSize = 4096
	maxBuildIDCache  = 256

	// ELF note types for GNU and Go build IDs.
	elfNoteTypeGNUBuildID = 3
	elfNoteTypeGoBuildID  = 4
)

// BuildIDs holds binary identifiers extracted for a mapped executable or library.
type BuildIDs struct {
	HTLHash    string
	GoBuildID  string
	GNUBuildID string
}

// BuildIDCache caches computed build IDs keyed by mapped file path so each
// executable is inspected at most once during the process lifetime.
type BuildIDCache struct {
	mu      sync.Mutex
	entries map[string]BuildIDs
}

// NewBuildIDCache returns a new empty BuildIDCache.
func NewBuildIDCache() *BuildIDCache {
	return &BuildIDCache{
		entries: make(map[string]BuildIDs),
	}
}

// Lookup returns cached build IDs for filename, computing them on first access.
// Non-file or virtual paths (such as "[vdso]") return zero BuildIDs.
func (c *BuildIDCache) Lookup(filename string) BuildIDs {
	if c == nil || !isFilePath(filename) {
		return BuildIDs{}
	}

	c.mu.Lock()
	if ids, ok := c.entries[filename]; ok {
		c.mu.Unlock()
		return ids
	}
	c.mu.Unlock()

	ids := inspectBinary(filename)

	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= maxBuildIDCache {
		clear(c.entries)
	}
	c.entries[filename] = ids
	return ids
}

func isFilePath(filename string) bool {
	if filename == "" || strings.HasPrefix(filename, "[") {
		return false
	}
	return filepath.IsAbs(filename)
}

func inspectBinary(filename string) BuildIDs {
	f, err := os.Open(filename) //nolint:gosec // Reads local mapped binary paths reported by runtime/pprof.
	if err != nil {
		return BuildIDs{}
	}
	defer func() { _ = f.Close() }()

	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return BuildIDs{}
	}

	var ids BuildIDs
	if h, err := ComputeHTLHash(f, info.Size()); err == nil {
		ids.HTLHash = h
	}
	ids.GoBuildID, ids.GNUBuildID = ExtractELFBuildIDs(f)
	return ids
}

// ComputeHTLHash computes the OpenTelemetry process.executable.build_id.htlhash
// digest for a binary of the given size.
//
// Per the OpenTelemetry Profiles mappings specification:
//
//	Input   = Concat(File[:4096], File[-4096:], BigEndianUInt64(Len(File)))
//	Digest  = SHA256(Input)
//	BuildID = hex(Digest[:16])
//
// When the file is smaller than 4096 bytes, both slices cover the full file
// without padding.
func ComputeHTLHash(r io.ReaderAt, size int64) (string, error) {
	if size < 0 {
		return "", io.ErrUnexpectedEOF
	}
	chunkLen := min(size, int64(htlHashChunkSize))
	buf := make([]byte, chunkLen)

	h := sha256.New()
	if chunkLen > 0 {
		if err := readFullAt(r, buf, 0); err != nil {
			return "", err
		}
		_, _ = h.Write(buf)

		tailOff := size - chunkLen
		if tailOff > 0 {
			if err := readFullAt(r, buf, tailOff); err != nil {
				return "", err
			}
		}
		_, _ = h.Write(buf)
	}

	var lenBuf [8]byte
	binary.BigEndian.PutUint64(lenBuf[:], uint64(size)) //nolint:gosec // Non-negative size verified above.
	_, _ = h.Write(lenBuf[:])

	var sum [sha256.Size]byte
	h.Sum(sum[:0])
	return hex.EncodeToString(sum[:16]), nil
}

func readFullAt(r io.ReaderAt, buf []byte, off int64) error {
	n, err := r.ReadAt(buf, off)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if n != len(buf) {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ExtractELFBuildIDs extracts Go and GNU build IDs from ELF SHT_NOTE sections.
// For non-ELF files it returns empty strings.
func ExtractELFBuildIDs(r io.ReaderAt) (goBuildID, gnuBuildID string) {
	ef, err := elf.NewFile(r)
	if err != nil {
		return "", ""
	}
	defer func() { _ = ef.Close() }()

	for _, sec := range ef.Sections {
		if sec.Type != elf.SHT_NOTE {
			continue
		}
		// Skip unexpectedly large note sections to avoid excessive memory usage.
		if sec.Size == 0 || sec.Size > 64*1024 {
			continue
		}
		data, err := sec.Data()
		if err != nil {
			continue
		}
		parseELFNotes(data, ef.ByteOrder, &goBuildID, &gnuBuildID)
		if goBuildID != "" && gnuBuildID != "" {
			break
		}
	}
	return goBuildID, gnuBuildID
}

func parseELFNotes(data []byte, order binary.ByteOrder, goBuildID, gnuBuildID *string) {
	const headerSize = 12
	off := 0
	for off+headerSize <= len(data) {
		nameSz := order.Uint32(data[off : off+4])
		descSz := order.Uint32(data[off+4 : off+8])
		noteType := order.Uint32(data[off+8 : off+12])
		off += headerSize

		nameAligned := align4(nameSz)
		descAligned := align4(descSz)
		if uint64(off)+nameAligned+descAligned > uint64(len(data)) {
			break
		}

		name := bytes.TrimRight(data[off:off+int(nameSz)], "\x00")
		descStart := off + int(nameAligned) //nolint:gosec // Bounded by len(data) above.
		desc := data[descStart : descStart+int(descSz)]
		off = descStart + int(descAligned) //nolint:gosec // Bounded by len(data) above.

		switch {
		case *goBuildID == "" && noteType == elfNoteTypeGoBuildID && bytes.Equal(name, []byte("Go")):
			*goBuildID = string(desc)
		case *gnuBuildID == "" && noteType == elfNoteTypeGNUBuildID && bytes.Equal(name, []byte("GNU")):
			*gnuBuildID = hex.EncodeToString(desc)
		}
	}
}

func align4(n uint32) uint64 {
	return (uint64(n) + 3) &^ 3
}
