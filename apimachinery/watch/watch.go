package watch

import (
	"github.com/slok/ragnarok/api"
)

// EventType is the type of an event.
type EventType int

const (
	// AddedEvent is an add event.
	AddedEvent EventType = iota
	// ModifiedEvent is a modify event.
	ModifiedEvent
	// DeletedEvent is a delete event.
	DeletedEvent
	// ErrorEvent is an error event.
	ErrorEvent
)

// Event represents a single event to a watched resource.
type Event struct {
	// Type is the type of the event.
	Type EventType
	// Object is the object.
	Object api.Object
}

// Watcher will be implemented by any one that wants to expose events on object.
type Watcher interface {
	// Stop will close the result chanel.
	Stop()
	// GetChan will return the channel that will notify the events
	GetChan() <-chan Event
}

// Multiplexer will multiplex the received events into multiple watchers.
type Multiplexer interface {
	// SendEvent will send an event on the to the desired watchers.
	SendEvent(Event)
	// StartWatcher will cretae new a watcher.
	StartWatcher() (Watcher, error)
	// CloseWatcher will stop a new watcher.
	StopWatcher(string)
	// StopAll will stop all the watchers.
	StopAll()
}
