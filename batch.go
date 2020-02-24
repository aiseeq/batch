// Package batch provides simple and convenient methods for delayed batch processing
package batch

import (
	"time"
)

type Batch struct {
	queue chan interface{}
	flush func(chan interface{})
}

// Add inserts new element into batch. Returns true on success or false if queue is overflowed
func (b *Batch) Add(row interface{}) bool {
	select {
	case b.queue <- row:
		// message sent
		return true
	default:
		// message dropped
		return false
	}
}

// Wait can be used to wait until all elements in batch will be processed
// For example on application shutdown when Add() isn't called anymore
func (b *Batch) Wait() {
	for len(b.queue) > 1 {
		time.Sleep(time.Millisecond)
	}
}

// New returns Batch object configured to call flush function with flushPeriod
func New(flush func(chan interface{}), queueLen int, flushPeriod time.Duration) *Batch {
	b := Batch{
		queue: make(chan interface{}, queueLen),
		flush: flush,
	}

	go func() {
		for {
			time.Sleep(flushPeriod)
			b.flush(b.queue)
		}
	}()

	return &b
}
