package watch

import (
	"fmt"
	"sync"
	"time"

	"github.com/slok/ragnarok/log"
)

const sendEventTimeout = 10 * time.Millisecond

// broadcastdWatcher is a simple watcher for the multiplexer with a filter that
// the multiplexer should call before adding a new event.
type broadcastWatcher struct {
	id     string
	filter ObjectFilter
	eventC chan Event
	logger log.Logger
	muxer  Multiplexer
	stop   sync.Once
}

// newFilteredWatcher returns a new broadcastWatcher.
func newFilteredWatcher(id string, filter ObjectFilter, eventC chan Event, muxer Multiplexer, logger log.Logger) *broadcastWatcher {
	return &broadcastWatcher{
		id:     id,
		filter: filter,
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
		// Check if the event based on the filter applied to an object.
		if w.filter.Filter(e.Object) {
			continue
		}

		e := e.DeepCopy() // Make a copy of each event so each watcher has its own event and don't share the objects.
		select {
		// Don't block if event can't be sent.
		case <-time.After(sendEventTimeout):
			b.logger.Warnf("timeout sending event to %s watcher", w.id)
		case w.eventC <- e:
		}
	}
}

// StartWatcher will create a new watchers and register. Satisfies multiplexer interface.
func (b *Broadcaster) StartWatcher(f ObjectFilter) (Watcher, error) {
	// Create the watcher.
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	c := make(chan Event)
	w := newFilteredWatcher(id, f, c, b, b.logger)

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

// BroadcasterFactory is the default broadcaster creator factory based on the object types. This
// is used in order to reuse the broadcasters of events.
type BroadcasterFactory struct {
	registry map[string]*Broadcaster
	lock     sync.Mutex
	logger   log.Logger
}

// NewDefaultBroadcasterFactory creates a new BroadcasterFactory instance with an empty registry.
func NewDefaultBroadcasterFactory(logger log.Logger) *BroadcasterFactory {
	return NewBroadcasterFactory(map[string]*Broadcaster{}, logger)
}

// NewBroadcasterFactory creates a new BroadcasterFactory instance.
func NewBroadcasterFactory(registry map[string]*Broadcaster, logger log.Logger) *BroadcasterFactory {
	return &BroadcasterFactory{
		registry: registry,
		logger:   logger,
	}
}

// Get returns a new broadcaster based on the obj type. If there is already present it returns the
// previously created. Satisfies MultiplexerFactory interface.
func (b *BroadcasterFactory) Get(id string) Multiplexer {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.registry[id]; !ok {
		b.registry[id] = NewBroadcaster(b.logger)
	}

	return b.registry[id]
}
