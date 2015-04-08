package bimg

import (
	"github.com/dustin/go-humanize"
	. "github.com/tj/go-debug"
	"runtime"
	"strconv"
	"time"
)

var debug = Debug("bimg")

// Print Go memory and garbage collector stats. Useful for debugging
func PrintMemoryStats() {
	log := Debug("memory")
	mem := memoryStats()

	log("\u001b[33m---- Memory Dump Stats ----\u001b[39m")
	log("Allocated: %s", humanize.Bytes(mem.Alloc))
	log("Total Allocated: %s", humanize.Bytes(mem.TotalAlloc))
	log("Memory Allocations: %d", mem.Mallocs)
	log("Memory Frees: %d", mem.Frees)
	log("Heap Allocated: %s", humanize.Bytes(mem.HeapAlloc))
	log("Heap System: %s", humanize.Bytes(mem.HeapSys))
	log("Heap In Use: %s", humanize.Bytes(mem.HeapInuse))
	log("Heap Idle: %s", humanize.Bytes(mem.HeapIdle))
	log("Heap OS Related: %s", humanize.Bytes(mem.HeapReleased))
	log("Heap Objects: %s", humanize.Bytes(mem.HeapObjects))
	log("Stack In Use: %s", humanize.Bytes(mem.StackInuse))
	log("Stack System: %s", humanize.Bytes(mem.StackSys))
	log("Stack Span In Use: %s", humanize.Bytes(mem.MSpanInuse))
	log("Stack Cache In Use: %s", humanize.Bytes(mem.MCacheInuse))
	log("Next GC cycle: %s", humanizeNano(mem.NextGC))
	log("Last GC cycle: %s", humanize.Time(time.Unix(0, int64(mem.LastGC))))
	log("\u001b[33m---- End Memory Dump ----\u001b[39m")
}

func memoryStats() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem
}

func humanizeNano(n uint64) string {
	var suffix string

	switch {
	case n > 1e9:
		n /= 1e9
		suffix = "s"
	case n > 1e6:
		n /= 1e6
		suffix = "ms"
	case n > 1e3:
		n /= 1e3
		suffix = "us"
	default:
		suffix = "ns"
	}

	return strconv.Itoa(int(n)) + suffix
}
