// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "cpu" namespace.
package cpuconv

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
)

// ModeAttr is an attribute conforming to the cpu.mode semantic conventions. It
// represents the mode of the CPU.
type ModeAttr string

var (
	// ModeUser is the none.
	ModeUser ModeAttr = "user"
	// ModeSystem is the none.
	ModeSystem ModeAttr = "system"
	// ModeNice is the none.
	ModeNice ModeAttr = "nice"
	// ModeIdle is the none.
	ModeIdle ModeAttr = "idle"
	// ModeIOWait is the none.
	ModeIOWait ModeAttr = "iowait"
	// ModeInterrupt is the none.
	ModeInterrupt ModeAttr = "interrupt"
	// ModeSteal is the none.
	ModeSteal ModeAttr = "steal"
	// ModeKernel is the none.
	ModeKernel ModeAttr = "kernel"
)