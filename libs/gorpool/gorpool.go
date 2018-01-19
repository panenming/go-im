package gorpool

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrPoolNotRunning = errors.New("the pool is not running")
	ErrJobNotFunc     = errors.New("generic worker not given a func()")
	ErrWorkerClosed   = errors.New("worker was closed")
	ErrJobTimeOut     = errors.New("job request time out")
)

// 执行的worker定义
type Worker interface {
	// 执行一个job ，返回结果
	Process(interface{}) interface{}
	//
	BlockUntilReady()
	Interrupt()
	Terminate()
}

type closureWorker struct {
	processor func(interface{}) interface{}
}

func (w *closureWorker) Process(payload interface{}) interface{} {
	return w.processor(payload)
}

func (w *closureWorker) BlockUntilReady() {}
func (w *closureWorker) Interrupt()       {}
func (w *closureWorker) Terminate()       {}

type callbackWorker struct{}

func (w *callbackWorker) Process(payload interface{}) interface{} {
	f, ok := payload.(func())
	if !ok {
		return ErrJobNotFunc
	}
	f()
	return nil
}

func (w *callbackWorker) BlockUntilReady() {}
func (w *callbackWorker) Interrupt()       {}
func (w *callbackWorker) Terminate()       {}

type Pool struct {
	ctor      func() Worker
	workers   []*workerWrapper
	reqChan   chan workRequest
	workerMut sync.Mutex
}

func (p *Pool) SetSize(n int) {
	p.workerMut.Lock()
	defer p.workerMut.Unlock()

	lWorkers := len(p.workers)
	if lWorkers == n {
		return
	}

	// Add extra workers if N > len(workers)
	for i := lWorkers; i < n; i++ {
		p.workers = append(p.workers, newWorkerWrapper(p.reqChan, p.ctor()))
	}

	// Asynchronously stop all workers > N
	for i := n; i < lWorkers; i++ {
		p.workers[i].stop()
	}

	// Synchronously wait for all workers > N to stop
	for i := n; i < lWorkers; i++ {
		p.workers[i].join()
	}

	// Remove stopped workers from slice
	p.workers = p.workers[:n]
}

// GetSize returns the current size of the pool.
func (p *Pool) GetSize() int {
	p.workerMut.Lock()
	defer p.workerMut.Unlock()

	return len(p.workers)
}

// Close will terminate all workers and close the job channel of this Pool.
func (p *Pool) Close() {
	p.SetSize(0)
	close(p.reqChan)
}

func New(n int, ctor func() Worker) *Pool {
	p := &Pool{
		ctor:    ctor,
		reqChan: make(chan workRequest),
	}
	p.SetSize(n)
	return p
}

func NewFunc(n int, f func(interface{}) interface{}) *Pool {
	return New(n, func() Worker {
		return &closureWorker{
			processor: f,
		}
	})
}

func NewCallback(n int) *Pool {
	return New(n, func() Worker {
		return &callbackWorker{}
	})
}

func (p *Pool) Process(payload interface{}) interface{} {
	request, open := <-p.reqChan
	if !open {
		panic(ErrPoolNotRunning)
	}

	request.jobChan <- payload
	payload, open = <-request.retChan
	if !open {
		panic(ErrWorkerClosed)
	}
	return payload
}

func (p *Pool) ProcessTimed(
	payload interface{},
	timeout time.Duration,
) (interface{}, error) {
	tout := time.NewTimer(timeout)

	var request workRequest
	var open bool

	select {
	case request, open = <-p.reqChan:
		if !open {
			return nil, ErrPoolNotRunning
		}
	case <-tout.C:
		return nil, ErrJobTimeOut
	}

	select {
	case request.jobChan <- payload:
	case <-tout.C:
		request.interruptFunc()
		return nil, ErrJobTimeOut
	}

	select {
	case payload, open = <-request.retChan:
		if !open {
			return nil, ErrWorkerClosed
		}
	case <-tout.C:
		request.interruptFunc()
		return nil, ErrJobTimeOut
	}

	tout.Stop()
	return payload, nil
}
