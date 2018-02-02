package bytes

import (
	"testing"
)

func TestBuffer(t *testing.T) {
	p := NewPool(2000, 10)
	b := p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}

	b = p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
	b = p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}

	p.Put(b)

	for i := 1; i <= 100000; i++ {
		go func() {
			b = p.Get()
			if b.Bytes() == nil || len(b.Bytes()) == 0 {
				t.FailNow()
			}
		}()
	}

}
