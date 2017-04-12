// +build integration

package memory

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryAllocationAttack(t *testing.T) {
	assert := assert.New(t)
	var size uint64 = 200 * MiB

	ma := NewMemAllocation(size)
	// Get current memory
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	startMem := mem.Alloc

	// Allocate memory and test if increased.
	ma.Apply(context.TODO())
	time.Sleep(1 * time.Millisecond)
	runtime.ReadMemStats(&mem)
	endMem := mem.Alloc

	// Let 5% margin delta from the wanted size
	assert.InDelta((endMem - startMem), size, float64(size)*0.05, "current memory allocation should be wanted allocation (5% deviation)")
	// Free memory and test if released.
	ma.Revert()
	time.Sleep(1 * time.Millisecond)
	runtime.ReadMemStats(&mem)

	// Let 5% margin delta from the wanted size
	assert.InDelta(startMem, mem.Alloc, float64(size)*0.05, "current memory and initial memory should be equal (5% deviation)")
}
