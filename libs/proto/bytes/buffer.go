package bytes

import (
	"sync"
)

type Buffer struct {
	buf  []byte
	next *Buffer
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}

// buf pool
type Pool struct {
	lock sync.Mutex
	free *Buffer
	max  int
	num  int // 每次grow的大小，不代表此pool最终的大小
	size int
}

func NewPool(num, size int) (p *Pool) {
	p = new(Pool)
	p.Init(num, size)
	return
}

func (p *Pool) Init(num, size int) {
	p.init(num, size)
}

func (p *Pool) init(num, size int) {
	p.num = num
	p.size = size
	p.max = num * size
	p.grow()
}

// 初始化pool
func (p *Pool) grow() {
	var (
		i   int
		b   *Buffer
		bs  []Buffer
		buf []byte
	)
	buf = make([]byte, p.max)
	bs = make([]Buffer, p.num)
	p.free = &bs[0]
	b = p.free
	for i = 1; i < p.num; i++ {
		b.buf = buf[(i-1)*p.size : i*p.size]
		b.next = &bs[i]
		b = b.next
	}
	b.buf = buf[(i-1)*p.size : i*p.size]
	b.next = nil
	return
}

// 取一个未使用memery buf
func (p *Pool) Get() (b *Buffer) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if b = p.free; b == nil {
		p.grow()
		b = p.free
	}
	p.free = b.next
	return
}

// 还回一个使用结束的memery buf
func (p *Pool) Put(b *Buffer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	b.next = p.free
	p.free = b
	return
}
