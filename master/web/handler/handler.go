package handler

import "net/http"

// ResourceHandler is the handler that every resource needs to implement will
// serve the server.
type ResourceHandler interface {
	// Ceate will create a new resource.
	Create(w http.ResponseWriter, r *http.Request)
	// Update will update an existing resource.
	Update(w http.ResponseWriter, r *http.Request, id string)
	// Delete will delete a resource.
	Delete(w http.ResponseWriter, r *http.Request, id string)
	// Get will retrieve an existing resource.
	Get(w http.ResponseWriter, r *http.Request, id string)
	// List will list resources based on options.
	List(w http.ResponseWriter, r *http.Request, opts map[string]string)
	// Watch will wait for watch events based on options.
	Watch(w http.ResponseWriter, r *http.Request, opts map[string]string)
	// TODO: Patch.

	// GetRoute returns the route where this resource will be handling requests.
	GetRoute() string
}
