package collections

import (
	"errors"
	"sync"
	"time"

	"github.com/panenming/go-im/libs/commonspool/concurrent"
)

type InterruptedErr struct{}

func NewInterruptedErr() *InterruptedErr {
	return &InterruptedErr{}
}

func (err *InterruptedErr) Error() string {
	return "Interrupted"
}

type Node struct {
	item interface{}
	prev *Node
	next *Node
}

func newNode(item interface{}, prev, next *Node) *Node {
	return &Node{
		item: item,
		prev: prev,
		next: next,
	}
}

type LinkedBlockingDeque struct {
	first    *Node
	last     *Node
	count    int
	capacity int
	lock     *sync.Mutex
	notEmpty *concurrent.TimeoutCond
	notFull  *concurrent.TimeoutCond
}

func NewDeque(capacity int) *LinkedBlockingDeque {
	if capacity < 0 {
		panic(errors.New("capacity must > 0"))
	}
	lock := new(sync.Mutex)
	return &LinkedBlockingDeque{capacity: capacity, lock: lock, notEmpty: concurrent.NewTimeoutCond(lock), notFull: concurrent.NewTimeoutCond(lock)}
}

func (q *LinkedBlockingDeque) linkFirst(e interface{}) bool {
	if q.count >= q.capacity {
		return false
	}
	f := q.first
	x := newNode(e, nil, f)
	q.first = x
	if q.last == nil {
		q.last = x
	} else {
		f.prev = x
	}
	q.count++
	q.notEmpty.Signal()
	return true
}

func (q *LinkedBlockingDeque) linkLast(e interface{}) bool {
	if q.count >= q.capacity {
		return false
	}
	l := q.last
	x := newNode(e, l, nil)
	q.last = x
	if q.first == nil {
		q.first = x
	} else {
		l.next = x
	}
	q.count++
	q.notEmpty.Signal()
	return true
}

func (q *LinkedBlockingDeque) unlinkFirst() interface{} {
	f := q.first
	if f == nil {
		return nil
	}
	n := f.next
	item := f.item
	f.item = nil
	f.next = f // help gc
	q.first = n
	if n == nil {
		q.last = nil
	} else {
		n.prev = nil
	}
	q.count--
	q.notFull.Signal()
	return item
}

func (q *LinkedBlockingDeque) unlinkLast() interface{} {
	l := q.last
	if l == nil {
		return nil
	}
	p := l.prev
	item := l.item
	l.item = nil
	l.prev = l // help GC
	q.last = p
	if p == nil {
		q.first = nil
	} else {
		p.next = nil
	}
	q.count = q.count - 1
	q.notFull.Signal()
	return item
}

func (q *LinkedBlockingDeque) unlink(x *Node) {
	// assert lock.isHeldByCurrentThread();
	p := x.prev
	n := x.next
	if p == nil {
		q.unlinkFirst()
	} else if n == nil {
		q.unlinkLast()
	} else {
		p.next = n
		n.prev = p
		x.item = nil
		// Don't mess with x's links.  They may still be in use by
		// an iterator.
		q.count = q.count - 1
		q.notFull.Signal()
	}
}

func (q *LinkedBlockingDeque) AddFirst(e interface{}) error {
	if e == nil {
		return errors.New("e is nil")
	}
	if !q.OfferFirst(e) {
		return errors.New("Deque full")
	}
	return nil
}

func (q *LinkedBlockingDeque) AddLast(e interface{}) error {
	if e == nil {
		return errors.New("e is nil")
	}
	if !q.OfferLast(e) {
		return errors.New("Deque full")
	}
	return nil
}

// OfferFirst inserts the specified element at the front of this deque unless it would violate capacity restrictions.
// return if the element was added to this deque
func (q *LinkedBlockingDeque) OfferFirst(e interface{}) bool {
	if e == nil {
		return false
	}
	q.lock.Lock()
	result := q.linkFirst(e)
	q.lock.Unlock()
	return result
}

// OfferLast inserts the specified element at the end of this deque unless it would violate capacity restrictions.
// return if the element was added to this deque
func (q *LinkedBlockingDeque) OfferLast(e interface{}) bool {
	if e == nil {
		return false
	}
	q.lock.Lock()
	result := q.linkLast(e)
	q.lock.Unlock()
	return result
}

