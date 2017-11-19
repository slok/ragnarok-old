package repository

import (
	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/apimachinery/watch"
)

// Client knows how to "store" & "retrieve" api objects to repositories, for example:
// Rest, Grpc, memory, Redis repository... In other words an store adapter.
type Client interface {
	Create(obj api.Object) (api.Object, error)
	Update(obj api.Object) (api.Object, error)
	Delete(id string) error
	Get(id string) (api.Object, error)
	List(opts api.ListOptions) ([]api.Object, error)
	Watch(opts api.ListOptions) (watch.Watch, error)
	// TODO Patch
}
