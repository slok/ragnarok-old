package watch

import (
	"fmt"
	"sync"
	"time"

	"github.com/slok/ragnarok/log"
)

const sendEventTimeout = 10 * time.Millisecond

// broadcastWatcher is a simple watcher for the Broadcaster multiplexer.
type broadcastWatcher struct {
	id     string
	eventC chan Event
	logger log.Logger
	muxer  Multiplexer
	stop   sync.Once
}

// newBroadcastWatcher returns a new BroadcastWatcher.
func newBroadcastWatcher(id string, eventC chan Event, muxer Multiplexer, logger log.Logger) *broadcastWatcher {
	return &broadcastWatcher{
		id:     id,
		eventC: eventC,
		muxer:  muxer,
		logger: logger,
	}
}

// Stop satisfies Watcher interface.
func (b *broadcastWatcher) Stop() {
	b.stop.Do(func() {
		b.muxer.StopWatcher(b.id)
	})
	b.logger.Infof("watcher stopped")
}

// GetChan satisfies Watcher interface.
func (b *broadcastWatcher) GetChan() <-chan Event {
	return b.eventC
}

// Broadcaster will multiplex events among all the watchers registered.
type Broadcaster struct {
	lock     sync.Mutex
	watchers map[string]*broadcastWatcher
	logger   log.Logger
}

// NewBroadcaster returns a new broadcaster.
func NewBroadcaster(logger log.Logger) *Broadcaster {
	logger = logger.WithField("multiplexer", "broadcaster")
	return &Broadcaster{
		watchers: map[string]*broadcastWatcher{},
		logger:   logger,
	}
}

// SendEvent will send events to all the registered watchers. Satisfies multiplexer interface.
func (b *Broadcaster) SendEvent(e Event) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, w := range b.watchers {
		e := e // Make a copy of each event so each watcher has its own event.
		select {
		// Don't block if event can't be sent.
		case <-time.After(sendEventTimeout):
			b.logger.Warnf("timeout sending event to %s watcher", w.id)
		case w.eventC <- e:
		}
	}
}

// StartWatcher will create a new watchers and register. Satisfies multiplexer interface.
func (b *Broadcaster) StartWatcher() (Watcher, error) {
	// Create the watcher.
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	c := make(chan Event)
	w := newBroadcastWatcher(id, c, b, b.logger)

	// Register our watcher.
	b.lock.Lock()
	b.watchers[id] = w
	b.lock.Unlock()

	b.logger.Infof("watcher %s registered", id)
	return w, nil
}

// StopWatcher will stop the watcher an unregister. Satisfies multiplexer interface.
func (b *Broadcaster) StopWatcher(id string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	w, ok := b.watchers[id]
	if !ok {
		return
	}
	close(w.eventC)
	delete(b.watchers, id)
	b.logger.Infof("watcher %s unregistered", id)
}

// StopAll will stop and unregister all the watchers. Satisfies multiplexer interface.
func (b *Broadcaster) StopAll() {
	for _, w := range b.watchers {
		w.Stop()
	}
}
