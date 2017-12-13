package informer

import (
	"fmt"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/client/repository"
)

// ListerWatcher is an interface that knows hot to list adn watch objects.
type ListerWatcher interface {
	List(opts api.ListOptions) ([]api.Object, error)
	Watch(opts api.ListOptions) (watch.Watcher, error)
}

// ListerWatcherFuncs implements ListerWatcher using functions, very helpful so we don't need
// to decaler new types to implement a lister watcher.
type ListerWatcherFuncs struct {
	ListFunc  func(opts api.ListOptions) ([]api.Object, error)
	WatchFunc func(opts api.ListOptions) (watch.Watcher, error)
}

// List will list resources from a repository. List satisfies ListerWatcher interface.
func (l *ListerWatcherFuncs) List(opts api.ListOptions) ([]api.Object, error) {
	if l.ListFunc != nil {
		return l.ListFunc(opts)
	}
	return nil, fmt.Errorf("list function can't be nil")
}

// Watch satisfies ListerWatcher interface.
func (l *ListerWatcherFuncs) Watch(opts api.ListOptions) (watch.Watcher, error) {
	if l.WatchFunc != nil {
		return l.WatchFunc(opts)
	}
	return nil, fmt.Errorf("watch function can't be nil")
}

// RepositoryListerWatcher will list and watch resources using a respoitory client interface.
type RepositoryListerWatcher struct {
	c repository.Client
}

// NewRepositoryListerWatcher returns a new RepositoryListerWatcher.
func NewRepositoryListerWatcher(client repository.Client) *RepositoryListerWatcher {
	return &RepositoryListerWatcher{
		c: client,
	}
}

// List will list resources from a repository. Satisfies ListerWatcher interface.
func (m *RepositoryListerWatcher) List(opts api.ListOptions) ([]api.Object, error) {
	return m.c.List(opts)
}

// Watch will watch resources from repository. Satisfies ListerWatcher interface.
func (m *RepositoryListerWatcher) Watch(opts api.ListOptions) (watch.Watcher, error) {
	return m.c.Watch(opts)
}
