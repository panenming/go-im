package disruptor

import (
	"testing"
)

func BenchmarkCursorStore(b *testing.B) {
	interations := int64(b.N)

	cursor := NewCursor()
	b.ReportAllocs()
	b.ResetTimer()
	for i := int64(0); i < interations; i++ {
		cursor.Store(i)
	}
}

func BenchmarkCursorLoad(b *testing.B) {
	iterations := int64(b.N)

	cursor := NewCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		cursor.Load()
	}
}
func BenchmarkCursorRead(b *testing.B) {
	iterations := int64(b.N)

	cursor := NewCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		cursor.Read(i)
	}
}

func BenchmarkCursorReadAsBarrier(b *testing.B) {
	var barrier Barrier = NewCursor()

	interations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()
	for i := int64(0); i < interations; i++ {
		barrier.Read(0)
	}
}
