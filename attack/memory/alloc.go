package memory

import (
	"bytes"
	"context"
	"errors"
	"runtime/debug"

	"github.com/slok/ragnarok/log"
)

const (
	// KiB refers to kibibyte.
	KiB = 1 << (10 * (iota + 1))
	// MiB refers to mebibyte.
	MiB = 1 << (10 * (iota + 1))
	// GiB refers to gibibyte.
	GiB = 1 << (10 * (iota + 1))
	// TiB refers to tebibyte.
	TiB = 1 << (10 * (iota + 1))
	// PiB refers to pebibyte.
	PiB = 1 << (10 * (iota + 1))
	// EiB refers to exbibyte.
	EiB = 1 << (10 * (iota + 1))
)

// MemAllocation failer will apply a failure allocating memory.
type MemAllocation struct {
	b    bytes.Buffer //  Structure that will allocate memory.
	size uint64       // The size of the buffer in bytes.
	log  log.Logger   // Logger.
	done bool         // Flag that marks the allocation was already made.
}

// NewMemAllocation returns a new default memory allocation failer.
func NewMemAllocation(size uint64) *MemAllocation {
	return &MemAllocation{
		size: size,
		log:  log.Base(),
	}
}

// Apply will allocate and use the memory.
func (m *MemAllocation) Apply(ctx context.Context) error {
	// Check if its done before continuing.
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	if m.done {
		return errors.New("memory allocation already applied")
	}
	log.With("MiB", m.size/MiB).Infof("Allocated memory")
	m.b.Write(make([]byte, m.size))
	m.done = true
	return nil
}

// Revert will clean the memory usage letting the GC do its work.
func (m *MemAllocation) Revert() error {
	// Reset our buffer to be able to free the memory.
	m.b.Reset()
	m.b = bytes.Buffer{}

	// Return memory to the OS (Way better to do?)
	debug.FreeOSMemory()
	log.With("MiB", m.size/MiB).Infof("Reverted allocated memory")
	m.done = false
	return nil
}
