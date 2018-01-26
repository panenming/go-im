package disruptor

import (
	"runtime"
)

const SpinMask = 1024*16 - 1

type Writer struct {
	written  *Cursor
	upstream Barrier
	capacity int64
	previous int64
	gate     int64
}

func NewWriter(written *Cursor, upstream Barrier, capacity int64) *Writer {
	assertPowerOfTwo(capacity)
	return &Writer{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: InitialSequenceValue,
		gate:     InitialSequenceValue,
	}
}

func assertPowerOfTwo(value int64) {
	if value > 0 && (value&(value-1)) != 0 {
		panic("The ring capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}
}

func (w *Writer) Reserve(count int64) int64 {
	w.previous += count
	for spin := int64(0); w.previous-w.capacity > w.gate; spin++ {
		if spin&SpinMask == 0 {
			runtime.Gosched()
		}
		w.gate = w.upstream.Read(0)
	}
	return w.previous
}

func (w *Writer) Await(next int64) {
	for next-w.capacity > w.gate {
		w.gate = w.upstream.Read(0)
	}
}
