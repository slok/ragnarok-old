package watch

import (
	"sync"

	"github.com/slok/ragnarok/log"
)

// BroadcastWatcher is a simple watcher for the Broadcaster multiplexer.
type BroadcastWatcher struct {
	id     string
	eventC chan Event
	logger log.Logger
	muxer  Multiplexer
	stop   sync.Once
}

// NewBroadcastWatcher returns a new BroadcastWatcher.
func NewBroadcastWatcher(id string, eventC chan Event, muxer Multiplexer, logger log.Logger) *BroadcastWatcher {
	return &BroadcastWatcher{
		id:     id,
		eventC: eventC,
		muxer:  muxer,
		logger: logger,
	}
}

// Stop satisfies Watcher interface.
func (b *BroadcastWatcher) Stop() {
	b.stop.Do(func() {
		b.muxer.StopWatcher(b.id)
	})
	b.logger.Infof("watcher stopped")
}

// GetChan satisfies Watcher interface.
func (b *BroadcastWatcher) GetChan() <-chan Event {
	return b.eventC
}

// SendEvent will send an event over the channel. Blocking operation.
func (b *BroadcastWatcher) SendEvent(e Event) {
	b.eventC <- e
}

// Broadcaster will multiplex events among all the watchers registered.
type Broadcaster struct{}

// SendEvent will send events to all the registered watchers. Satisfies multiplexer interface.
func (b *Broadcaster) SendEvent(Event) {
}

// StartWatcher will create a new watchers and register. Satisfies multiplexer interface.
func (b *Broadcaster) StartWatcher() (Watcher, error) {
	return nil, nil
}

// StopWatcher will stop the watcher an unregister. Satisfies multiplexer interface.
func (b *Broadcaster) StopWatcher(id string) {
}

// StopAll will stop and unregister all the watchers. Satisfies multiplexer interface.
func (b *Broadcaster) StopAll() {
}
