package queue

import (
	"fmt"
	"sync"
)

// Queue is a simple queue that will be used to get the work of resources that needs to
// be processed .
type Queue interface {
	// Pop will retrieve the first object in the queue if there are no objects it will
	// get blocked until there are more items to process or the work queue is shut down.
	// The queue will return also if the queue has been shut down.
	Pop() (v interface{}, shutDown bool)
	// Push will push an object in the queue. Will return an error if the queue is shut down.
	Push(interface{}) error
	// Len will get the length of the queue.
	Len() int
	// ShutDown will stop and close the queue.
	ShutDown() error
	// IsShutDown will return true if the queue has been shut down.
	IsShutDown() bool
}

// SimpleQueue is a simple and regular queue that is concurrently safe.
type SimpleQueue struct {
	q        []interface{}
	cond     *sync.Cond // cond will be used to block goroutines on the pop when there are no items and unblock goroutine by goroutine.
	lock     sync.Mutex // lock is the global object lock, to access q and shutdown attrs.
	shutDown bool
}

// NewSimpleQueue returns a new simple queue object.
func NewSimpleQueue() *SimpleQueue {
	return &SimpleQueue{
		q:    []interface{}{},
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// Pop satisfies Queue interface.
func (s *SimpleQueue) Pop() (interface{}, bool) {
	if s.IsShutDown() {
		return nil, true
	}

	// Block until there is an item in te queue.
	len := s.Len()
	if len == 0 {
		s.cond.L.Lock()
		defer s.cond.L.Unlock()
		s.cond.Wait()
	}

	// Check if we have been waiting and being unblocked after a shut down.
	if s.IsShutDown() {
		return nil, true
	}

	// Get the item
	s.lock.Lock()
	item, newQ := s.q[0], s.q[1:]
	s.q = newQ
	s.lock.Unlock()

	return item, false
}

// Push satisfies Queue interface.
func (s *SimpleQueue) Push(v interface{}) error {
	if s.IsShutDown() {
		return fmt.Errorf("can't push items on a shut down queue")
	}
	s.lock.Lock()
	s.q = append(s.q, v)
	s.lock.Unlock()

	// Send signal and wake up one pop waiter if there is any.
	s.cond.Signal()

	return nil
}

// Len satisfies Queue interface.
func (s *SimpleQueue) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.q)
}

// ShutDown satisfies Queue interface.
func (s *SimpleQueue) ShutDown() error {
	if s.IsShutDown() {
		return fmt.Errorf("queue already shut down")
	}
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	// Unblock all blocked pops.
	s.cond.Broadcast()
	s.shutDown = true

	return nil
}

// IsShutDown satisfies Queue interface.
func (s *SimpleQueue) IsShutDown() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.shutDown
}

// TODO Rate limit queue, so we can requeue multiple times.
