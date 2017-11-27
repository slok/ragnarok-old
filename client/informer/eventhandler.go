package informer

import (
	"github.com/slok/ragnarok/api"
)

// ResourceEventHandler knows how to handle all the resource events
type ResourceEventHandler interface {
	// OnAdd processes an added object resource.
	OnAdd(obj api.Object)
	// OnUpdate processes an updated object resource.
	OnUpdate(oldObj, newObj api.Object)
	// OnDelete processes an deleted object resource.
	OnDelete(obj api.Object)
}

// ResourceEventHandlerFuncs implements ResourceEventHandler using functions, very helpful so we don't need
// to decaler new types to implement a handler.
type ResourceEventHandlerFuncs struct {
	OnAddFunc    func(obj api.Object)
	OnUpdateFunc func(oldObj, newObj api.Object)
	OnDeleteFunc func(obj api.Object)
}

// OnAdd satisfies ResourceEventHandler.
func (r *ResourceEventHandlerFuncs) OnAdd(obj api.Object) {
	if r.OnAddFunc != nil {
		r.OnAddFunc(obj)
	}
}

// OnUpdate satisfies ResourceEventHandler.
func (r *ResourceEventHandlerFuncs) OnUpdate(oldObj, newObj api.Object) {
	if r.OnUpdateFunc != nil {
		r.OnUpdateFunc(oldObj, newObj)
	}
}

// OnDelete satisfies ResourceEventHandler.
func (r *ResourceEventHandlerFuncs) OnDelete(obj api.Object) {
	if r.OnDeleteFunc != nil {
		r.OnDeleteFunc(obj)
	}
}
