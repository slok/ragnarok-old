package queue_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/slok/ragnarok/client/util/queue"
	"github.com/stretchr/testify/assert"
)

func TestSimpleQueueBlocking(t *testing.T) {
	tests := []struct {
		name string
		item string
	}{
		{
			name: "Getting an item on an empty queue should block the caller until a push operation is made on that queue.",
			item: "this is a blocking test.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			q := queue.NewSimpleQueue()

			// Caller1 pop in a goroutine.
			c1 := make(chan interface{})
			go func() {
				res, _ := q.Pop()
				c1 <- res
			}()

			// Caller2 pop in a goroutine.
			c2 := make(chan interface{})
			go func() {
				res, _ := q.Pop()
				c2 <- res
			}()

			// Wait until we have items on the queue.
			select {
			case <-time.After(10 * time.Millisecond):
			case <-c1:
				assert.Fail("pop should get blocked, it didn't")
			case <-c2:
				assert.Fail("pop should get blocked, it didn't")
			}

			// Add an item & wait until we have items on the queue.
			var gotItem1, gotItem2 interface{}
			err := q.Push(test.item)
			assert.NoError(err)
			select {
			case <-time.After(10 * time.Millisecond):
				assert.Fail("pop shouldn't get blocked, we had push an item")
			case gotItem1 = <-c1:
			case gotItem2 = <-c2:
			}

			// Add a 2nd item & wait until we have items on the queue.
			err = q.Push(test.item)
			assert.NoError(err)

			select {
			case <-time.After(10 * time.Millisecond):
				assert.Fail("pop shouldn't get blocked, we had push an item")
			case gotItem1 = <-c1:
			case gotItem2 = <-c2:
			}

			// At this moment both should have the items.
			assert.Equal(test.item, gotItem1.(string))
			assert.Equal(test.item, gotItem2.(string))
		})
	}
}

func TestSimpleQueueShutDown(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Shutting down the queue when there are blocking pops should unblock everyon.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			q := queue.NewSimpleQueue()

			// Caller1 pop in a goroutine.
			c1 := make(chan bool)
			go func() {
				_, sd := q.Pop()
				c1 <- sd
			}()

			// Caller2 pop in a goroutine.
			c2 := make(chan bool)
			go func() {
				_, sd := q.Pop()
				c2 <- sd
			}()

			// Wait until we have items on the queue.
			select {
			case <-time.After(10 * time.Millisecond):
			case <-c1:
				assert.Fail("pop should get blocked, it didn't")
			case <-c2:
				assert.Fail("pop should get blocked, it didn't")
			}

			// Shut down the queue.
			err := q.ShutDown()
			assert.NoError(err)
			assert.True(q.IsShutDown())

			// Check if had a result
			select {
			case <-time.After(10 * time.Millisecond):
				assert.Fail("shut down should had unblock the pops, it didn't")
			case shutDown := <-c1:
				assert.True(shutDown) // check if unlock sent us shutDown.
			}

			select {
			case <-time.After(10 * time.Millisecond):
				assert.Fail("shut down should had unblock the pops, it didn't")
			case shutDown := <-c2:
				assert.True(shutDown) // check if unlock sent us shutDown.
			}

			// Pushing on a shut down queue should return error.
			assert.Error(q.Push("test"))

			// Shutting down an already shut donw queue should error.
			assert.Error(q.ShutDown())

			// Pop on a shutdown queue should return shut down flag.
			_, sd := q.Pop()
			assert.True(sd)
		})
	}
}

func TestSimpleQueueMultiplePushAndPop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Pushing and poping multiple times should not block the pops and get the items in order.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			times := 100
			q := queue.NewSimpleQueue()
			defer q.ShutDown() // Clean.

			// Items to queue.
			for i := 0; i < times; i++ {
				item := fmt.Sprintf("this is a test %d", i)
				assert.NoError(q.Push(item))
			}
			fmt.Println(q.Len())
			c := make(chan interface{})
			pop := func() {
				res, _ := q.Pop()
				c <- res
			}

			// Dequeue items.
			for i := 0; i < times; i++ {
				go pop()
				// Check if had a result.
				select {
				case <-time.After(10 * time.Millisecond):
					assert.Fail("shut down should had unblock the pops, it didn't")
				case res := <-c:
					expItem := fmt.Sprintf("this is a test %d", i)
					assert.Equal(expItem, res)
				}
			}

			// No items should block.
			go pop()
			select {
			case <-time.After(10 * time.Millisecond):
			case <-c:
				assert.Fail("pop should get blocked, it didn't")
			}
		})
	}
}
