package watch_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/log"
	testapi "github.com/slok/ragnarok/test/api"
)

func TestBroadcasSendEventOnWatchers(t *testing.T) {
	tests := []struct {
		name     string
		expEvent watch.Event
	}{
		{
			name: "Starting a watcher and sending an event should be received by multiple watchers.",
			expEvent: watch.Event{
				Type: watch.AddedEvent,
				Object: &testapi.TestObj{
					Labels: map[string]string{"test-event": "test1"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			numberWatchers := 5
			watchers := make([]watch.Watcher, numberWatchers)
			gotEvents := make([]watch.Event, numberWatchers)

			// Create the broadcaster and the watchers.
			b := watch.NewBroadcaster(log.Dummy)
			for i := 0; i < numberWatchers; i++ {
				w, err := b.StartWatcher()
				require.NoError(err)
				require.NotNil(w)
				watchers[i] = w
			}

			var wg sync.WaitGroup
			wg.Add(numberWatchers)
			// Start getting events from the watchers.
			for i := 0; i < numberWatchers; i++ {
				i := i
				go func() {
					defer wg.Done()
					c := watchers[i].GetChan()
					select {
					case <-time.After(10 * time.Millisecond): // If timeout don't add to the got events.
					case ev := <-c:
						gotEvents[i] = ev
					}
				}()
			}

			// Add an event.
			b.SendEvent(test.expEvent)

			// Wait until all events consumed
			wg.Wait()
			// Check every watcher has received the event.
			for i := 0; i < numberWatchers; i++ {
				assert.Equal(test.expEvent, gotEvents[i])
			}
		})
	}
}

func TestBroadcasStopAllWatchers(t *testing.T) {
	tests := []struct {
		name     string
		expEvent watch.Event
	}{
		{
			name: "Stopping all watchers Should stop all the registered watchers.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			numberWatchers := 5
			watchers := make([]watch.Watcher, numberWatchers)

			// Create the broadcaster and the watchers.
			b := watch.NewBroadcaster(log.Dummy)
			for i := 0; i < numberWatchers; i++ {
				w, err := b.StartWatcher()
				require.NoError(err)
				require.NotNil(w)
				watchers[i] = w
			}

			// Stop all watchers.
			b.StopAll()

			// Check it has unregistered all the watchers (should not panic sending an event).
			assert.NotPanics(func() {
				ev := watch.Event{
					Type: watch.AddedEvent,
					Object: &testapi.TestObj{
						Labels: map[string]string{"test-event": "test1"},
					},
				}
				b.SendEvent(ev)
			})

			// Check all the channels are closed.
			for _, w := range watchers {
				watcher := w.GetChan()
				select {
				case <-time.After(10 * time.Millisecond):
					assert.Fail("not closed channel")
				case <-watcher: // Closed channel is instant, so if it's closed should enter here (and should be closed).
				}
			}
		})
	}
}
