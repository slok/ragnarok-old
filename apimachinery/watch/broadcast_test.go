package watch_test

import (
	"testing"
	"time"

	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/log"
	mwatch "github.com/slok/ragnarok/mocks/apimachinery/watch"
	"github.com/stretchr/testify/assert"
)

func TestBroadcastWatcherSendEvent(t *testing.T) {
	tests := []struct {
		name     string
		expEvent watch.Event
	}{
		{
			name:     "Sending and event over channel should be sent over the channel.",
			expEvent: watch.Event{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// mocks.
			mm := &mwatch.Multiplexer{}

			// Buffered channel so we don't block waiting for the result at the end of the test.
			ec := make(chan watch.Event, 1)

			// Create the watcher.
			w := watch.NewBroadcastWatcher("id1", ec, mm, log.Dummy)

			// Send the event.
			go func() {
				w.SendEvent(test.expEvent)
			}()
			select {
			case <-time.After(5 * time.Millisecond):
				assert.Fail("timeout receiving result")
			case gotEvent := <-ec:
				assert.Equal(test.expEvent, gotEvent)
			}
		})
	}
}

func TestBroadcastWatcherStop(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{
			name: "Sending and event over channel should be sent over the channel.",
			id:   "test1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// mocks.
			mm := &mwatch.Multiplexer{}
			mm.On("StopWatcher", test.id).Once()

			// Buffered channel so we don't block waiting for the result at the end of the test.
			ec := make(chan watch.Event, 1)

			// Create the watcher.
			w := watch.NewBroadcastWatcher(test.id, ec, mm, log.Dummy)

			// Send the event.
			w.Stop()
			mm.AssertExpectations(t)
		})
	}
}
