package disruptor

import (
	"runtime"
	"time"
)

type Reader struct {
	read     *Cursor
	written  *Cursor
	upstream Barrier
	consumer Consumer
	ready    bool
}

func NewReader(read, written *Cursor, upstream Barrier, consumer Consumer) *Reader {
	return &Reader{
		read:     read,
		written:  written,
		upstream: upstream,
		consumer: consumer,
		ready:    false,
	}
}

func (r *Reader) Start() {
	r.ready = true
	go r.receive()
}

func (this *Reader) Stop() {
	this.ready = false
}

func (r *Reader) receive() {
	previous := r.read.Load()
	idling, gating := 0, 0

	for {
		lower := previous + 1
		upper := r.upstream.Read(lower)

		if lower <= upper {
			r.consumer.Consume(lower, upper)
			r.read.Store(upper)
			previous = upper
		} else if upper = r.written.Load(); lower <= upper {
			time.Sleep(time.Microsecond)
			gating++
			idling = 0
		} else if r.ready {
			time.Sleep(time.Millisecond)
			idling++
			gating = 0
		} else {
			break
		}

		runtime.Gosched()
	}
}