// PutFirst link the provided element as the first in the queue, waiting until there
// is space to do so if the queue is full.
func (q *LinkedBlockingDeque) PutFirst(e interface{}) {
	if e == nil {
		return
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	for !q.linkFirst(e) {
		q.notFull.Wait()
	}
}

// PutLast Link the provided element as the last in the queue, waiting until there
// is space to do so if the queue is full.
func (q *LinkedBlockingDeque) PutLast(e interface{}) {
	if e == nil {
		return
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	for !q.linkLast(e) {
		q.notFull.Wait()
	}
}

// PollFirst retrieves and removes the first element of this deque,
// or returns nil if this deque is empty.
func (q *LinkedBlockingDeque) PollFirst() (e interface{}) {
	q.lock.Lock()
	result := q.unlinkFirst()
	q.lock.Unlock()
	return result
}

// PollFirstWithTimeout retrieves and removes the first element of this deque, waiting
// up to the specified wait time if necessary for an element to become available.
// return NewInterruptedErr when waiting bean interrupted
func (q *LinkedBlockingDeque) PollFirstWithTimeout(timeout time.Duration) (interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	var x interface{}
	interrupt := false
	for x = q.unlinkFirst(); x == nil; x = q.unlinkFirst() {
		if timeout <= 0 {
			break
		}
		if interrupt {
			return nil, NewInterruptedErr()
		}
		timeout, interrupt = q.notEmpty.WaitWithTimeout(timeout)
	}
	return x, nil
}

// PollLast retrieves and removes the last element of this deque,
// or returns nil if this deque is empty.
func (q *LinkedBlockingDeque) PollLast() interface{} {
	q.lock.Lock()
	result := q.unlinkLast()
	q.lock.Unlock()
	return result
}

// PollLastWithTimeout retrieves and removes the last element of this deque, waiting
// up to the specified wait time if necessary for an element to become available.
// return NewInterruptedErr when waiting bean interrupted
func (q *LinkedBlockingDeque) PollLastWithTimeout(timeout time.Duration) (interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	var x interface{}
	interrupt := false
	for x = q.unlinkLast(); x == nil; x = q.unlinkLast() {
		if timeout <= 0 {
			break
		}
		if interrupt {
			return nil, NewInterruptedErr()
		}
		timeout, interrupt = q.notEmpty.WaitWithTimeout(timeout)
	}
	return x, nil
}

// TakeFirst unlink the first element in the queue, waiting until there is an element
// to unlink if the queue is empty.
// return NewInterruptedErr if wait condition is interrupted
func (q *LinkedBlockingDeque) TakeFirst() (interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	var x interface{}
	interrupt := false
	for x = q.unlinkFirst(); x == nil; x = q.unlinkFirst() {
		if interrupt {
			return nil, NewInterruptedErr()
		}
		interrupt = q.notEmpty.Wait()
	}
	return x, nil
}

// TakeLast unlink the last element in the queue, waiting until there is an element
// to unlink if the queue is empty.
// return NewInterruptedErr if wait condition is interrupted
func (q *LinkedBlockingDeque) TakeLast() (interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	var x interface{}
	interrupt := false
	for x = q.unlinkLast(); x == nil; x = q.unlinkLast() {
		if interrupt {
			return nil, NewInterruptedErr()
		}
		interrupt = q.notEmpty.Wait()
	}
	return x, nil
}

// PeekFirst retrieves, but does not remove, the first element of this deque,
// or returns nil if this deque is empty.
func (q *LinkedBlockingDeque) PeekFirst() interface{} {
	var result interface{}
	q.lock.Lock()
	if q.first == nil {
		result = nil
	} else {
		result = q.first.item
	}
	q.lock.Unlock()
	return result
}

// PeekLast retrieves, but does not remove, the last element of this deque,
// or returns nil if this deque is empty.
func (q *LinkedBlockingDeque) PeekLast() interface{} {
	var result interface{}
	q.lock.Lock()
	if q.last == nil {
		result = nil
	} else {
		result = q.last.item
	}
	q.lock.Unlock()
	return result
}

// RemoveFirstOccurrence removes the first occurrence of the specified element from this deque.
// If the deque does not contain the element, it is unchanged.
// More formally, removes the first element item such that
//		o == item
// (if such an element exists).
// Returns true if this deque contained the specified element
// (or equivalently, if this deque changed as a result of the call).
func (q *LinkedBlockingDeque) RemoveFirstOccurrence(item interface{}) bool {
	if item == nil {
		return false
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	for p := q.first; p != nil; p = p.next {
		if item == p.item {
			q.unlink(p)
			return true
		}
	}
	return false
}

// RemoveLastOccurrence removes the last occurrence of the specified element from this deque.
// If the deque does not contain the element, it is unchanged.
// More formally, removes the last element item such that
//		o == item
// (if such an element exists).
// Returns true if this deque contained the specified element
// (or equivalently, if this deque changed as a result of the call).
func (q *LinkedBlockingDeque) RemoveLastOccurrence(item interface{}) bool {
	if item == nil {
		return false
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	for p := q.last; p != nil; p = p.prev {
		if item == p.item {
			q.unlink(p)
			return true
		}
	}
	return false
}

// InterruptTakeWaiters interrupts the goroutine currently waiting to take an object from the pool.
func (q *LinkedBlockingDeque) InterruptTakeWaiters() {
	q.notEmpty.Interrupt()
}

// HasTakeWaiters returns true if there are goroutine waiting to take instances from this deque.
// See disclaimer on accuracy in  TimeoutCond.HasWaiters()
func (q *LinkedBlockingDeque) HasTakeWaiters() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.notEmpty.HasWaiters()
}

// ToSlice returns an slice containing all of the elements in this deque, in
// proper sequence (from first to last element).
func (q *LinkedBlockingDeque) ToSlice() []interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	a := make([]interface{}, q.count)
	for p, k := q.first, 0; p != nil; p, k = p.next, k+1 {
		a[k] = p.item
	}
	return a
}

// Size return this LinkedBlockingDeque current elements len, is concurrent safe
func (q *LinkedBlockingDeque) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.size()
}

func (q *LinkedBlockingDeque) size() int {
	return q.count
}

// Iterator return a asc iterator of this deque
func (q *LinkedBlockingDeque) Iterator() Iterator {
	return newIterator(q, false)
}

// DescendingIterator return a desc iterator of this deque
func (q *LinkedBlockingDeque) DescendingIterator() Iterator {
	return newIterator(q, true)
}

func newIterator(q *LinkedBlockingDeque, desc bool) *LinkedBlockingDequeIterator {
	q.lock.Lock()
	defer q.lock.Unlock()
	iterator := LinkedBlockingDequeIterator{q: q, desc: desc}
	iterator.next = iterator.firstNode()
	if iterator.next == nil {
		iterator.nextItem = nil
	} else {
		iterator.nextItem = iterator.next.item
	}
	return &iterator
}

// LinkedBlockingDequeIterator is iterator implements for LinkedBlockingDeque
type LinkedBlockingDequeIterator struct {
	q        *LinkedBlockingDeque
	next     *Node
	nextItem interface{}
	lastRet  *Node
	desc     bool
}

func (iterator *LinkedBlockingDequeIterator) firstNode() *Node {
	if iterator.desc {
		return iterator.q.last
	}
	return iterator.q.first
}

func (iterator *LinkedBlockingDequeIterator) nextNode(node *Node) *Node {
	if iterator.desc {
		return node.prev
	}
	return node.next
}

// HasNext return is exist next element
func (iterator *LinkedBlockingDequeIterator) HasNext() bool {
	return iterator.next != nil
}

// Next return next element, if not exist will return nil
func (iterator *LinkedBlockingDequeIterator) Next() interface{} {
	if iterator.next == nil {
		//TODO error or nil ?
		//panic(errors.New("NoSuchElement"))
		return nil
	}
	iterator.lastRet = iterator.next
	x := iterator.nextItem
	iterator.advance()
	return x
}

func (iterator *LinkedBlockingDequeIterator) advance() {
	lock := iterator.q.lock
	lock.Lock()
	defer lock.Unlock()
	iterator.next = iterator.succ(iterator.next)
	if iterator.next == nil {
		iterator.nextItem = nil
	} else {
		iterator.nextItem = iterator.next.item
	}
}

func (iterator *LinkedBlockingDequeIterator) succ(n *Node) *Node {
	for {
		s := iterator.nextNode(n)
		if s == nil {
			return nil
		} else if s.item != nil {
			return s
		} else if s == n {
			return iterator.firstNode()
		}
		n = s
	}
}

// Remove current element from dequeue
func (iterator *LinkedBlockingDequeIterator) Remove() {
	n := iterator.lastRet
	if n == nil {
		panic(errors.New("IllegalStateException"))
	}
	iterator.lastRet = nil
	lock := iterator.q.lock
	lock.Lock()
	if n.item != nil {
		iterator.q.unlink(n)
	}
	lock.Unlock()
}
