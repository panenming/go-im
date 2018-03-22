package concurrent

import (
	"sync/atomic"
)

type AtomicInteger int32

func (i *AtomicInteger) IncrementAndGet() int32 {
	return atomic.AddInt32((*int32)(i), int32(1))
}

func (i *AtomicInteger) GetAndIncrement() int32 {
	ret := atomic.LoadInt32((*int32)(i))
	atomic.AddInt32((*int32)(i), int32(1))
	return ret
}

func (i *AtomicInteger) DecrementAndGet() int32 {
	return atomic.AddInt32((*int32)(i), int32(-1))
}

// GetAndDecrement decrement wrapped int32 with 1 and return old value.
func (i *AtomicInteger) GetAndDecrement() int32 {
	ret := atomic.LoadInt32((*int32)(i))
	atomic.AddInt32((*int32)(i), int32(-1))
	return ret
}

// Get current value
func (i *AtomicInteger) Get() int32 {
	return atomic.LoadInt32((*int32)(i))
}
