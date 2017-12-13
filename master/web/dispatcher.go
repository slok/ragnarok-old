package web

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/web/handler"
)

const (
	labelSelectorQueryParam = "labelSelector"
	labelSelectorSeparator  = ","
	labelSelectorEqual      = "="
	watchQueryParam         = "watch"
)

// ResourceHandlerDispatcherInterface is the dispatcher interface that will be implemented in order
// to dispatch requests to the required handlers of the resources.
type ResourceHandlerDispatcherInterface interface {
	// Dispathc will dispatch to the required
	Dispatch(w http.ResponseWriter, r *http.Request)
}

// ResourceHandlerDispatcher will dispatch the request to the correct resource handlers.
type ResourceHandlerDispatcher struct {
	handler handler.ResourceHandler
	logger  log.Logger
}

// NewResourceHandlerDispatcher returns a new ResourceHandlerDispatcher.
func NewResourceHandlerDispatcher(h handler.ResourceHandler) *ResourceHandlerDispatcher {
	return &ResourceHandlerDispatcher{
		handler: h,
	}
}

func (d *ResourceHandlerDispatcher) getIDFromURL(url *url.URL) string {
	id := strings.TrimLeft(url.Path, d.handler.GetRoute())
	id = strings.TrimRight(id, "/")
	return id
}

func (d *ResourceHandlerDispatcher) getOptsFromURL(url *url.URL) map[string]string {
	// TODO Move all the query selector stuff to api machinery.
	res := map[string]string{}
	v := url.Query().Get(labelSelectorQueryParam)
	if v == "" {
		return res
	}

	labels := strings.Split(v, labelSelectorSeparator)
	for _, label := range labels {
		spl := strings.Split(label, labelSelectorEqual)
		if len(spl) != 2 {
			d.logger.Warnf("ignoring wrong list option label: %s", label)
		} else {
			res[spl[0]] = spl[1]
		}
	}
	return res
}
func (d *ResourceHandlerDispatcher) isWatcherFromURL(url *url.URL) bool {
	v := url.Query().Get(watchQueryParam)
	return strings.ToLower(v) == "true"
}

// Dispatch will dispatch the requests to the correct resource handler.
func (d *ResourceHandlerDispatcher) Dispatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// a post always means a creation.
		d.handler.Create(w, r)
	case "GET":
		// Get and ID means Get otherwise a List or a watch.
		if id := d.getIDFromURL(r.URL); id != "" {
			d.handler.Get(w, r, id)
			// if watch parameter then watch action otherw.
		} else if d.isWatcherFromURL(r.URL) {
			opts := d.getOptsFromURL(r.URL)
			d.handler.Watch(w, r, opts)
		} else {
			opts := d.getOptsFromURL(r.URL)
			d.handler.List(w, r, opts)
		}
	case "PUT":
		id := d.getIDFromURL(r.URL)
		d.handler.Update(w, r, id)
	case "DELETE":
		id := d.getIDFromURL(r.URL)
		d.handler.Delete(w, r, id)
	}
}
